<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Wallet,
  TrendingUp,
  TrendingDown,
  DollarSign,
  Plus
} from 'lucide-vue-next'
import { api } from '@/api'
import type { FinanceSummary } from '@/types/finance'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'


const { t } = useI18n()

const loading = ref(false)
const summary = ref<FinanceSummary | null>(null)

async function loadData() {
  loading.value = true
  try {
    summary.value = await api.finance.summary.get()
  } catch (e) {
    console.error('Failed to load finance summary:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="h-full flex flex-col min-h-0">
    <div class="flex-1 min-h-0 overflow-auto">
      <div class="finance-overview space-y-6 p-1">
        <!-- Header -->
        <div class="overview-header flex justify-between items-start">
          <div>
            <h1 class="page-title text-3xl font-bold tracking-tight mb-2">{{ t('finance.overview.title') }}</h1>
            <p class="page-subtitle text-muted-foreground">{{ t('finance.overview.subtitle') }}</p>
          </div>
          <div class="flex gap-2">
            <Button variant="outline">
              <Plus class="w-4 h-4 mr-2" />
              {{ t('finance.actions.addTransaction') }}
            </Button>
            <Button>
              <Wallet class="w-4 h-4 mr-2" />
              {{ t('finance.actions.addAccount') }}
            </Button>
          </div>
        </div>

        <!-- Summary Cards -->
        <div class="summary-section grid gap-4 md:grid-cols-2 lg:grid-cols-4">
          <Card class="hover:shadow-md transition-shadow">
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">
                {{ t('finance.summary.totalBalance') }}
              </CardTitle>
              <div class="icon-box p-2 rounded-md bg-[#8D7B68]/10 text-[#8D7B68]">
                <Wallet class="w-4 h-4" />
              </div>
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">${{ summary?.totalBalance.toFixed(2) || '0.00' }}</div>
            </CardContent>
          </Card>

          <Card class="hover:shadow-md transition-shadow">
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">
                {{ t('finance.summary.totalIncome') }}
              </CardTitle>
              <div
                class="icon-box p-2 rounded-md bg-emerald-100 text-emerald-600 dark:bg-emerald-900/30 dark:text-emerald-400">
                <TrendingUp class="w-4 h-4" />
              </div>
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">${{ summary?.totalIncome.toFixed(2) || '0.00' }}</div>
            </CardContent>
          </Card>

          <Card class="hover:shadow-md transition-shadow">
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">
                {{ t('finance.summary.totalExpenses') }}
              </CardTitle>
              <div class="icon-box p-2 rounded-md bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400">
                <TrendingDown class="w-4 h-4" />
              </div>
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">${{ summary?.totalExpense.toFixed(2) || '0.00' }}</div>
            </CardContent>
          </Card>

          <Card class="hover:shadow-md transition-shadow">
            <CardHeader class="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle class="text-sm font-medium">
                {{ t('finance.summary.netCashFlow') }}
              </CardTitle>
              <div class="icon-box p-2 rounded-md bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400">
                <DollarSign class="w-4 h-4" />
              </div>
            </CardHeader>
            <CardContent>
              <div class="text-2xl font-bold">${{ summary?.cashFlow.toFixed(2) || '0.00' }}</div>
            </CardContent>
          </Card>
        </div>

        <!-- Quick Actions -->
        <div class="actions-section">
          <Card>
            <CardHeader>
              <CardTitle>{{ t('finance.overview.quickActions') }}</CardTitle>
            </CardHeader>
            <CardContent class="p-0">
              <div class="flex flex-col">
                <button
                  class="flex items-center w-full p-4 hover:bg-muted/50 transition-colors text-left border-b last:border-0">
                  <div class="flex-1">
                    <div class="font-medium text-sm">{{ t('finance.overview.viewAllAccounts') }}</div>
                    <div class="text-xs text-muted-foreground">{{ t('finance.overview.manageYourAccounts') }}</div>
                  </div>
                </button>
                <button
                  class="flex items-center w-full p-4 hover:bg-muted/50 transition-colors text-left border-b last:border-0">
                  <div class="flex-1">
                    <div class="font-medium text-sm">{{ t('finance.overview.viewAllTransactions') }}</div>
                    <div class="text-xs text-muted-foreground">{{ t('finance.overview.browseRecentTransactions') }}
                    </div>
                  </div>
                </button>
                <button
                  class="flex items-center w-full p-4 hover:bg-muted/50 transition-colors text-left border-b last:border-0">
                  <div class="flex-1">
                    <div class="font-medium text-sm">{{ t('finance.overview.importStatement') }}</div>
                    <div class="text-xs text-muted-foreground">{{ t('finance.overview.importBankStatement') }}</div>
                  </div>
                </button>
                <button class="flex items-center w-full p-4 hover:bg-muted/50 transition-colors text-left">
                  <div class="flex-1">
                    <div class="font-medium text-sm">{{ t('finance.overview.viewReports') }}</div>
                    <div class="text-xs text-muted-foreground">{{ t('finance.overview.analyzeYourFinance') }}</div>
                  </div>
                </button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Scoped styles replaced by Tailwind classes */
</style>
