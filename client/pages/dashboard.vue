<template>
  <div class="space-y-6">
    <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-4">
        Dashboard
      </h1>
      <p class="text-gray-600 dark:text-gray-400">
        Welcome to your dashboard, {{ user?.name }}!
      </p>
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <UCard class="p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-blue-600 dark:text-blue-400 mb-2">
            0
          </div>
          <div class="text-gray-600 dark:text-gray-400">
            Total Time Records
          </div>
        </div>
      </UCard>
      
      <UCard class="p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-green-600 dark:text-green-400 mb-2">
            0
          </div>
          <div class="text-gray-600 dark:text-gray-400">
            Hours This Week
          </div>
        </div>
      </UCard>
      
      <UCard class="p-6">
        <div class="text-center">
          <div class="text-3xl font-bold text-purple-600 dark:text-purple-400 mb-2">
            0
          </div>
          <div class="text-gray-600 dark:text-gray-400">
            Active Projects
          </div>
        </div>
      </UCard>
    </div>
    
    <div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
        Recent Activity
      </h2>
      <div class="text-gray-600 dark:text-gray-400">
        No recent activity found. Start tracking your time!
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const { user } = useAuth()
</script>