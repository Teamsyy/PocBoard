import axios from 'axios'
import type { AxiosInstance, AxiosResponse, AxiosError } from 'axios'
import type { ApiResponse } from '@/types'

// Create axios instance with base configuration
const apiClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor to add edit token if available
apiClient.interceptors.request.use(
  (config) => {
    console.log('API Request:', config.method?.toUpperCase(), config.url, config.data)
    
    // Get edit token from URL params or stored state
    const urlParams = new URLSearchParams(window.location.search)
    const editToken = urlParams.get('edit_token')
    
    if (editToken && config.params) {
      config.params.edit_token = editToken
    } else if (editToken) {
      config.params = { edit_token: editToken }
    }
    
    return config
  },
  (error) => {
    console.error('API Request Error:', error)
    return Promise.reject(error)
  }
)

// Response interceptor for consistent error handling
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    console.log('API Response:', response.status, response.config.url, response.data)
    return response
  },
  (error: AxiosError) => {
    console.error('API Error:', error.response?.status, error.config?.url, error.response?.data)
    
    // Transform axios errors to our API response format
    const apiError: ApiResponse = {
      error: {
        code: 'NETWORK_ERROR',
        message: 'Network request failed',
      }
    }

    if (error.response?.data) {
      // Backend returned an error response
      apiError.error = error.response.data as any
    } else if (error.code === 'ECONNABORTED') {
      apiError.error = {
        code: 'TIMEOUT',
        message: 'Request timeout',
      }
    } else if (error.message) {
      apiError.error!.message = error.message
    }

    return Promise.reject(apiError)
  }
)

export default apiClient