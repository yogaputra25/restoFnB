<template>
  <div class="flex flex-col items-center justify-center py-16 px-6 text-center">
    <!-- Illustration slot -->
    <div v-if="$slots.illustration" class="mb-6">
      <slot name="illustration" />
    </div>
    <div v-else class="mb-6 w-32 h-32 rounded-full bg-[var(--color-surface-secondary)] flex items-center justify-center">
      <svg class="w-16 h-16 text-[var(--color-border)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
      </svg>
    </div>

    <!-- Title -->
    <h3 class="font-heading text-[var(--font-size-section)] text-[var(--color-text-primary)] mb-2">
      {{ title }}
    </h3>

    <!-- Description -->
    <p class="text-[var(--color-text-secondary)] max-w-md mb-6">
      {{ description }}
    </p>

    <!-- Action button -->
    <div v-if="$slots.action || actionLabel">
      <slot name="action">
        <button
          class="px-6 py-2.5 rounded-[var(--radius-button)] bg-[var(--color-primary)] text-white font-semibold text-sm transition-all duration-[var(--transition-base)] hover:bg-[#a84e31] hover:scale-[1.03]"
          @click="$emit('action')"
        >
          {{ actionLabel }}
        </button>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
  description?: string
  actionLabel?: string
}>()

defineEmits<{ action: [] }>()
</script>
