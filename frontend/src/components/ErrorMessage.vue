<template>
  <div 
    class="rounded-lg p-4 mb-4"
    :class="[
      variant === 'error' ? 'bg-red-50 border border-red-200' : 'bg-amber-50 border border-amber-200',
      containerClass
    ]"
  >
    <div class="flex items-start">
      <div class="flex-shrink-0">
        <svg 
          v-if="variant === 'error'"
          class="h-5 w-5 text-red-400" 
          viewBox="0 0 20 20" 
          fill="currentColor"
        >
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        <svg 
          v-else
          class="h-5 w-5 text-amber-400" 
          viewBox="0 0 20 20" 
          fill="currentColor"
        >
          <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
      </div>
      <div class="ml-3 flex-1">
        <h3 
          v-if="title"
          class="text-sm font-medium"
          :class="variant === 'error' ? 'text-red-800' : 'text-amber-800'"
        >
          {{ title }}
        </h3>
        <div 
          class="text-sm"
          :class="[
            variant === 'error' ? 'text-red-700' : 'text-amber-700',
            title ? 'mt-1' : ''
          ]"
        >
          <p>{{ message }}</p>
        </div>
        <div v-if="$slots.actions" class="mt-3">
          <slot name="actions" />
        </div>
      </div>
      <div v-if="dismissible" class="ml-auto pl-3">
        <button
          @click="$emit('dismiss')"
          class="inline-flex rounded-md p-1.5 focus:outline-none focus:ring-2 focus:ring-offset-2"
          :class="[
            variant === 'error' 
              ? 'text-red-500 hover:bg-red-100 focus:ring-red-600' 
              : 'text-amber-500 hover:bg-amber-100 focus:ring-amber-600'
          ]"
        >
          <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'error' | 'warning'
  title?: string
  message: string
  dismissible?: boolean
  containerClass?: string
}

withDefaults(defineProps<Props>(), {
  variant: 'error',
  dismissible: false,
  containerClass: '',
})

defineEmits<{
  dismiss: []
}>()
</script>