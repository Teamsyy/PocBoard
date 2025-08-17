import apiClient from './client'
import type { RecapData, DateFilter, ApiResponse } from '@/types'

export interface RecapParams {
  filter: DateFilter
  date?: string // ISO date string for the reference date
}

export const recapApi = {
  // Get recap data for a board
  async getRecap(boardId: string, params: RecapParams): Promise<RecapData> {
    const response = await apiClient.get<ApiResponse<RecapData>>(
      `/boards/${boardId}/recap`,
      { params }
    )
    
    if (response.data.error) {
      throw response.data
    }
    return response.data.data!
  },
}