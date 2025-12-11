// API Layer - Central service registry
// Automatically uses Wails bindings in production, mock services in dev

import {
  mockClientService,
  mockProjectService,
  mockTimeEntryService,
  mockInvoiceService,
} from "./mock";
import * as WailsClientService from "@/wailsjs/go/services/ClientService";
import * as WailsProjectService from "@/wailsjs/go/services/ProjectService";
import * as WailsTimesheetService from "@/wailsjs/go/services/TimesheetService";
import * as WailsInvoiceService from "@/wailsjs/go/services/InvoiceService";
import * as WailsSettingsService from "@/wailsjs/go/services/SettingsService";
import * as WailsInvoiceEmailSettingsService from "@/wailsjs/go/services/InvoiceEmailSettingsService";
import { useAuthStore } from "@/stores/auth";
import { dto } from "@/wailsjs/go/models";
import type { InvoiceEmailSettings } from "@/types";
import type {
  IClientService,
  IProjectService,
  ITimeEntryService,
  IInvoiceService,
} from "@/types";

// Check if we're in Wails runtime (window.go exists)
const isWailsRuntime = typeof window !== "undefined" && "go" in window;

const getUserId = () => {
  try {
    const store = useAuthStore();
    if (store.userId > 0) {
      return store.userId;
    }
  } catch {
    // Ignore store access errors (e.g., before pinia init)
  }
  const stored = localStorage.getItem("currentUserId");
  const parsed = stored ? parseInt(stored, 10) : 0;
  return Number.isFinite(parsed) ? parsed : 0;
};

// Wails service adapters - types are now aligned via re-exports
const wailsClientService: IClientService = {
  list: () => WailsClientService.List(getUserId()),
  get: (id) => WailsClientService.Get(getUserId(), id),
  create: (input) => WailsClientService.Create(getUserId(), input),
  update: (input) => WailsClientService.Update(getUserId(), input),
  delete: (id) => WailsClientService.Delete(getUserId(), id),
};

const wailsProjectService: IProjectService = {
  list: () => WailsProjectService.List(getUserId()),
  listByClient: (clientId) =>
    WailsProjectService.ListByClient(getUserId(), clientId),
  get: (id) => WailsProjectService.Get(getUserId(), id),
  create: (input) => WailsProjectService.Create(getUserId(), input),
  update: (input) => WailsProjectService.Update(getUserId(), input),
  delete: (id) => WailsProjectService.Delete(getUserId(), id),
};

const wailsTimeEntryService: ITimeEntryService = {
  list: (projectId) =>
    WailsTimesheetService.List(getUserId(), projectId ?? 0),
  get: (id) => WailsTimesheetService.Get(getUserId(), id),
  create: (input) => WailsTimesheetService.Create(getUserId(), input),
  update: (input) => WailsTimesheetService.Update(getUserId(), input),
  delete: (id) => WailsTimesheetService.Delete(getUserId(), id),
};

const wailsInvoiceService: IInvoiceService = {
  list: () => WailsInvoiceService.List(getUserId()),
  get: (id) => WailsInvoiceService.Get(getUserId(), id),
  create: (input) => WailsInvoiceService.Create(getUserId(), input),
  update: (input) => WailsInvoiceService.Update(getUserId(), input),
  delete: (id) => WailsInvoiceService.Delete(getUserId(), id),
  generatePdf: (id, message) =>
    WailsInvoiceService.GeneratePDF(getUserId(), id, message ?? ""),
  sendEmail: (id) => WailsInvoiceService.SendEmail(getUserId(), id),
  setTimeEntries: (input) =>
    WailsInvoiceService.SetTimeEntries(getUserId(), input),
};

const wailsSettingsService = {
  get: () => WailsSettingsService.Get(getUserId()),
  update: (input: dto.UserSettings) =>
    WailsSettingsService.Update(getUserId(), input),
};

const wailsInvoiceEmailSettingsService = {
  get: () => WailsInvoiceEmailSettingsService.Get(getUserId()),
  update: (input: InvoiceEmailSettings) =>
    WailsInvoiceEmailSettingsService.Update(getUserId(), input),
  export: () => WailsInvoiceEmailSettingsService.ExportSettings(getUserId()),
};

const mockEmailSettings: InvoiceEmailSettings = {
  provider: "mailto",
  subjectTemplate: "Invoice {{number}}",
  bodyTemplate: "Please find attached invoice {{number}}.",
  signature: "",
};

// Export the appropriate service based on runtime environment
export const api = isWailsRuntime
  ? {
      clients: wailsClientService,
      projects: wailsProjectService,
      timeEntries: wailsTimeEntryService,
      invoices: wailsInvoiceService,
      settings: wailsSettingsService,
      invoiceEmailSettings: wailsInvoiceEmailSettingsService,
    }
  : {
      clients: mockClientService,
      projects: mockProjectService,
      timeEntries: mockTimeEntryService,
      invoices: mockInvoiceService,
      settings: {
        get: async () => ({
          currency: "USD",
          defaultTaxRate: 0,
          language: "en-US",
          theme: "light",
          dateFormat: "2006-01-02",
          timezone: "UTC",
          senderName: "",
          senderCompany: "",
          senderAddress: "",
          senderPhone: "",
          senderEmail: "",
          senderPostalCode: "",
          invoiceTerms: "Due upon receipt",
          defaultMessageTemplate: "",
        }),
        update: async (input: dto.UserSettings) => input,
      },
      invoiceEmailSettings: {
        get: async () => mockEmailSettings,
        update: async (input: InvoiceEmailSettings) => input,
        export: async () => JSON.stringify(mockEmailSettings),
      },
    };

// Re-export types for convenience
export type {
  IClientService,
  IProjectService,
  ITimeEntryService,
  IInvoiceService,
};
