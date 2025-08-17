<template>
  <BaseLayout :board-id="boardId">
    <div class="h-screen flex flex-col">
      <!-- Placeholder for recap view -->
      <div class="flex-1 flex items-center justify-center bg-amber-50">
        <div class="text-center">
          <h1 class="text-2xl font-bold text-amber-800 mb-4">
            Board Recap
          </h1>
          <p class="text-amber-600">
            Board ID: {{ boardId }}
          </p>
          <p class="text-amber-600 text-sm mt-2">
            Recap implementation coming in later tasks
          </p>
        </div>
      </div>
    </div>
  </BaseLayout>
</template>

<script setup lang="ts">
import { BaseLayout } from '@/components'

interface Props {
  boardId: string
}

defineProps<Props>()
</script>