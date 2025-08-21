<template>
  <div class="share-export-bar bg-white border-t border-gray-200 px-4 py-3">
    <div class="flex items-center justify-between max-w-6xl mx-auto">
      <!-- Left side - Sharing controls -->
      <div class="flex items-center space-x-4">
        <h3 class="text-sm font-medium text-gray-700">Share Board</h3>
        
        <!-- Edit URL -->
        <div class="flex items-center space-x-2">
          <div class="flex items-center bg-gray-50 rounded-md px-3 py-2 text-sm">
            <svg class="w-4 h-4 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
            <span class="text-gray-600 mr-2">Edit:</span>
            <input
              ref="editUrlInput"
              :value="editUrl"
              readonly
              class="bg-transparent border-none outline-none text-gray-800 font-mono text-xs flex-1 min-w-0"
              @focus="(event) => {
                if (event.target && 'select' in event.target) {
                  (event.target as HTMLInputElement).select()
                }
              }"
            />
          </div>
          <button
            @click="copyEditUrl"
            class="btn-icon"
            :class="{ 'text-green-600': editUrlCopied }"
            title="Copy edit URL"
          >
            <svg v-if="!editUrlCopied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        </div>

        <!-- Public URL -->
        <div class="flex items-center space-x-2">
          <div class="flex items-center bg-gray-50 rounded-md px-3 py-2 text-sm">
            <svg class="w-4 h-4 text-gray-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            <span class="text-gray-600 mr-2">View:</span>
            <input
              ref="publicUrlInput"
              :value="publicUrl"
              readonly
              class="bg-transparent border-none outline-none text-gray-800 font-mono text-xs flex-1 min-w-0"
              @focus="(event) => {
                if (event.target && 'select' in event.target) {
                  (event.target as HTMLInputElement).select()
                }
              }"
            />
          </div>
          <button
            @click="copyPublicUrl"
            class="btn-icon"
            :class="{ 'text-green-600': publicUrlCopied }"
            title="Copy public URL"
          >
            <svg v-if="!publicUrlCopied" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Right side - Export controls -->
      <div class="flex items-center space-x-3">
        <div class="h-6 w-px bg-gray-300"></div>
        
        <!-- Export dropdown -->
        <div class="relative" ref="exportDropdown">
          <button
            @click="showExportMenu = !showExportMenu"
            class="btn-secondary flex items-center space-x-2"
            :disabled="isExporting"
          >
            <svg v-if="!isExporting" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <svg v-else class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ isExporting ? 'Exporting...' : 'Export' }}</span>
            <svg class="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>

          <!-- Export menu -->
          <Transition name="fade">
            <div
              v-if="showExportMenu"
              class="absolute right-0 top-full mt-2 w-56 bg-white border border-gray-200 rounded-lg shadow-lg z-50"
            >
              <div class="py-1">
                <button
                  @click="exportCurrentPage"
                  class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center"
                  :disabled="isExporting"
                >
                  <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  Export Current Page as PNG
                </button>
                <button
                  @click="exportAllPages"
                  class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center"
                  :disabled="isExporting"
                >
                  <svg class="w-4 h-4 mr-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                  </svg>
                  Export All Pages as PNG
                </button>
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </div>

    <!-- Success notification -->
    <Transition name="slide-up">
      <div
        v-if="showSuccessMessage"
        class="fixed bottom-4 right-4 bg-green-50 border border-green-200 rounded-md px-4 py-3 shadow-sm"
      >
        <div class="flex items-center">
          <svg class="w-5 h-5 text-green-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <span class="text-sm text-green-800">{{ successMessage }}</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useBoardsStore } from '@/stores/boards'
import { useEditorStore } from '@/stores/editor'
import html2canvas from 'html2canvas'

interface Props {
  boardId: string
  editToken?: string
}

const props = defineProps<Props>()

// Stores
const boardsStore = useBoardsStore()
const editorStore = useEditorStore()

// Refs
const editUrlInput = ref<HTMLInputElement>()
const publicUrlInput = ref<HTMLInputElement>()
const exportDropdown = ref<HTMLDivElement>()

// State
const editUrlCopied = ref(false)
const publicUrlCopied = ref(false)
const showExportMenu = ref(false)
const isExporting = ref(false)
const showSuccessMessage = ref(false)
const successMessage = ref('')

// Computed URLs
const editUrl = computed(() => {
  const baseUrl = window.location.origin
  return `${baseUrl}/board/${props.boardId}?edit=${props.editToken}`
})

const publicUrl = computed(() => {
  const baseUrl = window.location.origin
  return `${baseUrl}/board/${props.boardId}/public?token=${boardsStore.currentBoard?.public_token}`
})

// Copy functionality
const copyEditUrl = async () => {
  try {
    await navigator.clipboard.writeText(editUrl.value)
    editUrlCopied.value = true
    showSuccess('Edit URL copied to clipboard!')
    setTimeout(() => {
      editUrlCopied.value = false
    }, 2000)
  } catch (error) {
    console.error('Failed to copy edit URL:', error)
    // Fallback for browsers without clipboard API
    fallbackCopy(editUrlInput.value, 'Edit URL copied!')
  }
}

const copyPublicUrl = async () => {
  try {
    await navigator.clipboard.writeText(publicUrl.value)
    publicUrlCopied.value = true
    showSuccess('Public URL copied to clipboard!')
    setTimeout(() => {
      publicUrlCopied.value = false
    }, 2000)
  } catch (error) {
    console.error('Failed to copy public URL:', error)
    // Fallback for browsers without clipboard API
    fallbackCopy(publicUrlInput.value, 'Public URL copied!')
  }
}

const fallbackCopy = (input: HTMLInputElement | undefined, message: string) => {
  if (input) {
    input.select()
    input.setSelectionRange(0, 99999) // For mobile devices
    try {
      document.execCommand('copy')
      showSuccess(message)
    } catch (error) {
      console.error('Fallback copy failed:', error)
    }
  }
}

// Export functionality
const exportCurrentPage = async () => {
  if (!boardsStore.currentBoard || !editorStore.currentPageId) {
    return
  }

  showExportMenu.value = false
  isExporting.value = true

  try {
    const currentPage = boardsStore.sortedPages.find(p => p.id === editorStore.currentPageId)
    const pageName = currentPage?.title || 'page'
    const boardName = boardsStore.currentBoard.title || 'board'
    const fileName = `${boardName}-${pageName}.png`

    await exportCanvasToPng(fileName)
    showSuccess(`${pageName} exported successfully!`)
  } catch (error) {
    console.error('Export failed:', error)
    showSuccess('Export failed. Please try again.')
  } finally {
    isExporting.value = false
  }
}

const exportAllPages = async () => {
  if (!boardsStore.currentBoard || boardsStore.sortedPages.length === 0) {
    return
  }

  showExportMenu.value = false
  isExporting.value = true

  try {
    const currentPageId = editorStore.currentPageId
    const boardName = boardsStore.currentBoard.title || 'board'
    
    for (let i = 0; i < boardsStore.sortedPages.length; i++) {
      const page = boardsStore.sortedPages[i]
      
      // Switch to the page if not current
      if (page.id !== editorStore.currentPageId) {
        editorStore.setCurrentPage(page.id)
        await editorStore.loadPageElements(page.id)
        // Wait for canvas to update
        await new Promise(resolve => setTimeout(resolve, 500))
      }
      
      const fileName = `${boardName}-${page.title || `page-${i + 1}`}.png`
      await exportCanvasToPng(fileName)
    }
    
    // Restore original page
    if (currentPageId && currentPageId !== editorStore.currentPageId) {
      editorStore.setCurrentPage(currentPageId)
      await editorStore.loadPageElements(currentPageId)
    }
    
    showSuccess(`All ${boardsStore.sortedPages.length} pages exported successfully!`)
  } catch (error) {
    console.error('Export all failed:', error)
    showSuccess('Export failed. Please try again.')
  } finally {
    isExporting.value = false
  }
}

const exportCanvasToPng = async (fileName: string) => {
  const canvasElement = document.querySelector('#fabric-canvas') as HTMLCanvasElement
  if (!canvasElement) {
    throw new Error('Canvas not found')
  }

  // Use html2canvas to capture the canvas with high resolution
  const canvas = await html2canvas(canvasElement, {
    scale: 2, // 2x resolution for better quality
    useCORS: true,
    allowTaint: true,
    backgroundColor: '#ffffff',
    logging: false
  })

  // Convert to blob and download
  canvas.toBlob((blob) => {
    if (blob) {
      const url = URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = fileName
      link.click()
      URL.revokeObjectURL(url)
    }
  }, 'image/png', 1.0)
}

// Success message handling
const showSuccess = (message: string) => {
  successMessage.value = message
  showSuccessMessage.value = true
  setTimeout(() => {
    showSuccessMessage.value = false
  }, 3000)
}

// Click outside to close export menu
const handleClickOutside = (event: Event) => {
  if (exportDropdown.value && !exportDropdown.value.contains(event.target as Node)) {
    showExportMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease;
}

.slide-up-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
