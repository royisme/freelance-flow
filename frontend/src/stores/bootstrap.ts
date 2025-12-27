import { defineStore } from "pinia";
import { ref } from "vue";
// @ts-ignore - Wails runtime
import { EventsOn } from "@/wailsjs/runtime/runtime";

export type BackendBootTimings = {
  processStart: string;
  dbInitMs: number;
  servicesInitMs: number;
  totalBeforeUiMs: number;
};

export type FrontendMarks = {
  appMountedMs?: number;
  splashMountedMs?: number;
  authInitMs?: number;
  updateInitMs?: number;
  totalToReadyMs?: number;
};

export const useBootstrapStore = defineStore("bootstrap", () => {
  const backendTimings = ref<BackendBootTimings | null>(null);
  const frontendMarks = ref<FrontendMarks>({});
  const startTs = ref<number>(typeof performance !== "undefined" ? performance.now() : 0);

  function mark(name: keyof FrontendMarks, valueMs?: number) {
    frontendMarks.value[name] =
      valueMs ?? (typeof performance !== "undefined" ? performance.now() - startTs.value : 0);
  }

  function init() {
    try {
      EventsOn("bootstrap:backend-timings", (t: BackendBootTimings) => {
        backendTimings.value = t;
      });
    } catch (e) {
      console.warn("Wails runtime not available or EventsOn failed:", e);
    }

    // Auto-mock for dev environment if not set within a short time (e.g. pure Vite dev)
    if (import.meta.env.DEV) {
      setTimeout(() => {
        if (!backendTimings.value) {
          console.info("Mocking backend timings for dev environment (no Wails backend detected)");
          backendTimings.value = {
            processStart: new Date().toISOString(),
            dbInitMs: 10,
            servicesInitMs: 20,
            totalBeforeUiMs: 30,
          };
        }
        if (frontendMarks.value.updateInitMs === undefined) {
          frontendMarks.value.updateInitMs = 40;
        }
      }, 500);
    }
  }

  return {
    backendTimings,
    frontendMarks,
    startTs,
    mark,
    init,
  };
});

