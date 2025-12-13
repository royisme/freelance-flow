<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NDatePicker, NTimePicker, NCheckbox, NButton, NSpace, useMessage } from 'naive-ui'
import type { TimeEntry, Project } from '@/types'
import type { FormInst, FormRules } from 'naive-ui'
import { useI18n } from 'vue-i18n'

interface Props {
  show: boolean
  entry?: TimeEntry | null
  projects: Project[]
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'submit', entry: Omit<TimeEntry, 'id'> | TimeEntry): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const message = useMessage()
const { t } = useI18n()

const formRef = ref<FormInst | null>(null)
const formValue = ref<Omit<TimeEntry, 'id'>>({
  projectId: 0,
  date: new Date().toISOString().split('T')[0] ?? '',
  startTime: '',
  endTime: '',
  durationSeconds: 0,
  description: '',
  invoiced: false
})

import { timeEntrySchema } from '@/schemas/timesheet'
import { useZodRule } from '@/utils/validation'

const rules = {
  projectId: useZodRule(timeEntrySchema.shape.projectId),
  date: useZodRule(timeEntrySchema.shape.date),
  durationSeconds: useZodRule(timeEntrySchema.shape.durationSeconds),
  description: useZodRule(timeEntrySchema.shape.description)
}

// Convert seconds to hours for input
const durationHours = computed({
  get: () => formValue.value.durationSeconds / 3600,
  set: (val: number) => { formValue.value.durationSeconds = Math.round(val * 3600) }
})

watch(() => props.entry, (newEntry) => {
  if (newEntry) {
    formValue.value = {
      projectId: newEntry.projectId,
      date: newEntry.date,
      startTime: newEntry.startTime || '',
      endTime: newEntry.endTime || '',
      durationSeconds: newEntry.durationSeconds,
      description: newEntry.description,
      invoiced: newEntry.invoiced
    }
  } else {
    formValue.value = {
      projectId: 0,
      date: new Date().toISOString().split('T')[0] ?? '',
      startTime: '',
      endTime: '',
      durationSeconds: 0,
      description: '',
      invoiced: false
    }
  }
}, { immediate: true })

function handleClose() {
  emit('update:show', false)
}

function handleSubmit() {
  formRef.value?.validate((errors) => {
    if (!errors) {
      if (props.entry) {
        emit('submit', { ...formValue.value, id: props.entry.id } as TimeEntry)
      } else {
        emit('submit', formValue.value)
      }
      handleClose()
    } else {
      message.error(t('form.saveError'))
    }
  })
}
</script>

<template>
  <n-modal :show="show" @update:show="handleClose" preset="card" :style="{ width: '600px' }"
    :title="entry ? t('timesheet.form.editTitle') : t('timesheet.form.createTitle')">
    <n-form ref="formRef" :model="formValue" :rules="rules" label-placement="top"
      require-mark-placement="right-hanging">
      <n-form-item :label="t('timesheet.form.project')" path="projectId">
        <n-select v-model:value="formValue.projectId" :options="projects.map(p => ({ label: p.name, value: p.id }))"
          :placeholder="t('timesheet.timer.selectProject')" />
      </n-form-item>

      <n-form-item :label="t('timesheet.form.date')" path="date">
        <n-date-picker v-model:formatted-value="formValue.date" type="date" value-format="yyyy-MM-dd"
          style="width: 100%;" />
      </n-form-item>

      <n-form-item :label="t('timesheet.form.duration')" path="durationSeconds">
        <n-input-number v-model:value="durationHours" :min="0" :step="0.25" placeholder="0.00" style="width: 100%;" />
      </n-form-item>

      <n-form-item :label="t('timesheet.form.description')" path="description">
        <n-input v-model:value="formValue.description" type="textarea"
          :placeholder="t('timesheet.form.descriptionPlaceholder')" :rows="3" />
      </n-form-item>

      <n-form-item :label="t('timesheet.form.timeRange')">
        <n-space>
          <n-time-picker v-model:formatted-value="formValue.startTime" value-format="HH:mm" format="HH:mm"
            :placeholder="t('timesheet.form.startTime')" />
          <span>to</span>
          <n-time-picker v-model:formatted-value="formValue.endTime" value-format="HH:mm" format="HH:mm"
            :placeholder="t('timesheet.form.endTime')" />
        </n-space>
      </n-form-item>

      <n-form-item>
        <n-checkbox v-model:checked="formValue.invoiced">{{ t('timesheet.form.alreadyInvoiced') }}</n-checkbox>
      </n-form-item>
    </n-form>

    <template #footer>
      <n-space justify="end">
        <n-button @click="handleClose">{{ t('common.cancel') }}</n-button>
        <n-button type="primary" @click="handleSubmit">
          {{ entry ? t('timesheet.form.update') : t('timesheet.logTime') }}
        </n-button>
      </n-space>
    </template>
  </n-modal>
</template>
