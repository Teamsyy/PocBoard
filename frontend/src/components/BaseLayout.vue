<template>
  <div class="min-h-screen bg-amber-50">
    <header v-if="showHeader" class="bg-white shadow-sm border-b border-amber-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center">
            <router-link to="/" class="text-xl font-bold text-amber-800 font-journal">
              Junk Journal & Board
            </router-link>
          </div>
          
          <nav v-if="boardId" class="flex items-center space-x-4">
            <router-link 
              :to="`/board/${boardId}/edit${editToken ? `?token=${editToken}` : ''}`"
              class="text-amber-600 hover:text-amber-800 font-medium"
              :class="{ 'text-amber-800 font-semibold': $route.name === 'board-editor' }"
            >
              Editor
            </router-link>
            <router-link 
              :to="`/board/${boardId}/public${publicToken ? `?token=${publicToken}` : ''}`"
              class="text-amber-600 hover:text-amber-800 font-medium"
              :class="{ 'text-amber-800 font-semibold': $route.name === 'board-public' }"
            >
              Public View
            </router-link>
            <router-link 
              :to="`/board/${boardId}/recap${editToken ? `?token=${editToken}` : ''}`"
              class="text-amber-600 hover:text-amber-800 font-medium"
              :class="{ 'text-amber-800 font-semibold': $route.name === 'board-recap' }"
            >
              Recap
            </router-link>
          </nav>
        </div>
      </div>
    </header>

    <main class="flex-1">
      <slot />
    </main>

    <footer v-if="showFooter" class="bg-white border-t border-amber-200 py-4">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center text-sm text-amber-600">
          <p>&copy; 2024 Junk Journal & Board. Create beautiful visual journals.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

interface Props {
  showHeader?: boolean
  showFooter?: boolean
  boardId?: string
  editToken?: string
  publicToken?: string
}

const props = withDefaults(defineProps<Props>(), {
  showHeader: true,
  showFooter: false,
})

const route = useRoute()

// Extract tokens from route query if not provided as props
const editToken = computed(() => props.editToken || route.query.token as string)
const publicToken = computed(() => props.publicToken || route.query.public_token as string)
const boardId = computed(() => props.boardId || route.params.boardId as string)
</script>