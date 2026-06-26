export default defineNuxtRouteMiddleware(async (to) => {
  const auth = useAuthStore()

  // Wait for session restore to complete
  if (auth.restoring) {
    await new Promise<void>((resolve) => {
      const unwatch = watch(
        () => auth.restoring,
        (val) => {
          if (!val) {
            unwatch()
            resolve()
          }
        },
        { immediate: true }
      )
      // Safety timeout: if restoring never completes, proceed anyway
      setTimeout(() => { unwatch(); resolve() }, 5000)
    })
  }

  if (!auth.isAuthenticated) {
    return navigateTo('/login')
  }

  if (to.path.startsWith('/admin') && !auth.isAdmin) {
    return navigateTo('/')
  }
})
