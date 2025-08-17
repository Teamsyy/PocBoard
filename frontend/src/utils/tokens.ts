/**
 * Token management utilities for URL-based access control
 */

export interface TokenInfo {
  editToken: string | null
  publicToken: string | null
  isEditMode: boolean
}

/**
 * Extract tokens from current URL parameters
 */
export const getTokensFromUrl = (): TokenInfo => {
  const urlParams = new URLSearchParams(window.location.search)
  const editToken = urlParams.get('edit_token')
  const publicToken = urlParams.get('public_token')
  
  return {
    editToken,
    publicToken,
    isEditMode: !!editToken,
  }
}

/**
 * Generate edit URL with edit token
 */
export const generateEditUrl = (editToken: string, path = '/editor'): string => {
  const url = new URL(path, window.location.origin)
  url.searchParams.set('edit_token', editToken)
  return url.toString()
}

/**
 * Generate public URL with public token
 */
export const generatePublicUrl = (publicToken: string, path = '/public'): string => {
  const url = new URL(path, window.location.origin)
  url.searchParams.set('public_token', publicToken)
  return url.toString()
}

/**
 * Update current URL with new token parameters
 */
export const updateUrlWithTokens = (editToken?: string, publicToken?: string): void => {
  const url = new URL(window.location.href)
  
  if (editToken) {
    url.searchParams.set('edit_token', editToken)
    url.searchParams.delete('public_token')
  } else if (publicToken) {
    url.searchParams.set('public_token', publicToken)
    url.searchParams.delete('edit_token')
  }
  
  window.history.replaceState({}, '', url.toString())
}

/**
 * Check if current URL has valid access tokens
 */
export const hasValidTokens = (): boolean => {
  const { editToken, publicToken } = getTokensFromUrl()
  return !!(editToken || publicToken)
}

/**
 * Remove tokens from URL
 */
export const clearTokensFromUrl = (): void => {
  const url = new URL(window.location.href)
  url.searchParams.delete('edit_token')
  url.searchParams.delete('public_token')
  window.history.replaceState({}, '', url.toString())
}

/**
 * Copy URL to clipboard
 */
export const copyUrlToClipboard = async (url: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(url)
  } catch (err) {
    // Fallback for older browsers
    const textArea = document.createElement('textarea')
    textArea.value = url
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
  }
}