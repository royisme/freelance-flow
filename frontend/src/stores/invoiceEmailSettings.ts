import { defineStore } from "pinia";
import { ref } from "vue";
import { api } from "@/api";
import type { InvoiceEmailSettings } from "@/types";

export const useInvoiceEmailSettingsStore = defineStore(
  "invoiceEmailSettings",
  () => {
    const settings = ref<InvoiceEmailSettings | null>(null);
    const loading = ref(false);
    const error = ref<string | null>(null);

    async function fetchSettings() {
      loading.value = true;
      error.value = null;
      try {
        settings.value = await api.invoiceEmailSettings.get();
      } catch (e) {
        error.value = e instanceof Error ? e.message : "Failed to load settings";
      } finally {
        loading.value = false;
      }
    }

    async function saveSettings(input: InvoiceEmailSettings) {
      loading.value = true;
      error.value = null;
      try {
        settings.value = await api.invoiceEmailSettings.update(input);
      } catch (e) {
        error.value = e instanceof Error ? e.message : "Failed to save settings";
        throw e;
      } finally {
        loading.value = false;
      }
    }

    return {
      settings,
      loading,
      error,
      fetchSettings,
      saveSettings,
    };
  }
);
