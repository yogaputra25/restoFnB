<script setup lang="ts">
const route = useRoute()
const router = useRouter()

const branch = route.query.branch as string
const table = route.query.table as string

if (branch && table) {
  router.replace(`/order?branch=${branch}&table=${table}`)
} else {
  router.replace('/')
}
</script>

<template>
  <div class="min-h-screen bg-dark-900 flex items-center justify-center">
    <p class="text-white">Redirecting...</p>
  </div>
</template>
