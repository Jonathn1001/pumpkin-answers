import { test, expect, type Page } from '@playwright/test'

// Pick an option from the Nth antd <Select> on the page. Scope to the open
// dropdown and match the option by its title attr (avoids hidden virtual-list rows).
async function pickSelect(page: Page, index: number, optionName: string) {
  await page.getByRole('combobox').nth(index).click()
  // The just-opened dropdown is appended last; .last() avoids a stale sibling dropdown.
  await page.locator(`.ant-select-item-option[title="${optionName}"]:visible`).last().click()
}

test('lists the three seeded tenants', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('link', { name: 'SafeGuard Insurance' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'HealthFirst' })).toBeVisible()
  await expect(page.getByRole('link', { name: 'GovHealth' })).toBeVisible()
})

async function preview(page: Page, slug: string) {
  await page.goto(`/t/${slug}`)
  await page.getByRole('tab', { name: 'Preview' }).click()
  await page.getByRole('button', { name: 'Run preview' }).click()
}

test('same canonical claim → different outcomes per tenant', async ({ page }) => {
  await preview(page, 'safeguard')
  await expect(page.getByText('auto-approved', { exact: true })).toBeVisible()
  await expect(page.getByText(/5 days →/)).toBeVisible()

  await preview(page, 'healthfirst')
  await expect(page.getByText(/routed →/).first()).toBeVisible()
  await expect(page.getByText(/7 days →/)).toBeVisible()

  await preview(page, 'govhealth')
  await expect(page.getByText(/Government Claims Committee/).first()).toBeVisible()
  await expect(page.getByText(/15 days →/)).toBeVisible()
})

test('config diff shows differences between two tenants', async ({ page }) => {
  await page.goto('/compare')
  await pickSelect(page, 0, 'SafeGuard Insurance (safeguard)')
  await pickSelect(page, 1, 'GovHealth (govhealth)')
  await page.getByRole('button', { name: 'Diff' }).click()
  await expect(page.getByText('changed').or(page.getByText('added')).first()).toBeVisible()
})

test('zero-code 4th tenant onboarding via the create wizard', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Create Tenant' }).click()
  const dialog = page.getByRole('dialog')
  await dialog.getByPlaceholder('e.g. SafeGuard Insurance').fill('Demo Four')
  await dialog.getByRole('button', { name: 'Next' }).click()
  await dialog.getByRole('button', { name: 'Create Tenant' }).click()
  await expect(page.getByText('Tenant Created Successfully!')).toBeVisible()
  await dialog.getByRole('button', { name: 'Go to Tenant' }).click()
  await page.getByRole('tab', { name: 'Preview' }).click()
  await page.getByRole('button', { name: 'Run preview' }).click()
  await expect(page.getByText('Required documents')).toBeVisible()
})

test('runtime processes a claim against a tenant active config', async ({ page }) => {
  await page.goto('/runtime')
  await pickSelect(page, 0, 'SafeGuard Insurance (safeguard)')
  await page.getByRole('button', { name: 'Process claim' }).click()
  await expect(page.getByText('auto-approved', { exact: true }).first()).toBeVisible()
})

test('archiving a tenant flips its status, and it can be reactivated', async ({ page }) => {
  await page.goto('/')
  const row = page.getByRole('row', { name: /SafeGuard Insurance/ })
  await row.getByRole('button', { name: 'Archive' }).click()
  await expect(row.getByText('archived')).toBeVisible()
  await row.getByRole('button', { name: 'Activate' }).click()
  await expect(row.getByText('active', { exact: true })).toBeVisible()
})
