<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-white mb-6">Menu Management</h1>

    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">Categories</h2>
      <form @submit.prevent="addCategory" class="flex gap-2 mb-4 max-w-md">
        <input v-model="newCategory" placeholder="Category name"
          class="flex-1 px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none" />
        <button type="submit" class="px-4 py-2 bg-primary-500 text-white rounded hover:bg-primary-600">Add</button>
      </form>
      <ul class="space-y-2 max-w-md">
        <li v-for="cat in categories" :key="cat.id"
          class="bg-dark-700 px-4 py-3 rounded flex items-center justify-between">
          <span class="text-white">{{ cat.name }}</span>
          <button @click="deleteCategory(cat.id)" class="text-red-400 text-sm hover:text-red-300">Delete</button>
        </li>
      </ul>
    </section>

    <section>
      <h2 class="text-lg font-semibold text-white mb-4">Menu Items</h2>
      <div class="grid gap-3">
        <div v-for="item in menuItems" :key="item.id"
          class="bg-dark-700 px-4 py-3 rounded flex items-center justify-between">
          <div>
            <span class="text-white">{{ item.name }}</span>
            <span class="ml-2 text-dark-400 text-sm">Rp{{ formatPrice(item.price) }}</span>
          </div>
          <div class="flex gap-2">
            <button @click="deleteItem(item.id)" class="text-red-400 text-sm">Delete</button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'dashboard', middleware: 'auth' })

const { $api } = useNuxtApp()
const toast = useToast()
const newCategory = ref('')
const categories = ref<Array<{ id: string; name: string }>>([])
const menuItems = ref<Array<{ id: string; name: string; price: number }>>([])

async function fetchCategories() {
  try {
    const res = await $api('/api/admin/categories/list')
    categories.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load categories')
  }
}

async function fetchMenuItems() {
  try {
    const res = await $api('/api/admin/menu-items/list')
    menuItems.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load menu items')
  }
}

async function addCategory() {
  if (!newCategory.value) return
  try {
    await $api('/api/admin/categories', {
      method: 'POST',
      body: { name: newCategory.value, display_order: 0 },
    })
    newCategory.value = ''
    await fetchCategories()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to add category')
  }
}

async function deleteCategory(id: string) {
  try {
    await $api('/api/admin/categories', { method: 'POST', body: { id } })
    await fetchCategories()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete category')
  }
}

async function deleteItem(id: string) {
  try {
    await $api('/api/admin/menu-items', { method: 'POST', body: { id } })
    await fetchMenuItems()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete item')
  }
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

fetchCategories()
fetchMenuItems()
</script>
