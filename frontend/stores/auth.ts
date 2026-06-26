import { defineStore } from 'pinia'

interface User {
  id: string
  email: string
  name: string
  role: 'admin' | 'cashier' | 'customer'
  chainId: string
}

function parseJwt(token: string) {
  try {
    const base64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    return JSON.parse(atob(base64))
  } catch {
    return null
  }
}

function isExpired(token: string): boolean {
  const payload = parseJwt(token)
  if (!payload || !payload.exp) return true
  return Date.now() / 1000 > payload.exp
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const restoring = ref(true) // true during initial session restore

  // --- Hydrate from cookie on both server AND client ---
  function hydrateFromCookie() {
    const cookie = useCookie('auth_token')
    const token = cookie.value
    if (!token || isExpired(token)) {
      if (token) cookie.value = null // clear expired
      return
    }
    accessToken.value = token
    const payload = parseJwt(token)
    if (payload) {
      user.value = {
        id: payload.sub,
        email: payload.email || '',
        name: payload.name || payload.email || '',
        role: payload.role,
        chainId: payload.chain,
      }
    }
  }

  // SSR: hydrate immediately
  if (import.meta.server) {
    hydrateFromCookie()
  }

  // Client: hydrate from Pinia SSR state first, then fallback to cookie
  // (SSR state may be lost on some navigations; cookie is the backup)
  if (import.meta.client && !accessToken.value) {
    hydrateFromCookie()
  }

  const isAuthenticated = computed(() => !!accessToken.value && !isExpired(accessToken.value))
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isCashier = computed(() => user.value?.role === 'cashier')

  async function login(email: string, password: string) {
    const res = await $fetch('/api/auth/login', {
      method: 'POST',
      body: { email, password },
    })

    const data = res as { data: { access_token: string; refresh_token: string } }
    accessToken.value = data.data.access_token
    refreshToken.value = data.data.refresh_token

    const payload = parseJwt(accessToken.value)
    user.value = payload
      ? { id: payload.sub, email: payload.email || '', name: payload.name || payload.email || '', role: payload.role, chainId: payload.chain }
      : null

    const cookie = useCookie('auth_token', {
      maxAge: 60 * 15,
      sameSite: 'lax',
      path: '/',
    })
    cookie.value = accessToken.value
  }

  function setSession(token: string, userData: User) {
    accessToken.value = token
    user.value = userData
    const cookie = useCookie('auth_token', {
      maxAge: 60 * 15,
      sameSite: 'lax',
      path: '/',
    })
    cookie.value = token
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    const cookie = useCookie('auth_token')
    cookie.value = null
    restoring.value = false
  }

  return {
    user, accessToken, refreshToken, restoring,
    isAuthenticated, isAdmin, isCashier,
    login, logout, setSession, hydrateFromCookie,
  }
})
