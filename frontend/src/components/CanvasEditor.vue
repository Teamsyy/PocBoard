<template>
  <div class="canvas-editor-container relative w-full h-full">
    <!-- Canvas Container -->
    <div 
      ref="canvasContainer" 
      class="canvas-wrapper relative w-full h-full overflow-hidden"
      :class="{ 'snap-to-grid': snapToGrid }"
      @click="handleCanvasClick"
      @keydown="handleKeyDown"
      tabindex="0"
    >
      <!-- Grid Overlay (when snap to grid is enabled) -->
      <div 
        v-if="snapToGrid" 
        class="grid-overlay absolute inset-0 pointer-events-none"
        :style="gridStyle"
      />
      
      <!-- Render Elements -->
      <TextElement
        v-for="element in textElements"
        :key="element.id"
        :element="element"
        :isSelected="selectedElementIds.includes(element.id)"
        :isEditMode="boardsStore.isEditMode"
        :snapToGrid="snapToGrid"
        :gridSize="GRID_SIZE"
        @select="selectElement"
        @deselect="deselectElement"
        @update="updateElement"
        @delete="deleteElement"
      />
      
      <ImageElement
        v-for="element in imageElements"
        :key="element.id"
        :element="element"
        :isSelected="selectedElementIds.includes(element.id)"
        :isEditMode="boardsStore.isEditMode"
        :snapToGrid="snapToGrid"
        :gridSize="GRID_SIZE"
        @select="selectElement"
        @deselect="deselectElement"
        @update="updateElement"
        @delete="deleteElement"
        @openViewer="openImageViewer"
      />
      
      <ShapeElement
        v-for="element in shapeElements"
        :key="element.id"
        :element="element"
        :isSelected="selectedElementIds.includes(element.id)"
        :isEditMode="boardsStore.isEditMode"
        :snapToGrid="snapToGrid"
        :gridSize="GRID_SIZE"
        @select="selectElement"
        @deselect="deselectElement"
        @update="updateElement"
        @delete="deleteElement"
      />
      
      <StickerElement
        v-for="element in stickerElements"
        :key="element.id"
        :element="element"
        :isSelected="selectedElementIds.includes(element.id)"
        :isEditMode="boardsStore.isEditMode"
        :snapToGrid="snapToGrid"
        :gridSize="GRID_SIZE"
        @select="selectElement"
        @deselect="deselectElement"
        @update="updateElement"
        @delete="deleteElement"
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
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useEditorStore } from '@/stores/editor'
import { useBoardsStore } from '@/stores/boards'
import LoadingSpinner from './LoadingSpinner.vue'
import ErrorMessage from './ErrorMessage.vue'
import ImageViewer from './ImageViewer.vue'
import TextElement from './elements/TextElement.vue'
import ImageElement from './elements/ImageElement.vue'
import ShapeElement from './elements/ShapeElement.vue'
import StickerElement from './elements/StickerElement.vue'
import type { Element, TextPayload, ImagePayload, ShapePayload, StickerPayload } from '@/types'

// Props
interface Props {
  width?: number
  height?: number
}

withDefaults(defineProps<Props>(), {
  width: 800,
  height: 600
})

// Stores
const editorStore = useEditorStore()
const boardsStore = useBoardsStore()

// Refs
const canvasContainer = ref<HTMLDivElement>()

// Image viewer state
const imageViewer = ref({
  isVisible: false,
  imageData: null as ImagePayload | null,
  elementId: null as string | null
})

// Computed
const snapToGrid = computed(() => editorStore.snapToGrid)
const loading = computed(() => editorStore.loading)
const error = computed(() => editorStore.error)
const elements = computed(() => editorStore.sortedElements)
const selectedElementIds = computed(() => editorStore.selectedElementIds)

// Element type filters
const textElements = computed(() => elements.value.filter((el: Element) => el.kind === 'text'))
const imageElements = computed(() => elements.value.filter((el: Element) => el.kind === 'image'))
const shapeElements = computed(() => elements.value.filter((el: Element) => el.kind === 'shape'))
const stickerElements = computed(() => elements.value.filter((el: Element) => el.kind === 'sticker'))

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

// Element management methods
const selectElement = (elementId: string) => {
  editorStore.selectElements([elementId])
}

const deselectElement = () => {
  editorStore.clearSelection()
}

const updateElement = (elementId: string, updates: any) => {
  editorStore.updateElement(elementId, updates)
}

const deleteElement = async (elementId: string) => {
  await editorStore.deleteElement(elementId)
}

// Canvas event handlers
const handleCanvasClick = (event: MouseEvent) => {
  // If clicking on canvas background (not on an element), clear selection
  if (event.target === canvasContainer.value) {
    editorStore.clearSelection()
  }
}

const handleKeyDown = (event: KeyboardEvent) => {
  // Handle global keyboard shortcuts
  if (event.ctrlKey || event.metaKey) {
    switch (event.key) {
      case 'z':
        if (!event.shiftKey) {
          event.preventDefault()
          editorStore.undo()
        } else {
          event.preventDefault()
          editorStore.redo()
        }
        break
      case 'y':
        event.preventDefault()
        editorStore.redo()
        break
      case 't':
        if (boardsStore.isEditMode) {
          event.preventDefault()
          addTextElement()
        }
        break
      case 'i':
        if (boardsStore.isEditMode) {
          event.preventDefault()
          const placeholderUrl = 'https://via.placeholder.com/200x200/3B82F6/FFFFFF?text=Image'
          addImageElement(placeholderUrl, 100, 100, 200, 200)
        }
        break
      case 's':
        if (boardsStore.isEditMode) {
          event.preventDefault()
          addShapeElement('rectangle', 100, 100)
        }
        break
      case 'a':
        event.preventDefault()
        selectAllElements()
        break
    }
  }
  
  // Grid toggle
  if (event.key === 'g' && !event.ctrlKey && !event.metaKey) {
    event.preventDefault()
    editorStore.toggleSnapToGrid()
  }
}// Public methods for element creation
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
    
    // Auto-select the new element
    editorStore.selectElements([element.id])
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
    description: '',
  }
  
  try {
    console.log('Creating image element on page:', editorStore.currentPageId)
    const element = await editorStore.createElement('image', x, y, originalWidth, originalHeight, payload)
    
    if (!element) {
      console.error('Image element creation returned undefined')
      return
    }
    
    // Auto-select the new element
    editorStore.selectElements([element.id])
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
    
    // Auto-select the new element
    editorStore.selectElements([element.id])
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
    url: stickerUrl,
    stickerType,
    category,
  }
  
  try {
    console.log('Creating sticker element on page:', editorStore.currentPageId)
    const element = await editorStore.createElement('sticker', x, y, 80, 80, payload)
    
    if (!element) {
      console.error('Sticker element creation returned undefined')
      return
    }
    
    // Auto-select the new element
    editorStore.selectElements([element.id])
    console.log('Sticker element created successfully:', element.id)
  } catch (error) {
    console.error('Failed to add sticker element:', error)
    return
  }
}// Selection and deletion methods
const selectAllElements = () => {
  const allElementIds = elements.value.map((el: Element) => el.id)
  editorStore.selectElements(allElementIds)
}

const deleteSelectedElements = async () => {
  if (selectedElementIds.value.length === 0) return
  
  for (const elementId of selectedElementIds.value) {
    await editorStore.deleteElement(elementId)
  }
  editorStore.clearSelection()
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
  const element = editorStore.elements.find((el: Element) => el.id === elementId)
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

// Utility methods
const clearError = () => {
  editorStore.clearError()
}

// Lifecycle
onMounted(() => {
  // Focus canvas container to enable keyboard events
  if (canvasContainer.value) {
    canvasContainer.value.focus()
  }
  console.log('CanvasEditor mounted with vue-drag-resize')
})

onUnmounted(() => {
  console.log('CanvasEditor unmounted')
})

// Expose public methods for external use
defineExpose({
  addTextElement,
  addImageElement,
  addShapeElement,
  addStickerElement,
  selectAllElements,
  deleteSelectedElements
})
</script>

<style scoped>
.canvas-editor-container {
  position: relative;
  overflow: hidden;
}

.canvas-wrapper {
  position: relative;
  background: #ffffff;
  outline: none;
}

.canvas-wrapper.snap-to-grid {
  cursor: crosshair;
}

.grid-overlay {
  background-size: 8px 8px;
  background-position: 0 0;
  opacity: 0.3;
}

/* Vue Drag Resize global styles */
:deep(.vdr) {
  border: 2px solid transparent;
  position: absolute;
}

:deep(.vdr.active) {
  border-color: #2563EB;
}

:deep(.vdr-handle) {
  position: absolute;
  background: #2563EB;
  border: 1px solid #1E40AF;
  border-radius: 2px;
  box-sizing: border-box;
}

:deep(.vdr-handle-tl) {
  top: -6px;
  left: -6px;
  cursor: nw-resize;
}

:deep(.vdr-handle-tm) {
  top: -6px;
  left: 50%;
  margin-left: -6px;
  cursor: n-resize;
}

:deep(.vdr-handle-tr) {
  top: -6px;
  right: -6px;
  cursor: ne-resize;
}

:deep(.vdr-handle-ml) {
  top: 50%;
  left: -6px;
  margin-top: -6px;
  cursor: w-resize;
}

:deep(.vdr-handle-mr) {
  top: 50%;
  right: -6px;
  margin-top: -6px;
  cursor: e-resize;
}

:deep(.vdr-handle-bl) {
  bottom: -6px;
  left: -6px;
  cursor: sw-resize;
}

:deep(.vdr-handle-bm) {
  bottom: -6px;
  left: 50%;
  margin-left: -6px;
  cursor: s-resize;
}

:deep(.vdr-handle-br) {
  bottom: -6px;
  right: -6px;
  cursor: se-resize;
}

:deep(.vdr-handle) {
  width: 12px;
  height: 12px;
  background: #2563EB;
  border: 2px solid #ffffff;
  border-radius: 2px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

:deep(.vdr-handle:hover) {
  background: #1E40AF;
}

/* Element specific styles */
.text-element,
.image-element,
.shape-element,
.sticker-element {
  outline: none;
}

.text-element.selected,
.image-element.selected,
.shape-element.selected,
.sticker-element.selected {
  outline: 2px solid #2563EB;
  outline-offset: -2px;
}
</style>