<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center" @click.self="cancel">
        <div class="absolute inset-0 bg-black/40" />
        <div class="relative bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-lg)] w-full max-w-md mx-4 p-6" @click.stop>
          <h3 class="font-heading text-[var(--font-size-section)] text-[var(--color-text-primary)] mb-2">{{ title }}</h3>
          <div class="text-[var(--color-text-secondary)] text-sm mb-6">
            <slot />
          </div>
          <div class="flex gap-3 justify-end">
            <button
              class="px-4 py-2 rounded-[var(--radius-button)] border border-[var(--color-border)] text-[var(--color-text-primary)] text-sm font-medium transition-all duration-[var(--transition-fast)] hover:bg-[var(--color-surface-secondary)]"
              @click="cancel"
            >
              {{ cancelText }}
            </button>
            <button
              :class="['px-4 py-2 rounded-[var(--radius-button)] text-white text-sm font-semibold transition-all duration-[var(--transition-fast)] hover:scale-[1.03]',
                variant === 'danger' ? 'bg-[var(--color-danger)] hover:bg-[#a33225]' : 'bg-[var(--color-primary)] hover:bg-[#a84e31]']"
              :disabled="loading"
              @click="confirm"
            >
              <span v-if="loading" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin mr-2 align-middle" />
              {{ loading ? loadingText || confirmText : confirmText }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  open: boolean
  title?: string
  confirmText?: string
  cancelText?: string
  variant?: 'default' | 'danger'
  loading?: boolean
  loadingText?: string
}>(), {
  title: 'Confirm',
  confirmText: 'Confirm',
  cancelText: 'Cancel',
  variant: 'default',
  loading: false,
})

const emit = defineEmits<{ confirm: []; cancel: [] }>()

function confirm() { emit('confirm') }
function cancel() { emit('cancel') }

onMounted(() => {
  const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') cancel() }
  document.addEventListener('keydown', handler)
  onUnmounted(() => document.removeEventListener('keydown', handler))
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-active > div:last-child,
.modal-leave-active > div:last-child {
  transition: transform 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from > div:last-child {
  transform: scale(0.95);
}
.modal-leave-to > div:last-child {
  transform: scale(0.95);
}
</style>
