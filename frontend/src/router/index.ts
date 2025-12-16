import { createRouter, createWebHashHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { allModules } from "@/modules/registry";
import { toRouteRecordRaw } from "@/modules/types";
import type { ModuleID } from "@/modules/types";
import { useSettingsStore } from "@/stores/settings";
import {
  isModuleIDEnabled,
  normalizeModuleOverrides,
} from "@/modules/registry";

// Static imports for auth views
import Splash from "@/views/Splash.vue";
import Login from "@/views/Login.vue";
import Register from "@/views/Register.vue";

const routes = [
  // Auth routes (no auth required)
  {
    path: "/",
    redirect: "/splash",
  },
  {
    path: "/splash",
    component: Splash,
    meta: { requiresAuth: false, layout: "auth" },
  },
  {
    path: "/login",
    component: Login,
    meta: { requiresAuth: false, layout: "auth" },
  },
  {
    path: "/register",
    component: Register,
    meta: { requiresAuth: false, layout: "auth" },
  },

  // Main app routes (auth required) are provided by module registry
  ...allModules.flatMap((m) => m.routes.map(toRouteRecordRaw)),
];

const router = createRouter({
  // Use Hash history for Wails/Desktop compatibility
  history: createWebHashHistory(),
  routes,
});

// Persist last authenticated route for session restore.
router.afterEach((to) => {
  if (to.meta.requiresAuth) {
    localStorage.setItem("lastAuthedRoute", to.fullPath);
  }
});

// Navigation guard
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();
  const settingsStore = useSettingsStore();

  // If route requires auth and user is not authenticated
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // Redirect to splash if not initialized, otherwise to login
    if (!authStore.isInitialized) {
      next("/splash");
    } else if (authStore.usersList.length > 0) {
      next("/login");
    } else {
      next("/register");
    }
  }
  // If authenticated user tries to access auth pages, redirect to dashboard
  else if (
    !to.meta.requiresAuth &&
    authStore.isAuthenticated &&
    to.path !== "/splash"
  ) {
    next("/dashboard");
  } else {
    const moduleID = to.matched
      .map((r) => r.meta?.moduleID as ModuleID | undefined)
      .find((v) => v !== undefined);

    if (to.meta.requiresAuth && moduleID) {
      const overrides = normalizeModuleOverrides(
        settingsStore.settings?.moduleOverrides
      );
      if (!isModuleIDEnabled(moduleID, overrides)) {
        next("/dashboard");
        return;
      }
    }
    next();
  }
});

export default router;
