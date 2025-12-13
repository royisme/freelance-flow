import { defineStore, storeToRefs } from "pinia";
import { ref, computed } from "vue";
import { api } from "@/api";
import type {
  TimeEntry,
  CreateTimeEntryInput,
  UpdateTimeEntryInput,
  Project,
} from "@/types";
import { useProjectStore } from "./projects";
import { useStatusBarStore } from "./statusBar";

// Grouped Data Interface for View
export interface DailyTimeGroup {
  date: string;
  totalSeconds: number;
  entries: (TimeEntry & { project?: Project })[];
}

export const useTimesheetStore = defineStore("timesheet", () => {
  const entries = ref<TimeEntry[]>([]);
  const loading = ref(false);
  const projectStore = useProjectStore();
  const { projects } = storeToRefs(projectStore);
  const statusBarStore = useStatusBarStore();

  // Actions
  async function fetchTimesheet(projectId?: number) {
    loading.value = true;
    try {
      // Ensure we have project data for joining
      if (projects.value.length === 0) {
        await projectStore.fetchProjects();
      }
      entries.value = await api.timeEntries.list(projectId);
      await statusBarStore.refresh();
    } catch (error) {
      console.error("Failed to fetch timesheet", error);
    } finally {
      loading.value = false;
    }
  }

  async function createTimeEntry(input: CreateTimeEntryInput) {
    loading.value = true;
    try {
      await api.timeEntries.create(input);
      await fetchTimesheet();
    } finally {
      loading.value = false;
    }
  }

  async function updateTimeEntry(input: UpdateTimeEntryInput) {
    loading.value = true;
    try {
      await api.timeEntries.update(input);
      await fetchTimesheet();
    } finally {
      loading.value = false;
    }
  }

  async function deleteTimeEntry(id: number) {
    loading.value = true;
    try {
      await api.timeEntries.delete(id);
      await fetchTimesheet();
    } finally {
      loading.value = false;
    }
  }

  // Continue Timer - creates a new time entry with same project and description
  async function continueTimer(entry: TimeEntry & { project?: Project }) {
    loading.value = true;
    try {
      await api.timeEntries.create({
        projectId: entry.projectId,
        description: entry.description,
        date: new Date().toISOString().split('T')[0],
        startTime: '',
        endTime: '',
        durationSeconds: 0,
        billable: entry.billable,
        invoiced: false,
      });
      await fetchTimesheet();
    } finally {
      loading.value = false;
    }
  }

  // CSV Export - exports all entries to CSV format
  function exportToCSV(enrichedEntries: EnrichedTimeEntry[]) {
    const headers = [
      'Date',
      'Project',
      'Description',
      'Start Time',
      'End Time',
      'Duration (hours)',
      'Billable',
      'Invoiced',
    ];

    const rows = enrichedEntries.map((entry) => {
      const durationHours = (entry.durationSeconds / 3600).toFixed(2);
      return [
        entry.date,
        entry.project?.name || 'Unknown Project',
        entry.description || '',
        entry.startTime || '',
        entry.endTime || '',
        durationHours,
        entry.billable ? 'Yes' : 'No',
        entry.invoiced ? 'Yes' : 'No',
      ];
    });

    const csvContent = [
      headers.join(','),
      ...rows.map((row) =>
        row.map((cell) => `"${String(cell).replace(/"/g, '""')}"`).join(',')
      ),
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', `timesheet-${new Date().toISOString().split('T')[0]}.csv`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }

  // Getters (Computed)

  // Enriched Entry type for views that need project info
  type EnrichedTimeEntry = TimeEntry & { project?: Project };

  // Enriched Entries (Join with Project)
  const enrichedEntries = computed<EnrichedTimeEntry[]>(() => {
    return entries.value.map((entry): EnrichedTimeEntry => {
      const proj = projects.value.find((p) => p.id === entry.projectId);
      return { ...entry, project: proj };
    });
  });

  // Group by Date for Linear/Journal View
  const groupedByDay = computed<DailyTimeGroup[]>(() => {
    const groups: Record<string, DailyTimeGroup> = {};

    enrichedEntries.value.forEach((entry) => {
      if (!groups[entry.date]) {
        groups[entry.date] = { date: entry.date, totalSeconds: 0, entries: [] };
      }
      // Non-null assertion safe: we just created the group above if it didn't exist
      const group = groups[entry.date]!;
      group.entries.push(entry);
      group.totalSeconds += entry.durationSeconds;
    });

    // Sort descending by date
    return Object.values(groups).sort((a, b) => b.date.localeCompare(a.date));
  });

  const totalHours = computed(() => {
    const totalSecs = entries.value.reduce(
      (acc, curr) => acc + curr.durationSeconds,
      0
    );
    return (totalSecs / 3600).toFixed(1);
  });

  return {
    entries,
    enrichedEntries,
    groupedByDay,
    totalHours,
    loading,
    fetchTimesheet,
    createTimeEntry,
    updateTimeEntry,
    deleteTimeEntry,
    continueTimer,
    exportToCSV,
  };
});
