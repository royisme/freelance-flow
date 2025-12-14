<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import {
  NButton, NDataTable, NTag, NSpace, NProgress, NText, NIcon, type DataTableColumns,
  useMessage, useDialog
} from 'naive-ui'
import PageContainer from '@/components/PageContainer.vue'
import PageHeader from '@/components/PageHeader.vue'
import ProjectFormModal from '@/components/ProjectFormModal.vue'
import { useProjectStore } from '@/stores/projects'
import { useClientStore } from '@/stores/clients'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { Project } from '@/types'
import { PlusOutlined, EditOutlined, DeleteOutlined, FolderOpenOutlined } from '@vicons/antd'

const message = useMessage()
const dialog = useDialog()
const projectStore = useProjectStore()
const clientStore = useClientStore()
const { projects, loading } = storeToRefs(projectStore)
const { clients } = storeToRefs(clientStore)
const { t } = useI18n()

const showModal = ref(false)
const editingProject = ref<Project | null>(null)

function handleNewProject() {
  editingProject.value = null
  showModal.value = true
}

function handleEditProject(project: Project) {
  editingProject.value = project
  showModal.value = true
}

function handleDeleteProject(project: Project) {
  dialog.warning({
    title: t('projects.deleteTitle'),
    content: t('projects.deleteConfirm', { name: project.name }),
    positiveText: t('common.delete'),
    negativeText: t('common.cancel'),
    onPositiveClick: async () => {
      try {
        await projectStore.deleteProject(project.id)
        message.success(t('projects.deleteSuccess'))
      } catch (error) {
        message.error(t('projects.deleteError'))
      }
    }
  })
}

async function handleSubmitProject(project: Omit<Project, 'id'> | Project) {
  try {
    if ('id' in project) {
      await projectStore.updateProject(project)
      message.success(t('projects.updateSuccess'))
    } else {
      await projectStore.createProject(project)
      message.success(t('projects.createSuccess'))
    }
  } catch (error) {
    message.error(t('projects.saveError'))
  }
}

onMounted(() => {
  projectStore.fetchProjects()
  clientStore.fetchClients()
})

const columns: DataTableColumns<Project> = [
  {
    title: () => t('projects.columns.projectName'),
    key: 'name',
    width: 250,
    render(row) {
      return h('div', [
        h('div', { style: 'font-weight: 600;' }, row.name),
        h(NText, { depth: 3, style: 'font-size: 12px;' }, { default: () => row.description || '-' })
      ])
    }
  },
  {
    title: () => t('projects.columns.status'),
    key: 'status',
    width: 120,
    render(row) {
      let type: 'success' | 'warning' | 'default' = 'default'
      if (row.status === 'active') type = 'success'
      if (row.status === 'archived') type = 'warning'

      return h(
        NTag,
        { type, bordered: false, round: true, size: 'small' },
        { default: () => t(`projects.status.${row.status}`) }
      )
    }
  },
  {
    title: () => t('projects.columns.progress'),
    key: 'progress', // Mock progress for UI demo
    render() {
      const randomProgress = Math.floor(Math.random() * 40) + 30
      return h(
        NProgress,
        {
          type: 'line',
          percentage: randomProgress,
          color: 'var(--n-primary-color)',
          height: 6,
          'show-indicator': false
        }
      )
    }
  },
  {
    title: () => t('projects.columns.deadline'),
    key: 'deadline',
    render(row) {
      return row.deadline || t('projects.columns.noDeadline')
    }
  },
  {
    title: () => t('projects.columns.hourlyRate'),
    key: 'hourlyRate',
    render(row) {
      return h(NText, { depth: 1 }, { default: () => `${row.currency} $${row.hourlyRate}/hr` })
    }
  },
  {
    title: () => t('projects.columns.actions'),
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
              onClick: (e) => {
                e.stopPropagation()
                message.info(t('projects.manage', { name: row.name }))
              }
            },
            { icon: () => h(FolderOpenOutlined) }
          ),
          h(
            NButton,
            {
              size: 'small',
              quaternary: true,
              circle: true,
              onClick: (e) => {
                e.stopPropagation()
                handleEditProject(row)
              }
            },
            { icon: () => h(EditOutlined) }
          )
        ]
      })
    }
  }
]
</script>

<template>
  <PageContainer>
    <PageHeader :title="t('projects.title')" :subtitle="t('projects.subtitle')">
      <template #extra>
        <n-button type="primary" @click="handleNewProject">
          <template #icon>
            <n-icon>
              <PlusOutlined />
            </n-icon>
          </template>
          {{ t('projects.addProject') }}
        </n-button>
      </template>
    </PageHeader>

    <ProjectFormModal v-model:show="showModal" :project="editingProject" :clients="clients"
      @submit="handleSubmitProject" />

    <n-data-table :columns="columns" :data="projects" :loading="loading" :bordered="false" class="project-table" />
  </PageContainer>
</template>

<style scoped>
.project-table :deep(.n-data-table-th) {
  font-weight: 600;
  color: var(--n-text-color-2);
}
</style>
