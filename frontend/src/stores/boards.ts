import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { boardsApi, pagesApi } from '@/api'
import type { Board, Page } from '@/types'

export const useBoardsStore = defineStore('boards', () => {
  // State
  const currentBoard = ref<Board | null>(null)
  const pages = ref<Page[]>([])
  const editToken = ref<string | null>(null)
  const publicToken = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const isEditMode = computed(() => !!editToken.value)
  const sortedPages = computed(() => (pages.value || []).slice().sort((a, b) => a.orderIdx - b.orderIdx))

  // Actions
  const setTokensFromUrl = () => {
    const urlParams = new URLSearchParams(window.location.search)
    editToken.value = urlParams.get('edit_token')
    publicToken.value = urlParams.get('public_token')
  }

  const createBoard = async (title: string, skin?: string) => {
    loading.value = true
    error.value = null

    try {
      const createData = skin ? { title, skin } : { title }
      const response = await boardsApi.create(createData)
      currentBoard.value = response.board

      // Extract tokens from board object directly (more reliable)
      editToken.value = response.board.edit_token
      publicToken.value = response.board.public_token

      return response
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to create board'
      throw err
    } finally {
      loading.value = false
    }
  }

  const loadBoard = async () => {
    if (!editToken.value && !publicToken.value) {
      setTokensFromUrl()
    }

    if (!editToken.value && !publicToken.value) {
      error.value = 'No access token found'
      return
    }

    loading.value = true
    error.value = null

    try {
      let board: Board

      if (editToken.value) {
        board = await boardsApi.getByEditToken(editToken.value)
      } else {
        board = await boardsApi.getByPublicToken(publicToken.value!)
      }

      currentBoard.value = board

      // Initialize pages array
      pages.value = []

      // Load pages if board has them
      if (board.pages && Array.isArray(board.pages)) {
        pages.value = board.pages
      } else if (editToken.value || publicToken.value) {
        await loadPages()
      }
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to load board'
      pages.value = [] // Ensure pages is always an array
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateBoard = async (updates: { title?: string; skin?: string }) => {
    if (!currentBoard.value || !editToken.value) {
      throw new Error('No board loaded or edit access required')
    }

    loading.value = true
    error.value = null

    try {
      const updatedBoard = await boardsApi.update(currentBoard.value.id, updates)
      currentBoard.value = updatedBoard
      return updatedBoard
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to update board'
      throw err
    } finally {
      loading.value = false
    }
  }

  const loadPages = async () => {
    if (!currentBoard.value) {
      throw new Error('No board loaded')
    }

    try {
      const loadedPages = await pagesApi.list(currentBoard.value.id)
      pages.value = loadedPages || []
      return loadedPages || []
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to load pages'
      pages.value = [] // Ensure pages is always an array
      throw err
    }
  }

  const createPage = async (title: string, date: string, orderIdx?: number) => {
    if (!currentBoard.value || !editToken.value) {
      throw new Error('No board loaded or edit access required')
    }

    const finalOrderIdx = orderIdx ?? (pages.value || []).length

    try {
      // Convert date string to ISO format for backend
      const dateObj = new Date(date + 'T00:00:00.000Z')
      
      const newPage = await pagesApi.create(currentBoard.value.id, {
        title,
        date: dateObj.toISOString(),
        orderIdx: finalOrderIdx,
      })

      if (!pages.value) {
        pages.value = []
      }
      pages.value.push(newPage)
      return newPage
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to create page'
      throw err
    }
  }

  const updatePage = async (
    pageId: string,
    updates: { title?: string; date?: string; orderIdx?: number }
  ) => {
    if (!currentBoard.value || !editToken.value) {
      throw new Error('No board loaded or edit access required')
    }

    try {
      const updatedPage = await pagesApi.update(currentBoard.value.id, pageId, updates)

      const index = pages.value.findIndex(p => p.id === pageId)
      if (index !== -1) {
        pages.value[index] = updatedPage
      }

      return updatedPage
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to update page'
      throw err
    }
  }

  const deletePage = async (pageId: string) => {
    if (!currentBoard.value || !editToken.value) {
      throw new Error('No board loaded or edit access required')
    }

    try {
      await pagesApi.delete(currentBoard.value.id, pageId)
      pages.value = pages.value.filter(p => p.id !== pageId)
    } catch (err: any) {
      error.value = err.error?.message || 'Failed to delete page'
      throw err
    }
  }

  const clearError = () => {
    error.value = null
  }

  const reset = () => {
    currentBoard.value = null
    pages.value = []
    editToken.value = null
    publicToken.value = null
    loading.value = false
    error.value = null
  }

  return {
    // State
    currentBoard,
    pages,
    editToken,
    publicToken,
    loading,
    error,

    // Getters
    isEditMode,
    sortedPages,

    // Actions
    setTokensFromUrl,
    createBoard,
    loadBoard,
    updateBoard,
    loadPages,
    createPage,
    updatePage,
    deletePage,
    clearError,
    reset,
  }
})
