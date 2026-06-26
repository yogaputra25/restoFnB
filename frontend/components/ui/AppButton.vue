<template>
  <component
    :is="to ? 'RouterLink' : 'button'"
    :to="to"
    :disabled="disabled || loading"
    :class="[
      'inline-flex items-center justify-center gap-2 font-semibold transition-all duration-[var(--transition-base)] focus:outline-none focus-visible:ring-2 focus-visible:ring-[var(--color-primary)] focus-visible:ring-offset-2',
      variantClasses[variant],
      sizeClasses[size],
      (disabled || loading) ? 'opacity-50 cursor-not-allowed pointer-events-none' : 'hover:scale-[1.03] active:scale-[0.98]',
    ]"
    @click="$emit('click', $event)"
  >
    <!-- Loading spinner -->
    <span v-if="loading" class="inline-block w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin" />
    <!-- Icon slot -->
    <span v-if="icon && !loading" class="inline-flex shrink-0">
      <slot name="icon">{{ icon }}</slot>
    </span>
    <slot />
  </component>
</template>

<script setup lang="ts">
defineProps<{
  variant?: 'primary' | 'secondary' | 'danger' | 'success'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  icon?: string
  to?: string | Record<string, unknown>
}>()

defineEmits<{ click: [e: MouseEvent] }>()

const variantClasses: Record<string, string> = {
  primary: 'bg-[var(--color-primary)] text-white hover:bg-[#a84e31]',
  secondary: 'bg-transparent text-[var(--color-text-primary)] border-2 border-[var(--color-border)] hover:border-[var(--color-primary)] hover:text-[var(--color-primary)]',
  danger: 'bg-[var(--color-danger)] text-white hover:bg-[#a33225]',
  success: 'bg-[var(--color-success)] text-white hover:bg-[#256d29]',
}

const sizeClasses: Record<string, string> = {
  sm: 'px-4 py-1.5 text-sm rounded-[var(--radius-button)]',
  md: 'px-6 py-2.5 text-base rounded-[var(--radius-button)]',
  lg: 'px-8 py-3 text-lg rounded-[var(--radius-button)]',
}
</script>
