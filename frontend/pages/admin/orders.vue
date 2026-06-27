<template>
  <div class="p-6">
    <h1 class="font-heading text-[var(--font-size-page-title)] text-[var(--color-text-primary)] mb-6">All Orders</h1>

    <div v-if="loading" class="grid gap-3">
      <AppSkeleton v-for="i in 5" :key="i" variant="card" />
    </div>

    <AppEmptyState
      v-else-if="orders.length === 0"
      title="No Orders"
      description="Orders will appear here once customers start ordering."
    />

    <div v-else class="grid gap-3">
      <div v-for="order in orders" :key="order.id"
        class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-4 flex items-center justify-between transition-all hover:-translate-y-1 hover:shadow-[var(--shadow-md)]">
        <div class="flex items-center gap-3">
          <span class="text-[var(--color-text-primary)] font-semibold">#{{ order.id?.slice(0, 8) }}</span>
          <StatusBadge :status="order.status" variant="order" />
          <span class="text-[var(--color-text-secondary)] text-sm">{{ order.order_type }}</span>
          <span v-if="order.branch_id" class="text-[10px] px-1.5 py-0.5 rounded bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)]"><Building2 class="w-2.5 h-2.5 inline mr-0.5" />{{ getBranchName(order.branch_id) }}</span>
          <span v-if="order.notes" class="text-[10px] px-1.5 py-0.5 rounded bg-[var(--color-accent)]/10 text-[var(--color-accent)] max-w-[120px] truncate inline-block align-middle" :title="order.notes">{{ order.notes }}</span>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-[var(--color-text-secondary)] text-sm capitalize">{{ order.payment_status }}</span>
          <span class="text-[var(--color-primary)] font-bold">Rp{{ formatPrice(order.total_amount) }}</span>
        </div>
      </div>
    </div>
    <Paginator :current-page="page" :total-items="total" :limit="limit" @page-change="onPageChange" />
  </div>
</template>

<script setup lang="ts">
import { Building2 } from 'lucide-vue-next'
definePageMeta({ layout: 'admin', middleware: 'auth' })

const { $api } = useNuxtApp()
const toast = useToast()
const { selectedBranchId, getBranchName } = useBranchSelector()

interface Order {
  id: string
  status: string
  order_type: string
  total_amount: number
  payment_status: string
  branch_id?: string
  notes?: string
}

interface PaginatedData {
  items: Order[]
  total: number
  page: number
  limit: number
}

const orders = ref<Order[]>([])
const page = ref(1)
const total = ref(0)
const limit = 20
const loading = ref(false)

async function fetchOrders() {
  loading.value = true
  try {
    const endpoint = selectedBranchId.value
      ? `/api/orders/by-branch?branch_id=${selectedBranchId.value}&page=${page.value}&limit=${limit}`
      : `/api/orders/by-chain?page=${page.value}&limit=${limit}`
    const res = await $api(endpoint)
    const data = (res as any).data as PaginatedData
    orders.value = data.items || []
    total.value = data.total
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load orders')
  }
  loading.value = false
}

// Watch branch changes to refetch
watch(selectedBranchId, () => { page.value = 1; fetchOrders() })

function onPageChange(p: number) {
  page.value = p
  fetchOrders()
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

fetchOrders()
</script>
