<template>
  <div :class="{ 'modal-blur': modalStore.isAnyModalOpen }">
    <div class="flex flex-col items-center mt-12"></div>
  </div>
</template>

<script setup>
import { useModalStore } from '@/stores/modal'
const modalStore = useModalStore()
// Add logic for profile here if needed
</script>

<style scoped></style>
