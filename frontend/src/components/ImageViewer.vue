<template>
  <div v-if="isVisible" class="image-viewer-overlay fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50" @click="closeViewer">
    <div class="image-viewer-modal bg-white rounded-lg max-w-4xl max-h-[90vh] overflow-hidden shadow-2xl" @click.stop>
      <!-- Header -->
      <div class="flex items-center justify-between p-4 border-b">
        <h3 class="text-lg font-semibold text-gray-900">Image Details</h3>
        <button 
          @click="closeViewer"
          class="text-gray-400 hover:text-gray-600 transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="p-6">
        <!-- Image Display -->
        <div class="mb-6 text-center">
          <img 
            v-if="imageData?.url"
            :src="imageData.url" 
            :alt="imageData.description || 'Image'" 
            class="max-w-full max-h-96 object-contain rounded-lg shadow-md mx-auto"
          />
        </div>

        <!-- Description Section -->
        <div class="space-y-4">
          <label class="block text-sm font-medium text-gray-700">
            Description
          </label>
          
          <div v-if="!isEditingDescription">
            <p 
              v-if="imageData?.description" 
              class="text-gray-900 bg-gray-50 p-3 rounded-lg border cursor-pointer hover:bg-gray-100 transition-colors select-none"
              @dblclick="startEditingDescription"
              @click="handleSingleClick"
            >
              {{ imageData.description }}
            </p>
            <p 
              v-else 
              class="text-gray-500 italic bg-gray-50 p-3 rounded-lg border cursor-pointer hover:bg-gray-100 transition-colors select-none"
              @dblclick="startEditingDescription"
              @click="handleSingleClick"
            >
              Double-click to add a description...
            </p>
          </div>

          <div v-else class="space-y-3">
            <textarea
              ref="descriptionInput"
              v-model="editingDescription"
              @keydown.escape="cancelEditingDescription"
              @keydown.ctrl.enter="saveDescription"
              @keydown.meta.enter="saveDescription"
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
              rows="4"
              placeholder="Enter a description for this image..."
            />
            <div class="flex justify-end space-x-2">
              <button 
                @click="cancelEditingDescription"
                class="px-3 py-1 text-sm text-gray-600 hover:text-gray-800 transition-colors"
              >
                Cancel
              </button>
              <button 
                @click="saveDescription"
                class="px-4 py-2 bg-blue-500 text-white text-sm rounded-lg hover:bg-blue-600 transition-colors"
              >
                Save
              </button>
            </div>
          </div>
          
          <p class="text-xs text-gray-500">
            Tip: Double-click the description to edit, or press Ctrl+Enter to save
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import type { ImagePayload } from '@/types'

interface Props {
  isVisible: boolean
  imageData: ImagePayload | null
  elementId: string | null
  isEditMode: boolean
}

interface Emits {
  (e: 'close'): void
  (e: 'update-description', elementId: string, description: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Description editing state
const isEditingDescription = ref(false)
const editingDescription = ref('')
const descriptionInput = ref<HTMLTextAreaElement>()

// Watch for image data changes to reset editing state
watch(() => props.imageData, (newData) => {
  if (newData) {
    editingDescription.value = newData.description || ''
    isEditingDescription.value = false
  }
})

// Watch for visibility changes to reset state
watch(() => props.isVisible, (visible) => {
  if (!visible) {
    isEditingDescription.value = false
    editingDescription.value = ''
  }
})

const closeViewer = () => {
  emit('close')
}

const handleSingleClick = (event: MouseEvent) => {
  event.preventDefault()
  event.stopPropagation()
  // Do nothing on single click - only double-click should trigger editing
}

const startEditingDescription = () => {
  if (!props.isEditMode) {
    return
  }
  
  isEditingDescription.value = true
  editingDescription.value = props.imageData?.description || ''
  
  nextTick(() => {
    descriptionInput.value?.focus()
  })
}

const cancelEditingDescription = () => {
  isEditingDescription.value = false
  editingDescription.value = props.imageData?.description || ''
}

const saveDescription = () => {
  if (!props.elementId) return
  
  const trimmedDescription = editingDescription.value.trim()
  emit('update-description', props.elementId, trimmedDescription)
  isEditingDescription.value = false
}
</script>

<style scoped>
.image-viewer-overlay {
  backdrop-filter: blur(4px);
}

.image-viewer-modal {
  animation: slideIn 0.2s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

/* Ensure textarea doesn't resize horizontally */
textarea {
  resize: vertical;
  min-height: 100px;
}
</style>
