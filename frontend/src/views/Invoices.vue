<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import InvoiceFormModal from '@/components/InvoiceFormModal.vue'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from '@/components/ui/dialog'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Textarea } from '@/components/ui/textarea'
import { useInvoiceStore, type EnrichedInvoice } from '@/stores/invoices'
import { useClientStore } from '@/stores/clients'
import { useTimesheetStore } from '@/stores/timesheet'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { Invoice, CreateInvoiceInput, UpdateInvoiceInput } from '@/types'
import { toast } from 'vue-sonner'
import {
  Plus,
  Download,
  FileText,
  DollarSign,
  Mail,
  Search,
  Check,
  Undo,
  Trash2,
  Loader2
} from 'lucide-vue-next'
import { api } from '@/api'

const invoiceStore = useInvoiceStore()
const clientStore = useClientStore()
const timesheetStore = useTimesheetStore()
const { enrichedInvoices, stats, loading } = storeToRefs(invoiceStore)
const { clients } = storeToRefs(clientStore)

const { t } = useI18n()

const showModal = ref(false)
const editingInvoice = ref<Invoice | null>(null)
const pdfLoading = ref(false)
const sendLoading = ref(false)
const messageModalVisible = ref(false)
const messageDraft = ref("")
const exportingInvoice = ref<EnrichedInvoice | null>(null)
const searchQuery = ref('')
const statusFilter = ref<string | null>(null)

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

const hasActiveFilters = computed(() => {
  return searchQuery.value !== '' || statusFilter.value !== null
})

function handleNewInvoice() {
  editingInvoice.value = null
  showModal.value = true
}

function handleEditInvoice(invoice: Invoice) {
  editingInvoice.value = invoice
  showModal.value = true
}

async function handleDownload(invoice: EnrichedInvoice) {
  try {
    pdfLoading.value = true
    exportingInvoice.value = invoice
    messageDraft.value = await invoiceStore.getDefaultMessage(invoice.id)
    messageModalVisible.value = true
  } catch {
    toast.error(t('invoices.downloadError'))
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
    // Use native save dialog and auto-open
    const { SaveAndOpenPDF } = await import('@/wailsjs/go/main/App')
    const savedPath = await SaveAndOpenPDF(base64, `INV-${exportingInvoice.value.number}.pdf`)
    if (savedPath) {
      toast.success(t('invoices.downloaded'))
    }
    // If savedPath is empty, user cancelled - no error
    messageModalVisible.value = false
  } catch (e) {
    console.error('PDF download error:', e)
    toast.error(t('invoices.downloadError'))
  } finally {
    pdfLoading.value = false
  }
}

async function handleSend(invoice: EnrichedInvoice) {
  try {
    sendLoading.value = true
    await invoiceStore.sendEmail(invoice.id)
    toast.success(t('invoices.sendSuccess'))
  } catch (e: any) {
    if (e && typeof e === 'string') {
      toast.error(e)
    } else if (e instanceof Error) {
      toast.error(e.message)
    } else {
      toast.error(t('invoices.sendError'))
    }
  } finally {
    sendLoading.value = false
  }
}

async function handleDelete(invoice: EnrichedInvoice) {
  try {
    await invoiceStore.deleteInvoice(invoice.id)
    toast.success(t('invoices.deleteSuccess'))
  } catch {
    toast.error(t('invoices.deleteError'))
  }
}

async function handleTogglePaymentStatus(invoice: EnrichedInvoice) {
  try {
    const newStatus = invoice.status === 'paid' ? 'sent' : 'paid'
    await invoiceStore.updateStatus(invoice.id, newStatus)
    toast.success(t('invoices.updateSuccess'))
  } catch {
    toast.error(t('invoices.saveError'))
  }
}

async function handleUpdateInvoice(payload: { input: UpdateInvoiceInput; timeEntryIds: number[] }) {
  try {
    await invoiceStore.updateInvoice(payload.input)
    await api.invoices.setTimeEntries({ invoiceId: payload.input.id, timeEntryIds: payload.timeEntryIds })
    await invoiceStore.fetchInvoices() // Refresh to get updated totals/status
    toast.success(t('invoices.updateSuccess'))
  } catch {
    toast.error(t('invoices.saveError'))
  }
}

async function handleCreateInvoiceFromEntries(payload: { input: CreateInvoiceInput; timeEntryIds: number[] }) {
  try {
    const created = await api.invoices.create(payload.input)
    await api.invoices.setTimeEntries({ invoiceId: created.id, timeEntryIds: payload.timeEntryIds })
    await invoiceStore.fetchInvoices()
    await timesheetStore.fetchTimesheet()
    toast.success(t('invoices.createSuccess'))
  } catch {
    toast.error(t('invoices.saveError'))
  }
}

onMounted(() => {
  invoiceStore.fetchInvoices()
  clientStore.fetchClients()
  timesheetStore.fetchTimesheet()
})

function getStatusColor(status: string) {
  if (status === 'paid') return 'text-green-600 bg-green-600/10'
  if (status === 'sent') return 'text-amber-600 bg-amber-600/10'
  if (status === 'overdue') return 'text-red-600 bg-red-600/10'
  return 'text-muted-foreground bg-muted'
}

function getStatusDotColor(status: string) {
  if (status === 'paid') return 'bg-green-600'
  if (status === 'sent') return 'bg-amber-600'
  if (status === 'overdue') return 'bg-red-600'
  return 'bg-muted-foreground'
}

function formatCurrency(value: number) {
  return value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}
</script>

<template>
  <div class="h-full flex flex-col min-h-0 gap-4">
    <!-- Header Actions -->
    <div class="shrink-0 flex items-center justify-end">
      <Button class="shadow-lg shadow-primary/30" @click="handleNewInvoice">
        <Plus class="w-4 h-4 mr-2" />
        {{ t('invoices.createInvoice') }}
      </Button>
    </div>

    <InvoiceFormModal v-model:show="showModal" :invoice="editingInvoice" :clients="clients"
      @create="handleCreateInvoiceFromEntries" @update="handleUpdateInvoice" />

    <!-- Stats Grid -->
    <div class="shrink-0 grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        class="stat-card primary flex items-center gap-5 p-5 rounded-2xl border bg-linear-to-br from-blue-50 to-white dark:from-blue-950/20 dark:to-card border-blue-200/50 dark:border-blue-900/50 shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-md">
        <div
          class="flex items-center justify-center size-12 rounded-xl bg-white dark:bg-blue-900/30 shadow-sm text-primary">
          <DollarSign class="size-6" />
        </div>
        <div>
          <p class="text-sm text-muted-foreground">{{ t('invoices.stats.outstandingAmount') }}</p>
          <p class="text-2xl font-bold tabular-nums">
            {{ formatCurrency(stats.totalDue) }}
            <span class="text-sm font-medium text-muted-foreground ml-1">USD</span>
          </p>
        </div>
      </div>

      <div
        class="stat-card secondary flex items-center gap-5 p-5 rounded-2xl border bg-linear-to-br from-gray-50 to-white dark:from-gray-900/20 dark:to-card border-border shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-md">
        <div
          class="flex items-center justify-center size-12 rounded-xl bg-white dark:bg-gray-800/50 shadow-sm text-muted-foreground">
          <FileText class="size-6" />
        </div>
        <div>
          <p class="text-sm text-muted-foreground">{{ t('invoices.stats.totalInvoices') }}</p>
          <p class="text-2xl font-bold tabular-nums">{{ enrichedInvoices.length }}</p>
        </div>
      </div>
    </div>

    <!-- Table Section -->
    <div class="flex-1 min-h-0 flex flex-col overflow-hidden">
      <!-- Search and Filter Bar -->
      <div class="shrink-0 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-4">
        <div class="relative max-w-sm w-full">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input v-model="searchQuery" :placeholder="t('invoices.searchPlaceholder')" class="pl-9" />
        </div>
        <div class="flex gap-2 bg-muted p-1 rounded-lg">
          <Button variant="ghost" size="sm" :class="{ 'bg-background shadow-sm': statusFilter === null }"
            @click="statusFilter = null">
            {{ t('invoices.filter.all') }}
          </Button>
          <Button variant="ghost" size="sm" :class="{ 'bg-background shadow-sm': statusFilter === 'draft' }"
            @click="statusFilter = 'draft'">
            {{ t('invoices.filter.draft') }}
          </Button>
          <Button variant="ghost" size="sm" :class="{ 'bg-background shadow-sm': statusFilter === 'sent' }"
            @click="statusFilter = 'sent'">
            {{ t('invoices.filter.sent') }}
          </Button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex-1 flex justify-center items-center">
        <Loader2 class="size-8 animate-spin text-muted-foreground" />
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredInvoices.length === 0"
        class="flex-1 flex flex-col items-center justify-center border-2 border-dashed rounded-lg bg-muted/30">
        <FileText class="size-12 text-muted-foreground/50 mb-4" />
        <p class="text-muted-foreground mb-4">
          {{ hasActiveFilters ? t('invoices.empty.noMatch') : t('invoices.empty.description') }}
        </p>
        <Button v-if="!hasActiveFilters" variant="outline" size="sm" @click="handleNewInvoice">
          {{ t('invoices.empty.action') }}
        </Button>
      </div>

      <!-- Table -->
      <div v-else class="flex-1 min-h-0 overflow-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead class="w-[160px]">{{ t('invoices.columns.invoiceNumber') }}</TableHead>
              <TableHead class="w-[120px]">{{ t('invoices.columns.client') }}</TableHead>
              <TableHead class="w-[130px]">{{ t('invoices.columns.issueDate') }}</TableHead>
              <TableHead class="w-[130px] text-right">{{ t('invoices.columns.amount') }}</TableHead>
              <TableHead class="w-[140px] text-center">{{ t('invoices.columns.status') }}</TableHead>
              <TableHead class="w-[140px] text-right">{{ t('invoices.columns.actions') }}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-for="row in filteredInvoices" :key="row.id"
              class="cursor-pointer hover:bg-muted/50 transition-colors">
              <TableCell>
                <div class="flex items-center gap-2 cursor-pointer hover:text-primary transition-colors"
                  @click="handleEditInvoice(row)">
                  <FileText class="size-4 text-muted-foreground" />
                  <span class="font-semibold">{{ row.number }}</span>
                </div>
              </TableCell>
              <TableCell class="font-medium">{{ row.clientName }}</TableCell>
              <TableCell class="text-muted-foreground">{{ row.issueDate }}</TableCell>
              <TableCell class="text-right font-semibold tabular-nums">
                {{ row.clientCurrency }} {{ formatCurrency(row.total) }}
              </TableCell>
              <TableCell class="text-center">
                <span
                  :class="['inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold', getStatusColor(row.status)]">
                  <span :class="['size-1.5 rounded-full', getStatusDotColor(row.status)]"></span>
                  {{ t(`invoices.status.${row.status}`).toUpperCase() }}
                </span>
              </TableCell>
              <TableCell class="text-right">
                <div class="flex gap-1 justify-end items-center">
                  <Button variant="ghost" size="icon" class="size-8"
                    :disabled="pdfLoading && exportingInvoice?.id === row.id" @click.stop="handleDownload(row)">
                    <Download class="size-4" />
                  </Button>

                  <AlertDialog>
                    <AlertDialogTrigger as-child>
                      <Button variant="ghost" size="icon" class="size-8 text-destructive hover:text-destructive"
                        @click.stop>
                        <Trash2 class="size-4" />
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>{{ t('common.confirmDelete') }}</AlertDialogTitle>
                        <AlertDialogDescription>
                          This action cannot be undone. This will permanently delete invoice {{ row.number }}.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction @click="handleDelete(row)">Delete</AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>

                  <Button variant="ghost" size="icon" class="size-8" :disabled="sendLoading"
                    @click.stop="handleSend(row)">
                    <Mail class="size-4" />
                  </Button>

                  <Button variant="ghost" size="icon" class="size-8"
                    :title="row.status === 'paid' ? t('invoices.actions.markAsUnpaid') : t('invoices.actions.markAsPaid')"
                    @click.stop="handleTogglePaymentStatus(row)">
                    <component :is="row.status === 'paid' ? Undo : Check" class="size-4" />
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </div>
    </div>

    <!-- PDF Message Dialog -->
    <Dialog v-model:open="messageModalVisible">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ t('invoices.preparePdf.title') }}</DialogTitle>
          <DialogDescription>Customize the message to include in the PDF invoice.</DialogDescription>
        </DialogHeader>
        <Textarea v-model="messageDraft" :placeholder="t('invoices.preparePdf.messagePlaceholder')"
          class="min-h-[120px]" />
        <DialogFooter>
          <Button variant="outline" @click="messageModalVisible = false">
            {{ t('invoices.preparePdf.negative') }}
          </Button>
          <Button :disabled="pdfLoading" @click="confirmDownload">
            <Loader2 v-if="pdfLoading" class="size-4 mr-2 animate-spin" />
            {{ t('invoices.preparePdf.positive') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
