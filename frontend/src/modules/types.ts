import type { Component } from "vue";
import type { RouteRecordRaw } from "vue-router";
import type { Locale } from "@/locales/types";

/**
 * Available module identifiers in the application
 */
export type ModuleID =
  | "dashboard"
  | "clients"
  | "projects"
  | "timesheet"
  | "invoices"
  | "reports"
  | "finance"
  | "settings";

export type ModuleMessages = Partial<Record<Locale, Record<string, unknown>>>;

/**
 * Navigation item for sidebar menu
 */
export type ModuleNavItem = {
  key: string;
  labelKey: string;
  icon?: Component;
  moduleID?: ModuleID;
  children?: ModuleNavItem[];
};

/**
 * Settings page contributed by a module
 */
export type ModuleSettingsPage = {
  key: string;
  labelKey: string;
  component: Component;
  order: number;
  moduleID: ModuleID;
};

/**
 * Route definition with static component
 */
export type ModuleRoute = {
  path: string;
  component: Component;
  meta?: {
    requiresAuth?: boolean;
    layout?: string;
    moduleID?: ModuleID;
  };
  children?: ModuleRoute[];
  redirect?: string;
};

/**
 * Application module definition
 */
export type AppModule = {
  id: ModuleID;
  enabledByDefault: boolean;
  toggleable: boolean;
  nav?: ModuleNavItem;
  routes: ModuleRoute[];
  settingsPages?: ModuleSettingsPage[];
  messages?: ModuleMessages;
};

/**
 * Convert ModuleRoute to vue-router's RouteRecordRaw
 */
export function toRouteRecordRaw(route: ModuleRoute): RouteRecordRaw {
  // Handle redirect routes
  if (route.redirect) {
    return {
      path: route.path,
      redirect: route.redirect,
    } as RouteRecordRaw;
  }

  // Handle component routes with children
  if (route.children && route.children.length > 0) {
    return {
      path: route.path,
      component: route.component,
      meta: route.meta,
      children: route.children.map(toRouteRecordRaw),
    } as RouteRecordRaw;
  }

  // Handle simple component routes
  return {
    path: route.path,
    component: route.component,
    meta: route.meta,
  } as RouteRecordRaw;
}
