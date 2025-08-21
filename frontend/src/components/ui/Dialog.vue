<template>
  <Teleport to="body">
    <div v-if="open" class="dialog-overlay" @click="handleOverlayClick">
      <div class="dialog-content" @click.stop>
        <slot />
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
interface Props {
  open?: boolean
  onOpenChange?: (open: boolean) => void
}

const props = withDefaults(defineProps<Props>(), {
  open: false
})

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const handleOverlayClick = () => {
  emit('update:open', false)
  props.onOpenChange?.(false)
}
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  z-index: 50;
  background: rgba(0, 0, 0, 0.8);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  animation: dialog-overlay-show 150ms cubic-bezier(0.16, 1, 0.3, 1);
}

.dialog-content {
  background: white;
  border-radius: 16px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  max-width: 90vw;
  max-height: 85vh;
  overflow: hidden;
  animation: dialog-content-show 150ms cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes dialog-overlay-show {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes dialog-content-show {
  from {
    opacity: 0;
    transform: scale(0.96) translateY(-4px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
</style>
