<template>
  <header class="h-[72px] bg-[var(--color-surface)] border-b border-[var(--color-border)] sticky top-0 z-40 shadow-[var(--shadow-sm)]">
    <div class="flex items-center justify-between h-full px-6">
      <!-- Left: hamburger + breadcrumb -->
      <div class="flex items-center gap-4">
        <button class="lg:hidden p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors" @click="$emit('toggleMobile')">
          <Menu class="w-5 h-5 text-[var(--color-text-primary)]" />
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

        <!-- Branch selector (admin only — cashier auto-assigned) -->
        <div v-if="auth.isAdmin" class="hidden md:flex items-center gap-1.5 text-sm bg-[var(--color-surface-secondary)] px-3 py-1.5 rounded-[var(--radius-button)]">
          <Building2 class="w-4 h-4 text-[var(--color-text-secondary)]" />
          <select
            :value="selectedBranchId || ''"
            @change="setBranch(($event.target as HTMLSelectElement).value || null)"
            class="bg-transparent text-[var(--color-text-primary)] outline-none text-sm max-w-[160px] cursor-pointer"
          >
            <option value="">All Branches</option>
            <option v-for="b in branches" :key="b.id" :value="b.id">{{ b.name }}</option>
          </select>
        </div>
        <!-- Branch label for cashier (read-only) -->
        <div v-else-if="auth.isCashier && auth.user?.branchId" class="hidden md:flex items-center gap-1.5 text-sm text-[var(--color-text-secondary)] bg-[var(--color-surface-secondary)] px-3 py-1.5 rounded-[var(--radius-button)]">
          <Building2 class="w-4 h-4" />
          <span>{{ getBranchName(auth.user.branchId) }}</span>
        </div>

        <!-- Search (UI only) -->
        <button class="p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors">
          <Search class="w-5 h-5 text-[var(--color-text-secondary)]" />
        </button>

        <!-- Notifications (UI only) -->
        <button class="relative p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors">
          <Bell class="w-5 h-5 text-[var(--color-text-secondary)]" />
          <span class="absolute top-1 right-1 w-2 h-2 rounded-full bg-[var(--color-danger)]" />
        </button>

        <!-- Dark mode toggle (UI only) -->
        <button class="hidden md:flex p-2 rounded-full hover:bg-[var(--color-surface-secondary)] transition-colors">
          <Moon class="w-5 h-5 text-[var(--color-text-secondary)]" />
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
                <LogOut class="w-4 h-4" />
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
import { Menu, Building2, Search, Bell, Moon, LogOut } from 'lucide-vue-next'

const auth = useAuthStore()
const router = useRouter()
const route = useRoute()
const { selectedBranchId, branches, setBranch, fetchBranches } = useBranchSelector()

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
  fetchBranches()
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
