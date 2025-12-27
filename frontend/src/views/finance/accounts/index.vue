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
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogFooter,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Plus } from 'lucide-vue-next'
import { api } from '@/api'
import type { FinanceAccount } from '@/types/finance'
import { toast } from 'vue-sonner'

const { t } = useI18n()

const loading = ref(false)
const accounts = ref<FinanceAccount[]>([])
const showCreateModal = ref(false)

const newAccount = ref({
  name: '',
  type: 'checking',
  currency: 'CAD',
  balance: 0,
  bankName: '',
})

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

async function createAccount() {
  try {
    await api.finance.accounts.create(newAccount.value)
    toast.success(t('common.saved'))
    showCreateModal.value = false
    loadAccounts()
    // Reset form
    newAccount.value = {
      name: '',
      type: 'checking',
      currency: 'CAD',
      balance: 0,
      bankName: '',
    }
  } catch (e) {
    toast.error(t('common.error'))
  }
}

function getTypeConfig(type: string) {
  const typeMap: Record<string, { label: string; variant: 'default' | 'secondary' | 'destructive' | 'outline' }> = {
    checking: { label: t('finance.accounts.types.checking'), variant: 'default' },
    savings: { label: t('finance.accounts.types.savings'), variant: 'secondary' },
    credit: { label: t('finance.accounts.types.credit'), variant: 'destructive' },
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
      <Dialog v-model:open="showCreateModal">
        <DialogTrigger as-child>
          <Button>
            <Plus class="mr-2 h-4 w-4" />
            {{ t('finance.accounts.addAccount') }}
          </Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>{{ t('finance.accounts.addAccount') }}</DialogTitle>
          </DialogHeader>
          <div class="grid gap-4 py-4">
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="name" class="text-right">{{ t('finance.accounts.table.name') }}</Label>
              <Input id="name" v-model="newAccount.name" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="type" class="text-right">{{ t('finance.accounts.table.type') }}</Label>
              <Select v-model="newAccount.type">
                <SelectTrigger class="col-span-3">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="checking">Checking</SelectItem>
                  <SelectItem value="savings">Savings</SelectItem>
                  <SelectItem value="credit">Credit Card</SelectItem>
                  <SelectItem value="investment">Investment</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="bank" class="text-right">Bank</Label>
              <Input id="bank" v-model="newAccount.bankName" class="col-span-3" placeholder="e.g. CIBC" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="currency" class="text-right">{{ t('finance.accounts.table.currency') }}</Label>
              <Input id="currency" v-model="newAccount.currency" class="col-span-3" />
            </div>
            <div class="grid grid-cols-4 items-center gap-4">
              <Label for="balance" class="text-right">{{ t('finance.accounts.table.balance') }}</Label>
              <Input id="balance" type="number" step="0.01" v-model.number="newAccount.balance" class="col-span-3" />
            </div>
          </div>
          <DialogFooter>
            <Button @click="createAccount">{{ t('common.save') }}</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>

    <Card>
      <CardContent class="p-0">
        <div class="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{{ t('finance.accounts.table.name') }}</TableHead>
                <TableHead class="w-[120px]">{{ t('finance.accounts.table.type') }}</TableHead>
                <TableHead>{{ t('finance.accounts.table.bankName') || 'Bank' }}</TableHead>
                <TableHead class="w-[100px]">{{ t('finance.accounts.table.currency') }}</TableHead>
                <TableHead class="w-[150px] text-right">{{ t('finance.accounts.table.balance') }}</TableHead>
                <TableHead class="w-[150px]">{{ t('common.actions') }}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <template v-if="accounts.length > 0">
                <TableRow v-for="account in accounts" :key="account.id">
                  <TableCell class="font-medium">{{ account.name }}</TableCell>
                  <TableCell>
                    <Badge :variant="getTypeConfig(account.type).variant" class="rounded-full">
                      {{ getTypeConfig(account.type).label }}
                    </Badge>
                  </TableCell>
                  <TableCell>{{ account.bankName }}</TableCell>
                  <TableCell>{{ account.currency }}</TableCell>
                  <TableCell class="text-right">${{ account.balance.toFixed(2) }}</TableCell>
                  <TableCell>
                    <div class="flex items-center gap-2">
                      <!-- Edit/Delete implementation can be added later -->
                      <Button size="sm" variant="ghost" class="h-8 px-2 text-primary">
                        {{ t('common.edit') }}
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              </template>
              <TableRow v-else>
                <TableCell colspan="6" class="h-24 text-center">
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
