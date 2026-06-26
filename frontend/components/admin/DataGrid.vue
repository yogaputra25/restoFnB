<template>
  <div class="w-full">
    <!-- Toolbar slot -->
    <div v-if="$slots.toolbar" class="mb-3">
      <slot name="toolbar" />
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="space-y-3">
      <div v-for="i in 5" :key="i" class="h-12 rounded-[var(--radius-card)] bg-gradient-to-r from-[var(--color-surface-secondary)] via-white/40 to-[var(--color-surface-secondary)] bg-[length:200%_100%] animate-shimmer" />
    </div>

    <!-- Empty state -->
    <div v-else-if="rows.length === 0" class="text-center py-12">
      <slot name="empty">
        <p class="text-[var(--color-text-secondary)] text-sm">{{ emptyMessage || 'No data' }}</p>
      </slot>
    </div>

    <!-- Data rows -->
    <div v-else class="space-y-2">
      <!-- Header row -->
      <div class="flex items-center gap-3 px-4 py-2 text-xs font-semibold text-[var(--color-text-secondary)] uppercase tracking-wider">
        <div
          v-for="col in columns"
          :key="col.key"
          class="flex items-center gap-1 cursor-pointer select-none hover:text-[var(--color-text-primary)] transition-colors"
          :style="{ flex: col.flex || '1' }"
          @click="col.sortable && toggleSort(col.key)"
        >
          <span>{{ col.label }}</span>
          <svg v-if="col.sortable && sortKey === col.key" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="sortDir === 'asc'" stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 15l7-7 7 7" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>

      <!-- Data row -->
      <div
        v-for="(row, i) in sortedRows"
        :key="i"
        class="flex items-center gap-3 px-4 py-3 bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] transition-all duration-[var(--transition-fast)] hover:shadow-[var(--shadow-md)] hover:-translate-y-0.5 cursor-pointer"
        @click="$emit('rowClick', row)"
      >
        <div
          v-for="col in columns"
          :key="col.key"
          class="text-sm"
          :style="{ flex: col.flex || '1' }"
        >
          <slot :name="`cell-${col.key}`" :row="row" :value="row[col.key]">
            <span class="text-[var(--color-text-primary)]">{{ row[col.key] }}</span>
          </slot>
        </div>
        <div v-if="$slots.actions" class="flex-shrink-0">
          <slot name="actions" :row="row" />
        </div>
      </div>

      <!-- Pagination slot -->
      <div v-if="$slots.pagination" class="pt-2">
        <slot name="pagination" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Column {
  key: string
  label: string
  sortable?: boolean
  flex?: string
}

const props = withDefaults(defineProps<{
  columns: Column[]
  rows: Record<string, any>[]
  loading?: boolean
  emptyMessage?: string
  sortKey?: string
  sortDir?: 'asc' | 'desc'
}>(), {
  loading: false,
  sortKey: '',
  sortDir: 'asc',
})

const emit = defineEmits<{
  sort: [key: string, dir: 'asc' | 'desc']
  rowClick: [row: Record<string, any>]
}>()

const internalSortKey = ref(props.sortKey)
const internalSortDir = ref<'asc' | 'desc'>(props.sortDir)

function toggleSort(key: string) {
  if (internalSortKey.value === key) {
    internalSortDir.value = internalSortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    internalSortKey.value = key
    internalSortDir.value = 'asc'
  }
  emit('sort', internalSortKey.value, internalSortDir.value)
}

const sortedRows = computed(() => {
  if (!internalSortKey.value) return props.rows
  return [...props.rows].sort((a, b) => {
    const va = a[internalSortKey.value]
    const vb = b[internalSortKey.value]
    if (typeof va === 'number' && typeof vb === 'number') {
      return internalSortDir.value === 'asc' ? va - vb : vb - va
    }
    const sa = String(va || '').toLowerCase()
    const sb = String(vb || '').toLowerCase()
    return internalSortDir.value === 'asc' ? sa.localeCompare(sb) : sb.localeCompare(sa)
  })
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
