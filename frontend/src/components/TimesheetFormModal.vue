<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useI18n } from 'vue-i18n'
import type { TimeEntry, Project } from '@/types'
import { timeEntrySchema } from '@/schemas/timesheet'
import { Calendar as CalendarIcon } from 'lucide-vue-next'
import { cn } from '@/lib/utils'
import {
  parseDate,
  type DateValue,
} from '@internationalized/date'

import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Calendar } from '@/components/ui/calendar'
import { Checkbox } from '@/components/ui/checkbox'

interface Props {
  show?: boolean
  entry?: TimeEntry | null
  projects: Project[]
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'submit', entry: Omit<TimeEntry, 'id'> | TimeEntry): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
const { t } = useI18n()

// Extend schema if needed, or use as is
const formSchema = toTypedSchema(timeEntrySchema)

const form = useForm({
  validationSchema: formSchema,
})

// Access values
const formValues = form.values
const isFixedBilling = computed(() => formValues.billingMode === 'fixed')

// Duration Hours Computed Helper
const durationHours = computed({
  get: () => {
    const secs = formValues.durationSeconds
    return typeof secs === 'number' ? secs / 3600 : 0
  },
  set: (val: string | number) => {
    const num = Number(val)
    if (!isNaN(num)) {
      form.setFieldValue('durationSeconds', Math.round(num * 3600))
    }
  }
})

// Helper date getter/setter
const dateValue = computed({
  get: (): DateValue | undefined => {
    const d = formValues.date as string | undefined
    if (!d) return undefined
    try {
      return parseDate(d)
    } catch {
      return undefined
    }
  },
  set: (val: DateValue | undefined) => {
    form.setFieldValue('date', val ? val.toString() : undefined)
  }
})

// Billable state (not in schema but needed for form)
const billable = ref(true)

// Initialize form based on entry prop (runs once on component mount due to v-if)
watch(() => props.entry, (newEntry) => {
  if (newEntry) {
    form.setValues({
      projectId: newEntry.projectId,
      date: newEntry.date,
      startTime: newEntry.startTime || undefined,
      endTime: newEntry.endTime || undefined,
      durationSeconds: newEntry.durationSeconds,
      description: newEntry.description,
      invoiced: newEntry.invoiced,
      billingMode: newEntry.billingMode || 'hourly',
      manualAmount: newEntry.manualAmount ?? null,
    })
    billable.value = !!newEntry.billable
  } else {
    resetForm()
  }
}, { immediate: true })

function resetForm() {
  form.resetForm({
    values: {
      projectId: undefined,
      date: new Date().toISOString().split('T')[0],
      startTime: undefined,
      endTime: undefined,
      durationSeconds: 0,
      description: '',
      invoiced: false,
      billingMode: 'hourly',
      manualAmount: null,
    }
  })
  billable.value = true
}

function handleUpdateShow(value: boolean) {
  emit('update:show', value)
}

function handleClose() {
  handleUpdateShow(false)
}

const onSubmit = form.handleSubmit((values) => {
  // Construct payload
  const payload: any = {
    projectId: values.projectId,
    date: values.date,
    startTime: values.startTime || '',
    endTime: values.endTime || '',
    durationSeconds: values.durationSeconds,
    description: values.description,
    invoiced: values.invoiced,
    billingMode: values.billingMode,
    manualAmount: values.billingMode === 'fixed' ? values.manualAmount ?? 0 : null,
    billable: billable.value,
    invoiceId: props.entry?.invoiceId || 0
  }

  if (props.entry) {
    emit('submit', { ...payload, id: props.entry.id } as TimeEntry)
  } else {
    emit('submit', payload as Omit<TimeEntry, 'id'>)
  }
  handleClose()
})
</script>

<template>
  <Dialog :open="show ?? true" @update:open="handleUpdateShow">
    <DialogContent class="sm:max-w-[600px]">
      <DialogHeader>
        <DialogTitle>{{ entry ? t('timesheet.form.editTitle') : t('timesheet.form.createTitle') }}</DialogTitle>
      </DialogHeader>

      <form @submit="onSubmit" class="space-y-6">
        <!-- Project -->
        <FormField v-slot="{ componentField }" name="projectId">
          <FormItem>
            <FormLabel>{{ t('timesheet.form.project') }}</FormLabel>
            <Select v-bind="componentField" :model-value="componentField.modelValue?.toString()">
              <FormControl>
                <SelectTrigger>
                  <SelectValue :placeholder="t('timesheet.timer.selectProject')" />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectItem v-for="p in projects" :key="p.id" :value="p.id.toString()">
                  {{ p.name }}
                </SelectItem>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        </FormField>

        <!-- Date -->
        <FormField name="date">
          <FormItem class="flex flex-col">
            <FormLabel>{{ t('timesheet.form.date') }}</FormLabel>
            <Popover>
              <PopoverTrigger as-child>
                <FormControl>
                  <Button variant="outline"
                    :class="cn('w-full pl-3 text-left font-normal', !dateValue && 'text-muted-foreground')">
                    <span>{{ dateValue ? dateValue.toString() : t('common.pickDate') }}</span>
                    <CalendarIcon class="ml-auto h-4 w-4 opacity-50" />
                  </Button>
                </FormControl>
              </PopoverTrigger>
              <PopoverContent class="w-auto p-0" align="start">
                <Calendar v-model="dateValue" mode="single" />
              </PopoverContent>
            </Popover>
            <FormMessage />
          </FormItem>
        </FormField>

        <!-- Duration -->
        <FormField name="durationSeconds">
          <FormItem>
            <FormLabel>{{ t('timesheet.form.duration') }}</FormLabel>
            <FormControl>
              <Input type="number" step="0.25" min="0" v-model="durationHours" placeholder="0.00" :disabled="isFixedBilling" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-slot="{ componentField }" name="billingMode">
          <FormItem>
            <FormLabel>{{ t('timesheet.form.billingMode') }}</FormLabel>
            <Select v-bind="componentField" :model-value="componentField.modelValue">
              <FormControl>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
              </FormControl>
              <SelectContent>
                <SelectItem value="hourly">{{ t('timesheet.form.billingModes.hourly') }}</SelectItem>
                <SelectItem value="fixed">{{ t('timesheet.form.billingModes.fixed') }}</SelectItem>
              </SelectContent>
            </Select>
            <FormMessage />
          </FormItem>
        </FormField>

        <FormField v-if="isFixedBilling" v-slot="{ componentField }" name="manualAmount">
          <FormItem>
            <FormLabel>{{ t('timesheet.form.manualAmount') }}</FormLabel>
            <FormControl>
              <Input type="number" step="0.01" min="0" v-bind="componentField" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <!-- Description -->
        <FormField v-slot="{ componentField }" name="description">
          <FormItem>
            <FormLabel>{{ t('timesheet.form.description') }}</FormLabel>
            <FormControl>
              <Textarea v-bind="componentField" :rows="3" :placeholder="t('timesheet.form.descriptionPlaceholder')" />
            </FormControl>
            <FormMessage />
          </FormItem>
        </FormField>

        <!-- Time Range -->
        <div class="grid grid-cols-2 gap-4">
          <FormField v-slot="{ componentField }" name="startTime">
            <FormItem>
              <FormLabel>{{ t('timesheet.form.startTime') }}</FormLabel>
              <FormControl>
                <Input type="time" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <FormField v-slot="{ componentField }" name="endTime">
            <FormItem>
              <FormLabel>{{ t('timesheet.form.endTime') }}</FormLabel>
              <FormControl>
                <Input type="time" v-bind="componentField" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>

        <!-- Checkboxes -->
        <div class="flex items-center gap-6">
          <div class="flex items-center gap-2">
            <Checkbox id="billable" :checked="billable" @update:checked="(v: boolean) => billable = !!v" />
            <label for="billable"
              class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              {{ t('timesheet.entries.billable') }}
            </label>
          </div>

          <FormField v-slot="{ value, handleChange }" name="invoiced">
            <FormItem class="flex flex-row items-center gap-2 space-y-0">
              <FormControl>
                <Checkbox :checked="value" @update:checked="handleChange" />
              </FormControl>
              <FormLabel class="font-normal">
                {{ t('timesheet.form.alreadyInvoiced') }}
              </FormLabel>
            </FormItem>
          </FormField>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" @click="handleClose">
            {{ t('common.cancel') }}
          </Button>
          <Button type="submit">
            {{ entry ? t('timesheet.form.update') : t('timesheet.logTime') }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
