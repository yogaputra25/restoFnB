<template>
  <div v-if="repeat > 1" class="space-y-3">
    <div
      v-for="i in repeat"
      :key="i"
      :style="variantStyle"
      class="animate-shimmer rounded-[var(--radius-card)] bg-gradient-to-r from-[var(--color-surface-secondary)] via-[color-mix(in_srgb,var(--color-surface-secondary),white_40%)] to-[var(--color-surface-secondary)] bg-[length:200%_100%]"
    />
  </div>
  <div
    v-else
    :style="variantStyle"
    class="animate-shimmer rounded-[var(--radius-card)] bg-gradient-to-r from-[var(--color-surface-secondary)] via-[color-mix(in_srgb,var(--color-surface-secondary),white_40%)] to-[var(--color-surface-secondary)] bg-[length:200%_100%]"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  variant?: 'text' | 'card' | 'image' | 'circle'
  width?: string
  height?: string
  borderRadius?: string
  repeat?: number
}>(), {
  variant: 'text',
  repeat: 1,
})

const variantStyle = computed(() => {
  switch (props.variant) {
    case 'card':
      return {
        width: props.width || '100%',
        height: props.height || '280px',
        borderRadius: props.borderRadius || 'var(--radius-card)',
      }
    case 'image':
      return {
        width: props.width || '100%',
        height: props.height || '180px',
        borderRadius: props.borderRadius || 'var(--radius-image)',
      }
    case 'circle':
      return {
        width: props.width || '48px',
        height: props.height || '48px',
        borderRadius: '50%',
      }
    case 'text':
    default:
      return {
        width: props.width || '100%',
        height: props.height || '16px',
        borderRadius: props.borderRadius || '8px',
      }
  }
})
</script>

<style scoped>
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}
.animate-shimmer {
  animation: shimmer 1.5s ease-in-out infinite;
}
</style>
