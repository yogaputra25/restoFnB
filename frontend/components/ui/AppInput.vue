<template>
  <div class="relative w-full">
    <div
      :class="[
        'relative rounded-[var(--radius-input)] border transition-all duration-[var(--transition-fast)]',
        focusState ? 'border-[var(--color-primary)] ring-1 ring-[var(--color-primary)]' : 'border-[var(--color-border)]',
        error ? 'border-[var(--color-danger)] ring-1 ring-[var(--color-danger)]' : '',
        success ? 'border-[var(--color-success)]' : '',
      ]"
    >
      <input
        :id="inputId"
        :type="type"
        :value="modelValue"
        :placeholder="floating ? ' ' : placeholder"
        :disabled="disabled"
        :required="required"
        class="peer w-full bg-transparent px-4 pt-5 pb-2 text-[var(--font-size-body)] text-[var(--color-text-primary)] outline-none disabled:opacity-50 rounded-[var(--radius-input)]"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        @focus="focusState = true"
        @blur="focusState = false"
      />

      <!-- Floating label -->
      <label
        v-if="floating"
        :for="inputId"
        :class="[
          'absolute left-4 transition-all duration-[var(--transition-fast)] pointer-events-none',
          'peer-focus:top-2 peer-focus:text-xs peer-focus:text-[var(--color-primary)]',
          'peer-placeholder-shown:top-4 peer-placeholder-shown:text-base peer-placeholder-shown:text-[var(--color-text-secondary)]',
          'top-2 text-xs',
          modelValue ? 'top-2 text-xs text-[var(--color-text-secondary)]' : '',
          error ? 'text-[var(--color-danger)]' : '',
        ]"
      >
        {{ label }}
      </label>

      <!-- Static label (non-floating) -->
      <label
        v-else-if="label"
        :for="inputId"
        class="block text-xs font-medium text-[var(--color-text-secondary)] mb-1"
      >
        {{ label }}
      </label>
    </div>

    <!-- Error message -->
    <p v-if="error" class="mt-1 text-xs text-[var(--color-danger)] flex items-center gap-1">
      <AlertTriangle class="w-3 h-3" /> {{ error }}
    </p>

    <!-- Hint text -->
    <p v-else-if="hint" class="mt-1 text-xs text-[var(--color-text-secondary)]">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { AlertTriangle } from 'lucide-vue-next'
const props = defineProps<{
  modelValue: string
  label?: string
  type?: string
  placeholder?: string
  floating?: boolean
  error?: string
  hint?: string
  disabled?: boolean
  required?: boolean
  success?: boolean
}>()

defineEmits<{ 'update:modelValue': [value: string] }>()

const focusState = ref(false)
const inputId = `input-${Math.random().toString(36).slice(2, 9)}`
</script>
