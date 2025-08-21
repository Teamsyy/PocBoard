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
    class="shape-element"
    :class="{ 'edit-mode': isEditMode, 'selected': isSelected }"
  >
    <svg
      :width="element.w"
      :height="element.h"
      class="shape-svg"
      viewBox="0 0 100 100"
      preserveAspectRatio="none"
    >
      <!-- Rectangle -->
      <rect
        v-if="payload.shapeType === 'rectangle'"
        x="0"
        y="0"
        width="100"
        height="100"
        :fill="payload.fill"
        :stroke="payload.stroke"
        :stroke-width="strokeWidthPercent"
      />
      
      <!-- Circle -->
      <circle
        v-else-if="payload.shapeType === 'circle'"
        cx="50"
        cy="50"
        r="45"
        :fill="payload.fill"
        :stroke="payload.stroke"
        :stroke-width="strokeWidthPercent"
      />
      
      <!-- Triangle -->
      <polygon
        v-else-if="payload.shapeType === 'triangle'"
        points="50,5 95,85 5,85"
        :fill="payload.fill"
        :stroke="payload.stroke"
        :stroke-width="strokeWidthPercent"
      />
    </svg>
  </Vue3DraggableResizable>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import Vue3DraggableResizable from 'vue3-draggable-resizable'
import type { Element, ShapePayload } from '@/types'

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

const payload = computed(() => props.element.payload as ShapePayload)

// Calculate stroke width as percentage of viewBox for consistent appearance
const strokeWidthPercent = computed(() => {
  const baseSize = Math.min(props.element.w, props.element.h)
  const strokePercent = (payload.value.strokeWidth / baseSize) * 100
  return Math.max(0.5, Math.min(10, strokePercent)) // Clamp between 0.5% and 10%
})

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
.shape-element {
  position: absolute;
  cursor: pointer;
}

.shape-element.edit-mode {
  cursor: move;
}

.shape-element.selected {
  outline: 2px solid #2563EB;
  outline-offset: -2px;
}

.shape-svg {
  display: block;
  width: 100%;
  height: 100%;
}
</style>
