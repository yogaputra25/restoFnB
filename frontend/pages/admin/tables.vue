<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-white">Tables</h1>
      <div class="flex items-center gap-3">
        <select v-model="selectedBranchId" @change="fetchTables"
          class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none">
          <option value="" disabled>Select branch</option>
          <option v-for="b in branches" :key="b.id" :value="b.id">{{ b.name }}</option>
        </select>
        <button v-if="selectedBranchId" @click="showForm = true"
          class="px-4 py-2 bg-primary-500 text-white rounded hover:bg-primary-600 text-sm">
          Add Table
        </button>
      </div>
    </div>

    <div v-if="!selectedBranchId" class="text-dark-400 text-center py-8">Select a branch to view tables</div>

    <div v-else class="grid gap-3">
      <div v-for="t in tables" :key="t.id" class="bg-dark-700 rounded-lg p-4 flex items-center justify-between">
        <div>
          <p class="text-white font-medium">{{ t.label }}</p>
          <p class="text-dark-400 text-sm">Branch: {{ branchName(t.branch_id) }}</p>
        </div>
        <button @click="deleteTable(t.id)" class="text-red-400 text-sm hover:text-red-300">Delete</button>
      </div>
      <p v-if="tables.length === 0" class="text-dark-400 text-center py-8">No tables in this branch</p>
    </div>

    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showForm = false">
      <form @submit.prevent="addTable" class="bg-dark-700 p-6 rounded-lg w-full max-w-md mx-4">
        <h2 class="text-lg font-semibold text-white mb-4">New Table</h2>
        <div class="space-y-3">
          <input v-model="form.label" placeholder="Table label (e.g. Meja 1)" required
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-500 focus:border-primary-500 outline-none" />
        </div>
        <div class="flex gap-2 mt-4 justify-end">
          <button type="button" @click="showForm = false" class="px-4 py-2 bg-dark-600 text-white rounded text-sm">Cancel</button>
          <button type="submit" class="px-4 py-2 bg-primary-500 text-white rounded text-sm">Save</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { $api } = useNuxtApp()
const toast = useToast()
const showForm = ref(false)
const branches = ref<Array<{ id: string; name: string }>>([])
const selectedBranchId = ref('')
const tables = ref<Array<{ id: string; label: string; branch_id: string }>>([])
const form = reactive({ label: '' })

function branchName(id: string) {
  return branches.value.find(b => b.id === id)?.name || id.slice(0, 8)
}

async function fetchBranches() {
  try {
    const res = await $api('/api/admin/branches/list')
    branches.value = (res as any).data || []
    if (branches.value.length > 0 && !selectedBranchId.value) {
      selectedBranchId.value = branches.value[0].id
    }
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load branches')
  }
}

async function fetchTables() {
  if (!selectedBranchId.value) return
  try {
    const res = await $api(`/api/admin/tables/list?branch_id=${selectedBranchId.value}`)
    tables.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load tables')
  }
}

async function addTable() {
  if (!selectedBranchId.value) return
  try {
    await $api('/api/admin/tables', { method: 'POST', body: { branch_id: selectedBranchId.value, label: form.label } })
    showForm.value = false
    form.label = ''
    await fetchTables()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to add table')
  }
}

async function deleteTable(id: string) {
  try {
    await $api(`/api/admin/tables/delete?id=${id}`, { method: 'POST' })
    await fetchTables()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete table')
  }
}

fetchBranches()
</script>
