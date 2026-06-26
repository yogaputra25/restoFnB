<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold text-white mb-6">Carousel Management</h1>

    <form @submit.prevent="addSlide" class="mb-6 p-4 bg-dark-800 rounded-lg max-w-lg space-y-3">
      <h3 class="text-white font-medium text-sm">Add Carousel Slide</h3>
      <div class="grid grid-cols-2 gap-3">
        <input v-model="newSlide.title" placeholder="Title"
          class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
        <input v-model.number="newSlide.display_order" type="number" placeholder="Display order"
          class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
      </div>
      <input v-model="newSlide.description" placeholder="Description"
        class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
      <div class="flex gap-2 items-start">
        <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" @change="onAddImageSelect"
          class="flex-1 text-sm text-white file:mr-2 file:py-1 file:px-3 file:rounded file:border-0 file:bg-primary-500 file:text-white file:text-xs file:font-semibold file:cursor-pointer" />
        <img v-if="addImagePreview" :src="addImagePreview" class="w-12 h-8 rounded object-cover shrink-0" />
      </div>
      <input v-if="!addImagePreview" v-model="newSlide.image_url" placeholder="...or paste Image URL"
        class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
      <div class="flex gap-3">
        <label class="flex items-center gap-2 text-sm text-white">
          <input v-model="newSlide.is_active" type="checkbox" class="accent-primary-500" />
          Active
        </label>
        <input v-model="newSlide.bg_color" placeholder="~or bg color (e.g. #FF6B00)"
          class="flex-1 px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
      </div>
      <button type="submit" class="w-full py-2 bg-primary-500 text-white rounded hover:bg-primary-600 text-sm font-semibold">
        Add Slide
      </button>
    </form>

    <div class="space-y-3 max-w-lg">
      <div v-for="slide in slides" :key="slide.id"
        class="bg-dark-700 px-4 py-3 rounded flex items-center justify-between gap-4">
        <div class="flex items-center gap-3 flex-1 min-w-0">
          <img v-if="slide.image_url" :src="slide.image_url" class="w-12 h-8 rounded object-cover shrink-0" />
          <div class="min-w-0">
            <div class="flex items-center gap-2">
              <span class="text-white font-medium text-sm">{{ slide.title }}</span>
              <span :class="['text-xs px-1.5 py-0.5 rounded', slide.is_active ? 'bg-green-900 text-green-300' : 'bg-red-900 text-red-300']">
                {{ slide.is_active ? 'Active' : 'Inactive' }}
              </span>
            </div>
            <div class="text-dark-400 text-xs truncate">
              {{ slide.description || 'No description' }} &middot; Order: {{ slide.display_order }}
            </div>
          </div>
        </div>
        <div class="flex gap-2 shrink-0">
          <button @click="startEditSlide(slide)" class="text-blue-400 text-sm hover:text-blue-300">Edit</button>
          <button @click="deleteSlide(slide.id)" class="text-red-400 text-sm hover:text-red-300">Delete</button>
        </div>
      </div>
      <p v-if="slides.length === 0" class="text-dark-400 text-sm">No carousel slides yet.</p>
    </div>

    <Teleport to="body">
      <div v-if="editModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/60" @click.self="editModal = false">
        <div class="bg-dark-800 rounded-lg p-6 w-full max-w-md mx-4 space-y-3">
          <h3 class="text-white font-semibold">Edit Carousel Slide</h3>
          <div class="grid grid-cols-2 gap-3">
            <input v-model="editForm.title" placeholder="Title"
              class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
            <input v-model.number="editForm.display_order" type="number" placeholder="Display order"
              class="px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
          </div>
          <input v-model="editForm.description" placeholder="Description"
            class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
          <div class="flex gap-2 items-start">
            <input type="file" accept="image/png,image/jpeg,image/gif,image/webp" @change="onEditImageSelect"
              class="flex-1 text-sm text-white file:mr-2 file:py-1 file:px-3 file:rounded file:border-0 file:bg-primary-500 file:text-white file:text-xs file:font-semibold file:cursor-pointer" />
            <img v-if="editImagePreview" :src="editImagePreview" class="w-12 h-8 rounded object-cover shrink-0" />
          </div>
          <input v-if="!editFileSelected" v-model="editForm.image_url" placeholder="Image URL"
            class="w-full px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
          <div class="flex gap-3">
            <label class="flex items-center gap-2 text-sm text-white">
              <input v-model="editForm.is_active" type="checkbox" class="accent-primary-500" />
              Active
            </label>
            <input v-model="editForm.bg_color" placeholder="~or bg color"
              class="flex-1 px-3 py-2 rounded bg-dark-700 text-white border border-dark-500 focus:border-primary-500 outline-none text-sm" />
          </div>
          <div class="flex gap-2 pt-2">
            <button @click="saveSlide" class="flex-1 py-2 bg-primary-500 text-white rounded hover:bg-primary-600 text-sm font-semibold">Save</button>
            <button @click="editModal = false" class="flex-1 py-2 bg-dark-500 text-white rounded hover:bg-dark-400 text-sm">Cancel</button>
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
  try {
    const res = await $api('/api/admin/carousel/list')
    slides.value = (res as any).data || []
  } catch (e: any) {
    toast.show(e?.data?.message || 'Failed to load slides')
  }
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
