// Locales index - re-export all locale messages
import enUS from "./en-US";
import zhCN from "./zh-CN";
import zodZhCN from "zod-i18n-map/locales/zh-CN/zod.json";
import zodEn from "zod-i18n-map/locales/en/zod.json";

// Helper to convert i18next style placeholders {{param}} to vue-i18n style {param}
function convertToVueI18n(messages: any): any {
  const newMessages: any = {};
  for (const key in messages) {
    const value = messages[key];
    if (typeof value === "string") {
      // Replace {{- param}} with {param} and {{param}} with {param}
      newMessages[key] = value
        .replace(/{{\s*-\s*([^}]+)\s*}}/g, "{$1}")
        .replace(/{{\s*([^}]+)\s*}}/g, "{$1}");
    } else if (typeof value === "object" && value !== null) {
      newMessages[key] = convertToVueI18n(value);
    } else {
      newMessages[key] = value;
    }
  }
  return newMessages;
}

export const messages = {
  "en-US": { ...enUS, ...convertToVueI18n(zodEn) },
  "zh-CN": { ...zhCN, ...convertToVueI18n(zodZhCN) },
};

export type Locale = "en-US" | "zh-CN";

export { enUS, zhCN };
