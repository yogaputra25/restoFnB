<template>
  <div class="min-h-screen bg-[var(--color-background)] flex">
    <AdminSidebar
      :collapsed="sidebarCollapsed"
      :mobile-open="sidebarMobileOpen"
      @toggle-collapse="sidebarCollapsed = !sidebarCollapsed"
      @close-mobile="sidebarMobileOpen = false"
    />

    <div class="flex-1 flex flex-col min-w-0">
      <AdminHeader @toggle-mobile="sidebarMobileOpen = !sidebarMobileOpen" />

      <main class="flex-1 p-6 overflow-auto">
        <slot />
      </main>

      <AdminFooter />
    </div>

    <Toast />
  </div>
</template>

<script setup lang="ts">
const sidebarCollapsed = ref(false)
const sidebarMobileOpen = ref(false)
</script>
