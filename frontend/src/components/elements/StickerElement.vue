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
    class="sticker-element"
    :class="{ 'edit-mode': isEditMode, 'selected': isSelected }"
  >
    <img
      :src="payload.url"
      :alt="`${payload.stickerType} sticker`"
      @load="handleStickerLoad"
      @error="handleStickerError"
      class="sticker-content"
      :style="stickerStyle"
    />
    
    <!-- Loading indicator -->
    <div v-if="isLoading" class="loading-overlay">
      <div class="loading-spinner"></div>
    </div>
    
    <!-- Error indicator -->
    <div v-if="hasError" class="error-overlay">
      <div class="error-icon">😵</div>
      <div class="error-text">Sticker not found</div>
    </div>
  </Vue3DraggableResizable>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Vue3DraggableResizable from 'vue3-draggable-resizable'
import type { Element, StickerPayload } from '@/types'

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
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isLoading = ref(true)
const hasError = ref(false)

const payload = computed(() => props.element.payload as StickerPayload)

const stickerStyle = computed(() => ({
  width: '100%',
  height: '100%',
  objectFit: 'contain' as 'contain',
  display: 'block',
  maxWidth: '100%',
  maxHeight: '100%',
  userSelect: 'none' as 'none',
  pointerEvents: 'none' as 'none'
}))

// Drag and resize handlers
const handleActivated = () => {
  emit('select', props.element.id)
}

const handleDeactivated = () => {
  emit('deselect')
}

const handleDragging = (position: { x: number; y: number }) => {
  emit('update', props.element.id, {
    x: position.x,
    y: position.y
  })
}

const handleResizing = (rect: { x: number; y: number; width: number; height: number }) => {
  emit('update', props.element.id, {
    w: rect.width,
    h: rect.height,
    x: rect.x,
    y: rect.y
  })
}

const handleDragEnd = (position: { x: number; y: number }) => {
  emit('update', props.element.id, {
    x: position.x,
    y: position.y
  })
}

const handleResizeEnd = (rect: { x: number; y: number; width: number; height: number }) => {
  emit('update', props.element.id, {
    w: rect.width,
    h: rect.height,
    x: rect.x,
    y: rect.y
  })
}

const handleStickerLoad = () => {
  isLoading.value = false
  hasError.value = false
}

const handleStickerError = () => {
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
.sticker-element {
  position: absolute;
  cursor: pointer;
}

.sticker-element.edit-mode {
  cursor: move;
}

.sticker-element.selected {
  outline: 2px solid #2563EB;
  outline-offset: -2px;
}

.sticker-content {
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
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
  width: 20px;
  height: 20px;
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
  font-size: 20px;
  margin-bottom: 4px;
}

.error-text {
  font-size: 10px;
  color: #ef4444;
  text-align: center;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}
</style>
