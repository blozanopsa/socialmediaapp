<template>
  <div>
    <div :class="{ 'modal-blur': showModal || showCommentModal }">
      <div class="flex flex-col items-center ml-auto mr-auto w-[50%] mt-4">
        <div class="w-full flex justify-start mb-8">
          <button
            @click="showModal = true"
            class="flex items-center justify-center ml-auto mr-auto gap-2 bg-gradient-to-r from-blue-500 to-blue-700 text-white px-5 py-2 rounded-full shadow-lg hover:from-blue-600 hover:to-blue-800 hover:scale-105 transition-all duration-200 text-base focus:outline-none focus:ring-2 focus:ring-blue-400"
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
                d="M12 4v16m8-8H4"
              />
            </svg>
            <span>Create Post</span>
          </button>
        </div>
        <div v-if="posts.length > 0" class="w-full max-w-xl mb-8">
          <div class="flex justify-end gap-2 mb-4">
            <button
              @click="sortOrder = 'newest'"
              :class="{
                'bg-blue-600 text-white': sortOrder === 'newest',
                'bg-gray-200 hover:bg-gray-300': sortOrder !== 'newest',
              }"
              class="px-3 py-1 rounded-lg text-sm transition-colors"
            >
              Sort by Newest
            </button>
            <button
              @click="sortOrder = 'oldest'"
              :class="{
                'bg-blue-600 text-white': sortOrder === 'oldest',
                'bg-gray-200 hover:bg-gray-300': sortOrder !== 'oldest',
              }"
              class="px-3 py-1 rounded-lg text-sm transition-colors"
            >
              Sort by Oldest
            </button>
          </div>
          <div v-if="loading" class="flex justify-center py-8">
            <Loader />
          </div>
          <PostList
            v-else
            :posts="sortedPosts"
            :liked-by-user-map="likedByUserMap"
            @like="handleLike"
            @delete="handleDelete"
            @comment="openCommentModal"
          />
        </div>
      </div>
    </div>
    <BaseModal :show="showModal" @close="showModal = false">
      <form @submit.prevent="handleNewPost" class="flex flex-col gap-4">
        <h2 class="text-xl mb-2">Create a Post</h2>
        <input
          v-model="newPostContent"
          type="text"
          placeholder="What's on your mind?"
          class="border border-gray-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-400"
        />
        <label
          class="flex items-center gap-2 cursor-pointer w-fit px-4 py-2 bg-blue-50 border border-blue-300 rounded-lg hover:bg-blue-100 transition text-blue-700 font-medium"
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
              d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2M7 10l5 5 5-5M12 15V3"
            />
          </svg>
          <span>Upload Image</span>
          <input
            type="file"
            accept="image/*"
            @change="onImageChange"
            class="hidden"
            aria-label="Upload image"
          />
        </label>
        <!-- Image Preview -->
        <div v-if="imagePreview" class="mt-2">
          <img :src="imagePreview" alt="Image preview" class="max-h-32 rounded border" />
          <button
            type="button"
            @click="removeImage"
            class="block mt-1 text-xs text-red-600 hover:underline"
          >
            Remove
          </button>
        </div>
        <div class="flex gap-2 justify-end">
          <button
            type="button"
            @click="showModal = false"
            class="px-4 py-2 rounded-lg bg-gray-200 hover:bg-gray-300"
          >
            Cancel
          </button>
          <button
            type="submit"
            class="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition"
            :disabled="imageUploading"
          >
            Post
          </button>
        </div>
      </form>
    </BaseModal>
    <PostModal
      :show="showCommentModal"
      :selected-post="selectedPost"
      :comments="comments"
      :loading-comments="loadingComments"
      :editing-post="editingPost"
      v-model:editPostContent="editPostContent"
      :editing-comment-id="editingCommentId"
      :edit-comment-content="editCommentContent"
      :open-comment-menu-id="openCommentMenuId"
      :menu-position="menuPosition"
      :liked-by-user-map="likedByUserMap"
      @close="showCommentModal = false"
      @like="handleLike"
      @add-comment="addComment"
      @cancel-edit-post="cancelEditPost"
      @save-edit-post="saveEditPost"
      @start-edit-post="startEditPost"
      @edit-comment="editComment"
      @cancel-edit-comment="cancelEditComment"
      @save-edit-comment="saveEditComment"
      @start-edit-comment="startEditComment"
      @delete-comment="deleteComment"
      @toggle-comment-menu="toggleCommentMenu"
    />
  </div>
</template>
<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick, watch, computed } from 'vue' // Added computed
import Loader from '@/components/Loader.vue'
import PostList from '@/components/PostList.vue'
import BaseModal from '@/components/modals/BaseModal.vue'
import PostModal from '@/components/modals/PostModal.vue'
import axios from 'axios'
import { useUserStore } from '@/stores/user'
import { useModalStore } from '@/stores/modal'

const userStore = useUserStore()
const modalStore = useModalStore()
const posts = ref([])
const likedByUserMap = ref({})
const loading = ref(false)
const newPostContent = ref('')
const showModal = ref(false)
const showCommentModal = ref(false)
const commentInput = ref('')
const selectedPost = ref(null)
const comments = ref([])
const loadingComments = ref(false)
const openCommentMenuId = ref(null)
const menuPosition = ref({ top: 0, left: 0 })
const editingPost = ref(false)
const editPostContent = ref('')
const editingCommentId = ref(null)
const editCommentContent = ref('')
const imageFile = ref(null)
const imagePreview = ref(null)
const sortOrder = ref('newest') // Default sort order

// Watch for modal open/close and update global modal state
watch([showModal, showCommentModal], ([modal, commentModal]) => {
  modalStore.isAnyModalOpen = modal || commentModal
})

const sortedPosts = computed(() => {
  return [...posts.value].sort((a, b) => {
    const dateA = new Date(a.CreatedAt).getTime()
    const dateB = new Date(b.CreatedAt).getTime()
    return sortOrder.value === 'newest' ? dateB - dateA : dateA - dateB
  })
})

function toggleCommentMenu(id, event) {
  if (openCommentMenuId.value === id) {
    openCommentMenuId.value = null
    return
  }
  openCommentMenuId.value = id
  nextTick(() => {
    const rect = event.target.getBoundingClientRect()
    menuPosition.value = {
      top: window.scrollY + rect.top - 8, // 8px above the button
      left: window.scrollX + rect.right - 96, // align right edge, 96px is menu width
    }
  })
}

function handleClickOutsideMenu(e) {
  // Only close if a menu is open and the click is outside any open menu
  if (openCommentMenuId.value !== null) {
    const menus = document.querySelectorAll('.comment-menu-tooltip')
    let inside = false
    menus.forEach((menu) => {
      if (menu.contains(e.target)) inside = true
    })
    if (!inside) openCommentMenuId.value = null
  }
}

onMounted(() => {
  fetchPosts()
  fetchLikes()
  document.addEventListener('mousedown', handleClickOutsideMenu)
  window.addEventListener('scroll', () => {
    openCommentMenuId.value = null
  })
  window.addEventListener('resize', () => {
    openCommentMenuId.value = null
  })
})

onBeforeUnmount(() => {
  document.removeEventListener('mousedown', handleClickOutsideMenu)
  window.removeEventListener('scroll', () => {
    openCommentMenuId.value = null
  })
  window.removeEventListener('resize', () => {
    openCommentMenuId.value = null
  })
})

async function fetchPosts() {
  loading.value = true
  try {
    const res = await axios.get('http://localhost:8080/api/posts', { withCredentials: true })
    posts.value = res.data
  } finally {
    loading.value = false
  }
}

async function fetchLikes() {
  if (!userStore.user?.id) return
  const res = await axios.get(`http://localhost:8080/api/posts?likedBy=${userStore.user.id}`, {
    withCredentials: true,
  })
  likedByUserMap.value = Object.fromEntries(res.data.map((post) => [String(post.ID), true]))
}

async function handleLike(post) {
  if (!userStore.user?.id) return
  const liked = likedByUserMap.value[post.ID]
  if (liked) {
    await axios.delete(`http://localhost:8080/api/posts/${post.ID}/like`, {
      headers: { 'Content-Type': 'application/json' },
      data: { userId: userStore.user.id },
      withCredentials: true,
    })
    likedByUserMap.value[post.ID] = false
  } else {
    await axios.post(
      `http://localhost:8080/api/posts/${post.ID}/like`,
      { userId: userStore.user.id },
      { withCredentials: true },
    )
    likedByUserMap.value = true
  }
  await fetchPosts()
  await fetchLikes()
  // If the modal is open and this is the selected post, refresh its data
  if (showCommentModal.value && selectedPost.value && selectedPost.value.ID === post.ID) {
    const res = await axios.get(`http://localhost:8080/api/posts/${post.ID}`, {
      withCredentials: true,
    })
    selectedPost.value = res.data
  }
}

async function handleDelete(post) {
  await axios.delete(`http://localhost:8080/api/posts/${post.ID}`, { withCredentials: true })
  fetchPosts()
}

async function handleNewPost() {
  const formData = new FormData()
  formData.append('description', newPostContent.value)
  formData.append('userId', userStore.user?.id)
  if (imageFile.value) {
    formData.append('image', imageFile.value)
  }
  await axios.post('http://localhost:8080/api/posts', formData, {
    withCredentials: true,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  newPostContent.value = ''
  imageFile.value = null
  imagePreview.value = null
  showModal.value = false
  fetchPosts()
}

function onImageChange(e) {
  imageFile.value = e.target.files[0]
  if (imageFile.value) {
    imagePreview.value = URL.createObjectURL(imageFile.value)
  } else {
    imagePreview.value = null
  }
}

function removeImage() {
  imageFile.value = null
  imagePreview.value = null
}

function openCommentModal(post) {
  selectedPost.value = post
  commentInput.value = ''
  fetchComments(post.ID)
  showCommentModal.value = true
}

async function fetchComments(postId) {
  loadingComments.value = true
  try {
    const res = await axios.get(`http://localhost:8080/api/posts/${postId}`, {
      withCredentials: true,
    })
    // Comments from API now include UserName directly
    comments.value = res.data.Comments || []
  } catch (error) {
    console.error('Failed to fetch comments:', error)
    comments.value = []
  } finally {
    loadingComments.value = false
  }
}

async function addComment(commentContent) {
  if (!commentContent || !commentContent.trim() || !selectedPost.value) return

  if (!userStore.user?.id) {
    console.error('User not logged in, cannot comment.')
    return
  }

  try {
    const response = await axios.post(
      // Capture the response
      `http://localhost:8080/api/posts/${selectedPost.value.ID}/comments`,
      {
        content: commentContent,
        userId: userStore.user.id,
      },
      { withCredentials: true },
    )

    const newComment = response.data // The new comment from the API

    // Append the new comment to the existing list
    if (!Array.isArray(comments.value)) {
      comments.value = []
    }
    comments.value.push(newComment) // Add to the end. PostModal will reverse for display.

    // No longer re-fetching all comments:
    // await fetchComments(selectedPost.value.ID)
  } catch (error) {
    console.error('Failed to add comment:', error)
    // Optionally, inform the user that adding the comment failed.
  }
}

async function deleteComment(comment) {
  await axios.delete(
    `http://localhost:8080/api/posts/${selectedPost.value.ID}/comments/${comment.ID}`,
    { withCredentials: true },
  )
  fetchComments(selectedPost.value.ID)
}

function canEditOrDeleteComment(comment) {
  const user = JSON.parse(localStorage.getItem('user'))
  return user && user.id === comment.UserID
}

// Handles editing a comment (edit-in-place, robust, from scratch, with logging)
async function editComment({ comment, content }) {
  console.log('[editComment] called with:', {
    comment,
    content,
    editCommentContent: editCommentContent.value,
  })
  // Use the content from the event, or fallback to editCommentContent
  let newContent = content !== undefined ? content : editCommentContent.value
  console.log('[editComment] newContent before trim:', newContent)
  // Defensive: trim and check
  if (!newContent || !newContent.trim() || newContent === comment.Content) {
    console.warn('[editComment] Invalid or unchanged content, aborting.', { newContent, comment })
    editingCommentId.value = null
    editCommentContent.value = ''
    return
  }
  newContent = newContent.trim()
  console.log('[editComment] newContent after trim:', newContent)
  try {
    // Call the API to update the comment
    console.log(
      '[editComment] Sending PUT to API:',
      `http://localhost:8080/api/posts/${selectedPost.value.ID}/comments/${comment.ID}`,
      { content: newContent },
    )
    const response = await axios.put(
      `http://localhost:8080/api/posts/${selectedPost.value.ID}/comments/${comment.ID}`,
      { content: newContent },
      { withCredentials: true },
    )
    console.log('[editComment] API response:', response)
    // Optionally update the comment in the local array for instant UI feedback
    const idx = comments.value.findIndex((c) => c.ID === comment.ID)
    if (idx !== -1) {
      comments.value[idx].Content = newContent
      console.log('[editComment] Updated local comment:', comments.value[idx])
    }
  } catch (e) {
    // Optionally show error to user
    console.error('[editComment] Failed to edit comment', e)
  } finally {
    editingCommentId.value = null
    editCommentContent.value = ''
    // Optionally re-fetch comments for consistency
    console.log('[editComment] Fetching comments for post', selectedPost.value.ID)
    fetchComments(selectedPost.value.ID)
  }
}

function canEditPost(post) {
  return userStore.user && post.User?.ID === userStore.user.id
}

function startEditPost() {
  editingPost.value = true
  editPostContent.value = selectedPost.value.Description
}

function cancelEditPost() {
  editingPost.value = false
}

async function saveEditPost() {
  if (!editPostContent.value.trim()) return
  await axios.put(
    `http://localhost:8080/api/posts/${selectedPost.value.ID}`,
    { ...selectedPost.value, description: editPostContent.value },
    { withCredentials: true },
  )
  selectedPost.value.Description = editPostContent.value
  editingPost.value = false
  fetchPosts()
}

function startEditComment(comment) {
  editingCommentId.value = comment.ID
  editCommentContent.value = comment.Content
}

function cancelEditComment() {
  editingCommentId.value = null
  editCommentContent.value = ''
}

async function saveEditComment(comment) {
  if (!editCommentContent.value.trim() || editCommentContent.value === comment.Content) {
    editingCommentId.value = null
    return
  }
  await axios.put(
    `http://localhost:8080/api/posts/${selectedPost.value.ID}/comments/${comment.ID}`,
    { content: editCommentContent.value },
    { withCredentials: true },
  )
  editingCommentId.value = null
  editCommentContent.value = ''
  fetchComments(selectedPost.value.ID)
}
</script>
<style scoped></style>
