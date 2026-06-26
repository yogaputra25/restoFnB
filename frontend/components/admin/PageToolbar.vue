<template>
  <div class="flex items-center justify-between mb-4 pb-3 border-b border-[var(--color-border)]">
    <div class="flex items-center gap-3 flex-1">
      <slot name="left" />
    </div>
    <div class="flex items-center gap-2 flex-shrink-0">
      <slot name="right" />
    </div>
  </div>
</template>
