declare module 'vue3-draggable-resizable' {
  import { DefineComponent } from 'vue'
  
  export interface Vue3DraggableResizableProps {
    initW?: number
    initH?: number
    x?: number
    y?: number
    w?: number
    h?: number
    minW?: number
    minH?: number
    maxW?: number
    maxH?: number
    parent?: boolean
    draggable?: boolean
    resizable?: boolean
    active?: boolean
    preventActiveBehavior?: boolean
    lockAspectRatio?: boolean
    disabledX?: boolean
    disabledY?: boolean
    disabledW?: boolean
    disabledH?: boolean
    boundary?: boolean
    'onUpdate:x'?: (x: number) => void
    'onUpdate:y'?: (y: number) => void
    'onUpdate:w'?: (w: number) => void
    'onUpdate:h'?: (h: number) => void
    'onUpdate:active'?: (active: boolean) => void
    onActivated?: () => void
    onDeactivated?: () => void
    onDragStart?: (position: { x: number; y: number }) => void
    onDragEnd?: (position: { x: number; y: number }) => void
    onResizeStart?: (rect: { x: number; y: number; width: number; height: number }) => void
    onResizeEnd?: (rect: { x: number; y: number; width: number; height: number }) => void
    onDragging?: (position: { x: number; y: number }) => void
    onResizing?: (rect: { x: number; y: number; width: number; height: number }) => void
  }

  const Vue3DraggableResizable: DefineComponent<Vue3DraggableResizableProps>
  export default Vue3DraggableResizable
}
