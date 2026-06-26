<template>
  <header class="h-[72px] bg-[var(--color-surface)] border-b border-[var(--color-border)] sticky top-0 z-40 shadow-[var(--shadow-sm)]">
    <div class="flex items-center justify-between h-full px-6">
      <!-- Left: hamburger + breadcrumb -->
      <div class="flex items-center gap-4">
        <button class="lg:hidden p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors" @click="$emit('toggleMobile')">
          <svg class="w-5 h-5 text-[var(--color-text-primary)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
        </button>
        <div class="hidden sm:flex items-center gap-2 text-sm text-[var(--color-text-secondary)]">
          <span v-for="(crumb, i) in breadcrumbs" :key="i">
            <RouterLink v-if="crumb.to" :to="crumb.to" class="hover:text-[var(--color-primary)] transition-colors">{{ crumb.label }}</RouterLink>
            <span v-else class="text-[var(--color-text-primary)] font-medium">{{ crumb.label }}</span>
            <span v-if="i < breadcrumbs.length - 1" class="mx-1">/</span>
          </span>
        </div>
      </div>

      <!-- Right: branch, search, notifications, profile -->
      <div class="flex items-center gap-3">
        <!-- Current Date -->
        <span class="hidden xl:block text-xs text-[var(--color-text-secondary)] font-mono">{{ currentDate }}</span>

        <!-- Branch selector (UI only) -->
        <div class="hidden md:flex items-center gap-1.5 text-sm text-[var(--color-text-secondary)] bg-[var(--color-surface-secondary)] px-3 py-1.5 rounded-[var(--radius-button)]">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
          </svg>
          <span>All Branches</span>
        </div>

        <!-- Search (UI only) -->
        <button class="p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors">
          <svg class="w-5 h-5 text-[var(--color-text-secondary)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>

        <!-- Notifications (UI only) -->
        <button class="relative p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors">
          <svg class="w-5 h-5 text-[var(--color-text-secondary)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
          </svg>
          <span class="absolute top-1 right-1 w-2 h-2 rounded-full bg-[var(--color-danger)]" />
        </button>

        <!-- Dark mode toggle (UI only) -->
        <button class="hidden md:flex p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors">
          <svg class="w-5 h-5 text-[var(--color-text-secondary)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        </button>

        <!-- Profile dropdown -->
        <div class="relative" @click.outside="profileOpen = false">
          <button class="flex items-center gap-2 p-1.5 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors" @click="profileOpen = !profileOpen">
            <div class="w-8 h-8 rounded-full bg-[var(--color-primary)] flex items-center justify-center text-white text-sm font-bold">
              {{ initials }}
            </div>
            <div class="hidden lg:block text-left text-sm">
              <p class="text-[var(--color-text-primary)] font-medium leading-tight">{{ auth.user?.name || 'User' }}</p>
              <p class="text-[var(--color-text-secondary)] text-xs capitalize">{{ auth.user?.role || '' }}</p>
            </div>
          </button>

          <Transition name="fade-down">
            <div v-if="profileOpen" class="absolute right-0 top-full mt-2 w-48 bg-[var(--color-surface)] rounded-[var(--radius-card)] shadow-[var(--shadow-lg)] border border-[var(--color-border)] py-1 overflow-hidden">
              <div class="px-4 py-2 border-b border-[var(--color-border)]">
                <p class="text-sm font-medium text-[var(--color-text-primary)]">{{ auth.user?.name }}</p>
                <p class="text-xs text-[var(--color-text-secondary)]">{{ auth.user?.email }}</p>
              </div>
              <button class="w-full text-left px-4 py-2 text-sm text-[var(--color-text-secondary)] hover:bg-[var(--color-surface-secondary)] transition-colors flex items-center gap-2" @click="handleLogout">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" /></svg>
                Logout
              </button>
            </div>
          </Transition>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()
const route = useRoute()

const profileOpen = ref(false)

const initials = computed(() => {
  if (!auth.user?.name) return '?'
  return auth.user.name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
})

const breadcrumbs = computed(() => {
  const parts = route.path.split('/').filter(Boolean)
  const crumbs: { label: string; to?: string }[] = []
  let path = ''
  for (const part of parts) {
    path += `/${part}`
    crumbs.push({
      label: part.charAt(0).toUpperCase() + part.slice(1),
      to: path,
    })
  }
  return crumbs
})

const currentDate = ref('')
onMounted(() => {
  const d = new Date()
  currentDate.value = d.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric', year: 'numeric' })
})

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.fade-down-enter-active,
.fade-down-leave-active { transition: all 0.15s ease; }
.fade-down-enter-from,
.fade-down-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
