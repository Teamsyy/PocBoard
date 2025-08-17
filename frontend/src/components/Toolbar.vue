<template>
  <div class="toolbar bg-white border-b border-amber-200 shadow-sm">
    <div class="px-4 py-2 flex items-center justify-between">
      <!-- Left: Element Creation Tools -->
      <div class="flex items-center space-x-1">
        <div class="toolbar-group">
          <button 
            @click="$emit('addText')"
            :disabled="!isEditMode"
            class="toolbar-btn"
            title="Add Text (T)"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
            </svg>
          </button>
          
          <!-- Image Upload Button -->
          <ImageUploader :auto-create="true" />
          
          <button 
            @click="$emit('addShape')"
            :disabled="!isEditMode"
            class="toolbar-btn"
            title="Add Shape (S)"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
          </button>
          
          <button 
            @click="$emit('addSticker')"
            :disabled="!isEditMode"
            class="toolbar-btn"
            title="Add Sticker"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1.01M15 10h1.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
        </div>

        <!-- Text Styling Controls (shown when text is selected) -->
        <div v-if="hasTextSelection" class="toolbar-group ml-4">
          <select 
            v-model="fontFamily"
            @change="updateTextStyle('fontFamily', fontFamily)"
            class="toolbar-select"
            title="Font Family"
          >
            <option value="Arial">Arial</option>
            <option value="Helvetica">Helvetica</option>
            <option value="Times New Roman">Times New Roman</option>
            <option value="Georgia">Georgia</option>
            <option value="Verdana">Verdana</option>
            <option value="Courier New">Courier New</option>
            <option value="Comic Sans MS">Comic Sans MS</option>
          </select>
          
          <select 
            v-model="fontSize"
            @change="updateTextStyle('fontSize', fontSize)"
            class="toolbar-select"
            title="Font Size"
          >
            <option value="12">12px</option>
            <option value="14">14px</option>
            <option value="16">16px</option>
            <option value="18">18px</option>
            <option value="20">20px</option>
            <option value="24">24px</option>
            <option value="28">28px</option>
            <option value="32">32px</option>
            <option value="36">36px</option>
            <option value="48">48px</option>
          </select>
          
          <input 
            v-model="textColor"
            @change="updateTextStyle('fill', textColor)"
            type="color"
            class="toolbar-color-input"
            title="Text Color"
          />
          
          <button 
            @click="toggleBold"
            :class="['toolbar-btn', { 'active': isBold }]"
            title="Bold (Ctrl+B)"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 4h8a4 4 0 014 4 4 4 0 01-4 4H6z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 12h9a4 4 0 014 4 4 4 0 01-4 4H6z" />
            </svg>
          </button>
          
          <button 
            @click="toggleItalic"
            :class="['toolbar-btn', { 'active': isItalic }]"
            title="Italic (Ctrl+I)"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 4H9m11 16H9m-5-16l4 16" />
            </svg>
          </button>
        </div>

        <!-- Z-Order Controls (shown when elements are selected) -->
        <div v-if="hasSelection" class="toolbar-group ml-4">
          <button 
            @click="bringToFront"
            :disabled="!isEditMode"
            class="toolbar-btn"
            title="Bring to Front"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 11l5-5m0 0l5 5m-5-5v12" />
            </svg>
          </button>
          
          <button 
            @click="bringForward"
            :disabled="!isEditMode"
            class="toolbar-btn"
            title="Bring Forward"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
            </svg>
          </button>
          
          <button 
            @click="sendBackward"
            :disabled="!isEditMode"
            class="toolbar-btn"
            title="Send Backward"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
          
          <button 
            @click="sendToBack"
            :disabled="!isEditMode"
            class="toolbar-btn"
            title="Send to Back"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 13l-5 5m0 0l-5-5m5 5V6" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Right: Canvas Controls -->
      <div class="flex items-center space-x-1">
        <div class="toolbar-group">
          <button 
            @click="toggleSnapToGrid"
            :class="['toolbar-btn', { 'active': snapToGrid }]"
            title="Snap to Grid (8px)"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
            </svg>
          </button>
        </div>
        
        <div class="toolbar-group ml-4">
          <button 
            @click="undo"
            :disabled="!canUndo"
            class="toolbar-btn"
            title="Undo (Ctrl+Z)"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
            </svg>
          </button>
          
          <button 
            @click="redo"
            :disabled="!canRedo"
            class="toolbar-btn"
            title="Redo (Ctrl+Y)"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 10h-10a8 8 0 00-8 8v2m18-10l-6 6m6-6l-6-6" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useEditorStore } from '@/stores/editor'
import { ImageUploader } from '@/components'

interface Props {
  isEditMode: boolean
}

const props = defineProps<Props>()

// Emits for parent component
defineEmits<{
  addText: []
  addShape: []
  addSticker: []
}>()

// Store
const editorStore = useEditorStore()

// Computed properties
const hasSelection = computed(() => editorStore.selectedElementIds.length > 0)
const hasTextSelection = computed(() => {
  return editorStore.selectedElements.some(el => el.kind === 'text')
})
const snapToGrid = computed(() => editorStore.snapToGrid)
const canUndo = computed(() => editorStore.canUndo)
const canRedo = computed(() => editorStore.canRedo)

// Text styling state
const fontFamily = ref('Arial')
const fontSize = ref(16)
const textColor = ref('#000000')
const isBold = ref(false)
const isItalic = ref(false)

// Watch for selection changes to update text styling controls
watch(() => editorStore.selectedElements, (selectedElements) => {
  if (selectedElements.length === 1 && selectedElements[0].kind === 'text') {
    const textElement = selectedElements[0]
    const payload = textElement.payload as any
    
    if (payload) {
      fontFamily.value = payload.fontFamily || 'Arial'
      fontSize.value = payload.fontSize || 16
      textColor.value = payload.fill || '#000000'
      isBold.value = payload.fontWeight === 'bold' || payload.fontWeight >= 700
      isItalic.value = payload.fontStyle === 'italic'
    }
  }
}, { immediate: true })

// Methods
const toggleSnapToGrid = () => {
  editorStore.toggleSnapToGrid()
}

const undo = () => {
  editorStore.undo()
}

const redo = () => {
  editorStore.redo()
}

const updateTextStyle = (property: string, value: any) => {
  if (!props.isEditMode || !hasTextSelection.value) return
  
  editorStore.updateSelectedElementsStyle(property, value)
}

const toggleBold = () => {
  const newWeight = isBold.value ? 'normal' : 'bold'
  isBold.value = !isBold.value
  updateTextStyle('fontWeight', newWeight)
}

const toggleItalic = () => {
  const newStyle = isItalic.value ? 'normal' : 'italic'
  isItalic.value = !isItalic.value
  updateTextStyle('fontStyle', newStyle)
}

const bringToFront = () => {
  if (!props.isEditMode) return
  editorStore.bringToFront()
}

const bringForward = () => {
  if (!props.isEditMode) return
  editorStore.bringForward()
}

const sendBackward = () => {
  if (!props.isEditMode) return
  editorStore.sendBackward()
}

const sendToBack = () => {
  if (!props.isEditMode) return
  editorStore.sendToBack()
}
</script>

<style scoped>
.toolbar-group {
  display: flex;
  align-items: center;
  border: 1px solid #f3e8ff;
  border-radius: 0.375rem;
  background-color: white;
}

.toolbar-btn {
  padding: 0.5rem 0.75rem;
  color: #d97706;
  transition: all 0.2s;
  border-right: 1px solid #f3e8ff;
}

.toolbar-btn:hover:not(:disabled) {
  background-color: #fef3c7;
  color: #92400e;
}

.toolbar-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.toolbar-btn:first-child {
  border-top-left-radius: 0.375rem;
  border-bottom-left-radius: 0.375rem;
}

.toolbar-btn:last-child {
  border-top-right-radius: 0.375rem;
  border-bottom-right-radius: 0.375rem;
  border-right: none;
}

.toolbar-btn.active {
  background-color: #fde68a;
  color: #92400e;
}

.toolbar-select {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
  border-right: 1px solid #f3e8ff;
  background-color: white;
  color: #d97706;
  outline: none;
  min-width: 80px;
}

.toolbar-select:focus {
  background-color: #fef3c7;
}

.toolbar-select:first-child {
  border-top-left-radius: 0.375rem;
  border-bottom-left-radius: 0.375rem;
}

.toolbar-select:last-child {
  border-top-right-radius: 0.375rem;
  border-bottom-right-radius: 0.375rem;
  border-right: none;
}

.toolbar-color-input {
  width: 2rem;
  height: 2rem;
  border-right: 1px solid #f3e8ff;
  cursor: pointer;
  -webkit-appearance: none;
  appearance: none;
  background: none;
  border: none;
  outline: none;
}

.toolbar-color-input:first-child {
  border-top-left-radius: 0.375rem;
  border-bottom-left-radius: 0.375rem;
}

.toolbar-color-input:last-child {
  border-top-right-radius: 0.375rem;
  border-bottom-right-radius: 0.375rem;
  border-right: none;
}

.toolbar-color-input::-webkit-color-swatch-wrapper {
  padding: 0;
  border: none;
}

.toolbar-color-input::-webkit-color-swatch {
  border: none;
  border-radius: 2px;
}
</style>
