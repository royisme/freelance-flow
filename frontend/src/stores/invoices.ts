import { defineStore, storeToRefs } from "pinia";
import { ref, computed } from "vue";
import { api } from "@/api";
import type {
  Invoice,
  CreateInvoiceInput,
  UpdateInvoiceInput,
} from "@/types";
import { useClientStore } from "./clients";

// Enriched Invoice for View (adds Client Name directly)
export type EnrichedInvoice = Invoice & {
  clientName?: string;
  clientCurrency?: string;
};

export const useInvoiceStore = defineStore("invoices", () => {
  const invoices = ref<Invoice[]>([]);
  const loading = ref(false);
  const clientStore = useClientStore();
  const { clients } = storeToRefs(clientStore);

  async function fetchInvoices() {
    loading.value = true;
    try {
      // Ensure we have client data for joining
      if (clients.value.length === 0) {
        await clientStore.fetchClients();
      }
      invoices.value = await api.invoices.list();
    } catch (error) {
      console.error("Failed to fetch invoices", error);
    } finally {
      loading.value = false;
    }
  }

  async function createInvoice(input: CreateInvoiceInput) {
    loading.value = true;
    try {
      await api.invoices.create(input);
      await fetchInvoices();
    } finally {
      loading.value = false;
    }
  }

  async function updateInvoice(input: UpdateInvoiceInput) {
    loading.value = true;
    try {
      await api.invoices.update(input);
      await fetchInvoices();
    } finally {
      loading.value = false;
    }
  }

  async function deleteInvoice(id: number) {
    loading.value = true;
    try {
      await api.invoices.delete(id);
      await fetchInvoices();
    } finally {
      loading.value = false;
    }
  }

  // Enriched Getters
  const enrichedInvoices = computed<EnrichedInvoice[]>(() => {
    return invoices.value.map((inv) => {
      const client = clients.value.find((c) => c.id === inv.clientId);
      return {
        ...inv,
        clientName: client?.name || "Unknown Client",
        clientCurrency: client?.currency || "USD",
      };
    });
  });

  // Grouped by Status (for Dashboard widgets)
  const stats = computed(() => {
    const totalDue = invoices.value
      .filter((i) => i.status === "sent" || i.status === "overdue")
      .reduce((acc, curr) => acc + curr.total, 0);

    return {
      totalDue,
    };
  });

  async function setTimeEntries(invoiceId: number, timeEntryIds: number[]) {
    loading.value = true;
    try {
      await api.invoices.setTimeEntries({ invoiceId, timeEntryIds });
      await fetchInvoices();
    } finally {
      loading.value = false;
    }
  }

  async function generatePdf(id: number, message?: string) {
    return api.invoices.generatePdf(id, message);
  }

  async function sendEmail(id: number) {
    return api.invoices.sendEmail(id);
  }

  return {
    invoices,
    enrichedInvoices,
    stats,
    loading,
    fetchInvoices,
    createInvoice,
    updateInvoice,
    deleteInvoice,
    setTimeEntries,
    generatePdf,
    sendEmail,
  };
});
