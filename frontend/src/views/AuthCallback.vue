<script setup>
import { onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import Loader from '@/components/Loader.vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

onMounted(async () => {
  try {
    // Add a 0.5 second delay so the loader is visible
    await new Promise((resolve) => setTimeout(resolve, 500))

    // Fetch user data from the backend
    const response = await fetch('http://localhost:8080/api/user', {
      credentials: 'include',
    })

    if (!response.ok) throw new Error('Failed to fetch user data')

    const userData = await response.json()

    // Store user data in localStorage (now includes id)
    localStorage.setItem('user', JSON.stringify(userData))
    if (userData && userData.name) {
      localStorage.setItem('userName', userData.name)
    } else {
      router.replace({ name: 'home' })
      return
    }

    // Save to Pinia store and localStorage
    userStore.setToken(userData.token || '')
    userStore.setUser(userData || {})

    // Redirect to postfeed
    router.replace({ name: 'postfeed' })
  } catch (e) {
    router.replace({ name: 'home' })
  }

  // Immediately read cookies and store session ID
  const cookies = document.cookie.split('; ').reduce((acc, cookie) => {
    const [key, value] = cookie.split('=')
    acc[key] = value
    return acc
  }, {})
  if (cookies.sessionID) {
    localStorage.setItem('sessionID', cookies.sessionID)
  }
})
</script>

<template>
  <Loader />
</template>
