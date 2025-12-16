<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent } from '@/components/ui/card'
import { Plus } from 'lucide-vue-next'
import { api } from '@/api'
import type { FinanceAccount } from '@/types/finance'
import { toast } from 'vue-sonner'

const { t } = useI18n()

const loading = ref(false)
const accounts = ref<FinanceAccount[]>([])

async function loadAccounts() {
  loading.value = true
  try {
    accounts.value = await api.finance.accounts.list()
  } catch (e) {
    toast.error(t('finance.accounts.errors.loadFailed'))
  } finally {
    loading.value = false
  }
}

function getTypeConfig(type: string) {
  const typeMap: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
    checking: { label: t('finance.accounts.types.checking'), variant: 'default' }, // info-like
    savings: { label: t('finance.accounts.types.savings'), variant: 'secondary' }, // success-like (approx)
    credit: { label: t('finance.accounts.types.credit'), variant: 'destructive' }, // warning-like (approx)
    investment: { label: t('finance.accounts.types.investment'), variant: 'outline' },
  }
  return typeMap[type] || { label: type, variant: 'outline' }
}

onMounted(() => {
  loadAccounts()
})
</script>

<template>
  <div class="accounts-page space-y-6">
    <div class="page-header flex justify-between items-center">
      <h1 class="page-title text-2xl font-bold">{{ t('finance.accounts.title') }}</h1>
      <Button>
        <Plus class="mr-2 h-4 w-4" />
        {{ t('finance.accounts.addAccount') }}
      </Button>
    </div>

    <Card>
      <CardContent class="p-0">
        <div class="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{{ t('finance.accounts.table.name') }}</TableHead>
                <TableHead class="w-[120px]">{{ t('finance.accounts.table.type') }}</TableHead>
                <TableHead class="w-[100px]">{{ t('finance.accounts.table.currency') }}</TableHead>
                <TableHead class="w-[150px] text-right">{{ t('finance.accounts.table.balance') }}</TableHead>
                <TableHead class="w-[150px]">{{ t('common.actions') }}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <template v-if="accounts.length > 0">
                <TableRow v-for="account in accounts" :key="account.id">
                  <TableCell>{{ account.name }}</TableCell>
                  <TableCell>
                    <Badge :variant="getTypeConfig(account.type).variant" class="rounded-full">
                      {{ getTypeConfig(account.type).label }}
                    </Badge>
                  </TableCell>
                  <TableCell>{{ account.currency }}</TableCell>
                  <TableCell class="text-right">${{ account.balance.toFixed(2) }}</TableCell>
                  <TableCell>
                    <div class="flex items-center gap-2">
                      <Button size="sm" variant="ghost" class="h-8 px-2 text-primary">
                        {{ t('common.edit') }}
                      </Button>
                      <Button size="sm" variant="ghost" class="h-8 px-2 text-destructive hover:text-destructive/90">
                        {{ t('common.delete') }}
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              </template>
              <TableRow v-else>
                <TableCell colspan="5" class="h-24 text-center">
                  {{ loading ? 'Loading...' : t('common.noData') }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<style scoped>
/* No extra styles needed with utility classes */
</style>
