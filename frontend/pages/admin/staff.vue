<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h1 class="font-heading text-[var(--font-size-page-title)] text-[var(--color-text-primary)]">Staff</h1>
      <AppButton variant="primary" size="sm" @click="showForm = true">
        Add Staff
      </AppButton>
    </div>

    <div v-if="loading" class="grid gap-3">
      <AppSkeleton v-for="i in 3" :key="i" variant="card" />
    </div>

    <template v-else>
      <AppEmptyState
        v-if="users.length === 0"
        title="No Staff"
        description="Add your first staff member to get started."
        action-label="Add Staff"
        @action="showForm = true"
      />

      <div v-else class="grid gap-3">
        <div v-for="u in users" :key="u.id" class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-4 flex items-center justify-between transition-all hover:-translate-y-1 hover:shadow-[var(--shadow-md)]" :class="u.is_active === false ? 'opacity-60' : ''">
          <div>
            <p class="text-[var(--color-text-primary)] font-medium">{{ u.name }}</p>
            <p class="text-[var(--color-text-secondary)] text-sm">{{ u.email }} <span class="capitalize ml-2 px-1.5 py-0.5 rounded-[var(--radius-button)] text-xs bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)]">{{ u.role }}</span> <span v-if="u.is_active === false" class="px-1.5 py-0.5 rounded-[var(--radius-button)] text-xs bg-[var(--color-danger)]/10 text-[var(--color-danger)]">Inactive</span></p>
            <p v-if="u.branch_id" class="text-[10px] text-[var(--color-text-secondary)] mt-0.5">{{ getBranchName(u.branch_id) }}</p>
          </div>
          <button v-if="u.role === 'cashier' && u.is_active !== false" @click="deactivate(u.id)" class="text-[var(--color-danger)] text-sm hover:underline">Deactivate</button>
          <button v-if="u.role === 'cashier' && u.is_active === false" @click="reactivate(u.id)" class="text-[var(--color-success)] text-sm hover:underline">Reactivate</button>
        </div>
      </div>
    </template>

    <Teleport v-if="showForm" to="body">
      <div class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/40" @click="showForm = false" />
        <form @submit.prevent="addStaff" class="relative bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-lg)] w-full max-w-md mx-4 p-6">
          <h2 class="font-heading text-[var(--font-size-section)] text-[var(--color-text-primary)] mb-4">New Staff</h2>
          <div class="space-y-3">
            <AppInput v-model="form.name" placeholder="Name" required floating />
            <AppInput v-model="form.email" type="email" placeholder="Email" required floating />
            <AppInput v-model="form.password" type="password" placeholder="Password" required floating />
            <div>
              <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">Branch</label>
              <select v-model="form.branch_id" class="w-full px-3 py-2.5 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]">
                <option value="">Select branch</option>
                <option v-for="b in branches" :key="b.id" :value="b.id">{{ b.name }}</option>
              </select>
            </div>
          </div>
          <div class="flex gap-2 mt-4 justify-end">
            <AppButton variant="secondary" @click="showForm = false">Cancel</AppButton>
            <AppButton variant="primary" type="submit">Save</AppButton>
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
const { branches, fetchBranches, getBranchName } = useBranchSelector()

const showForm = ref(false)
const loading = ref(true)
const users = ref<Array<{ id: string; name: string; email: string; role: string; branch_id?: string; is_active?: boolean }>>([])
const form = reactive({ name: '', email: '', password: '', branch_id: '' })

async function fetchStaff() {
  loading.value = true
  try {
    const res = await $api('/api/admin/staff/list')
    users.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load staff')
  }
  loading.value = false
}

async function addStaff() {
  try {
    await $api('/api/auth/register', {
      method: 'POST',
      body: { name: form.name, email: form.email, password: form.password, role: 'cashier', branch_id: form.branch_id || undefined },
    })
    showForm.value = false
    form.name = ''; form.email = ''; form.password = ''; form.branch_id = ''
    await fetchStaff()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to add staff')
  }
}

async function deactivate(id: string) {
  try {
    await $api(`/api/admin/staff/deactivate?id=${id}`, { method: 'POST' })
    await fetchStaff()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to deactivate staff')
  }
}

async function reactivate(id: string) {
  try {
    await $api(`/api/admin/staff/reactivate?id=${id}`, { method: 'POST' })
    await fetchStaff()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to reactivate staff')
  }
}

fetchBranches()
fetchStaff()
</script>
