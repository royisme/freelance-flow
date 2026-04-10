```
<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useI18n } from 'vue-i18n'
import { api } from '@/api'
import type { Invoice, Client, Project, TimeEntry } from '@/types'
import { dto } from '@/wailsjs/go/models'
import { invoiceSchema } from '@/schemas/invoice'
import { useUserTaxSettingsStore } from '@/stores/userTaxSettings'
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
import { SelectableDataTable } from '@/components/ui/data-table'
import { Separator } from '@/components/ui/separator'
import type { ColumnDef } from '@tanstack/vue-table'

interface Props {
  show?: boolean
  invoice?: Invoice | null
  clients: Client[]
}

interface CreateInvoiceFromEntriesPayload {
  input: dto.CreateInvoiceInput
  timeEntryIds: number[]
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'create', payload: CreateInvoiceFromEntriesPayload): void
  (e: 'update', payload: { input: dto.UpdateInvoiceInput; timeEntryIds: number[] }): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()
// const message = useMessage() // Removed: use toast or form errors instead, or mostly validation errors show inline
const { t } = useI18n()
const taxSettingsStore = useUserTaxSettingsStore()

// HST/Tax settings from user preferences
const taxEnabled = computed(() => taxSettingsStore.settings?.taxEnabled ?? false)
const defaultTaxRate = computed(() => {
  if (!taxEnabled.value) return 0
  return 0.13 // Default to 13% as per new tax logic (custom rate removed for now)
})

type InvoiceStatus = 'draft' | 'sent' | 'paid' | 'overdue'

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

// Form setup
const formSchema = toTypedSchema(invoiceSchema)
const form = useForm({
  validationSchema: formSchema,
})

// Access values for calculations
const formValues = form.values

const statusOptions = [
  { label: t('invoices.status.draft'), value: 'draft' },
  { label: t('invoices.status.sent'), value: 'sent' },
  { label: t('invoices.status.paid'), value: 'paid' },
  { label: t('invoices.status.overdue'), value: 'overdue' }
]

// Project selection
const selectedProjectId = ref<number | null>(null)
const selectedProjectRate = ref<number>(0)
const projectOptions = ref<Array<{ label: string; value: number }>>([])
const projectsLoading = ref(false)

// Time Entries
const timeEntries = ref<TimeEntry[]>([])
const timeEntriesLoading = ref(false)
const selectedTimeEntryIds = ref<number[]>([])

function setSelectionFromIds(ids: number[]) {
  selectedTimeEntryIds.value = ids
}

function clearSelection() {
  selectedTimeEntryIds.value = []
}

// Eligible time entries: billable entries that are either:
// 1. Already linked to this invoice (for editing)
// 2. Not linked to any invoice yet (available for selection)
const eligibleTimeEntries = computed(() =>
  timeEntries.value.filter((e) => {
    if (!e.billable) return false
    // Entry is linked to THIS invoice
    if (props.invoice && e.invoiceId === props.invoice.id) return true
    // Entry is not linked to any invoice (available)
    if (!e.invoiceId || e.invoiceId === 0) return true
    return false
  })
)

function getTimeEntryAmount(entry: TimeEntry): number {
  if (entry.billingMode === 'fixed') {
    return entry.manualAmount ?? 0
  }
  const project = projectsData.value.find((p) => p.id === entry.projectId)
  const rate = project?.hourlyRate ?? 0
  return (entry.durationSeconds / 3600) * rate
}

const selectedHours = computed(() => {
  const selectedSet = new Set(selectedTimeEntryIds.value)
  const totalSeconds = eligibleTimeEntries.value
    .filter((e) => selectedSet.has(e.id))
    .reduce((sum, e) => sum + e.durationSeconds, 0)
  return totalSeconds / 3600
})

const selectedSubtotal = computed(() => {
  const selectedSet = new Set(selectedTimeEntryIds.value)
  return eligibleTimeEntries.value
    .filter((e) => selectedSet.has(e.id))
    .reduce((sum, entry) => sum + getTimeEntryAmount(entry), 0)
})

// Time entries table columns
const timeEntryColumns: ColumnDef<TimeEntry>[] = [
  {
    accessorKey: 'date',
    header: () => t('invoices.selectEntries.columns.date'),
    cell: ({ row }) => row.original.date,
  },
  {
    accessorKey: 'description',
    header: () => t('timesheet.form.description'),
    cell: ({ row }) => row.original.description,
  },
  {
    accessorKey: 'durationSeconds',
    header: () => t('invoices.selectEntries.columns.hours'),
    cell: ({ row }) => row.original.billingMode === 'fixed' ? '-' : (row.original.durationSeconds / 3600).toFixed(2),
    meta: { class: 'text-right' },
  },
  {
    id: 'amount',
    header: () => t('invoices.columns.amount'),
    cell: ({ row }) => getTimeEntryAmount(row.original).toFixed(2),
    meta: { class: 'text-right' },
  },
]

// Helper date getter/setter factory
const createDateComputed = (fieldName: 'issueDate' | 'dueDate') => computed({
  get: (): DateValue | undefined => {
    const d = formValues[fieldName] as string | undefined
    if (!d) return undefined
    try {
      return parseDate(d)
    } catch {
      return undefined
    }
  },
  set: (val: DateValue | undefined) => {
    form.setFieldValue(fieldName, val ? val.toString() : undefined)
  }
})

const issueDateValue = createDateComputed('issueDate')
const dueDateValue = createDateComputed('dueDate')


// Save loaded projects
const projectsData = ref<Project[]>([])

async function loadProjects(clientId: number) {
  projectsLoading.value = true
  try {
    const projects = await api.projects.listByClient(clientId)
    projectsData.value = projects
    projectOptions.value = projects.map((p: Project) => ({
      label: p.name,
      value: p.id
    }))
  } catch {
    projectOptions.value = []
    projectsData.value = []
    // message.error(t('projects.loadError'))
    console.error('Failed to load projects')
  } finally {
    projectsLoading.value = false
  }
}

const isInitializing = ref(false)

// Watch Client ID change
watch(
  () => form.values.clientId,
  async (clientId) => {
    if (isInitializing.value) return

    // Reset project and time entries on client change
    selectedProjectId.value = null
    selectedProjectRate.value = 0
    projectOptions.value = []
    projectsData.value = []
    timeEntries.value = []
    clearSelection()

    if (!clientId) return

    // Ensure clientId is a number (may come as string from Select component)
    await loadProjects(Number(clientId))
  }
)

// Watch Project ID change
watch(
  () => selectedProjectId.value,
  async (projectId) => {
    // Only fetch if changed (or force fetch if needed, but usually watch triggers on change)
    // If we just loaded an invoice with project ID, this might trigger. 
    // We handle "isInitializing" to avoid clearing data during load, 
    // BUT checking logic below:

    // We WANT to fetch time entries whenever project changes, even during init?
    // During init we set selectedProjectId. If we wait for that tick, we fetch.

    // However, if we are in init, we might want to manually fetch inside the main watch(props.invoice)
    // to control the flow better?
    // Let's rely on this watcher but be careful not to reset selectedTimeEntryIds if we are initializing?
    // Actually, "Reset" happens at the top here:
    if (!isInitializing.value) {
      timeEntries.value = []
      clearSelection()
      selectedProjectRate.value = 0
    }

    if (!projectId) return

    const project = projectsData.value.find(p => p.id === projectId)
    if (project) {
      selectedProjectRate.value = project.hourlyRate
    }

    timeEntriesLoading.value = true
    try {
      timeEntries.value = await api.timeEntries.list(projectId)

      // Auto-select logic
      if (!isInitializing.value) { // Only auto-select if user actively changed project? 
        // Or if it's a new invoice creation flow. 
        // If editing, we rely on existing connections.
      }

      // Since existing selectedTimeEntryIds are reset above if !isInitializing,
      // here we re-populate.
      // If IS initializing, we might want to populate existing invoice entries.

      if (isEditMode.value && props.invoice) {
        // Editing mode: auto-select entries already linked to this invoice
        console.log('[InvoiceFormModal] timeEntries loaded:', timeEntries.value.length)
        console.log('[InvoiceFormModal] eligibleTimeEntries:', eligibleTimeEntries.value.length)

        // Auto-select entries that are linked to this invoice
        const linkedIds = eligibleTimeEntries.value
          .filter(e => e.invoiceId === props.invoice?.id)
          .map(e => e.id)

        console.log('[InvoiceFormModal] Auto-selecting linked entries:', linkedIds)
        setSelectionFromIds(linkedIds)
      } else {
        // New invoice: select all eligible entries by default
        setSelectionFromIds(eligibleTimeEntries.value.map((e) => e.id))
        console.log('[InvoiceFormModal] New invoice, selecting all:', selectedTimeEntryIds.value)
      }
    } catch {
      timeEntries.value = []
      // message.error
    } finally {
      timeEntriesLoading.value = false
    }
  }
)

// Watch Time Selection -> Update Totals
watch(
  () => selectedTimeEntryIds.value,
  () => {
    const subtotal = selectedSubtotal.value
    form.setFieldValue('subtotal', subtotal)

    const taxRate = form.values.taxRate || 0
    const taxAmount = subtotal * taxRate
    form.setFieldValue('taxAmount', taxAmount)
    form.setFieldValue('total', subtotal + taxAmount)
  },
  { deep: true }
)

// Also watch Project Rate change -> recalculate totals
watch(
  () => selectedProjectRate.value,
  () => {
    const subtotal = selectedSubtotal.value
    form.setFieldValue('subtotal', subtotal)

    const taxRate = form.values.taxRate || 0
    const taxAmount = subtotal * taxRate
    form.setFieldValue('taxAmount', taxAmount)
    form.setFieldValue('total', subtotal + taxAmount)
  }
)

// Also watch Tax Rate change
watch(
  () => form.values.taxRate,
  (newRate) => {
    const subtotal = form.values.subtotal || 0
    const taxAmount = subtotal * (newRate || 0)
    form.setFieldValue('taxAmount', taxAmount)
    form.setFieldValue('total', subtotal + taxAmount)
  }
)

// Initialize form based on invoice prop (runs once on component mount due to v-if)
watch(() => props.invoice, async (newInvoice) => {
  if (newInvoice) {
    isInitializing.value = true
    try {
      form.setValues({
        clientId: newInvoice.clientId,
        number: newInvoice.number,
        issueDate: newInvoice.issueDate,
        dueDate: newInvoice.dueDate || defaultDueDateFromIssueDate(newInvoice.issueDate),
        subtotal: newInvoice.subtotal,
        taxRate: newInvoice.taxRate,
        taxAmount: newInvoice.taxAmount,
        total: newInvoice.total,
        status: coerceInvoiceStatus(newInvoice.status),
        items: newInvoice.items || [],
      })

      // Load projects for the client
      if (newInvoice.clientId) {
        await loadProjects(Number(newInvoice.clientId))
      }

      // Set project ID and project rate if available
      const inv = newInvoice as any
      console.log('[InvoiceFormModal] Initializing invoice:', newInvoice)
      console.log('[InvoiceFormModal] Invoice projectId:', inv.projectId)
      console.log('[InvoiceFormModal] Loaded projects:', projectsData.value)

      if (inv.projectId) {
        selectedProjectId.value = inv.projectId
        // Also set the project rate from loaded projects
        const project = projectsData.value.find(p => p.id === inv.projectId)

        if (project) {
          selectedProjectRate.value = project.hourlyRate
        }
      } else {

        selectedProjectId.value = null
        selectedProjectRate.value = 0
      }
    } finally {
      setTimeout(() => {
        isInitializing.value = false
      }, 0)
    }
  } else {
    // New Invoice
    const issueDate = new Date().toISOString().split('T')[0]!
    const dueDate = defaultDueDateFromIssueDate(issueDate)

    // Generate Invoice Number
    const now = new Date()
    const yyyy = now.getFullYear()
    const mm = String(now.getMonth() + 1).padStart(2, '0')
    const dd = String(now.getDate()).padStart(2, '0')
    const hh = String(now.getHours()).padStart(2, '0')
    const min = String(now.getMinutes()).padStart(2, '0')
    const number = `INV-${yyyy}${mm}${dd}-${hh}${min}`

    form.resetForm({
      values: {
        clientId: undefined,
        number,
        issueDate,
        dueDate,
        subtotal: 0,
        taxRate: defaultTaxRate.value,
        taxAmount: 0,
        total: 0,
        status: 'draft',
        items: []
      }
    })

    selectedProjectId.value = null
    projectOptions.value = []
    timeEntries.value = []
    clearSelection()
  }
}, { immediate: true })

function handleClose() {
  handleUpdateShow(false)
}

function handleUpdateShow(value: boolean) {
  emit('update:show', value)
}

const onSubmit = form.handleSubmit((values) => {
  console.log('[InvoiceFormModal] Form submitted with values:', values)
  console.log('[InvoiceFormModal] selectedTimeEntryIds:', selectedTimeEntryIds.value)

  if (isEditMode.value && props.invoice) {
    console.log('[InvoiceFormModal] Emitting update event')
    emit('update', {
      input: dto.UpdateInvoiceInput.createFrom({
        id: props.invoice.id,
        clientId: values.clientId || 0,
        number: values.number,
        issueDate: values.issueDate,
        dueDate: values.dueDate,
        subtotal: values.subtotal,
        taxRate: values.taxRate,
        taxAmount: values.taxAmount,
        total: values.total,
        status: values.status,
        items: (values.items || []).map((i) => ({
          description: i.description,
          quantity: i.quantity,
          unitPrice: i.unitPrice,
          amount: i.amount
        }))
      }),
      timeEntryIds: selectedTimeEntryIds.value
    })
    handleUpdateShow(false)
    return
  }

  if (!values.clientId) {
    // Validation handled by form, but check custom checks
    return
  }
  if (!selectedProjectId.value) {
    // message.warning(t('form.validation.select', { field: t('invoices.form.project') }))
    return
  }
  if (selectedTimeEntryIds.value.length === 0) {
    // message.warning(t('invoices.form.validation.selectEntries'))
    return
  }


  emit('create', {
    input: dto.CreateInvoiceInput.createFrom({
      clientId: values.clientId,
      number: values.number,
      issueDate: values.issueDate,
      dueDate: values.dueDate,
      subtotal: values.subtotal || 0,
      taxRate: values.taxRate,
      taxAmount: values.taxAmount || 0,
      total: values.total || 0,
      status: values.status,
      items: []
    }),
    timeEntryIds: selectedTimeEntryIds.value
  })
  handleUpdateShow(false)
})
</script>

<template>
  <Dialog :open="show ?? true" @update:open="handleUpdateShow">
    <DialogContent class="sm:max-w-[900px] max-h-[85vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>{{ invoice ? t('invoices.form.editTitle') : t('invoices.form.newTitle') }}</DialogTitle>
      </DialogHeader>

      <form @submit="onSubmit" class="space-y-6">
        <!-- Top Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- Client -->
          <FormField v-slot="{ componentField }" name="clientId">
            <FormItem>
              <FormLabel>{{ t('invoices.form.client') }}</FormLabel>
              <Select v-bind="componentField" :model-value="componentField.modelValue?.toString()">
                <FormControl>
                  <SelectTrigger>
                    <SelectValue :placeholder="t('invoices.form.selectClient')" />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem v-for="client in clients" :key="client.id" :value="client.id.toString()">
                    {{ client.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Project (not a form field, just a local state selector) -->
          <div class="space-y-2">
            <label class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70">
              {{ t('invoices.form.project') }}
            </label>
            <Select :model-value="selectedProjectId?.toString()"
              @update:model-value="(v) => selectedProjectId = Number(v)"
              :disabled="!form.values.clientId || projectsLoading">
              <SelectTrigger>
                <SelectValue :placeholder="t('invoices.form.selectProject')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem v-for="p in projectsData" :key="p.id" :value="p.id.toString()">
                  {{ p.name }}
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          <!-- Status -->
          <FormField v-slot="{ componentField }" name="status">
            <FormItem>
              <FormLabel>{{ t('invoices.form.status') }}</FormLabel>
              <Select v-bind="componentField">
                <FormControl>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                </FormControl>
                <SelectContent>
                  <SelectItem v-for="opt in statusOptions" :key="opt.value" :value="opt.value">
                    {{ opt.label }}
                  </SelectItem>
                </SelectContent>
              </Select>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Invoice Number -->
          <FormField v-slot="{ componentField }" name="number">
            <FormItem>
              <FormLabel>{{ t('invoices.form.invoiceNumber') }}</FormLabel>
              <FormControl>
                <Input v-bind="componentField" :placeholder="t('invoices.form.invoiceNumberPlaceholder')" />
              </FormControl>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Issue Date -->
          <FormField name="issueDate">
            <FormItem class="flex flex-col">
              <FormLabel>{{ t('invoices.form.issueDate') }}</FormLabel>
              <Popover>
                <PopoverTrigger as-child>
                  <FormControl>
                    <Button variant="outline"
                      :class="cn('w-full pl-3 text-left font-normal', !issueDateValue && 'text-muted-foreground')">
                      <span>{{ issueDateValue ? issueDateValue.toString() : t('common.pickDate') }}</span>
                      <CalendarIcon class="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent class="w-auto p-0" align="start">
                  <Calendar v-model="issueDateValue" mode="single" />
                </PopoverContent>
              </Popover>
              <FormMessage />
            </FormItem>
          </FormField>

          <!-- Due Date -->
          <FormField name="dueDate">
            <FormItem class="flex flex-col">
              <FormLabel>{{ t('invoices.form.dueDate') }}</FormLabel>
              <Popover>
                <PopoverTrigger as-child>
                  <FormControl>
                    <Button variant="outline"
                      :class="cn('w-full pl-3 text-left font-normal', !dueDateValue && 'text-muted-foreground')">
                      <span>{{ dueDateValue ? dueDateValue.toString() : t('common.pickDate') }}</span>
                      <CalendarIcon class="ml-auto h-4 w-4 opacity-50" />
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent class="w-auto p-0" align="start">
                  <Calendar v-model="dueDateValue" mode="single" />
                </PopoverContent>
              </Popover>
              <FormMessage />
            </FormItem>
          </FormField>
        </div>

        <Separator />

        <!-- Time Entries Selection -->
        <div>
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <span class="font-semibold">{{ t('invoices.form.timeEntries.title') }}</span>
              <span class="text-xs text-muted-foreground">
                {{ t('invoices.form.timeEntries.selectedHours', { hours: selectedHours.toFixed(2) }) }}
              </span>
            </div>
          </div>

          <div v-if="!selectedProjectId" class="border rounded-md p-8 text-center text-muted-foreground">
            {{ t('invoices.form.timeEntries.selectProjectHint') }}
          </div>
          <SelectableDataTable
            v-else
            :columns="timeEntryColumns"
            :data="eligibleTimeEntries"
            :loading="timeEntriesLoading"
            :selectable="true"
            v-model:selection="selectedTimeEntryIds"
          />
        </div>

        <!-- Totals -->
        <div class="flex flex-col items-end space-y-2 pt-4 border-t">
          <div class="flex items-center gap-4 text-sm">
            <span class="text-muted-foreground">{{ t('invoices.form.subtotal') }}</span>
            <span class="font-medium min-w-[80px] text-right">${{ form.values.subtotal?.toFixed(2) || '0.00' }}</span>
          </div>

          <div v-if="taxEnabled" class="flex items-center gap-4 text-sm">
            <!-- Tax Rate Input -->
            <div class="flex items-center gap-2">
              <span class="text-muted-foreground">{{ t('invoices.form.taxRate') }}</span>
              <div class="relative w-[80px]">
                <Input type="number" step="0.01" min="0" max="1" class="h-8 pr-6 text-right"
                  :model-value="form.values.taxRate"
                  @update:model-value="(v) => form.setFieldValue('taxRate', Number(v))" />
                <span class="absolute right-2 top-0 bottom-0 flex items-center text-muted-foreground text-xs">%</span>
              </div>
            </div>
            <!-- Tax Amount -->
            <div class="flex items-center gap-4">
              <span class="text-muted-foreground">{{ t('invoices.form.taxAmount') }}</span>
              <span class="font-medium min-w-[80px] text-right">${{ form.values.taxAmount?.toFixed(2) || '0.00'
              }}</span>
            </div>
          </div>

          <div class="flex items-center gap-4 text-lg font-bold border-l-4 border-primary pl-4">
            <span>{{ t('invoices.form.total') }}</span>
            <span class="text-primary min-w-[80px] text-right">${{ form.values.total?.toFixed(2) || '0.00' }}</span>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" type="button" @click="handleClose">
            {{ t('invoices.form.cancel') }}
          </Button>
          <Button type="submit">
            {{ invoice ? t('invoices.form.update') : t('invoices.form.create') }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
