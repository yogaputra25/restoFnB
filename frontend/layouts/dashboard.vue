<template>
  <div class="min-h-screen bg-dark-900 flex">
    <Sidebar />
    <main class="flex-1 overflow-auto">
      <slot />
    </main>
    <Toast />
  </div>
</template>

<script setup lang="ts">
import Sidebar from '~/components/Sidebar.vue'
</script>
