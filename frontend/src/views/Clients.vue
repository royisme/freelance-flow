<script setup lang="ts">
import { h, onMounted, computed, ref } from 'vue'
import {
  NButton, NDataTable, NTag, NSpace, NAvatar, NText, NCard, NGrid, NGridItem, NIcon,
  type DataTableColumns, useMessage, useDialog
} from 'naive-ui'
import PageContainer from '@/components/PageContainer.vue'
import PageHeader from '@/components/PageHeader.vue'
import ClientFormModal from '@/components/ClientFormModal.vue'
import { useClientStore } from '@/stores/clients'
import { useProjectStore } from '@/stores/projects'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { Client, Project } from '@/types'
import { PlusOutlined, EditOutlined, DeleteOutlined, FolderOpenOutlined } from '@vicons/antd'

const message = useMessage()
const dialog = useDialog()
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

// Sub-Component: Project List (Rendered in Expanded Row)
const ProjectListRenderer = (props: { projects: Project[] }) => {
  if (props.projects.length === 0) {
    return h(NText, { depth: 3, style: 'padding: 12px 0; display: block;' }, { default: () => t('projects.noProjects') })
  }

  return h(NGrid, { cols: 2, xGap: 12, yGap: 12 }, {
    default: () => props.projects.map(project =>
      h(NGridItem, null, {
        default: () => h(NCard, { size: 'small', bordered: true, style: 'background-color: var(--n-action-color);' }, {
          header: () => h(NText, { strong: true }, { default: () => project.name }),
          headerExtra: () => h(NTag, {
            type: project.status === 'active' ? 'success' : 'default',
            size: 'tiny',
            bordered: false,
            round: true
          }, { default: () => t(`projects.status.${project.status}`) }),
          default: () => h(NText, { depth: 3, style: 'font-size: 12px' }, { default: () => project.description || '-' })
        })
      })
    )
  })
}

// Table Columns Definition
const columns: DataTableColumns<Client> = [
  {
    type: 'expand',
    renderExpand: (row) => {
      const clientProjects = getClientProjects(row.id)
      return h(
        'div',
        { style: 'padding: 12px 24px 24px 60px;' }, // Indent content
        [
          h(NText, { strong: true, style: 'margin-bottom: 8px; display: block;' }, { default: () => t('clients.columns.associatedProjects') }),
          h(ProjectListRenderer, { projects: clientProjects })
        ]
      )
    }
  },
  {
    title: () => t('clients.columns.clientName'),
    key: 'name',
    width: 280,
    render(row) {
      const projectCount = getClientProjects(row.id).length
      return h(
        NSpace,
        { align: 'center', size: 'small' },
        {
          default: () => [
            h(NAvatar, {
              size: 'small',
              src: row.avatar,
              fallbackSrc: 'https://07akioni.oss-cn-hangzhou.aliyuncs.com/07akioni.jpeg',
              round: true,
              style: 'margin-right: 8px'
            }),
            h('div', [
              h('div', { style: 'font-weight: 600;' }, row.name),
              h(NText, { depth: 3, style: 'font-size: 11px;' }, { default: () => t('clients.columns.projectsCount', { count: projectCount }) })
            ])
          ]
        }
      )
    }
  },
  {
    title: () => t('clients.columns.contactPerson'),
    key: 'contactPerson',
    render(row) {
      return row.contactPerson || '-'
    }
  },
  {
    title: () => t('clients.columns.status'),
    key: 'status',
    width: 100,
    render(row) {
      return h(
        NTag,
        {
          type: row.status === 'active' ? 'success' : 'default',
          bordered: false,
          round: true,
          size: 'small'
        },
        { default: () => t(`clients.status.${row.status}`) }
      )
    }
  },
  {
    title: () => t('clients.columns.actions'),
    key: 'actions',
    width: 140,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              circle: true,
              onClick: (e) => { e.stopPropagation(); handleEditClient(row); }
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
              onClick: (e) => { e.stopPropagation(); handleDeleteClient(row); }
            },
            { icon: () => h(DeleteOutlined) }
          )
        ]
      })
    }
  }
]
</script>

<template>
  <PageContainer>
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

    <n-data-table :columns="columns" :data="clients" :loading="loading" :bordered="false" class="client-table"
      :row-key="(row) => row.id" />
  </PageContainer>
</template>

<style scoped>
/* Scoped styles if needed */
.client-table :deep(.n-data-table-th) {
  font-weight: 600;
  color: var(--n-text-color-2);
}

/* Ensure clean nested look */
.client-table :deep(.n-data-table-td--expand) {
  background-color: var(--n-body-color) !important;
}
</style>
