<template>
  <div
    class="bg-[var(--color-surface)] rounded-[var(--radius-card)] p-5 shadow-[var(--shadow-sm)] transition-all duration-[var(--transition-base)] hover:-translate-y-0.5 hover:shadow-[var(--shadow-md)]"
    data-aos="fade-up"
    :style="{ transitionDelay: `${delay}ms` }"
  >
    <div class="flex items-start justify-between mb-3">
      <div
        class="w-10 h-10 rounded-full flex items-center justify-center"
        :style="{ backgroundColor: `${color}15`, color }"
      >
        <component :is="iconComponent" class="w-5 h-5" />
      </div>
      <span v-if="trend" class="flex items-center gap-1 text-xs font-semibold" :class="trendUp ? 'text-[var(--color-success)]' : 'text-[var(--color-danger)]'">
        <svg v-if="trendUp" class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 10l7-7m0 0l7 7m-7-7v18" /></svg>
        <svg v-else class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M19 14l-7 7m0 0l-7-7m7 7V3" /></svg>
        {{ trend }}
      </span>
    </div>
    <p class="font-heading text-[var(--font-size-page-title)] text-[var(--color-text-primary)] leading-tight">{{ value }}</p>
    <p class="text-[var(--font-size-small)] text-[var(--color-text-secondary)] mt-1">{{ label }}</p>
  </div>
</template>

<script setup lang="ts">
import { h } from 'vue'
import * as LucideIcons from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  icon: string
  label: string
  value: string
  trend?: string
  trendUp?: boolean
  color?: string
  delay?: number
}>(), {
  color: 'var(--color-primary)',
  delay: 0,
})

const iconComponent = computed(() => {
  const icons = LucideIcons as Record<string, any>
  const name = props.icon.charAt(0).toUpperCase() + props.icon.slice(1) + 'Icon'
  return icons[name] || icons['CircleIcon'] || 'span'
})
</script>
