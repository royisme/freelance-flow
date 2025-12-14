import { createApp } from "vue";
import { setupZodI18n } from "@/utils/validation";
import { createPinia } from "pinia";
import { createI18n } from "vue-i18n";
import router from "./router";
import App from "./App.vue";
import "./style.css";
// Naive UI recommended fonts
import "vfonts/Lato.css";
import "vfonts/FiraCode.css";
import { applyThemeToRoot } from "@/theme/tokens";

// ECharts global registration for desktop app
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { LineChart, BarChart } from "echarts/charts";
import {
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
} from "echarts/components";

use([
  CanvasRenderer,
  LineChart,
  BarChart,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  TitleComponent,
]);

// Import locale messages
import { messages } from "./locales";

if (typeof performance !== "undefined") {
  performance.mark("app:main:start");
}

// Apply a default theme early to avoid unstyled flash; AppProvider will keep it in sync.
applyThemeToRoot("light");

// I18n Setup
const i18n = createI18n({
  legacy: false, // Vue 3 Composition API
  locale: "zh-CN",
  fallbackLocale: "en-US",
  messages,
});

const pinia = createPinia();
const app = createApp(App);

app.use(pinia);
app.use(router);
app.use(i18n);
setupZodI18n(i18n);

app.mount("#app");

// Prevent duplicate Naive UI global styles during Vite HMR (dev only).
if (import.meta.hot) {
  import.meta.hot.dispose(() => {
    app.unmount();
    const root = document.getElementById("app");
    if (root) root.innerHTML = "";
  });
}
