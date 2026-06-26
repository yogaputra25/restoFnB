<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-white">Sales Report</h1>
      <button @click="downloadCSV"
        class="px-4 py-2 bg-primary-500 text-white rounded-lg text-sm hover:bg-primary-600 transition">
        Download CSV
      </button>
    </div>

    <div class="flex flex-wrap gap-4 mb-6 items-end">
      <div>
        <label class="block text-xs text-dark-400 mb-1">Period</label>
        <select v-model="period" @change="fetchReport"
          class="bg-dark-700 text-white border border-dark-600 rounded-lg px-3 py-2 text-sm">
          <option value="daily">Daily</option>
          <option value="monthly">Monthly</option>
          <option value="yearly">Yearly</option>
        </select>
      </div>
      <div>
        <label class="block text-xs text-dark-400 mb-1">From</label>
        <input type="date" v-model="fromDate" @change="fetchReport"
          class="bg-dark-700 text-white border border-dark-600 rounded-lg px-3 py-2 text-sm" />
      </div>
      <div>
        <label class="block text-xs text-dark-400 mb-1">To</label>
        <input type="date" v-model="toDate" @change="fetchReport"
          class="bg-dark-700 text-white border border-dark-600 rounded-lg px-3 py-2 text-sm" />
      </div>
    </div>

    <div v-if="loading" class="text-dark-400 text-center py-16">Loading...</div>

    <template v-if="report && !loading">
      <div class="grid grid-cols-3 gap-4 mb-6">
        <div class="bg-dark-700 rounded-lg p-4">
          <p class="text-dark-400 text-xs uppercase tracking-wide">Total Revenue</p>
          <p class="text-2xl font-bold text-primary-500 mt-1">Rp{{ formatNumber(report.summary.total_revenue) }}</p>
        </div>
        <div class="bg-dark-700 rounded-lg p-4">
          <p class="text-dark-400 text-xs uppercase tracking-wide">Total Orders</p>
          <p class="text-2xl font-bold text-white mt-1">{{ report.summary.total_orders }}</p>
        </div>
        <div class="bg-dark-700 rounded-lg p-4">
          <p class="text-dark-400 text-xs uppercase tracking-wide">Total Items Sold</p>
          <p class="text-2xl font-bold text-white mt-1">{{ report.summary.total_items }}</p>
        </div>
      </div>

      <div class="bg-dark-700 rounded-lg p-4 mb-6">
        <h2 class="text-white font-semibold mb-4">Revenue Trend</h2>
        <Line :data="lineChartData" :options="chartOptions" />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <div v-for="cat in report.categories" :key="cat.category_id"
          class="bg-dark-700 rounded-lg p-4">
          <h3 class="text-white font-semibold mb-3">{{ cat.category_name }}</h3>
          <Bar :data="barChartData(cat)" :options="barOptions" />
        </div>
      </div>

      <div class="bg-dark-700 rounded-lg p-4">
        <h2 class="text-white font-semibold mb-4">Item Details</h2>
        <div class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="text-dark-400 text-left border-b border-dark-600">
                <th class="pb-2 pr-4">Category</th>
                <th class="pb-2 pr-4">Menu Item</th>
                <th class="pb-2 pr-4 text-right">Quantity</th>
                <th class="pb-2 pr-4 text-right">Revenue</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="row in tableRows" :key="row.menu_item_id"
                class="border-b border-dark-600 last:border-0">
                <td class="py-2 pr-4 text-dark-400">{{ row.category }}</td>
                <td class="py-2 pr-4 text-white">{{ row.name }}</td>
                <td class="py-2 pr-4 text-white text-right">{{ row.quantity }}</td>
                <td class="py-2 pr-4 text-primary-500 text-right">Rp{{ formatNumber(row.revenue) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { Line, Bar } from 'vue-chartjs'

definePageMeta({ layout: 'admin', middleware: 'auth' })

const { $api } = useNuxtApp()

interface SalesReportItem {
  menu_item_id: string
  name: string
  quantity: number
  revenue: number
}

interface SalesReportCategory {
  category_id: string
  category_name: string
  items: SalesReportItem[]
}

interface SalesReportDaily {
  date: string
  revenue: number
  orders: number
}

interface SalesReportResponse {
  period: string
  from: string
  to: string
  summary: { total_revenue: number; total_orders: number; total_items: number }
  categories: SalesReportCategory[]
  daily_breakdown: SalesReportDaily[]
}

const period = ref('daily')
const now = new Date()
const fromDate = ref(new Date(now.getFullYear(), now.getMonth(), 1).toISOString().slice(0, 10))
const toDate = ref(now.toISOString().slice(0, 10))
const report = ref<SalesReportResponse | null>(null)
const loading = ref(false)

async function fetchReport() {
  loading.value = true
  try {
    const res: any = await $api(`/api/reports/sales?period=${period.value}&from=${fromDate.value}&to=${toDate.value}`)
    report.value = res.data
  } catch (e: any) {
    report.value = null
  }
  loading.value = false
}

function downloadCSV() {
  window.open(`/api/reports/sales?period=${period.value}&from=${fromDate.value}&to=${toDate.value}&format=csv`, '_blank')
}

const chartOptions = {
  responsive: true,
  plugins: {
    legend: { display: false },
  },
  scales: {
    x: {
      ticks: { color: '#888', maxTicksLimit: 12 },
      grid: { color: '#333' },
    },
    y: {
      ticks: { color: '#888' },
      grid: { color: '#333' },
    },
  },
}

const barOptions = {
  responsive: true,
  plugins: {
    legend: { display: false },
  },
  scales: {
    x: {
      ticks: { color: '#888' },
      grid: { display: false },
    },
    y: {
      ticks: { color: '#888' },
      grid: { color: '#333' },
    },
  },
}

const lineChartData = computed(() => {
  if (!report.value) return { labels: [], datasets: [] }
  return {
    labels: report.value.daily_breakdown.map(d => d.date),
    datasets: [
      {
        label: 'Revenue',
        data: report.value.daily_breakdown.map(d => d.revenue),
        borderColor: '#f59e0b',
        backgroundColor: 'rgba(245, 158, 11, 0.1)',
        fill: true,
        tension: 0.3,
      },
    ],
  }
})

function barChartData(cat: SalesReportCategory) {
  return {
    labels: cat.items.map(i => i.name),
    datasets: [
      {
        label: 'Sold',
        data: cat.items.map(i => i.quantity),
        backgroundColor: '#3b82f6',
        borderRadius: 4,
      },
    ],
  }
}

interface TableRow {
  menu_item_id: string
  category: string
  name: string
  quantity: number
  revenue: number
}

const tableRows = computed(() => {
  if (!report.value) return []
  const rows: TableRow[] = []
  for (const cat of report.value.categories) {
    for (const item of cat.items) {
      rows.push({
        menu_item_id: item.menu_item_id,
        category: cat.category_name,
        name: item.name,
        quantity: item.quantity,
        revenue: item.revenue,
      })
    }
  }
  return rows
})

function formatNumber(n: number) {
  return n.toLocaleString('id-ID')
}

fetchReport()
</script>
