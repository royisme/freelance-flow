<script setup lang="ts">
import { NElement } from 'naive-ui'

defineProps<{
  fill?: boolean  // When true, content fills available space with flexbox
}>()
</script>

<template>
  <n-element tag="div" class="page-container" :class="{ 'fill-mode': fill }">
    <slot></slot>
  </n-element>
</template>

<style scoped>
.page-container {
  box-sizing: border-box;
  height: 100%;
  padding: var(--space-6);
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-xl);
}

/* Fill mode - content uses flexbox to fill available space */
.page-container.fill-mode {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
