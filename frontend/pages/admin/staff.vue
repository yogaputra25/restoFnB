<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-white">Staff</h1>
      <button @click="showForm = true" class="px-4 py-2 bg-primary-500 text-white rounded hover:bg-primary-600 text-sm">
        Add Staff
      </button>
    </div>

    <div class="grid gap-3">
      <div v-for="u in users" :key="u.id" class="bg-dark-700 rounded-lg p-4 flex items-center justify-between">
        <div>
          <p class="text-white font-medium">{{ u.name }}</p>
          <p class="text-dark-400 text-sm">{{ u.email }} <span class="capitalize ml-2 px-1.5 py-0.5 rounded text-xs bg-dark-500 text-dark-200">{{ u.role }}</span></p>
        </div>
        <button v-if="u.role === 'cashier'" @click="deactivate(u.id)" class="text-red-400 text-sm hover:text-red-300">Deactivate</button>
      </div>
      <p v-if="users.length === 0" class="text-dark-400 text-center py-8">No staff</p>
    </div>

    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showForm = false">
      <form @submit.prevent="addStaff" class="bg-dark-700 p-6 rounded-lg w-full max-w-md mx-4">
        <h2 class="text-lg font-semibold text-white mb-4">New Staff</h2>
        <div class="space-y-3">
          <input v-model="form.name" placeholder="Name" required
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-500 focus:border-primary-500 outline-none" />
          <input v-model="form.email" type="email" placeholder="Email" required
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-500 focus:border-primary-500 outline-none" />
          <input v-model="form.password" type="password" placeholder="Password" required
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
const users = ref<Array<{ id: string; name: string; email: string; role: string }>>([])
const form = reactive({ name: '', email: '', password: '' })

async function fetchStaff() {
  try {
    const res = await $api('/api/admin/staff/list')
    users.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load staff')
  }
}

async function addStaff() {
  try {
    await $api('/api/auth/register', {
      method: 'POST',
      body: { ...form, role: 'cashier' },
    })
    showForm.value = false
    form.name = ''; form.email = ''; form.password = ''
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

fetchStaff()
</script>
