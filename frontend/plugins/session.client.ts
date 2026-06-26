/**
 * Session restore plugin.
 * Runs on every client-side app load to validate the stored token
 * against the backend. If valid, restores user session. If not,
 * clears the token without redirecting (middleware handles redirect).
 */
export default defineNuxtPlugin({
  name: 'session-restore',
  enforce: 'pre', // run before other plugins
  async setup() {
    const auth = useAuthStore()

    // SSR already hydrated; no need for /auth/me call on server
    if (import.meta.server) {
      auth.restoring = false
      return
    }

    // No token — nothing to restore
    if (!auth.accessToken) {
      auth.restoring = false
      return
    }

    try {
      const res = await $fetch('/api/auth/me', {
        headers: { Authorization: `Bearer ${auth.accessToken}` },
      })
      const profile = (res as any)?.data
      if (profile) {
        auth.setSession(auth.accessToken, {
          id: profile.id,
          email: profile.email || '',
          name: profile.name || '',
          role: profile.role || '',
          chainId: profile.chain_id || '',
        })
      }
    } catch {
      // Token invalid or expired — clear it
      auth.logout()
    } finally {
      auth.restoring = false
    }
  },
})
