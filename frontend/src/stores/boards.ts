import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
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

  // Watch pages to ensure it's always an array
  watch(pages, (newPages: any) => {
    if (!Array.isArray(newPages)) {
      console.warn('Pages was set to non-array value:', newPages, 'resetting to empty array')
      pages.value = []
    }
  }, { immediate: true })

  // Getters
  const isEditMode = computed(() => !!editToken.value)
  const sortedPages = computed(() => {
    const pagesArray = Array.isArray(pages.value) ? pages.value : []
    return pagesArray.slice().sort((a, b) => a.orderIdx - b.orderIdx)
  })

  // Actions
  const setTokensFromUrl = () => {
    const urlParams = new URLSearchParams(window.location.search)
    editToken.value = urlParams.get('edit_token')
    publicToken.value = urlParams.get('public_token')
  }

  const createBoard = async (title: string) => {
    loading.value = true
    error.value = null

    try {
      const response = await boardsApi.create({ title })
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

      // Always initialize pages as empty array first
      pages.value = []
      console.log('Initialized pages as empty array')

      // Load pages from API (they should come empty for new boards)
      if (editToken.value || publicToken.value) {
        console.log('Loading pages from API...')
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

  const updateBoard = async (updates: { title?: string }) => {
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
      console.log('Loading pages for board:', currentBoard.value.id)
      const loadedPages = await pagesApi.list(currentBoard.value.id)
      console.log('Pages loaded from API:', loadedPages)
      console.log('First page elements:', loadedPages[0]?.elements)
      
      // Ensure pages is always an array
      pages.value = Array.isArray(loadedPages) ? loadedPages : []
      console.log('Pages set to:', pages.value)
      console.log('Total elements in first page:', pages.value[0]?.elements?.length)
      
      return pages.value
    } catch (err: any) {
      console.error('Failed to load pages:', err)
      error.value = err.error?.message || 'Failed to load pages'
      pages.value = [] // Ensure pages is always an array
      throw err
    }
  }

  const createPage = async (title: string, date: string, orderIdx?: number) => {
    if (!currentBoard.value || !editToken.value) {
      throw new Error('No board loaded or edit access required')
    }

    // Ensure pages is an array before accessing it
    if (!Array.isArray(pages.value)) {
      pages.value = []
    }

    const finalOrderIdx = orderIdx ?? pages.value.length

    try {
      console.log('Creating page:', { title, date, orderIdx: finalOrderIdx })
      
      // Convert date string to ISO format for backend
      const dateObj = new Date(date + 'T00:00:00.000Z')
      
      const newPage = await pagesApi.create(currentBoard.value.id, {
        title,
        date: dateObj.toISOString(),
        orderIdx: finalOrderIdx,
      })

      console.log('Page created:', newPage)
      
      // Ensure pages is still an array before pushing
      if (!Array.isArray(pages.value)) {
        pages.value = []
      }
      pages.value.push(newPage)
      
      console.log('Pages after creation:', pages.value)
      return newPage
    } catch (err: any) {
      console.error('Failed to create page:', err)
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
