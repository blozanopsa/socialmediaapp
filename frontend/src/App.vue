<script setup>
import { RouterLink, RouterView, useRoute } from 'vue-router'
import Loader from './components/Loader.vue'
import Navbar from './components/Navbar.vue'
import PublicChat from './components/PublicChat.vue'
import { useModalStore } from './stores/modal'

const route = useRoute()
const modalStore = useModalStore()
</script>

<template>
  <transition name="fade-navbar">
    <div
      :class="{ 'modal-blur': modalStore.isAnyModalOpen }"
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
      <div v-if="route.name !== 'auth-callback'" class="pt-24">
        <RouterView />
      </div>
      <div v-else>
        <RouterView />
      </div>
    </div>
  </transition>
</template>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Cal+Sans&display=swap');

.app-bg {
  font-family: 'Cal Sans', 'Inter', 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, 'sans-serif';
  color: #2d3748;
  background: linear-gradient(to top, #64748b 0%, #f8fafc 100%);
  min-height: 100vh;
  transition:
    filter 0.3s cubic-bezier(0.4, 0, 0.2, 1),
    opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  filter: none;
  opacity: 1;
}

body,
.app-bg {
  font-weight: 400;
  letter-spacing: 0.01em;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

.modal-blur {
  transition:
    filter 0.3s cubic-bezier(0.4, 0, 0.2, 1),
    opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  filter: blur(8px);
  opacity: 0.7;
}

.fade-navbar-enter-active,
.fade-navbar-leave-active {
  transition: opacity 0.3s ease;
}
.fade-navbar-enter-from,
.fade-navbar-leave-to {
  opacity: 0;
}
.fade-navbar-enter-to,
.fade-navbar-leave-from {
  opacity: 1;
}
</style>
