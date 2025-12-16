<script setup lang="ts">
defineProps<{
  fill?: boolean  // When true, content fills available space with flexbox
}>()
</script>

<template>
  <div class="h-full flex flex-col min-h-0"
    :class="{ 'overflow-hidden': fill }">
    <slot></slot>
  </div>
</template>

<style scoped>
/* Scoped styles removed */
</style>

