<template>
  <nav
    class="w-full flex items-center justify-between px-8 py-4 bg-white shadow-md fixed top-0 left-0 z-50"
  >
    <div class="text-4xl text-blue-600 select-none">Social Media App</div>
    <div class="flex-1 flex justify-center">
      <nav class="flex gap-4">
        <RouterLink to="/postfeed" class="nav-btn" active-class="nav-btn-active">Home</RouterLink>
        <RouterLink to="/likes" class="nav-btn" active-class="nav-btn-active">Likes</RouterLink>
        <RouterLink to="/profile" class="nav-btn" active-class="nav-btn-active">Profile</RouterLink>
      </nav>
    </div>
    <div class="flex flex-row items-center gap-4">
      <div>{{ user?.Name || user?.name || 'Guest' }}</div>
      <button
        @click="logout"
        class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 transition duration-200"
      >
        Logout
      </button>
    </div>
  </nav>
</template>

<script setup>
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
  await fetch('http://localhost:8080/auth/session/logout', {
    method: 'POST',
    credentials: 'include',
  })
  console.log('NavBar: Called backend to clear sessionID cookie and session')
  // Redirect to homepage
  router.replace({ name: 'home' })
}
</script>

<style scoped>
.nav-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-width: 1px;
  border-radius: 0.375rem;
  padding: 0.5rem 1.25rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  background: #fff;
  color: #2563eb;
  border-color: #2563eb;
  transition:
    background 0.2s,
    color 0.2s,
    border 0.2s;
  text-decoration: none;
}
.nav-btn:hover {
  background: #2563eb;
  color: #fff;
  border-color: #2563eb;
}
.nav-btn-active {
  background: #2563eb;
  color: #fff;
  border-color: #2563eb;
}
</style>
