<template>
  <BaseModal :show="show" @close="onClose">
    <transition name="modal-fade">
      <div
        v-if="selectedPost"
        class="w-[1000px] min-h-[550px] max-w-full flex"
        :class="{ 'gap-x-6 p-4': selectedPost.ImageURL }"
        key="modal-content"
      >
        <!-- Image Section (Left) -->
        <div
          v-if="selectedPost.ImageURL"
          class="w-1/2 flex-shrink-0 bg-white rounded-lg flex items-center justify-center p-2 shadow-md"
        >
          <img
            :src="
              selectedPost.ImageURL.startsWith('/')
                ? selectedPost.ImageURL // Assuming 'post.ImageURL' was a typo and should be selectedPost.ImageURL
                : 'http://localhost:8080/' + selectedPost.ImageURL
            "
            alt="Post image"
            class="max-h-[530px] w-full object-contain rounded"
          />
        </div>

        <!-- Content Section (Right or Full Width) -->
        <div :class="['flex flex-col', selectedPost.ImageURL ? 'w-1/2' : 'w-full']">
          <!-- Post details: user, date, edit button -->
          <div
            class="pb-4 border-b border-gray-200"
            :class="{ 'mb-4': selectedPost.ImageURL, 'border-b border-gray-200': !selectedPost.ImageURL }"
          >
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
            <!-- Editing Post UI -->
            <div v-if="editingPost" class="flex flex-col gap-2 mt-2">
              <textarea
                :value="editPostContent"
                @input="(e) => emit('update:editPostContent', e.target.value)"
                rows="3"
                class="border border-gray-300 rounded-lg px-3 py-2 text-base focus:outline-none focus:ring-2 focus:ring-blue-400 resize-none"
              ></textarea>
              <div class="flex gap-2">
                <button @click="cancelEditPost" class="px-3 py-1 rounded bg-gray-200 hover:bg-gray-300">Cancel</button>
                <button @click="saveEditPost" class="px-3 py-1 rounded bg-blue-600 text-white hover:bg-blue-700">Save</button>
              </div>
            </div>
            <!-- Post Description -->
            <div v-else class="text-xl text-gray-800 font-medium my-2">
              {{ selectedPost.Description }}
            </div>

            <!-- Like Button and Count -->
            <div class="flex items-center gap-2 mt-2">
              <button
                @click="onLikeClick"
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

          <!-- Comments Section -->
          <div class="flex-grow flex flex-col min-h-0">
            <h2 class="text-lg mb-2 sticky top-0 bg-white pb-2 z-10 pt-4 border-t border-gray-200">Comments</h2>
            <div class="flex-grow overflow-y-auto pr-1 space-y-2">
              <div v-if="loadingComments" class="py-4 flex justify-center"><Loader /></div>
              <div v-else>
                <div v-if="comments.length > 0" class="space-y-2">
                  <div
                    v-for="comment in comments.slice().reverse()"
                    :key="comment.ID"
                    class="bg-gray-100 rounded-lg px-3 py-2 text-sm flex items-center justify-between"
                  >
                    <div class="flex-1">
                      <span class="text-blue-700">{{ comment?.UserName }}</span>
                      <span class="text-xs text-gray-400 ml-2">{{
                        new Date(comment.CreatedAt).toLocaleString()
                      }}</span>
                      <div v-if="editingCommentId === comment.ID && userStore.user?.id === comment.UserID">
                        <textarea
                          :value="localEditCommentContent"
                          @input="onEditCommentInput"
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
                    <div v-if="userStore.user?.id === comment.UserID" class="relative ml-2">
                      <button
                        @click="(e) => toggleCommentMenu(comment.ID, e)"
                        class="p-1 rounded hover:bg-gray-200"
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
                            @click="() => startEditComment(comment)"
                            class="px-3 py-2 text-left text-blue-600 hover:bg-blue-50"
                          >
                            Edit
                          </button>
                          <button
                            @click="() => deleteComment(comment)"
                            class="px-3 py-2 text-left text-red-600 hover:bg-red-50"
                          >
                            Delete
                          </button>
                        </div>
                      </teleport>
                    </div>
                  </div>
                </div>
                <div v-if="!loadingComments && comments.length === 0" class="text-gray-400 text-center py-4">No comments yet.</div>
              </div>
            </div>
            <!-- Add Comment Form -->
            <form @submit.prevent="addComment" class="flex gap-2 mt-auto pt-4 border-t border-gray-200">
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
        </div>
      </div>
      <div v-else class="flex items-center justify-center min-h-[400px]">
        <Loader />
      </div>
    </transition>
  </BaseModal>
</template>

<script setup>
import BaseModal from './BaseModal.vue'
import Loader from '../Loader.vue'
import { useUserStore } from '@/stores/user'
import { ref, computed, watch } from 'vue' // Removed 'watch' if no other watchers are left, but keeping for now as other watchers exist

const props = defineProps({
  show: Boolean,
  selectedPost: Object,
  comments: Array,
  loadingComments: Boolean,
  // userNames: Object, // Removed userNames prop
  editingPost: Boolean,
  editPostContent: String,
  editingCommentId: Number,
  editCommentContent: String,
  openCommentMenuId: Number,
  menuPosition: Object,
  likedByUserMap: Object,
})
const commentInput = ref('')
const emit = defineEmits([
  'close',
  'like',
  'add-comment',
  'edit-post',
  'cancel-edit-post',
  'save-edit-post',
  'start-edit-post',
  'edit-comment',
  'cancel-edit-comment',
  'save-edit-comment',
  'start-edit-comment',
  'delete-comment',
  'toggle-comment-menu',
  'update:editCommentContent',
])
const userStore = useUserStore()

function onClose() {
  emit('close')
}
function addComment() {
  if (!commentInput.value.trim()) return
  emit('add-comment', commentInput.value)
  commentInput.value = ''
}
function onLikeClick(e) {
  if (props.selectedPost.User?.ID === userStore.user?.id) return
  emit('like', props.selectedPost)
  if (e) e.stopPropagation()
}
function startEditPost() {
  emit('start-edit-post')
}
function cancelEditPost() {
  emit('cancel-edit-post')
}
function saveEditPost() {
  emit('save-edit-post')
}
function startEditComment(comment) {
  emit('start-edit-comment', comment)
}
function cancelEditComment() {
  emit('cancel-edit-comment')
}
// Local state for editing comment content
const localEditCommentContent = ref('')

// Watch for prop changes to sync local state
watch(
  () => props.editCommentContent,
  (val) => {
    localEditCommentContent.value = val || ''
  },
  { immediate: true },
)

function onEditCommentInput(e) {
  localEditCommentContent.value = e.target.value
  emit('update:editCommentContent', localEditCommentContent.value)
}
function saveEditComment(comment) {
  // Emit both the comment object and the new content for parent to handle
  emit('edit-comment', { comment, content: localEditCommentContent.value })
}
function deleteComment(comment) {
  emit('delete-comment', comment)
}
function toggleCommentMenu(id, event) {
  emit('toggle-comment-menu', id, event)
}
function canEditOrDeleteComment(comment) {
  const user = JSON.parse(localStorage.getItem('user'))
  return user && user.id === comment.UserID
}
function canEditPost(post) {
  return userStore.user && post.User?.ID === userStore.user.id
}
</script>
