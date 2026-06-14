import { test, expect } from "@playwright/test";

test("lists the three seeded tenants", async ({ page }) => {
  await page.goto("/");
  await expect(
    page.getByRole("link", { name: "SafeGuard Insurance" }),
  ).toBeVisible();
  await expect(page.getByRole("link", { name: "HealthFirst" })).toBeVisible();
  await expect(page.getByRole("link", { name: "GovHealth" })).toBeVisible();
});

async function preview(page: import("@playwright/test").Page, slug: string) {
  await page.goto(`/t/${slug}`);
  await page.getByRole("button", { name: "Preview", exact: true }).click();
  await page.getByRole("button", { name: "Run preview" }).click();
}

test("same canonical claim → different outcomes per tenant", async ({
  page,
}) => {
  await preview(page, "safeguard");
  await expect(page.getByText("auto-approved", { exact: true })).toBeVisible();
  await expect(page.getByText(/5 days →/)).toBeVisible();

  await preview(page, "healthfirst");
  await expect(page.getByText(/^routed →/)).toBeVisible();
  await expect(page.getByText(/7 days →/)).toBeVisible();

  await preview(page, "govhealth");
  await expect(page.getByText(/^routed → Government Claims Committee/)).toBeVisible();
  await expect(page.getByText(/15 days →/)).toBeVisible();
});

test("config diff shows differences between two tenants", async ({ page }) => {
  await page.goto("/compare");
  await page.locator("select").first().selectOption("safeguard");
  await page.locator("select").nth(1).selectOption("govhealth");
  await page.getByRole("button", { name: "Diff" }).click();
  await expect(
    page.getByText("changed").or(page.getByText("added")).first(),
  ).toBeVisible();
});

test("zero-code 4th tenant onboarding via UI", async ({ page }) => {
  await page.goto("/");
  await page.getByPlaceholder("slug").fill("demo4");
  await page.getByPlaceholder("name").fill("Demo Four");
  await page.getByRole("button", { name: "Create" }).click();
  await expect(page.getByRole("link", { name: "Demo Four" })).toBeVisible();
  await preview(page, "demo4"); // default config enables OUTPATIENT → processable
  await expect(
    page.getByText("Required documents:", { exact: true }),
  ).toBeVisible();
});
