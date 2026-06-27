export interface CartItem {
  id: string
  name: string
  price: number
  quantity: number
  image_url?: string
  variant_id?: string
  variant_name?: string
}

const items = ref<CartItem[]>([])

const total = computed(() =>
  items.value.reduce((sum, item) => sum + item.price * item.quantity, 0)
)

const itemCount = computed(() =>
  items.value.reduce((sum, item) => sum + item.quantity, 0)
)

export function useCart() {
  function addItem(menuItem: { id: string; name: string; price: number; image_url?: string }, opts?: { variant_id?: string; variant_name?: string }) {
    const existing = items.value.find(i => i.id === menuItem.id && i.variant_id === (opts?.variant_id || undefined))
    if (existing) {
      existing.quantity++
    } else {
      items.value.push({
        id: menuItem.id,
        name: menuItem.name + (opts?.variant_name ? ` (${opts.variant_name})` : ''),
        price: menuItem.price,
        quantity: 1,
        image_url: menuItem.image_url,
        variant_id: opts?.variant_id,
        variant_name: opts?.variant_name,
      })
    }
  }

  function updateQuantity(itemId: string, qty: number) {
    if (qty <= 0) {
      items.value = items.value.filter(i => i.id !== itemId)
    } else {
      const item = items.value.find(i => i.id === itemId)
      if (item) item.quantity = qty
    }
  }

  function clear() {
    items.value = []
  }

  return { items, total, itemCount, addItem, updateQuantity, clear }
}
