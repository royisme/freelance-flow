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

// Import locale messages
import { messages } from "./locales";

if (typeof performance !== "undefined") {
  performance.mark("app:main:start");
}

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
