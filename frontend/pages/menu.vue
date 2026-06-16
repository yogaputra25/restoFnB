<template>
  <div class="min-h-screen bg-gray-50">
    <div class="bg-gradient-to-r from-orange-500 to-red-500 text-white px-6 py-8 text-center">
      <p class="text-lg font-bold">Special Promo</p>
      <p class="text-sm opacity-90">Diskon 20% semua menu hari ini!</p>
    </div>

    <main class="max-w-5xl mx-auto p-4 space-y-6">
      <section v-for="cat in categories" :key="cat.id">
        <h2 class="text-lg font-bold text-gray-800 border-b border-gray-300 pb-1 mb-3">{{ cat.name }}</h2>
        <div class="flex gap-4 overflow-x-auto pb-2">
          <div v-for="item in getItemsByCategory(cat.id)" :key="item.id"
            class="min-w-[160px] max-w-[160px] bg-white rounded-xl shadow flex-shrink-0 overflow-hidden">
            <div class="h-32 bg-gray-200 flex items-center justify-center">
              <img v-if="item.image_url" :src="item.image_url" :alt="item.name" class="w-full h-full object-cover" />
              <span v-else class="text-gray-400 text-3xl">🍽</span>
            </div>
            <div class="p-3">
              <h3 class="font-semibold text-gray-800 text-sm truncate">{{ item.name }}</h3>
              <p v-if="item.description" class="text-gray-400 text-xs truncate mt-0.5">{{ item.description }}</p>
              <span class="text-orange-600 font-bold text-sm mt-1 block">Rp{{ formatPrice(item.price) }}</span>
            </div>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
interface Category { id: string; name: string }
interface MenuItem { id: string; category_id: string; name: string; description: string; price: number; image_url?: string }

const route = useRoute()
const chainId = (route.query.chain_id as string) || ''
const categories = ref<Category[]>([])
const items = ref<MenuItem[]>([])

async function fetchData() {
  if (!chainId) return
  const catRes = await $fetch(`/api/categories/list?chain_id=${chainId}`)
  categories.value = (catRes as any).data || []
  const itemRes = await $fetch(`/api/menu-items/list?chain_id=${chainId}`)
  items.value = (itemRes as any).data || []
}

function getItemsByCategory(categoryId: string) {
  return items.value.filter(i => i.category_id === categoryId)
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

fetchData()
</script>
