<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h1 class="font-heading text-[var(--font-size-page-title)] text-[var(--color-text-primary)]">Tables</h1>
      <div class="flex items-center gap-3">
        <select v-model="selectedBranchId" @change="fetchTables"
          class="px-3 py-2 rounded-[var(--radius-input)] bg-[var(--color-surface)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]">
          <option value="" disabled>Select branch</option>
          <option v-for="b in branches" :key="b.id" :value="b.id">{{ b.name }}</option>
        </select>
        <AppButton v-if="selectedBranchId" variant="primary" size="sm" @click="openCreate">
          Add Table
        </AppButton>
      </div>
    </div>

    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <AppSkeleton v-for="i in 8" :key="i" variant="card" />
    </div>

    <AppEmptyState
      v-else-if="!selectedBranchId"
      title="Select Branch"
      description="Select a branch to view tables"
    />

    <AppEmptyState
      v-else-if="tables.length === 0"
      title="No Tables"
      description="Add your first table to this branch."
      action-label="Add Table"
      @action="openCreate"
    />

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
      <div v-for="t in tables" :key="t.id" class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-4 flex flex-col items-center gap-3 transition-all hover:-translate-y-1 hover:shadow-[var(--shadow-md)]">
        <div class="w-full flex items-center justify-between">
          <p class="text-[var(--color-text-primary)] font-semibold truncate">{{ t.label }}</p>
          <button @click="editTable(t)" class="text-[var(--color-primary)] hover:underline text-sm shrink-0 ml-2">
            Edit
          </button>
        </div>

        <div v-if="t.qr_code_url" class="bg-white rounded-[var(--radius-card)] p-2">
          <img :src="`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(t.qr_code_url)}`"
            :alt="`QR ${t.label}`" class="w-36 h-36 object-contain" />
        </div>
        <p v-else class="text-[var(--color-text-secondary)] text-xs">QR not available</p>

        <div class="flex gap-2 w-full">
          <button @click="regenerateQr(t.id)" class="flex-1 px-3 py-1.5 bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)] text-xs rounded-[var(--radius-button)] hover:bg-[var(--color-border)] transition-all hover:scale-[1.03]">
            Regen QR
          </button>
          <a v-if="t.qr_code_url"
            :href="`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(t.qr_code_url)}`"
            download :title="`QR ${t.label}`"
            class="flex-1 px-3 py-1.5 bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)] text-xs rounded-[var(--radius-button)] hover:bg-[var(--color-border)] text-center transition-all">
            Download
          </a>
          <button @click="deleteTable(t.id)" class="flex-1 px-3 py-1.5 bg-[var(--color-danger)]/10 text-[var(--color-danger)] text-xs rounded-[var(--radius-button)] hover:bg-[var(--color-danger)]/20 transition-all hover:scale-[1.03]">
            Delete
          </button>
        </div>
      </div>
    </div>

    <Teleport v-if="showForm" to="body">
      <div class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/40" @click="showForm = false" />
        <form @submit.prevent="editing ? updateTable() : addTable()" class="relative bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-lg)] w-full max-w-md mx-4 p-6">
          <h2 class="font-heading text-[var(--font-size-section)] text-[var(--color-text-primary)] mb-4">{{ editing ? 'Edit Table' : 'New Table' }}</h2>
          <div class="space-y-3">
            <AppInput v-model="form.label" placeholder="Table label (e.g. Meja 1)" required floating />
          </div>
          <div v-if="editing && editing.qr_code_url" class="mt-3 flex justify-center">
            <div class="bg-white rounded-[var(--radius-card)] p-1">
              <img :src="`https://api.qrserver.com/v1/create-qr-code/?size=120x120&data=${encodeURIComponent(editing.qr_code_url)}`"
                alt="QR" class="w-20 h-20 object-contain" />
            </div>
          </div>
          <div class="flex gap-2 mt-4 justify-end">
            <AppButton variant="secondary" @click="showForm = false">Cancel</AppButton>
            <AppButton variant="primary" type="submit">{{ editing ? 'Update' : 'Save' }}</AppButton>
          </div>
        </form>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: 'auth' })

const { $api } = useNuxtApp()
const toast = useToast()
const showForm = ref(false)
const editing = ref<{ id: string; label: string; qr_code_url?: string } | null>(null)
const loading = ref(false)
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
  loading.value = true
  try {
    const res = await $api(`/api/admin/tables/list?branch_id=${selectedBranchId.value}`)
    tables.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load tables')
  }
  loading.value = false
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
