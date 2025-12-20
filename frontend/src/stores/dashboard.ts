import { defineStore } from "pinia";
import { ref } from "vue";
import { useTimesheetStore } from "./timesheet";
import { useInvoiceStore } from "./invoices";
import type { RecentActivity } from "@/types";

export const useDashboardStore = defineStore("dashboard", () => {
  const loading = ref(false);

  // Dashboard Metrics
  const totalHoursWeek = ref(0);
  const totalRevenueMonth = ref(0);
  const pendingAmount = ref(0);

  // Recent Activity
  const recentActivities = ref<RecentActivity[]>([]);

  const timesheetStore = useTimesheetStore();
  const invoiceStore = useInvoiceStore();

  async function fetchDashboardData() {
    loading.value = true;
    try {
      // Parallel fetch from domain stores
      await Promise.all([
        timesheetStore.fetchTimesheet(),
        invoiceStore.fetchInvoices(),
      ]);

      // 1. Calculate Tasks (Mock logic: All entries = "This Week" for MVP demo)
      totalHoursWeek.value = parseFloat(timesheetStore.totalHours);

      // 2. Financials from Invoice Store
      pendingAmount.value = invoiceStore.stats.totalDue;

      // Mock Monthly Revenue (Paid Invoices) - In real app, filter by date
      const invoiceList = invoiceStore.invoices || [];
      totalRevenueMonth.value = invoiceList
        .filter((i) => i.status === "paid")
        .reduce((sum, i) => sum + i.total, 0);

      // 3. Recent Activity from Timesheet
      // Accessing internal Enriched Entries involves ensuring getter is reactive
      // We manually map top 5 entries for now
      recentActivities.value = timesheetStore.entries.slice(0, 5).map((e) => {
        // We know timesheetStore fetches projects internally for its own getters,
        // but exposed 'entries' are raw. We can use the getter if it's accessible or just map basic info.
        // For Dashboard, let's use the Raw data + strict types, or leverage the store's computed.
        // Since we are inside an action, we can access the store state.
        // Ideally we should use `timesheetStore.groupedByDay` or similar, but let's simple map:
        return {
          id: e.id,
          project: `Project #${e.projectId}`, // ideally lookup from ProjectStore, but keep simple for MVP
          date: e.date,
          hours: (e.durationSeconds / 3600).toFixed(1),
          description: e.description,
        };
      });
    } catch (error) {
      console.error("Failed to fetch dashboard data", error);
    } finally {
      loading.value = false;
    }
  }

  return {
    loading,
    totalHoursWeek,
    totalRevenueMonth,
    pendingAmount,
    recentActivities,
    fetchDashboardData,
  };
});
