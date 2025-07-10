export default defineNuxtRouteMiddleware((to, from) => {
  const { isAuthenticated, checkAuth } = useAuth()
  
  // Check authentication status
  if (!checkAuth()) {
    // Redirect to login if not authenticated
    return navigateTo('/login')
  }
})