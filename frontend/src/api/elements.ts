import apiClient from './client'
import type { Element, ElementKind, ElementPayload, ApiResponse } from '@/types'

export interface CreateElementRequest {
  kind: ElementKind
  x: number
  y: number
  w: number
  h: number
  rotation?: number
  z?: number
  payload: ElementPayload
}

export interface UpdateElementRequest {
  x?: number
  y?: number
  w?: number
  h?: number
  rotation?: number
  z?: number
  payload?: ElementPayload
}

export interface BatchUpdateZIndexRequest {
  updates: Array<{
    id: string
    z: number
  }>
}

export const elementsApi = {
  // Create a new element
  async create(boardId: string, pageId: string, data: CreateElementRequest): Promise<Element> {
    const response = await apiClient.post<ApiResponse<Element>>(
      `/boards/${boardId}/pages/${pageId}/elements`,
      data
    )
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Update an element
  async update(
    boardId: string,
    pageId: string,
    elementId: string,
    data: UpdateElementRequest
  ): Promise<Element> {
    const response = await apiClient.put<ApiResponse<Element>>(
      `/boards/${boardId}/pages/${pageId}/elements/${elementId}`,
      data
    )
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Delete an element
  async delete(boardId: string, pageId: string, elementId: string): Promise<void> {
    const response = await apiClient.delete<ApiResponse<void>>(
      `/boards/${boardId}/pages/${pageId}/elements/${elementId}`
    )
    if (response.data.error) {
      throw response.data
    }
  },

  // Batch update z-index for multiple elements
  async batchUpdateZIndex(
    boardId: string,
    pageId: string,
    data: BatchUpdateZIndexRequest
  ): Promise<void> {
    const response = await apiClient.put<ApiResponse<void>>(
      `/boards/${boardId}/pages/${pageId}/elements/reorder`,
      data
    )
    if (response.data.error) {
      throw response.data
    }
  },
}