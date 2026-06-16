<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-white mb-6">Orders</h1>

    <div class="flex gap-2 mb-4">
      <button
        v-for="s in statusFilters"
        :key="s.value"
        @click="filterStatus = s.value"
        class="px-3 py-1.5 rounded text-sm transition"
        :class="filterStatus === s.value ? 'bg-primary-500 text-white' : 'bg-dark-700 text-dark-300 hover:bg-dark-600'"
      >{{ s.label }}</button>
    </div>

    <div class="grid gap-3">
      <div v-for="order in filteredOrders" :key="order.id"
        class="bg-dark-700 rounded-lg p-4 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <span class="text-white font-semibold">#{{ order.id?.slice(0, 8) }}</span>
          <span class="px-2 py-0.5 rounded text-xs text-white"
            :class="statusClass(order.status)">{{ order.status }}</span>
          <span class="text-dark-400 text-sm">{{ order.order_type }}</span>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-primary-500 font-bold">Rp{{ formatPrice(order.total_amount) }}</span>
          <div class="flex gap-1">
            <button v-if="order.status === 'pending'" @click="updateStatus(order.id, 'confirmed')"
              class="px-3 py-1 bg-green-600 text-white rounded text-xs">Confirm</button>
            <button v-if="order.status === 'confirmed'" @click="updateStatus(order.id, 'preparing')"
              class="px-3 py-1 bg-blue-600 text-white rounded text-xs">Prepare</button>
            <button v-if="order.status === 'preparing'" @click="updateStatus(order.id, 'ready')"
              class="px-3 py-1 bg-yellow-600 text-white rounded text-xs">Ready</button>
            <button v-if="order.status === 'ready' && order.payment_status !== 'paid'" @click="processPayment(order.id)"
              class="px-3 py-1 bg-primary-500 text-white rounded text-xs">Pay</button>
            <button v-if="order.status === 'ready' && order.payment_status === 'paid'" @click="updateStatus(order.id, 'completed')"
              class="px-3 py-1 bg-green-600 text-white rounded text-xs">Complete</button>
          </div>
        </div>
      </div>
      <p v-if="filteredOrders.length === 0" class="text-dark-400 text-center py-8">No orders</p>
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

const filterStatus = ref('')
const orders = ref<Order[]>([])

const statusFilters = [
  { label: 'All', value: '' },
  { label: 'Pending', value: 'pending' },
  { label: 'Confirmed', value: 'confirmed' },
  { label: 'Preparing', value: 'preparing' },
  { label: 'Ready', value: 'ready' },
  { label: 'Completed', value: 'completed' },
]

const filteredOrders = computed(() =>
  filterStatus.value ? orders.value.filter(o => o.status === filterStatus.value) : orders.value
)

async function fetchOrders() {
  try {
    const res = await $api('/api/orders/by-branch')
    orders.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load orders')
  }
}

async function updateStatus(id: string, status: string) {
  try {
    await $api(`/api/orders/status?id=${id}`, { method: 'POST', body: { status } })
    await fetchOrders()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to update order')
  }
}

async function processPayment(id: string) {
  try {
    await $api(`/api/orders/pay?id=${id}`, { method: 'POST', body: { method: 'cash' } })
    await fetchOrders()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Payment failed')
  }
}

function statusClass(status: string) {
  const map: Record<string, string> = {
    pending: 'bg-yellow-600',
    confirmed: 'bg-blue-600',
    preparing: 'bg-purple-600',
    ready: 'bg-green-600',
    completed: 'bg-dark-500 text-dark-200',
  }
  return map[status] || 'bg-dark-500'
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

fetchOrders()
</script>
