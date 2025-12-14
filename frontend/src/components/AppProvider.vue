<script setup lang="ts">
import { computed, watch } from 'vue'
import { NConfigProvider, NMessageProvider, NDialogProvider, NNotificationProvider, NGlobalStyle, darkTheme } from 'naive-ui'
import { lightThemeOverrides, darkThemeOverrides } from '@/theme'
import { zhCN, dateZhCN, enUS, dateEnUS } from 'naive-ui'
import { useAppStore } from '@/stores/app'
import { useI18n } from 'vue-i18n'
import { applyThemeToRoot } from '@/theme/tokens'

const appStore = useAppStore()
const { locale } = useI18n()

// Theme Configuration
const theme = computed(() => appStore.theme === 'dark' ? darkTheme : null)
const themeOverrides = computed(() => appStore.theme === 'dark' ? darkThemeOverrides : lightThemeOverrides)

// Locale Configuration
const naiveLocale = computed(() => appStore.locale === 'zh-CN' ? zhCN : enUS)
const naiveDateLocale = computed(() => appStore.locale === 'zh-CN' ? dateZhCN : dateEnUS)

watch(
  () => appStore.theme,
  (newTheme) => {
    applyThemeToRoot(newTheme)
  },
  { immediate: true }
)

// Watch for locale changes to update vue-i18n
watch(
  () => appStore.locale,
  (newLocale) => {
    if (newLocale !== locale.value) {
      locale.value = newLocale
    }
  },
  { immediate: true }
)
</script>

<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides" :locale="naiveLocale"
    :date-locale="naiveDateLocale">
    <n-global-style />
    <n-notification-provider>
      <n-message-provider>
        <n-dialog-provider>
          <slot />
        </n-dialog-provider>
      </n-message-provider>
    </n-notification-provider>
  </n-config-provider>
</template>
