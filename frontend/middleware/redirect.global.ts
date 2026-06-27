export default defineNuxtRouteMiddleware((to) => {
  const auth = useAuthStore()

  if (to.path === '/') {
    return navigateTo('/login')
  }

  if (to.path === '/login' && auth.isAuthenticated) {
    if (auth.isAdmin) return navigateTo('/admin')
    if (auth.isCashier) return navigateTo('/cashier')
    return navigateTo('/menu')
  }
})
