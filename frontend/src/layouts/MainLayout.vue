<script setup lang="ts">
import {
    NLayout, NLayoutSider, NLayoutContent, NLayoutHeader, NLayoutFooter,
    NMenu, NIcon, NButton, NTooltip, NSpace, NDropdown, NAvatar
} from 'naive-ui'
import { ref, watch, h, computed, onMounted, type Component } from 'vue'
import { RouterView, useRouter, useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useStatusBarStore } from '@/stores/statusBar'
import { useI18n } from 'vue-i18n'
import type { MenuOption, DropdownOption } from 'naive-ui'
import {
    DashboardOutlined,
    UserOutlined,
    ProjectOutlined,
    ClockCircleOutlined,
    FileTextOutlined,
    BarChartOutlined,
    WalletOutlined,
    SettingOutlined,
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
const statusBarStore = useStatusBarStore()
const { t } = useI18n()

// Menu options using i18n
const menuOptions = computed<MenuOption[]>(() => [
    { label: t('nav.dashboard'), key: 'dashboard', icon: renderIcon(DashboardOutlined) },
    { label: t('nav.clients'), key: 'clients', icon: renderIcon(UserOutlined) },
    { label: t('nav.projects'), key: 'projects', icon: renderIcon(ProjectOutlined) },
    { label: t('nav.timesheet'), key: 'timesheet', icon: renderIcon(ClockCircleOutlined) },
    { label: t('nav.invoices'), key: 'invoices', icon: renderIcon(FileTextOutlined) },
    { label: t('nav.reports'), key: 'reports', icon: renderIcon(BarChartOutlined) },
    {
        label: t('nav.finance'),
        key: 'finance',
        icon: renderIcon(WalletOutlined),
        children: [
            {
                label: t('finance.nav.overview'),
                key: 'finance/overview',
            },
            {
                label: t('finance.nav.accounts'),
                key: 'finance/accounts',
            },
            {
                label: t('finance.nav.transactions'),
                key: 'finance/transactions',
            },
            {
                label: t('finance.nav.import'),
                key: 'finance/import',
            },
            {
                label: t('finance.nav.categories'),
                key: 'finance/categories',
            },
            {
                label: t('finance.nav.reports'),
                key: 'finance/reports',
            },
        ],
    },
    {
        label: t('nav.settings'),
        key: 'settings',
        icon: renderIcon(SettingOutlined),
        children: [
            {
                label: t('settings.general.title'),
                key: 'settings/general',
            },
            {
                label: t('settings.profile.title'),
                key: 'settings/profile',
            },
            {
                label: t('settings.invoice.title'),
                key: 'settings/invoice',
            },
            {
                label: t('settings.email.title'),
                key: 'settings/email',
            },
            {
                label: t('settings.finance.title'),
                key: 'settings/finance',
            },
        ],
    },
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

onMounted(() => {
    statusBarStore.refresh()
})
</script>

<template>
    <n-layout class="main-layout">

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
                </n-space>
            </div>
        </n-layout-header>

        <!-- 2. Main Body Area (Sider + Content) -->
        <n-layout position="absolute" style="top: 64px; bottom: 32px;" has-sider>
            <n-layout-sider bordered collapse-mode="width" :collapsed-width="64" :width="240" :collapsed="collapsed"
                show-trigger @collapse="collapsed = true" @expand="collapsed = false" class="app-sider">
                <div class="sider-content">
                    <n-menu :collapsed="collapsed" :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions"
                        :value="activeKey" @update:value="handleMenuUpdate" />
                </div>
                <!-- User Menu at bottom of sidebar -->
                <div class="sider-footer" :class="{ 'collapsed': collapsed }">
                    <n-dropdown :options="userMenuOptions" @select="handleUserMenuSelect" placement="right-start">
                        <div class="user-menu-trigger" :class="{ 'collapsed': collapsed }">
                            <n-avatar :size="collapsed ? 36 : 32"
                                :src="authStore.avatarUrl || `https://api.dicebear.com/9.x/avataaars/svg?seed=${authStore.username}`" />
                            <span v-if="!collapsed" class="username">{{ authStore.username }}</span>
                        </div>
                    </n-dropdown>
                </div>
            </n-layout-sider>

            <n-layout-content class="app-content">
                <RouterView />
            </n-layout-content>
        </n-layout>

        <!-- 3. Footer / Status Bar -->
        <n-layout-footer bordered position="absolute" class="app-footer">
            <div class="status-bar">
                <span class="status-item">{{ t('footer.statusBar') }}</span>
                <span class="status-item">{{ t('footer.monthlyHours') }} <strong>{{ statusBarStore.monthHoursLabel }}</strong></span>
                <span class="divider">|</span>
                <span class="status-item">{{ t('footer.uninvoiced') }} <strong>{{ statusBarStore.uninvoicedTotalLabel }}</strong></span>
                <span class="divider">|</span>
                <span class="status-item">{{ t('footer.pendingPayment') }} <strong>{{ statusBarStore.unpaidTotalLabel }}</strong></span>
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
    padding: 0 var(--space-6);
    z-index: 10;
}

/* Brand Logo Area */
.brand-logo {
    font-family: var(--font-display);
    font-size: var(--text-2xl);
    font-weight: 800;
    color: var(--n-primary-color);
    letter-spacing: -0.02em;
}

/* Content Styles */
.app-content {
    padding: var(--space-6);
    height: 100%;
    overflow: auto;
}

/* Footer Styles */
.app-footer {
    height: 32px;
    padding: 0 var(--space-6);
    display: flex;
    align-items: center;
    background-color: var(--n-color-modal);
    color: var(--n-text-color-3);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    border-top: 1px solid var(--n-border-color);
}

.main-layout {
    height: 100%;
    position: relative;
    overflow: hidden;
}

.status-bar {
    display: flex;
    align-items: center;
    gap: var(--space-3);
}

.divider {
    color: var(--n-text-color-3);
    opacity: 0.5;
}

strong {
    color: var(--n-primary-color);
    font-weight: 600;
}

/* Sidebar Layout */
.app-sider {
    display: flex;
    flex-direction: column;
}

:deep(.app-sider > .n-layout-sider-scroll-container) {
    display: flex !important;
    flex-direction: column !important;
}

.sider-content {
    flex: 1;
    overflow: auto;
}

.sider-footer {
    padding: 12px 20px;
    border-top: 1px solid var(--n-border-color);
}

.sider-footer.collapsed {
    padding: 12px 0;
    display: flex;
    justify-content: center;
}

/* User Menu - aligned with Naive UI menu item padding */
.user-menu-trigger {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 8px 0;
    cursor: pointer;
    transition: opacity var(--transition-normal);
    width: 100%;
}

.user-menu-trigger.collapsed {
    justify-content: center;
    padding: 8px;
}

.user-menu-trigger:hover {
    opacity: 0.8;
}

.username {
    font-size: var(--text-sm);
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
