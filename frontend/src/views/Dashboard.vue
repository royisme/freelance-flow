<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useDashboardStore } from '@/stores/dashboard'
import { useAuthStore } from '@/stores/auth'
import { Clock, DollarSign, TrendingUp, Calendar } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import { Card, CardContent } from '@/components/ui/card'

const { t } = useI18n()
const store = useDashboardStore()
const authStore = useAuthStore()

const username = computed(() => authStore.currentUser?.username || 'User')

onMounted(() => {
  store.fetchDashboardData()
})
</script>

<template>
  <div class="h-full flex flex-col min-h-0 gap-4">
    <!-- Greeting/Header Section -->
    <div class="shrink-0">
      <h2 class="text-2xl font-bold tracking-tight">{{ t('dashboard.greeting', { name: username }) }}</h2>
      <p class="text-muted-foreground text-sm">{{ t('dashboard.weekOverview') }}</p>
    </div>

    <!-- Key Metrics Cards -->
    <div class="shrink-0 grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Weekly Hours -->
      <Card class="hover:-translate-y-1 transition-all duration-300 shadow-sm hover:shadow-md">
        <CardContent class="p-4 flex items-center gap-4">
          <div
            class="h-10 w-10 rounded-full bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center shrink-0">
            <Clock class="h-5 w-5 text-orange-600 dark:text-orange-400" />
          </div>
          <div>
            <p class="text-sm font-medium text-muted-foreground">{{ t('dashboard.metrics.weeklyHours') }}</p>
            <div class="flex items-baseline gap-1">
              <span class="text-xl font-bold tracking-tight">{{ store.totalHoursWeek.toFixed(1) }}</span>
              <span class="text-sm text-muted-foreground">{{ t('dashboard.metrics.hoursUnit') }}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Monthly Revenue -->
      <Card class="hover:-translate-y-1 transition-all duration-300 shadow-sm hover:shadow-md">
        <CardContent class="p-4 flex items-center gap-4">
          <div
            class="h-10 w-10 rounded-full bg-green-100 dark:bg-green-900/30 flex items-center justify-center shrink-0">
            <DollarSign class="h-5 w-5 text-green-600 dark:text-green-400" />
          </div>
          <div>
            <p class="text-sm font-medium text-muted-foreground">{{ t('dashboard.metrics.monthlyRevenue') }}</p>
            <div class="flex items-baseline gap-1">
              <span class="text-xl font-bold tracking-tight">{{ store.totalRevenueMonth.toLocaleString() }}</span>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Pending Amount -->
      <Card class="hover:-translate-y-1 transition-all duration-300 shadow-sm hover:shadow-md">
        <CardContent class="p-4 flex items-center gap-4">
          <div class="h-10 w-10 rounded-full bg-rose-100 dark:bg-rose-900/30 flex items-center justify-center shrink-0">
            <TrendingUp class="h-5 w-5 text-rose-600 dark:text-rose-400" />
          </div>
          <div>
            <p class="text-sm font-medium text-muted-foreground">{{ t('dashboard.metrics.pendingAmount') }}</p>
            <div class="flex items-baseline gap-1">
              <span class="text-xl font-bold tracking-tight">{{ store.pendingAmount.toLocaleString() }}</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Recent Activity Section -->
    <div class="flex-1 min-h-0 flex flex-col">
      <div class="shrink-0 flex items-center justify-between mb-2">
        <h2 class="text-lg font-semibold tracking-tight">{{ t('dashboard.recentActivity.title') }}</h2>
      </div>

      <Card class="flex-1 min-h-0 flex flex-col overflow-hidden">
        <CardContent class="flex-1 min-h-0 p-0 overflow-auto">
          <div class="p-4 space-y-1">
            <template v-if="store.recentActivities.length > 0">
              <div v-for="activity in store.recentActivities" :key="activity.id"
                class="group flex items-center gap-4 p-3 rounded-lg hover:bg-muted/50 transition-colors cursor-pointer">
                <div
                  class="h-9 w-9 rounded-full bg-orange-100 dark:bg-orange-900/20 flex items-center justify-center shrink-0">
                  <Clock class="h-4.5 w-4.5 text-orange-600 dark:text-orange-400" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between">
                    <p class="font-medium truncate">{{ activity.project }}</p>
                    <span
                      class="text-xs font-medium px-2.5 py-0.5 rounded-full bg-primary/10 text-primary whitespace-nowrap">
                      {{ t('dashboard.recentActivity.hoursLabel', { hours: activity.hours }) }}
                    </span>
                  </div>
                  <p class="text-sm text-muted-foreground truncate">{{ activity.date }} · {{ activity.description }}</p>
                </div>
              </div>
            </template>

            <div v-else class="flex flex-col items-center justify-center text-muted-foreground p-8">
              <Calendar class="h-12 w-12 mb-4 opacity-20" />
              <p>{{ t('dashboard.recentActivity.empty') }}</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>
