<script setup lang="ts">
import { RouterView } from 'vue-router'
import { NButton, NDropdown, NIcon, NSpace, NTooltip } from 'naive-ui'
import { GlobalOutlined, BulbOutlined, BulbFilled } from '@vicons/antd'
import { useAppStore } from '@/stores/app'
import { useI18n } from 'vue-i18n'

const appStore = useAppStore()
const { t } = useI18n()

// Locale Options
const localeOptions = [
    { label: '中文 (简体)', key: 'zh-CN' },
    { label: 'English', key: 'en-US' }
]

function handleLocaleSelect(key: 'zh-CN' | 'en-US') {
    appStore.setLocale(key)
}
</script>

<template>
    <div class="auth-layout auth-gradient-bg">
        <!-- Transparent Header on gradient -->
        <header class="auth-header">
            <div class="brand">FreelanceFlow</div>
            <n-space>
                <n-dropdown :options="localeOptions" @select="handleLocaleSelect">
                    <n-button quaternary circle>
                        <template #icon><n-icon>
                                <GlobalOutlined />
                            </n-icon></template>
                    </n-button>
                </n-dropdown>
                <n-tooltip trigger="hover">
                    <template #trigger>
                        <n-button quaternary circle @click="appStore.toggleTheme()">
                            <template #icon>
                                <n-icon>
                                    <BulbFilled v-if="appStore.theme === 'dark'" />
                                    <BulbOutlined v-else />
                                </n-icon>
                            </template>
                        </n-button>
                    </template>
                    {{ appStore.theme === 'dark' ? t('theme.switchToLight') : t('theme.switchToDark') }}
                </n-tooltip>
            </n-space>
        </header>

        <!-- Main Content -->
        <div class="auth-content">
            <RouterView />
        </div>
    </div>
</template>

<style scoped>
.auth-layout {
    min-height: 100vh;
    width: 100vw;
    display: flex;
    flex-direction: column;
}

.auth-header {
    padding: var(--space-4) var(--space-8);
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.brand {
    font-family: var(--font-display);
    font-size: var(--text-xl);
    font-weight: 800;
    color: white;
    text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.auth-content {
    flex: 1;
    min-height: 0;
    display: flex;
    justify-content: center;
    align-items: center;
    padding: var(--space-4);
    overflow: hidden;
}

/* Override button colors for visibility on gradient */
.auth-header :deep(.n-button) {
    color: white;
}

.auth-header :deep(.n-button:hover) {
    background-color: rgba(255, 255, 255, 0.15);
}
</style>
