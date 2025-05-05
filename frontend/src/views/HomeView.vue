<script setup>
import { ref } from 'vue'
import Loader from '@/components/Loader.vue'

const loading = ref(false)

function handleMicrosoftLogin() {
  loading.value = true
  window.location.href = 'http://localhost:8080/auth/login'
}

function handleFacebookLogin() {
  loading.value = true
  window.location.href = 'http://localhost:8080/auth/facebook/login'
}
</script>

<template>
  <div class="flex flex-col items-center justify-center min-h-[80vh]">
    <div class="text-[3rem] opacity-0 animate-fadeinup mb-8 text-center">
      Welcome to Social Media App
    </div>
    <div class="flex flex-row gap-5 items-center justify-center text-[1.5rem]">
      <button @click="handleMicrosoftLogin" class="login-btn microsoft-btn animate-fadeinup">
        <img
          src="https://upload.wikimedia.org/wikipedia/commons/4/44/Microsoft_logo.svg"
          alt="Microsoft Logo"
          class="w-6 h-6"
        />
        Sign in with Microsoft
      </button>
      <button @click="handleFacebookLogin" class="login-btn facebook-btn animate-fadeinup">
        <img src="@/assets/facebook_logo.png" alt="Facebook Logo" class="w-6 h-6" />
        Sign in with Facebook
      </button>
    </div>
    <div
      v-if="loading"
      class="fixed inset-0 z-50 flex items-center justify-center"
      style="background: rgba(100, 116, 139, 0.5)"
    >
      <Loader />
    </div>
  </div>
</template>
<style scoped>
.login-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border-width: 1px;
  border-radius: 0.375rem;
  padding: 0.5rem 1rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition:
    background 0.2s,
    color 0.2s,
    border 0.2s;
}
.microsoft-btn {
  background: #fff;
  color: #5e5e5e;
  border-color: #8c8c8c;
}
.microsoft-btn:hover {
  background: #e5e7eb;
  color: #2563eb;
  border-color: #2563eb;
}
.facebook-btn {
  background: #1877f2;
  color: #fff;
  border-color: #1877f2;
}
.facebook-btn:hover {
  background: #145db2;
  color: #fff;
  border-color: #145db2;
}
@keyframes fadeinup {
  0% {
    opacity: 0;
    transform: translateY(2rem);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}
.animate-fadeinup {
  animation: fadeinup 1.5s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}
</style>
