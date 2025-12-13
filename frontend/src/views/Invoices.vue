<script setup lang="ts">
import { h, onMounted, ref, computed } from 'vue'
import {
  NButton, NDataTable, NSpace, NText, NNumberAnimation, NIcon,
  NModal, NInput, NAvatar, NRow, NCol, NEmpty, NStatistic,
  type DataTableColumns, useMessage
} from 'naive-ui'
import PageContainer from '@/components/PageContainer.vue'
import InvoiceFormModal from '@/components/InvoiceFormModal.vue'
import { useInvoiceStore, type EnrichedInvoice } from '@/stores/invoices'
import { useClientStore } from '@/stores/clients'
import { useTimesheetStore } from '@/stores/timesheet'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { Invoice, TimeEntry, Project } from '@/types'
import {
  PlusOutlined,
  DownloadOutlined,
  FileTextOutlined,
  DollarOutlined,
  MailOutlined,
  SearchOutlined
} from '@vicons/antd'

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
const messageModalVisible = ref(false)
const messageDraft = ref("")
const exportingInvoice = ref<EnrichedInvoice | null>(null)
const searchQuery = ref('')
const statusFilter = ref<string | null>(null)

type EntryRow = TimeEntry & { project?: Project }

// Computed
const filteredInvoices = computed(() => {
  let result = enrichedInvoices.value

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(inv =>
      inv.number.toLowerCase().includes(q) ||
      (inv.clientName && inv.clientName.toLowerCase().includes(q))
    )
  }

  if (statusFilter.value) {
    result = result.filter(inv => inv.status === statusFilter.value)
  }

  return result
})

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
    exportingInvoice.value = invoice
    messageDraft.value = await invoiceStore.getDefaultMessage(invoice.id)
    messageModalVisible.value = true
  } catch {
    message.error(t('invoices.downloadError'))
  } finally {
    pdfLoading.value = false
  }
}

async function confirmDownload() {
  if (!exportingInvoice.value) return
  try {
    pdfLoading.value = true
    const base64 = await invoiceStore.generatePdf(
      exportingInvoice.value.id,
      messageDraft.value
    )
    const bytes = Uint8Array.from(atob(base64), (c) => c.charCodeAt(0))
    const url = URL.createObjectURL(new Blob([bytes], { type: 'application/pdf' }))
    const a = document.createElement('a')
    a.href = url
    a.download = `INV-${exportingInvoice.value.number}.pdf`
    a.click()
    URL.revokeObjectURL(url)
    message.success(t('invoices.downloaded'))
    messageModalVisible.value = false
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
    width: 140,
    render(row) {
      return h(
        'div',
        {
          class: 'invoice-number-cell',
          onClick: () => handleEditInvoice(row)
        },
        [
          h(NIcon, { size: 16, component: FileTextOutlined, style: 'margin-right: 8px; color: var(--n-text-color-3);' }),
          h(NText, { strong: true, style: 'cursor: pointer; transition: color 0.2s;' }, { default: () => row.number })
        ]
      )
    }
  },
  {
    title: () => t('invoices.columns.client'),
    key: 'clientName',
    width: 250,
    render(row) {
      const initial = row.clientName ? row.clientName.charAt(0).toUpperCase() : '?'
      return h('div', { style: 'display: flex; align-items: center; gap: 12px;' }, [
        h(NAvatar, {
          size: 32,
          round: true,
          style: 'background-color: var(--n-primary-color-suppl); color: var(--n-primary-color); font-weight: 600;'
        }, { default: () => initial }),
        h(NText, { style: 'font-weight: 500;' }, { default: () => row.clientName })
      ])
    }
  },
  {
    title: () => t('invoices.columns.issueDate'),
    key: 'issueDate',
    width: 150,
    render(row) {
      return h(NText, { depth: 3 }, { default: () => row.issueDate })
    }
  },
  {
    title: () => t('invoices.columns.amount'),
    key: 'total',
    align: 'right',
    width: 150,
    render(row) {
      return h(
        NText,
        { strong: true, style: 'font-variant-numeric: tabular-nums; font-size: 14px;' },
        { default: () => `${row.clientCurrency} ${row.total.toLocaleString(undefined, { minimumFractionDigits: 2 })}` }
      )
    }
  },
  {
    title: () => t('invoices.columns.status'),
    key: 'status',
    width: 140,
    align: 'center',
    render(row) {
      let color = ''
      if (row.status === 'paid') { color = '#18a058' }
      else if (row.status === 'sent') { color = '#f0a020' }
      else if (row.status === 'overdue') { color = '#d03050' }
      else { color = '#666' }

      // Custom dot-style tag
      return h(
        'div',
        {
          style: {
            display: 'inline-flex',
            alignItems: 'center',
            padding: '4px 10px',
            borderRadius: '12px',
            backgroundColor: `${color}15`, // 15% opacity
            color: color,
            fontSize: '12px',
            fontWeight: '600'
          }
        },
        [
          h('div', {
            style: {
              width: '6px',
              height: '6px',
              borderRadius: '50%',
              backgroundColor: color,
              marginRight: '6px'
            }
          }),
          t(`invoices.status.${row.status}`).toUpperCase()
        ]
      )
    }
  },
  {
    title: '',
    key: 'actions',
    width: 120,
    align: 'right',
    render(row) {
      return h(NSpace, { size: 'small', justify: 'end' }, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              circle: true,
              class: 'action-btn',
              loading: pdfLoading.value && exportingInvoice.value?.id === row.id,
              onClick: (e) => { e.stopPropagation(); handleDownload(row); }
            },
            { icon: () => h(DownloadOutlined) }
          ),
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              circle: true,
              class: 'action-btn',
              loading: sendLoading.value,
              onClick: (e) => { e.stopPropagation(); handleSend(row); }
            },
            { icon: () => h(MailOutlined) }
          ),
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              circle: true,
              class: 'action-btn',
              onClick: (e) => { e.stopPropagation(); openEntrySelector(row); }
            },
            { icon: () => h(FileTextOutlined) }
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
      <n-button type="primary" class="create-btn" @click="handleNewInvoice">
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

    <!-- Stats Grid moved to scrollable content -->
    <div class="view-content">
      <!-- Stats Grid -->
      <div class="stats-grid-container">
        <n-row :gutter="24">
          <n-col :span="12">
            <div class="stat-card primary">
              <div class="stat-icon">
                <n-icon>
                  <DollarOutlined />
                </n-icon>
              </div>
              <n-statistic :label="t('invoices.stats.outstandingAmount')">
                <template #default>
                  <n-number-animation :from="0" :to="stats.totalDue" :precision="2" show-separator />
                </template>
                <template #suffix>
                  <span class="currency">USD</span>
                </template>
              </n-statistic>
            </div>
          </n-col>
          <n-col :span="12">
            <div class="stat-card secondary">
              <div class="stat-icon secondary">
                <n-icon>
                  <FileTextOutlined />
                </n-icon>
              </div>
              <n-statistic :label="t('invoices.stats.totalInvoices')">
                <n-number-animation :from="0" :to="enrichedInvoices.length" />
              </n-statistic>
            </div>
          </n-col>
        </n-row>
      </div>

      <div class="content-wrapper">
        <!-- Search and Filter Bar -->
        <div class="toolbar">
          <n-input v-model:value="searchQuery" :placeholder="t('invoices.searchPlaceholder')" class="search-input">
            <template #prefix>
              <n-icon :component="SearchOutlined" />
            </template>
          </n-input>
          <div class="filters">
            <n-button quaternary size="small" :type="statusFilter === null ? 'primary' : 'default'"
              @click="statusFilter = null">
              {{ t('invoices.filter.all') }}
            </n-button>
            <n-button quaternary size="small" :type="statusFilter === 'draft' ? 'primary' : 'default'"
              @click="statusFilter = 'draft'">
              {{ t('invoices.filter.draft') }}
            </n-button>
            <n-button quaternary size="small" :type="statusFilter === 'sent' ? 'primary' : 'default'"
              @click="statusFilter = 'sent'">
              {{ t('invoices.filter.sent') }}
            </n-button>
          </div>
        </div>

        <n-data-table :columns="columns" :data="filteredInvoices" :loading="loading" :bordered="true"
          :row-class-name="() => 'invoice-row'" class="invoice-table">
          <template #empty>
            <div class="empty-state">
              <n-empty :description="t('invoices.empty.description')" size="large">
                <template #extra>
                  <n-button dashed size="small" @click="handleNewInvoice">{{ t('invoices.empty.action') }}</n-button>
                </template>
              </n-empty>
            </div>
          </template>
        </n-data-table>
      </div>
    </div>

    <!-- Modals -->
    <n-modal v-model:show="entrySelectorVisible" preset="card" :title="t('invoices.selectEntries.title')"
      style="width: 720px" :segmented="true">
      <n-data-table :loading="timesLoading" :columns="[
        { title: t('invoices.selectEntries.columns.date'), key: 'date' },
        { title: t('invoices.selectEntries.columns.project'), key: 'project', render: (row: EntryRow) => row.project?.name || '-' },
        { title: t('invoices.selectEntries.columns.hours'), key: 'hours', render: (row: EntryRow) => (row.durationSeconds / 3600).toFixed(2) },
        { title: t('invoices.selectEntries.columns.linked'), key: 'linked', render: (row: EntryRow) => (row.invoiceId ? '✔' : '') }
      ]" :data="enrichedEntries.filter((e) => !activeInvoiceId || e.invoiceId === activeInvoiceId || !e.invoiceId)"
        :row-key="(row: EntryRow) => row.id" checkable :checked-row-keys="entrySelection"
        @update:checked-row-keys="(keys) => entrySelection = keys as number[]" :max-height="400" />
      <template #footer>
        <n-space justify="end">
          <n-button quaternary @click="entrySelectorVisible = false">{{ t('invoices.selectEntries.cancel') }}</n-button>
          <n-button type="primary" :loading="loading" @click="applyEntrySelection">{{ t('invoices.selectEntries.apply')
            }}</n-button>
        </n-space>
      </template>
    </n-modal>

    <n-modal v-model:show="messageModalVisible" preset="dialog" :title="t('invoices.preparePdf.title')"
      :positive-text="t('invoices.preparePdf.positive')" :negative-text="t('invoices.preparePdf.negative')"
      :loading="pdfLoading" @positive-click="confirmDownload">
      <n-input v-model:value="messageDraft" type="textarea" :autosize="{ minRows: 4, maxRows: 8 }"
        :placeholder="t('invoices.preparePdf.messagePlaceholder')" />
    </n-modal>
  </PageContainer>
</template>

<style scoped>
.view-content {
  /* No flex layout to avoid overlap issues with grid */
  display: block;
}

.stats-grid-container {
  margin-bottom: 24px;
  padding: 2px 12px;
  /* 12px horizontal padding to compensate for n-row gutter of 24px (-12px margin) */
}



.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.06);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba(32, 128, 240, 0.1);
  color: var(--n-primary-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.stat-icon.secondary {
  background: rgba(100, 100, 100, 0.1);
  color: var(--n-text-color-3);
}



.currency {
  font-size: 14px;
  font-weight: 500;
  color: var(--n-text-color-3);
  margin-left: 8px;
}

.content-wrapper {
  background: var(--n-card-color);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.02);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  gap: 16px;
}

.search-input {
  max-width: 300px;
}

.filters {
  display: flex;
  gap: 8px;
  background: var(--n-color-modal);
  padding: 4px;
  border-radius: 8px;
}

/* Table Styling */
.invoice-table :deep(.n-data-table-th) {
  font-weight: 600;
  font-size: 13px;
  color: var(--n-text-color-3);
  background: #fafafc;
  border-bottom: 1px solid var(--n-divider-color);
  padding: 12px 16px;
}

.invoice-table :deep(.n-data-table-td) {
  padding: 16px;
  border-bottom: 1px solid var(--n-divider-color) !important;
}

.invoice-table :deep(.invoice-row) {
  transition: background-color 0.2s;
  cursor: default;
}

.invoice-table :deep(.invoice-row:hover) {
  background-color: rgba(0, 0, 0, 0.015);
}

.invoice-table :deep(.invoice-number-cell) {
  display: flex;
  align-items: center;
}

.invoice-table :deep(.invoice-number-cell:hover .n-text) {
  color: var(--n-primary-color);
  text-decoration: underline;
}

.create-btn {
  border-radius: 8px;
  font-weight: 600;
  height: 40px;
  padding: 0 20px;
  box-shadow: 0 4px 14px rgba(32, 128, 240, 0.3);
}

.action-btn:hover {
  background-color: var(--n-action-color);
  color: var(--n-primary-color);
}


.empty-state {
  border: 2px dashed #d9d9d9;
  /* Explicit darker gray for visibility */
  border-radius: 8px;
  padding: 32px;
  margin: 16px 0;
  background-color: #fafafa;
  /* Slight background contrast */
  display: flex;
  justify-content: center;
}

/* Stat card gradients/colors */
.stat-card {
  background: var(--n-card-color);
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.02);
  border: 1px solid rgba(0, 0, 0, 0.05);
  /* slightly visible border */
  transition: transform 0.2s, box-shadow 0.2s;
  /* height: 100%; removed to prevent overflow issues in simple row layout */
}

.stat-card.primary {
  background: linear-gradient(135deg, #f0f9ff 0%, #ffffff 100%);
  border-color: rgba(32, 128, 240, 0.2);
}

.stat-card.secondary {
  background: linear-gradient(135deg, #f9f9f9 0%, #ffffff 100%);
  border-color: rgba(0, 0, 0, 0.08);

}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.06);
}

.stat-card.primary .stat-icon {
  background: rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 12px rgba(32, 128, 240, 0.15);
}

.stat-card.secondary .stat-icon {
  background: rgba(255, 255, 255, 0.8);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}
</style>
