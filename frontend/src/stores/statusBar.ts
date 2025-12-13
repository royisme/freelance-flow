import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { api } from "@/api";
import type { StatusBarOutput } from "@/types";

function currencySymbol(currency: string): string {
  switch (currency) {
    case "USD":
    case "CAD":
      return "$";
    case "CNY":
      return "¥";
    case "EUR":
      return "€";
    case "GBP":
      return "£";
    default:
      return "";
  }
}

function formatHoursFromSeconds(seconds: number): string {
  const hours = seconds / 3600;
  if (!Number.isFinite(hours)) return "0h";
  const isInt = Math.abs(hours - Math.round(hours)) < 1e-9;
  return `${isInt ? Math.round(hours) : hours.toFixed(1)}h`;
}

function formatMoney(amount: number, currency: string): string {
  const symbol = currencySymbol(currency);
  const safe = Number.isFinite(amount) ? amount : 0;
  const formatted = safe.toLocaleString(undefined, {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2,
  });
  return symbol ? `${symbol}${formatted}` : `${formatted} ${currency}`;
}

export const useStatusBarStore = defineStore("statusBar", () => {
  const loading = ref(false);
  const data = ref<StatusBarOutput | null>(null);

  const monthHoursLabel = computed(() =>
    formatHoursFromSeconds(data.value?.monthSeconds ?? 0)
  );
  const unpaidTotalLabel = computed(() => {
    const currency = data.value?.currency ?? "USD";
    return formatMoney(data.value?.unpaidTotal ?? 0, currency);
  });
  const uninvoicedTotalLabel = computed(() => {
    const currency = data.value?.currency ?? "USD";
    return formatMoney(data.value?.uninvoicedTotal ?? 0, currency);
  });

  async function refresh() {
    loading.value = true;
    try {
      data.value = await api.statusBar.get();
    } finally {
      loading.value = false;
    }
  }

  return {
    loading,
    data,
    monthHoursLabel,
    uninvoicedTotalLabel,
    unpaidTotalLabel,
    refresh,
  };
});
