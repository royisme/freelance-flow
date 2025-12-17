<script setup lang="ts">
import { h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import type { ColumnDef } from '@tanstack/vue-table'
import { storeToRefs } from 'pinia'
import { toast } from 'vue-sonner'

import { useProjectStore } from '@/stores/projects'
import { useClientStore } from '@/stores/clients'
import type { Project } from '@/types'


import ProjectFormModal from '@/components/ProjectFormModal.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { DataTable, DataTableColumnHeader } from '@/components/ui/data-table'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Plus, Edit, FolderOpen, Trash2 } from 'lucide-vue-next'

const router = useRouter()
const projectStore = useProjectStore()
const clientStore = useClientStore()
const { projects } = storeToRefs(projectStore)
const { clients } = storeToRefs(clientStore)
const { t } = useI18n()

const showModal = ref(false)
const editingProject = ref<Project | null>(null)

// Delete Dialog State
const showDeleteDialog = ref(false)
const projectToDelete = ref<Project | null>(null)

function handleNewProject() {
  editingProject.value = null
  showModal.value = true
}

function handleEditProject(project: Project) {
  editingProject.value = project
  showModal.value = true
}

function confirmDeleteProject(project: Project) {
  projectToDelete.value = project
  showDeleteDialog.value = true
}

async function handleDeleteProject() {
  if (!projectToDelete.value) return
  try {
    await projectStore.deleteProject(projectToDelete.value.id)
    toast.success(t('projects.deleteSuccess'))
  } catch (error) {
    toast.error(t('projects.deleteError'))
  } finally {
    showDeleteDialog.value = false
    projectToDelete.value = null
  }
}

async function handleSubmitProject(project: Omit<Project, 'id'> | Project) {
  try {
    if ('id' in project) {
      await projectStore.updateProject(project)
      toast.success(t('projects.updateSuccess'))
    } else {
      await projectStore.createProject(project)
      toast.success(t('projects.createSuccess'))
    }
  } catch (error) {
    toast.error(t('projects.saveError'))
  }
}

onMounted(() => {
  projectStore.fetchProjects()
  clientStore.fetchClients()
})

const columns: ColumnDef<Project>[] = [
  {
    accessorKey: 'name',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('projects.columns.projectName') }),
    cell: ({ row }) => {
      const p = row.original
      return h('div', [
        h('div', { class: 'font-semibold' }, p.name),
        h('span', { class: 'text-xs text-muted-foreground' }, p.description || '-')
      ])
    }
  },
  {
    accessorKey: 'status',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('projects.columns.status') }),
    cell: ({ row }) => {
      const status = row.getValue('status') as string
      let variant: 'default' | 'secondary' | 'outline' = 'outline'
      if (status === 'active') variant = 'default'
      if (status === 'archived') variant = 'secondary'

      return h(Badge, { variant, class: 'capitalize rounded-full' }, () => t(`projects.status.${status}`))
    }
  },
  {
    accessorKey: 'progress',
    header: t('projects.columns.progress'),
    cell: () => {
      const randomProgress = Math.floor(Math.random() * 40) + 30
      return h(Progress, { modelValue: randomProgress, class: 'h-2 w-[100px]' })
    }
  },
  {
    accessorKey: 'deadline',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('projects.columns.deadline') }),
    cell: ({ row }) => {
      return row.original.deadline || t('projects.columns.noDeadline')
    }
  },
  {
    accessorKey: 'hourlyRate',
    header: ({ column }) => h(DataTableColumnHeader, { column: column as any, title: t('projects.columns.hourlyRate') }),
    cell: ({ row }) => {
      return h('span', { class: 'font-medium' }, `${row.original.currency} $${row.original.hourlyRate}/hr`)
    }
  },
  {
    id: 'actions',
    header: t('projects.columns.actions'),
    cell: ({ row }) => {
      const project = row.original
      return h('div', { class: 'flex gap-1' }, [
        h(Button, {
          size: 'icon',
          variant: 'ghost',
          class: 'h-8 w-8',
          onClick: (e: MouseEvent) => {
            e.stopPropagation() // Prevent row click if any
            router.push(`/projects/${project.id}`)
          }
        }, () => h(FolderOpen, { class: 'w-4 h-4' })),
        h(Button, {
          size: 'icon',
          variant: 'ghost',
          class: 'h-8 w-8',
          onClick: (e: MouseEvent) => {
            e.stopPropagation()
            handleEditProject(project)
          }
        }, () => h(Edit, { class: 'w-4 h-4' })),
        h(Button, {
          size: 'icon',
          variant: 'ghost',
          class: 'h-8 w-8 text-destructive hover:text-destructive',
          onClick: (e: MouseEvent) => {
            e.stopPropagation()
            confirmDeleteProject(project)
          }
        }, () => h(Trash2, { class: 'w-4 h-4' }))
      ])
    }
  }
]
</script>

<template>
  <div class="h-full flex flex-col min-h-0">
    <!-- Header Actions -->
    <div class="shrink-0 flex items-center justify-end mb-4">
      <Button @click="handleNewProject">
        <Plus class="w-4 h-4 mr-2" />
        {{ t('projects.addProject') }}
      </Button>
    </div>

    <!-- Content -->
    <div class="flex-1 min-h-0 overflow-hidden">
      <DataTable :columns="columns" :data="projects" />
    </div>

    <ProjectFormModal v-model:show="showModal" :project="editingProject" :clients="clients"
      @submit="handleSubmitProject" />

    <!-- Delete Confirmation Dialog -->
    <Dialog v-model:open="showDeleteDialog">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{{ t('projects.deleteTitle') }}</DialogTitle>
          <DialogDescription>
            {{ t('projects.deleteConfirm', { name: projectToDelete?.name }) }}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="showDeleteDialog = false">
            {{ t('common.cancel') }}
          </Button>
          <Button variant="destructive" @click="handleDeleteProject">
            {{ t('common.delete') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>
