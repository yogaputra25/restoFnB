<template>
  <div v-if="totalPages > 1" class="flex items-center justify-center gap-3 py-4">
    <button
      @click="$emit('page-change', currentPage - 1)"
      :disabled="currentPage <= 1"
      class="px-3 py-1.5 rounded-[var(--radius-button)] text-sm transition-all duration-[var(--transition-base)] hover:scale-[1.03]"
      :class="currentPage <= 1 ? 'bg-[var(--color-surface-secondary)] text-[var(--color-border)] cursor-not-allowed' : 'bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)] hover:bg-[var(--color-border)] hover:text-[var(--color-text-primary)]'"
    >Prev</button>

    <span class="text-[var(--color-text-secondary)] text-sm">
      Page {{ currentPage }} of {{ totalPages }}
    </span>

    <button
      @click="$emit('page-change', currentPage + 1)"
      :disabled="currentPage >= totalPages"
      class="px-3 py-1.5 rounded-[var(--radius-button)] text-sm transition-all duration-[var(--transition-base)] hover:scale-[1.03]"
      :class="currentPage >= totalPages ? 'bg-[var(--color-surface-secondary)] text-[var(--color-border)] cursor-not-allowed' : 'bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)] hover:bg-[var(--color-border)] hover:text-[var(--color-text-primary)]'"
    >Next</button>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  currentPage: number
  totalItems: number
  limit: number
}>()

defineEmits<{
  'page-change': [page: number]
}>()

const totalPages = computed(() => Math.max(1, Math.ceil(props.totalItems / props.limit)))
</script>
