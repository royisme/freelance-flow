<script setup lang="ts">
import { Sidebar, SidebarContent, SidebarFooter, SidebarHeader, SidebarRail } from '@/components/ui/sidebar'
import AppSidebarHeader from './AppSidebarHeader.vue'
import AppSidebarNav, { type NavItem } from './AppSidebarNav.vue'
import AppSidebarFooter from './AppSidebarFooter.vue'

interface UserData {
    name: string
    email: string
    avatar: string
}

defineProps<{
    features: NavItem[]
    settings: NavItem[]
    user: UserData
}>()

const emit = defineEmits<{
    (e: 'logout'): void
    (e: 'profile'): void
}>()
</script>

<template>
    <Sidebar collapsible="icon">
        <SidebarHeader>
            <AppSidebarHeader />
        </SidebarHeader>
        <SidebarContent>
            <AppSidebarNav :features="features" :settings="settings" />
        </SidebarContent>
        <SidebarFooter>
            <AppSidebarFooter :user="user" @logout="emit('logout')" @profile="emit('profile')" />
        </SidebarFooter>
        <SidebarRail />
    </Sidebar>
</template>
