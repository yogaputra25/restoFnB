<template>
  <span
    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-semibold"
    :style="{ backgroundColor: bgColor, color: textColor }"
  >
    {{ displayLabel }}
  </span>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  status: string
  variant?: 'order' | 'payment' | 'generic'
}>(), {
  variant: 'order',
})

const statusMap: Record<string, { bg: string; text: string; label: string }> = {
  pending:    { bg: '#FDEBD0', text: '#C65D3B', label: 'Pending' },
  confirmed:  { bg: '#E5E7EB', text: '#2F2F2F', label: 'Confirmed' },
  preparing:  { bg: '#FEF1E8', text: '#D97706', label: 'Preparing' },
  ready:      { bg: '#E8F5E9', text: '#2E7D32', label: 'Ready' },
  completed:  { bg: '#E8F5E9', text: '#2E7D32', label: 'Completed' },
  cancelled:  { bg: '#FDECEA', text: '#C0392B', label: 'Cancelled' },
  paid:       { bg: '#E8F5E9', text: '#2E7D32', label: 'Paid' },
  unpaid:     { bg: '#FDEBD0', text: '#C65D3B', label: 'Unpaid' },
  active:     { bg: '#E8F5E9', text: '#2E7D32', label: 'Active' },
  inactive:   { bg: '#F5F3EF', text: '#6B6B6B', label: 'Inactive' },
}

const mapped = computed(() => statusMap[props.status] || { bg: '#F5F3EF', text: '#6B6B6B', label: props.status })
const bgColor = computed(() => mapped.value.bg)
const textColor = computed(() => mapped.value.text)
const displayLabel = computed(() => mapped.value.label)
</script>
