<template>
  <div class="min-h-screen bg-gray-50">
    <PromoCarousel v-if="promos.length > 0" :promos="promos" />

    <div class="lg:flex max-w-7xl mx-auto">
      <div class="flex-1 min-w-0 p-3 sm:p-4 lg:p-6 pb-16 lg:pb-6">
        <section v-for="cat in categories" :key="cat.id" class="mb-6 sm:mb-8">
          <h2 class="text-lg font-bold text-gray-800 border-b border-gray-300 pb-1 mb-3">{{ cat.name }}</h2>
          <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
            <div v-for="item in getItemsByCategory(cat.id)" :key="item.id"
              class="bg-white rounded-xl shadow-sm overflow-hidden flex flex-col hover:shadow-md transition-shadow">
              <div class="h-28 sm:h-32 bg-gray-100 flex items-center justify-center">
                <img v-if="item.image_url" :src="item.image_url" :alt="item.name" class="w-full h-full object-cover" />
                <span v-else class="text-gray-300 text-4xl">🍽</span>
              </div>
              <div class="p-3 flex-1 flex flex-col">
                <h3 class="font-semibold text-gray-800 text-sm truncate">{{ item.name }}</h3>
                <p v-if="item.description" class="text-gray-400 text-xs truncate mt-0.5">{{ item.description }}</p>
                <span class="text-orange-600 font-bold text-sm mt-1 block">Rp{{ formatPrice(item.price) }}</span>
                <div class="mt-auto pt-2">
                  <button v-if="!getCartItem(item.id)"
                    @click="addItem(item)"
                    class="w-full py-1.5 rounded-lg bg-primary-500 text-white text-xs font-semibold hover:bg-primary-600 transition-colors">
                    + Tambah
                  </button>
                  <div v-else class="flex items-center justify-between bg-gray-50 rounded-lg border border-gray-200 px-2 py-1">
                    <button @click="updateQuantity(item.id, getCartItem(item.id)!.quantity - 1)"
                      class="w-6 h-6 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 hover:bg-gray-300 transition-colors text-xs font-bold">
                      -
                    </button>
                    <span class="font-semibold text-sm text-gray-800">{{ getCartItem(item.id)!.quantity }}</span>
                    <button @click="updateQuantity(item.id, getCartItem(item.id)!.quantity + 1)"
                      class="w-6 h-6 rounded-full bg-primary-500 flex items-center justify-center text-white hover:bg-primary-600 transition-colors text-xs font-bold">
                      +
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>

      <aside class="hidden lg:block w-80 xl:w-96 flex-shrink-0 border-l border-gray-200 bg-white sticky top-0 h-screen">
        <CartPanel />
      </aside>
    </div>

    <Transition name="slide-up">
      <div v-if="isCartOpen" class="fixed inset-0 z-50 lg:hidden">
        <div class="absolute inset-0 bg-black/40" @click="isCartOpen = false" />
        <div class="absolute bottom-0 left-0 right-0 bg-white rounded-t-2xl max-h-[80vh] flex flex-col overflow-hidden">
          <CartPanel show-close @close="isCartOpen = false" />
        </div>
      </div>
    </Transition>

    <Transition name="bar-up">
      <div v-if="itemCount > 0"
        @click="isCartOpen = true"
        class="lg:hidden fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 shadow-[0_-4px_12px_rgba(0,0,0,0.08)] flex items-center px-3 sm:px-4 md:px-6 py-2.5 sm:py-3 gap-2 sm:gap-3 z-40 cursor-pointer active:bg-gray-50 transition-colors">
        <div class="flex items-center gap-1.5 sm:gap-2 flex-shrink-0 min-w-0">
          <div :class="['w-8 h-8 sm:w-9 sm:h-9 rounded-full flex items-center justify-center transition-colors shrink-0',
            isFabPulsing ? 'bg-primary-600 scale-110' : 'bg-primary-500']">
            <svg class="w-4 h-4 sm:w-5 sm:h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z" />
            </svg>
          </div>
          <span class="font-bold text-xs sm:text-sm text-gray-800 truncate">{{ itemCount }} item</span>
        </div>
        <span class="ml-auto font-bold text-sm sm:text-base text-primary-600 truncate">Rp{{ formatPrice(total) }}</span>
        <button class="bg-primary-500 text-white px-4 sm:px-5 py-1.5 sm:py-2 rounded-xl text-xs sm:text-sm font-bold hover:bg-primary-600 transition-colors active:bg-primary-700 shrink-0">
          Pesan
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useCart } from '~/composables/useCart'
import type { CartItem } from '~/composables/useCart'
import type { Promo } from '~/components/PromoCarousel.vue'

interface Category { id: string; name: string }
interface MenuItem { id: string; category_id: string; name: string; description: string; price: number; image_url?: string }

const promos = ref<Promo[]>([])

const route = useRoute()
const chainId = (route.query.chain_id as string) || ''
const categories = ref<Category[]>([])
const items = ref<MenuItem[]>([])
const isCartOpen = ref(false)
const isFabPulsing = ref(false)

const { items: cartItems, itemCount, total, addItem, updateQuantity } = useCart()

watch(itemCount, (newVal, oldVal) => {
  if (newVal > oldVal) {
    isFabPulsing.value = true
    setTimeout(() => isFabPulsing.value = false, 400)
  }
})

function getCartItem(itemId: string): CartItem | undefined {
  return cartItems.value.find(i => i.id === itemId)
}

async function fetchData() {
  if (!chainId) return
  const [catRes, itemRes, carRes] = await Promise.all([
    $fetch(`/api/categories/list?chain_id=${chainId}`),
    $fetch(`/api/menu-items/list?chain_id=${chainId}`),
    $fetch(`/api/carousel/list?chain_id=${chainId}`).catch(() => ({ data: [] })),
  ])
  categories.value = (catRes as any).data || []
  items.value = (itemRes as any).data || []
  const slides = (carRes as any).data || (carRes as any) || []
  promos.value = (Array.isArray(slides) ? slides : []).map((s: any) => ({
    title: s.title,
    description: s.description || '',
    image: s.image_url || undefined,
    bgColor: s.bg_color ? `bg-gradient-to-r ${s.bg_color}` : undefined,
  }))
}

function getItemsByCategory(categoryId: string) {
  return items.value.filter(i => i.category_id === categoryId)
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

fetchData()
</script>

<style scoped>
@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.15); }
}

.animate-pulse {
  animation: pulse 0.4s ease-in-out;
}

.bar-up-enter-active {
  transition: transform 0.25s ease-out;
}
.bar-up-leave-active {
  transition: transform 0.2s ease-in;
}
.bar-up-enter-from,
.bar-up-leave-to {
  transform: translateY(100%);
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.3s ease-out;
}
.slide-up-enter-from {
  opacity: 0;
}
.slide-up-enter-from > div:last-child {
  transform: translateY(100%);
}
.slide-up-leave-to {
  opacity: 0;
}
.slide-up-leave-to > div:last-child {
  transform: translateY(100%);
}


</style>
