<template>
  <div>
    <PostCard
      v-for="post in posts"
      :key="post.ID"
      :post="post"
      :liked-by-user="!!likedByUserMap[String(post.ID)]"
      @like="$emit('like', post)"
      @delete="$emit('delete', post)"
      @comment="$emit('comment', post)"
    />
  </div>
</template>

<script setup>
import PostCard from './PostCard.vue'
const props = defineProps({
  posts: Array,
  likedByUserMap: Object, // { [postId]: true/false }
})
</script>
