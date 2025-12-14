<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import {
  NAlert,
  NButton,
  NCard,
  NDataTable,
  NDatePicker,
  NEmpty,
  NSelect,
  NSpace,
  NSpin,
  NStatistic,
  type DataTableColumns,
} from "naive-ui";
import PageContainer from "@/components/PageContainer.vue";
import PageHeader from "@/components/PageHeader.vue";
import { api } from "@/api";
import type { Client, Project, ReportFilter, ReportOutput, ReportRow } from "@/types";
import VChart from "vue-echarts";
import type { EChartsOption } from "echarts";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

const loading = ref(false);
const error = ref<string | null>(null);
const report = ref<ReportOutput | null>(null);

const clients = ref<Client[]>([]);
const projects = ref<Project[]>([]);

const dateRange = ref<[number, number] | null>(null);
const selectedClientId = ref<number | null>(null);
const selectedProjectId = ref<number | null>(null);

const clientOptions = computed(() =>
  clients.value.map((c) => ({ label: c.name, value: c.id }))
);
const filteredProjectOptions = computed(() => {
  const filtered = selectedClientId.value
    ? projects.value.filter((p) => p.clientId === selectedClientId.value)
    : projects.value;
  return filtered.map((p) => ({ label: p.name, value: p.id }));
});

function formatDate(ts: number): string {
  const d = new Date(ts);
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
    if (dateRange.value) {
      filter.startDate = formatDate(dateRange.value[0]);
      filter.endDate = formatDate(dateRange.value[1]);
    }
    if (selectedClientId.value) filter.clientId = selectedClientId.value;
    if (selectedProjectId.value) filter.projectId = selectedProjectId.value;

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

const columns = computed<DataTableColumns<ReportRow>>(() => [
  { title: t("reports.table.date"), key: "date", width: 110 },
  { title: t("reports.table.client"), key: "clientName", width: 100 },
  { title: t("reports.table.project"), key: "projectName", ellipsis: { tooltip: true } },
  {
    title: t("reports.table.hours"),
    key: "hours",
    width: 70,
    align: "right",
    render: (row) => row.hours.toFixed(1),
  },
  {
    title: t("reports.table.income"),
    key: "income",
    width: 90,
    align: "right",
    render: (row) => row.income.toFixed(2),
  },
]);

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

// DataTable pagination config - show 5 items per page
const tablePagination = {
  pageSize: 5,
};
</script>

<template>
  <PageContainer fill>
    <PageHeader :title="t('reports.title')" :subtitle="t('reports.subtitle')">
      <template #extra>
        <n-space size="large">
          <n-statistic :label="t('reports.stats.totalHours')">
            <template #default>
              {{ report?.totalHours || 0 }}
            </template>
          </n-statistic>
          <n-statistic :label="t('reports.stats.totalIncome')">
            <template #default>
              {{ report?.totalIncome || 0 }}
            </template>
          </n-statistic>
        </n-space>
      </template>
    </PageHeader>

    <div class="reports-root">
      <n-card class="filters-card" size="small" :content-style="{ padding: '8px 12px' }">
        <n-space align="center" :wrap="true" :size="8">
          <n-date-picker v-model:value="dateRange" type="daterange" clearable size="small"
            :placeholder="t('reports.filters.dateRange')" class="filter-date" />
          <n-select v-model:value="selectedClientId" :options="clientOptions" clearable size="small"
            :placeholder="t('reports.filters.client')" class="filter-select" />
          <n-select v-model:value="selectedProjectId" :options="filteredProjectOptions" clearable size="small"
            :placeholder="t('reports.filters.project')" class="filter-select" />
          <n-button type="primary" size="small" @click="loadReport">
            {{ t("reports.filters.apply") }}
          </n-button>
        </n-space>
      </n-card>

      <div class="reports-body">
        <div v-if="loading && !report" class="loading-center">
          <n-spin size="large" />
        </div>

        <n-alert v-else-if="error" type="error" :title="error" />

        <template v-else-if="report">
          <div class="reports-content">
            <n-card v-if="report.chart.dates.length" size="small" :content-style="{ padding: '8px' }">
              <v-chart :option="chartOption" autoresize class="report-chart" />
            </n-card>

            <n-empty v-if="report.rows.length === 0" :description="t('reports.empty')" size="small" />

            <n-card v-else class="table-card" :bordered="true" size="small"
              :content-style="{ padding: 0, height: '100%', display: 'flex', flexDirection: 'column' }">
              <div class="table-scroll-container">
                <n-data-table :columns="columns" :data="report.rows" :bordered="false" :loading="loading"
                  size="small" />
              </div>
            </n-card>
          </div>
        </template>
      </div>
    </div>
  </PageContainer>
</template>

<style scoped>
.reports-root {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-height: 0;
}

.reports-body {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  position: relative;
}

.loading-center {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.reports-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 0;
  flex: 1;
}

.report-chart {
  height: clamp(180px, 32vh, 320px);
}

.table-card {
  flex: 1;
  min-height: 0;
}

.table-scroll-container {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

/* Custom scrollbar styling */
.table-scroll-container::-webkit-scrollbar {
  width: 6px;
}

.table-scroll-container::-webkit-scrollbar-track {
  background: transparent;
}

.table-scroll-container::-webkit-scrollbar-thumb {
  background: var(--n-scrollbar-color);
  border-radius: 3px;
}

.table-scroll-container::-webkit-scrollbar-thumb:hover {
  background: var(--n-scrollbar-color-hover);
}

.filter-date {
  width: 200px;
}

.filter-select {
  width: 140px;
}

@media (max-width: 520px) {

  .filter-date,
  .filter-select {
    width: 100%;
  }
}
</style>
