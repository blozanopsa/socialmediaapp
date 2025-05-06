<template>
  <nav
    class="w-full flex items-center justify-between px-8 py-4 bg-white shadow-md fixed top-0 left-0 z-50"
  >
    <div class="text-4xl text-blue-600 select-none">Social Media App</div>
    <div class="flex mr-auto ml-10 justify-center">
      <nav class="flex gap-2">
        <RouterLink
          to="/postfeed"
          class="group flex items-center px-3 py-2 text-blue-600 rounded-xl transition-all duration-150 gap-2 focus:outline-none"
          :class="{
            'bg-blue-600 text-white shadow hover:bg-blue-500/60 font-semibold':
              $route.path === '/postfeed',
            ' hover:bg-blue-300/60': $route.path !== '/postfeed',
            '': '',
          }"
        >
          <span>Home</span>
        </RouterLink>
      </nav>
    </div>
    <div class="flex flex-row items-center gap-4">
      <div>{{ user?.Name || user?.name || 'Guest' }}</div>
      <button
        @click="logout"
        class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-red-500 to-red-700 text-white rounded-full shadow-md hover:from-red-600 hover:to-red-800 hover:scale-105 transition-all duration-200 font-semibold focus:outline-none focus:ring-2 focus:ring-red-400"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a2 2 0 01-2 2H7a2 2 0 01-2-2V7a2 2 0 012-2h4a2 2 0 012 2v1"
          />
        </svg>
        <span>Logout</span>
      </button>
    </div>
  </nav>
</template>

<script setup>
import axios from 'axios'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

const router = useRouter()
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

async function logout() {
  // Clear user state
  userStore.clearUser()
  localStorage.removeItem('userName')
  // Call the new backend logout endpoint to clear the sessionID cookie
  await axios.post(
    'http://localhost:8080/auth/session/logout',
    {},
    {
      withCredentials: true,
    },
  )
  console.log('NavBar: Called backend to clear sessionID cookie and session')
  // Redirect to homepage
  router.replace({ name: 'home' })
}
</script>

<style scoped></style>
