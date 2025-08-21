<template>
  <div class="save-status-container">
    <!-- Save Status Indicator -->
    <Transition name="fade" mode="out-in">
      <div
        v-if="shouldShowStatus"
        :key="status"
        class="save-status flex items-center gap-2 px-3 py-2 rounded-full text-sm font-medium shadow-sm border"
        :class="statusClasses"
      >
        <!-- Icon -->
        <div class="save-icon" :class="iconClasses">
          <!-- Saving spinner -->
          <svg
            v-if="status === 'saving'"
            class="animate-spin h-4 w-4"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            />
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
          
          <!-- Saved checkmark -->
          <svg
            v-else-if="status === 'saved'"
            class="h-4 w-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M5 13l4 4L19 7"
            />
          </svg>
          
          <!-- Error exclamation -->
          <svg
            v-else-if="status === 'error'"
            class="h-4 w-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
            />
          </svg>
        </div>
        
        <!-- Status Text -->
        <span class="save-text">{{ statusText }}</span>
        
        <!-- Pending saves count -->
        <span v-if="pendingCount > 0" class="save-count text-xs opacity-75">
          ({{ pendingCount }})
        </span>
      </div>
    </Transition>
    
    <!-- Detailed Error Message -->
    <Transition name="fade">
      <div
        v-if="status === 'error' && saveError"
        class="save-error-details mt-2 p-3 bg-red-50 border border-red-200 rounded-md"
      >
        <div class="flex items-start gap-2">
          <svg
            class="h-5 w-5 text-red-400 mt-0.5 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
            />
          </svg>
          <div class="flex-1">
            <h4 class="text-sm font-medium text-red-800 mb-1">Save Error</h4>
            <p class="text-sm text-red-700">{{ saveError }}</p>
            <button
              @click="retryAll"
              class="mt-2 text-xs text-red-800 hover:text-red-900 underline font-medium"
            >
              Retry All
            </button>
          </div>
          <button
            @click="dismissError"
            class="text-red-400 hover:text-red-600"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useEditorStore } from '@/stores/editor'

const editorStore = useEditorStore()

// Computed properties
const status = computed(() => editorStore.saveStatus)
const saveError = computed(() => editorStore.saveError)
const pendingCount = computed(() => editorStore.pendingSaves.size)
const lastSaveTime = computed(() => editorStore.lastSaveTime)

const shouldShowStatus = computed(() => {
  if (status.value === 'idle') return false
  if (status.value === 'saved') {
    return Date.now() - (lastSaveTime.value?.getTime() || 0) < 3000
  }
  return true
})

const statusClasses = computed(() => {
  switch (status.value) {
    case 'saving':
      return 'bg-blue-50 border-blue-200 text-blue-800'
    case 'saved':
      return 'bg-green-50 border-green-200 text-green-800'
    case 'error':
      return 'bg-red-50 border-red-200 text-red-800'
    default:
      return 'bg-gray-50 border-gray-200 text-gray-600'
  }
})

const iconClasses = computed(() => {
  switch (status.value) {
    case 'saving':
      return 'text-blue-600'
    case 'saved':
      return 'text-green-600'
    case 'error':
      return 'text-red-600'
    default:
      return 'text-gray-400'
  }
})

const statusText = computed(() => {
  switch (status.value) {
    case 'saving':
      return pendingCount.value > 1 ? 'Saving changes...' : 'Saving...'
    case 'saved':
      return 'All changes saved'
    case 'error':
      return 'Save failed'
    default:
      return ''
  }
})

// Actions
const retryAll = async () => {
  try {
    // Clear error state
    editorStore.saveError = null
    editorStore.saveStatus = 'saving'
    
    // Get all pending element IDs
    const pendingIds = Array.from(editorStore.pendingSaves)
    
    // Retry saving each pending element
    for (const elementId of pendingIds) {
      const element = editorStore.elements.find(el => el.id === elementId)
      if (element) {
        await editorStore.saveElement(element)
      }
    }
  } catch (error) {
    console.error('Failed to retry saves:', error)
  }
}

const dismissError = () => {
  editorStore.saveError = null
  if (pendingCount.value === 0) {
    editorStore.saveStatus = 'idle'
  }
}
</script>

<style scoped>
.save-status-container {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 50;
}

.save-status {
  transition: all 0.3s ease-in-out;
  backdrop-filter: blur(8px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
