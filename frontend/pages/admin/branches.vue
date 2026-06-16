<template>
  <div class="p-6">
    <div class="flex items-center justify-between mb-6">
      <h1 class="text-2xl font-bold text-white">Branches</h1>
      <button @click="showForm = true" class="px-4 py-2 bg-primary-500 text-white rounded hover:bg-primary-600 text-sm">
        Add Branch
      </button>
    </div>

    <div class="grid gap-3">
      <div v-for="b in branches" :key="b.id"
        class="bg-dark-700 rounded-lg p-4 flex items-center justify-between">
        <div>
          <p class="text-white font-medium">{{ b.name }}</p>
          <p class="text-dark-400 text-sm">{{ b.address || 'No address' }}</p>
        </div>
        <button @click="deleteBranch(b.id)" class="text-red-400 text-sm hover:text-red-300">Delete</button>
      </div>
      <p v-if="branches.length === 0" class="text-dark-400 text-center py-8">No branches</p>
    </div>

    <div v-if="showForm" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="showForm = false">
      <form @submit.prevent="addBranch" class="bg-dark-700 p-6 rounded-lg w-full max-w-md mx-4">
        <h2 class="text-lg font-semibold text-white mb-4">New Branch</h2>
        <div class="space-y-3">
          <input v-model="form.name" placeholder="Branch name" required
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-500 focus:border-primary-500 outline-none" />
          <input v-model="form.address" placeholder="Address"
            class="w-full px-3 py-2 rounded bg-dark-600 text-white border border-dark-500 focus:border-primary-500 outline-none" />
          <input v-model="form.contact_info" placeholder="Contact info"
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
const branches = ref<Array<{ id: string; name: string; address: string }>>([])
const form = reactive({ name: '', address: '', contact_info: '' })

async function fetchBranches() {
  try {
    const res = await $api('/api/admin/branches/list')
    branches.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load branches')
  }
}

async function addBranch() {
  try {
    await $api('/api/admin/branches', { method: 'POST', body: form })
    showForm.value = false
    form.name = ''; form.address = ''; form.contact_info = ''
    await fetchBranches()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to add branch')
  }
}

async function deleteBranch(id: string) {
  try {
    await $api('/api/admin/branches', { method: 'POST', body: { id } })
    await fetchBranches()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete branch')
  }
}

fetchBranches()
</script>
