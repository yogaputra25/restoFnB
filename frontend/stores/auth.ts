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

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)

  if (import.meta.server) {
    const cookie = useCookie('auth_token')
    if (cookie.value) {
      const payload = parseJwt(cookie.value)
      if (payload && payload.sub && payload.exp > Date.now() / 1000) {
        accessToken.value = cookie.value
        user.value = {
          id: payload.sub,
          email: payload.email || '',
          name: payload.name || payload.email || '',
          role: payload.role,
          chainId: payload.chain,
        }
      }
    }
  }

  const isAuthenticated = computed(() => !!accessToken.value)
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

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    const cookie = useCookie('auth_token')
    cookie.value = null
  }

  return {
    user, accessToken, refreshToken,
    isAuthenticated, isAdmin, isCashier,
    login, logout,
  }
})
