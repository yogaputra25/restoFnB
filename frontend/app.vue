<template>
  <NuxtLayout>
    <Transition name="page-fade" mode="out-in">
      <NuxtPage />
    </Transition>
  </NuxtLayout>
</template>

<style>
.page-fade-enter-active,
.page-fade-leave-active {
  transition: opacity var(--transition-base);
}
.page-fade-enter-from,
.page-fade-leave-to {
  opacity: 0;
}
</style>
