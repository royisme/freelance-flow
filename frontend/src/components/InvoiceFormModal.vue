<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import {
  NModal,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSelect,
  NDatePicker,
  NButton,
  NSpace,
  NText,
  useMessage,
  NGrid,
  NGi,
  NDivider,
  NDataTable,
  NEmpty
} from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { api } from '@/api'
import type { Invoice, InvoiceItem, Client, Project, TimeEntry, CreateInvoiceInput, UpdateInvoiceInput } from '@/types'
import type { FormInst, FormRules, DataTableColumns } from 'naive-ui'
import { invoiceSchema } from '@/schemas/invoice'
import { useZodRule } from '@/utils/validation'

interface Props {
  show: boolean
  invoice?: Invoice | null
  clients: Client[]
}

interface CreateInvoiceFromEntriesPayload {
  input: CreateInvoiceInput
  timeEntryIds: number[]
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'create', payload: CreateInvoiceFromEntriesPayload): void
  (e: 'update', invoice: UpdateInvoiceInput): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const message = useMessage()
const { t } = useI18n()

const formRef = ref<FormInst | null>(null)

type InvoiceStatus = 'draft' | 'sent' | 'paid' | 'overdue'

interface InvoiceFormData {
  clientId: number | null
  number: string
  issueDate: string
  dueDate: string
  items: InvoiceItem[]
  subtotal: number
  taxRate: number
  taxAmount: number
  total: number
  status: InvoiceStatus
}

function coerceInvoiceStatus(status: string): InvoiceStatus {
  if (status === 'draft' || status === 'sent' || status === 'paid' || status === 'overdue') return status
  return 'draft'
}

function defaultDueDateFromIssueDate(issueDate: string): string {
  const ts = Date.parse(issueDate)
  if (!Number.isFinite(ts)) return new Date().toISOString().split('T')[0]!
  const plusDays = 14
  return new Date(ts + plusDays * 24 * 60 * 60 * 1000).toISOString().split('T')[0]!
}

const isEditMode = computed(() => Boolean(props.invoice))

const formValue = ref<InvoiceFormData>({
  clientId: null,
  number: '',
  issueDate: new Date().toISOString().split('T')[0],
  dueDate: defaultDueDateFromIssueDate(new Date().toISOString().split('T')[0]!),
  items: [],
  subtotal: 0,
  taxRate: 0.13,
  taxAmount: 0,
  total: 0,
  status: 'draft'
})

const rules: FormRules = {
  clientId: {
    required: true,
    type: 'number',
    trigger: ['blur', 'change'],
    message: t('form.validation.select', { field: t('invoices.form.client') })
  },
  number: useZodRule(invoiceSchema.shape.number),
  issueDate: useZodRule(invoiceSchema.shape.issueDate),
  dueDate: useZodRule(invoiceSchema.shape.dueDate),
  taxRate: useZodRule(invoiceSchema.shape.taxRate),
  status: useZodRule(invoiceSchema.shape.status)
}

const statusOptions = [
  { label: t('invoices.status.draft'), value: 'draft' },
  { label: t('invoices.status.sent'), value: 'sent' },
  { label: t('invoices.status.paid'), value: 'paid' },
  { label: t('invoices.status.overdue'), value: 'overdue' }
]

const projects = ref<Project[]>([])
const projectsLoading = ref(false)
const selectedProjectId = ref<number | null>(null)

const timeEntries = ref<TimeEntry[]>([])
const timeEntriesLoading = ref(false)
const selectedTimeEntryIds = ref<number[]>([])

const eligibleTimeEntries = computed(() =>
  timeEntries.value.filter((e) => e.billable && !e.invoiceId && !e.invoiced)
)

const selectedHours = computed(() => {
  const selectedSet = new Set(selectedTimeEntryIds.value)
  const totalSeconds = eligibleTimeEntries.value
    .filter((e) => selectedSet.has(e.id))
    .reduce((sum, e) => sum + e.durationSeconds, 0)
  return totalSeconds / 3600
})

const timeEntryColumns = computed<DataTableColumns<TimeEntry>>(() => [
  { title: t('invoices.selectEntries.columns.date'), key: 'date', width: 110 },
  { title: t('timesheet.entries.description'), key: 'description', ellipsis: true },
  {
    title: t('invoices.selectEntries.columns.hours'),
    key: 'durationSeconds',
    width: 90,
    align: 'right',
    render: (row) => (row.durationSeconds / 3600).toFixed(2)
  }
])

watch(() => props.invoice, (newInvoice) => {
  if (newInvoice) {
    formValue.value = {
      clientId: newInvoice.clientId,
      number: newInvoice.number,
      issueDate: newInvoice.issueDate,
      dueDate: newInvoice.dueDate || defaultDueDateFromIssueDate(newInvoice.issueDate),
      items: newInvoice.items.map(i => ({ ...i })), // Deep copy items to avoid reactive issues
      subtotal: newInvoice.subtotal,
      taxRate: newInvoice.taxRate,
      taxAmount: newInvoice.taxAmount,
      total: newInvoice.total,
      status: coerceInvoiceStatus(newInvoice.status)
    }
    projects.value = []
    selectedProjectId.value = null
    timeEntries.value = []
    selectedTimeEntryIds.value = []
  } else {
    const issueDate = new Date().toISOString().split('T')[0]!
    formValue.value = {
      clientId: null,
      number: `INV-${Date.now().toString().slice(-6)}`,
      issueDate,
      dueDate: defaultDueDateFromIssueDate(issueDate),
      items: [],
      subtotal: 0,
      taxRate: 0.13,
      taxAmount: 0,
      total: 0,
      status: 'draft'
    }
    projects.value = []
    selectedProjectId.value = null
    timeEntries.value = []
    selectedTimeEntryIds.value = []
  }
}, { immediate: true })

watch(
  () => formValue.value.clientId,
  async (clientId) => {
    if (isEditMode.value) return
    selectedProjectId.value = null
    timeEntries.value = []
    selectedTimeEntryIds.value = []
    projects.value = []

    if (!clientId) return

    projectsLoading.value = true
    try {
      projects.value = await api.projects.listByClient(clientId)
    } catch {
      projects.value = []
      message.error(t('projects.loadError'))
    } finally {
      projectsLoading.value = false
    }
  }
)

watch(
  () => selectedProjectId.value,
  async (projectId) => {
    if (isEditMode.value) return
    timeEntries.value = []
    selectedTimeEntryIds.value = []
    if (!projectId) return

    timeEntriesLoading.value = true
    try {
      timeEntries.value = await api.timeEntries.list(projectId)
      selectedTimeEntryIds.value = eligibleTimeEntries.value.map((e) => e.id)
    } catch {
      timeEntries.value = []
      message.error(t('timesheet.loadError'))
    } finally {
      timeEntriesLoading.value = false
    }
  }
)

function handleClose() {
  emit('update:show', false)
}

function handleUpdateShow(value: boolean) {
  emit('update:show', value)
}

function handleSubmit() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (isEditMode.value && props.invoice) {
        emit('update', {
          id: props.invoice.id,
          clientId: formValue.value.clientId ?? 0,
          number: formValue.value.number,
          issueDate: formValue.value.issueDate,
          dueDate: formValue.value.dueDate,
          subtotal: formValue.value.subtotal,
          taxRate: formValue.value.taxRate,
          taxAmount: formValue.value.taxAmount,
          total: formValue.value.total,
          status: formValue.value.status,
          items: formValue.value.items.map((i) => ({
            description: i.description,
            quantity: i.quantity,
            unitPrice: i.unitPrice,
            amount: i.amount
          }))
        })
        handleClose()
        return
      }

      if (!formValue.value.clientId) {
        message.warning(t('form.validation.select', { field: t('invoices.form.client') }))
        return
      }
      if (!selectedProjectId.value) {
        message.warning(t('form.validation.select', { field: t('invoices.form.project') }))
        return
      }
      if (selectedTimeEntryIds.value.length === 0) {
        message.warning(t('invoices.form.validation.selectEntries'))
        return
      }

      emit('create', {
        input: {
          clientId: formValue.value.clientId,
          number: formValue.value.number,
          issueDate: formValue.value.issueDate,
          dueDate: formValue.value.dueDate,
          subtotal: 0,
          taxRate: formValue.value.taxRate,
          taxAmount: 0,
          total: 0,
          status: formValue.value.status,
          items: []
        },
        timeEntryIds: selectedTimeEntryIds.value
      })
      handleClose()
    } else {
      message.error(t('invoices.form.validation.fixErrors'))
    }
  })
}
</script>

<template>
  <n-modal :show="show" @update:show="handleUpdateShow" preset="card" :style="{ width: '900px', maxWidth: '95vw' }"
    :title="invoice ? t('invoices.form.editTitle') : t('invoices.form.newTitle')"
    :segmented="{ content: 'soft', footer: 'soft' }" size="huge">
    <n-form ref="formRef" :model="formValue" :rules="rules" label-placement="top" size="medium">

      <!-- Top Section: Client & Dates -->
      <n-grid :x-gap="24" :cols="2">
        <n-gi>
          <n-form-item :label="t('invoices.form.client')" path="clientId">
            <n-select v-model:value="formValue.clientId" :options="clients.map(c => ({ label: c.name, value: c.id }))"
              :placeholder="t('invoices.form.selectClient')" filterable />
          </n-form-item>
          <n-form-item v-if="!invoice" :label="t('invoices.form.project')" path="projectId">
            <n-select v-model:value="selectedProjectId"
              :options="projects.map(p => ({ label: p.name, value: p.id }))"
              :loading="projectsLoading"
              :disabled="!formValue.clientId"
              :placeholder="t('invoices.form.selectProject')"
              filterable />
          </n-form-item>
          <n-form-item :label="t('invoices.form.status')" path="status">
            <n-select v-model:value="formValue.status" :options="statusOptions" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <div class="invoice-meta-box">
            <n-grid :x-gap="12" :cols="2">
              <n-gi :span="2">
                <n-form-item :label="t('invoices.form.invoiceNumber')" path="number">
                  <n-input v-model:value="formValue.number" placeholder="INV-001" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item :label="t('invoices.form.issueDate')" path="issueDate">
                  <n-date-picker v-model:formatted-value="formValue.issueDate" type="date" value-format="yyyy-MM-dd"
                    style="width: 100%;" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item :label="t('invoices.form.dueDate')" path="dueDate">
                  <n-date-picker v-model:formatted-value="formValue.dueDate" type="date" value-format="yyyy-MM-dd"
                    style="width: 100%;" />
                </n-form-item>
              </n-gi>
            </n-grid>
          </div>
        </n-gi>
      </n-grid>

      <n-divider style="margin: 24px 0" />

      <!-- Create Mode: Pick time entries -->
      <template v-if="!invoice">
        <div class="items-header">
          <n-text strong>{{ t('invoices.form.timeEntries.title') }}</n-text>
          <n-text depth="3" style="font-size: 12px; margin-left: 8px;">
            {{ t('invoices.form.timeEntries.selectedHours', { hours: selectedHours.toFixed(2) }) }}
          </n-text>
        </div>

        <div class="items-container">
          <n-empty v-if="!selectedProjectId" :description="t('invoices.form.timeEntries.selectProjectHint')" size="small" />
          <n-empty v-else-if="!timeEntriesLoading && eligibleTimeEntries.length === 0"
            :description="t('invoices.form.timeEntries.empty')" size="small" />
          <n-data-table
            v-else
            :loading="timeEntriesLoading"
            :columns="timeEntryColumns"
            :data="eligibleTimeEntries"
            :row-key="(row: TimeEntry) => row.id"
            checkable
            :checked-row-keys="selectedTimeEntryIds"
            @update:checked-row-keys="(keys) => selectedTimeEntryIds = keys as number[]"
            :max-height="320"
          />
        </div>
      </template>

      <!-- Edit Mode: show derived items (read-only) -->
      <template v-else>
        <div class="items-header">
          <n-text strong>{{ t('invoices.form.items.title') }}</n-text>
          <n-text depth="3" style="font-size: 12px; margin-left: 8px;">
            ({{ formValue.items.length }} {{ t('invoices.form.items.countSuffix') }})
          </n-text>
        </div>
        <div class="items-container">
          <n-empty v-if="formValue.items.length === 0" :description="t('invoices.form.items.empty')" size="small" />
          <n-data-table
            v-else
            :columns="[
              { title: t('invoices.form.items.description'), key: 'description' },
              { title: t('invoices.form.items.qty'), key: 'quantity', width: 90, align: 'right', render: (row: InvoiceItem) => row.quantity.toFixed(2) },
              { title: t('invoices.form.items.rate'), key: 'unitPrice', width: 110, align: 'right', render: (row: InvoiceItem) => row.unitPrice.toFixed(2) },
              { title: t('invoices.form.items.amount'), key: 'amount', width: 110, align: 'right', render: (row: InvoiceItem) => row.amount.toFixed(2) }
            ]"
            :data="formValue.items"
            :row-key="(row: InvoiceItem) => row.id"
            :max-height="260"
          />
          <n-text depth="3" style="display:block; margin-top: 8px;">
            {{ t('invoices.form.items.editHint') }}
          </n-text>
        </div>
      </template>

      <!-- Totals Section -->
      <div class="totals-section">
        <div class="totals-grid">
          <div class="total-row">
            <span class="label">{{ t('invoices.form.subtotal') }}</span>
            <span class="value">${{ formValue.subtotal.toFixed(2) }}</span>
          </div>
          <div class="total-row">
            <span class="label">{{ t('invoices.form.taxRate') }}</span>
            <div class="value" style="width: 100px;">
              <n-input-number v-model:value="formValue.taxRate" :step="0.01" :min="0" :max="1" size="small"
                :show-button="false">
                <template #suffix>%</template>
              </n-input-number>
            </div>
          </div>
          <div class="total-row">
            <span class="label">{{ t('invoices.form.taxAmount') }}</span>
            <span class="value">${{ formValue.taxAmount.toFixed(2) }}</span>
          </div>
          <div class="total-row grand-total">
            <span class="label">{{ t('invoices.form.total') }}</span>
            <span class="value">${{ formValue.total.toFixed(2) }}</span>
          </div>
        </div>
      </div>

    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="handleClose" size="large">{{ t('invoices.form.cancel') }}</n-button>
        <n-button type="primary" @click="handleSubmit" size="large" style="padding-left: 32px; padding-right: 32px;">
          {{ invoice ? t('invoices.form.update') : t('invoices.form.create') }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>

<style scoped>
.invoice-meta-box {
  background-color: var(--n-color-modal);
  padding: 16px;
  border-radius: 8px;
  border: 1px solid var(--n-divider-color);
}

.items-header {
  margin-bottom: 12px;
}

.items-container {
  background-color: rgba(0, 0, 0, 0.02);
  border-radius: 8px;
  padding: 16px;
  border: 1px solid var(--n-divider-color);
}

.items-grid-header {
  display: flex;
  gap: 12px;
  padding: 0 46px 8px 12px;
  /* 46px right padding accounts for delete button area */
  font-size: 12px;
  font-weight: 600;
  color: var(--n-text-color-3);
  text-transform: uppercase;
}

.item-row {
  display: flex;
  gap: 12px;
  width: 100%;
  align-items: center;
}

.totals-section {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}

.totals-grid {
  width: 300px;
}

.total-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.total-row .label {
  color: var(--n-text-color-3);
}

.total-row .value {
  font-weight: 500;
  color: var(--n-text-color-1);
}

.total-row.grand-total {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--n-divider-color);
}

.total-row.grand-total .label {
  font-size: 16px;
  font-weight: 700;
  color: var(--n-text-color-1);
}

.total-row.grand-total .value {
  font-size: 20px;
  font-weight: 700;
  color: var(--n-primary-color);
}

/* Override default dynamic input styling to fit table look */
:deep(.n-dynamic-input .n-dynamic-input-item__action) {
  margin-left: 12px;
}
</style>
