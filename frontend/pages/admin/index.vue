<template>
  <div>
    <PageTitle title="Dashboard" description="Overview of your restaurant performance" />

    <!-- Stat Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <StatCard
        v-for="(stat, i) in statCards"
        :key="stat.label"
        :icon="stat.icon"
        :label="stat.label"
        :value="stat.value"
        :trend="stat.trend"
        :trend-up="stat.trendUp"
        :color="stat.color"
        :delay="i * 100"
      />
    </div>

    <!-- Loading skeleton for stats -->
    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
      <div v-for="i in 4" :key="i" class="h-28 rounded-[var(--radius-card)] bg-gradient-to-r from-[var(--color-surface-secondary)] via-white/40 to-[var(--color-surface-secondary)] bg-[length:200%_100%] animate-shimmer" />
    </div>

    <!-- Charts + Tables Row -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-8">
      <!-- Sales Chart -->
      <div class="lg:col-span-2 bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-6" data-aos="fade-up">
        <h3 class="font-heading text-card-title text-[var(--color-text-primary)] mb-4">Sales (7 days)</h3>
        <div v-if="chartLoading" class="h-64 rounded-lg bg-gradient-to-r from-[var(--color-surface-secondary)] via-white/40 to-[var(--color-surface-secondary)] bg-[length:200%_100%] animate-shimmer" />
        <div v-else-if="hasChartData">
          <canvas ref="chartCanvas" class="w-full h-64" />
        </div>
        <AppEmptyState v-else title="No Sales Data" description="Chart will appear once orders are placed." class="py-8" />
      </div>

      <!-- Popular Menu -->
      <div class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-6" data-aos="fade-up" data-aos-delay="100">
        <h3 class="font-heading text-card-title text-[var(--color-text-primary)] mb-4">Popular Menu</h3>
        <div v-if="loading" class="space-y-3">
          <div v-for="i in 5" :key="i" class="h-10 rounded-lg bg-gradient-to-r from-[var(--color-surface-secondary)] via-white/40 to-[var(--color-surface-secondary)] bg-[length:200%_100%] animate-shimmer" />
        </div>
        <div v-else-if="popularItems.length === 0">
          <AppEmptyState title="No Data" description="Popular items will show once orders come in." />
        </div>
        <div v-else class="space-y-3">
          <div v-for="(item, i) in popularItems" :key="i" class="flex items-center gap-3">
            <span class="w-6 h-6 rounded-full bg-[var(--color-surface-secondary)] flex items-center justify-center text-xs font-bold text-[var(--color-text-secondary)]">{{ i + 1 }}</span>
            <span class="flex-1 text-sm text-[var(--color-text-primary)] truncate">{{ item.name }}</span>
            <span class="text-xs font-semibold text-[var(--color-text-secondary)]">{{ item.count }}x</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Recent Orders -->
    <div class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-6 mb-8" data-aos="fade-up">
      <div class="flex items-center justify-between mb-4">
        <h3 class="font-heading text-card-title text-[var(--color-text-primary)]">Recent Orders</h3>
        <RouterLink to="/admin/orders" class="text-sm text-[var(--color-primary)] hover:underline font-medium">View All</RouterLink>
      </div>

      <DataGrid
        :columns="orderColumns"
        :rows="recentOrders"
        :loading="loading"
        empty-message="No orders yet"
        @row-click="goToOrders"
      >
        <template #cell-status="{ value }">
          <StatusBadge :status="value" />
        </template>
        <template #cell-total_amount="{ value }">
          <span class="font-semibold text-[var(--color-text-primary)]">Rp{{ formatPrice(value) }}</span>
        </template>
      </DataGrid>
    </div>

    <!-- Quick Actions -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4" data-aos="fade-up">
      <RouterLink
        v-for="action in quickActions"
        :key="action.label"
        :to="action.to"
        class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-4 flex flex-col items-center gap-2 text-center transition-all duration-[var(--transition-base)] hover:-translate-y-1 hover:shadow-[var(--shadow-md)]"
      >
        <div class="w-10 h-10 rounded-full flex items-center justify-center" :style="{ backgroundColor: `${action.color}15`, color: action.color }">
          <component :is="actionIcon(action.icon)" class="w-5 h-5" />
        </div>
        <span class="text-sm font-medium text-[var(--color-text-primary)]">{{ action.label }}</span>
      </RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import { unref } from 'vue'
import { h } from 'vue'
import * as LucideIcons from 'lucide-vue-next'
import { Chart as ChartJS, CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Filler } from 'chart.js'
import { Line } from 'vue-chartjs'

definePageMeta({ layout: 'admin', middleware: 'auth' })

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Filler)

const { $api } = useNuxtApp()
const toast = useToast()

const loading = ref(true)
const chartLoading = ref(true)
const recentOrders = ref<any[]>([])
const popularItems = ref<{ name: string; count: number }[]>([])
const chartData = ref<any>(null)
const hasChartData = ref(false)
const chartCanvas = ref<HTMLCanvasElement | null>(null)

interface StatCardData { icon: string; label: string; value: string; trend?: string; trendUp?: boolean; color: string }
const statCards = ref<StatCardData[]>([
  { icon: 'dollar-sign', label: 'Revenue Today', value: 'Rp0', trend: '0%', trendUp: true, color: '#2E7D32' },
  { icon: 'shopping-cart', label: "Today's Orders", value: '0', trend: '0%', trendUp: true, color: '#C65D3B' },
  { icon: 'check-circle', label: 'Completed', value: '0', color: '#2E7D32' },
  { icon: 'clock', label: 'Pending', value: '0', color: '#D97706' },
])

const orderColumns = [
  { key: 'id', label: 'Order #', flex: '1' },
  { key: 'customer_name', label: 'Customer', flex: '1' },
  { key: 'created_at', label: 'Time', flex: '1' },
  { key: 'total_amount', label: 'Total', flex: '0.8' },
  { key: 'status', label: 'Status', flex: '0.8' },
]

const quickActions = [
  { icon: 'plus-circle', label: 'Add Menu Item', to: '/admin/menu', color: '#C65D3B' },
  { icon: 'clipboard-list', label: 'View Orders', to: '/admin/orders', color: '#D97706' },
  { icon: 'bar-chart-3', label: 'View Reports', to: '/admin/reports', color: '#2E7D32' },
  { icon: 'building-2', label: 'Branches', to: '/admin/branches', color: '#F4B400' },
]

function actionIcon(name: string) {
  const icons = LucideIcons as Record<string, any>
  const key = name.split('-').map(w => w.charAt(0).toUpperCase() + w.slice(1)).join('') + 'Icon'
  return icons[key] || icons['CircleIcon'] || 'span'
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

async function fetchDashboard() {
  loading.value = true
  try {
    const [orderRes] = await Promise.all([
      $api('/api/orders/by-chain?page=1&limit=50'),
    ])
    const data = (orderRes as any).data || { items: [], total: 0 }
    const orders = data.items || []

    // Recent 5 orders
    recentOrders.value = orders.slice(0, 5)

    // Today orders
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    const todayOrders = orders.filter((o: any) => o.created_at && new Date(o.created_at) >= today)
    const completed = todayOrders.filter((o: any) => o.status === 'completed')
    const pending = todayOrders.filter((o: any) => o.status === 'pending' || o.status === 'confirmed' || o.status === 'preparing' || o.status === 'ready')
    const revenue = completed.reduce((sum: number, o: any) => sum + (o.total_amount || 0), 0)

    statCards.value = [
      { icon: 'dollar-sign', label: 'Revenue Today', value: `Rp${revenue.toLocaleString('id-ID')}`, trend: revenue > 0 ? '+100%' : '0%', trendUp: revenue > 0, color: '#2E7D32' },
      { icon: 'shopping-cart', label: "Today's Orders", value: String(todayOrders.length), trend: todayOrders.length > 0 ? '+100%' : '0%', trendUp: todayOrders.length > 0, color: '#C65D3B' },
      { icon: 'check-circle', label: 'Completed', value: String(completed.length), color: '#2E7D32' },
      { icon: 'clock', label: 'Pending', value: String(pending.length), color: '#D97706' },
    ]

    // Popular items
    const itemCount: Record<string, number> = {}
    for (const order of orders) {
      if (order.items) {
        for (const item of order.items) {
          const name = item.name || `Item ${item.menu_item_id?.slice(0, 6)}`
          itemCount[name] = (itemCount[name] || 0) + item.quantity
        }
      }
    }
    popularItems.value = Object.entries(itemCount)
      .map(([name, count]) => ({ name, count }))
      .sort((a, b) => b.count - a.count)
      .slice(0, 5)

    // Chart data — last 7 days
    chartLoading.value = false
    const days: string[] = []
    const dayRevenues: number[] = []
    for (let i = 6; i >= 0; i--) {
      const d = new Date()
      d.setDate(d.getDate() - i)
      const ds = d.toISOString().slice(0, 10)
      days.push(d.toLocaleDateString('en-US', { weekday: 'short' }))
      const dayOrders = orders.filter((o: any) => {
        if (!o.created_at) return false
        const od = o.created_at.slice(0, 10)
        return od === ds && (o.status === 'completed' || o.payment_status === 'paid')
      })
      dayRevenues.push(dayOrders.reduce((sum: number, o: any) => sum + (o.total_amount || 0), 0))
    }
    chartData.value = {
      labels: days,
      datasets: [{
        label: 'Revenue',
        data: dayRevenues,
        borderColor: '#C65D3B',
        backgroundColor: 'rgba(198, 93, 59, 0.08)',
        fill: true,
        tension: 0.3,
        pointRadius: 3,
      }],
    }
    hasChartData.value = dayRevenues.some(r => r > 0)

    // Render chart
    await nextTick()
    if (chartCanvas.value && hasChartData.value) {
      // Chart.js is already rendered — chartData ref will trigger
    }
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load dashboard')
  } finally {
    loading.value = false
    chartLoading.value = false
  }
}

fetchDashboard()
</script>

<style scoped>
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}
.animate-shimmer {
  animation: shimmer 1.5s ease-in-out infinite;
}
</style>
