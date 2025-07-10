<template>
  <div class="text-center space-y-8">
    <div class="max-w-3xl mx-auto">
      <h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-4">
        Welcome to Logcha
      </h1>
      <p class="text-xl text-gray-600 dark:text-gray-400 mb-8">
        Your comprehensive time tracking and logging solution
      </p>
      
      <div class="flex justify-center space-x-4">
        <template v-if="isAuthenticated">
          <NuxtLink to="/dashboard">
            <UButton size="lg">Go to Dashboard</UButton>
          </NuxtLink>
        </template>
        <template v-else>
          <NuxtLink to="/register">
            <UButton size="lg">Get Started</UButton>
          </NuxtLink>
          <NuxtLink to="/login">
            <UButton variant="outline" size="lg">Sign In</UButton>
          </NuxtLink>
        </template>
      </div>
    </div>
    
    <div class="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8">
      <div class="text-center">
        <div class="bg-blue-100 dark:bg-blue-900 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          Time Tracking
        </h3>
        <p class="text-gray-600 dark:text-gray-400">
          Track your time efficiently with our intuitive interface
        </p>
      </div>
      
      <div class="text-center">
        <div class="bg-green-100 dark:bg-green-900 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          Analytics
        </h3>
        <p class="text-gray-600 dark:text-gray-400">
          Get insights into your productivity and time usage patterns
        </p>
      </div>
      
      <div class="text-center">
        <div class="bg-purple-100 dark:bg-purple-900 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
          <svg class="w-8 h-8 text-purple-600 dark:text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
          Reports
        </h3>
        <p class="text-gray-600 dark:text-gray-400">
          Generate detailed reports for better project management
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { isAuthenticated } = useAuth()
</script>