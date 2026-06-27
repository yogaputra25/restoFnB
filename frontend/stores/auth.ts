import { defineStore } from 'pinia'

interface User {
  id: string
  email: string
  name: string
  role: 'admin' | 'cashier' | 'customer'
  chainId: string
  branchId?: string
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

  const isAuthenticated = computed(() => !!accessToken.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isCashier = computed(() => user.value?.role === 'cashier')

  function persist() {
    if (import.meta.client) {
      try {
        if (accessToken.value) {
          localStorage.setItem('auth_token', accessToken.value)
          localStorage.setItem('refresh_token', refreshToken.value || '')
          const enc = (v: string) => encodeURIComponent(v)
          document.cookie = `auth_token=${enc(accessToken.value)}; path=/; max-age=${60*60*24*7}; SameSite=Lax`
          document.cookie = `refresh_token=${enc(refreshToken.value || '')}; path=/; max-age=${60*60*24*7}; SameSite=Lax`
        } else {
          localStorage.removeItem('auth_token')
          localStorage.removeItem('refresh_token')
          document.cookie = 'auth_token=; path=/; max-age=0; SameSite=Lax'
          document.cookie = 'refresh_token=; path=/; max-age=0; SameSite=Lax'
        }
      } catch { /* ignore */ }
    }
  }

  function tryRestore() {
    if (import.meta.client && !accessToken.value) {
      try {
        const stored = localStorage.getItem('auth_token')
        if (stored && !isExpired(stored)) {
          const payload = parseJwt(stored)
          if (payload) {
            accessToken.value = stored
            user.value = {
              id: payload.sub,
              email: payload.email || '',
              name: payload.name || payload.email || '',
              role: payload.role,
              chainId: payload.chain,
              branchId: payload.branch || undefined,
            }
            const rt = localStorage.getItem('refresh_token')
            if (rt) refreshToken.value = rt
            return true
          }
        } else if (stored) {
          localStorage.removeItem('auth_token')
          localStorage.removeItem('refresh_token')
        }
      } catch { /* ignore */ }
    }
    return false
  }

  function setSession(token: string, userData: User) {
    accessToken.value = token
    user.value = userData
    persist()
  }

  async function login(email: string, password: string) {
    const res = await $fetch('/api/auth/login', {
      method: 'POST',
      body: { email, password },
    })

    const data = res as { data: { access_token: string; refresh_token: string } }
    accessToken.value = data.data.access_token
    refreshToken.value = data.data.refresh_token

    const payload = parseJwt(accessToken.value)
    if (payload) {
      user.value = {
        id: payload.sub,
        email: payload.email || '',
        name: payload.name || payload.email || '',
        role: payload.role,
        chainId: payload.chain,
        branchId: payload.branch || undefined,
      }
    }
    persist()
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    persist()
  }

  return {
    user, accessToken, refreshToken,
    isAuthenticated, isAdmin, isCashier,
    setSession, login, logout,
  }
})
