import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { elementsApi, uploadsApi } from '@/api'
import { useBoardsStore } from './boards'
import type { Element, CanvasState, ElementKind, ElementPayload } from '@/types'

export const useEditorStore = defineStore('editor', () => {
  // State
  const selectedElementIds = ref<string[]>([])
  const undoStack = ref<CanvasState[]>([])
  const redoStack = ref<CanvasState[]>([])
  const snapToGrid = ref(true)
  const canvas = ref<HTMLElement | null>(null)
  const currentPageId = ref<string | null>(null)
  const elements = ref<Element[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  
  // Auto-save state
  const saveStatus = ref<'idle' | 'saving' | 'saved' | 'error'>('idle')
  const lastSaveTime = ref<Date | null>(null)
  const saveError = ref<string | null>(null)
  const pendingSaves = ref<Set<string>>(new Set())
  const retryCount = ref<Map<string, number>>(new Map())

  // Constants
  const MAX_HISTORY_SIZE = 50
  const DEBOUNCE_DELAY = 300 // Default debounce
  const POSITION_DEBOUNCE_DELAY = 1000 // Longer delay for position/size changes
  const MAX_RETRY_ATTEMPTS = 3
  const RETRY_DELAY_BASE = 1000 // Base delay for retries (exponential backoff)
  const SAVE_SUCCESS_DISPLAY_TIME = 2000 // How long to show "saved" status

  // Getters
  const selectedElements = computed(() => 
    elements.value.filter(el => selectedElementIds.value.includes(el.id))
  )
  
  const canUndo = computed(() => undoStack.value.length > 0)
  const canRedo = computed(() => redoStack.value.length > 0)
  
  const sortedElements = computed(() => {
    return [...elements.value].sort((a, b) => a.z - b.z)
  })

  // Store for debounced saves
  let saveTimeout: NodeJS.Timeout | null = null
  const positionSaveTimeouts: Map<string, NodeJS.Timeout> = new Map()
  
  const debouncedSave = (element: Element, isPositionUpdate = false) => {
    const delay = isPositionUpdate ? POSITION_DEBOUNCE_DELAY : DEBOUNCE_DELAY
    
    if (isPositionUpdate) {
      // Handle position updates per element with longer delay
      const elementId = element.id
      const existingTimeout = positionSaveTimeouts.get(elementId)
      
      if (existingTimeout) {
        clearTimeout(existingTimeout)
      }
      
      const timeout = setTimeout(async () => {
        console.log('Saving position update for element:', elementId)
        await saveElement(element)
        positionSaveTimeouts.delete(elementId)
      }, delay)
      
      positionSaveTimeouts.set(elementId, timeout)
    } else {
      // Handle immediate updates (text changes, style changes, etc.)
      if (saveTimeout) {
        clearTimeout(saveTimeout)
      }
      
      saveTimeout = setTimeout(async () => {
        console.log('Saving immediate update for element:', element.id)
        await saveElement(element)
      }, delay)
    }
  }

  // Actions
  const setCanvas = (canvasElement: HTMLElement | null) => {
    canvas.value = canvasElement
  }

  const setCurrentPage = (pageId: string) => {
    currentPageId.value = pageId
    elements.value = []
    selectedElementIds.value = []
  }

  const loadPageElements = async (pageId: string) => {
    const boardsStore = useBoardsStore()
    console.log('Loading elements for page:', pageId)
    console.log('Available pages:', boardsStore.pages.length)
    
    const page = boardsStore.pages.find(p => p.id === pageId)
    console.log('Found page:', page?.title)
    console.log('Page elements:', page?.elements?.length || 0)
    
    if (page?.elements) {
      // Ensure all elements have proper default values for visible and locked
      elements.value = page.elements.map(element => ({
        ...element,
        visible: element.visible ?? true,  // Default to true if not set
        locked: element.locked ?? false   // Default to false if not set
      }))
      console.log('Elements loaded and normalized:', elements.value.length)
      
      if (elements.value.length > 0) {
        console.log('Sample element after normalization:', JSON.stringify(elements.value[0], null, 2))
      }
    } else {
      elements.value = []
      console.log('No elements found, setting empty array')
    }
  }

  const selectElements = (elementIds: string[]) => {
    selectedElementIds.value = elementIds
  }

  const selectElement = (elementId: string) => {
    selectedElementIds.value = [elementId]
  }

  const toggleElementSelection = (elementId: string) => {
    const index = selectedElementIds.value.indexOf(elementId)
    if (index > -1) {
      selectedElementIds.value.splice(index, 1)
    } else {
      selectedElementIds.value.push(elementId)
    }
  }

  const clearSelection = () => {
    selectedElementIds.value = []
  }

  const saveCanvasState = () => {
    const state: CanvasState = {
      elements: JSON.parse(JSON.stringify(elements.value)),
      timestamp: Date.now(),
    }
    
    undoStack.value.push(state)
    
    // Limit history size
    if (undoStack.value.length > MAX_HISTORY_SIZE) {
      undoStack.value.shift()
    }
    
    // Clear redo stack when new action is performed
    redoStack.value = []
  }

  const undo = () => {
    if (undoStack.value.length === 0) return
    
    // Save current state to redo stack
    const currentState: CanvasState = {
      elements: JSON.parse(JSON.stringify(elements.value)),
      timestamp: Date.now(),
    }
    redoStack.value.push(currentState)
    
    // Restore previous state
    const previousState = undoStack.value.pop()!
    elements.value = previousState.elements
    
    // Trigger canvas reload if available
    if (canvas.value) {
      // Emit event to trigger canvas reload
      const event = new CustomEvent('editor:reload-elements')
      window.dispatchEvent(event)
    }
  }

  const redo = () => {
    if (redoStack.value.length === 0) return
    
    // Save current state to undo stack
    const currentState: CanvasState = {
      elements: JSON.parse(JSON.stringify(elements.value)),
      timestamp: Date.now(),
    }
    undoStack.value.push(currentState)
    
    // Restore next state
    const nextState = redoStack.value.pop()!
    elements.value = nextState.elements
    
    // Trigger canvas reload if available
    if (canvas.value) {
      // Emit event to trigger canvas reload
      const event = new CustomEvent('editor:reload-elements')
      window.dispatchEvent(event)
    }
  }

  const createElement = async (
    kind: ElementKind,
    x: number,
    y: number,
    w: number,
    h: number,
    payload: ElementPayload,
    rotation = 0
  ) => {
    const boardsStore = useBoardsStore()
    
    if (!boardsStore.currentBoard || !currentPageId.value || !boardsStore.editToken) {
      throw new Error('No board loaded or edit access required')
    }

    loading.value = true
    error.value = null

    try {
      // Save current state for undo
      saveCanvasState()
      
      const newElement = await elementsApi.create(
        boardsStore.currentBoard.id,
        currentPageId.value,
        { 
          kind, 
          x, 
          y, 
          w, 
          h, 
          rotation, 
          visible: true,  // Default to visible
          locked: false,  // Default to unlocked
          payload 
        }
      )
      
      // Only add to local state if server creation was successful
      elements.value.push(newElement)
      console.log('Element created successfully:', newElement.id)
      return newElement
    } catch (err: any) {
      console.error('Failed to create element:', err)
      error.value = err.error?.message || 'Failed to create element'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateElement = async (elementId: string, updates: Partial<Element>) => {
    console.log('updateElement called:', elementId, updates)
    
    const boardsStore = useBoardsStore()
    
    if (!boardsStore.currentBoard || !currentPageId.value || !boardsStore.editToken) {
      throw new Error('No board loaded or edit access required')
    }

    const element = elements.value.find(el => el.id === elementId)
    if (!element) {
      throw new Error('Element not found')
    }

    // Update local state immediately for responsiveness
    const index = elements.value.findIndex(el => el.id === elementId)
    elements.value[index] = { ...element, ...updates }

    // Check if this is a position/size update
    const isPositionUpdate = 'x' in updates || 'y' in updates || 'w' in updates || 'h' in updates || 'rotation' in updates
    
    // Debounce API call with appropriate delay
    debouncedSave(elements.value[index], isPositionUpdate)
  }

  const saveElement = async (element: Element): Promise<void> => {
    const boardsStore = useBoardsStore()
    
    if (!boardsStore.currentBoard || !currentPageId.value || !boardsStore.editToken) {
      return
    }

    const elementId = element.id
    
    try {
      // Mark element as being saved
      pendingSaves.value.add(elementId)
      saveStatus.value = 'saving'
      saveError.value = null

      const updateData: any = {
        x: element.x,
        y: element.y,
        w: element.w,
        h: element.h,
        rotation: element.rotation,
        z: element.z,
        payload: element.payload,
      }
      
      // Always include visible/locked if they exist in the element
      if ('visible' in element) {
        updateData.visible = element.visible
      }
      if ('locked' in element) {
        updateData.locked = element.locked
      }

      await elementsApi.update(
        boardsStore.currentBoard.id,
        currentPageId.value,
        element.id,
        updateData
      )

      // Save successful
      pendingSaves.value.delete(elementId)
      retryCount.value.delete(elementId)
      saveStatus.value = 'saved'
      lastSaveTime.value = new Date()
      
      // Show "saved" status briefly, then return to idle
      setTimeout(() => {
        if (saveStatus.value === 'saved' && pendingSaves.value.size === 0) {
          saveStatus.value = 'idle'
        }
      }, SAVE_SUCCESS_DISPLAY_TIME)

    } catch (err: any) {
      pendingSaves.value.delete(elementId)
      
      // Handle 404 - element doesn't exist on server
      if (err.status === 404 || err.response?.status === 404) {
        console.warn('Element not found on server, removing from local state:', elementId)
        // Remove element from local state since it doesn't exist on server
        const elementIndex = elements.value.findIndex(el => el.id === elementId)
        if (elementIndex !== -1) {
          elements.value.splice(elementIndex, 1)
        }
        saveStatus.value = 'error'
        saveError.value = 'Element not found on server and was removed locally'
        error.value = saveError.value
        return
      }
      
      const currentRetryCount = retryCount.value.get(elementId) || 0
      
      // Handle conflicts (409 status)
      if (err.response?.status === 409) {
        console.warn('Conflict detected for element:', elementId, 'Attempting to resolve...')
        await handleConflict(element, err)
        return
      }
      
      // Retry logic for other errors
      if (currentRetryCount < MAX_RETRY_ATTEMPTS) {
        const nextAttempt = currentRetryCount + 1
        retryCount.value.set(elementId, nextAttempt)
        
        const delay = RETRY_DELAY_BASE * Math.pow(2, currentRetryCount) // Exponential backoff
        
        console.log(`Retrying save for element ${elementId} (attempt ${nextAttempt}/${MAX_RETRY_ATTEMPTS}) in ${delay}ms`)
        
        setTimeout(() => {
          saveElement(element)
        }, delay)
        
        return
      }
      
      // Max retries reached
      retryCount.value.delete(elementId)
      saveStatus.value = 'error'
      saveError.value = err.error?.message || `Failed to save element after ${MAX_RETRY_ATTEMPTS} attempts`
      error.value = saveError.value
      
      console.error('Failed to save element after retries:', elementId, err)
    }
  }

  // Handle conflicts when concurrent edits occur
  const handleConflict = async (_localElement: Element, _conflictError: any) => {
    try {
      // Fetch the latest version of the element from the server
      console.log('Fetching latest element version to resolve conflict...')
      
      // For now, we'll use a simple strategy: server wins
      // In a more sophisticated implementation, you could:
      // 1. Show a conflict resolution dialog to the user
      // 2. Merge changes intelligently based on timestamps
      // 3. Allow user to choose which version to keep
      
      // Reload the page elements to get the latest state
      if (currentPageId.value) {
        await loadPageElements(currentPageId.value)
        
        // Trigger canvas reload
        if (canvas.value) {
          const event = new CustomEvent('editor:reload-elements')
          window.dispatchEvent(event)
        }
      }
      
      saveError.value = 'Conflict resolved - loaded latest version from server'
      saveStatus.value = 'idle'
      
    } catch (err) {
      console.error('Failed to resolve conflict:', err)
      saveError.value = 'Failed to resolve conflict with server'
      saveStatus.value = 'error'
    }
  }

  const deleteElements = async (elementIds: string[]) => {
    const boardsStore = useBoardsStore()
    
    if (!boardsStore.currentBoard || !currentPageId.value || !boardsStore.editToken) {
      throw new Error('No board loaded or edit access required')
    }

    loading.value = true
    error.value = null

    try {
      // Save current state for undo
      saveCanvasState()
      
      // Delete elements from backend
      await Promise.all(
        elementIds.map(id =>
          elementsApi.delete(boardsStore.currentBoard!.id, currentPageId.value!, id)
        )
      )
      
      // Remove from local state
      elements.value = elements.value.filter(el => !elementIds.includes(el.id))
      selectedElementIds.value = selectedElementIds.value.filter(id => !elementIds.includes(id))
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to delete elements'
      throw err
    } finally {
      loading.value = false
    }
  }

  const reorderElements = async (updates: Array<{ id: string; z: number }>) => {
    const boardsStore = useBoardsStore()
    
    if (!boardsStore.currentBoard || !currentPageId.value || !boardsStore.editToken) {
      throw new Error('No board loaded or edit access required')
    }

    try {
      // Save current state for undo
      saveCanvasState()
      
      // Update backend
      await elementsApi.batchUpdateZIndex(
        boardsStore.currentBoard.id,
        currentPageId.value,
        { updates }
      )
      
      // Update local state
      updates.forEach(({ id, z }) => {
        const element = elements.value.find(el => el.id === id)
        if (element) {
          element.z = z
        }
      })
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to reorder elements'
      throw err
    }
  }

  const uploadImage = async (file: File) => {
    const boardsStore = useBoardsStore()
    
    if (!boardsStore.currentBoard || !boardsStore.editToken) {
      throw new Error('No board loaded or edit access required')
    }

    loading.value = true
    error.value = null

    try {
      const uploadResponse = await uploadsApi.uploadImage(boardsStore.currentBoard.id, file)
      return uploadResponse
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to upload image'
      throw err
    } finally {
      loading.value = false
    }
  }

  const toggleSnapToGrid = () => {
    snapToGrid.value = !snapToGrid.value
  }

  // Text styling methods
  const updateSelectedElementsStyle = (property: string, value: any) => {
    if (selectedElementIds.value.length === 0) return
    
    selectedElementIds.value.forEach(async (elementId) => {
      const element = elements.value.find(el => el.id === elementId)
      if (element && element.kind === 'text') {
        // Save canvas state before modification
        saveCanvasState()
        
        // Update element payload
        const payload = element.payload as any
        payload[property] = value
        
        // Update in store
        await updateElement(element.id, {
          payload: payload
        })
      }
    })
  }

  // Z-order methods
  const bringToFront = () => {
    if (selectedElementIds.value.length === 0) return
    
    saveCanvasState()
    
    selectedElementIds.value.forEach(async (elementId) => {
      const element = elements.value.find(el => el.id === elementId)
      if (element) {
        const maxZ = Math.max(...elements.value.map(el => el.z))
        await updateElement(element.id, { z: maxZ + 1 })
      }
    })
  }

  const bringForward = () => {
    if (selectedElementIds.value.length === 0) return
    
    saveCanvasState()
    
    selectedElementIds.value.forEach(async (elementId) => {
      const element = elements.value.find(el => el.id === elementId)
      if (element) {
        const higherElements = elements.value.filter(el => el.z > element.z)
        if (higherElements.length > 0) {
          const nextZ = Math.min(...higherElements.map(el => el.z))
          await updateElement(element.id, { z: nextZ + 1 })
        }
      }
    })
  }

  const sendBackward = () => {
    if (selectedElementIds.value.length === 0) return
    
    saveCanvasState()
    
    selectedElementIds.value.forEach(async (elementId) => {
      const element = elements.value.find(el => el.id === elementId)
      if (element) {
        const lowerElements = elements.value.filter(el => el.z < element.z)
        if (lowerElements.length > 0) {
          const prevZ = Math.max(...lowerElements.map(el => el.z))
          await updateElement(element.id, { z: prevZ - 1 })
        }
      }
    })
  }

  const sendToBack = () => {
    if (selectedElementIds.value.length === 0) return
    
    saveCanvasState()
    
    selectedElementIds.value.forEach(async (elementId) => {
      const element = elements.value.find(el => el.id === elementId)
      if (element) {
        const minZ = Math.min(...elements.value.map(el => el.z))
        await updateElement(element.id, { z: minZ - 1 })
      }
    })
  }

  const clearError = () => {
    error.value = null
  }

  const setElementVisibility = async (elementId: string, visible: boolean) => {
    console.log('Store setElementVisibility called:', elementId, visible)
    const element = elements.value.find(el => el.id === elementId)
    if (element) {
      console.log('Element found, updating visibility from', element.visible, 'to', visible)
      try {
        await updateElement(elementId, { visible })
        console.log('Visibility update completed successfully')
      } catch (error) {
        console.error('Failed to update element visibility:', error)
        throw error
      }
    } else {
      console.error('Element not found in store:', elementId)
      throw new Error('Element not found')
    }
  }

  const setElementLocked = async (elementId: string, locked: boolean) => {
    console.log('Store setElementLocked called:', elementId, locked)
    const element = elements.value.find(el => el.id === elementId)
    if (element) {
      console.log('Element found, updating locked state from', element.locked, 'to', locked)
      try {
        await updateElement(elementId, { locked })
        console.log('Lock state update completed successfully')
      } catch (error) {
        console.error('Failed to update element lock state:', error)
        throw error
      }
    } else {
      console.error('Element not found in store:', elementId)
      throw new Error('Element not found')
    }
  }

  const deleteElement = async (elementId: string) => {
    await deleteElements([elementId])
  }

  const reset = () => {
    selectedElementIds.value = []
    undoStack.value = []
    redoStack.value = []
    canvas.value = null
    currentPageId.value = null
    elements.value = []
    loading.value = false
    error.value = null
    
    // Reset auto-save state
    saveStatus.value = 'idle'
    lastSaveTime.value = null
    saveError.value = null
    pendingSaves.value.clear()
    retryCount.value.clear()
    
    // Clear timeouts
    if (saveTimeout) {
      clearTimeout(saveTimeout)
      saveTimeout = null
    }
    
    // Clear position save timeouts
    positionSaveTimeouts.forEach(timeout => clearTimeout(timeout))
    positionSaveTimeouts.clear()
  }

  return {
    // State
    selectedElementIds,
    undoStack,
    redoStack,
    snapToGrid,
    canvas,
    currentPageId,
    elements,
    loading,
    error,
    
    // Auto-save state
    saveStatus,
    lastSaveTime,
    saveError,
    pendingSaves,
    
    // Getters
    selectedElements,
    canUndo,
    canRedo,
    sortedElements,
    
    // Actions
    setCanvas,
    setCurrentPage,
    loadPageElements,
    selectElements,
    selectElement,
    toggleElementSelection,
    clearSelection,
    saveCanvasState,
    undo,
    redo,
    createElement,
    updateElement,
    deleteElements,
    deleteElement,
    saveElement,
    reorderElements,
    uploadImage,
    toggleSnapToGrid,
    updateSelectedElementsStyle,
    setElementVisibility,
    setElementLocked,
    bringToFront,
    bringForward,
    sendBackward,
    sendToBack,
    clearError,
    reset,
  }
})