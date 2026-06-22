<template>
  <div v-if="totalPages > 1" class="flex items-center justify-center gap-3 py-4">
    <button
      @click="$emit('page-change', currentPage - 1)"
      :disabled="currentPage <= 1"
      class="px-3 py-1.5 rounded text-sm transition"
      :class="currentPage <= 1 ? 'bg-dark-700 text-dark-500 cursor-not-allowed' : 'bg-dark-700 text-dark-300 hover:bg-dark-600 hover:text-white'"
    >Prev</button>

    <span class="text-dark-300 text-sm">
      Page {{ currentPage }} of {{ totalPages }}
    </span>

    <button
      @click="$emit('page-change', currentPage + 1)"
      :disabled="currentPage >= totalPages"
      class="px-3 py-1.5 rounded text-sm transition"
      :class="currentPage >= totalPages ? 'bg-dark-700 text-dark-500 cursor-not-allowed' : 'bg-dark-700 text-dark-300 hover:bg-dark-600 hover:text-white'"
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
