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
        <button v-if="selectedBranchId" @click="openCreate"
          class="px-4 py-2 bg-primary-500 text-white rounded hover:bg-primary-600 text-sm">
          Add Table
        </button>
      </div>
    </div>

    <div v-if="!selectedBranchId" class="text-dark-400 text-center py-8">Select a branch to view tables</div>

    <div v-else-if="tables.length === 0" class="text-dark-400 text-center py-8">No tables in this branch</div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <div v-for="t in tables" :key="t.id" class="bg-dark-700 rounded-lg p-4 flex flex-col items-center gap-3">
        <div class="w-full flex items-center justify-between">
          <p class="text-white font-semibold truncate">{{ t.label }}</p>
          <button @click="editTable(t)" class="text-primary-400 hover:text-primary-300 text-sm shrink-0 ml-2">
            Edit
          </button>
        </div>

        <div v-if="t.qr_code_url" class="bg-white rounded-lg p-2">
          <img :src="`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(t.qr_code_url)}`"
            :alt="`QR ${t.label}`" class="w-36 h-36 object-contain" />
        </div>
        <p v-else class="text-dark-400 text-xs">QR not available</p>

        <div class="flex gap-2 w-full">
          <button @click="regenerateQr(t.id)" class="flex-1 px-3 py-1.5 bg-blue-900/50 text-blue-400 text-xs rounded hover:bg-blue-900">
            Regen QR
          </button>
          <a v-if="t.qr_code_url"
            :href="`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(t.qr_code_url)}`"
            download :title="`QR ${t.label}`"
            class="flex-1 px-3 py-1.5 bg-dark-600 text-white text-xs rounded hover:bg-dark-500 text-center">
            Download
          </a>
          <button @click="deleteTable(t.id)" class="flex-1 px-3 py-1.5 bg-red-900/50 text-red-400 text-xs rounded hover:bg-red-900">
            Delete
          </button>
        </div>
      </div>
    </div>

    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showForm = false">
      <form @submit.prevent="editing ? updateTable() : addTable()" class="bg-dark-700 p-6 rounded-lg w-full max-w-md mx-4">
        <h2 class="text-lg font-semibold text-white mb-4">{{ editing ? 'Edit Table' : 'New Table' }}</h2>
        <div class="space-y-3">
          <input v-model="form.label" placeholder="Table label (e.g. Meja 1)" required
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-500 focus:border-primary-500 outline-none" />
        </div>
        <div v-if="editing && editing.qr_code_url" class="mt-3 flex justify-center">
          <div class="bg-white rounded-lg p-1">
            <img :src="`https://api.qrserver.com/v1/create-qr-code/?size=120x120&data=${encodeURIComponent(editing.qr_code_url)}`"
              alt="QR" class="w-20 h-20 object-contain" />
          </div>
        </div>
        <div class="flex gap-2 mt-4 justify-end">
          <button type="button" @click="showForm = false" class="px-4 py-2 bg-dark-600 text-white rounded text-sm">Cancel</button>
          <button type="submit" class="px-4 py-2 bg-primary-500 text-white rounded text-sm">{{ editing ? 'Update' : 'Save' }}</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: 'auth' })

const { $api } = useNuxtApp()
const toast = useToast()
const showForm = ref(false)
const editing = ref<{ id: string; label: string; qr_code_url?: string } | null>(null)
const branches = ref<Array<{ id: string; name: string }>>([])
const selectedBranchId = ref('')
const tables = ref<Array<{ id: string; label: string; branch_id: string; qr_code_url?: string }>>([])
const form = reactive({ label: '' })

function openCreate() {
  editing.value = null
  form.label = ''
  showForm.value = true
}

function editTable(t: { id: string; label: string; qr_code_url?: string }) {
  editing.value = t
  form.label = t.label
  showForm.value = true
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

async function updateTable() {
  if (!editing.value) return
  try {
    await $api(`/api/admin/tables/update?id=${editing.value.id}`, { method: 'POST', body: { label: form.label } })
    showForm.value = false
    editing.value = null
    form.label = ''
    await fetchTables()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to update table')
  }
}

async function regenerateQr(id: string) {
  try {
    await $api(`/api/admin/tables/regenerate-qr?id=${id}`, { method: 'POST' })
    await fetchTables()
    toast.show('QR code regenerated', 'success')
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to regenerate QR')
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
