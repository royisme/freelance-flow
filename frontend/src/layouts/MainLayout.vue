<script setup lang="ts">
import {
    NLayout, NLayoutSider, NLayoutContent, NLayoutHeader, NLayoutFooter,
    NMenu, NIcon, NButton, NTooltip, NSpace, NDropdown, NAvatar
} from 'naive-ui'
import { ref, watch, h, computed, type Component } from 'vue'
import { RouterView, useRouter, useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import type { MenuOption, DropdownOption } from 'naive-ui'
import {
    DashboardOutlined,
    UserOutlined,
    ProjectOutlined,
    ClockCircleOutlined,
    FileTextOutlined,
    BarChartOutlined,
    SettingOutlined,
    QuestionCircleOutlined,
    GlobalOutlined,
    BulbOutlined,
    BulbFilled,
    LogoutOutlined,
    SwapOutlined
} from '@vicons/antd'

function renderIcon(icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) })
}

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()
const { t } = useI18n()

// Menu options using i18n
const menuOptions = computed<MenuOption[]>(() => [
    { label: t('nav.dashboard'), key: 'dashboard', icon: renderIcon(DashboardOutlined) },
    { label: t('nav.clients'), key: 'clients', icon: renderIcon(UserOutlined) },
    { label: t('nav.projects'), key: 'projects', icon: renderIcon(ProjectOutlined) },
    { label: t('nav.timesheet'), key: 'timesheet', icon: renderIcon(ClockCircleOutlined) },
    { label: t('nav.invoices'), key: 'invoices', icon: renderIcon(FileTextOutlined) },
    { label: t('nav.reports'), key: 'reports', icon: renderIcon(BarChartOutlined) },
    { label: t('nav.settings'), key: 'settings/invoice', icon: renderIcon(SettingOutlined) },
])

const activeKey = ref<string>(route.path.substring(1) || 'dashboard')
const collapsed = ref(false)

// Locale Options
const localeOptions = [
    { label: '中文 (简体)', key: 'zh-CN' },
    { label: 'English', key: 'en-US' }
]

// User Menu Options
const userMenuOptions = computed<DropdownOption[]>(() => [
    {
        label: t('auth.switchUser'),
        key: 'switch',
        icon: renderIcon(SwapOutlined)
    },
    { type: 'divider', key: 'd1' },
    {
        label: t('auth.logout'),
        key: 'logout',
        icon: renderIcon(LogoutOutlined)
    },
])

function handleLocaleSelect(key: 'zh-CN' | 'en-US') {
    appStore.setLocale(key)
}

function handleUserMenuSelect(key: string) {
    if (key === 'logout') {
        authStore.logout()
        router.push('/splash')
    } else if (key === 'switch') {
        authStore.switchUser()
        router.push('/login')
    }
}

watch(() => route.path, (newPath) => {
    // Update active menu key based on path
    const key = newPath.substring(1).split('/')[0] // Get first segment
    activeKey.value = key || 'dashboard'
})

function handleMenuUpdate(key: string) {
    router.push('/' + key)
}
</script>

<template>
    <n-layout style="height: 100vh">

        <!-- 1. Header Area -->
        <n-layout-header bordered class="app-header">
            <div class="header-left">
                <div class="brand-logo">FreelanceFlow</div>
            </div>
            <div class="header-right">
                <n-space size="large" align="center">

                    <!-- Language Switcher -->
                    <n-dropdown :options="localeOptions" @select="handleLocaleSelect">
                        <n-button quaternary circle>
                            <template #icon><n-icon>
                                    <GlobalOutlined />
                                </n-icon></template>
                        </n-button>
                    </n-dropdown>

                    <!-- Theme Toggle -->
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

                    <div class="divider-vertical"></div>

                    <n-button quaternary circle>
                        <template #icon><n-icon>
                                <SettingOutlined />
                            </n-icon></template>
                    </n-button>
                    <n-button quaternary circle>
                        <template #icon><n-icon>
                                <QuestionCircleOutlined />
                            </n-icon></template>
                    </n-button>

                    <!-- User Menu -->
                    <n-dropdown :options="userMenuOptions" @select="handleUserMenuSelect">
                        <div class="user-menu-trigger">
                            <n-avatar :size="32"
                                :src="authStore.avatarUrl || `https://api.dicebear.com/9.x/avataaars/svg?seed=${authStore.username}`" />
                            <span class="username">{{ authStore.username }}</span>
                        </div>
                    </n-dropdown>
                </n-space>
            </div>
        </n-layout-header>

        <!-- 2. Main Body Area (Sider + Content) -->
        <n-layout position="absolute" style="top: 64px; bottom: 32px;" has-sider>
            <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="240" :collapsed="collapsed"
                show-trigger @collapse="collapsed = true" @expand="collapsed = false">
                <n-menu :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions"
                    :value="activeKey" @update:value="handleMenuUpdate" />
            </n-layout-sider>

            <n-layout-content class="app-content">
                <RouterView />
            </n-layout-content>
        </n-layout>

        <!-- 3. Footer / Status Bar -->
        <n-layout-footer bordered position="absolute" class="app-footer">
            <div class="status-bar">
                <span class="status-item">{{ t('footer.statusBar') }}</span>
                <span class="status-item">{{ t('footer.weeklyHours') }} <strong>32h</strong></span>
                <span class="divider">|</span>
                <span class="status-item">{{ t('footer.pendingPayment') }} <strong>$2,400</strong></span>
            </div>
        </n-layout-footer>

    </n-layout>
</template>

<style scoped>
/* Header Styles */
.app-header {
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
    z-index: 10;
}

/* Brand Logo Area */
.brand-logo {
    font-family: 'Varela Round', sans-serif;
    font-size: 1.5rem;
    font-weight: 800;
    color: var(--n-primary-color);
    letter-spacing: -0.02em;
}

/* Content Styles */
.app-content {
    padding: 24px;
}

/* Footer Styles */
.app-footer {
    height: 32px;
    padding: 0 24px;
    display: flex;
    align-items: center;
    background-color: var(--n-color-modal);
    color: var(--n-text-color-3);
    font-size: 0.85rem;
    font-family: 'Inter', sans-serif;
    border-top: 1px solid var(--n-border-color);
}

.status-bar {
    display: flex;
    align-items: center;
    gap: 12px;
}

.divider {
    color: var(--n-text-color-3);
    opacity: 0.5;
}

strong {
    color: var(--n-primary-color);
    font-weight: 600;
}

/* User Menu */
.user-menu-trigger {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 4px 12px 4px 4px;
    border-radius: 20px;
    cursor: pointer;
    transition: background-color 0.2s;
}

.user-menu-trigger:hover {
    background-color: rgba(0, 0, 0, 0.05);
}

.username {
    font-size: 0.9rem;
    font-weight: 500;
}
</style>
