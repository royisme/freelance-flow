<script setup lang="ts">
import { SidebarMenu, SidebarMenuItem, SidebarMenuButton } from '@/components/ui/sidebar'
import { GalleryVerticalEnd } from 'lucide-vue-next'

defineProps<{
    title?: string
    subtitle?: string
}>()
</script>

<template>
    <SidebarMenu>
        <SidebarMenuItem>
            <SidebarMenuButton size="lg" as-child>
                <a href="#">
                    <div
                        class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                        <GalleryVerticalEnd class="size-4" />
                    </div>
                    <div class="grid flex-1 text-left text-sm leading-tight">
                        <span class="truncate font-semibold">{{ title || 'FreelanceFlow' }}</span>
                        <span class="truncate text-xs">{{ subtitle || 'Premium' }}</span>
                    </div>
                </a>
            </SidebarMenuButton>
        </SidebarMenuItem>
    </SidebarMenu>
</template>
