<template>
  <div>
    <PageTitle title="Branches" description="Manage your restaurant branches" />

    <!-- Toolbar -->
    <PageToolbar>
      <template #left>
        <div class="relative flex-1 max-w-xs">
          <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-text-secondary)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input v-model="searchQuery" type="text" placeholder="Search branches..." class="w-full pl-9 pr-4 py-2 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]" />
        </div>
      </template>
      <template #right>
        <AppButton variant="primary" size="sm" @click="openAddDrawer">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" /></svg>
          Add Branch
        </AppButton>
      </template>
    </PageToolbar>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
      <div v-for="i in 6" :key="i" class="h-48 rounded-[var(--radius-card)] bg-gradient-to-r from-[var(--color-surface-secondary)] via-white/40 to-[var(--color-surface-secondary)] bg-[length:200%_100%] animate-shimmer" />
    </div>

    <!-- Empty -->
    <AppEmptyState
      v-else-if="filteredBranches.length === 0"
      title="No Branches"
      description="Add your first branch to get started."
      action-label="Add Branch"
      @action="openAddDrawer"
    />

    <!-- Branch cards -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4" data-aos="fade-up">
      <div
        v-for="b in filteredBranches"
        :key="b.id"
        class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] overflow-hidden transition-all duration-[var(--transition-base)] hover:-translate-y-1 hover:shadow-[var(--shadow-md)]"
      >
        <!-- Photo -->
        <div class="h-36 bg-[var(--color-surface-secondary)] overflow-hidden">
          <img v-if="b.image_url" :src="b.image_url" :alt="b.name" class="w-full h-full object-cover" loading="lazy" />
          <div v-else class="w-full h-full flex items-center justify-center">
            <svg class="w-12 h-12 text-[var(--color-border)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
        </div>

        <!-- Info -->
        <div class="p-4 space-y-2">
          <div class="flex items-start justify-between">
            <h4 class="font-heading text-card-title text-[var(--color-text-primary)]">{{ b.name }}</h4>
            <StatusBadge :status="getStatus(b)" variant="generic" />
          </div>

          <div class="space-y-1 text-xs text-[var(--color-text-secondary)]">
            <p v-if="b.address" class="flex items-center gap-1.5">
              <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" /><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
              {{ b.address }}
            </p>
            <p v-if="b.contact_info" class="flex items-center gap-1.5">
              <svg class="w-3.5 h-3.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" /></svg>
              {{ b.contact_info }}
            </p>
          </div>

          <!-- Stats row -->
          <div class="grid grid-cols-3 gap-2 pt-2 border-t border-[var(--color-border)] text-center">
            <div>
              <p class="font-bold text-sm text-[var(--color-text-primary)]">{{ b.tables_count || 0 }}</p>
              <p class="text-[10px] text-[var(--color-text-secondary)]">Tables</p>
            </div>
            <div>
              <p class="font-bold text-sm text-[var(--color-text-primary)]">{{ b.staff_count || 0 }}</p>
              <p class="text-[10px] text-[var(--color-text-secondary)]">Staff</p>
            </div>
            <div>
              <p class="font-bold text-sm text-[var(--color-success)]">Rp{{ formatPrice(b.today_revenue || 0) }}</p>
              <p class="text-[10px] text-[var(--color-text-secondary)]">Revenue</p>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex gap-2 pt-1">
            <button class="flex-1 py-1.5 rounded-lg text-xs font-medium text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] transition-colors" @click="openEditDrawer(b)">Edit</button>
            <button class="flex-1 py-1.5 rounded-lg text-xs font-medium text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] transition-colors" @click="viewOnMaps(b)">View on Maps</button>
            <button class="flex-1 py-1.5 rounded-lg text-xs font-medium text-[var(--color-danger)] hover:bg-[var(--color-danger)]/5 transition-colors" @click="confirmDelete(b)">Delete</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Drawer -->
    <Drawer :open="drawerOpen" :title="editingBranch ? 'Edit Branch' : 'Add Branch'" @close="closeDrawer">
      <form @submit.prevent="saveBranch" class="space-y-4">
        <AppInput v-model="form.name" label="Branch Name" floating required />
        <AppInput v-model="form.address" label="Address" floating />
        <AppInput v-model="form.contact_info" label="Contact Info" floating />
        <div class="flex gap-3 pt-2">
          <AppButton variant="secondary" class="flex-1" @click="closeDrawer">Cancel</AppButton>
          <AppButton variant="primary" class="flex-1" type="submit" :loading="saving">{{ editingBranch ? 'Save' : 'Add' }}</AppButton>
        </div>
      </form>
    </Drawer>

    <!-- Delete Modal -->
    <Modal
      :open="deleteModalOpen"
      title="Delete Branch"
      confirm-text="Delete"
      variant="danger"
      :loading="deleting"
      @confirm="deleteBranch"
      @cancel="deleteModalOpen = false"
    >
      <p>Are you sure you want to delete <strong>{{ deletingBranch?.name }}</strong>?</p>
    </Modal>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: 'auth' })

interface Branch {
  id: string
  name: string
  address?: string
  contact_info?: string
  image_url?: string
  is_active?: boolean
  tables_count?: number
  staff_count?: number
  today_revenue?: number
  [key: string]: any
}

const { $api } = useNuxtApp()
const toast = useToast()

const loading = ref(true)
const branches = ref<Branch[]>([])
const searchQuery = ref('')

// Drawer
const drawerOpen = ref(false)
const editingBranch = ref<Branch | null>(null)
const form = reactive({ name: '', address: '', contact_info: '' })
const saving = ref(false)

// Delete
const deleteModalOpen = ref(false)
const deletingBranch = ref<Branch | null>(null)
const deleting = ref(false)

const filteredBranches = computed(() => {
  if (!searchQuery.value) return branches.value
  const q = searchQuery.value.toLowerCase()
  return branches.value.filter(b => b.name.toLowerCase().includes(q) || (b.address || '').toLowerCase().includes(q))
})

function formatPrice(price: number) { return price.toLocaleString('id-ID') }

function getStatus(b: Branch) { return b.is_active !== false ? 'active' : 'inactive' }

function viewOnMaps(b: Branch) {
  if (b.address) window.open(`https://maps.google.com/?q=${encodeURIComponent(b.address)}`, '_blank')
}

// Drawer
function openAddDrawer() {
  editingBranch.value = null
  form.name = ''; form.address = ''; form.contact_info = ''
  drawerOpen.value = true
}

function openEditDrawer(b: Branch) {
  editingBranch.value = b
  form.name = b.name; form.address = b.address || ''; form.contact_info = b.contact_info || ''
  drawerOpen.value = true
}

function closeDrawer() { drawerOpen.value = false; editingBranch.value = null }

async function saveBranch() {
  saving.value = true
  try {
    if (editingBranch.value) {
      await $api('/api/admin/branches', {
        method: 'POST',
        body: { id: editingBranch.value.id, name: form.name, address: form.address, contact_info: form.contact_info },
      })
      toast.show('Branch updated')
    } else {
      await $api('/api/admin/branches', {
        method: 'POST',
        body: { name: form.name, address: form.address, contact_info: form.contact_info },
      })
      toast.show('Branch added')
    }
    closeDrawer()
    await fetchBranches()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to save branch')
  } finally { saving.value = false }
}

// Delete
function confirmDelete(b: Branch) { deletingBranch.value = b; deleteModalOpen.value = true }

async function deleteBranch() {
  if (!deletingBranch.value) return
  deleting.value = true
  try {
    await $api('/api/admin/branches', { method: 'POST', body: { id: deletingBranch.value.id } })
    toast.show('Branch deleted')
    deleteModalOpen.value = false
    await fetchBranches()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete branch')
  } finally { deleting.value = false }
}

async function fetchBranches() {
  try {
    const res = await $api('/api/admin/branches/list')
    branches.value = (res as any).data || []
  } catch (e: any) { toast.show(e?.data?.message || 'Failed to load branches') }
}

onMounted(async () => {
  loading.value = true
  await fetchBranches()
  loading.value = false
})
</script>

<style scoped>
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}
.animate-shimmer { animation: shimmer 1.5s ease-in-out infinite; }
</style>
