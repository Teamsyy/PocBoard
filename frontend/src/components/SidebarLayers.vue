<template>
  <div class="w-64 bg-white border-r border-gray-200 h-full flex flex-col">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-gray-200">
      <h2 class="text-sm font-semibold text-gray-900">Layers</h2>
      <p class="text-xs text-gray-500 mt-1">{{ elements.length }} elements</p>
    </div>

    <!-- Layer List -->
    <div class="flex-1 overflow-y-auto">
      <div class="p-2 space-y-1">
        <div
          v-for="(element, index) in sortedElements"
          :key="element.id"
          class="group relative"
          :class="{
            'bg-blue-50 border-blue-200': isSelected(element.id),
            'hover:bg-gray-50': !isSelected(element.id)
          }"
          draggable="true"
          @dragstart="handleDragStart(element, index, $event)"
          @dragover.prevent="handleDragOver($event, index)"
          @dragleave="handleDragLeave"
          @drop="handleDrop(element, $event)"
          @click="handleElementSelect(element.id, $event)"
        >
          <!-- Drop indicator above -->
          <div
            v-if="dragOverIndex === index && dragPosition === 'above'"
            class="absolute -top-0.5 left-2 right-2 h-0.5 bg-blue-500 rounded"
          ></div>

          <div class="flex items-center px-2 py-2 rounded border border-transparent cursor-pointer">
            <!-- Element Icon -->
            <div class="flex-shrink-0 w-6 h-6 flex items-center justify-center mr-2">
              <component :is="getElementIcon(element.kind)" class="w-4 h-4 text-gray-600" />
            </div>

            <!-- Element Info -->
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium text-gray-900 truncate">
                {{ getElementDisplayName(element) }}
              </div>
              <div class="text-xs text-gray-500">
                {{ element.kind }} • z:{{ element.z }}
              </div>
            </div>

            <!-- Layer Controls -->
            <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
              <!-- Visibility Toggle -->
              <button
                @click.stop="toggleVisibility(element.id)"
                class="p-1 rounded hover:bg-gray-200 transition-colors"
                :title="isVisible(element.id) ? 'Hide element' : 'Show element'"
              >
                <EyeIcon v-if="isVisible(element.id)" class="w-4 h-4 text-gray-600" />
                <EyeSlashIcon v-else class="w-4 h-4 text-gray-400" />
              </button>

              <!-- Lock Toggle -->
              <button
                @click.stop="toggleLock(element.id)"
                class="p-1 rounded hover:bg-gray-200 transition-colors"
                :title="isLocked(element.id) ? 'Unlock element' : 'Lock element'"
              >
                <LockClosedIcon v-if="isLocked(element.id)" class="w-4 h-4 text-gray-600" />
                <LockOpenIcon v-else class="w-4 h-4 text-gray-400" />
              </button>

              <!-- Delete Button -->
              <button
                @click.stop="handleDeleteElement(element.id)"
                class="p-1 rounded hover:bg-red-100 transition-colors"
                title="Delete element"
              >
                <TrashIcon class="w-4 h-4 text-gray-400 hover:text-red-600" />
              </button>
            </div>
          </div>

          <!-- Drop indicator below -->
          <div
            v-if="dragOverIndex === index && dragPosition === 'below'"
            class="absolute -bottom-0.5 left-2 right-2 h-0.5 bg-blue-500 rounded"
          ></div>
        </div>

        <!-- Empty State -->
        <div v-if="elements.length === 0" class="text-center py-8">
          <div class="text-gray-400 mb-2">
            <component :is="'div'" class="w-12 h-12 mx-auto rounded-full bg-gray-100 flex items-center justify-center">
              <LayersIcon class="w-6 h-6" />
            </component>
          </div>
          <p class="text-sm text-gray-500">No elements on this page</p>
          <p class="text-xs text-gray-400 mt-1">Add elements using the toolbar</p>
        </div>
      </div>
    </div>

    <!-- Layer Actions -->
    <div class="border-t border-gray-200 p-3">
      <div class="flex items-center justify-between text-xs text-gray-500">
        <span>{{ selectedElementIds.length }} selected</span>
        <div class="flex items-center space-x-2">
          <button
            v-if="selectedElementIds.length > 0"
            @click="handleDeleteSelected"
            class="px-2 py-1 rounded text-red-600 hover:bg-red-50 transition-colors"
            title="Delete selected elements"
          >
            Delete
          </button>
          <button
            v-if="selectedElementIds.length > 0"
            @click="clearSelection"
            class="px-2 py-1 rounded text-gray-600 hover:bg-gray-100 transition-colors"
            title="Clear selection"
          >
            Clear
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useEditorStore } from '@/stores/editor'
import type { Element } from '@/types'
import {
  EyeIcon,
  EyeSlashIcon,
  LockClosedIcon,
  LockOpenIcon,
  TrashIcon,
  RectangleStackIcon as LayersIcon,
  DocumentTextIcon,
  PhotoIcon,
  FaceSmileIcon,
  Square3Stack3DIcon
} from '@heroicons/vue/24/outline'

// Store
const editorStore = useEditorStore()

// State
const draggedElement = ref<Element | null>(null)
const draggedIndex = ref<number>(-1)
const dragOverIndex = ref<number>(-1)
const dragPosition = ref<'above' | 'below'>('below')
const hiddenElements = ref<Set<string>>(new Set())
const lockedElements = ref<Set<string>>(new Set())

// Computed
const elements = computed(() => editorStore.elements)
const selectedElementIds = computed(() => editorStore.selectedElementIds)

// Sort elements by z-index (top to bottom in sidebar = highest to lowest z)
const sortedElements = computed(() => 
  [...elements.value].sort((a, b) => b.z - a.z)
)

// Methods
const isSelected = (elementId: string) => {
  return selectedElementIds.value.includes(elementId)
}

const isVisible = (elementId: string) => {
  return !hiddenElements.value.has(elementId)
}

const isLocked = (elementId: string) => {
  return lockedElements.value.has(elementId)
}

const getElementIcon = (kind: string) => {
  switch (kind) {
    case 'text':
      return DocumentTextIcon
    case 'image':
      return PhotoIcon
    case 'sticker':
      return FaceSmileIcon
    case 'shape':
      return Square3Stack3DIcon
    default:
      return Square3Stack3DIcon
  }
}

const getElementDisplayName = (element: Element) => {
  switch (element.kind) {
    case 'text': {
      const textPayload = element.payload as { content?: string }
      const textContent = textPayload.content || 'Text'
      return textContent.length > 20 ? textContent.substring(0, 20) + '...' : textContent
    }
    case 'image':
      return 'Image'
    case 'sticker':
      return 'Sticker'
    case 'shape':
      return 'Shape'
    default:
      return 'Element'
  }
}

const handleElementSelect = (elementId: string, event: MouseEvent) => {
  if (event.ctrlKey || event.metaKey) {
    // Multi-select
    editorStore.toggleElementSelection(elementId)
  } else {
    // Single select
    editorStore.selectElement(elementId)
  }
}

const toggleVisibility = (elementId: string) => {
  if (hiddenElements.value.has(elementId)) {
    hiddenElements.value.delete(elementId)
    // Show element on canvas
    editorStore.setElementVisibility(elementId, true)
  } else {
    hiddenElements.value.add(elementId)
    // Hide element on canvas
    editorStore.setElementVisibility(elementId, false)
  }
}

const toggleLock = (elementId: string) => {
  if (lockedElements.value.has(elementId)) {
    lockedElements.value.delete(elementId)
    // Unlock element on canvas
    editorStore.setElementLocked(elementId, false)
  } else {
    lockedElements.value.add(elementId)
    // Lock element on canvas
    editorStore.setElementLocked(elementId, true)
  }
}

const handleDeleteElement = (elementId: string) => {
  editorStore.deleteElement(elementId)
  // Remove from local state
  hiddenElements.value.delete(elementId)
  lockedElements.value.delete(elementId)
}

const handleDeleteSelected = () => {
  const elementsToDelete = [...selectedElementIds.value]
  elementsToDelete.forEach(id => {
    editorStore.deleteElement(id)
    hiddenElements.value.delete(id)
    lockedElements.value.delete(id)
  })
}

const clearSelection = () => {
  editorStore.clearSelection()
}

// Drag and Drop functionality
const handleDragStart = (element: Element, index: number, event: DragEvent) => {
  draggedElement.value = element
  draggedIndex.value = index
  
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'move'
    event.dataTransfer.setData('text/plain', element.id)
  }
}

const handleDrop = (targetElement: Element, event: DragEvent) => {
  event.preventDefault()
  
  if (!draggedElement.value || draggedElement.value.id === targetElement.id) {
    return
  }

  const targetZ = targetElement.z
  
  // Calculate new z-index based on drop position
  let newZ: number
  
  if (dragPosition.value === 'above') {
    // Place above target (higher z-index)
    const elementsAbove = sortedElements.value.filter(el => el.z > targetZ)
    if (elementsAbove.length > 0) {
      newZ = Math.min(...elementsAbove.map(el => el.z)) + 1
    } else {
      newZ = targetZ + 1
    }
  } else {
    // Place below target (lower z-index)
    const elementsBelow = sortedElements.value.filter(el => el.z < targetZ)
    if (elementsBelow.length > 0) {
      newZ = Math.max(...elementsBelow.map(el => el.z)) - 1
    } else {
      newZ = targetZ - 1
    }
  }

  // Update element z-index
  editorStore.updateElement(draggedElement.value.id, { z: newZ })
  
  // Reset drag state
  draggedElement.value = null
  draggedIndex.value = -1
  dragOverIndex.value = -1
}

// Handle drag over to show drop indicators
const handleDragOver = (event: DragEvent, index: number) => {
  event.preventDefault()
  
  if (!event.currentTarget) return
  
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect()
  const midPoint = rect.top + rect.height / 2
  
  dragOverIndex.value = index
  dragPosition.value = event.clientY < midPoint ? 'above' : 'below'
}

const handleDragLeave = () => {
  dragOverIndex.value = -1
}
</script>

<style scoped>
/* Custom scrollbar for the layers list */
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: #f1f5f9;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Drag cursor */
[draggable="true"] {
  cursor: grab;
}

[draggable="true"]:active {
  cursor: grabbing;
}
</style>
