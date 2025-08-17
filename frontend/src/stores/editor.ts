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
  const canvas = ref<fabric.Canvas | null>(null)
  const currentPageId = ref<string | null>(null)
  const elements = ref<Element[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Constants
  const MAX_HISTORY_SIZE = 50
  const DEBOUNCE_DELAY = 300

  // Getters
  const selectedElements = computed(() => 
    elements.value.filter(el => selectedElementIds.value.includes(el.id))
  )
  
  const canUndo = computed(() => undoStack.value.length > 0)
  const canRedo = computed(() => redoStack.value.length > 0)
  
  const sortedElements = computed(() => 
    [...elements.value].sort((a, b) => a.z - b.z)
  )

  // Debounced save function
  let saveTimeout: NodeJS.Timeout | null = null
  
  const debouncedSave = (element: Element) => {
    if (saveTimeout) {
      clearTimeout(saveTimeout)
    }
    
    saveTimeout = setTimeout(async () => {
      await saveElement(element)
    }, DEBOUNCE_DELAY)
  }

  // Actions
  const setCanvas = (fabricCanvas: fabric.Canvas) => {
    canvas.value = fabricCanvas
  }

  const setCurrentPage = (pageId: string) => {
    currentPageId.value = pageId
    elements.value = []
    selectedElementIds.value = []
  }

  const loadPageElements = async (pageId: string) => {
    const boardsStore = useBoardsStore()
    const page = boardsStore.pages.find(p => p.id === pageId)
    
    if (page?.elements) {
      elements.value = page.elements
    } else {
      elements.value = []
    }
  }

  const selectElements = (elementIds: string[]) => {
    selectedElementIds.value = elementIds
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
    
    // Update canvas if available
    if (canvas.value) {
      // Canvas update logic would go here
      canvas.value.renderAll()
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
    
    // Update canvas if available
    if (canvas.value) {
      // Canvas update logic would go here
      canvas.value.renderAll()
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
        { kind, x, y, w, h, rotation, payload }
      )
      
      elements.value.push(newElement)
      return newElement
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to create element'
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateElement = async (elementId: string, updates: Partial<Element>) => {
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

    // Debounce API call
    debouncedSave(elements.value[index])
  }

  const saveElement = async (element: Element) => {
    const boardsStore = useBoardsStore()
    
    if (!boardsStore.currentBoard || !currentPageId.value || !boardsStore.editToken) {
      return
    }

    try {
      await elementsApi.update(
        boardsStore.currentBoard.id,
        currentPageId.value,
        element.id,
        {
          x: element.x,
          y: element.y,
          w: element.w,
          h: element.h,
          rotation: element.rotation,
          z: element.z,
          payload: element.payload,
        }
      )
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to save element'
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

  const clearError = () => {
    error.value = null
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
    
    if (saveTimeout) {
      clearTimeout(saveTimeout)
      saveTimeout = null
    }
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
    clearSelection,
    saveCanvasState,
    undo,
    redo,
    createElement,
    updateElement,
    deleteElements,
    reorderElements,
    uploadImage,
    toggleSnapToGrid,
    clearError,
    reset,
  }
})