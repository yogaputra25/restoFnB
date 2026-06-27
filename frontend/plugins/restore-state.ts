import { parseCookies } from 'h3'

export default defineNuxtPlugin(() => {
  const auth = useAuthStore()

  function parseJwt(token: string) {
    try {
      const base64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
      return JSON.parse(atob(base64))
    } catch { return null }
  }

  function isExpired(token: string): boolean {
    const p = parseJwt(token)
    return !p || !p.exp || Date.now() / 1000 > p.exp
  }

  function restoreFromPayload(token: string, payload: Record<string, unknown>) {
    auth.setSession(token, {
      id: payload.sub as string,
      email: (payload.email as string) || '',
      name: (payload.name as string) || (payload.email as string) || '',
      role: payload.role as 'admin' | 'cashier' | 'customer',
      chainId: payload.chain as string,
      branchId: (payload.branch as string) || undefined,
    })
  }

  if (import.meta.server) {
    try {
      const event = useRequestEvent()
      if (event) {
        const cookies = parseCookies(event)
        const token = cookies['auth_token']
        if (token && !isExpired(token)) {
          const payload = parseJwt(token)
          if (payload) restoreFromPayload(token, payload)
        }
      }
    } catch { /* cookie unavailable */ }
  }

  if (import.meta.client) {
    try {
      const token = localStorage.getItem('auth_token')
      if (token && !isExpired(token)) {
        const payload = parseJwt(token)
        if (payload) restoreFromPayload(token, payload)
      } else if (token) {
        localStorage.removeItem('auth_token')
        localStorage.removeItem('refresh_token')
      }
    } catch { /* localStorage unavailable */ }
  }
})
