<script setup lang="ts">
import { onMounted } from 'vue'
import { NCard, NStatistic, NGrid, NGridItem, NList, NListItem, NThing, NTag, NIcon, NSpace } from 'naive-ui'
import { useDashboardStore } from '@/stores/dashboard'
import PageContainer from '@/components/PageContainer.vue'
import { ClockCircleOutlined, DollarOutlined, RiseOutlined } from '@vicons/antd'

const store = useDashboardStore()

onMounted(() => {
  store.fetchDashboardData()
})
</script>

<template>
  <PageContainer title="早安，Roy 👋" subtitle="这里是你本周的工作概览">
    <!-- Key Metrics Cards -->
    <n-grid x-gap="24" y-gap="24" :cols="3">
      <n-grid-item>
        <n-card :bordered="true" class="metric-card">
          <n-statistic label="本周工时">
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
            <template #suffix><span class="metric-unit">小时</span></template>
          </n-statistic>
        </n-card>
      </n-grid-item>

      <n-grid-item>
        <n-card :bordered="true" class="metric-card">
          <n-statistic label="本月预计收入">
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
          <n-statistic label="待收金额">
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

    <!-- Recent Activity Section -->
    <div class="section-container">
      <div class="section-header">
        <h2 class="section-title">最近活动</h2>
        <n-tag :bordered="false" type="default" size="small" class="clickable-tag">View All</n-tag>
      </div>

      <n-card :bordered="true" class="activity-card-container" content-style="padding: 0;">
        <n-list hoverable clickable>
          <n-list-item v-for="activity in store.recentActivities" :key="activity.id">
            <div class="activity-item">
              <div class="activity-left">
                <div class="activity-icon-bg">
                  <n-icon size="18" color="#EA580C">
                    <ClockCircleOutlined />
                  </n-icon>
                </div>
                <div class="activity-content">
                  <div class="activity-title">{{ activity.project }}</div>
                  <div class="activity-desc">{{ activity.date }} · {{ activity.description }}</div>
                </div>
              </div>
              <div class="activity-right">
                <span class="hours-badge">{{ activity.hours }}h</span>
              </div>
            </div>
          </n-list-item>

          <div v-if="store.recentActivities.length === 0" class="empty-state">
            暂无活动记录
          </div>
        </n-list>
      </n-card>
    </div>
  </PageContainer>
</template>

<style scoped>
/* Metric Cards */
.metric-card {
  transition: all var(--transition-normal);
}

.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--color-warm-orange);
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

/* Activity Section */
.section-container {
  margin-top: var(--space-10);
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

/* Recent Activity List */
.activity-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-4) var(--space-6);
}

.activity-left {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.activity-icon-bg {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  background-color: var(--color-warning-bg);
  display: flex;
  align-items: center;
  justify-content: center;
}

.activity-title {
  font-weight: 600;
  color: var(--n-text-color);
  font-size: var(--text-base);
}

.activity-desc {
  color: var(--n-text-color-3);
  font-size: var(--text-sm);
  margin-top: 2px;
}

.hours-badge {
  font-family: var(--font-sans);
  font-weight: 600;
  font-size: var(--text-sm);
  color: var(--n-text-color-2);
  background-color: var(--n-action-color);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-full);
}

.empty-state {
  text-align: center;
  padding: var(--space-10);
  color: var(--n-text-color-3);
}
</style>
