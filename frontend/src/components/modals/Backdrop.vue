<template>
  <transition name="backdrop-fade">
    <div
      v-if="$attrs.show !== false"
      class="fixed inset-0 z-40 bg-black/40 backdrop-blur-sm"
      @click="$emit('click')"
    ></div>
  </transition>
</template>

<script setup>
// No props needed
</script>

<style scoped>
.backdrop-fade-enter-active,
.backdrop-fade-leave-active {
  transition:
    opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1),
    backdrop-filter 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.backdrop-fade-enter-from,
.backdrop-fade-leave-to {
  opacity: 0;
  backdrop-filter: blur(0px);
}
.backdrop-fade-enter-to,
.backdrop-fade-leave-from {
  opacity: 1;
  backdrop-filter: blur(8px);
}
</style>
