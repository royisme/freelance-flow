<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
// Naive UI components removed
import ProjectFormModal from '@/components/ProjectFormModal.vue'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table'
import { useProjectStore } from '@/stores/projects'
import { useClientStore } from '@/stores/clients'
import { useTimesheetStore } from '@/stores/timesheet'
import { storeToRefs } from 'pinia'
import { useI18n } from 'vue-i18n'
import type { Project, Client, TimeEntry } from '@/types'
import { ArrowLeft, Edit, Loader2 } from 'lucide-vue-next'
import { toast } from 'vue-sonner'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()

const projectStore = useProjectStore()
const clientStore = useClientStore()
const timesheetStore = useTimesheetStore()
const { projects, loading: projectLoading } = storeToRefs(projectStore)
const { clients } = storeToRefs(clientStore)
const { entries: timeEntries, loading: entriesLoading } = storeToRefs(timesheetStore)

// Modal state
const showModal = ref(false)

// Computed: Current project
const projectId = computed(() => Number(route.params.id))
const project = computed<Project | undefined>(() =>
    projects.value.find(p => p.id === projectId.value)
)
const client = computed<Client | undefined>(() =>
    project.value ? clients.value.find(c => c.id === project.value!.clientId) : undefined
)

// Computed: Time entries for this project
const projectTimeEntries = computed<TimeEntry[]>(() =>
    timeEntries.value.filter(e => e.projectId === projectId.value)
)
const totalHours = computed(() =>
    projectTimeEntries.value.reduce((sum, e) => sum + (e.durationSeconds / 3600), 0).toFixed(1)
)

const loading = computed(() => projectLoading.value || entriesLoading.value)

// Data fetch
onMounted(() => {
    projectStore.fetchProjects()
    clientStore.fetchClients()
    timesheetStore.fetchTimesheet()
})

function handleBack() {
    router.push('/projects')
}

function handleEdit() {
    showModal.value = true
}

async function handleSubmitProject(projectData: Omit<Project, 'id'> | Project) {
    try {
        if ('id' in projectData) {
            await projectStore.updateProject(projectData)
            toast.success(t('projects.updateSuccess'))
        }
    } catch (error) {
        toast.error(t('projects.saveError'))
    }
}
</script>

<template>
    <div class="h-full flex flex-col min-h-0">
        <div v-if="loading" class="flex items-center justify-center h-full">
            <Loader2 class="w-8 h-8 animate-spin text-muted-foreground" />
        </div>

        <template v-else-if="project">
            <!-- Secondary Page Header: Back + Breadcrumb -->
            <div class="flex items-center gap-2 mb-4 shrink-0">
                <Button variant="ghost" size="icon" class="h-8 w-8" @click="handleBack">
                    <ArrowLeft class="w-4 h-4" />
                </Button>
                <nav class="flex items-center text-sm text-muted-foreground">
                    <span class="hover:text-foreground cursor-pointer transition-colors" @click="handleBack">{{
                        t('nav.projects') }}</span>
                    <span class="mx-2">/</span>
                    <span class="font-medium text-foreground">{{ project.name }}</span>
                </nav>
            </div>

            <ProjectFormModal v-model:show="showModal" :project="project" :clients="clients"
                @submit="handleSubmitProject" />

            <div class="flex flex-col gap-4 flex-1 min-h-0">
                <!-- Project Info Card -->
                <Card class="shrink-0">
                    <CardHeader class="pb-2">
                        <div class="flex items-center justify-between w-full">
                            <h3 class="font-semibold leading-none tracking-tight">{{ t('projects.detail.info') }}
                            </h3>
                            <Button size="sm" @click="handleEdit">
                                <Edit class="w-4 h-4 mr-2" />
                                {{ t('common.edit') }}
                            </Button>
                        </div>
                    </CardHeader>
                    <CardContent class="pt-0">
                        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4 mt-4">
                            <div>
                                <p class="text-sm font-medium text-muted-foreground">{{ t('form.project.client') }}
                                </p>
                                <p class="text-sm font-medium mt-1">{{ client?.name || '-' }}</p>
                            </div>
                            <div>
                                <p class="text-sm font-medium text-muted-foreground">{{ t('form.project.status') }}
                                </p>
                                <Badge
                                    :variant="project.status === 'active' ? 'default' : (project.status === 'archived' ? 'secondary' : 'outline')"
                                    class="mt-1 capitalize rounded-full">
                                    {{ t(`projects.status.${project.status}`) }}
                                </Badge>
                            </div>
                            <div>
                                <p class="text-sm font-medium text-muted-foreground">{{
                                    t('projects.detail.totalHours') }}</p>
                                <p class="text-sm font-medium mt-1">{{ totalHours }}h</p>
                            </div>
                            <div>
                                <p class="text-sm font-medium text-muted-foreground">{{ t('form.project.hourlyRate')
                                }}</p>
                                <p class="text-sm font-medium mt-1">{{ project.currency }} ${{ project.hourlyRate
                                }}/hr</p>
                            </div>
                            <div>
                                <p class="text-sm font-medium text-muted-foreground">{{ t('form.project.deadline')
                                }}</p>
                                <p class="text-sm font-medium mt-1">{{ project.deadline ||
                                    t('projects.columns.noDeadline') }}</p>
                            </div>
                            <div class="col-span-full">
                                <p class="text-sm font-medium text-muted-foreground">{{
                                    t('form.project.description') }}</p>
                                <p class="text-sm text-muted-foreground mt-1 whitespace-pre-wrap">{{
                                    project.description || '-' }}</p>
                            </div>
                        </div>
                    </CardContent>
                </Card>

                <!-- Time Entries -->
                <Card class="flex-1 min-h-0 flex flex-col">
                    <CardHeader class="pb-2">
                        <h3 class="font-semibold leading-none tracking-tight">{{ t('projects.detail.timeEntries') }}
                        </h3>
                    </CardHeader>
                    <CardContent class="p-0 flex-1 min-h-0 overflow-auto">
                        <div class="rounded-md border m-4 mt-0">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead class="w-[120px]">{{ t('timesheet.columns.date') }}</TableHead>
                                        <TableHead>{{ t('timesheet.columns.task') }}</TableHead>
                                        <TableHead class="w-[100px]">{{ t('timesheet.columns.hours') }}</TableHead>
                                        <TableHead class="w-[100px]">{{ t('timesheet.columns.billable') }}
                                        </TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    <template v-if="projectTimeEntries.length > 0">
                                        <TableRow v-for="entry in projectTimeEntries" :key="entry.id">
                                            <TableCell>{{ entry.date ? entry.date.substring(0, 10) : '-' }}
                                            </TableCell>
                                            <TableCell>{{ entry.description }}</TableCell>
                                            <TableCell>{{ (entry.durationSeconds / 3600).toFixed(1) }}h</TableCell>
                                            <TableCell>
                                                <Badge :variant="entry.billable ? 'default' : 'secondary'"
                                                    class="rounded-full px-2">
                                                    {{ entry.billable ? t('timesheet.entries.billable') :
                                                        t('timesheet.entries.nonBillable') }}
                                                </Badge>
                                            </TableCell>
                                        </TableRow>
                                    </template>
                                    <TableRow v-else>
                                        <TableCell colspan="4" class="h-24 text-center">
                                            {{ t('common.noData') }}
                                        </TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </template>

        <template v-else>
            <div class="flex flex-col items-center justify-center py-12 gap-4">
                <p class="text-muted-foreground">{{ t('projects.detail.notFound') }}</p>
                <Button @click="handleBack">{{ t('nav.projects') }}</Button>
            </div>
        </template>
    </div>
</template>

<style scoped>
/* No extra styles needed */
</style>
