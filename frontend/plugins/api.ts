export default defineNuxtPlugin(() => {
  const auth = useAuthStore()
  const router = useRouter()

  function api(url: string, opts: Record<string, any> = {}) {
    const options = { ...opts }
    options.headers = { ...options.headers }

    let finalUrl = url
    if (auth.accessToken) {
      options.headers.Authorization = `Bearer ${auth.accessToken}`
    }
    if (auth.user?.chainId) {
      const sep = finalUrl.includes('?') ? '&' : '?'
      finalUrl += `${sep}chain_id=${auth.user.chainId}`
      console.log('[api] chain_id injected:', auth.user.chainId, '→', finalUrl)
      if (options.method && ['POST', 'PUT', 'PATCH'].includes(options.method)) {
        if (options.body && typeof options.body === 'object' && !Array.isArray(options.body)) {
          ;(options.body as Record<string, unknown>).chain_id = auth.user.chainId
        }
      }
    }

    return $fetch(finalUrl, options).catch((err: any) => {
      if (err?.response?.status === 401) {
        auth.logout()
        router.push('/login')
      }
      throw err
    })
  }

  return {
    provide: { api },
  }
})
