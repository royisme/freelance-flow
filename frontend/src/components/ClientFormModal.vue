<script setup lang="ts">
import { ref, watch } from 'vue'
import { NModal, NCard, NForm, NFormItem, NInput, NSelect, NButton, NSpace, useMessage } from 'naive-ui'
import type { Client } from '@/types'
import type { FormInst, FormRules } from 'naive-ui'
import { useI18n } from 'vue-i18n'

interface Props {
  show: boolean
  client?: Client | null
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'submit', client: Omit<Client, 'id'> | Client): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const message = useMessage()
const { t } = useI18n()

const formRef = ref<FormInst | null>(null)
const formValue = ref<Omit<Client, 'id'>>({
  name: '',
  email: '',
  website: '',
  avatar: '',
  contactPerson: '',
  address: '',
  currency: 'USD',
  status: 'active',
  notes: ''
})

import { clientSchema } from '@/schemas/client'
import { useZodRule } from '@/utils/validation'

const rules = {
  name: useZodRule(clientSchema.shape.name),
  email: useZodRule(clientSchema.shape.email),
  currency: useZodRule(clientSchema.shape.currency),
  status: useZodRule(clientSchema.shape.status)
}

const currencyOptions = [
  { label: 'USD - US Dollar', value: 'USD' },
  { label: 'CAD - Canadian Dollar', value: 'CAD' },
  { label: 'EUR - Euro', value: 'EUR' },
  { label: 'GBP - British Pound', value: 'GBP' }
]

const statusOptions = [
  { label: t('clients.status.active'), value: 'active' },
  { label: t('clients.status.inactive'), value: 'inactive' }
]

// Watch for client changes to populate form
watch(() => props.client, (newClient) => {
  if (newClient) {
    formValue.value = {
      name: newClient.name,
      email: newClient.email,
      website: newClient.website || '',
      avatar: newClient.avatar || '',
      contactPerson: newClient.contactPerson || '',
      address: newClient.address || '',
      currency: newClient.currency,
      status: newClient.status,
      notes: newClient.notes || ''
    }
  } else {
    // Reset for new client
    formValue.value = {
      name: '',
      email: '',
      website: '',
      avatar: '',
      contactPerson: '',
      address: '',
      currency: 'USD',
      status: 'active',
      notes: ''
    }
  }
}, { immediate: true })

function handleClose() {
  emit('update:show', false)
}

function handleSubmit() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (props.client) {
        // Update existing
        emit('submit', { ...formValue.value, id: props.client.id } as Client)
      } else {
        // Create new
        emit('submit', formValue.value)
      }
      handleClose()
    } else {
      message.error(t('form.saveError') || 'Please fix form errors')
    }
  })
}
</script>

<template>
  <n-modal :show="show" @update:show="handleClose" preset="card" :style="{ width: '600px' }"
    :title="client ? t('clients.editClient') : t('clients.newClient')">
    <n-form ref="formRef" :model="formValue" :rules="rules" label-placement="top"
      require-mark-placement="right-hanging">
      <n-form-item :label="t('form.client.name')" path="name">
        <n-input v-model:value="formValue.name" :placeholder="t('form.client.namePlaceholder')" />
      </n-form-item>

      <n-form-item :label="t('form.client.email')" path="email">
        <n-input v-model:value="formValue.email" :placeholder="t('form.client.emailPlaceholder')" />
      </n-form-item>

      <n-form-item :label="t('form.client.contactPerson')" path="contactPerson">
        <n-input v-model:value="formValue.contactPerson" :placeholder="t('form.client.contactPersonPlaceholder')" />
      </n-form-item>

      <n-form-item :label="t('form.client.website')" path="website">
        <n-input v-model:value="formValue.website" :placeholder="t('form.client.websitePlaceholder')" />
      </n-form-item>

      <n-form-item :label="t('form.client.address')" path="address">
        <n-input v-model:value="formValue.address" type="textarea" :placeholder="t('form.client.addressPlaceholder')"
          :rows="2" />
      </n-form-item>

      <n-space>
        <n-form-item :label="t('form.client.currency')" path="currency" style="flex: 1;">
          <n-select v-model:value="formValue.currency" :options="currencyOptions" />
        </n-form-item>

        <n-form-item :label="t('form.client.status')" path="status" style="flex: 1;">
          <n-select v-model:value="formValue.status" :options="statusOptions" />
        </n-form-item>
      </n-space>

      <n-form-item :label="t('form.client.notes')" path="notes">
        <n-input v-model:value="formValue.notes" type="textarea" :placeholder="t('form.client.notesPlaceholder')"
          :rows="3" />
      </n-form-item>
    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="handleClose">{{ t('form.cancel') }}</n-button>
        <n-button type="primary" @click="handleSubmit">
          {{ client ? t('form.update') : t('form.create') }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>
