// Board types
export interface Board {
  id: string
  title: string
  edit_token: string
  public_token: string
  created_at: string
  updated_at: string
  pages?: Page[]
}

// Page types
export interface Page {
  id: string
  boardId: string
  title: string
  date: string
  orderIdx: number
  createdAt: string
  updatedAt: string
  elements?: Element[]
}

// Element types
export interface Element {
  id: string
  pageId: string
  kind: ElementKind
  x: number
  y: number
  w: number
  h: number
  rotation: number
  z: number
  payload: ElementPayload
  visible: boolean
  locked: boolean
  createdAt: string
  updatedAt: string
}

export type ElementKind = 'text' | 'image' | 'sticker' | 'shape'

// Element payload types
export interface TextPayload {
  content: string
  fontFamily: string
  fontSize: number
  color: string
  bold: boolean
  italic: boolean
  textAlign: 'left' | 'center' | 'right'
}

export interface ImagePayload {
  url: string
  originalWidth: number
  originalHeight: number
  description?: string
  filters?: any[]
}

export interface ShapePayload {
  shapeType: 'rectangle' | 'circle' | 'triangle'
  fill: string
  stroke: string
  strokeWidth: number
}

export interface StickerPayload {
  stickerType: string
  url: string
  category: string
}

export type ElementPayload = TextPayload | ImagePayload | ShapePayload | StickerPayload

// API Response types
export interface ApiResponse<T = any> {
  data?: T
  error?: {
    code: string
    message: string
    details?: any
  }
}

// Recap types
export interface RecapData {
  pageCount: number
  elementCount: number
  pages: RecapPage[]
}

export interface RecapPage {
  id: string
  title: string
  date: string
  elementCount: number
  thumbnail?: string
}

// Filter types
export type DateFilter = 'day' | 'week' | 'month'

// Canvas state types
export interface CanvasState {
  elements: Element[]
  timestamp: number
}

// Fabric.js types are now provided by @types/fabric

// Upload types
export interface UploadResponse {
  url: string
  filename: string
}

// Error types
export interface AppError {
  code: string
  message: string
  details?: any
}