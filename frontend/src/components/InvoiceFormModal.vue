<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import {
  NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NDatePicker,
  NButton, NSpace, NDynamicInput, NText, useMessage, NGrid, NGi, NDivider
} from 'naive-ui'
import type { Invoice, InvoiceItem, Client } from '@/types'
import type { FormInst } from 'naive-ui'
import { invoiceSchema } from '@/schemas/invoice'
import { useZodRule } from '@/utils/validation'

interface Props {
  show: boolean
  invoice?: Invoice | null
  clients: Client[]
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'submit', invoice: Omit<Invoice, 'id'> | Invoice): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const message = useMessage()

const formRef = ref<FormInst | null>(null)

type InvoiceStatus = 'draft' | 'sent' | 'paid' | 'overdue'

interface InvoiceFormData {
  clientId: number
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

const formValue = ref<InvoiceFormData>({
  clientId: 0,
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

const rules = {
  clientId: useZodRule(invoiceSchema.shape.clientId),
  number: useZodRule(invoiceSchema.shape.number),
  issueDate: useZodRule(invoiceSchema.shape.issueDate),
  dueDate: useZodRule(invoiceSchema.shape.dueDate),
  items: useZodRule(invoiceSchema.shape.items)
}

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Sent', value: 'sent' },
  { label: 'Paid', value: 'paid' },
  { label: 'Overdue', value: 'overdue' }
]

const calculatedSubtotal = computed(() => {
  return formValue.value.items.reduce((sum: number, item: InvoiceItem) => sum + (item.amount || 0), 0)
})

const calculatedTaxAmount = computed(() => {
  return calculatedSubtotal.value * formValue.value.taxRate
})

const calculatedTotal = computed(() => {
  return calculatedSubtotal.value + calculatedTaxAmount.value
})

watch([calculatedSubtotal, calculatedTaxAmount, calculatedTotal], () => {
  formValue.value.subtotal = calculatedSubtotal.value
  formValue.value.taxAmount = calculatedTaxAmount.value
  formValue.value.total = calculatedTotal.value
})

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
  } else {
    const issueDate = new Date().toISOString().split('T')[0]!
    formValue.value = {
      clientId: 0,
      number: `INV-${Date.now().toString().slice(-6)}`,
      issueDate,
      dueDate: defaultDueDateFromIssueDate(issueDate),
      items: [createInvoiceItem()], // Start with one empty item
      subtotal: 0,
      taxRate: 0.13,
      taxAmount: 0,
      total: 0,
      status: 'draft'
    }
  }
}, { immediate: true })

function createInvoiceItem(): InvoiceItem {
  return {
    id: Date.now(),
    description: '',
    quantity: 1,
    unitPrice: 0,
    amount: 0
  }
}

function handleItemChange(item: InvoiceItem) {
  item.amount = (item.quantity || 0) * (item.unitPrice || 0)
}

function handleClose() {
  emit('update:show', false)
}

function handleUpdateShow(value: boolean) {
  emit('update:show', value)
}

function handleSubmit() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (props.invoice) {
        emit('submit', { ...formValue.value, id: props.invoice.id } as Invoice)
      } else {
        emit('submit', formValue.value as unknown as Omit<Invoice, 'id'>)
      }
      handleClose()
    } else {
      message.error('Please fix form errors')
    }
  })
}
</script>

<template>
  <n-modal :show="show" @update:show="handleUpdateShow" preset="card" :style="{ width: '900px', maxWidth: '95vw' }"
    :title="invoice ? 'Edit Invoice' : 'New Invoice'" :segmented="{ content: 'soft', footer: 'soft' }" size="huge">
    <n-form ref="formRef" :model="formValue" :rules="rules" label-placement="top" size="medium">

      <!-- Top Section: Client & Dates -->
      <n-grid :x-gap="24" :cols="2">
        <n-gi>
          <n-form-item label="Client" path="clientId">
            <n-select v-model:value="formValue.clientId" :options="clients.map(c => ({ label: c.name, value: c.id }))"
              placeholder="Select client" filterable />
          </n-form-item>
          <n-form-item label="Status" path="status">
            <n-select v-model:value="formValue.status" :options="statusOptions" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <div class="invoice-meta-box">
            <n-grid :x-gap="12" :cols="2">
              <n-gi :span="2">
                <n-form-item label="Invoice Number" path="number">
                  <n-input v-model:value="formValue.number" placeholder="INV-001" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="Issue Date" path="issueDate">
                  <n-date-picker v-model:formatted-value="formValue.issueDate" type="date" value-format="yyyy-MM-dd"
                    style="width: 100%;" />
                </n-form-item>
              </n-gi>
              <n-gi>
                <n-form-item label="Due Date" path="dueDate">
                  <n-date-picker v-model:formatted-value="formValue.dueDate" type="date" value-format="yyyy-MM-dd"
                    style="width: 100%;" />
                </n-form-item>
              </n-gi>
            </n-grid>
          </div>
        </n-gi>
      </n-grid>

      <n-divider style="margin: 24px 0" />

      <!-- Items Section -->
      <div class="items-header">
        <n-text strong>Line Items</n-text>
        <n-text depth="3" style="font-size: 12px; margin-left: 8px;">
          ({{ formValue.items.length }} items)
        </n-text>
      </div>

      <div class="items-container">
        <!-- Custom Header Request: A simple grid header -->
        <div class="items-grid-header">
          <div style="flex: 4;">Description</div>
          <div style="flex: 1;">Qty</div>
          <div style="flex: 1.5;">Rate</div>
          <div style="flex: 1.5; text-align: right;">Amount</div>
          <div style="width: 32px;"></div>
        </div>

        <n-dynamic-input v-model:value="formValue.items" :on-create="createInvoiceItem" #="{ value }"
          class="custom-dynamic-input">
          <div class="item-row">
            <div style="flex: 4;">
              <n-input v-model:value="value.description" placeholder="Item description"
                @update:value="handleItemChange(value)" />
            </div>
            <div style="flex: 1;">
              <n-input-number v-model:value="value.quantity" :min="0" placeholder="0" :show-button="false"
                @update:value="handleItemChange(value)" />
            </div>
            <div style="flex: 1.5;">
              <n-input-number v-model:value="value.unitPrice" :min="0" placeholder="0.00" :precision="2"
                :show-button="false" @update:value="handleItemChange(value)">
                <template #prefix>$</template>
              </n-input-number>
            </div>
            <div style="flex: 1.5; text-align: right; display: flex; align-items: center; justify-content: flex-end;">
              <n-text strong>${{ value.amount.toFixed(2) }}</n-text>
            </div>
          </div>
        </n-dynamic-input>
      </div>

      <!-- Totals Section -->
      <div class="totals-section">
        <div class="totals-grid">
          <div class="total-row">
            <span class="label">Subtotal</span>
            <span class="value">${{ calculatedSubtotal.toFixed(2) }}</span>
          </div>
          <div class="total-row">
            <span class="label">Tax Rate</span>
            <div class="value" style="width: 100px;">
              <n-input-number v-model:value="formValue.taxRate" :step="0.01" :min="0" :max="1" size="small"
                :show-button="false">
                <template #suffix>%</template>
              </n-input-number>
            </div>
          </div>
          <div class="total-row">
            <span class="label">Tax Amount</span>
            <span class="value">${{ calculatedTaxAmount.toFixed(2) }}</span>
          </div>
          <div class="total-row grand-total">
            <span class="label">Total</span>
            <span class="value">${{ calculatedTotal.toFixed(2) }}</span>
          </div>
        </div>
      </div>

    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="handleClose" size="large">Cancel</n-button>
        <n-button type="primary" @click="handleSubmit" size="large" style="padding-left: 32px; padding-right: 32px;">
          {{ invoice ? 'Update Invoice' : 'Create Invoice' }}
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
