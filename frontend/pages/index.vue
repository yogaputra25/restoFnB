<template>
  <div class="min-h-screen bg-[var(--color-background)]">
    <AppNavbar
      :cart-count="itemCount"
      @toggle-cart="toggleCart"
    />

    <!-- Hero Section -->
    <section
      data-aos="fade-up"
      class="relative overflow-hidden bg-[var(--color-primary)] text-white"
    >
      <div class="max-w-7xl mx-auto px-4 py-16 sm:py-20 lg:py-28 relative z-10">
        <div class="max-w-lg">
          <h1 class="font-heading text-hero leading-tight mb-4">
            Good Food,<br />Good Mood
          </h1>
          <p class="text-lg text-white/80 mb-8 max-w-md">
            Enjoy authentic flavors crafted with passion. Order your favorites and we'll handle the rest.
          </p>
          <div class="flex gap-4">
            <RouterLink
              to="/menu"
              class="inline-flex items-center gap-2 px-6 py-3 rounded-[var(--radius-button)] bg-white text-[var(--color-primary)] font-semibold text-sm transition-all duration-[var(--transition-base)] hover:scale-[1.03]"
            >
              Browse Menu
              <ArrowRight class="w-4 h-4" />
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Decorative background circles -->
      <div class="absolute -bottom-16 -right-16 w-64 h-64 rounded-full bg-white/5" />
      <div class="absolute -top-20 -right-8 w-48 h-48 rounded-full bg-white/5" />
    </section>

    <!-- Featured Section -->
    <section data-aos="fade-up" class="max-w-7xl mx-auto px-4 py-16">
      <h2 class="font-heading text-section text-[var(--color-text-primary)] mb-8">
        Featured Menu
      </h2>

      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <AppSkeleton v-for="i in 4" :key="i" variant="card" />
      </div>

      <div v-else-if="items.length === 0">
        <AppEmptyState
          title="No Items Yet"
          description="Check back soon for our featured dishes."
        />
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <AppCard
          v-for="item in items.slice(0, 8)"
          :key="item.id"
          :image="item.image_url"
          :title="item.name"
          :description="item.description"
          :price="formatPrice(item.price)"
          :show-add-to-cart="true"
          shadow
          @add-to-cart="addItem()"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ArrowRight } from 'lucide-vue-next'
interface MenuItem {
  id: string
  category_id?: string
  name: string
  description?: string
  price: number
  image_url?: string
}

const loading = ref(true)
const items = ref<MenuItem[]>([])
const itemCount = ref(0)

onMounted(async () => {
  try {
    const res = await $fetch('/api/menu-items/list?chain_id=demo')
    items.value = ((res as any).data || []).slice(0, 8)
  } catch {
    // silent fail — empty state handles it
  } finally {
    loading.value = false
  }
})

function formatPrice(price: number) {
  return `Rp${price.toLocaleString('id-ID')}`
}

function addItem() {
  itemCount.value++
}

function toggleCart() {
  // no-op for home page
}
</script>
