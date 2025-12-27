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
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { api } from '@/api'
import { toast } from 'vue-sonner'
import { format } from 'date-fns'
import { Badge } from '@/components/ui/badge'

const { t } = useI18n()

const loading = ref(false)
const transactions = ref<any[]>([])
const accounts = ref<any[]>([])
const categories = ref<any[]>([])

const filter = ref({
  accountId: 0,
  startDate: '',
  endDate: '',
})

async function loadData() {
  loading.value = true
  try {
    const [txs, accs, cats] = await Promise.all([
      api.finance.transactions.list({
        accountId: filter.value.accountId || undefined,
        startDate: filter.value.startDate || undefined,
        endDate: filter.value.endDate || undefined
      }),
      api.finance.accounts.list(),
      api.finance.categories.list()
    ])
    transactions.value = txs
    accounts.value = accs
    categories.value = cats
  } catch (e) {
    toast.error('Failed to load transactions')
  } finally {
    loading.value = false
  }
}

async function updateCategory(txId: number, categoryId: string) {
  try {
    await api.finance.transactions.update(txId, parseInt(categoryId))
    toast.success(t('common.saved'))
  } catch(e) {
    toast.error(t('common.error'))
  }
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <div class="transactions-page space-y-6">
    <div class="page-header flex justify-between items-center">
      <h1 class="page-title text-2xl font-bold">{{ t('finance.transactions.title') }}</h1>
      <div class="flex gap-2">
        <router-link to="/finance/import">
           <Button variant="outline">{{ t('finance.nav.import') }}</Button>
        </router-link>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex gap-4 items-end bg-card p-4 rounded-lg border">
      <div class="grid gap-2">
        <label class="text-sm font-medium">Account</label>
        <Select v-model="filter.accountId" :model-value="filter.accountId?.toString()" @update:model-value="(v) => { filter.accountId = Number(v); loadData() }">
          <SelectTrigger class="w-[200px]">
            <SelectValue placeholder="All Accounts" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="0">All Accounts</SelectItem>
            <SelectItem v-for="acc in accounts" :key="acc.id" :value="acc.id.toString()">
              {{ acc.name }}
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
      <div class="grid gap-2">
        <label class="text-sm font-medium">Date Range</label>
        <div class="flex gap-2">
           <Input type="date" v-model="filter.startDate" @change="loadData" />
           <Input type="date" v-model="filter.endDate" @change="loadData" />
        </div>
      </div>
      <Button variant="secondary" @click="loadData">Refresh</Button>
    </div>

    <Card>
      <CardContent class="p-0">
        <div class="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Date</TableHead>
                <TableHead>Description</TableHead>
                <TableHead>Amount</TableHead>
                <TableHead>Category</TableHead>
                <TableHead>Status</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <template v-if="transactions.length > 0">
                <TableRow v-for="tx in transactions" :key="tx.id">
                  <TableCell>{{ format(new Date(tx.date), 'yyyy-MM-dd') }}</TableCell>
                  <TableCell>{{ tx.description }}</TableCell>
                  <TableCell :class="tx.amount < 0 ? 'text-red-500' : 'text-green-500'">
                    {{ tx.amount.toFixed(2) }}
                  </TableCell>
                  <TableCell>
                    <Select :model-value="tx.categoryId?.toString()" @update:model-value="(v) => updateCategory(tx.id, v)">
                      <SelectTrigger class="w-[160px] h-8">
                        <SelectValue placeholder="Uncategorized" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem v-for="cat in categories" :key="cat.id" :value="cat.id.toString()">
                          <div class="flex items-center gap-2">
                             <div class="w-3 h-3 rounded-full" :style="{backgroundColor: cat.color}"></div>
                             {{ cat.name }}
                          </div>
                        </SelectItem>
                      </SelectContent>
                    </Select>
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline">{{ tx.status }}</Badge>
                  </TableCell>
                </TableRow>
              </template>
              <TableRow v-else>
                <TableCell colspan="5" class="h-24 text-center">
                  {{ loading ? t('common.loading') : t('common.noData') }}
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  </div>
</template>
