export default defineNuxtRouteMiddleware((to, from) => {
  const { isAuthenticated, checkAuth } = useAuth()
  
  // If already authenticated, redirect to dashboard
  if (checkAuth()) {
    return navigateTo('/dashboard')
  }
})