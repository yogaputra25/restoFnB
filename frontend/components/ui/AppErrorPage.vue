<template>
  <div class="min-h-screen bg-[var(--color-background)] flex items-center justify-center px-6">
    <div class="max-w-lg w-full text-center">
      <!-- Illustration -->
      <div v-if="$slots.illustration" class="mb-8">
        <slot name="illustration" />
      </div>
      <div v-else class="mb-8 w-40 h-40 mx-auto rounded-full bg-[var(--color-surface-secondary)] flex items-center justify-center">
        <svg
          v-if="type === '404'"
          class="w-20 h-20 text-[var(--color-border)]"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <svg
          v-else
          class="w-20 h-20 text-[var(--color-border)]"
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17.657 18.657A8 8 0 016.343 7.343S7 9 9 10c0-2 .5-5 2.986-7A14 14 0 0012 3a11 11 0 1011 11 11 11 0 00-1-3.5c-1 1-2.5 1.5-3.5 1.5 2 2 1.5 4.5 1.5 5.657z" />
        </svg>
      </div>

      <!-- Error code -->
      <p class="font-heading text-[100px] leading-none font-bold text-[var(--color-primary)] opacity-30 mb-4">
        {{ type === '404' ? '404' : '500' }}
      </p>

      <!-- Title -->
      <h1 class="font-heading text-[var(--font-size-page-title)] text-[var(--color-text-primary)] mb-3">
        {{ title }}
      </h1>

      <!-- Message -->
      <p class="text-[var(--color-text-secondary)] mb-8 max-w-sm mx-auto">
        {{ message }}
      </p>

      <!-- Actions -->
      <div class="flex items-center justify-center gap-4">
        <button
          class="px-6 py-2.5 rounded-[var(--radius-button)] bg-[var(--color-primary)] text-white font-semibold text-sm transition-all duration-[var(--transition-base)] hover:bg-[#a84e31] hover:scale-[1.03]"
          @click="goHome"
        >
          Back to Menu
        </button>
        <button
          v-if="type === '500'"
          class="px-6 py-2.5 rounded-[var(--radius-button)] border-2 border-[var(--color-border)] text-[var(--color-text-primary)] font-semibold text-sm transition-all duration-[var(--transition-base)] hover:border-[var(--color-primary)] hover:scale-[1.03]"
          @click="retry"
        >
          Try Again
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  type?: '404' | '500'
  title?: string
  message?: string
}>(), {
  type: '404',
  title: 'Dish Not Found',
  message: 'Looks like this dish wandered off the menu. Let\'s get you back to the good stuff.',
})

const router = useRouter()

function goHome() { router.push('/menu') }
function retry() { window.location.reload() }
</script>
