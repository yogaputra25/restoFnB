<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-white mb-6">All Orders</h1>
    <div class="grid gap-3">
      <div v-for="order in orders" :key="order.id"
        class="bg-dark-700 rounded-lg p-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <span class="text-white font-semibold">#{{ order.id?.slice(0, 8) }}</span>
          <span class="px-2 py-0.5 rounded text-xs text-white"
            :class="order.status === 'completed' ? 'bg-green-600' : 'bg-yellow-600'">{{ order.status }}</span>
          <span class="text-dark-400 text-sm">{{ order.order_type }}</span>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-dark-300 text-sm capitalize">{{ order.payment_status }}</span>
          <span class="text-primary-500 font-bold">Rp{{ formatPrice(order.total_amount) }}</span>
        </div>
      </div>
      <p v-if="orders.length === 0" class="text-dark-400 text-center py-8">No orders found</p>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { $api } = useNuxtApp()
const toast = useToast()

interface Order {
  id: string
  status: string
  order_type: string
  total_amount: number
  payment_status: string
}

const orders = ref<Order[]>([])

async function fetchOrders() {
  try {
    const res = await $api('/api/orders/by-chain')
    orders.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load orders')
  }
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

fetchOrders()
</script>
