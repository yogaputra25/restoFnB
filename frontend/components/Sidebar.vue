<template>
  <aside
    class="bg-dark-800 border-r border-dark-600 flex flex-col transition-all duration-300"
    :class="collapsed ? 'w-16' : 'w-60'"
  >
    <div class="flex items-center justify-between px-4 h-14 border-b border-dark-600">
      <span v-if="!collapsed" class="text-primary-500 font-bold text-lg">restoFnB</span>
      <button @click="toggleCollapse" class="text-dark-300 hover:text-white p-1">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            :d="collapsed ? 'M13 5l7 7-7 7M5 5l7 7-7 7' : 'M11 19l-7-7 7-7m8 14l-7-7 7-7'" />
        </svg>
      </button>
    </div>

    <nav class="flex-1 py-4 space-y-1 px-2">
      <NuxtLink
        v-for="item in menuItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors"
        :class="isActive(item.to) ? 'bg-primary-500/10 text-primary-500' : 'text-dark-300 hover:bg-dark-700 hover:text-white'"
      >
        <span class="w-5 h-5 flex-shrink-0" v-html="item.icon"></span>
        <span v-if="!collapsed" class="text-sm font-medium">{{ item.label }}</span>
      </NuxtLink>
    </nav>

    <div class="border-t border-dark-600 p-4">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-full bg-primary-500 flex items-center justify-center text-white text-sm font-bold">
          {{ initials }}
        </div>
        <div v-if="!collapsed" class="flex-1 min-w-0">
          <p class="text-sm text-white truncate">{{ auth.user?.name }}</p>
          <p class="text-xs text-dark-400 capitalize">{{ auth.user?.role }}</p>
        </div>
      </div>
      <button @click="handleLogout" v-if="!collapsed"
        class="mt-3 w-full text-xs text-dark-400 hover:text-red-400 text-left flex items-center gap-2">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
        Logout
      </button>
    </div>
  </aside>

  <div v-if="mobileOpen" class="fixed inset-0 bg-black/50 z-40 lg:hidden" @click="mobileOpen = false"></div>
  <aside
    v-if="mobileOpen"
    class="fixed inset-y-0 left-0 z-50 w-60 bg-dark-800 border-r border-dark-600 lg:hidden"
  >
    <div class="flex items-center justify-between px-4 h-14 border-b border-dark-600">
      <span class="text-primary-500 font-bold text-lg">restoFnB</span>
      <button @click="mobileOpen = false" class="text-dark-300 hover:text-white p-1">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
    <nav class="py-4 space-y-1 px-2">
      <NuxtLink
        v-for="item in menuItems"
        :key="item.to"
        :to="item.to"
        @click="mobileOpen = false"
        class="flex items-center gap-3 px-3 py-2.5 rounded-lg transition-colors"
        :class="isActive(item.to) ? 'bg-primary-500/10 text-primary-500' : 'text-dark-300 hover:bg-dark-700 hover:text-white'"
      >
        <span class="w-5 h-5 flex-shrink-0" v-html="item.icon"></span>
        <span class="text-sm font-medium">{{ item.label }}</span>
      </NuxtLink>
    </nav>
  </aside>
</template>

<script setup lang="ts">
const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const collapsed = ref(false)
const mobileOpen = ref(false)

function toggleCollapse() {
  if (window.innerWidth < 1024) {
    mobileOpen.value = !mobileOpen.value
  } else {
    collapsed.value = !collapsed.value
  }
}

function isActive(path: string) {
  return route.path.startsWith(path)
}

function handleLogout() {
  auth.logout()
  router.push('/login')
}

const initials = computed(() => {
  if (!auth.user?.name) return '?'
  return auth.user.name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
})

interface MenuItem {
  to: string
  label: string
  icon: string
}

const adminMenu: MenuItem[] = [
  { to: '/admin', label: 'Overview', icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" /></svg>' },
  { to: '/admin/menu', label: 'Menu', icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100-4m0 4v2m0-6V4" /></svg>' },
  { to: '/admin/branches', label: 'Branches', icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" /></svg>' },
  { to: '/admin/staff', label: 'Staff', icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" /></svg>' },
  { to: '/admin/tables', label: 'Tables', icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>' },
  { to: '/admin/orders', label: 'Orders', icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" /></svg>' },
]

const cashierMenu: MenuItem[] = [
  { to: '/cashier', label: 'Orders', icon: '<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>' },
]

const menuItems = computed(() => {
  if (auth.isAdmin) return adminMenu
  if (auth.isCashier) return cashierMenu
  return []
})
</script>
