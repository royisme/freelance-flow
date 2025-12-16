```vue
<script setup lang="ts">
import { h, onMounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import type { ColumnDef } from '@tanstack/vue-table'
import ClientFormModal from '@/components/ClientFormModal.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { DataTable } from '@/components/ui/data-table'
import { DataTableColumnHeader } from '@/components/ui/data-table'
import { useClientStore } from '@/stores/clients'
import { useProjectStore } from '@/stores/projects'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { Client, Project } from '@/types'
import { Plus, Edit, Trash2, ChevronRight, ChevronDown } from 'lucide-vue-next'

const router = useRouter()
const clientStore = useClientStore()
const projectStore = useProjectStore()
const { clients, loading: clientsLoading } = storeToRefs(clientStore)
const { projects, loading: projectsLoading } = storeToRefs(projectStore)
const { t } = useI18n()

// Modal State
const showModal = ref(false)
const editingClient = ref<Client | null>(null)

// Alert Dialog State
const showDeleteDialog = ref(false)
const clientToDelete = ref<Client | null>(null)

function handleNewClient() {
  editingClient.value = null
  showModal.value = true
}

function handleEditClient(client: Client) {
  editingClient.value = client
  showModal.value = true
}

function confirmDeleteClient(client: Client) {
  clientToDelete.value = client
  showDeleteDialog.value = true
}

async function handleDeleteClient() {
  if (!clientToDelete.value) return

  try {
    await clientStore.deleteClient(clientToDelete.value.id)
    toast.success(t('clients.deleteSuccess'))
    showDeleteDialog.value = false
  } catch (error) {
    toast.error(t('clients.deleteError'))
  }
}

async function handleSubmitClient(client: Omit<Client, 'id'> | Client) {
  try {
    if ('id' in client) {
      await clientStore.updateClient(client)
      toast.success(t('clients.updateSuccess'))
    } else {
      await clientStore.createClient(client)
      toast.success(t('clients.createSuccess'))
    }
  } catch (error) {
    toast.error(t('clients.saveError'))
  }
}

onMounted(() => {
  clientStore.fetchClients()
  projectStore.fetchProjects()
})

// Combined loading state
const loading = computed(() => clientsLoading.value || projectsLoading.value)

// Helper: Get projects for a specific client
function getClientProjects(clientId: number): Project[] {
  return projects.value.filter(p => p.clientId === clientId)
}

// Navigate to project detail
function handleViewProject(projectId: number) {
  router.push(`/projects/${projectId}`)
}

// Type for tree data rows
interface TreeRow {
  key: string
  type: 'client' | 'project'
  name: string
  contactPerson?: string
  status: string
  clientId?: number
  projectId?: number
  description?: string
  children?: TreeRow[]
}

// Build tree data: Clients as parents, Projects as children
const treeData = computed<TreeRow[]>(() => {
  return clients.value.map(client => {
    const clientProjects = getClientProjects(client.id)
    return {
      key: `client-${client.id}`,
      type: 'client' as const,
      name: client.name,
      contactPerson: client.contactPerson || '-',
      status: client.status,
      clientId: client.id,
      children: clientProjects.map(project => ({
        key: `project-${project.id}`,
        type: 'project' as const,
        name: project.name,
        description: project.description || '-',
        status: project.status,
        projectId: project.id,
      }))
    }
  })
})

// Table Columns for tree data
const columns: ColumnDef<TreeRow>[] = [
  {
    accessorKey: 'name',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('clients.columns.clientName') }),
    cell: ({ row }) => {
      const isClient = row.original.type === 'client'
      const projectCount = row.original.children?.length || 0

      return h('div', {
        class: 'flex items-center gap-2',
        style: { paddingLeft: `${row.depth * 2}rem` }
      }, [
        // Expand toggle button for clients
        isClient
          ? h(Button, {
            variant: 'ghost',
            size: 'icon',
            class: 'h-6 w-6 p-0 shrink-0',
            onClick: row.getToggleExpandedHandler(),
          }, {
            default: () => row.getIsExpanded()
              ? h(ChevronDown, { class: 'h-4 w-4' })
              : h(ChevronRight, { class: 'h-4 w-4' })
          })
          : h('div', { class: 'w-6 h-6 shrink-0' }), // Spacer for alignment

        h('div', { class: 'flex flex-col' }, [
          h('span', { class: isClient ? 'font-semibold' : 'font-medium' }, row.original.name),
          isClient
            ? h('span', { class: 'text-xs text-muted-foreground' }, t('clients.columns.projectsCount', { count: projectCount }))
            : h('span', { class: 'text-xs text-muted-foreground' }, row.original.description)
        ])
      ])
    }
  },
  {
    accessorKey: 'contactPerson',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('clients.columns.contactPerson') }),
    cell: ({ row }) => {
      // Only show for clients
      return row.original.type === 'client' ? (row.original.contactPerson || '-') : ''
    }
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('clients.columns.status') }),
    cell: ({ row }) => {
      const variant = row.original.status === 'active' ? 'default' : (row.original.status === 'archived' ? 'secondary' : 'outline')
      const statusKey = row.original.type === 'client'
        ? `clients.status.${row.original.status}`
        : `projects.status.${row.original.status}`

      return h(Badge, { variant, class: 'capitalize' }, { default: () => t(statusKey) })
    }
  },
  {
    id: 'actions',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('clients.columns.actions') }),
    cell: ({ row }) => {
      if (row.original.type === 'client') {
        const client = clients.value.find(c => c.id === row.original.clientId)
        if (!client) return null

        return h('div', { class: 'flex gap-1' }, [
          h(Button, {
            size: 'icon',
            variant: 'ghost',
            class: 'h-8 w-8',
            onClick: (e: MouseEvent) => {
              e.stopPropagation()
              handleEditClient(client)
            }
          }, { default: () => h(Edit, { class: 'w-4 h-4' }) }),

          h(Button, {
            size: 'icon',
            variant: 'ghost',
            class: 'h-8 w-8 text-destructive hover:text-destructive',
            onClick: (e: MouseEvent) => {
              e.stopPropagation()
              confirmDeleteClient(client)
            }
          }, { default: () => h(Trash2, { class: 'w-4 h-4' }) })
        ])
      } else {
        // Project actions
        const projectId = row.original.projectId
        if (!projectId) return null

        return h(Button, {
          size: 'sm',
          variant: 'ghost',
          onClick: () => handleViewProject(projectId)
        }, { default: () => t('common.view') })
      }
    }
  }
]
</script>

<template>
  <div class="space-y-6">
    <!-- Header Actions -->
    <div class="flex items-center justify-end">
      <Button @click="handleNewClient">
        <Plus class="w-4 h-4 mr-2" />
        {{ t('clients.addClient') }}
      </Button>
    </div>

    <!-- Content -->
    <div class="border rounded-xl bg-card text-card-foreground shadow-sm">
      <DataTable :columns="columns" :data="treeData" :loading="loading" :get-sub-rows="(row) => row.children" />
    </div>

    <!-- Client Form Modal -->
    <ClientFormModal v-model:show="showModal" :edit-client="editingClient" :grid-layout="true"
      @submit="handleSubmitClient" />

    <!-- Delete Confirmation Dialog -->
    <AlertDialog :open="showDeleteDialog" @update:open="showDeleteDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ t('clients.deleteTitle') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ t('clients.deleteConfirm', { name: clientToDelete?.name }) }}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel @click="showDeleteDialog = false">{{ t('common.cancel') }}</AlertDialogCancel>
          <AlertDialogAction @click="handleDeleteClient"
            class="bg-destructive text-destructive-foreground hover:bg-destructive/90">
            {{ t('common.delete') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>

<style scoped>
.client-table :deep(.n-data-table-th) {
  font-weight: 600;
  color: var(--n-text-color-2);
}
</style>

<style>
/* Non-scoped styles for row class styling - following official Naive UI pattern */
.project-row td {
  background-color: rgba(0, 128, 0, 0.04) !important;
}
</style>
