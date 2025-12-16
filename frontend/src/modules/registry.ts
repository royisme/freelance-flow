import type {
  AppModule,
  ModuleID,
  ModuleMessages,
  ModuleNavItem,
  ModuleRoute,
  ModuleSettingsPage,
} from "@/modules/types";
import { financeModule } from "@/modules/finance/module";
import {
  LayoutDashboard,
  User,
  FolderKanban,
  Clock,
  FileText,
  BarChart3,
  Settings,
} from "lucide-vue-next";

// Static imports for all view components
import Dashboard from "@/views/Dashboard.vue";
import Clients from "@/views/Clients.vue";
import Projects from "@/views/Projects.vue";
import ProjectDetail from "@/views/ProjectDetail.vue";
import Timesheet from "@/views/Timesheet.vue";
import Invoices from "@/views/Invoices.vue";
import Reports from "@/views/Reports.vue";
import SettingsLayout from "@/views/settings/SettingsLayout.vue";
import GeneralSettings from "@/views/settings/GeneralSettings.vue";
import ProfileSettings from "@/views/settings/ProfileSettings.vue";
import InvoiceSettings from "@/views/settings/InvoiceSettings.vue";
import EmailSettings from "@/views/settings/EmailSettings.vue";

// ============================================================================
// Module Input Types
// ============================================================================

export type EnabledModulesInput = {
  moduleOverrides?: Partial<Record<ModuleID, boolean>> | null;
};

export function normalizeModuleOverrides(
  values: Record<string, unknown> | null | undefined
): Partial<Record<ModuleID, boolean>> | null {
  if (!values) return null;
  const out: Partial<Record<ModuleID, boolean>> = {};
  for (const [k, v] of Object.entries(values)) {
    if (typeof v !== "boolean") continue;
    out[k as ModuleID] = v;
  }
  return Object.keys(out).length > 0 ? out : null;
}

// ============================================================================
// Core Modules Definition
// ============================================================================

const baseModules: AppModule[] = [
  {
    id: "dashboard",
    enabledByDefault: true,
    toggleable: false,
    nav: {
      labelKey: "nav.dashboard",
      key: "dashboard",
      icon: LayoutDashboard,
    },
    routes: [
      {
        path: "/dashboard",
        component: Dashboard,
        meta: { requiresAuth: true, layout: "main", moduleID: "dashboard" },
      },
    ],
  },
  {
    id: "clients",
    enabledByDefault: true,
    toggleable: false,
    nav: { labelKey: "nav.clients", key: "clients", icon: User },
    routes: [
      {
        path: "/clients",
        component: Clients,
        meta: { requiresAuth: true, layout: "main", moduleID: "clients" },
      },
    ],
  },
  {
    id: "projects",
    enabledByDefault: true,
    toggleable: false,
    nav: { labelKey: "nav.projects", key: "projects", icon: FolderKanban },
    routes: [
      {
        path: "/projects",
        component: Projects,
        meta: { requiresAuth: true, layout: "main", moduleID: "projects" },
      },
      {
        path: "/projects/:id",
        component: ProjectDetail,
        meta: { requiresAuth: true, layout: "main", moduleID: "projects" },
      },
    ],
  },
  {
    id: "timesheet",
    enabledByDefault: true,
    toggleable: false,
    nav: {
      labelKey: "nav.timesheet",
      key: "timesheet",
      icon: Clock,
    },
    routes: [
      {
        path: "/timesheet",
        component: Timesheet,
        meta: { requiresAuth: true, layout: "main", moduleID: "timesheet" },
      },
    ],
  },
  {
    id: "invoices",
    enabledByDefault: true,
    toggleable: false,
    nav: { labelKey: "nav.invoices", key: "invoices", icon: FileText },
    routes: [
      {
        path: "/invoices",
        component: Invoices,
        meta: { requiresAuth: true, layout: "main", moduleID: "invoices" },
      },
    ],
  },
  {
    id: "reports",
    enabledByDefault: true,
    toggleable: false,
    nav: { labelKey: "nav.reports", key: "reports", icon: BarChart3 },
    routes: [
      {
        path: "/reports",
        component: Reports,
        meta: { requiresAuth: true, layout: "main", moduleID: "reports" },
      },
    ],
  },
];

// ============================================================================
// Settings Module (Special - Aggregates settings from other modules)
// ============================================================================

/**
 * Core settings pages that belong to the settings module itself
 */
const coreSettingsPages: ModuleSettingsPage[] = [
  {
    key: "general",
    order: 10,
    labelKey: "settings.general.title",
    component: GeneralSettings,
    moduleID: "settings",
  },
  {
    key: "profile",
    order: 20,
    labelKey: "settings.profile.title",
    component: ProfileSettings,
    moduleID: "settings",
  },
  {
    key: "invoice",
    order: 30,
    labelKey: "settings.invoice.title",
    component: InvoiceSettings,
    moduleID: "settings",
  },
  {
    key: "email",
    order: 40,
    labelKey: "settings.email.title",
    component: EmailSettings,
    moduleID: "settings",
  },
];

/**
 * Build the settings module by aggregating core pages and contributed pages from other modules
 */
function createSettingsModule(contribPages: ModuleSettingsPage[]): AppModule {
  // Merge core and contributed pages, then sort by order
  const allPages = [...coreSettingsPages, ...contribPages].sort(
    (a, b) => a.order - b.order
  );

  // Build child routes for settings
  const childRoutes: ModuleRoute[] = allPages.map((page) => ({
    path: page.key,
    component: page.component,
    meta: { moduleID: page.moduleID },
  }));

  // Build nav children
  const navChildren: ModuleNavItem[] = allPages.map((page) => ({
    key: `settings/${page.key}`,
    labelKey: page.labelKey,
    moduleID: page.moduleID,
  }));

  return {
    id: "settings",
    enabledByDefault: true,
    toggleable: false,
    nav: {
      labelKey: "nav.settings",
      key: "settings",
      icon: Settings,
      children: navChildren,
    },
    routes: [
      {
        path: "/settings",
        component: SettingsLayout,
        meta: { requiresAuth: true, layout: "main", moduleID: "settings" },
        children: [
          {
            path: "",
            redirect: "/settings/general",
            component: GeneralSettings,
          },
          ...childRoutes,
        ],
      },
    ],
  };
}

// ============================================================================
// Module Registry
// ============================================================================

const nonSettingsModules: AppModule[] = [...baseModules, financeModule];

export const allModules: AppModule[] = [
  ...nonSettingsModules,
  createSettingsModule(
    nonSettingsModules.flatMap((m) => m.settingsPages ?? [])
  ),
];

// ============================================================================
// Module Query Functions
// ============================================================================

export function isModuleEnabled(
  module: AppModule,
  input: EnabledModulesInput | null
): boolean {
  const overrides = input?.moduleOverrides ?? null;
  if (overrides && overrides[module.id] !== undefined) {
    return overrides[module.id] === true;
  }
  return module.enabledByDefault;
}

export function collectModuleMessages(modules: AppModule[]): ModuleMessages[] {
  return modules.flatMap((m) => (m.messages ? [m.messages] : []));
}

export function collectNavItems(modules: AppModule[]): ModuleNavItem[] {
  return modules.flatMap((m) => (m.nav ? [m.nav] : []));
}

export function getModuleByID(moduleID: ModuleID): AppModule | null {
  return allModules.find((m) => m.id === moduleID) ?? null;
}

export function isModuleIDEnabled(
  moduleID: ModuleID,
  overrides: EnabledModulesInput["moduleOverrides"] | null
): boolean {
  const mod = getModuleByID(moduleID);
  if (!mod) return true;
  return isModuleEnabled(mod, { moduleOverrides: overrides });
}
