<template>
  <Vue3DraggableResizable
    :initW="element.w"
    :initH="element.h"
    :x="element.x"
    :y="element.y"
    :draggable="isEditMode"
    :resizable="isEditMode"
    :minW="20"
    :minH="16"
    :parent="true"
    :disableUserSelect="false"
    @activated="handleActivated"
    @deactivated="handleDeactivated"
    @dragging="handleDragging"
    @resizing="handleResizing"
    @drag-end="handleDragEnd"
    @resize-end="handleResizeEnd"
    class="text-element"
    :class="{ 'edit-mode': isEditMode, 'selected': isSelected, 'text-editing': isTextEditing }"
  >
    <div
      ref="textContent"
      :contenteditable="isEditMode && isSelected && isTextEditing"
      @blur="handleTextBlur"
      @input="handleTextInput"
      @dblclick="handleDoubleClick"
      @click="handleTextClick"
      @keydown="handleKeyDown"
      class="text-content"
      :style="textStyle"
      v-html="payload.content"
    ></div>
  </Vue3DraggableResizable>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import Vue3DraggableResizable from 'vue3-draggable-resizable'
import type { Element, TextPayload } from '@/types'

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

const textContent = ref<HTMLDivElement>()
const isTextEditing = ref(false)

const payload = computed(() => props.element.payload as TextPayload)

const textStyle = computed(() => ({
  fontFamily: payload.value.fontFamily,
  fontSize: `${payload.value.fontSize}px`,
  color: payload.value.color,
  fontWeight: payload.value.bold ? 'bold' : 'normal',
  fontStyle: payload.value.italic ? 'italic' : 'normal',
  textAlign: payload.value.textAlign as 'left' | 'center' | 'right',
  width: '100%',
  height: '100%',
  border: 'none',
  outline: 'none',
  background: 'transparent',
  overflow: 'hidden',
  wordWrap: 'break-word' as 'break-word',
  whiteSpace: 'pre-wrap' as 'pre-wrap',
  cursor: isTextEditing.value ? 'text' : 'default',
  lineHeight: '1.4'
}))

// Drag and resize handlers
const handleActivated = () => {
  emit('select', props.element.id)
}

const handleDeactivated = () => {
  // Don't deselect if we're in text editing mode
  if (!isTextEditing.value) {
    emit('deselect')
  }
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

const handleDoubleClick = (event: MouseEvent) => {
  if (props.isEditMode && props.isSelected) {
    event.preventDefault()
    event.stopPropagation()
    enterTextEditMode()
  }
}

const handleTextClick = (event: MouseEvent) => {
  // Prevent click from bubbling up and deselecting the element
  if (isTextEditing.value) {
    event.stopPropagation()
  } else if (props.isEditMode && props.isSelected) {
    // Single click on selected text element enters text editing mode
    event.preventDefault()
    event.stopPropagation()
    enterTextEditMode()
  }
}

const enterTextEditMode = async () => {
  isTextEditing.value = true
  await nextTick()
  if (textContent.value) {
    textContent.value.focus()
    
    // If content is empty or just "New Text", select all
    const content = textContent.value.textContent || ''
    if (content === '' || content === 'New Text') {
      const range = document.createRange()
      const selection = window.getSelection()
      range.selectNodeContents(textContent.value)
      selection?.removeAllRanges()
      selection?.addRange(range)
    } else {
      // Set cursor at the end of content
      const range = document.createRange()
      const selection = window.getSelection()
      range.selectNodeContents(textContent.value)
      range.collapse(false)
      selection?.removeAllRanges()
      selection?.addRange(range)
    }
  }
}

const exitTextEditMode = () => {
  isTextEditing.value = false
  if (textContent.value) {
    textContent.value.blur()
  }
}

// Expose method to parent component
defineExpose({
  exitTextEditMode,
  isTextEditing: () => isTextEditing.value
})

const handleTextInput = (event: Event) => {
  const target = event.target as HTMLDivElement
  const newContent = target.innerHTML || ''
  
  emit('update', props.element.id, {
    payload: {
      ...payload.value,
      content: newContent
    }
  })
}

const handleTextBlur = (event: FocusEvent) => {
  // Check if the focus is moving to a related element (like a toolbar)
  const relatedTarget = event.relatedTarget as HTMLElement
  
  // Don't exit edit mode if clicking on toolbar, panels, or other UI elements
  if (relatedTarget && (
    relatedTarget.closest('.toolbar') ||
    relatedTarget.closest('.sidebar') ||
    relatedTarget.closest('.panel') ||
    relatedTarget.closest('[data-ui-panel]') ||
    relatedTarget.closest('button') ||
    relatedTarget.closest('select') ||
    relatedTarget.closest('input')
  )) {
    // Refocus the text content after a short delay
    setTimeout(() => {
      if (textContent.value && isTextEditing.value) {
        textContent.value.focus()
      }
    }, 10)
    return
  }
  
  exitTextEditMode()
}

const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Escape') {
    event.preventDefault()
    exitTextEditMode()
    return
  }
  
  // F2 key to toggle text editing mode
  if (event.key === 'F2' && props.isSelected && props.isEditMode) {
    event.preventDefault()
    if (isTextEditing.value) {
      exitTextEditMode()
    } else {
      enterTextEditMode()
    }
    return
  }
  
  // Handle rich text formatting shortcuts only when in text editing mode
  if (isTextEditing.value && (event.ctrlKey || event.metaKey)) {
    switch (event.key.toLowerCase()) {
      case 'b':
        event.preventDefault()
        document.execCommand('bold')
        break
      case 'i':
        event.preventDefault()
        document.execCommand('italic')
        break
      case 'u':
        event.preventDefault()
        document.execCommand('underline')
        break
      case 'enter':
        event.preventDefault()
        exitTextEditMode()
        return
    }
  }
  
  // Allow other text editing keys to work normally
  event.stopPropagation()
}

// Watch for selection changes to handle editing
watch(() => props.isSelected, (selected) => {
  if (!selected && isTextEditing.value) {
    exitTextEditMode()
  }
})

// Listen for global delete key and F2 toggle
const handleGlobalKeyDown = (event: KeyboardEvent) => {
  // F2 key to toggle text editing mode when element is selected
  if (event.key === 'F2' && props.isSelected && props.isEditMode) {
    event.preventDefault()
    if (isTextEditing.value) {
      exitTextEditMode()
    } else {
      enterTextEditMode()
    }
    return
  }
  
  // Only handle delete if the element is selected but not in text editing mode
  if ((event.key === 'Delete' || event.key === 'Backspace') && 
      props.isSelected && props.isEditMode && !isTextEditing.value) {
    
    // Make sure no input elements are focused
    const activeElement = document.activeElement
    if (activeElement && (
      activeElement.tagName === 'INPUT' ||
      activeElement.tagName === 'TEXTAREA' ||
      activeElement.getAttribute('contenteditable') === 'true'
    )) {
      return // Don't delete if an input is focused
    }
    
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
.text-element {
  position: absolute;
  cursor: pointer;
}

.text-element.edit-mode {
  cursor: move;
}

.text-element.edit-mode.selected {
  cursor: text;
}

.text-element.edit-mode.selected .text-content:hover {
  background: rgba(255, 255, 255, 0.2);
}

.text-element.selected {
  outline: 2px solid #2563EB;
  outline-offset: -2px;
}

.text-element.editing {
  z-index: 1000;
}

.text-element.text-editing {
  z-index: 1001;
}

.text-element.text-editing .text-content {
  cursor: text !important;
}

.text-content {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  padding: 2px;
  min-height: 100%;
  box-sizing: border-box;
}

.text-element.selected .text-content {
  cursor: text;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  transition: background-color 0.2s ease;
}

.text-element.selected .text-content:hover {
  background: rgba(255, 255, 255, 0.2);
}

.text-content[contenteditable="true"] {
  cursor: text !important;
  background: rgba(255, 255, 255, 0.95) !important;
  border: 2px solid #2563EB;
  border-radius: 3px;
  box-shadow: 0 0 0 1px rgba(37, 99, 235, 0.2);
}

.text-content:focus {
  outline: none;
}

/* Rich text formatting styles */
.text-content b,
.text-content strong {
  font-weight: bold;
}

.text-content i,
.text-content em {
  font-style: italic;
}

.text-content u {
  text-decoration: underline;
}
</style>
