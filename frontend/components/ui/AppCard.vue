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
      <div v-else class="w-full h-full flex items-center justify-center text-[var(--color-border)]">
        <UtensilsCrossed class="w-8 h-8" />
      </div>

      <!-- Favorite button -->
      <button
        v-if="showFavorite"
        class="absolute top-3 right-3 w-9 h-9 rounded-full bg-white/80 backdrop-blur-sm flex items-center justify-center transition-all duration-[var(--transition-fast)] hover:scale-110"
        :class="{ 'animate-heart-bounce': isFavorited }"
        @click.stop="toggleFavorite"
      >
        <Heart
          class="w-5 h-5 transition-colors"
          :class="isFavorited ? 'text-[var(--color-danger)] fill-[var(--color-danger)]' : 'text-[var(--color-text-secondary)]'"
          stroke-width="2"
        />
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
          <Star class="w-4 h-4 fill-current" />
          <span class="font-semibold">{{ rating }}</span>
        </div>

        <span class="font-bold text-[var(--color-primary)] text-lg">{{ price }}</span>
      </div>

      <!-- Variants -->
      <div v-if="variants && variants.length > 0" class="mt-2">
        <div class="flex flex-wrap gap-1.5">
          <button
            v-for="v in variants"
            :key="v.id"
            @click.stop="selectedVariantId = v.id; selectedVariantName = v.name"
            class="px-2.5 py-1 rounded-[var(--radius-button)] text-xs font-medium border transition-all"
            :class="selectedVariantId === v.id ? 'bg-[var(--color-primary)] text-white border-[var(--color-primary)]' : 'bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)] border-[var(--color-border)] hover:border-[var(--color-primary)]'"
          >
            {{ v.name }}<span v-if="v.price_adjustment"> +Rp{{ v.price_adjustment }}</span>
          </button>
        </div>
      </div>

      <!-- Action slot / Add to Cart -->
      <div v-if="$slots.action || showAddToCart" class="mt-2">
        <slot name="action">
          <button
            v-if="showAddToCart"
            class="w-full py-2.5 rounded-[var(--radius-button)] bg-[var(--color-primary)] text-white font-semibold text-sm transition-all duration-[var(--transition-base)] hover:bg-[#a84e31] hover:scale-[1.03] active:scale-[0.98]"
            @click="$emit('addToCart', selectedVariantId)"
          >
            + Add to Cart
          </button>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { UtensilsCrossed, Heart, Star } from 'lucide-vue-next'

interface Variant { id: string; name: string; price_adjustment?: number }

const props = defineProps<{
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
  variants?: Variant[]
}>()

const emit = defineEmits<{
  favorite: [val: boolean]
  addToCart: [variantId?: string]
}>()

const selectedVariantId = ref<string | undefined>(undefined)
const selectedVariantName = ref('')

// Reset variant when variants change
watch(() => props.variants, () => { selectedVariantId.value = undefined; selectedVariantName.value = '' })

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
