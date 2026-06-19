<template>
  <div class="min-h-screen bg-green-50 flex flex-col">
    <main class="flex-1 max-w-lg mx-auto w-full p-6 flex flex-col">
      <div class="flex-1 flex flex-col items-center justify-center text-center pt-12">
        <div class="w-16 h-16 rounded-full bg-green-500 flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold text-gray-800 mb-1">Pesanan Berhasil!</h1>
        <p class="text-gray-500 text-sm mb-6">Pesanan Anda telah diterima</p>
      </div>

      <div v-if="order" class="bg-white rounded-2xl shadow-sm p-6 space-y-4">
        <div class="flex justify-between items-center">
          <span class="text-gray-500 text-sm">No. Order</span>
          <span class="text-gray-800 font-mono text-sm font-bold">{{ order.id }}</span>
        </div>
        <div v-if="order.table_id" class="flex justify-between items-center">
          <span class="text-gray-500 text-sm">Meja</span>
          <span class="text-gray-800 font-semibold text-sm">{{ order.table_id }}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-gray-500 text-sm">Tipe</span>
          <span class="text-gray-800 font-semibold text-sm">{{ order.order_type === 'dine-in' ? 'Makan di Tempat' : 'Bungkus' }}</span>
        </div>
        <div class="flex justify-between items-center">
          <span class="text-gray-500 text-sm">Waktu</span>
          <span class="text-gray-800 font-semibold text-sm">{{ formatTime(order.created_at) }}</span>
        </div>

        <hr class="border-gray-100" />

        <p class="font-semibold text-gray-800 text-sm">Pesanan</p>
        <div v-for="item in order.items" :key="item.menu_item_id" class="flex justify-between items-center text-sm">
          <span class="text-gray-600">{{ item.quantity }}x {{ item.name || 'Item' }}</span>
          <span class="text-gray-800 font-medium">Rp{{ formatPrice(item.subtotal) }}</span>
        </div>

        <hr class="border-gray-100" />

        <div class="flex justify-between items-center text-lg font-bold">
          <span class="text-gray-800">Total</span>
          <span class="text-primary-600">Rp{{ formatPrice(order.total_amount) }}</span>
        </div>
      </div>

      <div v-else class="bg-white rounded-2xl shadow-sm p-6 text-center">
        <p class="text-gray-500">Data pesanan tidak ditemukan</p>
      </div>

      <button @click="backToMenu"
        class="mt-6 w-full py-3 bg-primary-500 text-white rounded-xl font-bold text-sm hover:bg-primary-600 transition-colors active:bg-primary-700">
        Kembali ke Menu
      </button>
    </main>
  </div>
</template>

<script setup lang="ts">
import { unref } from 'vue'

interface OrderItem {
  menu_item_id: string
  name?: string
  quantity: number
  unit_price: number
  subtotal: number
}

interface Order {
  id: string
  chain_id?: string
  branch_id?: string
  table_id?: string
  customer_name?: string
  order_type: string
  status?: string
  total_amount: number
  items: OrderItem[]
  created_at: string
}

const router = useRouter()

const order = ref<Order | null>(null)

onMounted(() => loadOrder())

function loadOrder() {
  const history = window.history.state
  if (history?.order) {
    order.value = history.order as Order
    sessionStorage.setItem('lastOrder', JSON.stringify(history.order))
  } else {
    const stored = sessionStorage.getItem('lastOrder')
    if (stored) {
      try {
        order.value = JSON.parse(stored) as Order
      } catch {
        order.value = null
      }
    }
  }
}

function formatPrice(price: number) {
  return unref(price).toLocaleString('id-ID')
}

function formatTime(iso: string) {
  const d = new Date(iso)
  return d.toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })
}

function backToMenu() {
  const o = order.value
  if (!o?.chain_id) { router.push('/menu'); return }
  const params = new URLSearchParams({ chain_id: o.chain_id })
  if (o.branch_id) params.set('branch_id', o.branch_id)
  if (o.table_id) params.set('table_id', o.table_id)
  router.push(`/menu?${params}`)
}

</script>
