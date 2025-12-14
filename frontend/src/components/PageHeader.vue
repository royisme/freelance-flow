<script setup lang="ts">
defineProps<{
    title: string
    subtitle?: string
}>()
</script>

<template>
    <div class="page-header">
        <div class="header-main">
            <div class="header-titles">
                <h1 class="title">{{ title }}</h1>
                <p v-if="subtitle" class="subtitle">{{ subtitle }}</p>
            </div>
            <div v-if="$slots.extra" class="header-extra">
                <slot name="extra"></slot>
            </div>
        </div>
    </div>
</template>

<style scoped>
.page-header {
    flex-shrink: 0;
    margin-bottom: 24px;
}

.header-main {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
}

.header-titles {
    flex: 1;
    min-width: 0;
}

.title {
    font-family: 'Varela Round', sans-serif;
    font-size: 1.75rem;
    font-weight: 700;
    margin: 0;
    color: var(--n-text-color);
    line-height: 1.2;
}

.subtitle {
    margin: 8px 0 0;
    color: var(--n-text-color-3);
    font-size: 1rem;
}

.header-extra {
    flex-shrink: 0;
    margin-left: 24px;
}
</style>
