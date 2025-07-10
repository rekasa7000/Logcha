<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <nav class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <NuxtLink to="/" class="text-xl font-semibold text-gray-900 dark:text-white">
              Logcha
            </NuxtLink>
          </div>
          <div class="flex items-center space-x-4">
            <template v-if="isAuthenticated">
              <span class="text-sm text-gray-700 dark:text-gray-300">
                Welcome, {{ user?.name }}
              </span>
              <UButton @click="logout" variant="outline">
                Logout
              </UButton>
            </template>
            <template v-else>
              <NuxtLink to="/login">
                <UButton variant="outline">Login</UButton>
              </NuxtLink>
              <NuxtLink to="/register">
                <UButton>Register</UButton>
              </NuxtLink>
            </template>
          </div>
        </div>
      </div>
    </nav>
    
    <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
const { user, isAuthenticated, logout } = useAuth()
</script>