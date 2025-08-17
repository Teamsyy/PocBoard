import apiClient from './client'
import type { Page, ApiResponse } from '@/types'

export interface CreatePageRequest {
  title: string
  date: string
  orderIdx: number
}

export interface UpdatePageRequest {
  title?: string
  date?: string
  orderIdx?: number
}

export const pagesApi = {
  // Create a new page
  async create(boardId: string, data: CreatePageRequest): Promise<Page> {
    const response = await apiClient.post<ApiResponse<Page>>(`/boards/${boardId}/pages`, data)
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Get all pages for a board
  async list(boardId: string): Promise<Page[]> {
    const response = await apiClient.get<ApiResponse<Page[]>>(`/boards/${boardId}/pages`)
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Get a specific page
  async get(boardId: string, pageId: string): Promise<Page> {
    const response = await apiClient.get<ApiResponse<Page>>(`/boards/${boardId}/pages/${pageId}`)
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Update a page
  async update(boardId: string, pageId: string, data: UpdatePageRequest): Promise<Page> {
    const response = await apiClient.put<ApiResponse<Page>>(`/boards/${boardId}/pages/${pageId}`, data)
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Delete a page
  async delete(boardId: string, pageId: string): Promise<void> {
    const response = await apiClient.delete<ApiResponse<void>>(`/boards/${boardId}/pages/${pageId}`)
    if (response.data.error) {
      throw response.data
    }
  },
}