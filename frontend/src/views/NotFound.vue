<template>
  <BaseLayout :show-header="false">
    <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-amber-50 to-orange-100">
      <div class="text-center">
        <div class="mb-8">
          <h1 class="text-9xl font-bold text-amber-300 font-journal">404</h1>
        </div>
        
        <div class="mb-8">
          <h2 class="text-2xl font-bold text-amber-800 mb-2">
            Page Not Found
          </h2>
          <p class="text-amber-600 max-w-md mx-auto">
            The page you're looking for doesn't exist. It might have been moved, deleted, or you entered the wrong URL.
          </p>
        </div>

        <div class="space-y-4">
          <router-link to="/" class="btn-primary inline-block">
            Go Home
          </router-link>
          <div>
            <button @click="$router.back()" class="btn-secondary">
              Go Back
            </button>
          </div>
        </div>
      </div>
    </div>
  </BaseLayout>
</template>

<script setup lang="ts">
import { BaseLayout } from '@/components'
</script>