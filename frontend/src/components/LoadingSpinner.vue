<template>
  <div class="flex items-center justify-center" :class="containerClass">
    <div class="relative">
      <div 
        class="animate-spin rounded-full border-4 border-amber-200"
        :class="[
          sizeClass,
          `border-t-${color}-500`
        ]"
      ></div>
      <div v-if="text" class="mt-2 text-center">
        <p :class="`text-${color}-600 text-sm font-medium`">{{ text }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: 'sm' | 'md' | 'lg'
  color?: 'amber' | 'primary'
  text?: string
  fullScreen?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  color: 'amber',
  fullScreen: false,
})

const sizeClass = computed(() => {
  switch (props.size) {
    case 'sm': return 'w-6 h-6'
    case 'lg': return 'w-12 h-12'
    default: return 'w-8 h-8'
  }
})

const containerClass = computed(() => {
  return props.fullScreen ? 'min-h-screen' : 'p-4'
})
</script>