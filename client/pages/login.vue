<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
          Sign in to your account
        </h2>
      </div>
      
      <UCard class="p-6">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <UFormGroup label="Email address" name="email">
              <UInput
                v-model="form.email"
                type="email"
                placeholder="Enter your email"
                required
                :disabled="loading"
              />
            </UFormGroup>
          </div>
          
          <div>
            <UFormGroup label="Password" name="password">
              <UInput
                v-model="form.password"
                type="password"
                placeholder="Enter your password"
                required
                :disabled="loading"
              />
            </UFormGroup>
          </div>
          
          <div v-if="error" class="text-red-600 text-sm">
            {{ error }}
          </div>
          
          <div>
            <UButton
              type="submit"
              :loading="loading"
              :disabled="loading"
              class="w-full"
            >
              Sign in
            </UButton>
          </div>
          
          <div class="text-center">
            <p class="text-sm text-gray-600 dark:text-gray-400">
              Don't have an account?
              <NuxtLink to="/register" class="font-medium text-primary-600 hover:text-primary-500">
                Sign up
              </NuxtLink>
            </p>
          </div>
        </form>
      </UCard>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: false,
  middleware: 'guest'
})

const { login } = useAuth()

const form = reactive({
  email: '',
  password: ''
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  loading.value = true
  error.value = ''
  
  try {
    await login(form)
  } catch (err: any) {
    error.value = err.data?.message || 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>