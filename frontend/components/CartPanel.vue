<template>
  <div class="bg-[var(--color-surface)] flex flex-col h-full">
    <div class="flex items-center justify-between px-4 py-3 border-b border-[var(--color-border)]">
      <h2 class="font-heading text-card-title text-[var(--color-text-primary)]">Pesanan</h2>
      <button v-if="showClose" @click="$emit('close')" class="lg:hidden p-2 rounded-full hover:bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)] transition-colors">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <div class="flex-1 overflow-y-auto px-4 py-3 space-y-3">
      <AppEmptyState
        v-if="items.length === 0"
        title="Belum Ada Item"
        description="Tambahkan menu dari daftar"
      />

      <div v-for="item in items" :key="item.id"
        class="flex items-center gap-3 bg-[var(--color-surface-secondary)] rounded-[var(--radius-card)] p-3"
      >
        <div v-if="item.image_url" class="w-12 h-12 rounded-[var(--radius-image)] overflow-hidden flex-shrink-0 bg-[var(--color-border)]">
          <img :src="item.image_url" :alt="item.name" class="w-full h-full object-cover" loading="lazy" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="font-medium text-sm text-[var(--color-text-primary)] truncate">{{ item.name }}</p>
          <p class="text-xs text-[var(--color-text-secondary)] mt-0.5">Rp{{ formatPrice(item.price) }}</p>
        </div>
        <div class="flex items-center gap-2 flex-shrink-0">
          <button @click="cart.updateQuantity(item.id, item.quantity - 1)"
            class="w-7 h-7 rounded-full bg-[var(--color-border)] flex items-center justify-center text-[var(--color-text-secondary)] hover:bg-[var(--color-text-secondary)] hover:text-white transition-colors text-sm font-bold">
            -
          </button>
          <span class="w-6 text-center font-semibold text-sm text-[var(--color-text-primary)]">{{ item.quantity }}</span>
          <button @click="cart.updateQuantity(item.id, item.quantity + 1)"
            class="w-7 h-7 rounded-full bg-[var(--color-primary)] flex items-center justify-center text-white hover:bg-[#a84e31] transition-colors text-sm font-bold">
            +
          </button>
        </div>
      </div>
    </div>

    <div class="border-t border-[var(--color-border)] p-4 space-y-3">
      <div class="flex justify-between items-center text-lg font-bold text-[var(--color-text-primary)]">
        <span>Total</span>
        <span class="text-[var(--color-primary)]">Rp{{ formatPrice(cart.total) }}</span>
      </div>

      <div v-if="items.length > 0" class="space-y-2">
        <div>
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">Tipe Pesanan</label>
          <div class="flex gap-2">
            <button @click="orderType = 'dine-in'"
              :class="['flex-1 py-2 px-3 rounded-[var(--radius-button)] text-sm font-medium border transition-all duration-[var(--transition-fast)]',
                orderType === 'dine-in' ? 'bg-[var(--color-primary)] text-white border-[var(--color-primary)]' : 'bg-white text-[var(--color-text-secondary)] border-[var(--color-border)] hover:border-[var(--color-primary)]']">
              Makan di Tempat
            </button>
            <button @click="orderType = 'takeaway'"
              :class="['flex-1 py-2 px-3 rounded-[var(--radius-button)] text-sm font-medium border transition-all duration-[var(--transition-fast)]',
                orderType === 'takeaway' ? 'bg-[var(--color-primary)] text-white border-[var(--color-primary)]' : 'bg-white text-[var(--color-text-secondary)] border-[var(--color-border)] hover:border-[var(--color-primary)]']">
              Bungkus
            </button>
          </div>
        </div>

        <div v-if="orderType === 'dine-in'">
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">Nomor Meja</label>
          <AppInput v-model="tableId" :placeholder="tableIdFromUrl || 'Masukkan nomor meja'" />
          <p v-if="tableIdFromUrl && !tableId" class="text-xs text-[var(--color-text-secondary)] mt-1">
            Terdeteksi dari QR: meja {{ tableIdFromUrl }}
          </p>
        </div>

        <div>
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">Nama Pemesan</label>
          <AppInput v-model="customerName" placeholder="Masukkan nama Anda (opsional)" />
        </div>

        <div v-if="branchIdFromUrl">
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">Cabang</label>
          <div class="px-3 py-2 rounded-[var(--radius-input)] bg-green-50 border border-green-200 text-green-700 text-sm">
            {{ branchIdFromUrl }}
          </div>
        </div>
        <div v-else>
          <label class="text-xs font-medium text-[var(--color-text-secondary)] block mb-1">ID Cabang (Branch)</label>
          <AppInput v-model="branchId" placeholder="Masukkan ID cabang" />
        </div>
      </div>

      <AppButton
        variant="primary"
        size="lg"
        class="w-full"
        :loading="submitting"
        :disabled="items.length === 0"
        @click="submitOrder"
      >
        {{ submitting ? 'Memproses...' : 'Pesan Sekarang' }}
      </AppButton>
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
const router = useRouter()
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
    const data = res as any
    const order = data?.data || data
    order.items = order.items.map((item: any, i: number) => ({
      ...item,
      name: items.value.find(c => c.id === item.menu_item_id)?.name || `Item ${i + 1}`,
    }))
    const chainParam = order.chain_id || chainId
    cart.clear()
    router.push({ path: '/order-success', state: { order }, query: { chain_id: chainParam } })
  } catch (err: any) {
    const msg = err?.data?.message || err?.response?._data?.message || 'Gagal memproses pesanan'
    toast.show(msg, 'error')
  } finally {
    submitting.value = false
  }
}
</script>
