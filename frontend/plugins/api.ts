export default defineNuxtPlugin(() => {
  const auth = useAuthStore()
  const router = useRouter()

  async function refreshAccessToken(): Promise<string | null> {
    if (!auth.refreshToken) return null
    try {
      const res = await $fetch('/api/auth/refresh', {
        method: 'POST',
        body: { refresh_token: auth.refreshToken },
      })
      const data = (res as any)?.data
      if (data?.access_token) {
        const payload = parseJwt(data.access_token)
        auth.setSession(data.access_token, {
          id: payload?.sub || auth.user?.id || '',
          email: payload?.email || auth.user?.email || '',
          name: payload?.name || auth.user?.name || '',
          role: payload?.role || auth.user?.role || 'customer',
          chainId: payload?.chain || auth.user?.chainId || '',
          branchId: payload?.branch || auth.user?.branchId || undefined,
        })
        if (data.refresh_token) {
          auth.refreshToken = data.refresh_token
        }
        return data.access_token
      }
      return null
    } catch {
      return null
    }
  }

  async function api(url: string, opts: Record<string, any> = {}) {
    const options = { ...opts }
    options.headers = { ...options.headers }

    let finalUrl = url
    let token = auth.accessToken

    // Attach token
    if (token) {
      options.headers.Authorization = `Bearer ${token}`
    }

    // Inject chain_id
    if (auth.user?.chainId) {
      const sep = finalUrl.includes('?') ? '&' : '?'
      finalUrl += `${sep}chain_id=${auth.user.chainId}`
      if (options.method && ['POST', 'PUT', 'PATCH'].includes(options.method)) {
        if (options.body && typeof options.body === 'object' && !Array.isArray(options.body)) {
          ;(options.body as Record<string, unknown>).chain_id = auth.user.chainId
        }
      }
    }

    try {
      return await $fetch(finalUrl, options)
    } catch (err: any) {
      const status = err?.response?.status

      // 401 — try refresh once, then logout
      if (status === 401) {
        const newToken = await refreshAccessToken()
        if (newToken) {
          // Retry with new token
          options.headers.Authorization = `Bearer ${newToken}`
          return await $fetch(finalUrl, options)
        }
        // Refresh failed — logout
        auth.logout()
        router.push('/login')
      }

      // 403 — permission error (already handled by backend response)
      // 500 — server error (bubble up for toast handling)

      throw err
    }
  }

  return {
    provide: { api },
  }
})

function parseJwt(token: string) {
  try {
    const base64 = token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')
    return JSON.parse(atob(base64))
  } catch {
    return null
  }
}
