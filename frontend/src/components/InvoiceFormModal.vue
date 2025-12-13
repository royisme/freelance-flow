<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NDatePicker, NButton, NSpace, NDynamicInput, useMessage } from 'naive-ui'
import type { Invoice, InvoiceItem, Client } from '@/types'
import type { FormInst, FormRules } from 'naive-ui'

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
const formValue = ref<Omit<Invoice, 'id'>>({
  clientId: 0,
  number: '',
  issueDate: new Date().toISOString().split('T')[0],
  dueDate: '',
  items: [],
  subtotal: 0,
  taxRate: 0.13,
  taxAmount: 0,
  total: 0,
  status: 'draft'
})

import { invoiceSchema } from '@/schemas/invoice'
import { useZodRule } from '@/utils/validation'

const rules = {
  clientId: useZodRule(invoiceSchema.shape.clientId),
  number: useZodRule(invoiceSchema.shape.number),
  issueDate: useZodRule(invoiceSchema.shape.issueDate),
  dueDate: useZodRule(invoiceSchema.shape.dueDate),
  items: useZodRule(invoiceSchema.shape.items) // Validates the entire array
}

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Sent', value: 'sent' },
  { label: 'Paid', value: 'paid' },
  { label: 'Overdue', value: 'overdue' }
]

// Auto-calculate totals
const calculatedSubtotal = computed(() => {
  return formValue.value.items.reduce((sum, item) => sum + (item.amount || 0), 0)
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
      dueDate: newInvoice.dueDate,
      items: newInvoice.items,
      subtotal: newInvoice.subtotal,
      taxRate: newInvoice.taxRate,
      taxAmount: newInvoice.taxAmount,
      total: newInvoice.total,
      status: newInvoice.status
    }
  } else {
    formValue.value = {
      clientId: 0,
      number: `INV-${Date.now()}`,
      issueDate: new Date().toISOString().split('T')[0],
      dueDate: '',
      items: [],
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
  item.amount = item.quantity * item.unitPrice
}

function handleClose() {
  emit('update:show', false)
}

function handleSubmit() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (props.invoice) {
        emit('submit', { ...formValue.value, id: props.invoice.id } as Invoice)
      } else {
        emit('submit', formValue.value)
      }
      handleClose()
    } else {
      message.error('Please fix form errors')
    }
  })
}
</script>

<template>
  <n-modal :show="show" @update:show="handleClose" preset="card" :style="{ width: '700px' }"
    :title="invoice ? 'Edit Invoice' : 'New Invoice'">
    <n-form ref="formRef" :model="formValue" :rules="rules" label-placement="top"
      require-mark-placement="right-hanging">
      <n-space>
        <n-form-item label="Client" path="clientId" style="flex: 2;">
          <n-select v-model:value="formValue.clientId" :options="clients.map(c => ({ label: c.name, value: c.id }))"
            placeholder="Select client" />
        </n-form-item>

        <n-form-item label="Invoice Number" path="number" style="flex: 1;">
          <n-input v-model:value="formValue.number" placeholder="INV-001" />
        </n-form-item>
      </n-space>

      <n-space>
        <n-form-item label="Issue Date" path="issueDate" style="flex: 1;">
          <n-date-picker v-model:formatted-value="formValue.issueDate" type="date" value-format="yyyy-MM-dd"
            style="width: 100%;" />
        </n-form-item>

        <n-form-item label="Due Date" path="dueDate" style="flex: 1;">
          <n-date-picker v-model:formatted-value="formValue.dueDate" type="date" value-format="yyyy-MM-dd"
            style="width: 100%;" />
        </n-form-item>

        <n-form-item label="Status" path="status" style="flex: 1;">
          <n-select v-model:value="formValue.status" :options="statusOptions" />
        </n-form-item>
      </n-space>

      <n-form-item label="Line Items" path="items">
        <n-dynamic-input v-model:value="formValue.items" :on-create="createInvoiceItem" #="{ value }">
          <n-space>
            <n-input v-model:value="value.description" placeholder="Description" style="width: 200px;"
              @update:value="handleItemChange(value)" />
            <n-input-number v-model:value="value.quantity" :min="0" placeholder="Qty" style="width: 80px;"
              @update:value="handleItemChange(value)" />
            <n-input-number v-model:value="value.unitPrice" :min="0" placeholder="Rate" style="width: 100px;"
              @update:value="handleItemChange(value)" />
            <n-text strong>= ${{ value.amount.toFixed(2) }}</n-text>
          </n-space>
        </n-dynamic-input>
      </n-form-item>

      <n-space vertical style="width: 100%; margin-top: 16px;">
        <n-space justify="space-between">
          <n-text>Subtotal:</n-text>
          <n-text strong>${{ calculatedSubtotal.toFixed(2) }}</n-text>
        </n-space>
        <n-space justify="space-between">
          <n-text>Tax ({{ (formValue.taxRate * 100).toFixed(0) }}%):</n-text>
          <n-text strong>${{ calculatedTaxAmount.toFixed(2) }}</n-text>
        </n-space>
        <n-space justify="space-between">
          <n-text strong style="font-size: 1.1em;">Total:</n-text>
          <n-text strong style="font-size: 1.1em; color: var(--n-primary-color);">${{ calculatedTotal.toFixed(2)
          }}</n-text>
        </n-space>
      </n-space>
    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="handleClose">Cancel</n-button>
        <n-button type="primary" @click="handleSubmit">
          {{ invoice ? 'Update' : 'Create' }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>
