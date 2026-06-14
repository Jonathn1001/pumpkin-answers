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

test('zero-code 4th tenant onboarding via the full create page', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('button', { name: 'Create Tenant' }).click()
  // The create flow now lives on its own route, not a modal.
  await expect(page).toHaveURL(/\/tenants\/new$/)
  await page.getByPlaceholder('e.g. SafeGuard Insurance').fill('Demo Four')
  // Basics → seed the default config → first dimension step.
  await page.getByRole('button', { name: 'Next' }).click()
  // Step through the six dimension steps (Branding … Custom Fields) to Review.
  for (let i = 0; i < 6; i++) {
    await page.getByRole('button', { name: 'Next' }).click()
  }
  await page.getByRole('button', { name: 'Create Tenant' }).click()
  await expect(page.getByText('Tenant Created Successfully!')).toBeVisible()
  await page.getByRole('button', { name: 'Go to Tenant' }).click()
  await expect(page).toHaveURL(/\/t\/demo-four/)
  const slug = new URL(page.url()).pathname.split('/').pop()!
  await page.getByRole('tab', { name: 'Preview' }).click()
  await page.getByRole('button', { name: 'Run preview' }).click()
  await expect(page.getByText('Required documents')).toBeVisible()
  // Clean up so repeated runs don't accumulate tenants (archive = soft delete).
  await page.request.patch(`/api/tenants/${slug}`, { data: { name: 'Demo Four', status: 'archived' } })
})

test('editing a tenant publishes a new version; Save enables only on change', async ({ page }) => {
  await page.goto('/t/safeguard')
  await page.getByRole('button', { name: 'Edit configuration' }).click()
  await expect(page).toHaveURL(/\/tenants\/safeguard\/edit$/)

  // Save sits above the stepper and is disabled until the config changes.
  const save = page.getByRole('button', { name: 'Save & Publish' })
  await expect(save).toBeDisabled()

  // Toggle a cosmetic branding field so the change is dirty on every run.
  const field = page.getByRole('textbox').first()
  const cur = await field.inputValue()
  await field.fill(cur.endsWith(' *') ? cur.slice(0, -2) : cur + ' *')
  await expect(save).toBeEnabled()

  await save.click()
  await expect(page.getByText(/Published v\d+/)).toBeVisible()
  // Baseline resets after publish → Save disables again.
  await expect(save).toBeDisabled()
})

test('runtime processes a claim against a tenant active config', async ({ page }) => {
  await page.goto('/runtime')
  await pickSelect(page, 0, 'SafeGuard Insurance (safeguard)')
  await page.getByRole('button', { name: 'Process claim' }).click()
  await expect(page.getByText('auto-approved', { exact: true }).first()).toBeVisible()
})

test('deleting a tenant hides it from the list, then it can be restored', async ({ page }) => {
  await page.goto('/')
  await page
    .getByRole('row', { name: /SafeGuard Insurance/ })
    .getByRole('button', { name: 'Delete' })
    .click()
  // Toast + table reload happen only after the server confirms the delete.
  await page.getByRole('dialog').getByRole('button', { name: 'Delete' }).click()
  await expect(page.getByRole('link', { name: 'SafeGuard Insurance' })).toHaveCount(0)

  // It moves to the Archived view, where it can be restored.
  await page.getByText('Archived', { exact: true }).click()
  await page
    .getByRole('row', { name: /SafeGuard Insurance/ })
    .getByRole('button', { name: 'Restore' })
    .click()

  await page.getByText('Active', { exact: true }).click()
  await expect(page.getByRole('link', { name: 'SafeGuard Insurance' })).toBeVisible()
})
