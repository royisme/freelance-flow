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
import * as WailsUserPreferencesService from "@/wailsjs/go/services/UserPreferencesService";
import * as WailsUserTaxSettingsService from "@/wailsjs/go/services/UserTaxSettingsService";
import * as WailsUserInvoiceSettingsService from "@/wailsjs/go/services/UserInvoiceSettingsService";
// Note: FinanceService binding will be available at runtime.
// If typescript complains, we can declare it or cast window.go.
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
import type { FinanceAccount, FinanceSummary } from "@/types/finance";

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
  get: (id) => WailsClientService.Get(getUserId(), Number(id)),
  create: (input) => WailsClientService.Create(getUserId(), input),
  update: (input) => WailsClientService.Update(getUserId(), input),
  delete: (id) => WailsClientService.Delete(getUserId(), Number(id)),
};

const wailsProjectService: IProjectService = {
  list: () => WailsProjectService.List(getUserId()),
  listByClient: (clientId) =>
    WailsProjectService.ListByClient(getUserId(), Number(clientId)),
  get: (id) => WailsProjectService.Get(getUserId(), Number(id)),
  create: (input) => WailsProjectService.Create(getUserId(), input),
  update: (input) => WailsProjectService.Update(getUserId(), input),
  delete: (id) => WailsProjectService.Delete(getUserId(), Number(id)),
};

const wailsTimeEntryService: ITimeEntryService = {
  list: (projectId) =>
    WailsTimesheetService.List(getUserId(), projectId ? Number(projectId) : 0),
  get: (id) => WailsTimesheetService.Get(getUserId(), Number(id)),
  create: (input) => WailsTimesheetService.Create(getUserId(), input),
  update: (input) => WailsTimesheetService.Update(getUserId(), input),
  delete: (id) => WailsTimesheetService.Delete(getUserId(), Number(id)),
};

const wailsInvoiceService: IInvoiceService = {
  list: () => WailsInvoiceService.List(getUserId()),
  get: (id) => WailsInvoiceService.Get(getUserId(), Number(id)),
  create: (input) => WailsInvoiceService.Create(getUserId(), input),
  update: (input) => WailsInvoiceService.Update(getUserId(), input),
  delete: (id) => WailsInvoiceService.Delete(getUserId(), Number(id)),
  getDefaultMessage: (id) => {
    const wails = window as unknown as {
      go: {
        services: {
          InvoiceService: {
            GetDefaultMessage: (
              userId: number,
              invoiceId: number
            ) => Promise<string>;
          };
        };
      };
    };
    return wails.go.services.InvoiceService.GetDefaultMessage(
      getUserId(),
      Number(id)
    );
  },
  generatePdf: (id, message) =>
    WailsInvoiceService.GeneratePDF(getUserId(), Number(id), message ?? ""),
  sendEmail: (id) => WailsInvoiceService.SendEmail(getUserId(), Number(id)),
  setTimeEntries: (input) =>
    WailsInvoiceService.SetTimeEntries(getUserId(), input),
  updateStatus: (id, status) => {
    // Cast to any to bypass missing type definition until regeneration
    return (WailsInvoiceService as any).UpdateStatus(
      getUserId(),
      Number(id),
      status
    );
  },
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
  update: (input: InvoiceEmailSettings) => {
    const dtoSettings = new dto.InvoiceEmailSettings({
      provider: input.provider,
      from: input.from || "",
      replyTo: input.replyTo,
      subjectTemplate: input.subjectTemplate,
      bodyTemplate: input.bodyTemplate,
      resendApiKey: input.resendApiKey,
      smtpHost: input.smtpHost,
      smtpPort: input.smtpPort,
      smtpUsername: input.smtpUsername,
      smtpPassword: input.smtpPassword,
      smtpUseTls: input.smtpUseTls,
      signature: input.signature,
    });
    return WailsInvoiceEmailSettingsService.Update(getUserId(), dtoSettings);
  },
  export: () => WailsInvoiceEmailSettingsService.ExportSettings(getUserId()),
};

const wailsUserPreferencesService = {
  get: () => WailsUserPreferencesService.Get(getUserId()),
  update: (input: dto.UserPreferences) =>
    WailsUserPreferencesService.Update(getUserId(), input),
};

const wailsUserTaxSettingsService = {
  get: () => WailsUserTaxSettingsService.Get(getUserId()),
  update: (input: dto.UserTaxSettings) =>
    WailsUserTaxSettingsService.Update(getUserId(), input),
};

const wailsUserInvoiceSettingsService = {
  get: () => WailsUserInvoiceSettingsService.Get(getUserId()),
  update: (input: dto.UserInvoiceSettings) =>
    WailsUserInvoiceSettingsService.Update(getUserId(), input),
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

// --- Finance Service Adapter ---
const wailsFinanceService = {
  summary: {
    get: async (): Promise<FinanceSummary> => {
       const wails = window as any;
       return wails.go.services.FinanceService.GetSummary(getUserId());
    }
  },
  accounts: {
    list: async (): Promise<FinanceAccount[]> => {
      const wails = window as any;
      return wails.go.services.FinanceService.GetAccounts(getUserId());
    },
    create: async (input: any) => {
      const wails = window as any;
      return wails.go.services.FinanceService.CreateAccount(getUserId(), input);
    },
    update: async (input: any) => {
      const wails = window as any;
      return wails.go.services.FinanceService.UpdateAccount(getUserId(), input);
    },
    delete: async (id: number) => {
      const wails = window as any;
      return wails.go.services.FinanceService.DeleteAccount(getUserId(), id);
    }
  },
  categories: {
    list: async () => {
      const wails = window as any;
      return wails.go.services.FinanceService.GetCategories(getUserId());
    },
    create: async (input: any) => {
      const wails = window as any;
      return wails.go.services.FinanceService.CreateCategory(getUserId(), input);
    },
    update: async (input: any) => {
      const wails = window as any;
      return wails.go.services.FinanceService.UpdateCategory(getUserId(), input);
    },
    delete: async (id: number) => {
      const wails = window as any;
      return wails.go.services.FinanceService.DeleteCategory(getUserId(), id);
    }
  },
  transactions: {
    list: async (filter: any) => {
      const wails = window as any;
      return wails.go.services.FinanceService.GetTransactions(getUserId(), filter);
    },
    update: async (id: number, categoryId: number | null) => {
      const wails = window as any;
      return wails.go.services.FinanceService.UpdateTransaction(getUserId(), id, categoryId);
    },
    delete: async (id: number) => {
      const wails = window as any;
      return wails.go.services.FinanceService.DeleteTransaction(getUserId(), id);
    },
    import: async (input: any) => {
      const wails = window as any;
      return wails.go.services.FinanceService.ImportTransactions(getUserId(), input);
    }
  }
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
      userPreferences: wailsUserPreferencesService,
      userTaxSettings: wailsUserTaxSettingsService,
      userInvoiceSettings: wailsUserInvoiceSettingsService,
      reports: wailsReportService,
      invoiceEmailSettings: wailsInvoiceEmailSettingsService,
      statusBar: wailsStatusBarService,
      finance: wailsFinanceService,
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
        updateStatus: async () => {},
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
          const nextMonthStart = new Date(
            now.getFullYear(),
            now.getMonth() + 1,
            1
          );
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
      userPreferences: {
        get: async () =>
          new dto.UserPreferences({
            currency: "USD",
            language: "en-US",
            theme: "light",
            dateFormat: "2006-01-02",
            timezone: "UTC",
          }),
        update: async (input: dto.UserPreferences) => input,
      },
      userTaxSettings: {
        get: async () =>
          new dto.UserTaxSettings({
            hstRegistered: false,
            taxEnabled: false,
            defaultTaxRate: 0,
          }),
        update: async (input: dto.UserTaxSettings) => input,
      },
      userInvoiceSettings: {
        get: async () =>
          new dto.UserInvoiceSettings({
            defaultTerms: "Due upon receipt",
            defaultMessageTemplate: "Thank you for your business.",
          }),
        update: async (input: dto.UserInvoiceSettings) => input,
      },
      finance: {
        summary: {
          get: async () => ({
            totalBalance: 0,
            totalIncome: 0,
            totalExpense: 0,
            cashFlow: 0,
          }),
        },
        accounts: {
          list: async () => [],
          create: async (input: any) => ({ ...input, id: 1 }),
          update: async (input: any) => input,
          delete: async (id: number) => {},
        },
        categories: {
          list: async () => [],
          create: async (input: any) => input,
          update: async (input: any) => input,
          delete: async (id: number) => {},
        },
        transactions: {
          list: async (filter: any) => [],
          update: async (id: number, catId: number) => {},
          delete: async (id: number) => {},
          import: async (input: any) => 0,
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
