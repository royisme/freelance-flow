<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { api } from '@/api'
import { toast } from 'vue-sonner'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const router = useRouter()

const accounts = ref<any[]>([])
const selectedAccountId = ref<string>('')
const selectedBank = ref<string>('CIBC')
const fileContent = ref<string | null>(null)
const fileName = ref<string>('')
const loading = ref(false)

async function loadAccounts() {
  try {
    accounts.value = await api.finance.accounts.list()
  } catch (e) {
    toast.error('Failed to load accounts')
  }
}

function handleFileUpload(event: Event) {
  const target = event.target as HTMLInputElement
  if (target.files && target.files.length > 0) {
    const file = target.files[0]
    fileName.value = file.name
    const reader = new FileReader()
    reader.onload = (e) => {
      fileContent.value = e.target?.result as string
    }
    reader.readAsText(file)
  }
}

async function doImport() {
  if (!selectedAccountId.value || !fileContent.value) {
    toast.error('Please select an account and a file')
    return
  }

  loading.value = true
  try {
    const count = await api.finance.transactions.import({
      accountId: parseInt(selectedAccountId.value),
      bankType: selectedBank.value,
      fileContent: fileContent.value
    })
    toast.success(`Imported ${count} transactions`)
    router.push('/finance/transactions')
  } catch (e) {
    console.error(e)
    toast.error('Import failed. Please check the file format.')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadAccounts()
})
</script>

<template>
  <div class="import-page max-w-2xl mx-auto">
    <Card>
      <CardHeader>
        <CardTitle>{{ t('finance.import.title') }}</CardTitle>
      </CardHeader>
      <CardContent class="space-y-6">

        <div class="space-y-2">
          <Label>1. Select Target Account</Label>
          <Select v-model="selectedAccountId">
            <SelectTrigger>
              <SelectValue placeholder="Select Account" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="acc in accounts" :key="acc.id" :value="acc.id.toString()">
                {{ acc.name }} ({{ acc.currency }})
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div class="space-y-2">
          <Label>2. Select Bank Format</Label>
          <Select v-model="selectedBank">
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="CIBC">CIBC (CSV)</SelectItem>
              <SelectItem value="RBC">RBC (CSV)</SelectItem>
              <SelectItem value="TD">TD (CSV)</SelectItem>
              <SelectItem value="GENERIC">Generic (Date,Desc,Amount)</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div class="space-y-2">
          <Label>3. Upload File</Label>
          <div class="flex items-center gap-4">
            <Input type="file" accept=".csv" @change="handleFileUpload" />
          </div>
          <p class="text-sm text-muted-foreground" v-if="fileName">Selected: {{ fileName }}</p>
        </div>

        <Button class="w-full" @click="doImport" :disabled="loading || !fileContent || !selectedAccountId">
          {{ loading ? 'Importing...' : 'Import Transactions' }}
        </Button>

      </CardContent>
    </Card>
  </div>
</template>
