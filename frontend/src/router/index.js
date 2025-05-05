import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { requireAuth } from '@/utils/auth.guard.js'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import('../views/HomeView.vue'),
    },
    {
      path: '/postfeed',
      name: 'postfeed',
      component: () => import('../views/PostFeed.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/likes',
      name: 'likes',
      component: () => import('../views/LikesView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/auth/microsoft/callback',
      name: 'auth-callback',
      component: () => import('../views/AuthCallback.vue'),
    },
  ],
})

router.beforeEach(requireAuth)

export default router
