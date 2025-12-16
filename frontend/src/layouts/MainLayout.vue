<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView, useRouter, useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useStatusBarStore } from '@/stores/statusBar'
import { useSettingsStore } from '@/stores/settings'
import { useI18n } from 'vue-i18n'
import {
    Globe,
    Sun,
    Moon
} from 'lucide-vue-next'
import { allModules, isModuleEnabled, normalizeModuleOverrides } from '@/modules/registry'
import type { ModuleNavItem } from '@/modules/types'
import { Button } from '@/components/ui/button'
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import {
    Tooltip,
    TooltipContent,
    TooltipProvider,
    TooltipTrigger,
} from '@/components/ui/tooltip'
import {
    SidebarProvider,
    SidebarInset,
    SidebarTrigger
} from '@/components/ui/sidebar'
import { Separator } from '@/components/ui/separator'
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbList,
    BreadcrumbPage,
} from '@/components/ui/breadcrumb'
import AppSidebar from '@/components/app-sidebar/AppSidebar.vue'
import type { NavItem } from '@/components/app-sidebar/types'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()
const statusBarStore = useStatusBarStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()

const currentBreadcrumbTitle = computed(() => {
    // 1. Try module based title
    if (route.meta.moduleID) {
        // Find module
        const mod = allModules.find(m => m.id === route.meta.moduleID)
        // For sub-pages (e.g. settings/general), strict match is needed?
        // Actually moduleID is shared by all pages in module.
        // So for "settings", we might want the specific child?
        // But for now, let's just use the module label or child label if possible?

        // Improve: find the specific nav item that matches the route?
        if (mod) {
            // If module has direct nav
            if (mod.nav && !mod.nav.children) {
                return t(mod.nav.labelKey)
            }
            // If module has children logic (Settings)
            if (mod.id === 'settings') {
                // Return "Settings" or maybe specific page?
                // For "Settings", let's just return "Settings" for now, or match the child?
                // The settings children keys are `settings/general`.
                // For simplicity, return module label which is "Settings".
                return t(mod.nav?.labelKey || '')
            }
            if (mod.nav) {
                return t(mod.nav.labelKey)
            }
        }
    }
    // 2. Fallback
    return ''
})

// Locale Options
const localeOptions = [
    { label: '中文 (简体)', key: 'zh-CN' },
    { label: 'English', key: 'en-US' }
]

function handleLocaleSelect(key: 'zh-CN' | 'en-US') {
    appStore.setLocale(key)
}

function handleLogout() {
    authStore.logout()
    router.push('/splash')
}

function handleSwitchUser() {
    authStore.switchUser()
    router.push('/login')
}

const userData = computed(() => ({
    name: authStore.username || t('common.user'),
    email: t('common.freelancer'),
    avatar: authStore.avatarUrl || `https://api.dicebear.com/9.x/avataaars/svg?seed=${authStore.username}`
}))


// Navigation Items Mapping
function toNavItem(item: ModuleNavItem): NavItem {
    const url = '/' + item.key
    const isActive = route.path === url || route.path.startsWith(url + '/')

    const navItem: NavItem = {
        title: t(item.labelKey),
        url,
        icon: item.icon,
        isActive
    }

    // Map children recursively
    if (item.children && item.children.length > 0) {
        navItem.children = item.children.map(child => toNavItem(child))
    }

    return navItem
}

const platformItems = computed<NavItem[]>(() => {
    const overrides = normalizeModuleOverrides(settingsStore.settings?.moduleOverrides)
    return allModules
        .filter((m) => m.id !== 'settings') // Exclude settings
        .filter((m) => m.nav)
        .filter((m) => isModuleEnabled(m, { moduleOverrides: overrides }))
        .map((m) => toNavItem(m.nav!))
})

const configurationItems = computed<NavItem[]>(() => {
    const settingsModule = allModules.find((m) => m.id === 'settings')
    if (!settingsModule || !settingsModule.nav) return []

    // Return the settings module nav item itself (which contains children)
    return [toNavItem(settingsModule.nav)]
})

onMounted(() => {
    statusBarStore.refresh()
})
</script>

<template>
    <SidebarProvider>
        <AppSidebar :features="platformItems" :settings="configurationItems" :user="userData" @logout="handleLogout"
            @profile="handleSwitchUser" />

        <!-- Main Content Inset -->
        <SidebarInset>
            <!-- Header -->
            <header
                class="flex h-16 shrink-0 items-center gap-2 border-b px-4 transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-12">
                <div class="flex items-center gap-2 px-4">
                    <SidebarTrigger class="-ml-1" />
                    <Separator orientation="vertical" class="mr-2 h-4" />
                    <Breadcrumb>
                        <BreadcrumbList>
                            <BreadcrumbItem>
                                <BreadcrumbPage>
                                    {{ currentBreadcrumbTitle }}
                                </BreadcrumbPage>
                            </BreadcrumbItem>
                        </BreadcrumbList>
                    </Breadcrumb>
                </div>

                <div class="ml-auto flex items-center gap-4">
                    <!-- Language Switcher -->
                    <DropdownMenu>
                        <DropdownMenuTrigger as-child>
                            <Button variant="ghost" size="icon" class="rounded-full">
                                <Globe class="h-5 w-5" />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                            <DropdownMenuItem v-for="opt in localeOptions" :key="opt.key"
                                @click="handleLocaleSelect(opt.key as 'zh-CN' | 'en-US')">
                                {{ opt.label }}
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>

                    <!-- Theme Toggle -->
                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger as-child>
                                <Button variant="ghost" size="icon" class="rounded-full"
                                    @click="appStore.toggleTheme()">
                                    <Sun v-if="appStore.theme === 'dark'" class="h-5 w-5" />
                                    <Moon v-else class="h-5 w-5" />
                                </Button>
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>{{ appStore.theme === 'dark' ? t('theme.switchToLight') : t('theme.switchToDark') }}
                                </p>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>
                </div>
            </header>

            <!-- Main Scrollable Area -->
            <div class="flex-1 flex flex-col gap-4 p-4 pt-2 min-h-0 overflow-hidden">
                <RouterView />
            </div>

            <!-- Footer / Status Bar -->
            <footer class="h-8 border-t flex items-center px-6 text-xs text-muted-foreground bg-muted/40">
                <div class="flex items-center gap-3">
                    <span class="font-medium">{{ t('footer.statusBar') }}</span>
                    <span>{{ t('footer.monthlyHours') }} <strong class="text-primary">{{ statusBarStore.monthHoursLabel
                            }}</strong></span>
                    <span class="text-muted-foreground/40">|</span>
                    <span>{{ t('footer.uninvoiced') }} <strong class="text-primary">{{
                        statusBarStore.uninvoicedTotalLabel }}</strong></span>
                    <span class="text-muted-foreground/40">|</span>
                    <span>{{ t('footer.pendingPayment') }} <strong class="text-primary">{{
                        statusBarStore.unpaidTotalLabel }}</strong></span>
                </div>
            </footer>
        </SidebarInset>
    </SidebarProvider>
</template>

<style scoped>
/* Scoped styles removed in favor of Tailwind classes */
</style>
