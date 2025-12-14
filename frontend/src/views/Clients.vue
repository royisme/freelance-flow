<script setup lang="ts">
import { h, onMounted, computed, ref } from 'vue'
import {
  NButton, NDataTable, NTag, NSpace, NText, NIcon,
  type DataTableColumns, useMessage, useDialog
} from 'naive-ui'
import { useRouter } from 'vue-router'
import PageContainer from '@/components/PageContainer.vue'
import PageHeader from '@/components/PageHeader.vue'
import ClientFormModal from '@/components/ClientFormModal.vue'
import { useClientStore } from '@/stores/clients'
import { useProjectStore } from '@/stores/projects'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { Client, Project } from '@/types'
import { PlusOutlined, EditOutlined, DeleteOutlined, RightOutlined } from '@vicons/antd'

const message = useMessage()
const dialog = useDialog()
const router = useRouter()
const clientStore = useClientStore()
const projectStore = useProjectStore()
const { clients, loading: clientsLoading } = storeToRefs(clientStore)
const { projects, loading: projectsLoading } = storeToRefs(projectStore)
const { t } = useI18n()

// Modal State
const showModal = ref(false)
const editingClient = ref<Client | null>(null)

function handleNewClient() {
  editingClient.value = null
  showModal.value = true
}

function handleEditClient(client: Client) {
  editingClient.value = client
  showModal.value = true
}

function handleDeleteClient(client: Client) {
  dialog.warning({
    title: t('clients.deleteTitle'),
    content: t('clients.deleteConfirm', { name: client.name }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await clientStore.deleteClient(client.id)
        message.success(t('clients.deleteSuccess'))
      } catch (error) {
        message.error(t('clients.deleteError'))
      }
    }
  })
}

async function handleSubmitClient(client: Omit<Client, 'id'> | Client) {
  try {
    if ('id' in client) {
      await clientStore.updateClient(client)
      message.success(t('clients.updateSuccess'))
    } else {
      await clientStore.createClient(client)
      message.success(t('clients.createSuccess'))
    }
  } catch (error) {
    message.error(t('clients.saveError'))
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
const columns: DataTableColumns<TreeRow> = [
  {
    title: () => t('clients.columns.clientName'),
    key: 'name',
    render(row) {
      if (row.type === 'client') {
        const projectCount = row.children?.length || 0
        // Use inline-flex to keep content on the same line as the tree expand icon
        return h('span', { style: 'display: inline-flex; flex-direction: column;' }, [
          h('span', { style: 'font-weight: 600;' }, row.name),
          h(NText, { depth: 3, style: 'font-size: 11px;' }, { default: () => t('clients.columns.projectsCount', { count: projectCount }) })
        ])
      } else {
        // Project row - NDataTable handles indentation automatically
        return h('span', { style: 'display: inline-flex; flex-direction: column;' }, [
          h('span', { style: 'font-weight: 500;' }, row.name),
          h(NText, { depth: 3, style: 'font-size: 11px;' }, { default: () => row.description })
        ])
      }
    }
  },
  {
    title: () => t('clients.columns.contactPerson'),
    key: 'contactPerson',
    width: 180,
    render(row) {
      // Only show for clients
      return row.type === 'client' ? (row.contactPerson || '-') : ''
    }
  },
  {
    title: () => t('clients.columns.status'),
    key: 'status',
    width: 100,
    render(row) {
      const statusType = row.status === 'active' ? 'success' : (row.status === 'archived' ? 'warning' : 'default')
      const statusKey = row.type === 'client' ? `clients.status.${row.status}` : `projects.status.${row.status}`
      return h(
        NTag,
        {
          type: statusType,
          bordered: false,
          round: true,
          size: 'small'
        },
        { default: () => t(statusKey) }
      )
    }
  },
  {
    title: () => t('clients.columns.actions'),
    key: 'actions',
    width: 100,
    render(row) {
      if (row.type === 'client') {
        // Client actions: Edit & Delete
        const client = clients.value.find(c => c.id === row.clientId)
        if (!client) return null
        return h(NSpace, { size: 'small' }, {
          default: () => [
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                circle: true,
                onClick: (e) => { e.stopPropagation(); handleEditClient(client) }
              },
              { icon: () => h(EditOutlined) }
            ),
            h(
              NButton,
              {
                size: 'small',
                quaternary: true,
                circle: true,
                type: 'error',
                onClick: (e) => { e.stopPropagation(); handleDeleteClient(client) }
              },
              { icon: () => h(DeleteOutlined) }
            )
          ]
        })
      } else {
        // Project action: View details
        return h(
          NButton,
          {
            size: 'small',
            quaternary: true,
            circle: true,
            onClick: () => handleViewProject(row.projectId!)
          },
          { icon: () => h(RightOutlined) }
        )
      }
    }
  }
]

// Row class for styling
function rowClassName(row: TreeRow) {
  return row.type === 'project' ? 'project-row' : 'client-row'
}
</script>

<template>
  <PageContainer fill>
    <PageHeader :title="t('clients.title')" :subtitle="t('clients.subtitle')">
      <template #extra>
        <n-button type="primary" @click="handleNewClient">
          <template #icon>
            <n-icon>
              <PlusOutlined />
            </n-icon>
          </template>
          {{ t('clients.addClient') }}
        </n-button>
      </template>
    </PageHeader>

    <ClientFormModal v-model:show="showModal" :client="editingClient" @submit="handleSubmitClient" />

    <n-data-table :columns="columns" :data="treeData" :loading="loading" :bordered="false"
      :row-key="(row: TreeRow) => row.key" children-key="children" :row-class-name="rowClassName"
      class="client-table" />
  </PageContainer>
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
