<template>
  <div>
    <PageTitle title="Orders" description="Manage incoming orders" />

    <div class="flex gap-6 h-[calc(100vh-200px)]">
      <!-- Left: Orders List -->
      <div class="w-[55%] flex flex-col">
        <!-- Toolbar -->
        <div class="flex items-center gap-3 mb-4 pb-3 border-b border-[var(--color-border)]">
          <div class="relative flex-1 max-w-xs">
            <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-text-secondary)]" />
            <input v-model="searchQuery" type="text" placeholder="Search orders..." class="w-full pl-9 pr-4 py-2 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]" />
          </div>
          <select v-model="statusFilter" class="px-3 py-2 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none">
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="confirmed">Confirmed</option>
            <option value="preparing">Preparing</option>
            <option value="ready">Ready</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex-1 space-y-3 overflow-y-auto">
          <div v-for="i in 5" :key="i" class="h-24 rounded-[var(--radius-card)] bg-gradient-to-r from-[var(--color-surface-secondary)] via-white/40 to-[var(--color-surface-secondary)] bg-[length:200%_100%] animate-shimmer" />
        </div>

        <!-- Empty -->
        <AppEmptyState
          v-else-if="filteredOrders.length === 0"
          title="No Orders"
          description="Orders will appear here when customers place them."
          class="flex-1"
        />

        <!-- Order cards -->
        <div v-else class="flex-1 space-y-3 overflow-y-auto pr-2">
          <div
            v-for="order in filteredOrders"
            :key="order.id"
            class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-4 transition-all duration-[var(--transition-fast)] cursor-pointer border-l-4 hover:shadow-[var(--shadow-md)]"
            :class="[selectedOrder?.id === order.id ? 'shadow-[var(--shadow-md)]' : '', statusBorder(order.status)]"
            @click="selectOrder(order)"
          >
            <div class="flex items-start justify-between gap-3">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="font-mono text-xs font-bold text-[var(--color-text-secondary)]">#{{ order.id?.slice(0, 8) }}</span>
                  <StatusBadge :status="order.status" />
                  <span v-if="order.payment_status" class="text-[10px] px-1.5 py-0.5 rounded bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)]">{{ order.payment_status }}</span>
                </div>
                <p class="text-sm font-medium text-[var(--color-text-primary)] truncate">{{ order.customer_name || 'Guest' }}</p>
                <p class="text-xs text-[var(--color-text-secondary)] mt-0.5">
                  {{ getTimeAgo(order.created_at) }} · {{ order.items?.length || 0 }} items
                  <span v-if="order.branch_id" class="ml-1.5 px-1 py-0.5 rounded bg-[var(--color-surface-secondary)] text-[10px]"><Building2 class="w-2.5 h-2.5 inline mr-0.5" />{{ getBranchName(order.branch_id) }}</span>
                </p>
                <p v-if="order.notes" class="text-xs text-[var(--color-text-secondary)] mt-0.5 truncate max-w-[200px]" :title="order.notes">{{ order.notes }}</p>
              </div>
              <div class="text-right shrink-0">
                <p class="font-bold text-sm text-[var(--color-primary)]">Rp{{ formatPrice(order.total_amount) }}</p>
                <p class="text-[10px] text-[var(--color-text-secondary)]">{{ order.order_type }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Order Detail -->
      <div class="w-[45%]">
        <div v-if="!selectedOrder" class="h-full flex items-center justify-center">
          <div class="text-center">
            <ClipboardList class="w-16 h-16 mx-auto text-[var(--color-border)]" stroke-width="1.5" />
            <p class="mt-3 text-sm text-[var(--color-text-secondary)]">Select an order to view details</p>
          </div>
        </div>

        <div v-else class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-6 h-full overflow-y-auto space-y-4">
          <!-- Order header -->
          <div class="flex items-start justify-between">
            <div>
              <h3 class="font-heading text-card-title text-[var(--color-text-primary)]">Order #{{ selectedOrder.id?.slice(0, 8) }}</h3>
              <p class="text-xs text-[var(--color-text-secondary)]">{{ selectedOrder.customer_name || 'Guest' }} · {{ selectedOrder.order_type }} · {{ formatDate(selectedOrder.created_at) }}</p>
              <p v-if="selectedOrder.branch_id" class="text-xs text-[var(--color-text-secondary)] mt-1">
                <Building2 class="w-3 h-3 inline mr-0.5" />{{ getBranchName(selectedOrder.branch_id) }}
              </p>
            </div>
            <StatusBadge :status="selectedOrder.status" />
          </div>

          <hr class="border-[var(--color-border)]" />

          <!-- Items -->
          <div>
            <h4 class="text-sm font-semibold text-[var(--color-text-primary)] mb-2">Items</h4>
            <div class="space-y-2">
              <div v-for="(item, i) in selectedOrder.items" :key="i" class="flex items-center justify-between text-sm">
                <span class="text-[var(--color-text-primary)]">{{ item.quantity }}x {{ item.name || `Item ${i + 1}` }}<span v-if="item.variant_name" class="text-[var(--color-text-secondary)]"> ({{ item.variant_name }})</span></span>
                <span class="font-medium text-[var(--color-text-primary)]">Rp{{ formatPrice(item.subtotal || item.unit_price * item.quantity) }}</span>
              </div>
            </div>
          </div>

          <hr class="border-[var(--color-border)]" />

          <!-- Totals -->
          <div class="space-y-1">
            <div class="flex justify-between text-sm">
              <span class="text-[var(--color-text-secondary)]">Subtotal</span>
              <span class="text-[var(--color-text-primary)]">Rp{{ formatPrice(selectedOrder.total_amount) }}</span>
            </div>
            <div class="flex justify-between font-bold text-lg">
              <span class="text-[var(--color-text-primary)]">Total</span>
              <span class="text-[var(--color-primary)]">Rp{{ formatPrice(selectedOrder.total_amount) }}</span>
            </div>
          </div>

          <!-- Notes -->
          <div v-if="selectedOrder.notes">
            <h4 class="text-sm font-semibold text-[var(--color-text-primary)] mb-1">Notes</h4>
            <p class="text-sm text-[var(--color-text-secondary)] bg-[var(--color-surface-secondary)] p-3 rounded-[var(--radius-input)]">{{ selectedOrder.notes }}</p>
          </div>

          <!-- Workflow buttons -->
          <div class="flex flex-wrap gap-2 pt-2">
            <button
              v-if="selectedOrder.status === 'pending'"
              class="flex-1 min-w-[120px] py-2.5 rounded-[var(--radius-button)] bg-[var(--color-primary)] text-white text-sm font-semibold hover:bg-[#a84e31] transition-all hover:scale-[1.03] active:scale-[0.98]"
              @click="updateStatus('confirmed')"
            >
              Confirm
            </button>
            <button
              v-if="selectedOrder.status === 'confirmed'"
              class="flex-1 min-w-[120px] py-2.5 rounded-[var(--radius-button)] bg-[var(--color-secondary)] text-white text-sm font-semibold hover:bg-[#b86205] transition-all hover:scale-[1.03] active:scale-[0.98]"
              @click="updateStatus('preparing')"
            >
              Prepare
            </button>
            <button
              v-if="selectedOrder.status === 'preparing'"
              class="flex-1 min-w-[120px] py-2.5 rounded-[var(--radius-button)] bg-[var(--color-accent)] text-white text-sm font-semibold hover:bg-[#d69e00] transition-all hover:scale-[1.03] active:scale-[0.98]"
              @click="updateStatus('ready')"
            >
              Ready
            </button>
            <button
              v-if="selectedOrder.status === 'ready' && selectedOrder.payment_status !== 'paid'"
              class="flex-1 min-w-[120px] py-2.5 rounded-[var(--radius-button)] bg-[var(--color-primary)] text-white text-sm font-semibold hover:bg-[#a84e31] transition-all hover:scale-[1.03] active:scale-[0.98]"
              @click="processPayment"
            >
              Pay
            </button>
            <button
              v-if="selectedOrder.status === 'ready' && selectedOrder.payment_status === 'paid'"
              class="flex-1 min-w-[120px] py-2.5 rounded-[var(--radius-button)] bg-[var(--color-success)] text-white text-sm font-semibold hover:bg-[#256d29] transition-all hover:scale-[1.03] active:scale-[0.98]"
              @click="updateStatus('completed')"
            >
              Complete
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Search, ClipboardList, Building2 } from 'lucide-vue-next'
definePageMeta({ layout: 'admin', middleware: 'auth' })

interface Order {
  id: string
  branch_id?: string
  status: string
  order_type: string
  total_amount: number
  payment_status?: string
  customer_name?: string
  created_at?: string
  notes?: string
  items?: Array<{
    menu_item_id?: string
    name?: string
    quantity: number
    unit_price: number
    subtotal?: number
  }>
}

const { $api } = useNuxtApp()
const toast = useToast()
const auth = useAuthStore()
const { selectedBranchId, getBranchName, setBranch } = useBranchSelector()

// Auto-select branch from logged-in user
if (auth.user?.branchId) {
  setBranch(auth.user.branchId)
}

const loading = ref(true)
const orders = ref<Order[]>([])
const searchQuery = ref('')
const statusFilter = ref('')
const selectedOrder = ref<Order | null>(null)

const filteredOrders = computed(() => {
  let result = orders.value
  if (statusFilter.value) {
    result = result.filter(o => o.status === statusFilter.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(o =>
      (o.id || '').toLowerCase().includes(q) ||
      (o.customer_name || '').toLowerCase().includes(q)
    )
  }
  return result
})

function formatPrice(price: number) { return price.toLocaleString('id-ID') }

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function getTimeAgo(dateStr: string) {
  if (!dateStr) return ''
  const diff = Date.now() - new Date(dateStr).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return 'Just now'
  if (mins < 60) return `${mins}m ago`
  const hrs = Math.floor(mins / 60)
  return `${hrs}h ${mins % 60}m ago`
}

const statusColors: Record<string, string> = {
  pending: 'border-l-yellow-500',
  confirmed: 'border-l-blue-500',
  preparing: 'border-l-purple-500',
  ready: 'border-l-green-500',
  completed: 'border-l-gray-400',
  cancelled: 'border-l-red-500',
}

function statusBorder(status: string) {
  return statusColors[status] || 'border-l-[var(--color-border)]'
}

function selectOrder(order: Order) { selectedOrder.value = order }

async function fetchOrders() {
  try {
    const endpoint = selectedBranchId.value
      ? `/api/orders/by-branch?branch_id=${selectedBranchId.value}&page=1&limit=50`
      : '/api/orders/by-chain?page=1&limit=50'
    const res = await $api(endpoint)
    const data = (res as any).data || { items: [] }
    orders.value = (data.items || []).sort((a: Order, b: Order) => {
      if (!a.created_at) return 1
      if (!b.created_at) return -1
      return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    })
    // Auto-select first pending order
    if (!selectedOrder.value && orders.value.length > 0) {
      selectedOrder.value = orders.value[0]
    }
  } catch (e: any) { toast.show(e?.data?.message || 'Failed to load orders') }
}

async function updateStatus(status: string) {
  if (!selectedOrder.value) return
  try {
    await $api(`/api/orders/status?id=${selectedOrder.value.id}`, { method: 'POST', body: { status } })
    toast.show(`Order ${status}`)
    await fetchOrders()
    // Re-select to keep detail open
    const updated = orders.value.find(o => o.id === selectedOrder.value?.id)
    if (updated) selectedOrder.value = updated
  } catch (e: any) { toast.show(e?.data?.message || 'Failed to update order') }
}

async function processPayment() {
  if (!selectedOrder.value) return
  try {
    await $api(`/api/orders/pay?id=${selectedOrder.value.id}`, { method: 'POST', body: { method: 'cash' } })
    toast.show('Payment recorded')
    await fetchOrders()
    const updated = orders.value.find(o => o.id === selectedOrder.value?.id)
    if (updated) selectedOrder.value = updated
  } catch (e: any) { toast.show(e?.data?.message || 'Payment failed') }
}

onMounted(async () => {
  loading.value = true
  await fetchOrders()
  loading.value = false
})

watch(selectedBranchId, () => {
  fetchOrders()
})
</script>

<style scoped>
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}
.animate-shimmer { animation: shimmer 1.5s ease-in-out infinite; }
</style>
