import { jwtDecode } from 'jwt-decode'
import type { JwtPayload } from 'jwt-decode'

interface AuthUser {
  id: string
  email: string
  name: string
  role: string
}

interface CustomJwtPayload extends JwtPayload {
  id: string
  email: string
  name: string
  role: string
}

export const useAuth = () => {
  const user = ref<AuthUser | null>(null)
  const isAuthenticated = computed(() => !!user.value)

  const login = async (credentials: { email: string; password: string }) => {
    try {
      const { data } = await $fetch<{ token: string; user: AuthUser }>('/api/auth/login', {
        method: 'POST',
        body: credentials,
        baseURL: useRuntimeConfig().public.apiBase
      })

      const token = useCookie('auth-token', {
        httpOnly: true,
        secure: true,
        sameSite: 'strict',
        maxAge: 60 * 60 * 24 * 7 // 7 days
      })

      token.value = data.token
      user.value = data.user

      await navigateTo('/dashboard')
      
      return { success: true, data }
    } catch (error: any) {
      console.error('Login error:', error)
      throw createError({
        statusCode: error.response?.status || 401,
        statusMessage: error.response?.data?.message || 'Login failed'
      })
    }
  }

  const register = async (userData: { 
    email: string; 
    password: string; 
    name: string;
    role?: string;
  }) => {
    try {
      const { data } = await $fetch<{ token: string; user: AuthUser }>('/api/auth/register', {
        method: 'POST',
        body: userData,
        baseURL: useRuntimeConfig().public.apiBase
      })

      const token = useCookie('auth-token', {
        httpOnly: true,
        secure: true,
        sameSite: 'strict',
        maxAge: 60 * 60 * 24 * 7 // 7 days
      })

      token.value = data.token
      user.value = data.user

      await navigateTo('/dashboard')
      
      return { success: true, data }
    } catch (error: any) {
      console.error('Registration error:', error)
      throw createError({
        statusCode: error.response?.status || 400,
        statusMessage: error.response?.data?.message || 'Registration failed'
      })
    }
  }

  const logout = async () => {
    const token = useCookie('auth-token')
    token.value = null
    user.value = null
    await navigateTo('/login')
  }

  const checkAuth = () => {
    const token = useCookie('auth-token')
    
    if (!token.value) {
      user.value = null
      return false
    }

    try {
      const decoded = jwtDecode<CustomJwtPayload>(token.value)
      
      // Check if token is expired
      if (decoded.exp && decoded.exp * 1000 < Date.now()) {
        token.value = null
        user.value = null
        return false
      }

      user.value = {
        id: decoded.id,
        email: decoded.email,
        name: decoded.name,
        role: decoded.role
      }
      
      return true
    } catch (error) {
      console.error('Token validation error:', error)
      token.value = null
      user.value = null
      return false
    }
  }

  const refreshToken = async () => {
    try {
      const { data } = await $fetch<{ token: string }>('/api/auth/refresh', {
        method: 'POST',
        baseURL: useRuntimeConfig().public.apiBase
      })

      const token = useCookie('auth-token', {
        httpOnly: true,
        secure: true,
        sameSite: 'strict',
        maxAge: 60 * 60 * 24 * 7 // 7 days
      })

      token.value = data.token
      checkAuth()
      
      return true
    } catch (error) {
      console.error('Token refresh error:', error)
      await logout()
      return false
    }
  }

  // Initialize auth state
  onMounted(() => {
    checkAuth()
  })

  return {
    user: readonly(user),
    isAuthenticated,
    login,
    register,
    logout,
    checkAuth,
    refreshToken
  }
}