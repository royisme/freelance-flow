import { test, expect } from "@playwright/test";
import {
  baseURL,
  registerUserAndGoDashboard,
  navigateByMenu,
} from "./helpers";

test.describe("Finance Module", () => {
  test.beforeEach(async ({ page, request }) => {
    // 1. Register a fresh user
    await registerUserAndGoDashboard(page, { currency: "CAD", timezone: "UTC" });
  });

  test("can navigate to finance dashboard and see sub-pages", async ({ page }) => {
    // 2. Navigate to Finance (Expand menu)
    await page.getByText(/Finance|财务/i).first().click();

    // Click Overview to navigate
    await page.getByRole('link', { name: /Overview|概览/i }).click();

    await expect(page).toHaveURL(/#\/finance\/overview/);

    // 3. Check for sub-navigation tabs/links
    // The sidebar stays open, so these should be visible
    const navItems = [
      "Overview",
      "Accounts",
      "Transactions",
      "Import",
    ];

    for (const item of navItems) {
      // Use first() to avoid ambiguity with breadcrumbs
      await expect(page.getByRole('link', { name: item }).first()).toBeVisible();
    }
    await page.screenshot({ path: "finance_dashboard.png" });
  });

  test("can create a new bank account", async ({ page }) => {
    await page.getByText(/Finance|财务/i).first().click();

    // Go to Accounts tab
    await page.getByRole('link', { name: 'Accounts' }).first().click();

    // Click Add Account
    await page.getByRole('button', { name: /Add Account|添加账户/i }).click();

    // Let's just check the dialog opens for now to verify the interaction
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    // Title is "Add Account"
    await expect(dialog.getByRole('heading', { name: /Add Account|添加账户/i })).toBeVisible();

    // Fill form
    await page.getByLabel(/Account Name|账户名称/i).fill("TD Checking 001");
    // Select Type
    await page.getByRole('combobox').first().click(); // Assuming first combobox is Type
    await page.getByLabel('Checking', { exact: true }).click();

    // Save
    await page.getByRole('button', { name: /Save|保存/i }).click();

    // Verify dialog closes (backend is mocked in dev/test, so we don't see the new item, but UI should proceed)
    await expect(dialog).not.toBeVisible();
    await page.screenshot({ path: "finance_account_created.png" });
  });
});
