<template>
  <div class="min-h-screen bg-[var(--color-background)]">
    <header class="bg-[var(--color-surface)] border-b border-[var(--color-border)] px-4 py-3 sticky top-0 z-10">
      <div class="max-w-3xl mx-auto flex items-center justify-between">
        <h1 class="font-heading text-page-title text-[var(--color-text-primary)]">Order</h1>
        <button @click="showCart = !showCart" class="relative p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors">
          <svg class="w-5 h-5 text-[var(--color-text-primary)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 100 4 2 2 0 000-4z" />
          </svg>
          <span v-if="cart.length" class="absolute -top-0.5 -right-0.5 w-5 h-5 rounded-full bg-[var(--color-primary)] text-white text-[10px] font-bold flex items-center justify-center">
            {{ cart.length }}
          </span>
        </button>
      </div>
    </header>

    <main class="max-w-3xl mx-auto p-4 space-y-6">
      <div v-if="loading">
        <AppSkeleton v-for="i in 3" :key="i" variant="card" height="120px" />
      </div>

      <section v-for="cat in categories" :key="cat.id" data-aos="fade-up">
        <h2 class="font-heading text-card-title text-[var(--color-text-primary)] border-b border-[var(--color-border)] pb-1 mb-3">{{ cat.name }}</h2>
        <div class="grid gap-3">
          <div v-for="item in getItemsByCategory(cat.id)" :key="item.id"
            class="bg-[var(--color-surface)] rounded-[var(--radius-card)] p-4 flex items-start justify-between shadow-[var(--shadow-sm)]"
          >
            <div>
              <h3 class="font-semibold text-[var(--color-text-primary)]">{{ item.name }}</h3>
              <p v-if="item.description" class="text-[var(--color-text-secondary)] text-[var(--font-size-small)]">{{ item.description }}</p>
              <span class="text-[var(--color-primary)] font-bold text-[var(--font-size-small)]">Rp{{ formatPrice(item.price) }}</span>
            </div>
            <button @click="addToCart(item)"
              class="px-4 py-1.5 rounded-[var(--radius-button)] bg-[var(--color-primary)] text-white text-sm font-semibold hover:bg-[#a84e31] transition-all duration-[var(--transition-fast)] hover:scale-[1.03]">
              + Add
            </button>
          </div>
        </div>
      </section>
    </main>

    <!-- Cart Drawer -->
    <Transition name="slide-up">
      <div v-if="showCart" class="fixed inset-0 z-50">
        <div class="absolute inset-0 bg-black/40" @click="showCart = false" />
        <div class="absolute right-0 top-0 h-full w-full max-w-md bg-[var(--color-surface)] p-6 overflow-y-auto">
          <div class="flex items-center justify-between mb-6">
            <h2 class="font-heading text-section text-[var(--color-text-primary)]">Your Cart</h2>
            <button @click="showCart = false" class="p-2 rounded-full hover:bg-[var(--color-surface-secondary)] text-[var(--color-text-secondary)]">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <AppEmptyState
            v-if="cart.length === 0"
            title="Cart is Empty"
            description="Add some items from the menu to get started."
          />

          <div v-for="(item, i) in cart" :key="i" class="flex items-center justify-between py-3 border-b border-[var(--color-border)]">
            <div>
              <span class="text-[var(--color-text-primary)] font-medium">{{ item.name }}</span>
              <span class="text-[var(--color-text-secondary)] text-sm ml-2">x{{ item.qty }}</span>
            </div>
            <div class="flex items-center gap-3">
              <span class="text-[var(--color-primary)] font-semibold">Rp{{ formatPrice(item.price * item.qty) }}</span>
              <button @click="removeFromCart(i)" class="text-[var(--color-danger)] text-sm font-medium hover:underline">Remove</button>
            </div>
          </div>

          <div v-if="cart.length" class="mt-6 space-y-4">
            <p class="text-[var(--color-text-primary)] font-bold text-xl">Total: Rp{{ formatPrice(cartTotal) }}</p>
            <AppButton variant="primary" size="lg" class="w-full" @click="submitOrder">
              Submit Order
            </AppButton>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
interface Category { id: string; name: string }
interface MenuItem { id: string; category_id: string; name: string; description: string; price: number }
interface CartItem { id: string; name: string; price: number; qty: number }

const route = useRoute()
const chainId = (route.query.chain_id as string) || (route.query.chain as string) || 'demo'
const branchId = route.query.branch as string || 'demo'
const tableId = route.query.table as string || ''

const loading = ref(true)
const categories = ref<Category[]>([])
const items = ref<MenuItem[]>([])
const cart = ref<CartItem[]>([])
const showCart = ref(false)

async function fetchData() {
  loading.value = true
  const [catRes, itemRes] = await Promise.all([
    $fetch(`/api/categories/list?chain_id=${chainId}`),
    $fetch(`/api/menu-items/list?chain_id=${chainId}`),
  ])
  categories.value = (catRes as any).data || []
  items.value = (itemRes as any).data || []
  loading.value = false
}

function getItemsByCategory(categoryId: string) {
  return items.value.filter(i => i.category_id === categoryId)
}

function addToCart(item: MenuItem) {
  const existing = cart.value.find(c => c.id === item.id)
  if (existing) {
    existing.qty++
  } else {
    cart.value.push({ id: item.id, name: item.name, price: item.price, qty: 1 })
  }
}

function removeFromCart(index: number) {
  cart.value.splice(index, 1)
}

const cartTotal = computed(() => cart.value.reduce((sum, c) => sum + c.price * c.qty, 0))

async function submitOrder() {
  await $fetch('/api/orders', {
    method: 'POST',
    body: {
      chain_id: chainId,
      branch_id: branchId,
      table_id: tableId || undefined,
      order_type: tableId ? 'dine-in' : 'takeaway',
      items: cart.value.map(c => ({ menu_item_id: c.id, quantity: c.qty, unit_price: c.price })),
    },
  })
  cart.value = []
  showCart.value = false
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

fetchData()
</script>

<style scoped>
.slide-up-enter-active,
.slide-up-leave-active { transition: all 0.3s ease-out; }
.slide-up-enter-from { opacity: 0; }
.slide-up-enter-from > div:last-child { transform: translateY(100%); }
.slide-up-leave-to { opacity: 0; }
.slide-up-leave-to > div:last-child { transform: translateY(100%); }
</style>
