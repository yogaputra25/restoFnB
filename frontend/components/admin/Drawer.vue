<template>
  <Teleport to="body">
    <Transition name="drawer-backdrop">
      <div v-if="open" class="fixed inset-0 z-50" @click="close">
        <div class="absolute inset-0 bg-black/40" />
        <div
          class="absolute top-0 right-0 h-full bg-[var(--color-surface)] shadow-[var(--shadow-lg)] flex flex-col overflow-hidden"
          :style="{ width: typeof width === 'number' ? `${width}px` : width }"
          @click.stop
        >
          <!-- Header -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-[var(--color-border)] shrink-0">
            <h2 class="font-heading text-[var(--font-size-section)] text-[var(--color-text-primary)]">{{ title }}</h2>
            <button class="p-2 rounded-full hover:bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)] transition-colors" @click="close">
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Body -->
          <div class="flex-1 overflow-y-auto px-6 py-4">
            <slot />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { X } from 'lucide-vue-next'

withDefaults(defineProps<{
  open: boolean
  title?: string
  width?: number | string
}>(), {
  title: '',
  width: 480,
})

const emit = defineEmits<{ close: [] }>()

function close() {
  emit('close')
}

onMounted(() => {
  const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') close() }
  document.addEventListener('keydown', handler)
  onUnmounted(() => document.removeEventListener('keydown', handler))
})
</script>

<style scoped>
.drawer-backdrop-enter-active,
.drawer-backdrop-leave-active {
  transition: opacity 0.25s ease;
}
.drawer-backdrop-enter-from,
.drawer-backdrop-leave-to {
  opacity: 0;
}
.drawer-backdrop-enter-active > div:last-child {
  transition: transform 0.3s ease-out;
}
.drawer-backdrop-leave-active > div:last-child {
  transition: transform 0.2s ease-in;
}
.drawer-backdrop-enter-from > div:last-child,
.drawer-backdrop-leave-to > div:last-child {
  transform: translateX(100%);
}
</style>
