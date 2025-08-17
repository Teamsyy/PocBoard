<template>
  <div class="image-uploader">
    <!-- Upload Trigger Button -->
    <button 
      @click="openFileDialog"
      :class="props.variant === 'toolbar' ? 'toolbar-btn inline-flex items-center space-x-1' : 'w-full btn-ghost justify-start'"
      :disabled="uploading"
      title="Upload image (drag & drop or paste supported)"
    >
      <PhotoIcon class="w-4 h-4" :class="props.variant === 'sidebar' ? 'mr-2' : ''" />
      <span v-if="!uploading">
        {{ props.variant === 'toolbar' ? 'Upload' : 'Upload Image' }}
      </span>
      <span v-else>
        {{ props.variant === 'toolbar' ? 'Uploading...' : 'Uploading...' }}
      </span>
    </button>

    <!-- Hidden File Input -->
    <input
      ref="fileInput"
      type="file"
      accept="image/jpeg,image/jpg,image/png,image/gif"
      multiple
      @change="handleFileSelect"
      class="hidden"
    />

    <!-- Drag and Drop Zone (shown when dragging) -->
    <div
      v-if="isDragOver || showDropZone"
      @dragover.prevent="handleDragOver"
      @dragleave="handleDragLeave"
      @drop.prevent="handleDrop"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div 
        class="bg-white rounded-xl p-8 m-4 max-w-md w-full border-4 border-dashed"
        :class="isDragOver ? 'border-blue-500 bg-blue-50' : 'border-gray-300'"
      >
        <div class="text-center">
          <PhotoIcon class="w-16 h-16 mx-auto mb-4 text-gray-400" />
          <h3 class="text-lg font-medium text-gray-900 mb-2">
            Drop images here
          </h3>
          <p class="text-sm text-gray-500 mb-4">
            Supports JPG, PNG, GIF up to 10MB
          </p>
          <button 
            @click="showDropZone = false"
            class="btn-secondary"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>

    <!-- Upload Progress Modal -->
    <div
      v-if="uploading"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    >
      <div class="bg-white rounded-xl p-6 m-4 max-w-md w-full">
        <div class="text-center">
          <LoadingSpinner class="w-8 h-8 mx-auto mb-4" />
          <h3 class="text-lg font-medium text-gray-900 mb-2">
            Uploading Images
          </h3>
          <p class="text-sm text-gray-500 mb-4">
            {{ uploadStatus }}
          </p>
          <div class="w-full bg-gray-200 rounded-full h-2 mb-4">
            <div 
              class="bg-blue-600 h-2 rounded-full transition-all duration-300"
              :style="{ width: `${uploadProgress}%` }"
            ></div>
          </div>
          <p class="text-xs text-gray-400">
            {{ uploadedCount }}/{{ totalFiles }} files uploaded
          </p>
        </div>
      </div>
    </div>

    <!-- Error Toast -->
    <div
      v-if="error"
      class="fixed top-4 right-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded z-50 max-w-md"
    >
      <div class="flex items-start">
        <ExclamationTriangleIcon class="w-5 h-5 mr-2 mt-0.5 flex-shrink-0" />
        <div class="flex-1">
          <strong class="font-medium">Upload Error</strong>
          <p class="text-sm mt-1">{{ error }}</p>
        </div>
        <button @click="error = ''" class="ml-2 text-red-500 hover:text-red-700">
          <XMarkIcon class="w-4 h-4" />
        </button>
      </div>
    </div>

    <!-- Success Toast -->
    <div
      v-if="successMessage"
      class="fixed top-4 right-4 bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded z-50 max-w-md"
    >
      <div class="flex items-start">
        <CheckCircleIcon class="w-5 h-5 mr-2 mt-0.5 flex-shrink-0" />
        <div class="flex-1">
          <strong class="font-medium">Upload Successful</strong>
          <p class="text-sm mt-1">{{ successMessage }}</p>
        </div>
        <button @click="successMessage = ''" class="ml-2 text-green-500 hover:text-green-700">
          <XMarkIcon class="w-4 h-4" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useEditorStore } from '@/stores/editor'
import { useBoardsStore } from '@/stores/boards'
import { uploadsApi } from '@/api'
import { LoadingSpinner } from '@/components'
import {
  PhotoIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

// Props
interface Props {
  autoCreate?: boolean // Automatically create image element after upload
  variant?: 'toolbar' | 'sidebar' // Button style variant
}

const props = withDefaults(defineProps<Props>(), {
  autoCreate: true,
  variant: 'toolbar'
})

// Emits
const emit = defineEmits<{
  uploaded: [url: string, filename: string]
  error: [message: string]
}>()

// Stores
const editorStore = useEditorStore()
const boardsStore = useBoardsStore()

// State
const fileInput = ref<HTMLInputElement>()
const isDragOver = ref(false)
const showDropZone = ref(false)
const uploading = ref(false)
const uploadProgress = ref(0)
const uploadStatus = ref('')
const uploadedCount = ref(0)
const totalFiles = ref(0)
const error = ref('')
const successMessage = ref('')

// Constants
const MAX_FILE_SIZE = 10 * 1024 * 1024 // 10MB
const ALLOWED_TYPES = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']

// Methods
const openFileDialog = () => {
  fileInput.value?.click()
}

const validateFile = (file: File): string | null => {
  if (!ALLOWED_TYPES.includes(file.type)) {
    return `File type ${file.type} is not supported. Please use JPG, PNG, or GIF.`
  }
  
  if (file.size > MAX_FILE_SIZE) {
    const sizeMB = (file.size / (1024 * 1024)).toFixed(1)
    return `File size ${sizeMB}MB exceeds the 10MB limit.`
  }
  
  return null
}

const handleFileSelect = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = Array.from(target.files || [])
  
  if (files.length > 0) {
    await processFiles(files)
  }
  
  // Clear the input so the same file can be selected again
  target.value = ''
}

const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (event: DragEvent) => {
  // Only hide drag overlay if leaving the entire drop zone
  const target = event.currentTarget as HTMLElement
  const related = event.relatedTarget as HTMLElement | null
  if (!target?.contains(related)) {
    isDragOver.value = false
  }
}

const handleDrop = async (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
  showDropZone.value = false
  
  const files = Array.from(event.dataTransfer?.files || [])
  if (files.length > 0) {
    await processFiles(files)
  }
}

const processFiles = async (files: File[]) => {
  if (!boardsStore.currentBoard) {
    error.value = 'No board is currently open'
    return
  }

  // Validate all files first
  const validFiles: File[] = []
  for (const file of files) {
    const validationError = validateFile(file)
    if (validationError) {
      error.value = validationError
      return
    }
    validFiles.push(file)
  }

  if (validFiles.length === 0) return

  try {
    uploading.value = true
    totalFiles.value = validFiles.length
    uploadedCount.value = 0
    uploadProgress.value = 0
    error.value = ''

    const uploadedImages: Array<{ url: string; filename: string }> = []

    for (let i = 0; i < validFiles.length; i++) {
      const file = validFiles[i]
      uploadStatus.value = `Uploading ${file.name}...`
      
      try {
        const result = await uploadsApi.uploadImage(boardsStore.currentBoard.id, file)
        uploadedImages.push(result)
        uploadedCount.value++
        uploadProgress.value = ((i + 1) / validFiles.length) * 100
        
        // Emit uploaded event
        emit('uploaded', result.url, result.filename)
        
        // Auto-create image element if enabled
        if (props.autoCreate) {
          await createImageElement(result.url, file.name)
        }
      } catch (err: any) {
        console.error('Upload failed for file:', file.name, err)
        error.value = `Failed to upload ${file.name}: ${err.error?.message || 'Unknown error'}`
        break
      }
    }

    if (uploadedImages.length > 0) {
      const count = uploadedImages.length
      successMessage.value = `Successfully uploaded ${count} image${count > 1 ? 's' : ''}`
      
      // Auto-hide success message after 3 seconds
      setTimeout(() => {
        successMessage.value = ''
      }, 3000)
    }

  } catch (err: any) {
    console.error('Upload process failed:', err)
    error.value = err.error?.message || 'Upload failed'
    emit('error', error.value)
  } finally {
    uploading.value = false
    uploadProgress.value = 0
    uploadStatus.value = ''
    uploadedCount.value = 0
    totalFiles.value = 0
  }
}

const createImageElement = async (imageUrl: string, _filename: string) => {
  try {
    // Create a temporary image to get dimensions
    const img = new Image()
    img.onload = async () => {
      const maxWidth = 300
      const maxHeight = 300
      
      let width = img.naturalWidth
      let height = img.naturalHeight
      
      // Scale down if too large
      if (width > maxWidth || height > maxHeight) {
        const ratio = Math.min(maxWidth / width, maxHeight / height)
        width = width * ratio
        height = height * ratio
      }
      
      // Position image in center of canvas with some randomness
      const canvasWidth = 800
      const canvasHeight = 600
      const x = (canvasWidth / 2) - (width / 2) + (Math.random() - 0.5) * 100
      const y = (canvasHeight / 2) - (height / 2) + (Math.random() - 0.5) * 100
      
      const payload = {
        url: imageUrl,
        originalWidth: img.naturalWidth,
        originalHeight: img.naturalHeight
      }
      
      await editorStore.createElement('image', x, y, width, height, payload)
    }
    
    img.onerror = () => {
      console.error('Failed to load uploaded image:', imageUrl)
      error.value = 'Failed to load uploaded image'
    }
    
    img.src = imageUrl
  } catch (err) {
    console.error('Failed to create image element:', err)
    error.value = 'Failed to create image element'
  }
}

// Clipboard paste handling
const handlePaste = async (event: ClipboardEvent) => {
  if (!event.clipboardData) return
  
  const items = Array.from(event.clipboardData.items)
  const imageItems = items.filter(item => item.type.startsWith('image/'))
  
  if (imageItems.length === 0) return
  
  event.preventDefault()
  
  const files: File[] = []
  for (const item of imageItems) {
    const file = item.getAsFile()
    if (file) {
      // Create a proper filename for pasted images
      const timestamp = new Date().toISOString().replace(/[:.]/g, '-')
      const extension = file.type.split('/')[1] || 'png'
      const renamedFile = new File([file], `pasted-image-${timestamp}.${extension}`, {
        type: file.type
      })
      files.push(renamedFile)
    }
  }
  
  if (files.length > 0) {
    await processFiles(files)
  }
}

// Global drag and drop handling
const handleGlobalDragEnter = (event: DragEvent) => {
  // Check if files are being dragged
  if (event.dataTransfer?.types.includes('Files')) {
    event.preventDefault()
    showDropZone.value = true
  }
}

const handleGlobalDragOver = (event: DragEvent) => {
  if (event.dataTransfer?.types.includes('Files')) {
    event.preventDefault()
  }
}

// Lifecycle
onMounted(() => {
  // Add global paste listener
  document.addEventListener('paste', handlePaste)
  
  // Add global drag listeners
  document.addEventListener('dragenter', handleGlobalDragEnter)
  document.addEventListener('dragover', handleGlobalDragOver)
})

onUnmounted(() => {
  // Remove global listeners
  document.removeEventListener('paste', handlePaste)
  document.removeEventListener('dragenter', handleGlobalDragEnter)
  document.removeEventListener('dragover', handleGlobalDragOver)
})

// Auto-hide error messages after 5 seconds
const errorTimeout = ref<NodeJS.Timeout>()

// Watch error changes
const handleErrorChange = () => {
  if (errorTimeout.value) {
    clearTimeout(errorTimeout.value)
  }
  if (error.value) {
    errorTimeout.value = setTimeout(() => {
      error.value = ''
    }, 5000)
  }
}

// Watch for error changes
const unwatchError = () => {
  let currentError = error.value
  const checkError = () => {
    if (currentError !== error.value) {
      currentError = error.value
      handleErrorChange()
    }
    requestAnimationFrame(checkError)
  }
  checkError()
}

onMounted(() => {
  unwatchError()
})
</script>

<style scoped>
/* Custom file input styling */
.image-uploader input[type="file"] {
  display: none;
}

/* Ensure button inherits toolbar styles when in toolbar */
.image-uploader .toolbar-btn {
  padding: 0.5rem 0.75rem;
  color: #d97706;
  transition: all 0.2s;
  border-right: 1px solid #f3e8ff;
  background-color: white;
  border: 1px solid #f3e8ff;
  border-radius: 0.375rem;
}

.image-uploader .toolbar-btn:hover:not(:disabled) {
  background-color: #fef3c7;
  color: #92400e;
}

.image-uploader .toolbar-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Drag over animation */
.drag-over {
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
}

/* Upload progress animation */
.upload-progress {
  transition: width 0.3s ease-in-out;
}

/* Toast animations */
.toast-enter-active, .toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from, .toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
