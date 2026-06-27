export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuthStore()

  if (!auth.isAuthenticated) {
    return navigateTo('/login')
  }

  if (to.path.startsWith('/admin') && !auth.isAdmin) {
    return navigateTo('/')
  }
})
