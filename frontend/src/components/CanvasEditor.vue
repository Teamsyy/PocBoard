<template>
  <div class="canvas-editor-container relative w-full h-full">
    <!-- Canvas Container -->
    <div 
      ref="canvasContainer" 
      class="canvas-wrapper relative w-full h-full overflow-hidden"
      :class="{ 'snap-to-grid': snapToGrid }"
    >
      <canvas 
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, computed, nextTick } from 'vue'
import { fabric } from 'fabric'
import { useEditorStore } from '@/stores/editor'
import { useBoardsStore } from '@/stores/boards'
import LoadingSpinner from './LoadingSpinner.vue'
import ErrorMessage from './ErrorMessage.vue'
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
const initializeCanvas = () => {
  if (!canvasElement.value) return

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
  })

  // Store canvas reference in editor store
  editorStore.setCanvas(fabricCanvas.value)

  // Set up event listeners
  setupCanvasEvents()
  
  // Load existing elements
  loadElementsToCanvas()
}

// Canvas event handlers
const setupCanvasEvents = () => {
  if (!fabricCanvas.value) return

  const canvas = fabricCanvas.value

  // Selection events
  canvas.on('selection:created', handleSelectionChange)
  canvas.on('selection:updated', handleSelectionChange)
  canvas.on('selection:cleared', handleSelectionClear)

  // Object modification events
  canvas.on('object:modified', handleObjectModified)
  canvas.on('object:moving', handleObjectMoving)
  canvas.on('object:scaling', handleObjectScaling)
  canvas.on('object:rotating', handleObjectRotating)

  // Mouse events for element creation
  canvas.on('mouse:down', handleMouseDown)
  canvas.on('mouse:up', handleMouseUp)

  // Keyboard events
  document.addEventListener('keydown', handleKeyDown)
}

// Selection handlers
const handleSelectionChange = (_e: fabric.IEvent) => {
  const activeObjects = fabricCanvas.value?.getActiveObjects() || []
  const elementIds = activeObjects
    .map((obj: any) => obj.elementId)
    .filter(Boolean)
  
  editorStore.selectElements(elementIds)
}

const handleSelectionClear = () => {
  editorStore.clearSelection()
}

// Object modification handlers
const handleObjectModified = (e: fabric.IEvent) => {
  if (!e.target) return
  
  const obj = e.target
  const elementId = (obj as any).elementId
  
  if (elementId) {
    updateElementFromFabricObject(elementId, obj)
  }
}

const handleObjectMoving = (e: fabric.IEvent) => {
  if (!e.target || !snapToGrid.value) return
  
  const obj = e.target
  
  // Snap to grid
  obj.left = Math.round((obj.left || 0) / GRID_SIZE) * GRID_SIZE
  obj.top = Math.round((obj.top || 0) / GRID_SIZE) * GRID_SIZE
}

const handleObjectScaling = (e: fabric.IEvent) => {
  if (!e.target) return
  
  const obj = e.target
  const elementId = (obj as any).elementId
  
  if (elementId && snapToGrid.value) {
    // Snap dimensions to grid
    const width = (obj.width || 0) * (obj.scaleX || 1)
    const height = (obj.height || 0) * (obj.scaleY || 1)
    
    const snappedWidth = Math.max(MIN_ELEMENT_SIZE, Math.round(width / GRID_SIZE) * GRID_SIZE)
    const snappedHeight = Math.max(MIN_ELEMENT_SIZE, Math.round(height / GRID_SIZE) * GRID_SIZE)
    
    obj.scaleX = snappedWidth / (obj.width || 1)
    obj.scaleY = snappedHeight / (obj.height || 1)
  }
}

const handleObjectRotating = (e: fabric.IEvent) => {
  if (!e.target) return
  
  const obj = e.target
  const elementId = (obj as any).elementId
  
  if (elementId) {
    // Debounced update will be handled by handleObjectModified
  }
}

// Mouse event handlers
const handleMouseDown = (e: fabric.IEvent) => {
  // Handle canvas clicks for deselection
  if (!e.target) {
    editorStore.clearSelection()
  }
}

const handleMouseUp = (_e: fabric.IEvent) => {
  // Handle any post-interaction cleanup
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
  
  // Undo/Redo handling
  if (e.ctrlKey || e.metaKey) {
    if (e.key === 'z' && !e.shiftKey) {
      e.preventDefault()
      editorStore.undo()
    } else if ((e.key === 'z' && e.shiftKey) || e.key === 'y') {
      e.preventDefault()
      editorStore.redo()
    }
  }
  
  // Select all
  if ((e.ctrlKey || e.metaKey) && e.key === 'a') {
    e.preventDefault()
    selectAllElements()
  }
}

// Element management
const loadElementsToCanvas = async () => {
  if (!fabricCanvas.value) return
  
  // Clear existing objects
  fabricCanvas.value.clear()
  
  // Load elements from store
  for (const element of elements.value) {
    if (element) {
      await addElementToCanvas(element)
    }
  }
  
  fabricCanvas.value.renderAll()
}

const addElementToCanvas = async (element: Element): Promise<fabric.Object | null> => {
  if (!fabricCanvas.value || !element) return null
  
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
      
      // Add to canvas
      fabricCanvas.value.add(fabricObject)
      
      // Set z-index
      fabricObject.moveTo(element.z)
    }
  } catch (error) {
    console.error(`Failed to create ${element.kind} element:`, error)
  }
  
  return fabricObject
}

// Element creation methods
const createTextObject = async (element: Element): Promise<fabric.Text> => {
  const payload = element.payload as TextPayload
  
  const textObject = new fabric.Text(payload.content, {
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
    selectable: boardsStore.isEditMode,
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
      img.selectable = boardsStore.isEditMode
      
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
  
  shapeObject.selectable = boardsStore.isEditMode
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
      img.selectable = boardsStore.isEditMode
      
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
  
  const updates = {
    x: fabricObject.left || 0,
    y: fabricObject.top || 0,
    w: width,
    h: height,
    rotation: fabricObject.angle || 0,
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
  
  if (!editorStore.currentPageId) {
    console.error('Cannot add text element: no current page selected')
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
    console.log('Text element created successfully:', element.id)
  } catch (error) {
    console.error('Failed to add text element:', error)
    // Return early to prevent further execution
    return
  }
}

const addImageElement = async (imageUrl: string, x = 100, y = 100, originalWidth = 200, originalHeight = 200) => {
  if (!boardsStore.isEditMode) {
    console.warn('Cannot add image element: not in edit mode')
    return
  }
  
  if (!editorStore.currentPageId) {
    console.error('Cannot add image element: no current page selected')
    return
  }
  
  const payload: ImagePayload = {
    url: imageUrl,
    originalWidth,
    originalHeight,
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
  
  if (!editorStore.currentPageId) {
    console.error('Cannot add shape element: no current page selected')
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
  
  if (!editorStore.currentPageId) {
    console.error('Cannot add sticker element: no current page selected')
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
  await nextTick()
  await loadElementsToCanvas()
}, { deep: true })

watch(() => editorStore.currentPageId, async () => {
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
})

// Lifecycle
onMounted(async () => {
  await nextTick()
  initializeCanvas()
  
  // Handle window resize
  window.addEventListener('resize', resizeCanvas)
})

onUnmounted(() => {
  // Clean up event listeners
  document.removeEventListener('keydown', handleKeyDown)
  window.removeEventListener('resize', resizeCanvas)
  
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
}

:deep(.upper-canvas) {
  border-radius: 8px;
}
</style>