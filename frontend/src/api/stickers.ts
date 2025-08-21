import apiClient from './client'
import type { UploadResponse, ApiResponse } from '@/types'

export interface StickerUploadData {
  url: string
  stickerType: string
  category: string
}

export class StickerUploadService {
  static async uploadSticker(
    boardId: string,
    file: File,
    category: string = 'custom'
  ): Promise<StickerUploadData> {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('category', category)

    // Use the same apiClient as other uploads - it handles token authentication
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
      throw new Error(response.data.error.message || 'Upload failed')
    }

    const uploadData: UploadResponse = response.data.data!

    // Generate a sticker type from the filename
    const stickerType = file.name
      .toLowerCase()
      .replace(/\.[^/.]+$/, '') // Remove file extension
      .replace(/[^a-z0-9]/g, '-') // Replace non-alphanumeric with hyphens
      .replace(/-+/g, '-') // Replace multiple hyphens with single
      .replace(/^-|-$/g, '') // Remove leading/trailing hyphens

    return {
      url: uploadData.url,
      stickerType: stickerType || 'custom-sticker',
      category
    }
  }
}
