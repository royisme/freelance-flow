import type { APIRequestContext, Page } from "@playwright/test";
import { test, expect } from "@playwright/test";

export const baseURL =
  process.env.PLAYWRIGHT_BASE_URL || "http://localhost:5173";

export async function isDevServerReachable(
  request: APIRequestContext
): Promise<boolean> {
  try {
    const res = await request.get(baseURL);
    return res.ok();
  } catch {
    return false;
  }
}

export async function hasWailsRuntime(page: Page): Promise<boolean> {
  return page.evaluate(() => "go" in window);
}

export async function registerUserAndGoDashboard(
  page: Page,
  opts?: { currency?: string; timezone?: string }
): Promise<{ username: string; password: string; currency: string; timezone: string }> {
  const username = `e2e_${Date.now()}`;
  const password = "pass1234";
  const currency = opts?.currency ?? "USD";
  const timezone = opts?.timezone ?? "UTC";

  await page.goto(`${baseURL}/#/register`);

  // Step 1: Profile
  await page.getByPlaceholder(/username|用户名/i).fill(username);
  // Click Next
  await page.getByRole("button", { name: /next|下一步/i }).click();

  // Step 2: Password
  await page.getByPlaceholder(/password(?!.*confirm)|密码/i).first().fill(password);
  await page.getByPlaceholder(/Re-enter|confirm|确认/i).fill(password);
  // Click Next
  await page.getByRole("button", { name: /next|下一步/i }).click();

  // Step 3: Preferences
  // Currency is the second select (first is Language)
  // We need to match the select triggers. The form uses Select from shadcn-vue.
  // The structure is likely button role="combobox" inside the form items.

  // Note: SelectTrigger usually has role="combobox".
  // The first combobox is Language, second is Currency, third is Timezone.
  const selects = page.getByRole('combobox');

  // Select Currency (2nd combobox, index 1)
  await selects.nth(1).click();
  await page.getByRole('option', { name: new RegExp(currency, "i") }).click();

  // Select Timezone (3rd combobox, index 2)
  await selects.nth(2).click();
  await page.getByRole('option', { name: timezone }).click();

  // Create Profile
  await page.getByRole("button", { name: /create|注册|profile/i }).click();
  await expect(page).toHaveURL(/#\/dashboard/);

  return { username, password, currency, timezone };
}

export async function navigateByMenu(page: Page, label: RegExp): Promise<void> {
  await page.getByText(label).first().click();
  await test.step(`navigated to ${label}`, async () => {
    await expect(page.locator("main, .main-layout, .page-container")).toBeVisible();
  });
}

