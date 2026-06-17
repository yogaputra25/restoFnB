<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-white mb-6">Menu Management</h1>

    <section class="mb-8">
      <h2 class="text-lg font-semibold text-white mb-4">Categories</h2>
      <form @submit.prevent="addCategory" class="flex gap-2 mb-4 max-w-md">
        <input v-model="newCategoryName" placeholder="Category name"
          class="flex-1 px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none" />
        <button type="submit" class="px-4 py-2 bg-primary-500 text-white rounded hover:bg-primary-600">Add</button>
      </form>
      <div class="space-y-2 max-w-lg">
        <div v-for="cat in categories" :key="cat.id"
          class="bg-dark-700 px-4 py-3 rounded flex items-center justify-between gap-4">
          <div v-if="editingCategoryId !== cat.id" class="flex-1 min-w-0">
            <span class="text-white font-medium">{{ cat.name }}</span>
            <span v-if="cat.description" class="ml-2 text-dark-400 text-sm">{{ cat.description }}</span>
          </div>
          <div v-else class="flex-1 flex gap-2">
            <input v-model="editCategoryName" placeholder="Name"
              class="flex-1 px-2 py-1 rounded bg-dark-600 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
            <input v-model="editCategoryDesc" placeholder="Description"
              class="flex-1 px-2 py-1 rounded bg-dark-600 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
            <button @click="saveCategory(cat.id)" class="px-2 py-1 bg-green-600 text-white rounded text-xs">Save</button>
            <button @click="editingCategoryId = null" class="px-2 py-1 bg-dark-500 text-white rounded text-xs">Cancel</button>
          </div>
          <div class="flex gap-2 shrink-0">
            <button @click="startEditCategory(cat)" class="text-blue-400 text-sm hover:text-blue-300">Edit</button>
            <button @click="deleteCategory(cat.id)" class="text-red-400 text-sm hover:text-red-300">Delete</button>
          </div>
        </div>
      </div>
    </section>

    <section>
      <h2 class="text-lg font-semibold text-white mb-4">Menu Items</h2>

      <form @submit.prevent="addMenuItem" class="mb-6 p-4 bg-dark-800 rounded-lg max-w-lg space-y-3">
        <h3 class="text-white font-medium text-sm">Add Menu Item</h3>
        <div class="grid grid-cols-2 gap-3">
          <input v-model="newItem.name" placeholder="Name"
            class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
          <input v-model.number="newItem.price" type="number" placeholder="Price"
            class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
        </div>
        <select v-model="newItem.category_id"
          class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm">
          <option value="" disabled>Select category</option>
          <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
        </select>
        <input v-model="newItem.description" placeholder="Description"
          class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
        <input v-model="newItem.image_url" placeholder="Image URL (optional)"
          class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
        <label class="flex items-center gap-2 text-sm text-white">
          <input v-model="newItem.is_available" type="checkbox" class="accent-primary-500" />
          Available
        </label>
        <button type="submit" class="w-full py-2 bg-primary-500 text-white rounded hover:bg-primary-600 text-sm font-semibold">
          Add Menu Item
        </button>
      </form>

      <div class="grid gap-3">
        <div v-for="item in menuItems" :key="item.id"
          class="bg-dark-700 px-4 py-3 rounded flex items-center justify-between gap-4">
          <div class="flex items-center gap-3 flex-1 min-w-0">
            <img v-if="item.image_url" :src="item.image_url" class="w-10 h-10 rounded object-cover shrink-0" />
            <div class="min-w-0">
              <div class="flex items-center gap-2">
                <span class="text-white font-medium">{{ item.name }}</span>
                <span :class="['text-xs px-1.5 py-0.5 rounded', item.is_available ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-300']">
                  {{ item.is_available ? 'Active' : 'Inactive' }}
                </span>
              </div>
              <div class="text-dark-400 text-xs truncate">
                {{ getCategoryName(item.category_id) }} &middot; Rp{{ formatPrice(item.price) }}
              </div>
            </div>
          </div>
          <div class="flex gap-2 shrink-0">
            <button @click="startEditItem(item)" class="text-blue-400 text-sm hover:text-blue-300">Edit</button>
            <button @click="deleteItem(item.id)" class="text-red-400 text-sm hover:text-red-300">Delete</button>
          </div>
        </div>
      </div>
    </section>

    <Teleport to="body">
      <div v-if="editItemModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60" @click.self="editItemModal = false">
        <div class="bg-dark-800 rounded-lg p-6 w-full max-w-md mx-4 space-y-3">
          <h3 class="text-white font-semibold">Edit Menu Item</h3>
          <div class="grid grid-cols-2 gap-3">
            <input v-model="editItemForm.name" placeholder="Name"
              class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
            <input v-model.number="editItemForm.price" type="number" placeholder="Price"
              class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
          </div>
          <select v-model="editItemForm.category_id"
            class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm">
            <option value="" disabled>Select category</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
          <input v-model="editItemForm.description" placeholder="Description"
            class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
          <input v-model="editItemForm.image_url" placeholder="Image URL"
            class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
          <label class="flex items-center gap-2 text-sm text-white">
            <input v-model="editItemForm.is_available" type="checkbox" class="accent-primary-500" />
            Available
          </label>
          <div class="flex gap-2 pt-2">
            <button @click="saveItem" class="flex-1 py-2 bg-primary-500 text-white rounded hover:bg-primary-600 text-sm font-semibold">Save</button>
            <button @click="editItemModal = false" class="flex-1 py-2 bg-dark-500 text-white rounded hover:bg-dark-400 text-sm">Cancel</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'dashboard', middleware: 'auth' })

interface Category { id: string; name: string; description?: string; display_order?: number }
interface MenuItem { id: string; category_id: string; name: string; description?: string; price: number; image_url?: string; is_available: boolean }

const { $api } = useNuxtApp()
const toast = useToast()

const newCategoryName = ref('')
const editingCategoryId = ref<string | null>(null)
const editCategoryName = ref('')
const editCategoryDesc = ref('')

const categories = ref<Category[]>([])
const menuItems = ref<MenuItem[]>([])

const newItem = ref({ name: '', price: 0, category_id: '', description: '', image_url: '', is_available: true })
const editItemModal = ref(false)
const editItemForm = ref<MenuItem>({ id: '', category_id: '', name: '', description: '', price: 0, image_url: '', is_available: true })

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

function getCategoryName(categoryId: string): string {
  return categories.value.find(c => c.id === categoryId)?.name || 'Unknown'
}

async function addCategory() {
  if (!newCategoryName.value) return
  try {
    await $api('/api/admin/categories', {
      method: 'POST',
      body: { name: newCategoryName.value, display_order: 0 },
    })
    newCategoryName.value = ''
    await fetchCategories()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to add category')
  }
}

function startEditCategory(cat: Category) {
  editingCategoryId.value = cat.id
  editCategoryName.value = cat.name
  editCategoryDesc.value = cat.description || ''
}

async function saveCategory(id: string) {
  try {
    await $api('/api/admin/categories/update', {
      method: 'POST',
      body: { id, name: editCategoryName.value, description: editCategoryDesc.value, display_order: 0 },
    })
    editingCategoryId.value = null
    await fetchCategories()
    toast.show('Category updated')
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to update category')
  }
}

async function deleteCategory(id: string) {
  try {
    await $api('/api/admin/categories/delete', { method: 'POST', body: { id } })
    await fetchCategories()
    toast.show('Category deleted')
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete category')
  }
}

async function addMenuItem() {
  if (!newItem.value.name || !newItem.value.price || !newItem.value.category_id) {
    toast.show('Name, price, and category are required')
    return
  }
  try {
    await $api('/api/admin/menu-items', {
      method: 'POST',
      body: {
        name: newItem.value.name,
        price: newItem.value.price,
        category_id: newItem.value.category_id,
        description: newItem.value.description,
        image_url: newItem.value.image_url,
        is_available: newItem.value.is_available,
      },
    })
    newItem.value = { name: '', price: 0, category_id: '', description: '', image_url: '', is_available: true }
    await fetchMenuItems()
    toast.show('Menu item added')
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to add menu item')
  }
}

function startEditItem(item: MenuItem) {
  editItemForm.value = { ...item }
  editItemModal.value = true
}

async function saveItem() {
  try {
    await $api('/api/admin/menu-items/update', {
      method: 'POST',
      body: editItemForm.value,
    })
    editItemModal.value = false
    await fetchMenuItems()
    toast.show('Menu item updated')
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to update menu item')
  }
}

async function deleteItem(id: string) {
  try {
    await $api('/api/admin/menu-items/delete', { method: 'POST', body: { id } })
    await fetchMenuItems()
    toast.show('Menu item deleted')
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
