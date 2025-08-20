import apiClient from './client'
import type { Board, ApiResponse } from '@/types'

export interface CreateBoardRequest {
  title: string
}

export interface UpdateBoardRequest {
  title?: string
}

export interface BoardResponse {
  board: Board
  edit_url: string
  public_url: string
}

export const boardsApi = {
  // Create a new board
  async create(data: CreateBoardRequest): Promise<BoardResponse> {
    const response = await apiClient.post<ApiResponse<BoardResponse>>('/boards', data)
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Get board by edit token
  async getByEditToken(editToken: string): Promise<Board> {
    const response = await apiClient.get<ApiResponse<Board>>(`/boards/edit/${editToken}`)
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Get board by public token (read-only)
  async getByPublicToken(publicToken: string): Promise<Board> {
    const response = await apiClient.get<ApiResponse<Board>>(`/boards/public/${publicToken}`)
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },

  // Update board
  async update(boardId: string, data: UpdateBoardRequest): Promise<Board> {
    const response = await apiClient.put<ApiResponse<Board>>(`/boards/${boardId}`, data)
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },
}