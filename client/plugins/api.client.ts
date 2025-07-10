export default defineNuxtPlugin(() => {
  const { $fetch } = useNuxtApp()
  
  // Create a custom $fetch instance with auth headers
  const api = $fetch.create({
    baseURL: useRuntimeConfig().public.apiBase,
    onRequest({ request, options }) {
      const token = useCookie('auth-token')
      
      if (token.value) {
        options.headers = {
          ...options.headers,
          Authorization: `Bearer ${token.value}`
        }
      }
    },
    onResponseError({ response }) {
      // Handle 401 errors by logging out
      if (response.status === 401) {
        const { logout } = useAuth()
        logout()
      }
    }
  })

  return {
    provide: {
      api
    }
  }
})