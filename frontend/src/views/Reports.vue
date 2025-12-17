<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import { api } from "@/api";
import type { Client, Project, ReportFilter, ReportOutput } from "@/types";
import VChart from "vue-echarts";
import type { EChartsOption } from "echarts";
import { useI18n } from "vue-i18n";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from "@/components/ui/select";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import DateRangePicker from "@/components/DateRangePicker.vue";
import { Loader2, BoxSelect } from "lucide-vue-next";
import { getLocalTimeZone, type DateValue } from '@internationalized/date';

type DateRange = {
  start: DateValue | undefined
  end: DateValue | undefined
}

const { t } = useI18n();

const loading = ref(false);
const error = ref<string | null>(null);
const report = ref<ReportOutput | null>(null);

const clients = ref<Client[]>([]);
const projects = ref<Project[]>([]);

const dateRange = ref<DateRange | undefined>();
const selectedClientId = ref<string>("");
const selectedProjectId = ref<string>("");

const clientOptions = computed(() =>
  clients.value.map((c) => ({ label: c.name, value: String(c.id) }))
);
const filteredProjectOptions = computed(() => {
  const filtered = selectedClientId.value
    ? projects.value.filter((p) => String(p.clientId) === selectedClientId.value)
    : projects.value;
  return filtered.map((p) => ({ label: p.name, value: String(p.id) }));
});

function formatDate(date: Date | number): string {
  const d = new Date(date);
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

async function loadFilters() {
  clients.value = await api.clients.list();
  projects.value = await api.projects.list();
}

async function loadReport() {
  loading.value = true;
  error.value = null;
  try {
    const filter: ReportFilter = {};
    if (dateRange.value?.start && dateRange.value?.end) {
      filter.startDate = formatDate(dateRange.value.start.toDate(getLocalTimeZone()));
      filter.endDate = formatDate(dateRange.value.end.toDate(getLocalTimeZone()));
    }
    if (selectedClientId.value) filter.clientId = Number(selectedClientId.value);
    if (selectedProjectId.value) filter.projectId = Number(selectedProjectId.value);

    report.value = await api.reports.get(filter);
  } catch (e) {
    error.value = e instanceof Error ? e.message : "Failed to load report";
    report.value = null;
  } finally {
    loading.value = false;
  }
}

onMounted(async () => {
  await loadFilters();
  await loadReport();
});

const chartOption = computed<EChartsOption>(() => {
  const c = report.value?.chart;
  return {
    tooltip: { trigger: "axis" },
    legend: {
      data: [t("reports.chart.hours"), t("reports.chart.revenue")],
      top: 0,
    },
    grid: { top: 40, right: 60, bottom: 30, left: 50 },
    xAxis: { type: "category", data: c?.dates ?? [] },
    yAxis: [
      { type: "value", name: t("reports.chart.hours") },
      { type: "value", name: t("reports.chart.revenue") },
    ],
    series: [
      {
        name: t("reports.chart.hours"),
        type: "line",
        smooth: true,
        data: c?.hours ?? [],
        yAxisIndex: 0,
        itemStyle: { color: '#2080f0' },
      },
      {
        name: t("reports.chart.revenue"),
        type: "bar",
        barMaxWidth: 40,
        data: c?.revenue ?? [],
        yAxisIndex: 1,
        itemStyle: { color: '#18a058' },
      },
    ],
  };
});

// Pagination state
const currentPage = ref(1);
const pageSize = 10;

const paginatedRows = computed(() => {
  if (!report.value?.rows) return [];
  const start = (currentPage.value - 1) * pageSize;
  return report.value.rows.slice(start, start + pageSize);
});

const totalPages = computed(() => {
  if (!report.value?.rows) return 0;
  return Math.ceil(report.value.rows.length / pageSize);
});
</script>

<template>
  <div class="h-full flex flex-col min-h-0">
    <!-- Header Stats -->
    <div class="shrink-0 flex justify-end mb-4">
      <div class="flex gap-8">
        <div class="flex flex-col items-end">
          <span class="text-xs text-muted-foreground uppercase font-semibold">{{ t('reports.stats.totalHours')
          }}</span>
          <span class="text-2xl font-bold">{{ report?.totalHours || 0 }}</span>
        </div>
        <div class="flex flex-col items-end">
          <span class="text-xs text-muted-foreground uppercase font-semibold">{{ t('reports.stats.totalIncome')
          }}</span>
          <span class="text-2xl font-bold text-primary">{{ report?.totalIncome || 0 }}</span>
        </div>
      </div>
    </div>

    <div class="flex flex-col gap-4 flex-1 min-h-0">
      <Card class="shrink-0 bg-muted/30 border-0">
        <CardContent class="p-3 flex flex-wrap gap-2 items-center">
          <DateRangePicker v-model="dateRange" class="w-[260px]" />

          <Select v-model="selectedClientId">
            <SelectTrigger class="w-[180px]">
              <SelectValue :placeholder="t('reports.filters.client')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="c in clientOptions" :key="c.value" :value="c.value">
                {{ c.label }}
              </SelectItem>
            </SelectContent>
          </Select>

          <Select v-model="selectedProjectId" :disabled="!selectedClientId">
            <SelectTrigger class="w-[180px]">
              <SelectValue :placeholder="t('reports.filters.project')" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem v-for="p in filteredProjectOptions" :key="p.value" :value="p.value">
                {{ p.label }}
              </SelectItem>
            </SelectContent>
          </Select>

          <Button @click="loadReport">
            {{ t("reports.filters.apply") }}
          </Button>
        </CardContent>
      </Card>

      <div class="flex-1 flex flex-col min-h-0 relative">
        <div v-if="loading && !report" class="absolute inset-0 flex items-center justify-center z-10 bg-background/50">
          <Loader2 class="w-8 h-8 animate-spin text-muted-foreground" />
        </div>

        <Alert v-else-if="error" variant="destructive" class="mb-4">
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>{{ error }}</AlertDescription>
        </Alert>

        <template v-else-if="report">
          <div class="flex flex-col gap-4 flex-1 min-h-0">
            <Card v-if="report.chart.dates.length" class="shrink-0">
              <CardContent class="p-2">
                <v-chart :option="chartOption" autoresize class="report-chart" />
              </CardContent>
            </Card>

            <div v-if="report.rows.length === 0"
              class="flex flex-col items-center justify-center py-12 text-muted-foreground">
              <BoxSelect class="w-12 h-12 mb-4 opacity-50" />
              <p>{{ t('reports.empty') }}</p>
            </div>

            <!-- shadcn Table replacing NDataTable -->
            <div v-else class="flex-1 min-h-0 overflow-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead class="w-[110px]">{{ t("reports.table.date") }}</TableHead>
                    <TableHead class="w-[100px]">{{ t("reports.table.client") }}</TableHead>
                    <TableHead>{{ t("reports.table.project") }}</TableHead>
                    <TableHead class="w-[70px] text-right">{{ t("reports.table.hours") }}</TableHead>
                    <TableHead class="w-[90px] text-right">{{ t("reports.table.income") }}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="row in paginatedRows" :key="`${row.date}-${row.projectName}`">
                    <TableCell>{{ row.date }}</TableCell>
                    <TableCell>{{ row.clientName }}</TableCell>
                    <TableCell class="truncate max-w-[200px]" :title="row.projectName">{{ row.projectName }}</TableCell>
                    <TableCell class="text-right tabular-nums">{{ row.hours.toFixed(1) }}</TableCell>
                    <TableCell class="text-right tabular-nums">{{ row.income.toFixed(2) }}</TableCell>
                  </TableRow>
                </TableBody>
              </Table>

              <!-- Pagination -->
              <div v-if="totalPages > 1" class="flex items-center justify-between px-4 py-3 border-t">
                <span class="text-sm text-muted-foreground">
                  Page {{ currentPage }} of {{ totalPages }} ({{ report.rows.length }} rows)
                </span>
                <div class="flex gap-2">
                  <Button variant="outline" size="sm" :disabled="currentPage === 1" @click="currentPage--">
                    Previous
                  </Button>
                  <Button variant="outline" size="sm" :disabled="currentPage === totalPages" @click="currentPage++">
                    Next
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Fixed chart height to make calc() predictable */
.report-chart {
  height: 200px;
}
</style>
