<template>
  <NuxtLayout>
    <AppErrorPage
      :type="error.statusCode === 404 ? '404' : '500'"
      :title="error.statusCode === 404 ? 'Dish Not Found' : 'Kitchen is on Fire'"
      :message="error.statusCode === 404
        ? 'Looks like this dish wandered off the menu. Let\'s get you back to the good stuff.'
        : 'Something caught fire in the kitchen. Our team is working on it.'"
    />
  </NuxtLayout>
</template>

<script setup lang="ts">
import type { NuxtError } from '#app'

const props = defineProps<{ error: NuxtError }>()

const handleError = useError()
</script>
