import { createRouter, createWebHashHistory } from "vue-router";
import { useAuthStore } from "@/stores/auth";

// Auth views (no lazy loading for faster initial load)
import Splash from "@/views/Splash.vue";

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
    component: () => import("@/views/Login.vue"),
    meta: { requiresAuth: false, layout: "auth" },
  },
  {
    path: "/register",
    component: () => import("@/views/Register.vue"),
    meta: { requiresAuth: false, layout: "auth" },
  },

  // Main app routes (auth required)
  {
    path: "/dashboard",
    component: () => import("@/views/Dashboard.vue"),
    meta: { requiresAuth: true, layout: "main" },
  },
  {
    path: "/clients",
    component: () => import("@/views/Clients.vue"),
    meta: { requiresAuth: true, layout: "main" },
  },
  {
    path: "/projects",
    component: () => import("@/views/Projects.vue"),
    meta: { requiresAuth: true, layout: "main" },
  },
  {
    path: "/timesheet",
    component: () => import("@/views/Timesheet.vue"),
    meta: { requiresAuth: true, layout: "main" },
  },
  {
    path: "/invoices",
    component: () => import("@/views/Invoices.vue"),
    meta: { requiresAuth: true, layout: "main" },
  },
  {
    path: "/reports",
    component: () => import("@/views/Reports.vue"),
    meta: { requiresAuth: true, layout: "main" },
  },
  {
    path: "/settings/invoice",
    component: () => import("@/views/settings/InvoiceSettings.vue"),
    meta: { requiresAuth: true, layout: "main" },
  },
];

const router = createRouter({
  // Use Hash history for Wails/Desktop compatibility
  history: createWebHashHistory(),
  routes,
});

// Navigation guard
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();

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
    next();
  }
});

export default router;
