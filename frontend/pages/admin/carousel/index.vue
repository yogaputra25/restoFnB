<template>
  <div class="p-6">
    <h1 class="font-heading text-[var(--font-size-page-title)] text-[var(--color-text-primary)] mb-6">Carousel Management</h1>

    <form @submit.prevent="addSlide" class="mb-6 p-4 bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] max-w-lg space-y-3">
      <h3 class="font-heading text-sm font-semibold text-[var(--color-text-primary)]">Add Carousel Slide</h3>
      <div class="grid grid-cols-2 gap-3">
        <AppInput v-model="newSlide.title" placeholder="Title" required floating />
        <AppInput v-model.number="newSlide.display_order" type="number" placeholder="Display order" floating />
      </div>
      <AppInput v-model="newSlide.description" placeholder="Description" floating />
      <div class="flex gap-2 items-start">
        <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" @change="onAddImageSelect"
          class="flex-1 text-sm file:mr-2 file:py-1 file:px-3 file:rounded-[var(--radius-button)] file:border-0 file:bg-[var(--color-primary)] file:text-white file:text-xs file:font-semibold file:cursor-pointer" />
        <img v-if="addImagePreview" :src="addImagePreview" class="w-12 h-8 rounded-[var(--radius-image)] object-cover shrink-0" />
      </div>
      <AppInput v-if="!addImagePreview" v-model="newSlide.image_url" placeholder="...or paste Image URL" floating />
      <div class="flex gap-3">
        <label class="flex items-center gap-2 text-sm text-[var(--color-text-primary)]">
          <input v-model="newSlide.is_active" type="checkbox" class="accent-[var(--color-primary)]" />
          Active
        </label>
        <AppInput v-model="newSlide.bg_color" placeholder="~or bg color (e.g. #FF6B00)" floating />
      </div>
      <AppButton variant="primary" size="sm" class="w-full" type="submit">
        Add Slide
      </AppButton>
    </form>

    <div v-if="loading" class="space-y-3 max-w-lg">
      <AppSkeleton v-for="i in 3" :key="i" variant="card" />
    </div>

    <div v-else-if="slides.length === 0" class="max-w-lg">
      <AppEmptyState title="No Carousel Slides" description="Add your first carousel slide to get started." />
    </div>

    <div v-else class="space-y-3 max-w-lg">
      <div v-for="slide in slides" :key="slide.id"
        class="bg-[var(--color-surface)] px-4 py-3 rounded-[var(--radius-card)] shadow-[var(--shadow-sm)] flex items-center justify-between gap-4 transition-all hover:-translate-y-1 hover:shadow-[var(--shadow-md)]">
        <div class="flex items-center gap-3 flex-1 min-w-0">
          <img v-if="slide.image_url" :src="slide.image_url" class="w-12 h-8 rounded-[var(--radius-image)] object-cover shrink-0" />
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <span class="text-[var(--color-text-primary)] font-medium text-sm">{{ slide.title }}</span>
              <StatusBadge :status="slide.is_active ? 'active' : 'inactive'" variant="generic" />
            </div>
            <div class="text-[var(--color-text-secondary)] text-xs truncate">
              {{ slide.description || 'No description' }} &middot; Order: {{ slide.display_order }}
            </div>
          </div>
        </div>
        <div class="flex gap-2 shrink-0">
          <button @click="startEditSlide(slide)" class="text-[var(--color-primary)] text-sm hover:underline">Edit</button>
          <button @click="deleteSlide(slide.id)" class="text-[var(--color-danger)] text-sm hover:underline">Delete</button>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="editModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/40" @click="editModal = false" />
        <div class="relative bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-lg)] p-6 w-full max-w-md mx-4 space-y-3">
          <h3 class="font-heading text-[var(--font-size-section)] text-[var(--color-text-primary)]">Edit Carousel Slide</h3>
          <div class="grid grid-cols-2 gap-3">
            <AppInput v-model="editForm.title" placeholder="Title" required floating />
            <AppInput v-model.number="editForm.display_order" type="number" placeholder="Display order" floating />
          </div>
          <AppInput v-model="editForm.description" placeholder="Description" floating />
          <div class="flex gap-2 items-start">
            <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" @change="onEditImageSelect"
              class="flex-1 text-sm file:mr-2 file:py-1 file:px-3 file:rounded-[var(--radius-button)] file:border-0 file:bg-[var(--color-primary)] file:text-white file:text-xs file:font-semibold file:cursor-pointer" />
            <img v-if="editImagePreview" :src="editImagePreview" class="w-12 h-8 rounded-[var(--radius-image)] object-cover shrink-0" />
          </div>
          <AppInput v-if="!editFileSelected" v-model="editForm.image_url" placeholder="Image URL" floating />
          <div class="flex gap-3">
            <label class="flex items-center gap-2 text-sm text-[var(--color-text-primary)]">
              <input v-model="editForm.is_active" type="checkbox" class="accent-[var(--color-primary)]" />
              Active
            </label>
            <AppInput v-model="editForm.bg_color" placeholder="~or bg color" floating />
          </div>
          <div class="flex gap-2 pt-2">
            <AppButton variant="primary" class="flex-1" @click="saveSlide">Save</AppButton>
            <AppButton variant="secondary" class="flex-1" @click="editModal = false">Cancel</AppButton>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ layout: 'admin', middleware: 'auth' })

interface CarouselSlide {
  id: string
  title: string
  description?: string
  image_url?: string
  bg_color?: string
  display_order: number
  is_active: boolean
}

const { $api } = useNuxtApp()
const auth = useAuthStore()
const toast = useToast()

const slides = ref<CarouselSlide[]>([])
const loading = ref(false)
const newSlide = ref({ title: '', description: '', image_url: '', bg_color: '', display_order: 0, is_active: true })
const editModal = ref(false)
const editForm = ref<CarouselSlide>({ id: '', title: '', description: '', image_url: '', bg_color: '', display_order: 0, is_active: true })

const addSelectedFile = ref<File | null>(null)
const addImagePreview = ref('')
const editSelectedFile = ref<File | null>(null)
const editImagePreview = ref('')
const editFileSelected = computed(() => editSelectedFile.value !== null)

async function uploadImage(file: File): Promise<string> {
  const formData = new FormData()
  formData.append('file', file)
  const res = await $fetch(`/api/upload?chain_id=${auth.user?.chainId}`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${auth.accessToken}` },
    body: formData,
  })
  return (res as any).url || (res as any).data?.url || ''
}

function onAddImageSelect(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  addSelectedFile.value = file
  const reader = new FileReader()
  reader.onload = () => { addImagePreview.value = reader.result as string }
  reader.readAsDataURL(file)
}

function onEditImageSelect(e: Event) {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  editSelectedFile.value = file
  const reader = new FileReader()
  reader.onload = () => { editImagePreview.value = reader.result as string }
  reader.readAsDataURL(file)
}

async function fetchSlides() {
  loading.value = true
  try {
    const res = await $api('/api/admin/carousel/list')
    slides.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load slides')
  }
  loading.value = false
}

async function addSlide() {
  if (!newSlide.value.title) {
    toast.show('Title is required')
    return
  }
  let imageUrl = newSlide.value.image_url
  if (addSelectedFile.value) {
    try {
      imageUrl = await uploadImage(addSelectedFile.value)
    } catch {
      toast.show('Failed to upload image')
      return
    }
  }
  try {
    await $api('/api/admin/carousel', {
      method: 'POST',
      body: {
        title: newSlide.value.title,
        description: newSlide.value.description,
        image_url: imageUrl,
        bg_color: newSlide.value.bg_color,
        display_order: newSlide.value.display_order,
        is_active: newSlide.value.is_active,
      },
    })
    newSlide.value = { title: '', description: '', image_url: '', bg_color: '', display_order: 0, is_active: true }
    addSelectedFile.value = null
    addImagePreview.value = ''
    await fetchSlides()
    toast.show('Slide added')
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to add slide')
  }
}

function startEditSlide(slide: CarouselSlide) {
  editForm.value = { ...slide }
  editModal.value = true
}

async function saveSlide() {
  let imageUrl = editForm.value.image_url
  if (editSelectedFile.value) {
    try {
      imageUrl = await uploadImage(editSelectedFile.value)
    } catch {
      toast.show('Failed to upload image')
      return
    }
  }
  try {
    await $api('/api/admin/carousel/update', {
      method: 'POST',
      body: { ...editForm.value, image_url: imageUrl },
    })
    editModal.value = false
    editSelectedFile.value = null
    editImagePreview.value = ''
    await fetchSlides()
    toast.show('Slide updated')
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to update slide')
  }
}

async function deleteSlide(id: string) {
  if (!confirm('Delete this slide?')) return
  try {
    await $api('/api/admin/carousel/delete', { method: 'POST', body: { id } })
    await fetchSlides()
    toast.show('Slide deleted')
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to delete slide')
  }
}

fetchSlides()
</script>
