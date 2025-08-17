<template>
  <BaseLayout :show-header="false" :show-footer="true">
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-amber-50 to-orange-100">
      <div class="max-w-md w-full mx-4">
        <div class="text-center mb-8">
          <h1 class="text-4xl font-bold text-amber-800 mb-2 font-journal">
            Junk Journal & Board
          </h1>
          <p class="text-amber-600">
            Create beautiful visual journal entries with drag-and-drop editing
          </p>
        </div>

        <div class="bg-white rounded-xl shadow-lg p-6">
          <h2 class="text-xl font-semibold text-amber-800 mb-4">
            Create New Board
          </h2>

          <ErrorMessage 
            v-if="error"
            :message="error"
            @dismiss="error = ''"
            dismissible
          />

          <form @submit.prevent="createBoard" class="space-y-4">
            <div>
              <label for="title" class="block text-sm font-medium text-amber-700 mb-1">
                Board Title
              </label>
              <input 
                id="title" 
                v-model="boardTitle" 
                type="text" 
                required 
                placeholder="My Journal Board"
                class="input-field w-full"
                :disabled="loading"
              />
            </div>

            <button type="submit" :disabled="loading || !boardTitle.trim()" class="btn-primary w-full">
              <LoadingSpinner v-if="loading" size="sm" color="primary" class="inline-block mr-2" />
              {{ loading ? 'Creating...' : 'Create Board' }}
            </button>
          </form>
        </div>
      </div>
    </div>
  </BaseLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { BaseLayout, LoadingSpinner, ErrorMessage } from '@/components'
import { useBoardsStore } from '@/stores/boards'
import { getFriendlyErrorMessage } from '@/utils/api-errors'

const router = useRouter()
const boardsStore = useBoardsStore()
const boardTitle = ref('')
const loading = ref(false)
const error = ref('')

const createBoard = async () => {
  if (!boardTitle.value.trim()) {
    error.value = 'Please enter a board title'
    return
  }

  loading.value = true
  error.value = ''
  
  try {
    // Create board via store
    const response = await boardsStore.createBoard(
      boardTitle.value.trim(),
      'default'
    )

    console.log('Board creation response:', response)
    console.log('Full response structure:', JSON.stringify(response, null, 2))
    
    // Extract edit token from edit URL or board object
    let editToken = null
    
    // First try to get token from the board object
    if (response.board.edit_token) {
      editToken = response.board.edit_token
      console.log('Got token from board object:', editToken)
    } else if (response.edit_url) {
      // Try to extract from URL as fallback
      console.log('Edit URL string:', response.edit_url)
      try {
        const url = new URL(response.edit_url)
        editToken = url.searchParams.get('edit_token')
        console.log('Extracted token from URL:', editToken)
      } catch (urlError) {
        console.error('Failed to parse URL:', response.edit_url, urlError)
        // Try regex fallback for malformed URLs
        const tokenMatch = response.edit_url.match(/edit_token=([^&]+)/)
        if (tokenMatch) {
          editToken = tokenMatch[1]
          console.log('Extracted token via regex:', editToken)
        }
      }
    }

    if (!editToken) {
      throw new Error('No edit token found in response')
    }

    // Navigate to the board editor with the edit token
    await router.push({
      name: 'board-editor',
      params: { boardId: response.board.id },
      query: { edit_token: editToken }
    })
  } catch (err) {
    console.error('Failed to create board:', err)
    error.value = getFriendlyErrorMessage(err)
  } finally {
    loading.value = false
  }
}
</script>