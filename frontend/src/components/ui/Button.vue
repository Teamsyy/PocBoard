<template>
  <button 
    :class="buttonClasses" 
    :disabled="disabled ? true : false"
    v-bind="$attrs"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'default',
  disabled: false
})

const buttonClasses = computed(() => {
  const baseClasses = 'inline-flex items-center justify-center rounded-lg font-medium transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50'
  
  const variantClasses = {
    default: 'bg-primary-500 text-white hover:bg-primary-600 focus:ring-primary-500',
    destructive: 'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500',
    outline: 'border-2 border-primary-300 text-primary-700 hover:bg-primary-50 focus:ring-primary-500',
    secondary: 'bg-amber-200 text-amber-800 hover:bg-amber-300 focus:ring-amber-500',
    ghost: 'text-primary-700 hover:bg-primary-50 focus:ring-primary-500',
    link: 'text-primary-600 underline-offset-4 hover:underline focus:ring-primary-500'
  }
  
  const sizeClasses = {
    default: 'h-10 px-4 py-2',
    sm: 'h-9 px-3 text-sm',
    lg: 'h-11 px-8',
    icon: 'h-10 w-10'
  }
  
  return [
    baseClasses,
    variantClasses[props.variant],
    sizeClasses[props.size]
  ].join(' ')
})
</script>

<style scoped>
/* Additional custom styles using CSS variables from tailwind config */
.bg-primary-500 {
  background-color: #ed7420;
}

.hover\:bg-primary-600:hover {
  background-color: #de5a16;
}

.focus\:ring-primary-500:focus {
  --tw-ring-color: #ed7420;
}

.text-primary-700 {
  color: #b84315;
}

.border-primary-300 {
  border-color: #f6ba77;
}

.hover\:bg-primary-50:hover {
  background-color: #fef7ee;
}

.text-primary-600 {
  color: #de5a16;
}
</style>
