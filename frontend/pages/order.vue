<template>
  <div class="min-h-screen bg-dark-50">
    <header class="bg-dark-900 text-white px-6 py-4 sticky top-0 z-10">
      <div class="max-w-3xl mx-auto flex items-center justify-between">
        <h1 class="text-xl font-bold text-primary-500">Order</h1>
        <button @click="showCart = !showCart" class="relative">
          <span class="text-white">Cart</span>
          <span v-if="cart.length" class="absolute -top-2 -right-4 bg-primary-500 text-xs rounded-full w-5 h-5 flex items-center justify-center">
            {{ cart.length }}
          </span>
        </button>
      </div>
    </header>

    <main class="max-w-3xl mx-auto p-4 space-y-6">
      <section v-for="cat in categories" :key="cat.id">
        <h2 class="text-lg font-bold text-dark-800 border-b border-dark-300 pb-1 mb-3">{{ cat.name }}</h2>
        <div class="grid gap-3">
          <div v-for="item in getItemsByCategory(cat.id)" :key="item.id"
            class="bg-white rounded-lg shadow p-4 flex items-start justify-between">
            <div>
              <h3 class="font-semibold text-dark-800">{{ item.name }}</h3>
              <p v-if="item.description" class="text-dark-400 text-sm">{{ item.description }}</p>
              <span class="text-primary-600 font-bold text-sm">Rp{{ formatPrice(item.price) }}</span>
            </div>
            <button @click="addToCart(item)" class="px-3 py-1 bg-primary-500 text-white rounded text-sm hover:bg-primary-600">
              + Add
            </button>
          </div>
        </div>
      </section>
    </main>

    <div v-if="showCart" class="fixed inset-0 bg-black/50 z-20" @click="showCart = false">
      <div class="absolute right-0 top-0 h-full w-full max-w-md bg-dark-800 p-6 overflow-y-auto" @click.stop>
        <h2 class="text-xl font-bold text-white mb-4">Your Cart</h2>
        <div v-if="cart.length === 0" class="text-dark-400">Cart is empty</div>
        <div v-for="(item, i) in cart" :key="i" class="flex items-center justify-between py-2 border-b border-dark-600">
          <div>
            <span class="text-white">{{ item.name }}</span>
            <span class="text-dark-400 text-sm ml-2">x{{ item.qty }}</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-primary-500">Rp{{ formatPrice(item.price * item.qty) }}</span>
            <button @click="removeFromCart(i)" class="text-red-400 text-sm">Remove</button>
          </div>
        </div>
        <div v-if="cart.length" class="mt-4">
          <p class="text-white font-bold">Total: Rp{{ formatPrice(cartTotal) }}</p>
          <button @click="submitOrder" class="w-full mt-4 py-3 bg-primary-500 text-white rounded-lg font-bold hover:bg-primary-600">
            Submit Order
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Category { id: string; name: string }
interface MenuItem { id: string; category_id: string; name: string; description: string; price: number }
interface CartItem { id: string; name: string; price: number; qty: number }

const route = useRoute()
const branchId = route.query.branch as string || 'demo'
const tableId = route.query.table as string || ''

const categories = ref<Category[]>([])
const items = ref<MenuItem[]>([])
const cart = ref<CartItem[]>([])
const showCart = ref(false)

async function fetchData() {
  const catRes = await $fetch('/api/admin/categories/list?chain_id=demo')
  categories.value = (catRes as any).data || []
  const itemRes = await $fetch('/api/admin/menu-items/list?chain_id=demo')
  items.value = (itemRes as any).data || []
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
      branch_id: branchId,
      table_id: tableId || undefined,
      order_type: tableId ? 'dine-in' : 'takeaway',
      items: cart.value.map(c => ({ menu_item_id: c.id, quantity: c.qty, unit_price: c.price })),
    },
  })
  cart.value = []
  showCart.value = false
  alert('Order submitted!')
}

function formatPrice(price: number) {
  return price.toLocaleString('id-ID')
}

fetchData()
</script>
