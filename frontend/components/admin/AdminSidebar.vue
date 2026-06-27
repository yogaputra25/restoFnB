<template>
  <!-- Desktop sidebar -->
  <aside
    class="hidden lg:flex flex-col bg-[var(--color-surface)] border-r border-[var(--color-border)] transition-all duration-300 overflow-hidden"
    :class="collapsed ? 'w-[72px]' : 'w-[280px]'"
  >
    <!-- Logo area -->
    <div class="flex items-center h-[72px] px-4 border-b border-[var(--color-border)] shrink-0">
      <div v-if="!collapsed" class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-[var(--radius-button)] bg-[var(--color-primary)] flex items-center justify-center text-white font-bold text-sm">RF</div>
        <div>
          <p class="font-heading font-bold text-[var(--color-text-primary)] leading-tight">RestoFnB</p>
          <p class="text-[10px] text-[var(--color-text-secondary)]">Management</p>
        </div>
      </div>
      <div v-else class="mx-auto">
        <div class="w-8 h-8 rounded-[var(--radius-button)] bg-[var(--color-primary)] flex items-center justify-center text-white font-bold text-sm">RF</div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto py-4 px-3 space-y-6">
      <div v-for="group in navGroups" :key="group.label">
        <p v-if="!collapsed" class="px-3 text-[10px] font-semibold uppercase tracking-widest text-[var(--color-text-secondary)] mb-2">{{ group.label }}</p>
        <div class="space-y-0.5">
          <NuxtLink
            v-for="item in group.items"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-3 px-3 py-2.5 rounded-[var(--radius-button)] text-sm font-medium transition-all duration-[var(--transition-fast)] group"
            :class="isActive(item.to)
              ? 'bg-[var(--color-primary)]/5 text-[var(--color-primary)]'
              : 'text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] hover:text-[var(--color-text-primary)]'"
          >
            <!-- Active accent bar -->
            <span v-if="isActive(item.to) && !collapsed" class="absolute left-0 w-0.5 h-6 rounded-r bg-[var(--color-primary)]" />

            <!-- Icon -->
            <component :is="item.icon" class="w-5 h-5 shrink-0" />

            <!-- Label -->
            <span v-if="!collapsed" class="truncate">{{ item.label }}</span>
          </NuxtLink>
        </div>
      </div>
    </nav>

    <!-- Collapse + Logout -->
    <div class="border-t border-[var(--color-border)] p-3 space-y-1 shrink-0">
      <button
        class="flex items-center gap-3 w-full px-3 py-2.5 rounded-[var(--radius-button)] text-sm font-medium text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] hover:text-[var(--color-text-primary)] transition-all"
        @click="$emit('toggleCollapse')"
      >
        <component :is="icon('ChevronLeft')" class="w-5 h-5 shrink-0" />
        <span v-if="!collapsed">Collapse</span>
      </button>
      <button
        class="flex items-center gap-3 w-full px-3 py-2.5 rounded-[var(--radius-button)] text-sm font-medium text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] hover:text-[var(--color-danger)] transition-all"
        @click="handleLogout"
      >
        <component :is="icon('LogOut')" class="w-5 h-5 shrink-0" />
        <span v-if="!collapsed">Logout</span>
      </button>
    </div>
  </aside>

  <!-- Mobile drawer -->
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="mobileOpen" class="fixed inset-0 z-50 lg:hidden">
        <div class="absolute inset-0 bg-black/40" @click="$emit('closeMobile')" />
        <aside class="absolute top-0 left-0 h-full w-[280px] bg-[var(--color-surface)] border-r border-[var(--color-border)] flex flex-col overflow-hidden">
          <!-- Mobile logo -->
          <div class="flex items-center justify-between h-[72px] px-4 border-b border-[var(--color-border)] shrink-0">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-[var(--radius-button)] bg-[var(--color-primary)] flex items-center justify-center text-white font-bold text-sm">RF</div>
              <p class="font-heading font-bold text-[var(--color-text-primary)]">RestoFnB</p>
            </div>
            <button class="p-2 rounded-full hover:bg-[var(--color-surface-secondary)]" @click="$emit('closeMobile')">
              <component :is="icon('X')" class="w-5 h-5 text-[var(--color-text-secondary)]" />
            </button>
          </div>

          <!-- Mobile nav -->
          <nav class="flex-1 overflow-y-auto py-4 px-3 space-y-6">
            <div v-for="group in navGroups" :key="group.label">
              <p class="px-3 text-[10px] font-semibold uppercase tracking-widest text-[var(--color-text-secondary)] mb-2">{{ group.label }}</p>
              <div class="space-y-0.5">
                <NuxtLink
                  v-for="item in group.items"
                  :key="item.to"
                  :to="item.to"
                  class="flex items-center gap-3 px-3 py-2.5 rounded-[var(--radius-button)] text-sm font-medium transition-all"
                  :class="isActive(item.to) ? 'bg-[var(--color-primary)]/5 text-[var(--color-primary)]' : 'text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)]'"
                  @click="$emit('closeMobile')"
                >
                  <component :is="item.icon" class="w-5 h-5" />
                  <span>{{ item.label }}</span>
                </NuxtLink>
              </div>
            </div>
          </nav>

          <!-- Mobile logout -->
          <div class="border-t border-[var(--color-border)] p-3 shrink-0">
            <button class="flex items-center gap-3 w-full px-3 py-2.5 rounded-lg text-sm font-medium text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] hover:text-[var(--color-danger)] transition-all" @click="handleLogout">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
              <span>Logout</span>
            </button>
          </div>
        </aside>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { h } from 'vue'
import * as LucideIcons from 'lucide-vue-next'

defineProps<{
  collapsed?: boolean
  mobileOpen?: boolean
}>()

defineEmits<{
  toggleCollapse: []
  closeMobile: []
}>()

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

function isActive(path: string) {
  return route.path.startsWith(path)
}

function handleLogout() {
  auth.logout()
  router.push('/login')
}

function icon(name: string) {
  const icons = LucideIcons as Record<string, any>
  return icons[`${name.charAt(0).toUpperCase()}${name.slice(1)}Icon`] || icons['CircleIcon'] || 'span'
}

interface NavItem {
  to: string
  label: string
  icon: any
}

interface NavGroup {
  label: string
  items: NavItem[]
}

const adminNavGroups: NavGroup[] = [
  {
    label: 'Dashboard',
    items: [
      { to: '/admin', label: 'Overview', icon: icon('LayoutDashboard') },
    ],
  },
  {
    label: 'Management',
    items: [
      { to: '/admin/menu', label: 'Menu', icon: icon('UtensilsCrossed') },
      { to: '/admin/categories', label: 'Categories', icon: icon('Tags') },
      { to: '/admin/orders', label: 'Orders', icon: icon('ClipboardList') },
      { to: '/admin/tables', label: 'Tables', icon: icon('Grid3x3') },
      { to: '/admin/staff', label: 'Staff', icon: icon('Users') },
    ],
  },
  {
    label: 'Restaurant',
    items: [
      { to: '/admin/branches', label: 'Branches', icon: icon('Building2') },
      { to: '/admin/carousel', label: 'Carousel', icon: icon('Images') },
      { to: '/admin/reports', label: 'Reports', icon: icon('BarChart3') },
    ],
  },
  {
    label: 'Settings',
    items: [
      { to: '/admin/settings', label: 'Settings', icon: icon('Settings') },
    ],
  },
]

const cashierNavGroups: NavGroup[] = [
  {
    label: 'Orders',
    items: [
      { to: '/cashier', label: 'Orders', icon: icon('ClipboardList') },
    ],
  },
]

const navGroups = computed(() => auth.isAdmin ? adminNavGroups : cashierNavGroups)
</script>

<style scoped>
.drawer-enter-active,
.drawer-leave-active { transition: opacity 0.2s ease; }
.drawer-enter-active > aside,
.drawer-leave-active > aside { transition: transform 0.25s ease; }
.drawer-enter-from,
.drawer-leave-to { opacity: 0; }
.drawer-enter-from > aside { transform: translateX(-100%); }
.drawer-leave-to > aside { transform: translateX(-100%); }
</style>
