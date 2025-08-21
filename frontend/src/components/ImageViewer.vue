<template>
  <Dialog :open="isVisible" @update:open="(value) => !value && $emit('close')">
    <DialogContent class="max-w-4xl max-h-[90vh] w-full mx-4 overflow-hidden">
      <!-- Header -->
      <DialogHeader>
        <DialogTitle>Image Details</DialogTitle>
      </DialogHeader>

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
              class="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 resize-none"
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
                class="px-4 py-2 bg-amber-500 text-white text-sm rounded-lg hover:bg-amber-600 transition-colors"
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
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import type { ImagePayload } from '@/types'
import Dialog from '@/components/ui/Dialog.vue'
import DialogContent from '@/components/ui/DialogContent.vue'
import DialogHeader from '@/components/ui/DialogHeader.vue'
import DialogTitle from '@/components/ui/DialogTitle.vue'

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
/* Ensure textarea doesn't resize horizontally */
textarea {
  resize: vertical;
  min-height: 100px;
}
</style>
