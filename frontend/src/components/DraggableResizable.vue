<template>
  <div
    ref="elementRef"
    class="draggable-resizable"
    :class="{ 
      'active': isSelected,
      'edit-mode': isEditMode 
    }"
    :style="elementStyle"
    @mousedown="handleMouseDown"
    @click="handleClick"
  >
    <!-- Main content slot -->
    <div class="content-wrapper">
      <slot />
    </div>

    <!-- Resize handles (only show when selected and in edit mode) -->
    <template v-if="isSelected && isEditMode && resizable">
      <div class="resize-handle resize-handle-nw" @mousedown="handleResizeStart('nw', $event)"></div>
      <div class="resize-handle resize-handle-n" @mousedown="handleResizeStart('n', $event)"></div>
      <div class="resize-handle resize-handle-ne" @mousedown="handleResizeStart('ne', $event)"></div>
      <div class="resize-handle resize-handle-w" @mousedown="handleResizeStart('w', $event)"></div>
      <div class="resize-handle resize-handle-e" @mousedown="handleResizeStart('e', $event)"></div>
      <div class="resize-handle resize-handle-sw" @mousedown="handleResizeStart('sw', $event)"></div>
      <div class="resize-handle resize-handle-s" @mousedown="handleResizeStart('s', $event)"></div>
      <div class="resize-handle resize-handle-se" @mousedown="handleResizeStart('se', $event)"></div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Props {
  x: number
  y: number
  width: number
  height: number
  zIndex: number
  isSelected: boolean
  isEditMode: boolean
  draggable?: boolean
  resizable?: boolean
  minWidth?: number
  minHeight?: number
  snapToGrid?: boolean
  gridSize?: number
}

interface Emits {
  (e: 'activated'): void
  (e: 'deactivated'): void
  (e: 'dragging', position: { x: number; y: number }): void
  (e: 'resizing', rect: { x: number; y: number; width: number; height: number }): void
  (e: 'dragEnd', position: { x: number; y: number }): void
  (e: 'resizeEnd', rect: { x: number; y: number; width: number; height: number }): void
}

const props = withDefaults(defineProps<Props>(), {
  draggable: true,
  resizable: true,
  minWidth: 20,
  minHeight: 20,
  snapToGrid: false,
  gridSize: 8
})

const emit = defineEmits<Emits>()

const elementRef = ref<HTMLElement>()
const isDragging = ref(false)
const isResizing = ref(false)
const resizeDirection = ref('')

// Current position and size (local state for smooth dragging/resizing)
const currentX = ref(props.x)
const currentY = ref(props.y)
const currentWidth = ref(props.width)
const currentHeight = ref(props.height)

// Drag state
const dragStart = ref({ x: 0, y: 0, elementX: 0, elementY: 0 })
const resizeStart = ref({ x: 0, y: 0, width: 0, height: 0, elementX: 0, elementY: 0 })

const elementStyle = computed(() => ({
  position: 'absolute' as const,
  left: `${currentX.value}px`,
  top: `${currentY.value}px`,
  width: `${currentWidth.value}px`,
  height: `${currentHeight.value}px`,
  // Constrain z-index between 1 and 999 to ensure elements are always visible and toolbar stays above
  zIndex: Math.max(1, Math.min(props.zIndex, 999)),
  userSelect: isDragging.value || isResizing.value ? ('none' as const) : ('auto' as const)
}))

const snapToGrid = (value: number): number => {
  if (!props.snapToGrid) return value
  return Math.round(value / props.gridSize) * props.gridSize
}

const handleClick = (event: MouseEvent) => {
  event.stopPropagation()
  if (!isDragging.value && !isResizing.value) {
    emit('activated')
  }
}

const handleMouseDown = (event: MouseEvent) => {
  if (!props.draggable || !props.isEditMode) return
  if ((event.target as HTMLElement).classList.contains('resize-handle')) return

  event.preventDefault()
  event.stopPropagation()

  isDragging.value = true
  
  dragStart.value = {
    x: event.clientX,
    y: event.clientY,
    elementX: currentX.value,
    elementY: currentY.value
  }

  document.addEventListener('mousemove', handleDragMove)
  document.addEventListener('mouseup', handleDragEnd)
  
  emit('activated')
}

const handleDragMove = (event: MouseEvent) => {
  if (!isDragging.value) return

  const deltaX = event.clientX - dragStart.value.x
  const deltaY = event.clientY - dragStart.value.y
  
  let newX = dragStart.value.elementX + deltaX
  let newY = dragStart.value.elementY + deltaY

  if (props.snapToGrid) {
    newX = snapToGrid(newX)
    newY = snapToGrid(newY)
  }

  currentX.value = Math.max(0, newX)
  currentY.value = Math.max(0, newY)

  emit('dragging', { x: currentX.value, y: currentY.value })
}

const handleDragEnd = () => {
  if (!isDragging.value) return

  isDragging.value = false
  
  document.removeEventListener('mousemove', handleDragMove)
  document.removeEventListener('mouseup', handleDragEnd)

  emit('dragEnd', { x: currentX.value, y: currentY.value })
}

const handleResizeStart = (direction: string, event: MouseEvent) => {
  if (!props.resizable || !props.isEditMode) return

  event.preventDefault()
  event.stopPropagation()

  isResizing.value = true
  resizeDirection.value = direction

  resizeStart.value = {
    x: event.clientX,
    y: event.clientY,
    width: currentWidth.value,
    height: currentHeight.value,
    elementX: currentX.value,
    elementY: currentY.value
  }

  document.addEventListener('mousemove', handleResizeMove)
  document.addEventListener('mouseup', handleResizeEnd)
}

const handleResizeMove = (event: MouseEvent) => {
  if (!isResizing.value) return

  const deltaX = event.clientX - resizeStart.value.x
  const deltaY = event.clientY - resizeStart.value.y

  let newWidth = resizeStart.value.width
  let newHeight = resizeStart.value.height
  let newX = resizeStart.value.elementX
  let newY = resizeStart.value.elementY

  const direction = resizeDirection.value

  // Handle width changes
  if (direction.includes('e')) {
    newWidth = Math.max(props.minWidth, resizeStart.value.width + deltaX)
  } else if (direction.includes('w')) {
    newWidth = Math.max(props.minWidth, resizeStart.value.width - deltaX)
    newX = resizeStart.value.elementX + (resizeStart.value.width - newWidth)
  }

  // Handle height changes
  if (direction.includes('s')) {
    newHeight = Math.max(props.minHeight, resizeStart.value.height + deltaY)
  } else if (direction.includes('n')) {
    newHeight = Math.max(props.minHeight, resizeStart.value.height - deltaY)
    newY = resizeStart.value.elementY + (resizeStart.value.height - newHeight)
  }

  if (props.snapToGrid) {
    newWidth = snapToGrid(newWidth)
    newHeight = snapToGrid(newHeight)
    newX = snapToGrid(newX)
    newY = snapToGrid(newY)
  }

  currentX.value = Math.max(0, newX)
  currentY.value = Math.max(0, newY)
  currentWidth.value = newWidth
  currentHeight.value = newHeight

  emit('resizing', {
    x: currentX.value,
    y: currentY.value,
    width: currentWidth.value,
    height: currentHeight.value
  })
}

const handleResizeEnd = () => {
  if (!isResizing.value) return

  isResizing.value = false
  resizeDirection.value = ''

  document.removeEventListener('mousemove', handleResizeMove)
  document.removeEventListener('mouseup', handleResizeEnd)

  emit('resizeEnd', {
    x: currentX.value,
    y: currentY.value,
    width: currentWidth.value,
    height: currentHeight.value
  })
}

// Update local state when props change
const updateFromProps = () => {
  currentX.value = props.x
  currentY.value = props.y
  currentWidth.value = props.width
  currentHeight.value = props.height
}

// Watch for prop changes
import { watch } from 'vue'
watch([() => props.x, () => props.y, () => props.width, () => props.height], updateFromProps)

onMounted(() => {
  updateFromProps()
})

onUnmounted(() => {
  document.removeEventListener('mousemove', handleDragMove)
  document.removeEventListener('mouseup', handleDragEnd)
  document.removeEventListener('mousemove', handleResizeMove)
  document.removeEventListener('mouseup', handleResizeEnd)
})
</script>

<style scoped>
.draggable-resizable {
  border: 2px solid transparent;
  box-sizing: border-box;
}

.draggable-resizable.active {
  border-color: #2563EB;
}

.draggable-resizable.edit-mode {
  cursor: move;
}

.content-wrapper {
  width: 100%;
  height: 100%;
  pointer-events: auto;
}

.resize-handle {
  position: absolute;
  background: #2563EB;
  border: 2px solid #ffffff;
  border-radius: 2px;
  width: 12px;
  height: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  z-index: 10;
}

.resize-handle:hover {
  background: #1E40AF;
}

.resize-handle-nw {
  top: -6px;
  left: -6px;
  cursor: nw-resize;
}

.resize-handle-n {
  top: -6px;
  left: 50%;
  transform: translateX(-50%);
  cursor: n-resize;
}

.resize-handle-ne {
  top: -6px;
  right: -6px;
  cursor: ne-resize;
}

.resize-handle-w {
  top: 50%;
  left: -6px;
  transform: translateY(-50%);
  cursor: w-resize;
}

.resize-handle-e {
  top: 50%;
  right: -6px;
  transform: translateY(-50%);
  cursor: e-resize;
}

.resize-handle-sw {
  bottom: -6px;
  left: -6px;
  cursor: sw-resize;
}

.resize-handle-s {
  bottom: -6px;
  left: 50%;
  transform: translateX(-50%);
  cursor: s-resize;
}

.resize-handle-se {
  bottom: -6px;
  right: -6px;
  cursor: se-resize;
}
</style>
