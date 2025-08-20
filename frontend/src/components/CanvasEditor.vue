<template>
  <div class="canvas-editor-container relative w-full h-full">
    <!-- Canvas Container -->
    <div 
      ref="canvasContainer" 
      class="canvas-wrapper relative w-full h-full overflow-hidden"
      :class="{ 'snap-to-grid': snapToGrid }"
    >
      <canvas 
        id="fabric-canvas"
        ref="canvasElement"
        class="fabric-canvas"
        :width="canvasWidth"
        :height="canvasHeight"
      />
      
      <!-- Grid Overlay (when snap to grid is enabled) -->
      <div 
        v-if="snapToGrid" 
        class="grid-overlay absolute inset-0 pointer-events-none"
        :style="gridStyle"
      />
    </div>

    <!-- Loading Overlay -->
    <div 
      v-if="loading" 
      class="absolute inset-0 bg-white bg-opacity-75 flex items-center justify-center z-10"
    >
      <LoadingSpinner />
    </div>

    <!-- Error Message -->
    <ErrorMessage 
      v-if="error" 
      :message="error" 
      @dismiss="clearError"
      class="absolute top-4 right-4 z-20"
    />

    <!-- Image Viewer Modal -->
    <ImageViewer
      :is-visible="imageViewer.isVisible"
      :image-data="imageViewer.imageData"
      :element-id="imageViewer.elementId"
      :is-edit-mode="boardsStore.isEditMode"
      @close="closeImageViewer"
      @update-description="updateImageDescription"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed, nextTick } from 'vue'
import { fabric } from 'fabric'
import { useEditorStore } from '@/stores/editor'
import { useBoardsStore } from '@/stores/boards'
import LoadingSpinner from './LoadingSpinner.vue'
import ErrorMessage from './ErrorMessage.vue'
import ImageViewer from './ImageViewer.vue'
import type { Element, TextPayload, ImagePayload, ShapePayload, StickerPayload } from '@/types'

// Props
interface Props {
  width?: number
  height?: number
}

const props = withDefaults(defineProps<Props>(), {
  width: 800,
  height: 600
})

// Stores
const editorStore = useEditorStore()
const boardsStore = useBoardsStore()

// Refs
const canvasContainer = ref<HTMLDivElement>()
const canvasElement = ref<HTMLCanvasElement>()
const fabricCanvas = ref<fabric.Canvas>()
const isTextEditing = ref(false) // Track text editing state
const isUserInteracting = ref(false) // Track if user is actively interacting with elements
const reloadDebounceTimer = ref<NodeJS.Timeout | null>(null)

// Image viewer state
const imageViewer = ref({
  isVisible: false,
  imageData: null as ImagePayload | null,
  elementId: null as string | null
})

// Computed
const canvasWidth = computed(() => props.width)
const canvasHeight = computed(() => props.height)
const snapToGrid = computed(() => editorStore.snapToGrid)
const loading = computed(() => editorStore.loading)
const error = computed(() => editorStore.error)
const elements = computed(() => editorStore.sortedElements)
const selectedElementIds = computed(() => editorStore.selectedElementIds)

// Grid styling
const gridStyle = computed(() => ({
  backgroundImage: `
    linear-gradient(to right, rgba(0,0,0,0.1) 1px, transparent 1px),
    linear-gradient(to bottom, rgba(0,0,0,0.1) 1px, transparent 1px)
  `,
  backgroundSize: '8px 8px'
}))

// Constants
const GRID_SIZE = 8
const MIN_ELEMENT_SIZE = 16

// Canvas initialization
const initializeCanvas = async () => {
  if (!canvasElement.value) {
    console.log('Canvas element not ready, waiting...')
    await nextTick()
    if (!canvasElement.value) return
  }

  console.log('Initializing canvas...')
  fabricCanvas.value = new fabric.Canvas(canvasElement.value, {
    width: canvasWidth.value,
    height: canvasHeight.value,
    backgroundColor: 'transparent',
    selection: true,
    preserveObjectStacking: true,
    imageSmoothingEnabled: true,
    allowTouchScrolling: false,
    fireRightClick: true,
    stopContextMenu: true,
    // Ensure interactions are enabled
    interactive: true,
    moveCursor: 'move',
    hoverCursor: 'move',
    defaultCursor: 'default',
    rotationCursor: 'crosshair',
  })
  
  // Ensure global prototype settings are correct (fix for potential global overrides)
  fabric.Object.prototype.hasControls = true
  fabric.Object.prototype.hasBorders = true
  fabric.Object.prototype.selectable = true
  fabric.Object.prototype.evented = true
  
  // CRITICAL: Ensure global scaling settings are correct
  fabric.Object.prototype.lockScalingX = false
  fabric.Object.prototype.lockScalingY = false
  fabric.Object.prototype.lockUniScaling = false
  ;(fabric.Object.prototype as any).uniformScaling = false
  ;(fabric.Object.prototype as any).centeredScaling = false
  
  console.log('Canvas initialized with global prototype settings verified')

  // Store canvas reference in editor store
  editorStore.setCanvas(fabricCanvas.value)

  // Set up event listeners
  setupCanvasEvents()
  
  // Load existing elements after a short delay to ensure everything is ready
  await nextTick()
  await loadElementsToCanvas()
  
  console.log('Canvas initialized successfully')
}

// Canvas event handlers
const setupCanvasEvents = () => {
  if (!fabricCanvas.value) return

  const canvas = fabricCanvas.value

  // Selection events
  canvas.on('selection:created', handleSelectionChange)
  canvas.on('selection:updated', handleSelectionChange)
  canvas.on('selection:cleared', handleSelectionClear)
  
  // Add immediate control activation on selection
  canvas.on('path:created', (_e: fabric.IEvent) => {
    console.log('Path created event')
  })
  
  // Add event for when an object is selected (immediate response)
  canvas.on('object:selected', (e: fabric.IEvent) => {
    console.log('Object selected event - immediately fixing controls')
    if (e.target && boardsStore.isEditMode) {
      // Immediate control fix with no delay
      fixObjectControls(e.target)
      canvas.renderAll()
    }
  })

  // Object modification events
  canvas.on('object:modified', handleObjectModified)
  canvas.on('object:moving', handleObjectMoving)
  canvas.on('object:scaling', handleObjectScaling)
  canvas.on('object:rotating', handleObjectRotating)
  
  // Add more detailed event handling for debugging
  canvas.on('mouse:over', handleMouseOver)
  canvas.on('mouse:out', handleMouseOut)
  canvas.on('object:scaling', handleObjectScaling)
  canvas.on('before:selection:cleared', handleBeforeSelectionCleared)
  
  // Add specific control and corner events
  canvas.on('mouse:over:before', (e: fabric.IEvent) => {
    console.log('Mouse over before:', e.target?.type)
  })
  
  // Monitor for corner detection specifically
  canvas.on('mouse:down:before', (e: fabric.IEvent) => {
    const pointer = canvas.getPointer(e.e)
    const activeObject = canvas.getActiveObject()
    if (activeObject) {
      console.log('Mouse down before - checking corner detection:', {
        pointer,
        activeObjectType: activeObject.type,
        hasControls: activeObject.hasControls,
        cornerSize: activeObject.cornerSize
      })
    }
  })
  
  // Try to capture scaling events more directly
  canvas.on('object:scaling', (e: fabric.IEvent) => {
    console.log('SCALING EVENT DETECTED!', e.target?.type)
  })
  
  // Listen for transformation events
  canvas.on('before:transform', (e: fabric.IEvent) => {
    console.log('Before transform:', e.target?.type, e.transform?.action)
    
    // Ensure scaling is enabled for the target object
    if (e.target && e.transform?.action?.includes('scale')) {
      const obj = e.target as any
      obj.lockScalingX = false
      obj.lockScalingY = false
      obj.lockUniScaling = false
      console.log('Enabling scaling for transform on:', obj.elementId)
    }
  })
  
  canvas.on('object:resizing', (e: fabric.IEvent) => {
    console.log('Object resizing:', e.target?.type)
  })
  
  // Add a specific event for detecting when user tries to scale
  canvas.on('mouse:down', (e: fabric.IEvent) => {
    if (e.target) {
      const obj = e.target as any
      // Force enable scaling when user interacts with corners
      obj.lockScalingX = false
      obj.lockScalingY = false
      obj.lockUniScaling = false
      obj.uniformScaling = false
      obj.centeredScaling = false
      console.log('Mouse down - ensuring scaling is enabled for:', obj.elementId)
    }
  })
  
  // Add event to detect corner interactions specifically
  canvas.on('selection:created', (e: fabric.IEvent) => {
    const target = e.target || (e as any).selected?.[0]
    if (target) {
      const obj = target as any
      // Immediately ensure scaling is enabled when object is selected
      obj.lockScalingX = false
      obj.lockScalingY = false
      obj.lockUniScaling = false
      obj.uniformScaling = false
      obj.centeredScaling = false
      
      // Force update the controls to be scaling controls
      setTimeout(() => {
        fixObjectControls(obj)
      }, 10)
      
      console.log('Selection created - enabled scaling for:', obj.elementId)
    }
  })

  // Mouse events for element creation
  canvas.on('mouse:down', handleMouseDown)
  canvas.on('mouse:up', handleMouseUp)

  // Text editing events
  canvas.on('mouse:dblclick', handleDoubleClick)
  canvas.on('text:changed', handleTextChanged)
  canvas.on('text:editing:exited', handleTextEditingExited)

  // Keyboard events
  document.addEventListener('keydown', handleKeyDown)
}

// Selection handlers
const handleSelectionChange = (_e: fabric.IEvent) => {
  const activeObjects = fabricCanvas.value?.getActiveObjects() || []
  const elementIds = activeObjects
    .map((obj: any) => obj.elementId)
    .filter(Boolean)
  
  console.log('Selection changed:', {
    objectCount: activeObjects.length,
    elementIds,
    userInteracting: isUserInteracting.value,
    objects: activeObjects.map((obj: any) => ({
      type: obj.type,
      elementId: obj.elementId,
      hasControls: obj.hasControls,
      selectable: obj.selectable,
      evented: obj.evented,
      lockScalingX: obj.lockScalingX,
      lockScalingY: obj.lockScalingY
    }))
  })
  
  // If user is actively interacting, be more careful about selection changes
  if (isUserInteracting.value && activeObjects.length > 0) {
    console.log('Selection change during user interaction - maintaining focus')
  }
  
  // For each selected object, ensure controls are properly enabled
  activeObjects.forEach((obj: any) => {
    if (boardsStore.isEditMode) {
      // Force enable all controls and interactions
      obj.hasControls = true
      obj.hasBorders = true
      obj.hasRotatingPoint = true
      obj.lockScalingX = false
      obj.lockScalingY = false
      obj.lockMovementX = false
      obj.lockMovementY = false
      obj.lockRotation = false
      obj.lockUniScaling = false
      obj.lockScalingFlip = false
      obj.selectable = true
      obj.evented = true
      
      // Set visual properties for better control visibility
      obj.cornerSize = 12
      obj.transparentCorners = false
      obj.cornerColor = '#2563EB'
      obj.borderColor = '#2563EB'
      obj.cornerStrokeColor = '#1E40AF'
      
      // Explicitly set all controls to be visible
      obj.setControlsVisibility({
        mtr: true, // rotation control
        tr: true,  // top-right
        br: true,  // bottom-right
        bl: true,  // bottom-left
        tl: true,  // top-left
        mt: true,  // middle-top
        mb: true,  // middle-bottom
        ml: true,  // middle-left
        mr: true   // middle-right
      })
      
      // CRITICAL: Ensure corner controls have proper action handlers for scaling
      const corners = ['tl', 'tr', 'br', 'bl']
      corners.forEach(cornerKey => {
        if (obj.controls && obj.controls[cornerKey]) {
          obj.controls[cornerKey].actionName = 'scale'
          // Use the default scaling action handler from Fabric's prototype
          const defaultControl = fabric.Object.prototype.controls[cornerKey]
          if (defaultControl && defaultControl.actionHandler) {
            obj.controls[cornerKey].actionHandler = defaultControl.actionHandler
          }
          obj.controls[cornerKey].cursorStyle = 'crosshair'
        }
      })
      
      // Also ensure middle controls have proper action handlers
      const middles = ['mt', 'mb', 'ml', 'mr']
      middles.forEach(middleKey => {
        if (obj.controls && obj.controls[middleKey]) {
          obj.controls[middleKey].actionName = 'scale'
          // Use the default scaling action handler from Fabric's prototype
          const defaultControl = fabric.Object.prototype.controls[middleKey]
          if (defaultControl && defaultControl.actionHandler) {
            obj.controls[middleKey].actionHandler = defaultControl.actionHandler
          }
        }
      })
      
      // Update coordinates to ensure proper interaction
      obj.setCoords()
      
      console.log('Aggressively enabled controls for object:', {
        elementId: obj.elementId,
        type: obj.type,
        hasControls: obj.hasControls,
        cornerSize: obj.cornerSize,
        locks: {
          scalingX: obj.lockScalingX,
          scalingY: obj.lockScalingY,
          movement: obj.lockMovementX || obj.lockMovementY,
          rotation: obj.lockRotation
        }
      })
    }
  })
  
  if (activeObjects.length > 0) {
    fabricCanvas.value?.renderAll()
    
    // CRITICAL FIX: Add aggressive control reactivation after selection
    // This addresses the issue where controls don't work until group/ungroup
    setTimeout(() => {
      if (fabricCanvas.value && activeObjects.length > 0) {
        activeObjects.forEach((obj: any) => {
          if (boardsStore.isEditMode) {
            // Force complete control reinitialization
            fixObjectControls(obj)
            
            // Additional aggressive reactivation
            obj.setCoords()
            
            // Force the canvas to recognize this object as actively selected
            fabricCanvas.value?.setActiveObject(obj)
          }
        })
        
        // Force final render to ensure controls are visible and interactive
        fabricCanvas.value.renderAll()
        
        console.log('Aggressively reactivated controls for selected objects')
      }
    }, 10) // Small delay to ensure selection is fully processed
  }
  
  editorStore.selectElements(elementIds)
}

const handleSelectionClear = () => {
  editorStore.clearSelection()
}

// Object modification handlers
const handleObjectModified = (e: fabric.IEvent) => {
  if (!e.target) return
  
  // Don't update during text editing to prevent interrupting the editing session
  if (isTextEditing.value) return
  
  const obj = e.target
  const elementId = (obj as any).elementId
  
  console.log('Object modified event triggered for element:', elementId)
  console.log('Object properties:', {
    left: obj.left,
    top: obj.top,
    width: obj.width,
    height: obj.height,
    scaleX: obj.scaleX,
    scaleY: obj.scaleY,
    angle: obj.angle
  })
  
  if (elementId) {
    // Immediately fix controls to ensure they remain active after modification
    setTimeout(() => {
      fixObjectControls(obj)
      fabricCanvas.value?.renderAll()
    }, 10)
    
    // Use setTimeout to prevent immediate canvas reload during interaction
    setTimeout(() => {
      updateElementFromFabricObject(elementId, obj)
    }, 50)
  }
}

const handleObjectMoving = (e: fabric.IEvent) => {
  if (!e.target || !snapToGrid.value) return
  
  // Set user interaction flag during movement
  isUserInteracting.value = true
  
  const obj = e.target
  
  // Snap to grid
  obj.left = Math.round((obj.left || 0) / GRID_SIZE) * GRID_SIZE
  obj.top = Math.round((obj.top || 0) / GRID_SIZE) * GRID_SIZE
}

const handleObjectScaling = (e: fabric.IEvent) => {
  if (!e.target) return
  
  // Set user interaction flag during scaling
  isUserInteracting.value = true
  
  const obj = e.target
  const elementId = (obj as any).elementId
  
  console.log('Object scaling event triggered for element:', elementId)
  console.log('Current scale:', obj.scaleX, obj.scaleY)
  
  // Allow free scaling, only apply grid snap if enabled
  if (elementId && snapToGrid.value) {
    // Snap dimensions to grid
    const width = (obj.width || 0) * (obj.scaleX || 1)
    const height = (obj.height || 0) * (obj.scaleY || 1)
    
    const snappedWidth = Math.max(MIN_ELEMENT_SIZE, Math.round(width / GRID_SIZE) * GRID_SIZE)
    const snappedHeight = Math.max(MIN_ELEMENT_SIZE, Math.round(height / GRID_SIZE) * GRID_SIZE)
    
    obj.scaleX = snappedWidth / (obj.width || 1)
    obj.scaleY = snappedHeight / (obj.height || 1)
    
    console.log('Snapped to grid. New scale:', obj.scaleX, obj.scaleY)
  }
}

const handleBeforeSelectionCleared = () => {
  // Save any pending changes before selection is cleared
  console.log('Selection about to be cleared')
}

const handleMouseOver = (e: fabric.IEvent) => {
  if (e.target) {
    console.log('Mouse over object:', {
      type: e.target.type,
      elementId: (e.target as any).elementId,
      hasControls: e.target.hasControls,
      selectable: e.target.selectable
    })
  }
}

const handleMouseOut = (e: fabric.IEvent) => {
  if (e.target) {
    console.log('Mouse out object:', {
      type: e.target.type,
      elementId: (e.target as any).elementId
    })
  }
}

const handleObjectRotating = (e: fabric.IEvent) => {
  if (!e.target) return
  
  // Set user interaction flag during rotation
  isUserInteracting.value = true
  
  const obj = e.target
  const elementId = (obj as any).elementId
  
  if (elementId) {
    // Debounced update will be handled by handleObjectModified
  }
}

// Mouse event handlers
const handleMouseDown = (e: fabric.IEvent) => {
  console.log('Mouse down on:', e.target?.type)
  
  // Set user interaction flag to prevent canvas reloads
  isUserInteracting.value = true
  
  // Log detailed object properties when clicking on an object
  if (e.target) {
    const obj = e.target as any
    console.log('Object detailed properties:', {
      type: obj.type,
      elementId: obj.elementId,
      selectable: obj.selectable,
      evented: obj.evented,
      hasControls: obj.hasControls,
      hasBorders: obj.hasBorders,
      hasRotatingPoint: obj.hasRotatingPoint,
      lockScalingX: obj.lockScalingX,
      lockScalingY: obj.lockScalingY,
      lockMovementX: obj.lockMovementX,
      lockMovementY: obj.lockMovementY,
      cornerSize: obj.cornerSize,
      transparentCorners: obj.transparentCorners,
      cornerColor: obj.cornerColor,
      borderColor: obj.borderColor,
      controls: Object.keys(obj.controls || {}),
      canvas: !!obj.canvas,
      isOnCanvas: fabricCanvas.value?.getObjects().includes(obj)
    })
    
    // Check if controls are visible
    const controlsVisibility = (obj as any)._getControlsVisibility ? (obj as any)._getControlsVisibility() : 'method not available'
    console.log('Controls visibility:', controlsVisibility)
  }
  
  // Handle canvas clicks for deselection
  if (!e.target) {
    editorStore.clearSelection()
    return
  }

  // IMPORTANT: Immediately fix controls when clicking on any object
  if (e.target && boardsStore.isEditMode) {
    console.log('Mouse down on object - immediately fixing controls')
    
    // Force activate controls immediately on click
    const clickedObject = e.target as any
    if (clickedObject && clickedObject.elementId) {
      // Apply the control fix immediately - no delay
      fixObjectControls(clickedObject)
      
      // Additional aggressive control activation
      setTimeout(() => {
        if (fabricCanvas.value && clickedObject) {
          // Force the object to be active and interactive
          fabricCanvas.value.setActiveObject(clickedObject)
          
          // Double-check all control properties are set
          clickedObject.hasControls = true
          clickedObject.hasBorders = true
          clickedObject.selectable = true
          clickedObject.evented = true
          
          // Force coordinate update and render
          clickedObject.setCoords()
          fabricCanvas.value.renderAll()
          
          console.log('Controls immediately activated for clicked object:', clickedObject.elementId)
        }
      }, 5) // Very minimal delay for proper activation
    }
  }

  // Handle text editing with single click if already selected
  if (e.target.type === 'i-text' && boardsStore.isEditMode) {
    const target = e.target as any
    // If the text is already selected, enter edit mode on single click
    if (fabricCanvas.value?.getActiveObject() === target) {
      console.log('Single click on selected text, entering edit mode')
      setTimeout(() => {
        try {
          target.enterEditing()
          target.selectAll()
          fabricCanvas.value?.renderAll()
        } catch (error) {
          console.error('Failed to enter edit mode on single click:', error)
        }
      }, 100)
    }
  }

  // Note: Image viewer opening moved to double-click handler
  // Single click on images will only select them, not open the viewer
}

const handleMouseUp = (e: fabric.IEvent) => {
  // Clear user interaction flag after a delay to allow for operations to complete
  setTimeout(() => {
    isUserInteracting.value = false
    console.log('User interaction ended')
  }, 500) // Give enough time for any modification events to finish
  
  // Double-check that controls are active after mouse up
  if (e.target && boardsStore.isEditMode) {
    setTimeout(() => {
      const target = e.target as any
      if (target && target.elementId) {
        console.log('Mouse up - verifying controls are active')
        
        // Ensure the object is selected and controls are active
        fabricCanvas.value?.setActiveObject(target)
        fixObjectControls(target)
        fabricCanvas.value?.renderAll()
        
        console.log('Controls verified active on mouse up for:', target.elementId)
      }
    }, 50) // Small delay to ensure any selection events have fired
  }
}

// Text editing handlers
const handleDoubleClick = (e: fabric.IEvent) => {
  console.log('Double click detected on:', e.target?.type, 'Edit mode:', boardsStore.isEditMode)
  if (!boardsStore.isEditMode) return
  
  const target = e.target as any
  if (target && target.type === 'i-text') {
    console.log('Entering text editing mode for:', target.elementId)
    e.e.preventDefault() // Prevent default double-click behavior
    e.e.stopPropagation()
    
    // Set text editing flag
    isTextEditing.value = true
    
    // Enter text editing mode
    setTimeout(() => {
      try {
        target.enterEditing()
        target.selectAll()
        fabricCanvas.value?.renderAll()
        console.log('Successfully entered edit mode via double-click')
      } catch (error) {
        console.error('Failed to enter edit mode via double-click:', error)
      }
    }, 50)
  }

  // Handle image double-click to open viewer
  if (target && target.type === 'image') {
    console.log('Double-click on image - opening viewer')
    e.e.preventDefault()
    e.e.stopPropagation()
    
    const elementId = target.elementId
    if (elementId) {
      const element = editorStore.elements.find(el => el.id === elementId)
      if (element && element.kind === 'image') {
        openImageViewer(element.id, element.payload as ImagePayload)
      }
    }
  }
}

const handleTextChanged = (e: fabric.IEvent) => {
  console.log('Text changed detected')
  if (!boardsStore.isEditMode) return
  
  // Set text editing flag to prevent canvas updates
  isTextEditing.value = true
  
  const target = e.target as any
  if (target && target.elementId) {
    console.log('Updating text content:', target.text)
    // Update the element in the store with new text content
    const elementId = target.elementId
    const element = editorStore.elements.find(el => el.id === elementId)
    
    if (element && element.kind === 'text') {
      const updatedPayload = {
        ...element.payload,
        content: target.text || ''
      }
      
      editorStore.updateElement(elementId, {
        payload: updatedPayload
      })
    }
  }
}

const handleTextEditingExited = (e: fabric.IEvent) => {
  console.log('Text editing exited')
  if (!boardsStore.isEditMode) return
  
  // Clear text editing flag
  isTextEditing.value = false
  
  const target = e.target as any
  if (target && target.elementId) {
    console.log('Final text update on exit:', target.text)
    // Final update when exiting text editing mode
    const elementId = target.elementId
    const element = editorStore.elements.find(el => el.id === elementId)
    
    if (element && element.kind === 'text') {
      const updatedPayload = {
        ...element.payload,
        content: target.text || ''
      }
      
      editorStore.updateElement(elementId, {
        payload: updatedPayload
      })
    }
  }
}

// Image viewer methods
const openImageViewer = (elementId: string, imageData: ImagePayload) => {
  imageViewer.value = {
    isVisible: true,
    imageData,
    elementId
  }
}

const closeImageViewer = () => {
  imageViewer.value = {
    isVisible: false,
    imageData: null,
    elementId: null
  }
}

const updateImageDescription = async (elementId: string, description: string) => {
  const element = editorStore.elements.find(el => el.id === elementId)
  if (element && element.kind === 'image') {
    const updatedPayload = {
      ...element.payload as ImagePayload,
      description: description
    }
    
    await editorStore.updateElement(elementId, {
      payload: updatedPayload
    })

    // Update the viewer data
    if (imageViewer.value.elementId === elementId) {
      imageViewer.value.imageData = updatedPayload
    }

    console.log('Image description updated:', description)
  }
}

// Keyboard event handlers
const handleKeyDown = (e: KeyboardEvent) => {
  // Delete key handling
  if (e.key === 'Delete' || e.key === 'Backspace') {
    if (selectedElementIds.value.length > 0 && boardsStore.isEditMode) {
      e.preventDefault()
      deleteSelectedElements()
    }
  }
  
  // Enter key for text editing
  if (e.key === 'Enter' && boardsStore.isEditMode) {
    const activeObject = fabricCanvas.value?.getActiveObject() as any
    if (activeObject && activeObject.type === 'i-text') {
      console.log('Enter key pressed, entering text edit mode')
      e.preventDefault()
      try {
        activeObject.enterEditing()
        activeObject.selectAll()
        fabricCanvas.value?.renderAll()
        console.log('Successfully entered edit mode via Enter key')
      } catch (error) {
        console.error('Failed to enter edit mode via Enter key:', error)
      }
    }
  }
  
  // Undo/Redo handling
  if (e.ctrlKey || e.metaKey) {
    switch (e.key) {
      case 'z':
        if (!e.shiftKey) {
          e.preventDefault()
          editorStore.undo()
        } else {
          e.preventDefault()
          editorStore.redo()
        }
        break
      case 'y':
        e.preventDefault()
        editorStore.redo()
        break
      case 't':
        if (boardsStore.isEditMode) {
          e.preventDefault()
          addTextElement()
        }
        break
      case 'i':
        if (boardsStore.isEditMode) {
          e.preventDefault()
          const placeholderUrl = 'https://via.placeholder.com/200x200/3B82F6/FFFFFF?text=Image'
          addImageElement(placeholderUrl, 100, 100, 200, 200)
        }
        break
      case 's':
        if (boardsStore.isEditMode) {
          e.preventDefault()
          addShapeElement('rectangle', 100, 100)
        }
        break
      case 'b':
        if (selectedElementIds.value.length > 0 && boardsStore.isEditMode) {
          e.preventDefault()
          editorStore.bringToFront()
        }
        break
      case 'f':
        if (selectedElementIds.value.length > 0 && boardsStore.isEditMode) {
          e.preventDefault()
          editorStore.sendToBack()
        }
        break
    }
  }
  
  // Select all
  if ((e.ctrlKey || e.metaKey) && e.key === 'a') {
    e.preventDefault()
    selectAllElements()
  }
  
  // Grid toggle
  if (e.key === 'g' && !e.ctrlKey && !e.metaKey) {
    e.preventDefault()
    editorStore.toggleSnapToGrid()
  }
  
  // Force fix controls with 'R' key (for "Resize" controls)
  if (e.key === 'r' && !e.ctrlKey && !e.metaKey && boardsStore.isEditMode) {
    e.preventDefault()
    const activeObject = fabricCanvas.value?.getActiveObject()
    if (activeObject) {
      console.log('R key pressed - force fixing controls for active object')
      fixObjectControls(activeObject)
      fabricCanvas.value?.renderAll()
      console.log('Controls force-fixed via R key')
    } else {
      console.log('R key pressed - fixing controls for all objects')
      const allObjects = fabricCanvas.value?.getObjects() || []
      allObjects.forEach(obj => fixObjectControls(obj))
      fabricCanvas.value?.renderAll()
      console.log('Controls force-fixed for all objects via R key')
    }
  }
  
  // Emergency reset with 'Ctrl+R' - implements the full sanity checklist
  if (e.ctrlKey && e.key === 'r' && boardsStore.isEditMode) {
    e.preventDefault()
    console.log('Emergency reset triggered - applying full sanity checklist')
    
    if (fabricCanvas.value) {
      // Reset canvas zoom and viewport
      fabricCanvas.value.setZoom(1)
      fabricCanvas.value.viewportTransform = [1, 0, 0, 1, 0, 0]
      fabricCanvas.value.selection = true
      
      // Reset prototype controls if they were modified
      fabric.Object.prototype.hasControls = true
      fabric.Object.prototype.hasBorders = true
      
      const activeObject = fabricCanvas.value.getActiveObject()
      if (activeObject) {
        // Apply the full sanity checklist
        (activeObject as any).set({
          selectable: true,
          hasControls: true,
          hasBorders: true,
          lockScalingX: false,
          lockScalingY: false,
          evented: true,
          lockMovementX: false,
          lockMovementY: false,
          lockRotation: false,
          lockUniScaling: false,
          lockScalingFlip: false,
          uniformScaling: false,
          centeredScaling: false
        })
        
        // Exit editing mode if in IText
        if (activeObject.type === 'i-text' && (activeObject as any).isEditing) {
          (activeObject as any).exitEditing()
        }
        
        // Restore default controls
        Object.assign((activeObject as any).controls, fabric.Object.prototype.controls)
        
        // Set visual properties
        (activeObject as any).cornerSize = 14
        ;(activeObject as any).transparentCorners = false
        ;(activeObject as any).cornerColor = '#2563EB'
        ;(activeObject as any).borderColor = '#2563EB'
        
        // Ensure all controls are visible
        activeObject.setControlsVisibility({
          mtr: true, tr: true, br: true, bl: true, tl: true,
          mt: true, mb: true, ml: true, mr: true
        })
        
        fabricCanvas.value.setActiveObject(activeObject)
        activeObject.setCoords()
        fabricCanvas.value.renderAll()
        
        console.log('Emergency reset completed for active object')
      } else {
        console.log('No active object - reset applied to canvas only')
      }
    }
  }
}

// Element management
const loadElementsToCanvas = async () => {
  if (!fabricCanvas.value) {
    console.log('Canvas not ready for loading elements')
    return
  }
  
  // Don't reload canvas during text editing
  if (isTextEditing.value) {
    console.log('Skipping canvas reload during text editing')
    return
  }
  
  // Check if canvas context is valid
  try {
    const context = fabricCanvas.value.getContext()
    if (!context) {
      console.log('Canvas context not ready')
      return
    }
  } catch (error) {
    console.log('Canvas not properly initialized yet')
    return
  }
  
  console.log(`Loading ${elements.value.length} elements to canvas`)
  
  try {
    // Get existing objects on canvas
    const existingObjects = fabricCanvas.value.getObjects()
    const existingElementIds = existingObjects.map((obj: any) => obj.elementId).filter(Boolean)
    
    // Find elements that need to be added
    const elementsToAdd = elements.value.filter(element => 
      element && !existingElementIds.includes(element.id)
    )
    
    // Find objects that need to be removed (elements that no longer exist)
    const objectsToRemove = existingObjects.filter((obj: any) => 
      obj.elementId && !elements.value.some(el => el.id === obj.elementId)
    )
    
    // Remove obsolete objects
    objectsToRemove.forEach(obj => {
      fabricCanvas.value?.remove(obj)
    })
    
    // Add new elements
    for (const element of elementsToAdd) {
      await addElementToCanvas(element)
    }
    
    // Only render if we made changes
    if (elementsToAdd.length > 0 || objectsToRemove.length > 0) {
      fabricCanvas.value.renderAll()
      console.log(`Added ${elementsToAdd.length} elements, removed ${objectsToRemove.length} objects`)
    }
    
    // Fix controls for all objects after loading (ensure consistency)
    if (boardsStore.isEditMode) {
      setTimeout(() => {
        const allObjects = fabricCanvas.value?.getObjects() || []
        allObjects.forEach(obj => {
          fixObjectControls(obj)
        })
        
        // Preserve active selection if there was one
        const activeElementId = selectedElementIds.value[0]
        if (activeElementId) {
          const activeObject = allObjects.find((obj: any) => obj.elementId === activeElementId)
          if (activeObject) {
            fabricCanvas.value?.setActiveObject(activeObject)
            // Extra aggressive control fixing for the active object
            setTimeout(() => {
              if (activeObject && fabricCanvas.value) {
                fixObjectControls(activeObject)
                fabricCanvas.value.renderAll()
                console.log('Restored controls for active object after reload:', activeElementId)
              }
            }, 50)
          }
        }
        
        fabricCanvas.value?.renderAll()
        console.log('Fixed controls for all loaded objects')
      }, 100)
    }
    
    console.log('Elements loaded to canvas successfully')
  } catch (error) {
    console.error('Error loading elements to canvas:', error)
  }
}

const addElementToCanvas = async (element: Element): Promise<fabric.Object | null> => {
  if (!fabricCanvas.value || !element) return null
  
  // Check if element already exists on canvas
  const existingObject = fabricCanvas.value.getObjects().find((obj: any) => obj.elementId === element.id)
  if (existingObject) {
    console.log('Element already exists on canvas:', element.id)
    return existingObject
  }
  
  let fabricObject: fabric.Object | null = null
  
  try {
    switch (element.kind) {
      case 'text':
        fabricObject = await createTextObject(element)
        break
      case 'image':
        fabricObject = await createImageObject(element)
        break
      case 'shape':
        fabricObject = await createShapeObject(element)
        break
      case 'sticker':
        fabricObject = await createStickerObject(element)
        break
    }
    
    if (fabricObject) {
      // Set common properties
      fabricObject.left = element.x
      fabricObject.top = element.y
      fabricObject.angle = element.rotation
      fabricObject.setCoords()
      
      // Store element ID for reference
      ;(fabricObject as any).elementId = element.id
      
      // Ensure scaling and controls are properly enabled
      fabricObject.selectable = boardsStore.isEditMode
      fabricObject.evented = boardsStore.isEditMode
      fabricObject.hasControls = boardsStore.isEditMode
      fabricObject.hasBorders = boardsStore.isEditMode
      fabricObject.hasRotatingPoint = boardsStore.isEditMode
      fabricObject.lockUniScaling = false
      fabricObject.lockScalingFlip = false
      fabricObject.lockMovementX = false
      fabricObject.lockMovementY = false
      fabricObject.lockRotation = false
      fabricObject.lockScalingX = false
      fabricObject.lockScalingY = false
      
      // CRITICAL: Force scaling properties for proper corner behavior
      ;(fabricObject as any).uniformScaling = false
      ;(fabricObject as any).centeredScaling = false
      
      // Debug logging
      console.log(`Added ${element.kind} element to canvas:`, {
        id: element.id,
        selectable: fabricObject.selectable,
        hasControls: fabricObject.hasControls,
        hasBorders: fabricObject.hasBorders,
        lockScalingX: fabricObject.lockScalingX,
        lockScalingY: fabricObject.lockScalingY,
        lockMovementX: fabricObject.lockMovementX,
        lockMovementY: fabricObject.lockMovementY,
        lockRotation: fabricObject.lockRotation,
        evented: fabricObject.evented,
        editMode: boardsStore.isEditMode,
        cornerSize: fabricObject.cornerSize,
        transparentCorners: fabricObject.transparentCorners
      })
      
      // Add to canvas first
      fabricCanvas.value.add(fabricObject)
      
      // Set z-index
      fabricObject.moveTo(element.z)
      
      // Force enable controls after adding to canvas
      if (boardsStore.isEditMode) {
        fabricObject.selectable = true
        fabricObject.evented = true
        fabricObject.hasControls = true
        fabricObject.hasBorders = true
        fabricObject.hasRotatingPoint = true
        fabricObject.lockUniScaling = false
        fabricObject.lockScalingFlip = false
        fabricObject.lockMovementX = false
        fabricObject.lockMovementY = false
        fabricObject.lockRotation = false
        fabricObject.lockScalingX = false
        fabricObject.lockScalingY = false
        
        // Ensure corner controls are visible and interactive
        fabricObject.setControlsVisibility({
          mtr: true, // rotation control
          tr: true,  // top-right
          br: true,  // bottom-right
          bl: true,  // bottom-left
          tl: true,  // top-left
          mt: true,  // middle-top
          mb: true,  // middle-bottom
          ml: true,  // middle-left
          mr: true   // middle-right
        })
        
        // CRITICAL: Ensure corner controls are set up for scaling with proper action handlers
        const corners = ['tl', 'tr', 'br', 'bl']
        corners.forEach(cornerKey => {
          if (fabricObject && fabricObject.controls && fabricObject.controls[cornerKey]) {
            fabricObject.controls[cornerKey].actionName = 'scale'
            // Use the default scaling action handler from Fabric's prototype
            const defaultControl = fabric.Object.prototype.controls[cornerKey]
            if (defaultControl && defaultControl.actionHandler) {
              fabricObject.controls[cornerKey].actionHandler = defaultControl.actionHandler
            }
            fabricObject.controls[cornerKey].cursorStyle = 'crosshair'
          }
        })
        
        // Also ensure middle controls work for scaling with proper action handlers
        const middles = ['mt', 'mb', 'ml', 'mr']
        middles.forEach(middleKey => {
          if (fabricObject && fabricObject.controls && fabricObject.controls[middleKey]) {
            fabricObject.controls[middleKey].actionName = 'scale'
            // Use the default scaling action handler from Fabric's prototype
            const defaultControl = fabric.Object.prototype.controls[middleKey]
            if (defaultControl && defaultControl.actionHandler) {
              fabricObject.controls[middleKey].actionHandler = defaultControl.actionHandler
            }
          }
        })
        
        // Update coordinates to ensure proper interaction
        fabricObject.setCoords()
        
        // Force canvas to recognize the object for interactions
        fabricCanvas.value.renderAll()
        
        // Additional debugging
        console.log('Post-setup object state:', {
          elementId: element.id,
          type: fabricObject.type,
          selectable: fabricObject.selectable,
          evented: fabricObject.evented,
          hasControls: fabricObject.hasControls,
          controlsVisibility: (fabricObject as any)._getControlsVisibility ? (fabricObject as any)._getControlsVisibility() : 'method not available',
          canvas: !!fabricObject.canvas,
          canvasObjectCount: fabricCanvas.value.getObjects().length
        })
      }
    }
  } catch (error) {
    console.error(`Failed to create ${element.kind} element:`, error)
  }
  
  return fabricObject
}

// Element creation methods
const createTextObject = async (element: Element): Promise<fabric.IText> => {
  const payload = element.payload as TextPayload
  
  const textObject = new fabric.IText(payload.content, {
    left: element.x,
    top: element.y,
    width: element.w,
    height: element.h,
    fontFamily: payload.fontFamily,
    fontSize: payload.fontSize,
    fill: payload.color,
    fontWeight: payload.bold ? 'bold' : 'normal',
    fontStyle: payload.italic ? 'italic' : 'normal',
    textAlign: payload.textAlign,
    // Visual properties for controls
    transparentCorners: false,
    cornerColor: '#2563EB',
    cornerStrokeColor: '#1E40AF',
    borderColor: '#2563EB',
    cornerSize: 10,
    // Enable editing for text objects
    editable: boardsStore.isEditMode,
    // Explicitly enable all controls and interactions
    hasControls: true,
    hasBorders: true,
    hasRotatingPoint: true,
    selectable: true,
    evented: true,
    lockScalingX: false,
    lockScalingY: false,
    lockMovementX: false,
    lockMovementY: false,
    lockRotation: false,
    lockUniScaling: false,
  })
  
  return textObject
}

const createImageObject = async (element: Element): Promise<fabric.Image> => {
  const payload = element.payload as ImagePayload
  
  return new Promise((resolve, reject) => {
    fabric.Image.fromURL(payload.url, (img: fabric.Image | null) => {
      if (!img) {
        reject(new Error('Failed to load image'))
        return
      }
      
      img.scaleToWidth(element.w)
      img.scaleToHeight(element.h)
      
      // Visual properties for controls
      img.transparentCorners = false
      img.cornerColor = '#2563EB'
      img.cornerStrokeColor = '#1E40AF'
      img.borderColor = '#2563EB'
      img.cornerSize = 10
      
      // Explicitly enable all controls and interactions
      img.hasControls = true
      img.hasBorders = true
      img.hasRotatingPoint = true
      img.selectable = true
      img.evented = true
      img.lockScalingX = false
      img.lockScalingY = false
      img.lockMovementX = false
      img.lockMovementY = false
      img.lockRotation = false
      img.lockUniScaling = false
      
      resolve(img)
    }, { crossOrigin: 'anonymous' })
  })
}

const createShapeObject = async (element: Element): Promise<fabric.Object> => {
  const payload = element.payload as ShapePayload
  
  let shapeObject: fabric.Object
  
  switch (payload.shapeType) {
    case 'rectangle':
      shapeObject = new fabric.Rect({
        width: element.w,
        height: element.h,
        fill: payload.fill,
        stroke: payload.stroke,
        strokeWidth: payload.strokeWidth,
      })
      break
    case 'circle':
      shapeObject = new fabric.Circle({
        radius: Math.min(element.w, element.h) / 2,
        fill: payload.fill,
        stroke: payload.stroke,
        strokeWidth: payload.strokeWidth,
      })
      break
    case 'triangle':
      shapeObject = new fabric.Triangle({
        width: element.w,
        height: element.h,
        fill: payload.fill,
        stroke: payload.stroke,
        strokeWidth: payload.strokeWidth,
      })
      break
    default:
      throw new Error(`Unknown shape type: ${payload.shapeType}`)
  }
  
  // Visual properties for controls
  shapeObject.transparentCorners = false
  shapeObject.cornerColor = '#2563EB'
  shapeObject.cornerStrokeColor = '#1E40AF'
  shapeObject.borderColor = '#2563EB'
  shapeObject.cornerSize = 10
  
  // Explicitly enable all controls and interactions
  shapeObject.hasControls = true
  shapeObject.hasBorders = true
  shapeObject.hasRotatingPoint = true
  shapeObject.selectable = true
  shapeObject.evented = true
  shapeObject.lockScalingX = false
  shapeObject.lockScalingY = false
  shapeObject.lockMovementX = false
  shapeObject.lockMovementY = false
  shapeObject.lockRotation = false
  shapeObject.lockUniScaling = false
  
  return shapeObject
}

const createStickerObject = async (element: Element): Promise<fabric.Image> => {
  const payload = element.payload as StickerPayload
  
  return new Promise((resolve, reject) => {
    fabric.Image.fromURL(payload.url, (img: fabric.Image | null) => {
      if (!img) {
        reject(new Error('Failed to load sticker'))
        return
      }
      
      img.scaleToWidth(element.w)
      img.scaleToHeight(element.h)
      
      // Visual properties for controls
      img.transparentCorners = false
      img.cornerColor = '#2563EB'
      img.cornerStrokeColor = '#1E40AF'
      img.borderColor = '#2563EB'
      img.cornerSize = 8
      
      resolve(img)
    }, { crossOrigin: 'anonymous' })
  })
}

// Element update methods
const updateElementFromFabricObject = (elementId: string, fabricObject: fabric.Object) => {
  const element = elements.value.find(el => el.id === elementId)
  if (!element) return
  
  // Calculate actual dimensions considering scaling
  const width = (fabricObject.width || 0) * (fabricObject.scaleX || 1)
  const height = (fabricObject.height || 0) * (fabricObject.scaleY || 1)
  
  const updates: any = {
    x: fabricObject.left || 0,
    y: fabricObject.top || 0,
    w: width,
    h: height,
    rotation: fabricObject.angle || 0,
  }
  
  // Handle text content updates for IText objects
  if (fabricObject.type === 'i-text' && element.kind === 'text') {
    const textObj = fabricObject as any
    const updatedPayload = {
      ...element.payload,
      content: textObj.text || ''
    }
    updates.payload = updatedPayload
  }
  
  // Update element in store (this will trigger debounced API save)
  editorStore.updateElement(elementId, updates)
}

// Public methods for element creation
const addTextElement = async (x = 100, y = 100) => {
  if (!boardsStore.isEditMode) {
    console.warn('Cannot add text element: not in edit mode')
    return
  }
  
  // If no page is selected, automatically select the first page
  if (!editorStore.currentPageId && boardsStore.sortedPages.length > 0) {
    console.log('No page selected, selecting first page')
    const firstPage = boardsStore.sortedPages[0]
    editorStore.setCurrentPage(firstPage.id)
    await editorStore.loadPageElements(firstPage.id)
  }
  
  if (!editorStore.currentPageId) {
    console.error('Cannot add text element: no current page selected and no pages available')
    return
  }
  
  const payload: TextPayload = {
    content: 'New Text',
    fontFamily: 'Arial',
    fontSize: 16,
    color: '#000000',
    bold: false,
    italic: false,
    textAlign: 'left',
  }
  
  try {
    console.log('Creating text element on page:', editorStore.currentPageId)
    const element = await editorStore.createElement('text', x, y, 120, 24, payload)
    
    if (!element) {
      console.error('Element creation returned undefined')
      return
    }
    
    await addElementToCanvas(element)
    fabricCanvas.value?.renderAll()
    
    // Automatically select the new element and ensure controls are working
    setTimeout(() => {
      console.log('Selecting and enabling controls for new text element:', element.id)
      const allObjects = fabricCanvas.value?.getObjects() || []
      const textObject = allObjects.find((obj: any) => obj.elementId === element.id) as any
      
      if (textObject && fabricCanvas.value) {
        // Use our enhanced control fixing function
        fixObjectControls(textObject)
        
        // Additional aggressive activation for new elements
        fabricCanvas.value.setActiveObject(textObject)
        
        // Force a double render to ensure everything is properly activated
        fabricCanvas.value.renderAll()
        setTimeout(() => {
          if (fabricCanvas.value && textObject) {
            textObject.setCoords()
            fabricCanvas.value.renderAll()
          }
        }, 10)
        
        console.log('Controls enabled for new text element')
        
        // Automatically enter edit mode for text elements
        setTimeout(() => {
          try {
            textObject.enterEditing()
            textObject.selectAll()
            fabricCanvas.value?.renderAll()
            console.log('Successfully entered edit mode')
          } catch (error) {
            console.error('Failed to enter edit mode:', error)
          }
        }, 50)
      } else {
        console.log('Text object not found')
      }
    }, 100)
    
    console.log('Text element created successfully:', element.id)
  } catch (error) {
    console.error('Failed to add text element:', error)
    return
  }
}

const addImageElement = async (imageUrl: string, x = 100, y = 100, originalWidth = 200, originalHeight = 200) => {
  if (!boardsStore.isEditMode) {
    console.warn('Cannot add image element: not in edit mode')
    return
  }
  
  // If no page is selected, automatically select the first page
  if (!editorStore.currentPageId && boardsStore.sortedPages.length > 0) {
    console.log('No page selected, selecting first page')
    const firstPage = boardsStore.sortedPages[0]
    editorStore.setCurrentPage(firstPage.id)
    await editorStore.loadPageElements(firstPage.id)
  }
  
  if (!editorStore.currentPageId) {
    console.error('Cannot add image element: no current page selected and no pages available')
    return
  }
  
  const payload: ImagePayload = {
    url: imageUrl,
    originalWidth,
    originalHeight,
    description: '', // Initialize with empty description
  }
  
  try {
    console.log('Creating image element on page:', editorStore.currentPageId)
    const element = await editorStore.createElement('image', x, y, originalWidth, originalHeight, payload)
    
    if (!element) {
      console.error('Image element creation returned undefined')
      return
    }
    
    await addElementToCanvas(element)
    fabricCanvas.value?.renderAll()
    
    // Automatically select the new element and ensure controls are working
    setTimeout(() => {
      console.log('Selecting and enabling controls for new image element:', element.id)
      const allObjects = fabricCanvas.value?.getObjects() || []
      const imageObject = allObjects.find((obj: any) => obj.elementId === element.id) as any
      
      if (imageObject && fabricCanvas.value) {
        // Use our enhanced control fixing function
        fixObjectControls(imageObject)
        
        // Additional aggressive activation for new elements
        fabricCanvas.value.setActiveObject(imageObject)
        
        // Force a double render to ensure everything is properly activated
        fabricCanvas.value.renderAll()
        setTimeout(() => {
          if (fabricCanvas.value && imageObject) {
            imageObject.setCoords()
            fabricCanvas.value.renderAll()
          }
        }, 10)
        
        console.log('Controls enabled for new image element')
      } else {
        console.log('Image object not found')
      }
    }, 100)
    
    console.log('Image element created successfully:', element.id)
  } catch (error) {
    console.error('Failed to add image element:', error)
    return
  }
}

const addShapeElement = async (shapeType: 'rectangle' | 'circle' | 'triangle' = 'rectangle', x = 100, y = 100) => {
  if (!boardsStore.isEditMode) {
    console.warn('Cannot add shape element: not in edit mode')
    return
  }
  
  // If no page is selected, automatically select the first page
  if (!editorStore.currentPageId && boardsStore.sortedPages.length > 0) {
    console.log('No page selected, selecting first page')
    const firstPage = boardsStore.sortedPages[0]
    editorStore.setCurrentPage(firstPage.id)
    await editorStore.loadPageElements(firstPage.id)
  }
  
  if (!editorStore.currentPageId) {
    console.error('Cannot add shape element: no current page selected and no pages available')
    return
  }
  
  const payload: ShapePayload = {
    shapeType,
    fill: '#3B82F6',
    stroke: '#1E40AF',
    strokeWidth: 2,
  }
  
  const size = shapeType === 'circle' ? 100 : 120
  
  try {
    console.log('Creating shape element on page:', editorStore.currentPageId)
    const element = await editorStore.createElement('shape', x, y, size, size, payload)
    
    if (!element) {
      console.error('Shape element creation returned undefined')
      return
    }
    
    await addElementToCanvas(element)
    fabricCanvas.value?.renderAll()
    
    // Automatically select the new element and ensure controls are working
    setTimeout(() => {
      console.log('Selecting and enabling controls for new shape element:', element.id)
      const allObjects = fabricCanvas.value?.getObjects() || []
      const shapeObject = allObjects.find((obj: any) => obj.elementId === element.id) as any
      
      if (shapeObject && fabricCanvas.value) {
        // Use our enhanced control fixing function
        fixObjectControls(shapeObject)
        
        // Additional aggressive activation for new elements
        fabricCanvas.value.setActiveObject(shapeObject)
        
        // Force a double render to ensure everything is properly activated
        fabricCanvas.value.renderAll()
        setTimeout(() => {
          if (fabricCanvas.value && shapeObject) {
            shapeObject.setCoords()
            fabricCanvas.value.renderAll()
          }
        }, 10)
        
        console.log('Controls enabled for new shape element')
      } else {
        console.log('Shape object not found')
      }
    }, 100)
    
    console.log('Shape element created successfully:', element.id)
  } catch (error) {
    console.error('Failed to add shape element:', error)
    return
  }
}

const addStickerElement = async (stickerUrl: string, stickerType: string, category: string, x = 100, y = 100) => {
  if (!boardsStore.isEditMode) {
    console.warn('Cannot add sticker element: not in edit mode')
    return
  }
  
  // If no page is selected, automatically select the first page
  if (!editorStore.currentPageId && boardsStore.sortedPages.length > 0) {
    console.log('No page selected, selecting first page')
    const firstPage = boardsStore.sortedPages[0]
    editorStore.setCurrentPage(firstPage.id)
    await editorStore.loadPageElements(firstPage.id)
  }
  
  if (!editorStore.currentPageId) {
    console.error('Cannot add sticker element: no current page selected and no pages available')
    return
  }
  
  const payload: StickerPayload = {
    stickerType,
    url: stickerUrl,
    category,
  }
  
  try {
    console.log('Creating sticker element on page:', editorStore.currentPageId)
    const element = await editorStore.createElement('sticker', x, y, 80, 80, payload)
    
    if (!element) {
      console.error('Sticker element creation returned undefined')
      return
    }
    
    await addElementToCanvas(element)
    fabricCanvas.value?.renderAll()
    
    // Automatically select the new element and ensure controls are working
    setTimeout(() => {
      console.log('Selecting and enabling controls for new sticker element:', element.id)
      const allObjects = fabricCanvas.value?.getObjects() || []
      const stickerObject = allObjects.find((obj: any) => obj.elementId === element.id) as any
      
      if (stickerObject && fabricCanvas.value) {
        // Use our enhanced control fixing function
        fixObjectControls(stickerObject)
        
        // Additional aggressive activation for new elements
        fabricCanvas.value.setActiveObject(stickerObject)
        
        // Force a double render to ensure everything is properly activated
        fabricCanvas.value.renderAll()
        setTimeout(() => {
          if (fabricCanvas.value && stickerObject) {
            stickerObject.setCoords()
            fabricCanvas.value.renderAll()
          }
        }, 10)
        
        console.log('Controls enabled for new sticker element')
      } else {
        console.log('Sticker object not found')
      }
    }, 100)
    
    console.log('Sticker element created successfully:', element.id)
  } catch (error) {
    console.error('Failed to add sticker element:', error)
    return
  }
}

// Selection and deletion methods
const selectAllElements = () => {
  if (!fabricCanvas.value) return
  
  const allObjects = fabricCanvas.value.getObjects()
  if (allObjects.length > 0) {
    const selection = new fabric.ActiveSelection(allObjects, {
      canvas: fabricCanvas.value,
    })
    fabricCanvas.value.setActiveObject(selection)
    fabricCanvas.value.renderAll()
  }
}

const deleteSelectedElements = async () => {
  if (selectedElementIds.value.length === 0) return
  
  try {
    // Remove from canvas first for immediate feedback
    const activeObjects = fabricCanvas.value?.getActiveObjects() || []
    activeObjects.forEach((obj: fabric.Object) => {
      fabricCanvas.value?.remove(obj)
    })
    fabricCanvas.value?.discardActiveObject()
    fabricCanvas.value?.renderAll()
    
    // Delete from backend and store
    await editorStore.deleteElements(selectedElementIds.value)
  } catch (error) {
    console.error('Failed to delete elements:', error)
    // Reload elements to restore state on error
    await loadElementsToCanvas()
  }
}

// Utility methods
const clearError = () => {
  editorStore.clearError()
}

// Utility function to fix controls on an object
const fixObjectControls = (obj: any) => {
  if (!obj || !boardsStore.isEditMode) return
  
  console.log('Fixing controls for object:', obj.elementId)
  
  // 1. Exit IText editing mode if active (this can hide controls)
  if (obj.type === 'i-text' && obj.isEditing) {
    console.log('Exiting IText editing mode to show controls')
    obj.exitEditing()
  }
  
  // 2. Force enable all controls and interactions
  obj.hasControls = true
  obj.hasBorders = true
  obj.hasRotatingPoint = true
  obj.lockScalingX = false
  obj.lockScalingY = false
  obj.lockMovementX = false
  obj.lockMovementY = false
  obj.lockRotation = false
  obj.lockUniScaling = false
  obj.lockScalingFlip = false
  obj.selectable = true
  obj.evented = true
  obj.moveCursor = 'move'
  obj.hoverCursor = 'move'
  
  // 3. Ensure canvas-level settings don't interfere
  if (fabricCanvas.value) {
    fabricCanvas.value.selection = true
    
    // 4. Check and reset zoom/viewport if they're causing issues
    const currentZoom = fabricCanvas.value.getZoom()
    const viewport = fabricCanvas.value.viewportTransform
    console.log('Canvas state check:', {
      zoom: currentZoom,
      viewport: viewport,
      selection: fabricCanvas.value.selection
    })
  }
  
  // 5. Restore default controls if they were overridden/hidden
  if (!obj.controls || Object.keys(obj.controls).length === 0) {
    console.log('Restoring default controls for object')
    Object.assign(obj.controls, fabric.Object.prototype.controls)
  }
  
  // 6. Ensure prototype hasn't been modified globally
  if (!fabric.Object.prototype.hasControls || !fabric.Object.prototype.hasBorders) {
    console.log('Restoring global prototype controls')
    fabric.Object.prototype.hasControls = true
    fabric.Object.prototype.hasBorders = true
  }
  
  // 7. Set visual properties for better control visibility
  obj.cornerSize = 14
  obj.transparentCorners = false
  obj.cornerColor = '#2563EB'
  obj.borderColor = '#2563EB'
  obj.cornerStrokeColor = '#1E40AF'
  obj.borderDashArray = null
  obj.borderOpacityWhenMoving = 0.4
  
  // 8. Explicitly set all controls to be visible and ensure they're interactive
  obj.setControlsVisibility({
    mtr: true, // rotation control
    tr: true,  // top-right
    br: true,  // bottom-right
    bl: true,  // bottom-left
    tl: true,  // top-left
    mt: true,  // middle-top
    mb: true,  // middle-bottom
    ml: true,  // middle-left
    mr: true   // middle-right
  })
  
  // 9. CRITICAL: Ensure corner controls are set up for scaling, not just moving
  // Force corner controls to be scaling controls with proper action handlers
  const corners = ['tl', 'tr', 'br', 'bl']
  corners.forEach(cornerKey => {
    if (obj.controls && obj.controls[cornerKey]) {
      // Ensure corner controls handle scaling
      obj.controls[cornerKey].actionName = 'scale'
      // Use the default scaling action handler from Fabric's prototype
      const defaultControl = fabric.Object.prototype.controls[cornerKey]
      if (defaultControl && defaultControl.actionHandler) {
        obj.controls[cornerKey].actionHandler = defaultControl.actionHandler
      }
      // Make sure corner controls have the scaling cursor
      obj.controls[cornerKey].cursorStyle = 'crosshair'
      obj.controls[cornerKey].cursorStyleHandler = () => 'crosshair'
    }
  })
  
  // Also ensure middle controls are set up properly with action handlers
  const middles = ['mt', 'mb', 'ml', 'mr']
  middles.forEach(middleKey => {
    if (obj.controls && obj.controls[middleKey]) {
      obj.controls[middleKey].actionName = 'scale'
      // Use the default scaling action handler from Fabric's prototype
      const defaultControl = fabric.Object.prototype.controls[middleKey]
      if (defaultControl && defaultControl.actionHandler) {
        obj.controls[middleKey].actionHandler = defaultControl.actionHandler
      }
      if (middleKey === 'mt' || middleKey === 'mb') {
        obj.controls[middleKey].cursorStyle = 'ns-resize'
        obj.controls[middleKey].cursorStyleHandler = () => 'ns-resize'
      } else {
        obj.controls[middleKey].cursorStyle = 'ew-resize'
        obj.controls[middleKey].cursorStyleHandler = () => 'ew-resize'
      }
    }
  })
  
  // 10. Force Fabric.js to recognize this as a scalable object
  obj.lockScalingX = false
  obj.lockScalingY = false
  obj.lockUniScaling = false // Allow non-uniform scaling
  
  // 11. Set scaling behavior
  obj.uniformScaling = false // Allow free scaling
  obj.centeredScaling = false // Scale from corners, not center
  obj.centeredRotation = true // But allow centered rotation
  
  // 9. Force refresh of control coordinates and rendering
  obj.setCoords()
  
  // 10. Make sure the object is properly initialized for interactions
  if (fabricCanvas.value) {
    // Force canvas to recognize this object's interactive state
    fabricCanvas.value.setActiveObject(obj)
    
    // Trigger a re-render to ensure controls are visible and interactive
    fabricCanvas.value.renderAll()
    
    // Force update the object's interactive elements with multiple cycles
    setTimeout(() => {
      if (obj && fabricCanvas.value) {
        obj.setCoords()
        fabricCanvas.value.renderAll()
        
        // Additional aggressive re-activation
        setTimeout(() => {
          if (obj && fabricCanvas.value) {
            // Re-apply all critical properties
            obj.hasControls = true
            obj.hasBorders = true
            obj.selectable = true
            obj.evented = true
            obj.setCoords()
            fabricCanvas.value.renderAll()
          }
        }, 10)
      }
    }, 0)
  }
  
  console.log('Controls fixed for object:', obj.elementId, {
    hasControls: obj.hasControls,
    hasBorders: obj.hasBorders,
    selectable: obj.selectable,
    evented: obj.evented,
    lockScalingX: obj.lockScalingX,
    lockScalingY: obj.lockScalingY,
    cornerSize: obj.cornerSize,
    isEditing: obj.isEditing || false,
    controlsCount: Object.keys(obj.controls || {}).length,
    canvasSelection: fabricCanvas.value?.selection,
    zoom: fabricCanvas.value?.getZoom()
  })
}

const resizeCanvas = () => {
  if (!fabricCanvas.value || !canvasContainer.value) return
  
  const container = canvasContainer.value
  const rect = container.getBoundingClientRect()
  
  fabricCanvas.value.setDimensions({
    width: rect.width,
    height: rect.height,
  })
  
  fabricCanvas.value.renderAll()
}

// Watchers
watch(() => elements.value, async () => {
  // Only reload if we're not in the middle of text editing, user interactions, or other operations
  if (isTextEditing.value) {
    console.log('Skipping element reload during text editing')
    return
  }
  
  if (isUserInteracting.value) {
    console.log('Skipping element reload during user interaction')
    return
  }
  
  // Check if there's an active object being manipulated
  const activeObject = fabricCanvas.value?.getActiveObject()
  if (activeObject && (activeObject as any).__isBeingTransformed) {
    console.log('Skipping element reload during object transformation')
    return
  }
  
  // Clear any existing reload timer
  if (reloadDebounceTimer.value) {
    clearTimeout(reloadDebounceTimer.value)
  }
  
  // Debounce the reload to prevent excessive reloading
  reloadDebounceTimer.value = setTimeout(async () => {
    // Double-check that user is not interacting before reloading
    if (!isTextEditing.value && !isUserInteracting.value) {
      console.log('Elements changed, reloading canvas')
      await nextTick()
      await loadElementsToCanvas()
    } else {
      console.log('Skipping scheduled element reload due to ongoing user interaction')
    }
  }, 200) // Increased from 150ms to 200ms to give more time for interactions to complete
}, { deep: false }) // Changed from deep: true to deep: false

watch(() => editorStore.currentPageId, async () => {
  console.log('Page changed, reloading canvas')
  await nextTick()
  await loadElementsToCanvas()
})

// Expose methods to parent component
defineExpose({
  addTextElement,
  addImageElement,
  addShapeElement,
  addStickerElement,
  selectAllElements,
  deleteSelectedElements,
  resizeCanvas,
  refreshCanvas: () => {
    console.log('Manual canvas refresh triggered')
    loadElementsToCanvas()
  },
  // Debug methods for testing
  testScaling: () => {
    const activeObject = fabricCanvas.value?.getActiveObject()
    if (activeObject) {
      console.log('Testing manual scaling on active object')
      console.log('Current object properties:', {
        scaleX: activeObject.scaleX,
        scaleY: activeObject.scaleY,
        hasControls: activeObject.hasControls,
        lockScalingX: activeObject.lockScalingX,
        lockScalingY: activeObject.lockScalingY
      })
      
      // Try to manually scale the object
      activeObject.scaleX = (activeObject.scaleX || 1) * 1.2
      activeObject.scaleY = (activeObject.scaleY || 1) * 1.2
      activeObject.setCoords()
      fabricCanvas.value?.renderAll()
      
      console.log('After manual scaling:', {
        scaleX: activeObject.scaleX,
        scaleY: activeObject.scaleY
      })
    } else {
      console.log('No active object to test scaling')
    }
  },
  getCanvasInfo: () => {
    if (fabricCanvas.value) {
      const activeObject = fabricCanvas.value.getActiveObject()
      console.log('Canvas info:', {
        interactive: fabricCanvas.value.interactive,
        selection: fabricCanvas.value.selection,
        allowTouchScrolling: fabricCanvas.value.allowTouchScrolling,
        objects: fabricCanvas.value.getObjects().length,
        activeObject: !!activeObject,
        width: fabricCanvas.value.width,
        height: fabricCanvas.value.height
      })
      
      if (activeObject) {
        console.log('Active object detailed info:', {
          type: activeObject.type,
          hasControls: activeObject.hasControls,
          hasBorders: activeObject.hasBorders,
          hasRotatingPoint: activeObject.hasRotatingPoint,
          selectable: activeObject.selectable,
          evented: activeObject.evented,
          lockScalingX: activeObject.lockScalingX,
          lockScalingY: activeObject.lockScalingY,
          lockMovementX: activeObject.lockMovementX,
          lockMovementY: activeObject.lockMovementY,
          lockRotation: activeObject.lockRotation,
          cornerSize: activeObject.cornerSize,
          transparentCorners: activeObject.transparentCorners,
          cornerColor: activeObject.cornerColor,
          borderColor: activeObject.borderColor,
          controls: Object.keys((activeObject as any).controls || {}),
          controlsVisible: (activeObject as any)._controlsVisibility || 'no visibility info'
        })
        
        // Check each control specifically
        const obj = activeObject as any
        if (obj.controls) {
          console.log('Individual control visibility:')
          Object.keys(obj.controls).forEach(key => {
            const control = obj.controls[key]
            console.log(`  ${key}:`, {
              visible: control.visible !== false,
              x: control.x,
              y: control.y,
              offsetX: control.offsetX,
              offsetY: control.offsetY
            })
          })
        }
      }
    }
  },
  forceEnableControls: () => {
    const activeObject = fabricCanvas.value?.getActiveObject()
    if (activeObject && boardsStore.isEditMode) {
      console.log('Force enabling all controls...')
      fixObjectControls(activeObject)
      fabricCanvas.value?.renderAll()
    } else {
      console.log('No active object or not in edit mode')
    }
  },
  fixAllControls: () => {
    if (!fabricCanvas.value || !boardsStore.isEditMode) {
      console.log('Canvas not available or not in edit mode')
      return
    }
    
    console.log('Fixing controls for all objects...')
    const allObjects = fabricCanvas.value.getObjects()
    allObjects.forEach(obj => {
      fixObjectControls(obj)
    })
    fabricCanvas.value.renderAll()
    console.log(`Fixed controls for ${allObjects.length} objects`)
  },
  checkControlsState: () => {
    const activeObject = fabricCanvas.value?.getActiveObject()
    if (activeObject) {
      const obj = activeObject as any
      console.log('Current controls state:', {
        hasControls: obj.hasControls,
        hasBorders: obj.hasBorders,
        selectable: obj.selectable,
        evented: obj.evented,
        lockScalingX: obj.lockScalingX,
        lockScalingY: obj.lockScalingY,
        cornerSize: obj.cornerSize,
        controls: Object.keys(obj.controls || {})
      })
      
      // Check if controls are properly attached
      if (obj.controls) {
        console.log('Control details:')
        Object.entries(obj.controls).forEach(([key, control]: [string, any]) => {
          console.log(`  ${key}:`, {
            visible: control.visible,
            render: typeof control.render,
            actionHandler: typeof control.actionHandler,
            mouseUpHandler: typeof control.mouseUpHandler
          })
        })
      } else {
        console.log('No controls object found!')
      }
    } else {
      console.log('No active object')
    }
  },
  repairControls: () => {
    const activeObject = fabricCanvas.value?.getActiveObject()
    if (activeObject && boardsStore.isEditMode) {
      console.log('Attempting to repair controls...')
      
      // Remove and re-add the object to reset its controls
      if (fabricCanvas.value) {
        const objects = fabricCanvas.value.getObjects()
        const objectIndex = objects.indexOf(activeObject)
        
        if (objectIndex >= 0) {
          // Store object properties
          const objData = activeObject.toObject()
          const elementId = (activeObject as any).elementId
          
          // Remove the object
          fabricCanvas.value.remove(activeObject)
          
          // Recreate it with proper controls
          fabric.util.enlivenObjects([objData], (enlivenedObjects: fabric.Object[]) => {
            const newObject = enlivenedObjects[0]
            
            // Restore element ID and control properties
            ;(newObject as any).elementId = elementId
            newObject.hasControls = true
            newObject.hasBorders = true
            newObject.hasRotatingPoint = true
            newObject.selectable = true
            newObject.evented = true
            newObject.lockScalingX = false
            newObject.lockScalingY = false
            newObject.cornerSize = 12
            newObject.cornerColor = '#FF0000'
            newObject.transparentCorners = false
            
            // Set controls visibility
            newObject.setControlsVisibility({
              mtr: true, tr: true, br: true, bl: true, tl: true,
              mt: true, mb: true, ml: true, mr: true
            })
            
            // Add back to canvas
            if (fabricCanvas.value) {
              fabricCanvas.value.add(newObject)
              fabricCanvas.value.setActiveObject(newObject)
              newObject.setCoords()
              fabricCanvas.value.renderAll()
              
              console.log('Object recreated with fresh controls')
            }
          }, '')
        }
      }
    }
  }
})

// Lifecycle
onMounted(async () => {
  await nextTick()
  await initializeCanvas()
  
  // Handle window resize
  window.addEventListener('resize', resizeCanvas)
  
  // Listen for editor reload events
  const handleReloadElements = () => {
    console.log('Reloading elements from undo/redo')
    // Clear canvas and reload all elements
    if (fabricCanvas.value && fabricCanvas.value.getContext()) {
      try {
        fabricCanvas.value.clear()
        editorStore.elements.forEach(element => {
          addElementToCanvas(element)
        })
        fabricCanvas.value.renderAll()
      } catch (error) {
        console.warn('Canvas reload failed, canvas may not be ready:', error)
      }
    } else {
      console.warn('Canvas not ready for reload')
    }
  }
  
  window.addEventListener('editor:reload-elements', handleReloadElements)
})

onUnmounted(() => {
  // Clean up event listeners
  document.removeEventListener('keydown', handleKeyDown)
  window.removeEventListener('resize', resizeCanvas)
  window.removeEventListener('editor:reload-elements', () => {})
  
  // Dispose of fabric canvas
  if (fabricCanvas.value) {
    fabricCanvas.value.dispose()
  }
})
</script>

<style scoped>
.canvas-editor-container {
  position: relative;
  background: transparent;
}

.canvas-wrapper {
  position: relative;
}

.fabric-canvas {
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  /* Ensure no pointer events interference */
  pointer-events: auto;
}

.grid-overlay {
  opacity: 0.3;
  pointer-events: none;
}

.snap-to-grid .fabric-canvas {
  cursor: crosshair;
}

/* Fabric.js custom styles */
:deep(.canvas-container) {
  margin: 0 !important;
  /* Ensure no pointer events interference */
  pointer-events: auto;
}

:deep(.upper-canvas) {
  border-radius: 8px;
  /* Critical: ensure canvas can receive all mouse events */
  pointer-events: auto !important;
  touch-action: none !important;
}

/* Ensure fabric controls are properly styled and interactive */
:deep(.fabric-control) {
  pointer-events: auto !important;
}

:deep(.fabric-control-corner) {
  pointer-events: auto !important;
}
</style>