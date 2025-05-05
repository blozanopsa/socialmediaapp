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
          <div v-if="loading" class="flex justify-center py-8">
            <Loader />
          </div>
          <PostList
            v-else
            :posts="posts"
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
        <IKContext
          publicKey="public_GwIFTTS26dR3QXyxCoO+pKBdTXU="
          urlEndpoint="https://ik.imagekit.io/ph47nw0omb"
          :authenticator="authenticator"
        >
          <IKUpload
            :onSuccess="onSuccess"
            :onError="onError"
            :useUniqueFileName="true"
            :isPrivateFile="false"
            :showProgress="true"
            :validateFile="(file) => file.size < 5 * 1024 * 1024"
          />
        </IKContext>
        <div v-if="imageUrl" class="text-green-600 text-xs">Image uploaded!</div>
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
    <BaseModal :show="showCommentModal" @close="showCommentModal = false">
      <div v-if="selectedPost" class="w-[1000px] min-h-[400px] max-w-full">
        <!-- Post content and like button at the top -->
        <div class="mb-4 pb-4 border-b border-gray-200">
          <div class="flex items-center gap-2 mb-1">
            <span class="text-blue-700 text-lg">{{ selectedPost.User?.Name || 'Unknown' }}</span>
            <span class="text-xs text-gray-400 ml-auto">{{
              new Date(selectedPost.CreatedAt).toLocaleString()
            }}</span>
            <button
              v-if="canEditPost(selectedPost) && !editingPost"
              @click="startEditPost"
              class="ml-2 px-2 py-1 text-xs rounded bg-blue-100 text-blue-700 hover:bg-blue-200"
            >
              Edit
            </button>
          </div>
          <div v-if="editingPost" class="flex flex-col gap-2">
            <textarea
              v-model="editPostContent"
              rows="3"
              class="border border-gray-300 rounded-lg px-3 py-2 text-base focus:outline-none focus:ring-2 focus:ring-blue-400 resize-none"
            ></textarea>
            <div class="flex gap-2">
              <button
                @click="cancelEditPost"
                class="px-3 py-1 rounded bg-gray-200 hover:bg-gray-300"
              >
                Cancel
              </button>
              <button
                @click="saveEditPost"
                class="px-3 py-1 rounded bg-blue-600 text-white hover:bg-blue-700"
              >
                Save
              </button>
            </div>
          </div>
          <div v-else class="text-xl text-gray-800 font-medium mb-2">
            {{ selectedPost.Description }}
          </div>
          <div class="flex items-center gap-2">
            <button
              @click="
                (e) => {
                  if (selectedPost.User?.ID === userStore.user?.id) return
                  handleLike(selectedPost)
                  e.stopPropagation()
                }
              "
              :aria-label="likedByUserMap[String(selectedPost.ID)] ? 'Unlike' : 'Like'"
              class="relative group"
            >
              <svg
                v-if="
                  likedByUserMap[String(selectedPost.ID)] &&
                  selectedPost.User?.ID !== userStore.user?.id
                "
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
            </button>
            <span class="text-base text-gray-700 min-w-[2ch] text-center">{{
              selectedPost.likesCount || (selectedPost.Likes ? selectedPost.Likes.length : 0)
            }}</span>
          </div>
        </div>
        <!-- Comments section -->
        <h2 class="text-lg mb-2">Comments</h2>
        <div v-if="loadingComments" class="py-4 flex justify-center"><Loader /></div>
        <div v-else class="max-h-64 overflow-y-auto space-y-2 mb-4">
          <div
            v-for="comment in comments"
            :key="comment.ID"
            class="bg-gray-100 rounded-lg px-3 py-2 text-sm flex items-center justify-between"
          >
            <div class="flex-1">
              <span class="text-blue-700">{{
                userNames[comment.UserID] || 'User ' + comment.UserID
              }}</span>
              <span class="text-xs text-gray-400 ml-2">{{
                new Date(comment.CreatedAt).toLocaleString()
              }}</span>
              <div v-if="editingCommentId === comment.ID">
                <textarea
                  v-model="editCommentContent"
                  rows="2"
                  class="border border-gray-300 rounded-lg px-2 py-1 w-full text-sm focus:outline-none focus:ring-2 focus:ring-blue-400 resize-none mt-1"
                ></textarea>
                <div class="flex gap-2 mt-1">
                  <button
                    @click="cancelEditComment"
                    class="px-2 py-1 rounded bg-gray-200 hover:bg-gray-300 text-xs"
                  >
                    Cancel
                  </button>
                  <button
                    @click="saveEditComment(comment)"
                    class="px-2 py-1 rounded bg-blue-600 text-white hover:bg-blue-700 text-xs"
                  >
                    Save
                  </button>
                </div>
              </div>
              <div v-else>{{ comment.Content }}</div>
            </div>
            <div v-if="canEditOrDeleteComment(comment)" class="relative ml-2">
              <button
                @click="(e) => toggleCommentMenu(comment.ID, e)"
                class="p-1 rounded hover:bg-gray-200"
              >
                <!-- Burger icon -->
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
                    d="M4 8h16M4 16h16"
                  />
                </svg>
              </button>
              <teleport to="body">
                <div
                  v-if="openCommentMenuId === comment.ID"
                  class="comment-menu-tooltip fixed w-24 bg-white border rounded shadow z-[99999] flex flex-col"
                  :style="{ top: menuPosition.top + 'px', left: menuPosition.left + 'px' }"
                  @click.stop
                >
                  <button
                    @click="
                      () => {
                        startEditComment(comment)
                        openCommentMenuId = null
                      }
                    "
                    class="px-3 py-2 text-left text-blue-600 hover:bg-blue-50"
                  >
                    Edit
                  </button>
                  <button
                    @click="
                      () => {
                        deleteComment(comment)
                        openCommentMenuId = null
                      }
                    "
                    class="px-3 py-2 text-left text-red-600 hover:bg-red-50"
                  >
                    Delete
                  </button>
                </div>
              </teleport>
            </div>
          </div>
          <div v-if="comments.length === 0" class="text-gray-400 text-center">No comments yet.</div>
        </div>
        <form @submit.prevent="addComment" class="flex gap-2">
          <input
            v-model="commentInput"
            type="text"
            placeholder="Add a comment..."
            class="flex-1 border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-400"
          />
          <button
            type="submit"
            class="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition text-sm"
          >
            Post
          </button>
        </form>
      </div>
      <div v-else class="flex items-center justify-center min-h-[400px]">
        <Loader />
      </div>
    </BaseModal>
  </div>
</template>
<script setup>
import { ref, onMounted, onBeforeUnmount, nextTick, watch } from 'vue'
import Loader from '@/components/Loader.vue'
import PostList from '@/components/PostList.vue'
import BaseModal from '@/components/modals/BaseModal.vue'
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
const userNames = ref({})
const editingPost = ref(false)
const editPostContent = ref('')
const editingCommentId = ref(null)
const editCommentContent = ref('')
const imageFile = ref(null)

// Watch for modal open/close and update global modal state
watch([showModal, showCommentModal], ([modal, commentModal]) => {
  modalStore.isAnyModalOpen = modal || commentModal
})

async function fetchUserNames() {
  const res = await axios.get('http://localhost:8080/users/names', { withCredentials: true })
  userNames.value = Object.fromEntries(res.data.map((u) => [u.id, u.name]))
}

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
  fetchUserNames()
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

async function uploadImage() {
  if (!imageFile.value) return ''
  const formData = new FormData()
  formData.append('file', imageFile.value)
  const res = await fetch('https://uploadthing.com/api/upload', {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${import.meta.env.VITE_UPLOADTHING_TOKEN}`,
    },
    body: formData,
  })
  const data = await res.json()
  return data.url
}

async function handleNewPost() {
  let url = ''
  if (imageFile.value) {
    url = await uploadImage()
  }
  await axios.post(
    'http://localhost:8080/api/posts',
    {
      description: newPostContent.value,
      userId: userStore.user?.id,
      image_url: url,
    },
    { withCredentials: true },
  )
  newPostContent.value = ''
  imageFile.value = null
  showModal.value = false
  fetchPosts()
}

function onImageChange(e) {
  imageFile.value = e.target.files[0]
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
    comments.value = res.data.Comments || []
  } finally {
    loadingComments.value = false
  }
}

async function addComment() {
  if (!commentInput.value.trim() || !selectedPost.value) return
  await axios.post(
    `http://localhost:8080/api/posts/${selectedPost.value.ID}/comments`,
    {
      content: commentInput.value,
      userId: JSON.parse(localStorage.getItem('user'))?.id,
    },
    { withCredentials: true },
  )
  commentInput.value = ''
  fetchComments(selectedPost.value.ID)
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

async function editComment(comment) {
  const newContent = prompt('Edit your comment:', comment.Content)
  if (newContent && newContent.trim() && newContent !== comment.Content) {
    await axios.put(
      `http://localhost:8080/api/posts/${selectedPost.value.ID}/comments/${comment.ID}`,
      { content: newContent },
      { withCredentials: true },
    )
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
