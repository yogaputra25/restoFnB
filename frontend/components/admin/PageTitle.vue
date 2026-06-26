<template>
  <div class="mb-6">
    <h1 class="font-heading text-[var(--font-size-page-title)] text-[var(--color-text-primary)] leading-tight">{{ title }}</h1>
    <p v-if="description" class="text-[var(--color-text-secondary)] text-sm mt-1">{{ description }}</p>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string
  description?: string
}>()
</script>
