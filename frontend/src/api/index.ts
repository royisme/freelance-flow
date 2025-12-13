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
import { dateOnlySortKey } from "@/utils/date";
import type { InvoiceEmailSettings } from "@/types";
import type { ReportFilter, ReportOutput } from "@/types";
import type {
  IClientService,
  IProjectService,
  ITimeEntryService,
  IInvoiceService,
  StatusBarOutput,
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
  getDefaultMessage: (id) => {
    const wails = window as unknown as {
      go: {
        services: {
          InvoiceService: {
            GetDefaultMessage: (userId: number, invoiceId: number) => Promise<string>;
          };
        };
      };
    };
    return wails.go.services.InvoiceService.GetDefaultMessage(getUserId(), id);
  },
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

const wailsReportService = {
  get: (filter: ReportFilter) => {
    const wails = window as unknown as {
      go: {
        services: {
          ReportService: {
            Get: (userId: number, f: ReportFilter) => Promise<ReportOutput>;
          };
        };
      };
    };
    return wails.go.services.ReportService.Get(getUserId(), filter);
  },
};

const wailsInvoiceEmailSettingsService = {
  get: () => WailsInvoiceEmailSettingsService.Get(getUserId()),
  update: (input: InvoiceEmailSettings) =>
    WailsInvoiceEmailSettingsService.Update(getUserId(), input),
  export: () => WailsInvoiceEmailSettingsService.ExportSettings(getUserId()),
};

const wailsStatusBarService = {
  get: () => {
    const wails = window as unknown as {
      go: {
        services: {
          StatusBarService: {
            Get: (userId: number) => Promise<StatusBarOutput>;
          };
        };
      };
    };
    return wails.go.services.StatusBarService.Get(getUserId());
  },
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
      reports: wailsReportService,
      invoiceEmailSettings: wailsInvoiceEmailSettingsService,
      statusBar: wailsStatusBarService,
    }
  : {
      clients: mockClientService,
      projects: mockProjectService,
      timeEntries: mockTimeEntryService,
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
          defaultMessageTemplate: "Thank you for your business.",
        }),
        update: async (input: dto.UserSettings) => input,
      },
      reports: {
        get: async () => ({
          totalHours: 0,
          totalIncome: 0,
          rows: [],
          chart: { dates: [], revenue: [], hours: [] },
        }),
      },
      invoices: {
        ...mockInvoiceService,
        getDefaultMessage: async () => "Thank you for your business.",
      },
      invoiceEmailSettings: {
        get: async () => mockEmailSettings,
        update: async (input: InvoiceEmailSettings) => input,
        export: async () => JSON.stringify(mockEmailSettings),
      },
      statusBar: {
        get: async (): Promise<StatusBarOutput> => {
          const [entries, invoices, projects, settings] = await Promise.all([
            mockTimeEntryService.list(undefined),
            mockInvoiceService.list(),
            mockProjectService.list(),
            // Use mock settings to provide currency.
            // (Mock mode doesn't persist real settings.)
            (async () => ({
              currency: "USD",
            }))(),
          ]);

          const now = new Date();
          const monthStart = new Date(now.getFullYear(), now.getMonth(), 1);
          const nextMonthStart = new Date(now.getFullYear(), now.getMonth() + 1, 1);
          const monthStartKey =
            monthStart.getFullYear() * 10000 +
            (monthStart.getMonth() + 1) * 100 +
            monthStart.getDate();
          const nextMonthStartKey =
            nextMonthStart.getFullYear() * 10000 +
            (nextMonthStart.getMonth() + 1) * 100 +
            nextMonthStart.getDate();

          const monthSeconds = entries
            .filter((e) => {
              const key = dateOnlySortKey(e.date);
              return key >= monthStartKey && key < nextMonthStartKey;
            })
            .reduce((acc, e) => acc + (e.durationSeconds ?? 0), 0);

          const unpaidTotal = invoices
            .filter((inv) => inv.status === "sent" || inv.status === "overdue")
            .reduce((acc, inv) => acc + (inv.total ?? 0), 0);

          const projectRateById = new Map<number, number>(
            projects.map((p) => [p.id, p.hourlyRate ?? 0])
          );
          const uninvoicedTotal = entries
            .filter((e) => e.billable !== false && e.invoiced !== true)
            .reduce((acc, e) => {
              const rate = projectRateById.get(e.projectId) ?? 0;
              const hours = (e.durationSeconds ?? 0) / 3600;
              return acc + hours * rate;
            }, 0);

          return {
            monthSeconds,
            uninvoicedTotal,
            unpaidTotal,
            currency: settings.currency,
          };
        },
      },
    };

// Re-export types for convenience
export type {
  IClientService,
  IProjectService,
  ITimeEntryService,
  IInvoiceService,
};
