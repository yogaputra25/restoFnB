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
  pending:    { bg: '#FEF3C7', text: '#92400E', label: 'Pending' },
  confirmed:  { bg: '#DBEAFE', text: '#1E40AF', label: 'Confirmed' },
  preparing:  { bg: '#FED7AA', text: '#9A3412', label: 'Preparing' },
  ready:      { bg: '#D8B4FE', text: '#6B21A8', label: 'Ready' },
  completed:  { bg: '#D1FAE5', text: '#065F46', label: 'Completed' },
  cancelled:  { bg: '#FEE2E2', text: '#991B1B', label: 'Cancelled' },
  paid:       { bg: '#A7F3D0', text: '#064E3B', label: 'Paid' },
  unpaid:     { bg: '#FEF3C7', text: '#92400E', label: 'Unpaid' },
  active:     { bg: '#D1FAE5', text: '#065F46', label: 'Active' },
  inactive:   { bg: '#FEE2E2', text: '#991B1B', label: 'Inactive' },
}

const mapped = computed(() => statusMap[props.status] || { bg: '#F3F4F6', text: '#374151', label: props.status })
const bgColor = computed(() => mapped.value.bg)
const textColor = computed(() => mapped.value.text)
const displayLabel = computed(() => mapped.value.label)
</script>
