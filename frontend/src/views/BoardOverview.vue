<template>
  <div class="board-overview">
    <!-- Header -->
    <header class="overview-header">
      <div class="header-content">
        <h1 class="overview-title">My Boards</h1>
        <div class="header-actions">
          <button 
            v-if="boards.length > 0"
            @click="showDeleteAllModal = true"
            class="delete-all-btn"
            :disabled="loading"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
            </svg>
            Delete All
          </button>
          <button 
            @click="createNewBoard"
            class="create-board-btn"
            :disabled="loading"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
            </svg>
            Create New Board
          </button>
        </div>
      </div>
    </header>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p class="loading-text">Loading boards...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="error-container">
      <div class="error-icon">⚠️</div>
      <h3 class="error-title">Failed to load boards</h3>
      <p class="error-message">{{ error }}</p>
      <button @click="loadBoards" class="retry-btn">
        Try Again
      </button>
    </div>

    <!-- Empty State -->
    <div v-else-if="boards.length === 0" class="empty-container">
      <div class="empty-icon">📋</div>
      <h3 class="empty-title">No boards yet</h3>
      <p class="empty-message">Create your first board to get started</p>
      <button @click="createNewBoard" class="create-first-board-btn">
        Create Your First Board
      </button>
    </div>

    <!-- Boards Grid -->
    <div v-else class="boards-grid">
      <div
        v-for="board in boards"
        :key="board.id"
        @click="navigateToBoard(board)"
        class="board-card"
        :class="{ 'loading': boardLoading === board.id }"
      >
        <!-- Board Preview/Thumbnail -->
        <div class="board-thumbnail">
          <div class="thumbnail-placeholder">
            <svg class="thumbnail-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                    d="M9 17V7m0 10a2 2 0 01-2 2H5a2 2 0 01-2-2V7a2 2 0 012-2h2a2 2 0 012 2m0 10a2 2 0 002 2h2a2 2 0 002-2M9 7a2 2 0 012-2h2a2 2 0 012 2m0 10V7m0 10a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2h2a2 2 0 002-2z">
              </path>
            </svg>
          </div>
          <div v-if="boardLoading === board.id" class="loading-overlay">
            <div class="small-spinner"></div>
          </div>
        </div>

        <!-- Board Info -->
        <div class="board-info">
          <h3 class="board-title">{{ board.title }}</h3>
          <p class="board-description">{{ board.description || 'No description' }}</p>
          
          <div class="board-meta">
            <span class="board-pages">
              {{ board.pageCount || 0 }} {{ board.pageCount === 1 ? 'page' : 'pages' }}
            </span>
            <span class="board-updated">
              Updated {{ formatDate(board.updated_at) }}
            </span>
          </div>
        </div>

        <!-- Board Actions -->
        <div class="board-actions" @click.stop>
          <button 
            @click="duplicateBoard(board)"
            class="action-btn"
            title="Duplicate board"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z">
              </path>
            </svg>
          </button>
          <button 
            @click="deleteBoard(board)"
            class="action-btn delete-btn"
            title="Delete board"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                    d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16">
              </path>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Create Board Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal-container" @click.stop>
        <div class="modal-header">
          <div class="modal-icon">
            <svg class="w-12 h-12 text-primary-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
            </svg>
          </div>
          <h3 class="modal-title">Create New Board</h3>
          <p class="modal-description">
            Give your new board a name and description to get started with your creative journey.
          </p>
        </div>
        
        <div class="modal-content">
          <div class="form-group">
            <label for="board-title" class="form-label">Board Name *</label>
            <input
              id="board-title"
              v-model="newBoardTitle"
              type="text"
              class="form-input"
              placeholder="Enter board name..."
              maxlength="100"
              @keyup.enter="confirmCreateBoard"
              :disabled="createLoading"
            />
          </div>
          
          <div class="form-group">
            <label for="board-description" class="form-label">Description (Optional)</label>
            <textarea
              id="board-description"
              v-model="newBoardDescription"
              class="form-textarea"
              placeholder="Describe what this board is for..."
              rows="3"
              maxlength="500"
              :disabled="createLoading"
            ></textarea>
          </div>
          
          <div class="character-count">
            <span class="count-text">{{ newBoardTitle.length }}/100</span>
          </div>
        </div>
        
        <div class="modal-actions">
          <button 
            @click="showCreateModal = false" 
            class="modal-btn modal-btn-secondary"
            :disabled="createLoading"
          >
            Cancel
          </button>
          <button 
            @click="confirmCreateBoard" 
            class="modal-btn modal-btn-primary"
            :disabled="createLoading || !newBoardTitle.trim()"
          >
            <svg v-if="createLoading" class="animate-spin w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
            </svg>
            {{ createLoading ? 'Creating...' : 'Create Board' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete All Modal -->
    <div v-if="showDeleteAllModal" class="modal-overlay" @click="showDeleteAllModal = false">
      <div class="modal-container" @click.stop>
        <div class="modal-header">
          <div class="modal-icon">
            <svg class="w-12 h-12 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.5 0L4.268 6.5c-.77.833-.192 2.5 1.732 2.5z"></path>
            </svg>
          </div>
          <h3 class="modal-title">Delete All Boards</h3>
          <p class="modal-description">
            Are you absolutely sure you want to delete all {{ boards.length }} boards? This action cannot be undone and will permanently remove all your boards and their content.
          </p>
        </div>
        
        <div class="modal-content">
          <div class="boards-preview">
            <div class="preview-title">Boards to be deleted:</div>
            <div class="marquee">
              <div class="marquee-content">
                <span v-for="board in boards" :key="board.id" class="board-chip">
                  {{ board.title }}
                </span>
                <!-- Duplicate for seamless loop -->
                <span v-for="board in boards" :key="`duplicate-${board.id}`" class="board-chip">
                  {{ board.title }}
                </span>
              </div>
            </div>
          </div>
        </div>
        
        <div class="modal-actions">
          <button 
            @click="showDeleteAllModal = false" 
            class="modal-btn modal-btn-secondary"
            :disabled="deleteAllLoading"
          >
            Cancel
          </button>
          <button 
            @click="deleteAllBoards" 
            class="modal-btn modal-btn-danger"
            :disabled="deleteAllLoading"
          >
            <svg v-if="deleteAllLoading" class="animate-spin w-4 h-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
            </svg>
            {{ deleteAllLoading ? 'Deleting...' : 'Delete All Boards' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { boardsApi, type CreateBoardRequest } from '@/api/boards'
import type { Board } from '@/types'

const router = useRouter()

// State
const boards = ref<Board[]>([])
const loading = ref(true)
const error = ref<string | null>(null)
const boardLoading = ref<string | null>(null)
const showDeleteAllModal = ref(false)
const deleteAllLoading = ref(false)
const showCreateModal = ref(false)
const createLoading = ref(false)
const newBoardTitle = ref('')
const newBoardDescription = ref('')

// Load all boards
const loadBoards = async () => {
  loading.value = true
  error.value = null
  
  try {
    console.log('Loading boards...')
    const response = await boardsApi.getAll()
    boards.value = response
    console.log('Loaded boards:', response.length)
  } catch (err: any) {
    console.error('Failed to load boards:', err)
    error.value = err.message || 'Failed to load boards'
  } finally {
    loading.value = false
  }
}

// Navigate to a specific board
const navigateToBoard = async (board: Board) => {
  if (boardLoading.value) return
  
  boardLoading.value = board.id
  
  try {
    // Navigate to the board editor with the edit token
    await router.push(`/board/${board.id}/edit?edit_token=${board.edit_token}`)
  } catch (err) {
    console.error('Failed to navigate to board:', err)
  } finally {
    boardLoading.value = null
  }
}

// Create a new board
const createNewBoard = () => {
  showCreateModal.value = true
  newBoardTitle.value = ''
  newBoardDescription.value = ''
}

// Confirm create new board
const confirmCreateBoard = async () => {
  if (!newBoardTitle.value.trim()) {
    return
  }
  
  createLoading.value = true
  
  try {
    const createData: CreateBoardRequest = {
      title: newBoardTitle.value.trim()
    }
    
    if (newBoardDescription.value.trim()) {
      createData.description = newBoardDescription.value.trim()
    }
    
    const newBoardResponse = await boardsApi.create(createData)
    
    showCreateModal.value = false
    
    // Navigate to the new board using the edit_url
    const editUrl = new URL(newBoardResponse.edit_url)
    await router.push(editUrl.pathname + editUrl.search)
  } catch (err: any) {
    console.error('Failed to create board:', err)
    error.value = err.message || 'Failed to create board'
  } finally {
    createLoading.value = false
  }
}

// Duplicate a board
const duplicateBoard = async (board: Board) => {
  try {
    await boardsApi.create({
      title: `${board.title} (Copy)`,
      description: board.description || ''
    })
    
    // Refresh the boards list
    await loadBoards()
  } catch (err: any) {
    console.error('Failed to duplicate board:', err)
    error.value = err.message || 'Failed to duplicate board'
  }
}

// Delete a board
const deleteBoard = async (board: Board) => {
  if (!confirm(`Are you sure you want to delete "${board.title}"? This action cannot be undone.`)) {
    return
  }
  
  try {
    console.log('Deleting board:', board.id)
    await boardsApi.delete(board.id, board.edit_token)
    console.log('Deleted board:', board.id)
    
    // Remove from local state
    boards.value = boards.value.filter(b => b.id !== board.id)
  } catch (err: any) {
    console.error('Failed to delete board:', err)
    error.value = err.message || 'Failed to delete board'
  }
}

// Delete all boards
const deleteAllBoards = async () => {
  deleteAllLoading.value = true
  
  try {
    await boardsApi.deleteAll(boards.value)
    boards.value = []
    showDeleteAllModal.value = false
  } catch (err: any) {
    console.error('Failed to delete all boards:', err)
    error.value = err.message || 'Failed to delete all boards'
  } finally {
    deleteAllLoading.value = false
  }
}

// Format date for display
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = now.getTime() - date.getTime()
  const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    return 'Today'
  } else if (diffDays === 1) {
    return 'Yesterday'
  } else if (diffDays < 7) {
    return `${diffDays} days ago`
  } else {
    return date.toLocaleDateString()
  }
}

// Load boards on mount
onMounted(() => {
  loadBoards()
})
</script>

<style scoped>
.board-overview {
  min-height: 100vh;
  background: linear-gradient(135deg, #fef7ee 0%, #fdedd6 20%, #fad7ac 100%);
  padding: 2rem;
}

.overview-header {
  margin-bottom: 3rem;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.header-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.overview-title {
  font-size: 2.5rem;
  font-weight: bold;
  color: #762e17;
  margin: 0;
  text-shadow: 0 2px 4px rgba(118, 46, 23, 0.1);
}

.create-board-btn {
  display: flex;
  align-items: center;
  background: #ed7420;
  color: white;
  border: none;
  border-radius: 12px;
  padding: 0.75rem 1.5rem;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(237, 116, 32, 0.25);
}

.create-board-btn:hover {
  transform: translateY(-2px);
  background: #de5a16;
  box-shadow: 0 6px 20px rgba(237, 116, 32, 0.35);
}

.create-board-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.delete-all-btn {
  display: flex;
  align-items: center;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 12px;
  padding: 0.75rem 1.5rem;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.25);
}

.delete-all-btn:hover {
  background: #dc2626;
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(239, 68, 68, 0.35);
}

.delete-all-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.loading-container,
.error-container,
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  max-width: 500px;
  margin: 0 auto;
  text-align: center;
}

.loading-spinner {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(237, 116, 32, 0.2);
  border-top: 4px solid #ed7420;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

.loading-text {
  color: #762e17;
  font-size: 1.1rem;
  margin: 0;
}

.error-icon,
.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.error-title,
.empty-title {
  color: #762e17;
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
}

.error-message,
.empty-message {
  color: #933619;
  font-size: 1rem;
  margin: 0 0 2rem 0;
}

.retry-btn,
.create-first-board-btn {
  background: #ed7420;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 0.75rem 1.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(237, 116, 32, 0.25);
}

.retry-btn:hover,
.create-first-board-btn:hover {
  background: #de5a16;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(237, 116, 32, 0.35);
}

.boards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.board-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(237, 116, 32, 0.1);
  border: 1px solid #fad7ac;
  position: relative;
}

.board-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 32px rgba(237, 116, 32, 0.2);
  border-color: #f6ba77;
}

.board-card.loading {
  pointer-events: none;
  opacity: 0.7;
}

.board-thumbnail {
  height: 200px;
  background: linear-gradient(135deg, #fef7ee 0%, #fdedd6 50%, #fad7ac 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.thumbnail-placeholder {
  color: #933619;
}

.thumbnail-icon {
  width: 64px;
  height: 64px;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
}

.small-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid #fad7ac;
  border-top: 2px solid #ed7420;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

.board-info {
  padding: 1.5rem;
}

.board-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 0.5rem 0;
  line-height: 1.3;
}

.board-description {
  color: #6b7280;
  font-size: 0.9rem;
  margin: 0 0 1rem 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.board-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.8rem;
  color: #9ca3af;
}

.board-pages {
  font-weight: 500;
}

.board-updated {
  font-style: italic;
}

.board-actions {
  position: absolute;
  top: 1rem;
  right: 1rem;
  display: flex;
  gap: 0.5rem;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.board-card:hover .board-actions {
  opacity: 1;
}

.action-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  background: rgba(255, 255, 255, 0.9);
  color: #6b7280;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  backdrop-filter: blur(4px);
}

.action-btn:hover {
  background: white;
  color: #ed7420;
  transform: scale(1.1);
}

.delete-btn:hover {
  color: #ef4444;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
  animation: fade-in 0.2s ease-out;
}

.modal-container {
  background: white;
  border-radius: 24px;
  max-width: 600px;
  width: 90%;
  max-height: 80vh;
  overflow: hidden;
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
  animation: slide-up 0.3s ease-out;
}

.modal-header {
  padding: 2rem 2rem 1rem 2rem;
  text-align: center;
  border-bottom: 1px solid #fad7ac;
}

.modal-icon {
  margin: 0 auto 1rem auto;
  width: fit-content;
}

.modal-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #762e17;
  margin: 0 0 0.5rem 0;
}

.modal-description {
  color: #933619;
  font-size: 1rem;
  margin: 0;
  line-height: 1.5;
}

.modal-content {
  padding: 1.5rem 2rem;
}

.boards-preview {
  margin-bottom: 1rem;
}

.preview-title {
  font-weight: 600;
  color: #762e17;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

/* Magic UI Marquee Effect */
.marquee {
  width: 100%;
  overflow: hidden;
  background: linear-gradient(135deg, #fef7ee 0%, #fdedd6 100%);
  border-radius: 12px;
  padding: 1rem 0;
  border: 1px solid #fad7ac;
}

.marquee-content {
  display: flex;
  animation: marquee 20s linear infinite;
  gap: 1rem;
}

.board-chip {
  background: white;
  color: #762e17;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 500;
  white-space: nowrap;
  border: 1px solid #f6ba77;
  box-shadow: 0 2px 4px rgba(237, 116, 32, 0.1);
}

.modal-actions {
  padding: 1rem 2rem 2rem 2rem;
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.modal-btn {
  display: flex;
  align-items: center;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  font-weight: 600;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.modal-btn-secondary {
  background: #f6ba77;
  color: #762e17;
}

.modal-btn-secondary:hover {
  background: #f19340;
  transform: translateY(-1px);
}

.modal-btn-danger {
  background: #ef4444;
  color: white;
}

.modal-btn-danger:hover {
  background: #dc2626;
  transform: translateY(-1px);
}

.modal-btn-primary {
  background: #ed7420;
  color: white;
}

.modal-btn-primary:hover {
  background: #de5a16;
  transform: translateY(-1px);
}

.modal-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* Form Styles */
.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  font-weight: 600;
  color: #762e17;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.form-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 2px solid #fad7ac;
  border-radius: 12px;
  font-size: 1rem;
  background: white;
  color: #762e17;
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: #ed7420;
  box-shadow: 0 0 0 3px rgba(237, 116, 32, 0.1);
}

.form-input:disabled {
  background: #fef7ee;
  cursor: not-allowed;
}

.form-textarea {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 2px solid #fad7ac;
  border-radius: 12px;
  font-size: 1rem;
  background: white;
  color: #762e17;
  transition: all 0.2s ease;
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
}

.form-textarea:focus {
  outline: none;
  border-color: #ed7420;
  box-shadow: 0 0 0 3px rgba(237, 116, 32, 0.1);
}

.form-textarea:disabled {
  background: #fef7ee;
  cursor: not-allowed;
}

.character-count {
  text-align: right;
  margin-top: -1rem;
  margin-bottom: 1rem;
}

.count-text {
  font-size: 0.8rem;
  color: #933619;
}

@keyframes marquee {
  0% { transform: translateX(0); }
  100% { transform: translateX(-50%); }
}

@keyframes fade-in {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slide-up {
  from { 
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to { 
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* Responsive design */
@media (max-width: 768px) {
  .board-overview {
    padding: 1rem;
  }
  
  .overview-title {
    font-size: 2rem;
  }
  
  .header-content {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }
  
  .header-actions {
    flex-direction: column;
    width: 100%;
  }
  
  .delete-all-btn,
  .create-board-btn {
    width: 100%;
    justify-content: center;
  }
  
  .boards-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
  
  .modal-container {
    margin: 1rem;
    width: calc(100% - 2rem);
  }
  
  .modal-header,
  .modal-content,
  .modal-actions {
    padding-left: 1.5rem;
    padding-right: 1.5rem;
  }
  
  .modal-actions {
    flex-direction: column;
  }
  
  .modal-btn {
    width: 100%;
    justify-content: center;
  }
  
  .form-input,
  .form-textarea {
    font-size: 16px; /* Prevents zoom on iOS */
  }
}
</style>
