interface ApiResponse<T> {
  data: T
  message?: string
  success: boolean
}

export const useApi = () => {
  const config = useRuntimeConfig()
  
  const apiCall = async <T>(
    endpoint: string,
    options: {
      method?: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH'
      body?: any
      query?: Record<string, any>
      headers?: Record<string, string>
    } = {}
  ): Promise<ApiResponse<T>> => {
    const token = useCookie('auth-token')
    
    const defaultHeaders: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    
    if (token.value) {
      defaultHeaders.Authorization = `Bearer ${token.value}`
    }
    
    try {
      const response = await $fetch<ApiResponse<T>>(endpoint, {
        baseURL: config.public.apiBase,
        method: options.method || 'GET',
        body: options.body,
        query: options.query,
        headers: {
          ...defaultHeaders,
          ...options.headers
        }
      })
      
      return response
    } catch (error: any) {
      console.error('API call error:', error)
      
      // Handle 401 errors
      if (error.response?.status === 401) {
        const { logout } = useAuth()
        await logout()
      }
      
      throw createError({
        statusCode: error.response?.status || 500,
        statusMessage: error.response?.data?.message || 'API call failed'
      })
    }
  }

  // Specific API methods
  const get = <T>(endpoint: string, query?: Record<string, any>) =>
    apiCall<T>(endpoint, { method: 'GET', query })

  const post = <T>(endpoint: string, body?: any) =>
    apiCall<T>(endpoint, { method: 'POST', body })

  const put = <T>(endpoint: string, body?: any) =>
    apiCall<T>(endpoint, { method: 'PUT', body })

  const patch = <T>(endpoint: string, body?: any) =>
    apiCall<T>(endpoint, { method: 'PATCH', body })

  const del = <T>(endpoint: string) =>
    apiCall<T>(endpoint, { method: 'DELETE' })

  return {
    api: apiCall,
    get,
    post,
    put,
    patch,
    delete: del
  }
}