<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-[100] space-y-2">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="px-4 py-3 rounded-lg shadow-lg text-sm text-white transition-all"
          :class="{
            'bg-red-600': toast.type === 'error',
            'bg-green-600': toast.type === 'success',
            'bg-blue-600': toast.type === 'info',
          }"
        >
          {{ toast.message }}
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
const { toasts } = useToast()
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}
.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}
.toast-leave-to {
  opacity: 0;
  transform: translateX(30px);
}
</style>
