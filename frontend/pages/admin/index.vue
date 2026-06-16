<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-white mb-6">Overview</h1>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-dark-700 p-5 rounded-lg">
        <p class="text-3xl font-bold text-primary-500">{{ stats.menuItems }}</p>
        <p class="text-dark-300 text-sm mt-1">Menu Items</p>
      </div>
      <div class="bg-dark-700 p-5 rounded-lg">
        <p class="text-3xl font-bold text-primary-500">{{ stats.branches }}</p>
        <p class="text-dark-300 text-sm mt-1">Branches</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { $api } = useNuxtApp()
const toast = useToast()
const stats = ref({ menuItems: 0, branches: 0 })

async function fetchStats() {
  try {
    const [menuRes, branchRes] = await Promise.all([
      $api('/api/admin/menu-items/list'),
      $api('/api/admin/branches/list'),
    ])
    stats.value.menuItems = ((menuRes as any).data || []).length
    stats.value.branches = ((branchRes as any).data || []).length
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load stats')
  }
}

fetchStats()
</script>
