<template>
  <Teleport to="body">
    <transition name="modal-fade">
      <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center">
        <Backdrop :show="show" @click="$emit('close')" />
        <div
          class="bg-white rounded-xl shadow-lg p-6 w-auto relative z-50 flex flex-col items-center justify-center"
          @click.stop
        >
          <button
            @click="$emit('close')"
            class="absolute top-2 right-2 text-gray-400 hover:text-gray-700 text-2xl leading-none"
          >
            &times;
          </button>
          <slot />
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup>
import Backdrop from './Backdrop.vue'
const props = defineProps({
  show: Boolean,
})
</script>

<style scoped>
.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
.modal-fade-enter-to,
.modal-fade-leave-from {
  opacity: 1;
}
</style>
