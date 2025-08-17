import apiClient from './client'
import type { UploadResponse, ApiResponse } from '@/types'

export const uploadsApi = {
  // Upload an image file
  async uploadImage(boardId: string, file: File): Promise<UploadResponse> {
    const formData = new FormData()
    formData.append('file', file)

    const response = await apiClient.post<ApiResponse<UploadResponse>>(
      `/boards/${boardId}/upload`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    )
    
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },
}