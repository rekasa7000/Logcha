<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
          Create your account
        </h2>
      </div>
      
      <UCard class="p-6">
        <form @submit.prevent="handleRegister" class="space-y-6">
          <div>
            <UFormGroup label="Full Name" name="name">
              <UInput
                v-model="form.name"
                type="text"
                placeholder="Enter your full name"
                required
                :disabled="loading"
              />
            </UFormGroup>
          </div>
          
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
          
          <div>
            <UFormGroup label="Confirm Password" name="confirmPassword">
              <UInput
                v-model="form.confirmPassword"
                type="password"
                placeholder="Confirm your password"
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
              :disabled="loading || !isFormValid"
              class="w-full"
            >
              Create Account
            </UButton>
          </div>
          
          <div class="text-center">
            <p class="text-sm text-gray-600 dark:text-gray-400">
              Already have an account?
              <NuxtLink to="/login" class="font-medium text-primary-600 hover:text-primary-500">
                Sign in
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

const { register } = useAuth()

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const loading = ref(false)
const error = ref('')

const isFormValid = computed(() => {
  return form.password === form.confirmPassword && form.password.length >= 6
})

const handleRegister = async () => {
  if (!isFormValid.value) {
    error.value = 'Passwords do not match or are too short (minimum 6 characters)'
    return
  }
  
  loading.value = true
  error.value = ''
  
  try {
    await register({
      name: form.name,
      email: form.email,
      password: form.password
    })
  } catch (err: any) {
    error.value = err.data?.message || 'Registration failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>