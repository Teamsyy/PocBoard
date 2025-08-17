import type { AppError } from '@/types'

/**
 * Extract user-friendly error message from API response
 */
export const getErrorMessage = (error: any): string => {
  if (error?.error?.message) {
    return error.error.message
  }
  
  if (error?.message) {
    return error.message
  }
  
  if (typeof error === 'string') {
    return error
  }
  
  return 'An unexpected error occurred'
}

/**
 * Check if error is a specific API error code
 */
export const isApiError = (error: any, code: string): boolean => {
  return error?.error?.code === code
}

/**
 * Check if error is an authorization error
 */
export const isAuthError = (error: any): boolean => {
  return isApiError(error, 'UNAUTHORIZED') || isApiError(error, 'FORBIDDEN')
}

/**
 * Check if error is a validation error
 */
export const isValidationError = (error: any): boolean => {
  return isApiError(error, 'VALIDATION_ERROR') || isApiError(error, 'BAD_REQUEST')
}

/**
 * Check if error is a not found error
 */
export const isNotFoundError = (error: any): boolean => {
  return isApiError(error, 'NOT_FOUND')
}

/**
 * Transform API error to AppError format
 */
export const transformApiError = (error: any): AppError => {
  if (error?.error) {
    return error.error
  }
  
  return {
    code: 'UNKNOWN_ERROR',
    message: getErrorMessage(error),
  }
}

/**
 * Get user-friendly error message based on error code
 */
export const getFriendlyErrorMessage = (error: any): string => {
  const code = error?.error?.code
  
  switch (code) {
    case 'UNAUTHORIZED':
      return 'You do not have permission to perform this action'
    case 'FORBIDDEN':
      return 'Access denied. Please check your access token'
    case 'NOT_FOUND':
      return 'The requested resource was not found'
    case 'VALIDATION_ERROR':
    case 'BAD_REQUEST':
      return error.error.message || 'Invalid input data'
    case 'TIMEOUT':
      return 'Request timed out. Please try again'
    case 'NETWORK_ERROR':
      return 'Network error. Please check your connection'
    default:
      return getErrorMessage(error)
  }
}