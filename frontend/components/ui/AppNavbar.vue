<template>
  <header
    ref="navRef"
    class="sticky top-0 z-50 transition-all duration-[var(--transition-base)]"
    :class="scrolled ? 'bg-[var(--color-surface)] shadow-[var(--shadow-sm)]' : 'bg-[var(--color-surface)]/80 backdrop-blur-md'"
  >
    <div class="max-w-7xl mx-auto flex items-center justify-between px-4 py-3">
      <!-- Logo -->
      <RouterLink to="/" class="flex items-center gap-2 shrink-0">
        <span class="font-heading text-xl font-bold text-[var(--color-primary)]">restoFnB</span>
      </RouterLink>

      <!-- Search (desktop) -->
      <div class="hidden md:flex flex-1 max-w-md mx-6">
        <div class="relative w-full">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-text-secondary)]" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search menu..."
            class="w-full pl-9 pr-4 py-2 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)] transition-all duration-[var(--transition-fast)]"
          />
        </div>
      </div>

      <!-- Nav links -->
      <nav class="hidden lg:flex items-center gap-6 text-sm">
        <RouterLink
          v-for="link in navLinks"
          :key="link.label"
          :to="link.to"
          class="text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors duration-[var(--transition-fast)] font-medium"
          active-class="text-[var(--color-primary)]"
        >
          {{ link.label }}
        </RouterLink>
      </nav>

      <!-- Right actions -->
      <div class="flex items-center gap-3">
        <!-- Cart -->
        <button class="relative p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors" @click="$emit('toggleCart')">
          <ShoppingCart class="w-5 h-5 text-[var(--color-text-primary)]" />
          <span
            v-if="cartCount > 0"
            class="absolute -top-0.5 -right-0.5 w-5 h-5 rounded-full bg-[var(--color-primary)] text-white text-[10px] font-bold flex items-center justify-center animate-cart-bounce"
          >
            {{ cartCount > 99 ? '99+' : cartCount }}
          </span>
        </button>

        <!-- Profile -->
        <button class="p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors" @click="$emit('toggleProfile')">
          <User class="w-5 h-5 text-[var(--color-text-primary)]" />
        </button>

        <!-- Mobile menu toggle -->
        <button class="lg:hidden p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors" @click="$emit('toggleMobile')">
          <Menu class="w-5 h-5 text-[var(--color-text-primary)]" />
        </button>
      </div>
    </div>

    <!-- Mobile search -->
    <div v-if="mobileSearchOpen" class="md:hidden px-4 pb-3">
      <div class="relative w-full">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-text-secondary)]" />
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search menu..."
          class="w-full pl-9 pr-4 py-2 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
        />
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { Search, ShoppingCart, User, Menu } from 'lucide-vue-next'
defineProps<{
  cartCount?: number
}>()

defineEmits<{
  toggleCart: []
  toggleProfile: []
  toggleMobile: []
  search: [q: string]
}>()

const searchQuery = ref('')
const mobileSearchOpen = ref(false)
const scrolled = ref(false)
const navRef = ref<HTMLElement | null>(null)

const navLinks = [
  { label: 'Menu', to: '/menu' },
  { label: 'Promotions', to: '/promotions' },
]

onMounted(() => {
  const onScroll = () => { scrolled.value = window.scrollY > 50 }
  window.addEventListener('scroll', onScroll, { passive: true })
  onUnmounted(() => window.removeEventListener('scroll', onScroll))
})
</script>

<style scoped>
@keyframes cart-bounce {
  0%   { transform: scale(1); }
  40%  { transform: scale(1.3); }
  60%  { transform: scale(0.95); }
  100% { transform: scale(1); }
}
.animate-cart-bounce {
  animation: cart-bounce 0.3s ease;
}
</style>
