import { test, expect } from "@playwright/test";

const baseURL = process.env.PLAYWRIGHT_BASE_URL || "http://localhost:5173";
let reachable = false;

test.beforeAll(async ({ request }) => {
  try {
    const res = await request.get(baseURL);
    reachable = res.ok();
  } catch {
    reachable = false;
  }
});

test("invoice flow: select time entries, export pdf, trigger send", async ({
  page,
}) => {
  test.skip(!reachable, "dev server not running at baseURL");

  await page.goto(baseURL);

  // navigate to invoices (menu item text)
  await page.getByText(/invoices/i).first().click();

  // open entry selector for first row
  const firstRow = page.locator(".invoice-table").getByRole("row").nth(1);
  await firstRow.getByRole("button").nth(1).click(); // edit/select entries

  // select first available time entry if any
  const rows = page.locator(".n-data-table-tbody .n-data-table-tr");
  if (await rows.count()) {
    await rows.first().click();
  }

  // apply selection
  await page.getByRole("button", { name: /apply/i }).click();

  // download PDF
  await firstRow.getByRole("button").first().click();

  // trigger send
  await firstRow.getByRole("button").nth(2).click();

  // basic assertion: no error toast
  await expect(page.locator(".n-message__content").first()).not.toContainText(
    /error/i,
    { timeout: 2000 }
  );
});
