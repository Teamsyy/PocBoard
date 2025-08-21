<template>
  <Vue3DraggableResizable
    :initW="element.w"
    :initH="element.h"
    :x="element.x"
    :y="element.y"
    :draggable="isEditMode"
    :resizable="isEditMode"
    :minW="20"
    :minH="20"
    :parent="true"
    @activated="handleActivated"
    @deactivated="handleDeactivated"
    @dragging="handleDragging"
    @resizing="handleResizing"
    @drag-end="handleDragEnd"
    @resize-end="handleResizeEnd"
    class="image-element"
    :class="{ 'edit-mode': isEditMode, 'selected': isSelected }"
  >
    <img
      :src="payload.url"
      :alt="payload.description || 'Image'"
      @dblclick="handleDoubleClick"
      @load="handleImageLoad"
      @error="handleImageError"
      class="image-content"
      :style="imageStyle"
    />
    
    <!-- Loading indicator -->
    <div v-if="isLoading" class="loading-overlay">
      <div class="loading-spinner"></div>
    </div>
    
    <!-- Error indicator -->
    <div v-if="hasError" class="error-overlay">
      <div class="error-icon">⚠️</div>
      <div class="error-text">Failed to load image</div>
    </div>
  </Vue3DraggableResizable>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Vue3DraggableResizable from 'vue3-draggable-resizable'
import type { Element, ImagePayload } from '@/types'

interface Props {
  element: Element
  isSelected: boolean
  isEditMode: boolean
  snapToGrid: boolean
  gridSize: number
}

interface Emits {
  (e: 'select', elementId: string): void
  (e: 'deselect'): void
  (e: 'update', elementId: string, updates: any): void
  (e: 'delete', elementId: string): void
  (e: 'openViewer', elementId: string, imageData: ImagePayload): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isLoading = ref(true)
const hasError = ref(false)

const payload = computed(() => props.element.payload as ImagePayload)

const imageStyle = computed(() => ({
  width: '100%',
  height: '100%',
  objectFit: 'contain' as 'contain',
  display: 'block',
  maxWidth: '100%',
  maxHeight: '100%'
}))

// Drag and resize handlers
const handleActivated = () => {
  emit('select', props.element.id)
}

const handleDeactivated = () => {
  emit('deselect')
}

const handleDragging = (position: { x: number; y: number }) => {
  // Update element position during drag
  emit('update', props.element.id, {
    x: position.x,
    y: position.y
  })
}

const handleResizing = (rect: { x: number; y: number; width: number; height: number }) => {
  // Update element size during resize
  emit('update', props.element.id, {
    w: rect.width,
    h: rect.height,
    x: rect.x,
    y: rect.y
  })
}

const handleDragEnd = (position: { x: number; y: number }) => {
  // Final update after drag
  emit('update', props.element.id, {
    x: position.x,
    y: position.y
  })
}

const handleResizeEnd = (rect: { x: number; y: number; width: number; height: number }) => {
  // Final update after resize
  emit('update', props.element.id, {
    w: rect.width,
    h: rect.height,
    x: rect.x,
    y: rect.y
  })
}

const handleDoubleClick = () => {
  if (props.isEditMode) {
    emit('openViewer', props.element.id, payload.value)
  }
}

const handleImageLoad = () => {
  isLoading.value = false
  hasError.value = false
}

const handleImageError = () => {
  isLoading.value = false
  hasError.value = true
}

// Listen for global delete key
const handleGlobalKeyDown = (event: KeyboardEvent) => {
  if ((event.key === 'Delete' || event.key === 'Backspace') && 
      props.isSelected && props.isEditMode) {
    event.preventDefault()
    emit('delete', props.element.id)
  }
}

// Add global event listener when component is mounted
document.addEventListener('keydown', handleGlobalKeyDown)

// Clean up event listener when component is unmounted
import { onUnmounted } from 'vue'
onUnmounted(() => {
  document.removeEventListener('keydown', handleGlobalKeyDown)
})
</script>

<style scoped>
.image-element {
  position: absolute;
  cursor: pointer;
}

.image-element.edit-mode {
  cursor: move;
}

.image-element.selected {
  outline: 2px solid #2563EB;
  outline-offset: -2px;
}

.image-content {
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid #e5e7eb;
  border-top: 2px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.error-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  border: 2px dashed #ef4444;
}

.error-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.error-text {
  font-size: 12px;
  color: #ef4444;
  text-align: center;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
