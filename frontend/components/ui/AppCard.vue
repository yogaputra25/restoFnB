<template>
  <div
    class="group bg-[var(--color-surface)] rounded-[var(--radius-card)] overflow-hidden flex flex-col transition-all duration-[var(--transition-base)] hover:-translate-y-1 hover:shadow-[var(--shadow-lg)]"
    :class="[shadow ? 'shadow-[var(--shadow-sm)]' : '']"
  >
    <!-- Image -->
    <div class="relative overflow-hidden aspect-[4/3] bg-[var(--color-surface-secondary)]">
      <img
        v-if="image"
        :src="image"
        :alt="title"
        loading="lazy"
        class="w-full h-full object-cover transition-transform duration-[var(--transition-base)] group-hover:scale-105"
      />
      <div v-else class="w-full h-full flex items-center justify-center text-4xl text-[var(--color-border)]">
        🍽
      </div>

      <!-- Favorite button -->
      <button
        v-if="showFavorite"
        class="absolute top-3 right-3 w-9 h-9 rounded-full bg-white/80 backdrop-blur-sm flex items-center justify-center transition-all duration-[var(--transition-fast)] hover:scale-110"
        :class="{ 'animate-heart-bounce': isFavorited }"
        @click.stop="toggleFavorite"
      >
        <svg
          class="w-5 h-5 transition-colors"
          :class="isFavorited ? 'text-[var(--color-danger)] fill-[var(--color-danger)]' : 'text-[var(--color-text-secondary)] fill-none'"
          stroke="currentColor"
          stroke-width="2"
          viewBox="0 0 24 24"
          fill="none"
        >
          <path stroke-linecap="round" stroke-linejoin="round" d="M20.84 4.61a5.5 5.5 0 00-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 00-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 000-7.78z" />
        </svg>
      </button>
    </div>

    <!-- Content -->
    <div class="p-4 flex-1 flex flex-col gap-1.5">
      <!-- Category badge -->
      <span v-if="category" class="self-start px-2.5 py-0.5 rounded-full text-[11px] font-medium bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)]">
        {{ category }}
      </span>

      <!-- Title -->
      <h3 class="font-heading text-[var(--font-size-card-title)] text-[var(--color-text-primary)] leading-tight line-clamp-1">
        {{ title }}
      </h3>

      <!-- Description -->
      <p v-if="description" class="text-[var(--font-size-small)] text-[var(--color-text-secondary)] line-clamp-2 leading-relaxed">
        {{ description }}
      </p>

      <!-- Rating + Price row -->
      <div class="flex items-center justify-between mt-auto pt-2">
        <div v-if="rating" class="flex items-center gap-1 text-sm text-[var(--color-accent)]">
          <svg class="w-4 h-4 fill-current" viewBox="0 0 20 20">
            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
          </svg>
          <span class="font-semibold">{{ rating }}</span>
        </div>

        <span class="font-bold text-[var(--color-primary)] text-lg">{{ price }}</span>
      </div>

      <!-- Action slot / Add to Cart -->
      <div v-if="$slots.action || showAddToCart" class="mt-2">
        <slot name="action">
          <button
            v-if="showAddToCart"
            class="w-full py-2.5 rounded-[var(--radius-button)] bg-[var(--color-primary)] text-white font-semibold text-sm transition-all duration-[var(--transition-fast)] hover:bg-[#a84e31] hover:scale-[1.02] active:scale-[0.98]"
            @click="$emit('addToCart')"
          >
            + Add to Cart
          </button>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  image?: string
  title: string
  description?: string
  category?: string
  rating?: string | number
  price: string
  shadow?: boolean
  showFavorite?: boolean
  showAddToCart?: boolean
  favorited?: boolean
}>()

defineEmits<{
  favorite: [val: boolean]
  addToCart: []
}>()

const isFavorited = ref(false)
watch(() => useAttrs().favorited, (v) => { isFavorited.value = !!v }, { immediate: true })

function toggleFavorite() {
  isFavorited.value = !isFavorited.value
}
</script>

<style scoped>
.animate-heart-bounce {
  animation: heart-bounce 0.3s ease;
}
</style>
