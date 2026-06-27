<template>
  <div>
    <PageTitle title="Menu Management" description="Manage categories and menu items" />

    <div class="flex gap-6">
      <!-- Left: Categories panel -->
      <div class="w-[260px] shrink-0 hidden lg:block">
        <div class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] overflow-hidden">
          <div class="p-4 border-b border-[var(--color-border)]">
            <h4 class="font-heading text-sm font-semibold text-[var(--color-text-primary)]">Categories</h4>
          </div>
          <div class="p-2 space-y-0.5">
            <button
              v-for="cat in categories"
              :key="cat.id"
              class="w-full text-left px-3 py-2 rounded-[var(--radius-button)] text-sm transition-all"
              :class="selectedCategory === cat.id ? 'bg-[var(--color-primary)]/5 text-[var(--color-primary)] font-medium' : 'text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)]'"
              @click="selectedCategory = cat.id"
            >
              {{ cat.name }}
            </button>
          </div>
          <div class="p-3 border-t border-[var(--color-border)]">
            <button class="w-full py-2 rounded-[var(--radius-button)] border border-dashed border-[var(--color-border)] text-sm text-[var(--color-text-secondary)] hover:border-[var(--color-primary)] hover:text-[var(--color-primary)] transition-colors" @click="showAddCategory = true">
              + Add Category
            </button>
          </div>
        </div>
      </div>

      <!-- Right: Menu items -->
      <div class="flex-1 min-w-0">
        <!-- Toolbar -->
        <PageToolbar>
          <template #left>
            <div class="relative flex-1 max-w-xs">
              <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--color-text-secondary)]" />
              <input v-model="searchQuery" type="text" placeholder="Search menu items..." class="w-full pl-9 pr-4 py-2 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]" />
            </div>
            <select v-model="sortBy" class="px-3 py-2 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none">
              <option value="name">Name</option>
              <option value="price">Price</option>
            </select>
          </template>
          <template #right>
            <AppButton variant="primary" size="sm" @click="openAddDrawer()">
              <Plus class="w-4 h-4" />
              Add Menu
            </AppButton>
          </template>
        </PageToolbar>

        <!-- Loading state -->
        <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
          <div v-for="i in 6" :key="i" class="h-48 rounded-[var(--radius-card)] bg-gradient-to-r from-[var(--color-surface-secondary)] via-white/40 to-[var(--color-surface-secondary)] bg-[length:200%_100%] animate-shimmer" />
        </div>

        <!-- Empty state -->
        <AppEmptyState
          v-else-if="filteredItems.length === 0"
          title="No Menu Items"
          description="Add your first menu item to get started."
          action-label="Add Menu Item"
          @action="openAddDrawer()"
        />

        <!-- Card grid -->
        <div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4" data-aos="fade-up">
          <div
            v-for="item in filteredItems"
            :key="item.id"
            class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] overflow-hidden transition-all duration-[var(--transition-base)] hover:-translate-y-1 hover:shadow-[var(--shadow-md)]"
          >
            <!-- Image -->
            <div class="h-32 bg-[var(--color-surface-secondary)] overflow-hidden">
              <img v-if="item.image_url" :src="item.image_url" :alt="item.name" class="w-full h-full object-cover" loading="lazy" />
              <div v-else class="w-full h-full flex items-center justify-center text-3xl text-[var(--color-border)]">
                <UtensilsCrossed class="w-8 h-8" />
              </div>
            </div>

            <!-- Content -->
            <div class="p-4">
              <div class="flex items-start justify-between gap-2 mb-1">
                <h4 class="font-heading text-card-title text-[var(--color-text-primary)] truncate">{{ item.name }}</h4>
                <StatusBadge :status="item.is_available ? 'active' : 'inactive'" variant="generic" />
              </div>
              <p class="text-xs text-[var(--color-text-secondary)] mb-2">{{ getCategoryName(item.category_id) }}</p>
              <p class="font-bold text-[var(--color-primary)]">Rp{{ formatPrice(item.price) }}</p>

              <!-- Actions -->
              <div class="flex items-center gap-2 mt-3 pt-3 border-t border-[var(--color-border)]">
                <button class="flex-1 py-1.5 rounded-lg text-xs font-medium text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] transition-colors" @click="openEditDrawer(item)">Edit</button>
                <button class="flex-1 py-1.5 rounded-lg text-xs font-medium text-[var(--color-danger)] hover:bg-[var(--color-danger)]/5 transition-colors" @click="confirmDelete(item)">Delete</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Drawer -->
    <Drawer :open="drawerOpen" :title="editingItem ? 'Edit Menu Item' : 'Add Menu Item'" @close="closeDrawer">
      <form @submit.prevent="saveItem" class="space-y-4">
        <div>
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">Category</label>
          <select v-model="form.category_id" required class="w-full px-3 py-2.5 rounded-[var(--radius-input)] text-sm bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]">
            <option value="" disabled>Select category</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>

        <AppInput v-model="form.name" label="Name" floating required />
        <AppInput v-model="form.description" label="Description" floating />
        <AppInput v-model.number="form.price" label="Price" type="number" floating required />

        <div>
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">Image</label>
          <div class="flex gap-2 items-start">
            <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" @change="onFileSelect" class="flex-1 text-sm file:mr-2 file:py-1 file:px-3 file:rounded-[var(--radius-button)] file:border-0 file:bg-[var(--color-primary)] file:text-white file:text-xs file:font-semibold file:cursor-pointer" />
            <img v-if="imagePreview" :src="imagePreview" class="w-12 h-12 rounded-[var(--radius-image)] object-cover shrink-0" />
          </div>
        </div>

        <div v-if="form.image_url && !selectedFile" class="mt-1">
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">Or paste URL</label>
          <AppInput v-model="form.image_url" placeholder="https://..." />
        </div>

        <label class="flex items-center gap-2 text-sm text-[var(--color-text-primary)]">
          <input v-model="form.is_available" type="checkbox" class="accent-[var(--color-primary)]" />
          Available
        </label>

        <!-- Variants -->
        <div v-if="editingItem">
          <hr class="border-[var(--color-border)] my-2" />
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-2">Variants (optional)</label>
          <div v-if="variants.length === 0" class="text-xs text-[var(--color-text-secondary)] mb-2">No variants yet</div>
          <div v-for="(v, i) in variants" :key="v.id || i" class="flex items-center gap-2 mb-1.5 text-sm">
            <span class="flex-1 text-[var(--color-text-primary)]">{{ v.name }}</span>
            <span v-if="v.price_adjustment" class="text-[var(--color-text-secondary)] text-xs">+Rp{{ v.price_adjustment }}</span>
            <button @click="deleteVariant(v.id)" class="text-[var(--color-danger)] text-xs hover:underline">Delete</button>
          </div>
          <div class="flex gap-2 items-center mt-2">
            <input v-model="variantForm.name" placeholder="Variant name (e.g. Level 1)" class="flex-1 px-2.5 py-1.5 rounded-[var(--radius-input)] text-xs bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-1 focus:ring-[var(--color-primary)]" />
            <input v-model.number="variantForm.price_adjustment" type="number" placeholder="+Price" class="w-20 px-2 py-1.5 rounded-[var(--radius-input)] text-xs bg-[var(--color-surface-secondary)] text-[var(--color-text-primary)] border border-[var(--color-border)] focus:outline-none focus:ring-1 focus:ring-[var(--color-primary)]" />
            <AppButton variant="primary" size="sm" @click="addVariant">Add</AppButton>
          </div>
        </div>

        <div v-if="formError" class="text-xs text-[var(--color-danger)]">{{ formError }}</div>

        <div class="flex gap-3 pt-2">
          <AppButton variant="secondary" class="flex-1" @click="closeDrawer">Cancel</AppButton>
          <AppButton variant="primary" class="flex-1" type="submit" :loading="saving">{{ editingItem ? 'Save' : 'Add' }}</AppButton>
        </div>
      </form>
    </Drawer>

    <!-- Delete Modal -->
    <Modal
      :open="deleteModalOpen"
      title="Delete Menu Item"
      confirm-text="Delete"
      variant="danger"
      :loading="deleting"
      @confirm="deleteItem"
      @cancel="deleteModalOpen = false"
    >
      <p>Are you sure you want to delete <strong>{{ deletingItem?.name }}</strong>? This action cannot be undone.</p>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { Search, Plus, UtensilsCrossed } from 'lucide-vue-next'
definePageMeta({ layout: 'admin', middleware: 'auth' })

interface Category { id: string; name: string; description?: string }
interface MenuItem { id: string; category_id: string; name: string; description?: string; price: number; image_url?: string; is_available: boolean }

const { $api } = useNuxtApp()
const auth = useAuthStore()
const toast = useToast()

const loading = ref(true)
const categories = ref<Category[]>([])
const menuItems = ref<MenuItem[]>([])
const selectedCategory = ref<string | null>(null)
const searchQuery = ref('')
const sortBy = ref('name')

// Drawer
const drawerOpen = ref(false)
const editingItem = ref<MenuItem | null>(null)
const form = reactive({ category_id: '', name: '', description: '', price: 0 as any, image_url: '', is_available: true })
const selectedFile = ref<File | null>(null)
const imagePreview = ref('')
const saving = ref(false)
const formError = ref('')

// Variants
interface Variant { id: string; menu_item_id?: string; name: string; price_adjustment: number }
const variants = ref<Variant[]>([])
const variantForm = reactive({ name: '', price_adjustment: 0 })

async function addVariant() {
  if (!variantForm.name || !editingItem.value) return
  try {
    await $api('/api/admin/variants', {
      method: 'POST',
      body: { menu_item_id: editingItem.value.id, name: variantForm.name, price_adjustment: Number(variantForm.price_adjustment) || 0 },
    })
    variantForm.name = ''
    variantForm.price_adjustment = 0
    await fetchVariants(editingItem.value.id)
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to add variant')
  }
}

async function deleteVariant(id: string) {
  try {
    await $api('/api/admin/variants/delete', { method: 'POST', body: { id } })
    if (editingItem.value) await fetchVariants(editingItem.value.id)
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete variant')
  }
}

async function fetchVariants(menuItemId: string) {
  try {
    const res = await $api(`/api/admin/variants/list?menu_item_id=${menuItemId}`)
    variants.value = (res as any).data || []
  } catch { variants.value = [] }
}

// Delete
const deleteModalOpen = ref(false)
const deletingItem = ref<MenuItem | null>(null)
const deleting = ref(false)

const filteredItems = computed(() => {
  let items = menuItems.value
  if (selectedCategory.value) {
    items = items.filter(i => i.category_id === selectedCategory.value)
  }
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    items = items.filter(i => i.name.toLowerCase().includes(q))
  }
  if (sortBy.value === 'price') {
    items = [...items].sort((a, b) => a.price - b.price)
  } else {
    items = [...items].sort((a, b) => a.name.localeCompare(b.name))
  }
  return items
})

function getCategoryName(id: string) {
  return categories.value.find(c => c.id === id)?.name || 'Unknown'
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

// Upload
async function upload(file: File): Promise<string> {
  const fd = new FormData()
  fd.append('file', file)
  const res = await $fetch(`/api/upload?chain_id=${auth.user?.chainId}`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${auth.accessToken}` },
    body: fd,
  })
  return (res as any).url || (res as any).data?.url || ''
}

function onFileSelect(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  selectedFile.value = file
  const reader = new FileReader()
  reader.onload = () => { imagePreview.value = reader.result as string }
  reader.readAsDataURL(file)
}

// Drawer
function openAddDrawer() {
  editingItem.value = null
  form.category_id = ''
  form.name = ''
  form.description = ''
  form.price = ''
  form.image_url = ''
  form.is_available = true
  selectedFile.value = null
  imagePreview.value = ''
  formError.value = ''
  variants.value = []
  drawerOpen.value = true
}

function openEditDrawer(item: MenuItem) {
  editingItem.value = item
  form.category_id = item.category_id
  form.name = item.name
  form.description = item.description || ''
  form.price = item.price
  form.image_url = item.image_url || ''
  form.is_available = item.is_available
  selectedFile.value = null
  imagePreview.value = ''
  formError.value = ''
  drawerOpen.value = true
  fetchVariants(item.id)
}

function closeDrawer() {
  drawerOpen.value = false
  editingItem.value = null
  variants.value = []
}

async function saveItem() {
  if (!form.name || !form.price || !form.category_id) {
    formError.value = 'Name, price, and category are required'
    return
  }
  saving.value = true
  formError.value = ''
  let imageUrl = form.image_url
  if (selectedFile.value) {
    try { imageUrl = await upload(selectedFile.value) } catch {
      formError.value = 'Failed to upload image'; saving.value = false; return
    }
  }
  try {
    if (editingItem.value) {
      await $api('/api/admin/menu-items/update', {
        method: 'POST',
        body: { id: editingItem.value.id, category_id: form.category_id, name: form.name, description: form.description, price: Number(form.price), image_url: imageUrl, is_available: form.is_available },
      })
      toast.show('Menu item updated')
    } else {
      await $api('/api/admin/menu-items', {
        method: 'POST',
        body: { category_id: form.category_id, name: form.name, description: form.description, price: Number(form.price), image_url: imageUrl, is_available: form.is_available },
      })
      toast.show('Menu item added')
    }
    closeDrawer()
    await fetchMenuItems()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to save item')
  } finally {
    saving.value = false
  }
}

// Delete
function confirmDelete(item: MenuItem) {
  deletingItem.value = item
  deleteModalOpen.value = true
}

async function deleteItem() {
  if (!deletingItem.value) return
  deleting.value = true
  try {
    await $api('/api/admin/menu-items/delete', { method: 'POST', body: { id: deletingItem.value.id } })
    toast.show('Menu item deleted')
    deleteModalOpen.value = false
    await fetchMenuItems()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete item')
  } finally {
    deleting.value = false
  }
}

async function fetchMenuItems() {
  try {
    const res = await $api('/api/admin/menu-items/list')
    menuItems.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load items')
  }
}

async function fetchCategories() {
  try {
    const res = await $api('/api/admin/categories/list')
    categories.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load categories')
  }
}

onMounted(async () => {
  loading.value = true
  await Promise.all([fetchCategories(), fetchMenuItems()])
  loading.value = false
})
</script>

<style scoped>
@keyframes shimmer {
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
}
.animate-shimmer {
  animation: shimmer 1.5s ease-in-out infinite;
}
</style>
