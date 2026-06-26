<template>
  <footer class="h-12 flex items-center justify-center gap-3 border-t border-[var(--color-border)] bg-[var(--color-surface)] text-xs text-[var(--color-text-secondary)]">
    <span>© {{ year }} RestoFnB</span>
    <span class="w-1 h-1 rounded-full bg-[var(--color-border)]" />
    <span>v{{ version }}</span>
    <span class="w-1 h-1 rounded-full bg-[var(--color-border)]" />
    <span class="px-1.5 py-0.5 rounded text-[10px] font-medium bg-[var(--color-surface-secondary)]">{{ environment }}</span>
  </footer>
</template>

<script setup lang="ts">
const year = ref(new Date().getFullYear())
const version = ref(process.env.APP_VERSION || '1.0.0')
const environment = ref(process.env.NODE_ENV || 'development')
</script>
