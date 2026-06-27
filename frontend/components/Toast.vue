<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-[100] space-y-2">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          class="px-4 py-3 rounded-[var(--radius-button)] text-sm text-white transition-all"
          :style="{ boxShadow: 'var(--shadow-lg)' }"
          :class="{
            'bg-[var(--color-danger)]': toast.type === 'error',
            'bg-[var(--color-success)]': toast.type === 'success',
            'bg-[var(--color-primary)]': toast.type === 'info',
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
