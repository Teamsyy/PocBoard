<template>
  <div v-if="isVisible" class="sticker-picker-overlay fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="sticker-picker bg-white rounded-lg shadow-xl w-96 max-h-[80vh] flex flex-col">
      <!-- Header -->
      <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
        <h3 class="text-lg font-semibold text-gray-900">Add Sticker</h3>
        <button 
          @click="$emit('close')"
          class="text-gray-400 hover:text-gray-600"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Tabs -->
      <div class="px-6 pt-4">
        <div class="flex space-x-1 bg-gray-100 rounded-lg p-1">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            @click="activeTab = tab.id"
            :class="[
              'flex-1 px-3 py-2 text-sm font-medium rounded-md transition-colors',
              activeTab === tab.id 
                ? 'bg-white text-gray-900 shadow-sm' 
                : 'text-gray-600 hover:text-gray-900'
            ]"
          >
            {{ tab.label }}
          </button>
        </div>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-hidden">
        <!-- Upload Tab -->
        <div v-if="activeTab === 'upload'" class="p-6">
          <!-- Upload not available message -->
          <div v-if="!canUpload" class="text-center py-8">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 15.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
            <p class="mt-2 text-sm text-gray-500">Upload not available</p>
            <p class="text-xs text-gray-400">Please use preset stickers instead</p>
          </div>

          <div v-else class="space-y-4">
            <!-- Upload Area -->
            <div 
              @drop="handleDrop"
              @dragover.prevent
              @dragenter.prevent
              :class="[
                'border-2 border-dashed rounded-lg p-8 text-center transition-colors',
                isDragOver ? 'border-blue-500 bg-blue-50' : 'border-gray-300'
              ]"
            >
              <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
              <p class="mt-2 text-sm text-gray-600">
                <span class="font-medium">Drop your sticker here</span> or
                <button 
                  @click="triggerFileInput"
                  class="text-blue-600 hover:text-blue-500"
                >
                  browse files
                </button>
              </p>
              <p class="text-xs text-gray-500 mt-1">PNG, JPG, GIF, SVG up to 5MB</p>
            </div>

            <!-- File Input -->
            <input
              ref="fileInput"
              type="file"
              accept=".png,.jpg,.jpeg,.gif,.svg"
              @change="handleFileSelect"
              class="hidden"
            />

            <!-- Upload Preview -->
            <div v-if="uploadPreview" class="space-y-3">
              <div class="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg">
                <img 
                  :src="uploadPreview.url" 
                  :alt="uploadPreview.name"
                  class="w-12 h-12 object-cover rounded"
                />
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-gray-900 truncate">{{ uploadPreview.name }}</p>
                  <p class="text-xs text-gray-500">{{ formatFileSize(uploadPreview.size) }}</p>
                </div>
                <button 
                  @click="clearUpload"
                  class="text-gray-400 hover:text-gray-600"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>

              <!-- Category Selection -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Category</label>
                <select 
                  v-model="uploadCategory"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="custom">Custom</option>
                  <option value="emoji">Emoji</option>
                  <option value="decoration">Decoration</option>
                  <option value="shapes">Shapes</option>
                  <option value="animals">Animals</option>
                  <option value="food">Food</option>
                  <option value="travel">Travel</option>
                  <option value="nature">Nature</option>
                </select>
              </div>

              <!-- Upload Button -->
              <button 
                @click="uploadSticker"
                :disabled="isUploading"
                :class="[
                  'w-full px-4 py-2 rounded-md font-medium',
                  isUploading 
                    ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
                    : 'bg-blue-600 text-white hover:bg-blue-700'
                ]"
              >
                {{ isUploading ? 'Uploading...' : 'Upload & Add Sticker' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Presets Tab -->
        <div v-if="activeTab === 'presets'" class="p-6">
          <!-- Category Filter -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Category</label>
            <select 
              v-model="selectedCategory"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value="all">All Categories</option>
              <option value="emoji">Emoji</option>
              <option value="decoration">Decoration</option>
              <option value="shapes">Shapes</option>
              <option value="animals">Animals</option>
              <option value="food">Food</option>
              <option value="travel">Travel</option>
              <option value="nature">Nature</option>
            </select>
          </div>

          <!-- Sticker Grid -->
          <div class="grid grid-cols-4 gap-3 max-h-64 overflow-y-auto">
            <button
              v-for="sticker in filteredPresets"
              :key="sticker.id"
              @click="selectPresetSticker(sticker)"
              class="aspect-square p-2 border border-gray-200 rounded-lg hover:border-blue-500 hover:bg-blue-50 transition-colors"
              :title="sticker.name"
            >
              <img 
                :src="sticker.url" 
                :alt="sticker.name"
                class="w-full h-full object-contain"
              />
            </button>
          </div>

          <!-- Empty State -->
          <div v-if="filteredPresets.length === 0" class="text-center py-8">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6 8h6m6-12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <p class="mt-2 text-sm text-gray-500">No stickers found in this category</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { StickerUploadService } from '@/api/stickers'

interface Props {
  isVisible: boolean
  boardId?: string
}

interface Emits {
  (e: 'close'): void
  (e: 'select', sticker: { url: string, stickerType: string, category: string }): void
}

interface PresetSticker {
  id: string
  name: string
  url: string
  category: string
}

interface UploadPreview {
  name: string
  url: string
  size: number
  file: File
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// Refs
const activeTab = ref('presets')
const isDragOver = ref(false)
const fileInput = ref<HTMLInputElement>()
const uploadPreview = ref<UploadPreview | null>(null)
const uploadCategory = ref('custom')
const isUploading = ref(false)
const selectedCategory = ref('all')

// Tabs configuration
const tabs = [
  { id: 'presets', label: 'Presets' },
  { id: 'upload', label: 'Upload' }
]

// Preset stickers (you can expand this with more categories and stickers)
const presetStickers = ref<PresetSticker[]>([
  // Emoji category
  { id: 'emoji-1', name: 'Happy Face', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f600.png', category: 'emoji' },
  { id: 'emoji-2', name: 'Heart Eyes', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f60d.png', category: 'emoji' },
  { id: 'emoji-3', name: 'Thumbs Up', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f44d.png', category: 'emoji' },
  { id: 'emoji-4', name: 'Fire', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f525.png', category: 'emoji' },
  { id: 'emoji-5', name: 'Star', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/2b50.png', category: 'emoji' },
  { id: 'emoji-6', name: 'Party', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f389.png', category: 'emoji' },
  
  // Animals category
  { id: 'animal-1', name: 'Cat Face', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f431.png', category: 'animals' },
  { id: 'animal-2', name: 'Dog Face', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f436.png', category: 'animals' },
  { id: 'animal-3', name: 'Unicorn', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f984.png', category: 'animals' },
  { id: 'animal-4', name: 'Butterfly', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f98b.png', category: 'animals' },
  
  // Food category
  { id: 'food-1', name: 'Pizza', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f355.png', category: 'food' },
  { id: 'food-2', name: 'Cake', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f382.png', category: 'food' },
  { id: 'food-3', name: 'Coffee', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/2615.png', category: 'food' },
  { id: 'food-4', name: 'Apple', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f34e.png', category: 'food' },
  
  // Nature category
  { id: 'nature-1', name: 'Moon', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f319.png', category: 'nature' },
  { id: 'nature-2', name: 'Rainbow', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f308.png', category: 'nature' },
  { id: 'nature-3', name: 'Tree', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f333.png', category: 'nature' },
  
  // Travel category
  { id: 'travel-1', name: 'Car', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f697.png', category: 'travel' },
  { id: 'travel-2', name: 'Camera', url: 'https://cdn.jsdelivr.net/npm/emoji-datasource-apple@15.0.1/img/apple/64/1f4f7.png', category: 'travel' }
])

// Computed
const filteredPresets = computed(() => {
  if (selectedCategory.value === 'all') {
    return presetStickers.value
  }
  return presetStickers.value.filter(sticker => sticker.category === selectedCategory.value)
})

const canUpload = computed(() => {
  return props.boardId
})

// Methods
const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    processFile(file)
  }
}

const handleDrop = (event: DragEvent) => {
  event.preventDefault()
  isDragOver.value = false
  
  const file = event.dataTransfer?.files[0]
  if (file) {
    processFile(file)
  }
}

const processFile = (file: File) => {
  // Validate file type
  const allowedTypes = ['image/png', 'image/jpeg', 'image/gif', 'image/svg+xml']
  if (!allowedTypes.includes(file.type)) {
    alert('Please select a valid image file (PNG, JPG, GIF, or SVG)')
    return
  }
  
  // Validate file size (5MB limit)
  if (file.size > 5 * 1024 * 1024) {
    alert('File size must be less than 5MB')
    return
  }
  
  // Create preview
  const url = URL.createObjectURL(file)
  uploadPreview.value = {
    name: file.name,
    url,
    size: file.size,
    file
  }
}

const clearUpload = () => {
  if (uploadPreview.value?.url) {
    URL.revokeObjectURL(uploadPreview.value.url)
  }
  uploadPreview.value = null
}

const uploadSticker = async () => {
  if (!uploadPreview.value) return
  
  // Check if we have the required props for upload
  if (!props.boardId) {
    alert('Board ID missing. Please try again.')
    return
  }
  
  isUploading.value = true
  
  try {
    // Upload sticker using the real API
    const stickerData = await StickerUploadService.uploadSticker(
      props.boardId,
      uploadPreview.value.file,
      uploadCategory.value
    )
    
    // Emit the sticker selection
    emit('select', stickerData)
    
    clearUpload()
    emit('close')
  } catch (error) {
    console.error('Upload failed:', error)
    const errorMessage = error instanceof Error ? error.message : 'Upload failed. Please try again.'
    alert(errorMessage)
  } finally {
    isUploading.value = false
  }
}

const selectPresetSticker = (sticker: PresetSticker) => {
  emit('select', {
    url: sticker.url,
    stickerType: sticker.name.toLowerCase().replace(/\s+/g, '-'),
    category: sticker.category
  })
  emit('close')
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Cleanup
onMounted(() => {
  // Clean up any existing preview URLs when component unmounts
  return () => {
    if (uploadPreview.value?.url) {
      URL.revokeObjectURL(uploadPreview.value.url)
    }
  }
})
</script>

<style scoped>
.sticker-picker-overlay {
  backdrop-filter: blur(4px);
}

.sticker-picker {
  max-width: 90vw;
}

/* Custom scrollbar for sticker grid */
.max-h-64::-webkit-scrollbar {
  width: 6px;
}

.max-h-64::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.max-h-64::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.max-h-64::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
