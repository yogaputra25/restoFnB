<template>
  <div>
    <PageTitle title="Categories" description="Manage menu categories" />

    <!-- Toolbar -->
    <PageToolbar>
      <template #right>
        <AppButton variant="primary" size="sm" @click="openCreate()">
          <Plus class="w-4 h-4" />
          Add Category
        </AppButton>
      </template>
    </PageToolbar>

    <!-- Loading -->
    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <AppSkeleton v-for="i in 6" :key="i" variant="card" height="120px" />
    </div>

    <!-- Empty -->
    <AppEmptyState
      v-else-if="categories.length === 0"
      title="No Categories"
      description="Add your first category to organize menu items."
      action-label="Add Category"
      @action="openCreate()"
    />

    <!-- Category cards -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4" data-aos="fade-up">
      <div
        v-for="cat in categories"
        :key="cat.id"
        class="bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] p-5 transition-all duration-[var(--transition-base)] hover:-translate-y-1 hover:shadow-[var(--shadow-md)]"
      >
        <div class="flex items-start justify-between gap-3">
          <div class="flex-1 min-w-0">
            <h3 class="font-heading text-[var(--font-size-card-title)] text-[var(--color-text-primary)] truncate">{{ cat.name }}</h3>
            <p v-if="cat.description" class="text-[var(--color-text-secondary)] text-sm mt-1 line-clamp-2">{{ cat.description }}</p>
            <div class="flex items-center gap-3 mt-2 text-xs text-[var(--color-text-secondary)]">
              <span class="flex items-center gap-1"><ListOrdered class="w-3 h-3" /> Order: {{ cat.display_order }}</span>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-2 mt-3 pt-3 border-t border-[var(--color-border)]">
          <button class="flex-1 py-1.5 rounded-[var(--radius-button)] text-xs font-medium text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] transition-colors" @click="openEdit(cat)">Edit</button>
          <button class="flex-1 py-1.5 rounded-[var(--radius-button)] text-xs font-medium text-[var(--color-danger)] hover:bg-[var(--color-danger)]/5 transition-colors" @click="confirmDelete(cat)">Delete</button>
        </div>
      </div>
    </div>

    <!-- Add/Edit Modal -->
    <Modal
      :open="modalOpen"
      :title="editing ? 'Edit Category' : 'Add Category'"
      :confirm-text="editing ? 'Save' : 'Add'"
      :loading="saving"
      @confirm="saveCategory"
      @cancel="closeModal"
    >
      <form @submit.prevent="saveCategory" class="space-y-4">
        <AppInput v-model="form.name" label="Name" floating required />
        <AppInput v-model="form.description" label="Description" floating />
        <AppInput v-model.number="form.display_order" label="Display Order" type="number" floating />
        <div v-if="formError" class="text-xs text-[var(--color-danger)]">{{ formError }}</div>
      </form>
    </Modal>

    <!-- Delete Confirm Modal -->
    <Modal
      :open="deleteModalOpen"
      title="Delete Category"
      confirm-text="Delete"
      variant="danger"
      :loading="deleting"
      @confirm="deleteCategory"
      @cancel="deleteModalOpen = false"
    >
      <p>Are you sure you want to delete <strong>{{ deleting?.name }}</strong>? This action cannot be undone.</p>
    </Modal>
  </div>
</template>

<script setup lang="ts">
import { Plus, ListOrdered } from 'lucide-vue-next'
definePageMeta({ layout: 'admin', middleware: 'auth' })

interface Category { id: string; name: string; description?: string; display_order: number }

const { $api } = useNuxtApp()
const toast = useToast()

const loading = ref(true)
const categories = ref<Category[]>([])

// Modal
const modalOpen = ref(false)
const editing = ref<Category | null>(null)
const form = reactive({ name: '', description: '', display_order: 0 })
const saving = ref(false)
const formError = ref('')

// Delete
const deleteModalOpen = ref(false)
const deleting = ref<Category | null>(null)
const deletingLoading = ref(false)

function openCreate() {
  editing.value = null
  form.name = ''
  form.description = ''
  form.display_order = 0
  formError.value = ''
  modalOpen.value = true
}

function openEdit(cat: Category) {
  editing.value = cat
  form.name = cat.name
  form.description = cat.description || ''
  form.display_order = cat.display_order
  formError.value = ''
  modalOpen.value = true
}

function closeModal() {
  modalOpen.value = false
  editing.value = null
}

async function saveCategory() {
  if (!form.name) {
    formError.value = 'Name is required'
    return
  }
  saving.value = true
  formError.value = ''
  try {
    if (editing.value) {
      await $api('/api/admin/categories/update', {
        method: 'POST',
        body: { id: editing.value.id, name: form.name, description: form.description, display_order: Number(form.display_order) },
      })
      toast.show('Category updated')
    } else {
      await $api('/api/admin/categories', {
        method: 'POST',
        body: { name: form.name, description: form.description, display_order: Number(form.display_order) },
      })
      toast.show('Category added')
    }
    closeModal()
    await fetchCategories()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to save category')
  } finally {
    saving.value = false
  }
}

function confirmDelete(cat: Category) {
  deleting.value = cat
  deleteModalOpen.value = true
}

async function deleteCategory() {
  if (!deleting.value) return
  deletingLoading.value = true
  try {
    await $api('/api/admin/categories/delete', { method: 'POST', body: { id: deleting.value.id } })
    toast.show('Category deleted')
    deleteModalOpen.value = false
    await fetchCategories()
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete category')
  } finally {
    deletingLoading.value = false
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
  await fetchCategories()
  loading.value = false
})
</script>
