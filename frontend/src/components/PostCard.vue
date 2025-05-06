<template>
  <div
    class="bg-white/90 rounded-3xl shadow-xl p-8 mb-8 flex flex-col gap-4 border border-slate-200 transition-transform hover:scale-[1.01] hover:shadow-2xl relative cursor-pointer"
    @click="onCardClick"
  >
    <button
      v-if="canDelete"
      @click.stop="onDelete"
      class="absolute top-1 right-2 text-red-400 hover:text-red-700 text-xl px-2 py-1 transition"
      aria-label="Delete post"
    >
      ×
    </button>
    <div class="flex items-center gap-2">
      <span class="text-blue-700 text-lg">{{ post.User?.Name || 'Unknown' }}</span>
      <span class="text-xs text-gray-400 ml-auto">{{ formatDate(post.CreatedAt) }}</span>
    </div>
    <div class="text-xl text-gray-800 font-medium">{{ post.Description }}</div>
    <div v-if="post.ImageURL" class="my-2">
      <img
        :src="
          post.ImageURL.startsWith('/') ? post.ImageURL : 'http://localhost:8080/' + post.ImageURL
        "
        alt="Post image"
        class="max-h-80 rounded-lg border border-gray-200 object-contain mx-auto"
        style="max-width: 100%"
      />
    </div>
    <div class="flex items-center gap-6 mt-2">
      <div class="flex items-center gap-2">
        <button
          @click.stop="onLike"
          :aria-label="likedByUser ? 'Unlike' : 'Like'"
          class="relative group"
        >
          <svg
            v-if="likedByUser && post.User?.ID !== userId"
            xmlns="http://www.w3.org/2000/svg"
            fill="currentColor"
            viewBox="0 0 24 24"
            class="w-7 h-7 text-red-500 transition-transform duration-200 scale-110"
          >
            <path
              d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41 0.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
            />
          </svg>
          <svg
            v-else
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            class="w-7 h-7 text-gray-400 group-hover:text-red-400 transition-colors"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41 0.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"
            />
          </svg>
          <span
            v-if="showOwnLikeTooltip"
            class="absolute -top-8 left-1/2 -translate-x-1/2 bg-black text-white text-xs rounded px-2 py-1 whitespace-nowrap z-10"
          >
            You can't like your own post.
          </span>
        </button>
        <span class="text-base text-gray-700 min-w-[2ch] text-center">{{
          post.likesCount || (post.Likes ? post.Likes.length : 0)
        }}</span>
      </div>
      <div class="flex items-center gap-2">
        <button @click.stop="onComment" aria-label="Comment" class="relative group">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            class="w-7 h-7 text-gray-400 group-hover:text-blue-500 transition-colors"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M8 10h.01M12 10h.01M16 10h.01M21 12c0 4.418-4.03 8-9 8a9.77 9.77 0 01-4-.8l-4 1 1-4A8.96 8.96 0 013 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
            />
          </svg>
        </button>
        <span class="text-base text-gray-700 min-w-[2ch] text-center">{{
          post.commentsCount || (post.Comments ? post.Comments.length : 0)
        }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
const props = defineProps({
  post: Object,
  likedByUser: Boolean,
})
const emit = defineEmits(['like', 'delete', 'comment'])
const showOwnLikeTooltip = ref(false)
const userId = computed(() => {
  const userStore = JSON.parse(localStorage.getItem('user'))
  return userStore?.id
})
const canDelete = computed(() => {
  const userStore = JSON.parse(localStorage.getItem('user'))
  return userStore && userStore.id === props.post.User?.ID
})
function formatDate(date) {
  if (!date) return ''
  return new Date(date).toLocaleString()
}
function onLike(e) {
  if (props.post.User?.ID === undefined || props.post.User?.ID === null) return
  // If user is the author, show tooltip
  const userStore = JSON.parse(localStorage.getItem('user'))
  if (userStore && userStore.id === props.post.User.ID) {
    showOwnLikeTooltip.value = true
    setTimeout(() => (showOwnLikeTooltip.value = false), 1500)
    return
  }
  emit('like', props.post)
}
function onDelete() {
  emit('delete', props.post)
}
function onComment() {
  emit('comment', props.post)
}
function onCardClick() {
  emit('comment', props.post)
}
</script>
