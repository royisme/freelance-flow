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
  height: 100%;
  padding: 0 4px 0 0;
}

/* Fill mode - content uses flexbox to fill available space */
.page-container.fill-mode {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
