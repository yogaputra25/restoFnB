<template>
  <div class="bg-white flex flex-col h-full">
    <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200">
      <h2 class="font-bold text-lg text-gray-800">Pesanan</h2>
      <button v-if="showClose" @click="$emit('close')" class="lg:hidden text-gray-500 hover:text-gray-700">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <div class="flex-1 overflow-y-auto px-4 py-3 space-y-3">
      <div v-if="items.length === 0" class="text-center text-gray-400 py-8">
        <div class="text-5xl mb-3">🛒</div>
        <p class="text-sm">Belum ada item</p>
        <p class="text-xs mt-1">Tambahkan menu dari daftar</p>
      </div>

      <div v-for="item in items" :key="item.id"
        class="flex items-center gap-3 bg-gray-50 rounded-xl p-3">
        <div v-if="item.image_url" class="w-12 h-12 rounded-lg overflow-hidden flex-shrink-0 bg-gray-200">
          <img :src="item.image_url" :alt="item.name" class="w-full h-full object-cover" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="font-medium text-sm text-gray-800 truncate">{{ item.name }}</p>
          <p class="text-xs text-gray-500 mt-0.5">Rp{{ formatPrice(item.price) }}</p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="cart.updateQuantity(item.id, item.quantity - 1)"
            class="w-7 h-7 rounded-full bg-gray-200 flex items-center justify-center text-gray-600 hover:bg-gray-300 transition-colors text-sm font-bold">
            -
          </button>
          <span class="w-6 text-center font-semibold text-sm">{{ item.quantity }}</span>
          <button @click="cart.updateQuantity(item.id, item.quantity + 1)"
            class="w-7 h-7 rounded-full bg-primary-500 flex items-center justify-center text-white hover:bg-primary-600 transition-colors text-sm font-bold">
            +
          </button>
        </div>
      </div>
    </div>

    <div class="border-t border-gray-200 p-4 space-y-3">
      <div class="flex justify-between items-center text-lg font-bold text-gray-800">
        <span>Total</span>
        <span class="text-primary-600">Rp{{ formatPrice(cart.total) }}</span>
      </div>

      <div v-if="items.length > 0" class="space-y-2">
        <div>
          <label class="text-xs font-medium text-gray-600 block mb-1">Tipe Pesanan</label>
          <div class="flex gap-2">
            <button @click="orderType = 'dine-in'"
              :class="['flex-1 py-2 px-3 rounded-lg text-sm font-medium border transition-colors',
                orderType === 'dine-in' ? 'bg-primary-500 text-white border-primary-500' : 'bg-white text-gray-600 border-gray-300 hover:border-gray-400']">
              Makan di Tempat
            </button>
            <button @click="orderType = 'takeaway'"
              :class="['flex-1 py-2 px-3 rounded-lg text-sm font-medium border transition-colors',
                orderType === 'takeaway' ? 'bg-primary-500 text-white border-primary-500' : 'bg-white text-gray-600 border-gray-300 hover:border-gray-400']">
              Bungkus
            </button>
          </div>
        </div>

        <div v-if="orderType === 'dine-in'">
          <label class="text-xs font-medium text-gray-600 block mb-1">Nomor Meja</label>
          <input v-model="tableId"
            :placeholder="tableIdFromUrl || 'Masukkan nomor meja'"
            class="w-full px-3 py-2 rounded-lg border border-gray-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
          <p v-if="tableIdFromUrl && !tableId" class="text-xs text-gray-400 mt-1">
            Terdeteksi dari QR: meja {{ tableIdFromUrl }}
          </p>
        </div>

        <div>
          <label class="text-xs font-medium text-gray-600 block mb-1">Nama Pemesan</label>
          <input v-model="customerName"
            placeholder="Masukkan nama Anda (opsional)"
            class="w-full px-3 py-2 rounded-lg border border-gray-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
        </div>

        <div>
          <label class="text-xs font-medium text-gray-600 block mb-1">ID Cabang (Branch)</label>
          <input v-model="branchId"
            :placeholder="branchIdFromUrl || 'Masukkan ID cabang'"
            class="w-full px-3 py-2 rounded-lg border border-gray-300 text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
          />
          <p v-if="branchIdFromUrl && !branchId" class="text-xs text-gray-400 mt-1">
            Terdeteksi dari URL: {{ branchIdFromUrl }}
          </p>
        </div>
      </div>

      <button @click="submitOrder"
        :disabled="submitting || items.length === 0"
        :class="['w-full py-3 rounded-xl font-bold text-sm transition-colors',
          submitting ? 'bg-gray-300 text-gray-500 cursor-not-allowed' :
          items.length === 0 ? 'bg-gray-200 text-gray-400 cursor-not-allowed' :
          'bg-primary-500 text-white hover:bg-primary-600 active:bg-primary-700']">
        <span v-if="submitting">Memproses...</span>
        <span v-else>Pesan Sekarang</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { unref } from 'vue'
import { useCart } from '~/composables/useCart'

defineProps<{ showClose?: boolean }>()
defineEmits<{ close: [] }>()

const cart = useCart()
const { items, total } = cart
const route = useRoute()
const toast = useToast()

const chainId = (route.query.chain_id as string) || ''
const tableIdFromUrl = (route.query.table_id as string) || ''
const branchIdFromUrl = (route.query.branch_id as string) || ''

const tableId = ref(tableIdFromUrl)
const branchId = ref(branchIdFromUrl)
const orderType = ref(tableIdFromUrl ? 'dine-in' : 'takeaway')
const customerName = ref('')
const submitting = ref(false)

function formatPrice(price: number) {
  return unref(price).toLocaleString('id-ID')
}

async function submitOrder() {
  if (items.value.length === 0) return
  submitting.value = true

  const actualTableId = tableId.value || tableIdFromUrl || undefined
  const actualBranchId = branchId.value || branchIdFromUrl

  if (!actualBranchId) {
    toast.show('Harap masukkan ID cabang (branch)', 'error')
    submitting.value = false
    return
  }

  if (orderType.value === 'dine-in' && !actualTableId) {
    toast.show('Harap masukkan nomor meja', 'error')
    submitting.value = false
    return
  }

  const payload: Record<string, any> = {
    chain_id: chainId,
    branch_id: actualBranchId,
    order_type: orderType.value,
    customer_name: customerName.value || undefined,
    items: items.value.map(item => ({
      menu_item_id: item.id,
      quantity: item.quantity,
      unit_price: item.price,
    })),
  }

  if (actualTableId) {
    payload.table_id = actualTableId
  }

  try {
    const res = await $fetch('/api/orders', {
      method: 'POST',
      body: payload,
    })
    const order = res as any
    toast.show(`Pesanan berhasil! No. Order: ${order.id?.slice(0, 8) || ''}`, 'success')
    cart.clear()
  } catch (err: any) {
    const msg = err?.data?.message || err?.response?._data?.message || 'Gagal memproses pesanan'
    toast.show(msg, 'error')
  } finally {
    submitting.value = false
  }
}
</script>
