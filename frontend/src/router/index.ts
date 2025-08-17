import { createRouter, createWebHistory } from 'vue-router'
import Landing from '@/views/Landing.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: Landing,
    },
    {
      path: '/board/:boardId/edit',
      name: 'board-editor',
      component: () => import('@/views/BoardEditor.vue'),
      props: route => ({
        boardId: route.params.boardId,
        editToken: route.query.edit_token,
      }),
    },
    {
      path: '/board/:boardId/public',
      name: 'board-public',
      component: () => import('@/views/BoardPublic.vue'),
      props: route => ({
        boardId: route.params.boardId,
        publicToken: route.query.public_token,
      }),
    },
    {
      path: '/board/:boardId/recap',
      name: 'board-recap',
      component: () => import('@/views/Recap.vue'),
      props: route => ({
        boardId: route.params.boardId,
        editToken: route.query.edit_token,
      }),
    },
    // Catch-all route for 404 pages
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFound.vue'),
    },
  ],
})

// Navigation guards for token validation (will be implemented in later tasks)
router.beforeEach((_to, _from, next) => {
  // TODO: Add token validation logic in later tasks
  next()
})

export default router