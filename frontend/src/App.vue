<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import Loader from './components/Loader.vue'
import Navbar from './components/Navbar.vue'
import PublicChat from './components/PublicChat.vue'

const route = useRoute()
</script>

<template>
  <div
    class="flex flex-col min-h-screen w-screen bg-gradient-to-t from-slate-500 to-slate-50 app-bg"
  >
    <Navbar v-if="route.name && route.name !== 'home' && route.name !== 'auth-callback'" />
    <PublicChat
      v-if="
        route.name &&
        route.name !== 'auth-callback' &&
        route.name !== 'home' &&
        (route.name === 'home' ||
          route.name === 'likes' ||
          route.name === 'postfeed' ||
          route.name === 'profile')
      "
    />
    <div class="pt-24">
      <RouterView />
    </div>
  </div>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cal+Sans&display=swap');

.app-bg {
  font-family: 'Cal Sans', 'Inter', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'sans-serif';
  color: #2d3748;
  background: linear-gradient(to top, #64748b 0%, #f8fafc 100%);
  min-height: 100vh;
}

body,
.app-bg {
  font-weight: 400;
  letter-spacing: 0.01em;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
