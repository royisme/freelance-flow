import type { AppModule } from "@/modules/types";
import { financeMessages } from "@/modules/finance/i18n";
import { Wallet } from "lucide-vue-next";

// Static imports for finance views
import FinanceLayout from "@/views/finance/FinanceLayout.vue";
import FinanceOverview from "@/views/finance/index.vue";
import FinanceAccounts from "@/views/finance/accounts/index.vue";
import FinanceTransactions from "@/views/finance/transactions/index.vue";
import FinanceImport from "@/views/finance/import/index.vue";
import FinanceCategories from "@/views/finance/categories/index.vue";
import FinanceReports from "@/views/finance/reports/index.vue";
import FinanceSettings from "@/views/settings/FinanceSettings.vue";

export const financeModule: AppModule = {
  id: "finance",
  enabledByDefault: false,
  toggleable: true,
  nav: {
    labelKey: "nav.finance",
    key: "finance",
    icon: Wallet,
    children: [
      {
        labelKey: "finance.nav.overview",
        key: "finance/overview",
        moduleID: "finance",
      },
      {
        labelKey: "finance.nav.accounts",
        key: "finance/accounts",
        moduleID: "finance",
      },
      {
        labelKey: "finance.nav.transactions",
        key: "finance/transactions",
        moduleID: "finance",
      },
      {
        labelKey: "finance.nav.import",
        key: "finance/import",
        moduleID: "finance",
      },
      {
        labelKey: "finance.nav.categories",
        key: "finance/categories",
        moduleID: "finance",
      },
      {
        labelKey: "finance.nav.reports",
        key: "finance/reports",
        moduleID: "finance",
      },
    ],
  },
  routes: [
    {
      path: "/finance",
      component: FinanceLayout,
      meta: { requiresAuth: true, layout: "main", moduleID: "finance" },
      children: [
        { path: "", redirect: "/finance/overview", component: FinanceOverview },
        { path: "overview", component: FinanceOverview },
        { path: "accounts", component: FinanceAccounts },
        { path: "transactions", component: FinanceTransactions },
        { path: "import", component: FinanceImport },
        { path: "categories", component: FinanceCategories },
        { path: "reports", component: FinanceReports },
      ],
    },
  ],
  settingsPages: [
    {
      key: "finance",
      labelKey: "settings.finance.title",
      component: FinanceSettings,
      order: 50,
      moduleID: "finance",
    },
  ],
  messages: financeMessages,
};
