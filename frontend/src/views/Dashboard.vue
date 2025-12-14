<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { NCard, NStatistic, NGrid, NGridItem, NList, NListItem, NThing, NTag, NIcon } from 'naive-ui'
import { useDashboardStore } from '@/stores/dashboard'
import { useAuthStore } from '@/stores/auth'
import PageContainer from '@/components/PageContainer.vue'
import PageHeader from '@/components/PageHeader.vue'
import { ClockCircleOutlined, DollarOutlined, RiseOutlined } from '@vicons/antd'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const store = useDashboardStore()
const authStore = useAuthStore()

const username = computed(() => authStore.currentUser?.username || 'User')

onMounted(() => {
  store.fetchDashboardData()
})
</script>

<template>
  <PageContainer fill>
    <PageHeader :title="t('dashboard.greeting', { name: username })" :subtitle="t('dashboard.weekOverview')" />

    <!-- Key Metrics Cards -->
    <div class="metrics-section">
      <n-grid x-gap="24" y-gap="24" :cols="3">
        <n-grid-item>
          <n-card :bordered="true" class="metric-card">
            <n-statistic :label="t('dashboard.metrics.weeklyHours')">
              <template #prefix>
                <div class="icon-box orange">
                  <n-icon>
                    <ClockCircleOutlined />
                  </n-icon>
                </div>
              </template>
              <template #default>
                <span class="metric-value">{{ store.totalHoursWeek.toFixed(1) }}</span>
              </template>
              <template #suffix><span class="metric-unit">{{ t('dashboard.metrics.hoursUnit') }}</span></template>
            </n-statistic>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card :bordered="true" class="metric-card">
            <n-statistic :label="t('dashboard.metrics.monthlyRevenue')">
              <template #prefix>
                <div class="icon-box green">
                  <n-icon>
                    <DollarOutlined />
                  </n-icon>
                </div>
              </template>
              <template #default>
                <span class="metric-value">{{ store.totalRevenueMonth.toLocaleString() }}</span>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>

        <n-grid-item>
          <n-card :bordered="true" class="metric-card">
            <n-statistic :label="t('dashboard.metrics.pendingAmount')">
              <template #prefix>
                <div class="icon-box rose">
                  <n-icon>
                    <RiseOutlined />
                  </n-icon>
                </div>
              </template>
              <template #default>
                <span class="metric-value">{{ store.pendingAmount.toLocaleString() }}</span>
              </template>
            </n-statistic>
          </n-card>
        </n-grid-item>
      </n-grid>
    </div>

    <!-- Recent Activity Section -->
    <div class="section-container">
      <div class="section-header">
        <h2 class="section-title">{{ t('dashboard.recentActivity.title') }}</h2>

      </div>

      <n-card :bordered="true" class="activity-card-container"
        :content-style="{ padding: 'var(--space-2)', height: '100%', display: 'flex', flexDirection: 'column' }">
        <div class="activity-scroll-container">
          <n-list hoverable clickable>
            <n-list-item v-for="activity in store.recentActivities" :key="activity.id">
              <template #prefix>
                <div class="activity-icon-bg">
                  <n-icon size="18" color="#EA580C">
                    <ClockCircleOutlined />
                  </n-icon>
                </div>
              </template>
              <n-thing :title="activity.project" :description="`${activity.date} · ${activity.description}`" />
              <template #suffix>
                <span class="hours-badge">{{ t('dashboard.recentActivity.hoursLabel', { hours: activity.hours })
                }}</span>
              </template>
            </n-list-item>

            <div v-if="store.recentActivities.length === 0" class="empty-state">
              {{ t('dashboard.recentActivity.empty') }}
            </div>
          </n-list>
        </div>
      </n-card>
    </div>
  </PageContainer>
</template>

<style scoped>
/* Metrics Section - Fixed height with spacing for hover animation */
.metrics-section {
  flex-shrink: 0;
  padding-top: 8px;
  margin-bottom: 24px;
}

/* Metric Cards */
.metric-card {
  transition: all var(--transition-normal);
  overflow: visible;
}

.metric-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  border-color: var(--color-warm-orange);
  z-index: 1;
}

/* Metric Typography */
.metric-value {
  font-family: var(--font-sans);
  font-weight: 700;
  color: var(--n-text-color);
  letter-spacing: -0.03em;
}

.metric-unit {
  font-size: var(--text-base);
  font-weight: 500;
  color: var(--n-text-color-3);
  margin-left: var(--space-1);
}

/* Activity Section - Fills remaining space */
.section-container {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-4);
}

.section-title {
  font-size: var(--text-xl);
  font-weight: 700;
  color: var(--n-text-color-2);
  margin: 0;
}

.clickable-tag {
  cursor: pointer;
  background-color: var(--n-close-color-hover);
  color: var(--n-text-color-3);
}

.clickable-tag:hover {
  background-color: var(--n-close-color-pressed);
  color: var(--n-text-color-2);
}

/* Recent Activity List - prefix icon */
.activity-icon-bg {
  width: 36px;
  height: 36px;
  flex-shrink: 0;
  border-radius: var(--radius-full);
  background-color: var(--color-warning-bg);
  display: flex;
  align-items: center;
  justify-content: center;
}

.hours-badge {
  font-family: var(--font-sans);
  font-weight: 600;
  font-size: var(--text-sm);
  color: var(--n-text-color-2);
  background-color: var(--n-action-color);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-full);
  white-space: nowrap;
}

.empty-state {
  text-align: center;
  padding: var(--space-10);
  color: var(--n-text-color-3);
}

/* Activity Card Container - Flex fill */
.activity-card-container {
  flex: 1;
  min-height: 0;
}

/* Activity Scroll Container - Flex fill with scroll */
.activity-scroll-container {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

/* Custom scrollbar styling */
.activity-scroll-container::-webkit-scrollbar {
  width: 6px;
}

.activity-scroll-container::-webkit-scrollbar-track {
  background: transparent;
}

.activity-scroll-container::-webkit-scrollbar-thumb {
  background: var(--n-scrollbar-color);
  border-radius: 3px;
}

.activity-scroll-container::-webkit-scrollbar-thumb:hover {
  background: var(--n-scrollbar-color-hover);
}
</style>
