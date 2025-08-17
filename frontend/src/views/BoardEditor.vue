<template>
  <div class="h-screen flex flex-col bg-amber-50">
    <!-- Top Navigation Bar -->
    <header class="bg-white border-b border-amber-200 shadow-sm flex-shrink-0">
      <div class="px-4 py-3 flex items-center justify-between">
        <!-- Left: Board Title -->
        <div class="flex items-center space-x-4">
          <router-link 
            to="/" 
            class="text-amber-600 hover:text-amber-800 transition-colors"
            title="Back to Home"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
          </router-link>
          
          <div class="flex items-center space-x-2">
            <h1 class="text-xl font-bold text-amber-800 font-journal">
              {{ boardTitle }}
            </h1>
            <button 
              v-if="isEditMode"
              @click="editingTitle = true"
              class="btn-icon text-amber-600 hover:text-amber-800"
              title="Edit board title"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Mobile Menu Toggle -->
        <button 
          @click="showMobileSidebar = !showMobileSidebar"
          class="btn-icon md:hidden"
          title="Toggle sidebar"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>

        <!-- Right: Sharing Controls -->
        <div class="flex items-center space-x-2">
          <!-- Page Navigation -->
          <div class="flex items-center space-x-1 sm:space-x-2 mr-2 sm:mr-4">
            <button 
              @click="previousPage"
              :disabled="currentPageIndex <= 0"
              class="btn-icon"
              title="Previous page"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            
            <span class="text-xs sm:text-sm text-amber-600 min-w-[40px] sm:min-w-[60px] text-center">
              {{ currentPageIndex + 1 }}/{{ totalPages }}
            </span>
            
            <button 
              @click="nextPage"
              :disabled="currentPageIndex >= totalPages - 1"
              class="btn-icon"
              title="Next page"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>

          <!-- Share Button -->
          <button 
            @click="showShareModal = true"
            class="btn-secondary hidden sm:flex"
            title="Share board"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.367 2.684 3 3 0 00-5.367-2.684z" />
            </svg>
            Share
          </button>

          <!-- Mobile Share Button -->
          <button 
            @click="showShareModal = true"
            class="btn-icon sm:hidden"
            title="Share board"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.367 2.684 3 3 0 00-5.367-2.684z" />
            </svg>
          </button>

          <!-- Export Button -->
          <button 
            @click="exportPage"
            class="btn-ghost hidden sm:flex"
            title="Export current page as PNG"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Export
          </button>

          <!-- Mobile Export Button -->
          <button 
            @click="exportPage"
            class="btn-icon sm:hidden"
            title="Export current page as PNG"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </button>
        </div>
      </div>
    </header>

    <!-- Toolbar -->
    <Toolbar 
      :is-edit-mode="isEditMode"
      @add-text="addTextElement"
      @add-shape="addShapeElement"
      @add-sticker="addStickerElement"
    />

    <!-- Main Content Area -->
    <div class="flex-1 flex overflow-hidden">
      <!-- Left Sidebar -->
      <aside class="sidebar w-64 flex-shrink-0 overflow-y-auto hidden md:block">
        <!-- Tools Section -->
        <div class="sidebar-header">
          <h2 class="text-sm font-semibold text-amber-800">Tools</h2>
        </div>
        <div class="sidebar-content">
          <!-- Element Creation Tools -->
          <div class="space-y-2">
            <button 
              @click="addTextElement"
              class="w-full btn-ghost justify-start"
              :disabled="!isEditMode"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
              </svg>
              Add Text
            </button>
            
            <!-- Image Upload Component -->
            <div class="w-full">
              <ImageUploader :auto-create="true" variant="sidebar" />
            </div>
            
            <button 
              @click="addShapeElement"
              class="w-full btn-ghost justify-start"
              :disabled="!isEditMode"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              Add Shape
            </button>
            
            <button 
              @click="addStickerElement"
              class="w-full btn-ghost justify-start"
              :disabled="!isEditMode"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1.01M15 10h1.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Add Sticker
            </button>
          </div>

          <!-- Canvas Controls -->
          <div class="mt-6 pt-4 border-t border-amber-200">
            <h3 class="text-xs font-semibold text-amber-700 uppercase tracking-wide mb-3">
              Canvas
            </h3>
            <div class="space-y-2">
              <label class="flex items-center">
                <input 
                  v-model="snapToGrid" 
                  type="checkbox" 
                  class="rounded border-amber-300 text-primary-500 focus:ring-primary-500"
                >
                <span class="ml-2 text-sm text-amber-700">Snap to Grid</span>
              </label>
              
              <div class="flex space-x-2">
                <button 
                  @click="undo"
                  :disabled="!canUndo"
                  class="btn-icon flex-1"
                  title="Undo"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                  </svg>
                </button>
                
                <button 
                  @click="redo"
                  :disabled="!canRedo"
                  class="btn-icon flex-1"
                  title="Redo"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 10h-10a8 8 0 00-8 8v2m18-10l-6 6m6-6l-6-6" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Theme Selection -->
        <div class="sidebar-header border-t border-amber-200">
          <h2 class="text-sm font-semibold text-amber-800">Themes</h2>
        </div>
        <div class="sidebar-content">
          <div class="grid grid-cols-2 gap-2">
            <button 
              v-for="theme in availableThemes"
              :key="theme.name"
              @click="applyTheme(theme)"
              :class="[
                'aspect-square rounded-lg border-2 transition-colors',
                currentTheme === theme.name 
                  ? 'border-primary-500 ring-2 ring-primary-200' 
                  : 'border-amber-200 hover:border-amber-300'
              ]"
              :style="{ backgroundImage: `url(${theme.preview})`, backgroundSize: 'cover' }"
              :title="theme.name"
              :disabled="!isEditMode"
            >
              <span class="sr-only">{{ theme.name }}</span>
            </button>
          </div>
        </div>
      </aside>

      <!-- Main Canvas Area -->
      <main class="flex-1 flex flex-col overflow-hidden bg-amber-50">
        <!-- Canvas Container -->
        <div class="flex-1 p-4 overflow-auto">
          <div class="canvas-container mx-auto max-w-4xl h-full min-h-[600px]" :style="canvasStyle">
            <CanvasEditor
              v-if="currentPage"
              ref="canvasEditor"
              :width="800"
              :height="600"
              class="w-full h-full"
            />
            <div v-else class="w-full h-full flex items-center justify-center text-amber-600">
              <div class="text-center">
                <svg class="w-16 h-16 mx-auto mb-4 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                </svg>
                <p class="text-lg font-medium">No Page Selected</p>
                <p class="text-sm mt-1">Create or select a page to start editing</p>
              </div>
            </div>
          </div>
        </div>
      </main>

      <!-- Right Sidebar - Layers -->
      <aside v-if="isEditMode" class="w-64 flex-shrink-0 bg-white border-l border-gray-200 hidden lg:block">
        <SidebarLayers />
      </aside>
    </div>

    <!-- Mobile Sidebar Overlay -->
    <div 
      v-if="showMobileSidebar"
      class="fixed inset-0 bg-black bg-opacity-50 z-40 md:hidden"
      @click="showMobileSidebar = false"
    >
      <aside class="sidebar w-64 h-full overflow-y-auto" @click.stop>
        <!-- Tools Section -->
        <div class="sidebar-header">
          <div class="flex items-center justify-between">
            <h2 class="text-sm font-semibold text-amber-800">Tools</h2>
            <button 
              @click="showMobileSidebar = false"
              class="btn-icon"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
        <div class="sidebar-content">
          <!-- Element Creation Tools -->
          <div class="space-y-2">
            <button 
              @click="addTextElement"
              class="w-full btn-ghost justify-start"
              :disabled="!isEditMode"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h7" />
              </svg>
              Add Text
            </button>
            
            <!-- Mobile Image Upload Component -->
            <div class="w-full">
              <ImageUploader :auto-create="true" variant="sidebar" />
            </div>
            
            <button 
              @click="addShapeElement"
              class="w-full btn-ghost justify-start"
              :disabled="!isEditMode"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              Add Shape
            </button>
            
            <button 
              @click="addStickerElement"
              class="w-full btn-ghost justify-start"
              :disabled="!isEditMode"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1.01M15 10h1.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Add Sticker
            </button>
          </div>

          <!-- Canvas Controls -->
          <div class="mt-6 pt-4 border-t border-amber-200">
            <h3 class="text-xs font-semibold text-amber-700 uppercase tracking-wide mb-3">
              Canvas
            </h3>
            <div class="space-y-2">
              <label class="flex items-center">
                <input 
                  v-model="snapToGrid" 
                  type="checkbox" 
                  class="rounded border-amber-300 text-primary-500 focus:ring-primary-500"
                >
                <span class="ml-2 text-sm text-amber-700">Snap to Grid</span>
              </label>
              
              <div class="flex space-x-2">
                <button 
                  @click="undo"
                  :disabled="!canUndo"
                  class="btn-icon flex-1"
                  title="Undo"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                  </svg>
                </button>
                
                <button 
                  @click="redo"
                  :disabled="!canRedo"
                  class="btn-icon flex-1"
                  title="Redo"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 10h-10a8 8 0 00-8 8v2m18-10l-6 6m6-6l-6-6" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Theme Selection -->
        <div class="sidebar-header border-t border-amber-200">
          <h2 class="text-sm font-semibold text-amber-800">Themes</h2>
        </div>
        <div class="sidebar-content">
          <div class="grid grid-cols-2 gap-2">
            <button 
              v-for="theme in availableThemes"
              :key="theme.name"
              @click="applyTheme(theme)"
              :class="[
                'aspect-square rounded-lg border-2 transition-colors',
                currentTheme === theme.name 
                  ? 'border-primary-500 ring-2 ring-primary-200' 
                  : 'border-amber-200 hover:border-amber-300'
              ]"
              :style="{ backgroundImage: `url(${theme.preview})`, backgroundSize: 'cover' }"
              :title="theme.name"
              :disabled="!isEditMode"
            >
              <span class="sr-only">{{ theme.name }}</span>
            </button>
          </div>
        </div>
      </aside>
    </div>

    <!-- Share Modal (placeholder) -->
    <div 
      v-if="showShareModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="showShareModal = false"
    >
      <div 
        class="bg-white rounded-xl p-6 max-w-md w-full mx-4"
        @click.stop
      >
        <h3 class="text-lg font-semibold text-amber-800 mb-4">Share Board</h3>
        <p class="text-amber-600 text-sm mb-4">
          Share functionality will be implemented in later tasks
        </p>
        <div class="flex justify-end">
          <button @click="showShareModal = false" class="btn-secondary">
            Close
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Share and Export Bar -->
  <ShareExportBar 
    :board-id="boardId"
    :edit-token="editToken"
  />

  <!-- Save Status Indicator -->
  <SaveStatus />
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useBoardsStore } from '@/stores/boards'
import { useEditorStore } from '@/stores/editor'
import { CanvasEditor, Toolbar, SidebarLayers, ImageUploader, SaveStatus, ShareExportBar } from '@/components'
import type { Theme } from '@/types'

interface Props {
  boardId: string
  editToken?: string
}

const props = defineProps<Props>()
const route = useRoute()

// Stores
const boardsStore = useBoardsStore()
const editorStore = useEditorStore()

// Local state
const editingTitle = ref(false)
const showShareModal = ref(false)
const showMobileSidebar = ref(false)
const currentPageIndex = ref(0)
const canvasEditor = ref<InstanceType<typeof CanvasEditor>>()

// Available themes (placeholder data)
const availableThemes = ref<Theme[]>([
  { name: 'Default', skin: 'default', preview: '/skins/default-preview.svg' },
  { name: 'Wood', skin: 'wood', preview: '/skins/wood-preview.svg' },
  { name: 'Notebook', skin: 'notebook', preview: '/skins/notebook-preview.svg' },
  { name: 'Cork', skin: 'cork', preview: '/skins/cork-preview.svg' },
])

// Computed properties
const boardTitle = computed(() => boardsStore.currentBoard?.title || 'Untitled Board')
const isEditMode = computed(() => boardsStore.isEditMode)
const totalPages = computed(() => boardsStore.sortedPages.length || 1)
const currentPage = computed(() => boardsStore.sortedPages[currentPageIndex.value])
const currentTheme = computed(() => boardsStore.currentBoard?.skin || 'default')
const editToken = computed(() => props.editToken || route.query.edit as string)

// Canvas styling based on theme
const canvasStyle = computed(() => {
  const theme = availableThemes.value.find(t => t.name.toLowerCase() === currentTheme.value)
  if (theme && theme.preview) {
    return {
      backgroundImage: `url(${theme.preview})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
      backgroundRepeat: 'no-repeat'
    }
  }
  return {}
})

// Editor state
const snapToGrid = computed({
  get: () => editorStore.snapToGrid,
  set: (_value) => editorStore.toggleSnapToGrid()
})

const canUndo = computed(() => editorStore.canUndo)
const canRedo = computed(() => editorStore.canRedo)

// Methods
const loadBoard = async () => {
  try {
    console.log('Loading board...')
    
    // Set tokens from props or URL
    if (props.editToken) {
      boardsStore.editToken = props.editToken
      console.log('Using edit token from props:', props.editToken)
    } else {
      boardsStore.setTokensFromUrl()
      console.log('Edit token from URL:', boardsStore.editToken)
    }
    
    await boardsStore.loadBoard()
    console.log('Board loaded:', boardsStore.currentBoard?.title)
    console.log('Pages found:', boardsStore.sortedPages.length)
    console.log('Is edit mode:', boardsStore.isEditMode)
    
    // Create a default page if none exist (for new boards)
    if (boardsStore.sortedPages.length === 0 && boardsStore.isEditMode) {
      console.log('Creating default page...')
      try {
        // Ensure pages array is properly initialized
        if (!Array.isArray(boardsStore.pages)) {
          console.log('Pages not properly initialized, reinitializing...')
          boardsStore.pages = []
        }
        
        const today = new Date().toISOString().split('T')[0] // YYYY-MM-DD format
        const newPage = await boardsStore.createPage('Page 1', today, 0)
        console.log('Default page created:', newPage.id)
      } catch (error) {
        console.error('Failed to create default page:', error)
      }
    }
    
    // Load first page if available
    if (boardsStore.sortedPages.length > 0) {
      console.log('Switching to first page')
      currentPageIndex.value = 0
      await switchToPage()
      
      // Double-check that the page was selected
      if (!editorStore.currentPageId) {
        console.log('Page selection failed, trying again')
        const firstPage = boardsStore.sortedPages[0]
        editorStore.setCurrentPage(firstPage.id)
        await editorStore.loadPageElements(firstPage.id)
      }
    } else {
      console.warn('No pages available after board load')
    }
  } catch (error) {
    console.error('Failed to load board:', error)
  }
}

const previousPage = () => {
  if (currentPageIndex.value > 0) {
    currentPageIndex.value--
    switchToPage()
  }
}

const nextPage = () => {
  if (currentPageIndex.value < totalPages.value - 1) {
    currentPageIndex.value++
    switchToPage()
  }
}

const switchToPage = async () => {
  const page = boardsStore.sortedPages[currentPageIndex.value]
  if (page) {
    console.log('Switching to page:', page.id)
    editorStore.setCurrentPage(page.id)
    await editorStore.loadPageElements(page.id)
    
    // The watchers will handle canvas reload automatically
    await nextTick()
  }
}

// Tool methods
const addTextElement = () => {
  if (!isEditMode.value || !canvasEditor.value) return
  canvasEditor.value.addTextElement()
  showMobileSidebar.value = false
}

const addShapeElement = () => {
  if (!isEditMode.value || !canvasEditor.value) return
  canvasEditor.value.addShapeElement('rectangle')
  showMobileSidebar.value = false
}

const addStickerElement = () => {
  if (!isEditMode.value) return
  // For now, we'll add a placeholder sticker
  // In later tasks, this will open a sticker picker
  const placeholderUrl = 'https://via.placeholder.com/80x80/F59E0B/FFFFFF?text=😊'
  if (canvasEditor.value) {
    canvasEditor.value.addStickerElement(placeholderUrl, 'emoji', 'faces')
  }
  showMobileSidebar.value = false
}

const undo = () => {
  editorStore.undo()
}

const redo = () => {
  editorStore.redo()
}

const applyTheme = async (theme: Theme) => {
  if (!isEditMode.value) return
  
  try {
    await boardsStore.updateBoard({ skin: theme.skin })
  } catch (error) {
    console.error('Failed to apply theme:', error)
  }
}

const exportPage = () => {
  console.log('Export page - will be implemented in later tasks')
}

// Watch for route changes
watch(() => route.query.edit_token, () => {
  loadBoard()
})

// Watch for canvas editor availability and ensure elements are loaded
watch(() => canvasEditor.value, async (newCanvas) => {
  if (newCanvas && editorStore.currentPageId) {
    console.log('Canvas editor is now available, ensuring elements are loaded')
    await nextTick()
    // Only trigger reload if canvas is not already loading elements
    setTimeout(() => {
      if (editorStore.elements.length === 0) {
        const event = new CustomEvent('editor:reload-elements')
        window.dispatchEvent(event)
      }
    }, 50)
  }
})

// Initialize
onMounted(() => {
  loadBoard()
})
</script>