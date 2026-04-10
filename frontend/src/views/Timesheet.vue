<script setup lang="ts">
import { h, onMounted, ref, computed } from 'vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

import { DataTable, DataTableColumnHeader } from '@/components/ui/data-table'
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

import QuickTimeEntry from '@/components/QuickTimeEntry.vue'
import TimesheetFormModal from '@/components/TimesheetFormModal.vue'
import { useTimesheetStore } from '@/stores/timesheet'
import { useProjectStore } from '@/stores/projects'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { TimeEntry } from '@/types'
import { dateOnlySortKey, formatISODateOnlyForLocale } from '@/utils/date'
import { toast } from 'vue-sonner'
import {
  Clock,
  Edit,
  Trash2,
  Download,
  Loader2
} from 'lucide-vue-next'
import type { ColumnDef } from '@tanstack/vue-table'

const timesheetStore = useTimesheetStore()
const projectStore = useProjectStore()
const { entries, enrichedEntries, loading } = storeToRefs(timesheetStore)
const { projects } = storeToRefs(projectStore)
const { t, locale } = useI18n()

const showModal = ref(false)
const editingEntry = ref<TimeEntry | null>(null)

function handleCreateEntry() {
  editingEntry.value = null
  showModal.value = true
}

// Helpers
function formatDate(dateStr: string): string {
  return formatISODateOnlyForLocale(dateStr, locale.value === 'zh-CN' ? 'zh-CN' : 'en-US')
}

function formatHours(seconds: number): string {
  const hours = seconds / 3600
  return hours.toFixed(hours % 1 === 0 ? 0 : 1)
}

function getProjectColor(projectId: number): string {
  const colors = ['#10b981', '#3b82f6', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899']
  const project = projects.value.find(p => p.id === projectId)
  if (!project) return colors[0]!
  return colors[projectId % colors.length]!
}

function getProjectRate(projectId: number): number {
  const project = projects.value.find(p => p.id === projectId)
  return project?.hourlyRate || 0
}

function getEntryAmount(entry: TimeEntry): number {
  if (entry.billingMode === 'fixed') {
    return entry.manualAmount ?? 0
  }
  const rate = getProjectRate(entry.projectId)
  const hours = entry.durationSeconds / 3600
  return rate * hours
}

// Entry actions
function handleEdit(entry: TimeEntry) {
  editingEntry.value = entry
  showModal.value = true
}

async function handleDelete(id: number) {
  try {
    await timesheetStore.deleteTimeEntry(id)
    toast.success(t('timesheet.entry.deletedMsg'))
  } catch {
    toast.error(t('timesheet.messages.deleteError'))
  }
}

// Form handlers
async function handleSubmitEntry(entry: Omit<TimeEntry, 'id'> | TimeEntry) {
  try {
    // Check if id exists AND is a valid number (not undefined/null/0)
    const isUpdate = 'id' in entry && typeof entry.id === 'number' && entry.id > 0
    if (isUpdate) {
      await timesheetStore.updateTimeEntry(entry as TimeEntry)
      toast.success(t('timesheet.entry.updatedMsg'))
    } else {
      // Remove id property if it exists but is invalid
      const { id, ...createData } = entry as TimeEntry
      await timesheetStore.createTimeEntry(createData)
      toast.success(t('timesheet.messages.logged'))
    }
  } catch {
    toast.error(t('timesheet.messages.saveError'))
  }
}

async function handleQuickEntry(data: { projectId: number; description: string; durationSeconds: number; date: string; billable: boolean }) {
  try {
    await timesheetStore.createTimeEntry({
      projectId: data.projectId,
      description: data.description,
      durationSeconds: data.durationSeconds,
      date: data.date,
      startTime: '',
      endTime: '',
      billable: data.billable,
      invoiced: false,
      billingMode: 'hourly',
      manualAmount: null,
    })
  } catch {
    toast.error(t('timesheet.messages.saveError'))
  }
}

function handleExportCSV() {
  try {
    timesheetStore.exportToCSV(timesheetStore.enrichedEntries)
    toast.success(t('timesheet.messages.exportCsvSuccess'))
  } catch {
    toast.error(t('timesheet.messages.exportCsvError'))
  }
}

const columns = computed<ColumnDef<TimeEntry & { project?: { name: string } }>[]>(() => [
  {
    accessorKey: 'date',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('timesheet.columns.date') as string }),
    cell: ({ row }) => formatDate(row.getValue('date')),
    sortingFn: (a, b) => dateOnlySortKey(a.original.date) - dateOnlySortKey(b.original.date),
  },
  {
    accessorKey: 'project',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('timesheet.columns.project') as string }),
    cell: ({ row }) => {
      const project = row.original.project
      const projectName = project?.name || t('timesheet.entry.noProject')
      const projectColor = getProjectColor(row.original.projectId)
      return h('div', { class: 'project-cell flex items-center gap-2' }, [
        h('span', { class: 'project-dot w-2 h-2 rounded-full', style: { backgroundColor: projectColor } }),
        h('span', {}, projectName)
      ])
    }
  },
  {
    accessorKey: 'description',
    header: t('timesheet.columns.task'),
    cell: ({ row }) => row.original.description || '-'
  },
  {
    accessorKey: 'billable',
    header: t('timesheet.columns.status'),
    cell: ({ row }) => {
      const billable = row.original.billable
      return h(Badge, {
        variant: billable ? 'default' : 'secondary',
        class: 'rounded-full px-2'
      }, () => billable ? t('timesheet.entries.billable') : t('timesheet.entries.nonBillable'))
    }
  },
  {
    accessorKey: 'durationSeconds',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('timesheet.columns.hours') as string }),
    cell: ({ row }) => h('div', { class: 'text-right' }, formatHours(row.original.durationSeconds))
  },
  {
    id: 'billableAmount',
    header: () => h('div', { class: 'text-right' }, t('timesheet.columns.billable')),
    cell: ({ row }) => {
      if (!row.original.billable) return h('span', { class: 'text-muted-foreground block text-right' }, '-')
      const amount = getEntryAmount(row.original).toFixed(2)
      return h('span', { class: 'billable-amount block text-right font-medium text-primary' }, `$${amount}`)
    }
  },
  {
    id: 'actions',
    enableHiding: false,
    cell: ({ row }) => {
      const entry = row.original
      return h('div', { class: 'flex gap-1 justify-end' }, [
        h(Button, {
          variant: 'ghost',
          size: 'icon',
          class: 'h-8 w-8',
          onClick: (e: MouseEvent) => {
            e.stopPropagation()
            handleEdit(entry)
          }
        }, () => h(Edit, { class: 'w-4 h-4' })),
        h(AlertDialog, {}, {
          default: () => [
            h(AlertDialogTrigger, { asChild: true }, () => 
              h(Button, {
                variant: 'ghost',
                size: 'icon',
                class: 'h-8 w-8 text-destructive hover:text-destructive',
                onClick: (e: MouseEvent) => e.stopPropagation()
              }, () => h(Trash2, { class: 'w-4 h-4' }))
            ),
            h(AlertDialogContent, {}, {
              default: () => [
                h(AlertDialogHeader, {}, {
                  default: () => [
                    h(AlertDialogTitle, {}, () => t('common.confirmDelete')),
                    h(AlertDialogDescription, {}, () => t('timesheet.entry.deleteConfirm'))
                  ]
                }),
                h(AlertDialogFooter, {}, {
                  default: () => [
                    h(AlertDialogCancel, {}, () => t('common.cancel')),
                    h(AlertDialogAction, { onClick: () => handleDelete(entry.id) }, () => t('common.delete'))
                  ]
                })
              ]
            })
          ]
        })
      ])
    }
  }
])

onMounted(() => {
  timesheetStore.fetchTimesheet()
  projectStore.fetchProjects()
})
</script>

<template>
  <div class="h-full flex flex-col min-h-0">
    <!-- Edit Modal -->
    <TimesheetFormModal v-if="showModal" :show="showModal" :entry="editingEntry" :projects="projects"
      @submit="handleSubmitEntry" @update:show="showModal = $event" />


    <!-- Time Entries Section -->
    <div class="flex-1 min-h-0 flex flex-col overflow-hidden">
      <!-- Section Header -->
      <div class="shrink-0 flex justify-between items-center mb-2">
        <h3 class="text-lg font-semibold">{{ t('timesheet.entries.title') }}</h3>
        <div class="flex items-center gap-2">
          <Button size="sm" @click="handleCreateEntry">
            {{ t('timesheet.logTime') }}
          </Button>
          <Button variant="ghost" size="sm" @click="handleExportCSV">
            <Download class="w-4 h-4 mr-2" />
            {{ t('timesheet.entries.exportCSV') }}
          </Button>
        </div>
      </div>

      <!-- Quick Entry Bar -->
      <QuickTimeEntry :projects="projects" @submit="handleQuickEntry" class="shrink-0 mb-4" />

      <!-- Data Table -->
      <div v-if="loading" class="flex-1 flex items-center justify-center">
        <Loader2 class="w-6 h-6 animate-spin text-muted-foreground" />
      </div>

      <div v-else-if="entries.length === 0"
        class="flex-1 flex flex-col items-center justify-center border rounded-lg bg-muted/20">
        <Clock class="w-12 h-12 text-muted-foreground mb-4" />
        <p class="text-muted-foreground">{{ t('timesheet.noEntries') }}</p>
        <p class="text-xs text-muted-foreground mt-1">{{ t('timesheet.noEntriesHint') }}</p>
      </div>

      <div v-else class="flex-1 min-h-0 overflow-hidden">
        <DataTable :columns="columns" :data="enrichedEntries" />
      </div>
    </div>

  </div>
</template>

<style scoped>
/* Section Header */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.section-title {
  font-size: 1rem;
}

/* Entries Section */
.entries-section {
  margin-top: 16px;
}

.entries-section :deep(.n-card__content) {
  padding: 16px;
}

/* Table Styling */
.entries-table {
  margin-top: 16px;
}

.entries-table :deep(.n-data-table-th) {
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--n-text-color-3);
}

.entries-table :deep(.n-data-table-td) {
  font-size: 0.875rem;
}

/* Project Cell */
.project-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.project-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

/* Billable Amount */
.billable-amount {
  color: var(--n-primary-color);
  font-weight: 500;
}

.text-muted {
  color: var(--n-text-color-3);
}

/* Action Buttons */
.action-buttons {
  display: flex;
  gap: 4px;
}
</style>
