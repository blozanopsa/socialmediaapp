import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useModalStore = defineStore('modal', () => {
  const isAnyModalOpen = ref(false)
  return { isAnyModalOpen }
})
