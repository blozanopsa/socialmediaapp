<template>
  <div>
    <div v-if="retracted" class="fixed left-6 bottom-6 z-50">
      <button
        class="w-16 h-16 rounded-full bg-blue-600 text-white flex items-center justify-center shadow-lg text-3xl hover:bg-blue-800 transition-colors duration-200"
        @click="toggleRetract"
        aria-label="Open Public Chat"
      >
        💬
      </button>
    </div>
    <div
      v-else
      class="fixed left-6 bottom-6 z-50 w-80 h-[70vh] rounded-xl shadow-lg bg-white flex flex-col transition-all duration-200"
    >
      <div class="flex justify-between items-center bg-blue-600 text-white py-3 px-4 rounded-t-xl">
        <span class="text-2xl">Public Chat</span>
        <button
          class="w-8 h-8 rounded-full bg-blue-700 text-white flex items-center justify-center text-lg hover:bg-blue-900 transition-colors duration-200"
          @click="toggleRetract"
          aria-label="Close Public Chat"
        >
          ×
        </button>
      </div>
      <div class="flex-1 overflow-y-auto p-4 bg-gray-100" ref="messagesRef">
        <div
          v-for="(msg, idx) in messages"
          :key="idx"
          class="flex mb-2"
          :class="msg.user === userStore.user?.id ? 'justify-end' : 'justify-start'"
        >
          <div
            :class="[
              'max-w-[70%] px-4 py-2 break-words',
              msg.user === userStore.user?.id
                ? 'bg-blue-600 text-white self-end rounded-2xl rounded-br-md'
                : 'bg-gray-200 text-gray-900 self-start rounded-2xl rounded-bl-md',
            ]"
          >
            <span v-if="msg.user !== userStore.user?.id" class="text-blue-600 block mb-1">
              {{ userNames[msg.user] || `User ${msg.user}` }}
            </span>
            <span>{{ msg.text }}</span>
          </div>
        </div>
      </div>
      <form class="flex border-t border-gray-200 p-2 bg-white" @submit.prevent="sendMessage">
        <input
          v-model="input"
          type="text"
          placeholder="Type a message..."
          autocomplete="off"
          class="flex-1 border border-gray-200 rounded-md p-2 mr-2 outline-none"
        />
        <button
          type="submit"
          class="bg-blue-600 text-white rounded-md px-4 py-2 font-medium transition-colors duration-200 hover:bg-blue-800"
        >
          Send
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import axios from 'axios'

const retracted = ref(true)
const input = ref('')
const messages = ref([])
const messagesRef = ref(null)
const userNames = ref({})
const POLL_INTERVAL = 3000 // 3 seconds
let pollIntervalId = null

function toggleRetract() {
  retracted.value = !retracted.value
}

async function fetchMessages() {
  try {
    const res = await axios.get('http://localhost:8080/public-chat-messages')
    messages.value = res.data
    nextTick(() => {
      if (messagesRef.value) {
        messagesRef.value.scrollTop = messagesRef.value.scrollHeight
      }
    })
  } catch (e) {
    // ignore fetch errors
  }
}

async function fetchUserNames() {
  try {
    const res = await axios.get('http://localhost:8080/users/names')
    // Map userId to name
    userNames.value = Object.fromEntries(res.data.map((u) => [u.id, u.name]))
  } catch (e) {
    // ignore fetch errors
  }
}

import { useUserStore } from '@/stores/user'
const userStore = useUserStore()

async function sendMessage() {
  if (input.value.trim() === '') return
  const userId = userStore.user?.id || null
  await axios.post('http://localhost:8080/public-chat-messages', {
    user: userId,
    text: input.value,
  })
  input.value = ''
  await fetchMessages()
}

onMounted(() => {
  fetchMessages()
  fetchUserNames()
  pollIntervalId = setInterval(fetchMessages, POLL_INTERVAL)
})

onUnmounted(() => {
  if (pollIntervalId) clearInterval(pollIntervalId)
})
</script>

<style scoped>
/* No custom CSS needed, all handled by Tailwind */
</style>
