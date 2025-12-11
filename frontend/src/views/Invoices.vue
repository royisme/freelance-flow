<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import {
  NButton, NDataTable, NTag, NSpace, NText, NNumberAnimation, NStatistic, NCard,
  NModal,
  type DataTableColumns, useMessage
} from 'naive-ui'
import PageContainer from '@/components/PageContainer.vue'
import InvoiceFormModal from '@/components/InvoiceFormModal.vue'
import { useInvoiceStore, type EnrichedInvoice } from '@/stores/invoices'
import { useClientStore } from '@/stores/clients'
import { useTimesheetStore } from '@/stores/timesheet'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { Invoice } from '@/types'
import { PlusOutlined, DownloadOutlined, FileTextOutlined, DollarOutlined, MailOutlined } from '@vicons/antd'

const message = useMessage()
const invoiceStore = useInvoiceStore()
const clientStore = useClientStore()
const timesheetStore = useTimesheetStore()
const { enrichedInvoices, stats, loading } = storeToRefs(invoiceStore)
const { clients } = storeToRefs(clientStore)
const { enrichedEntries, loading: timesLoading } = storeToRefs(timesheetStore)
const { t } = useI18n()

const showModal = ref(false)
const editingInvoice = ref<Invoice | null>(null)
const entrySelectorVisible = ref(false)
const entrySelection = ref<number[]>([])
const activeInvoiceId = ref<number | null>(null)
const pdfLoading = ref(false)
const sendLoading = ref(false)

function handleNewInvoice() {
  editingInvoice.value = null
  showModal.value = true
}

function handleEditInvoice(invoice: Invoice) {
  editingInvoice.value = invoice
  showModal.value = true
}

function openEntrySelector(invoice: EnrichedInvoice) {
  activeInvoiceId.value = invoice.id
  entrySelection.value = enrichedEntries.value
    .filter((e) => e.invoiceId === invoice.id)
    .map((e) => e.id)
  entrySelectorVisible.value = true
}

async function applyEntrySelection() {
  if (!activeInvoiceId.value) return
  await invoiceStore.setTimeEntries(activeInvoiceId.value, entrySelection.value)
  entrySelectorVisible.value = false
  message.success(t('invoices.entriesUpdated'))
}

async function handleDownload(invoice: EnrichedInvoice) {
  try {
    pdfLoading.value = true
    const base64 = await invoiceStore.generatePdf(invoice.id)
    const bytes = Uint8Array.from(atob(base64), (c) => c.charCodeAt(0))
    const url = URL.createObjectURL(new Blob([bytes], { type: 'application/pdf' }))
    const a = document.createElement('a')
    a.href = url
    a.download = `INV-${invoice.number}.pdf`
    a.click()
    URL.revokeObjectURL(url)
    message.success(t('invoices.downloaded'))
  } catch {
    message.error(t('invoices.downloadError'))
  } finally {
    pdfLoading.value = false
  }
}

async function handleSend(invoice: EnrichedInvoice) {
  try {
    sendLoading.value = true
    const ok = await invoiceStore.sendEmail(invoice.id)
    if (ok) {
      message.success(t('invoices.sendSuccess'))
    } else {
      message.error(t('invoices.sendError'))
    }
  } catch {
    message.error(t('invoices.sendError'))
  } finally {
    sendLoading.value = false
  }
}

async function handleSubmitInvoice(invoice: Omit<Invoice, 'id'> | Invoice) {
  try {
    if ('id' in invoice) {
      await invoiceStore.updateInvoice(invoice)
      message.success(t('invoices.updateSuccess'))
    } else {
      await invoiceStore.createInvoice(invoice)
      message.success(t('invoices.createSuccess'))
    }
  } catch (error) {
    message.error(t('invoices.saveError'))
  }
}

onMounted(() => {
  invoiceStore.fetchInvoices()
  clientStore.fetchClients()
  timesheetStore.fetchTimesheet()
})

const columns: DataTableColumns<EnrichedInvoice> = [
  {
    title: () => t('invoices.columns.invoiceNumber'),
    key: 'number',
    width: 150,
    render(row) {
      return h(NText, { strong: true }, { default: () => row.number })
    }
  },
  {
    title: () => t('invoices.columns.client'),
    key: 'clientName',
    width: 200,
  },
  {
    title: () => t('invoices.columns.issueDate'),
    key: 'issueDate',
    width: 140
  },
  {
    title: () => t('invoices.columns.amount'),
    key: 'total',
    render(row) {
      return h(
        NText,
        { type: 'default', style: 'font-variant-numeric: tabular-nums;' },
        { default: () => `${row.clientCurrency} ${row.total.toLocaleString()}` }
      )
    }
  },
  {
    title: () => t('invoices.columns.status'),
    key: 'status',
    width: 120,
    render(row) {
      let type: 'success' | 'warning' | 'error' | 'default' = 'default'
      if (row.status === 'paid') type = 'success'
      if (row.status === 'sent') type = 'warning'
      if (row.status === 'overdue') type = 'error'

      return h(
        NTag,
        { type, bordered: false, round: true, size: 'small' },
        { default: () => t(`invoices.status.${row.status}`) }
      )
    }
  },
  {
    title: () => t('invoices.columns.actions'),
    key: 'actions',
    width: 140,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              circle: true,
              loading: pdfLoading.value,
              onClick: () => handleDownload(row)
            },
            { icon: () => h(DownloadOutlined) }
          ),
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              circle: true,
              onClick: () => openEntrySelector(row)
            },
            { icon: () => h(FileTextOutlined) }
          ),
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              circle: true,
              loading: sendLoading.value,
              onClick: () => handleSend(row)
            },
            { icon: () => h(MailOutlined) }
          )
        ]
      })
    }
  }
]
</script>

<template>
  <PageContainer :title="t('invoices.title')" :subtitle="t('invoices.subtitle')">
    <template #extra>
      <n-button type="primary" @click="handleNewInvoice">
        <template #icon>
          <n-icon>
            <PlusOutlined />
          </n-icon>
        </template>
        {{ t('invoices.createInvoice') }}
      </n-button>
    </template>

    <InvoiceFormModal v-model:show="showModal" :invoice="editingInvoice" :clients="clients"
      @submit="handleSubmitInvoice" />

    <template #headerContent>
      <n-space size="large" style="margin-top: 16px;">
        <n-card size="small" :bordered="false" class="stat-card">
          <n-statistic :label="t('invoices.stats.outstandingAmount')">
            <template #prefix>
              <n-icon>
                <DollarOutlined />
              </n-icon>
            </template>
            <n-number-animation :from="0" :to="stats.totalDue" :precision="2" />
          </n-statistic>
        </n-card>
      </n-space>
    </template>

    <n-data-table :columns="columns" :data="enrichedInvoices" :loading="loading" :bordered="false"
      class="invoice-table" />

  <n-modal v-model:show="entrySelectorVisible" preset="dialog" title="Select Time Entries" style="width: 720px">
    <n-data-table
      :loading="timesLoading"
      :columns="[
        { title: 'Date', key: 'date' },
        { title: 'Project', key: 'project', render: (row: any) => row.project?.name || '-' },
        { title: 'Hours', key: 'hours', render: (row: any) => (row.durationSeconds / 3600).toFixed(2) },
        { title: 'Linked', key: 'linked', render: (row: any) => (row.invoiceId ? '✔' : '') }
      ]"
      :data="enrichedEntries.filter((e) => !activeInvoiceId || e.invoiceId === activeInvoiceId || !e.invoiceId)"
      :row-key="(row: any) => row.id"
      checkable
      :checked-row-keys="entrySelection"
      @update:checked-row-keys="(keys: number[]) => entrySelection = keys"
    />
    <template #action>
      <n-space justify="end">
        <n-button quaternary @click="entrySelectorVisible = false">Cancel</n-button>
        <n-button type="primary" :loading="loading" @click="applyEntrySelection">Apply</n-button>
      </n-space>
    </template>
  </n-modal>
  </PageContainer>
</template>

<style scoped>
.invoice-table :deep(.n-data-table-th) {
  font-weight: 600;
  color: var(--n-text-color-2);
}

.stat-card {
  background: var(--n-card-color);
  min-width: 200px;
}
</style>
