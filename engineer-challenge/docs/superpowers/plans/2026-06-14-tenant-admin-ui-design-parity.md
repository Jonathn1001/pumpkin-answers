# Tenant Admin UI Design-Parity — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the 6 remaining admin-UI features so the frontend reaches `.claude/UI_DESIGN.png` parity, each wired to a real, already-tested backend endpoint.

**Architecture:** Frontend-only build on the existing Vite/React/TS/Tailwind SPA (`apps/web`). React Query hooks (`src/api/hooks.ts`) call the Go API. No backend changes expected. Verification via the existing Playwright e2e suite, plus a unit test for the one pure util (`slugify`).

**Tech Stack:** React 19, Vite, TypeScript, Tailwind 4, @tanstack/react-query, react-router-dom, Playwright. Backend: Go (gin) — unchanged.

**Spec:** `docs/superpowers/specs/2026-06-14-tenant-admin-ui-design-parity-design.md`

---

## File Structure

| File | Responsibility | Action |
|---|---|---|
| `apps/web/src/lib/slugify.ts` | pure name→slug util | create |
| `apps/web/src/lib/slugify.test.ts` | unit tests for slugify (vitest) | create |
| `apps/web/src/components/DiffList.tsx` | render `Change[]` (path/type/before→after) | create (extract from CompareTab) |
| `apps/web/src/api/hooks.ts` | add `useUpdateTenant`, `useProcess` | modify |
| `apps/web/src/pages/TenantList.tsx` | row actions (edit/archive/clone); swap inline create → wizard | modify |
| `apps/web/src/components/CreateTenantWizard.tsx` | 2-step create wizard + success modal | create |
| `apps/web/src/pages/tabs/VersionsTab.tsx` | selectable rows + Version Details diff panel | modify |
| `apps/web/src/pages/tabs/CompareTab.tsx` | reuse `DiffList` | modify |
| `apps/web/src/pages/Runtime.tsx` | live process-a-claim screen + session log table | create |
| `apps/web/src/App.tsx` (router/nav) | register `/runtime` route + nav link | modify |
| `apps/web/src/schemaform/widgets/LogoUrlWidget.tsx` | URL input + image preview + remove | create |
| `apps/web/tests/*.spec.ts` | Playwright coverage per feature | modify/create |

**Conventions to follow (read before editing):** `src/api/hooks.ts` (hook + query-key + invalidation pattern), `src/api/client.ts` (`request`, `ApiError`), `src/pages/tabs/CompareTab.tsx` (current `Change[]` rendering), `src/App.tsx` (routing/nav), `src/schemaform/` (widget registration). Match existing style.

---

## Task 0 — Shared building blocks

**Files:** Create `src/lib/slugify.ts`, `src/lib/slugify.test.ts`, `src/components/DiffList.tsx`; modify `src/api/hooks.ts`, `src/pages/tabs/CompareTab.tsx`.

- [ ] **Step 1 — Write failing slugify unit test.** Decide test runner first: if `vitest` absent, `npm i -D vitest` in `apps/web` and add `"test": "vitest run"` to `package.json` scripts. Create `src/lib/slugify.test.ts`:

```ts
import { describe, it, expect } from "vitest";
import { slugify } from "./slugify";

describe("slugify", () => {
  it("lowercases and hyphenates spaces", () => expect(slugify("SafeGuard Insurance")).toBe("safeguard-insurance"));
  it("strips Vietnamese diacritics", () => expect(slugify("Bảo Việt Hà Nội")).toBe("bao-viet-ha-noi"));
  it("drops invalid chars and collapses hyphens", () => expect(slugify("A & B  --  C!")).toBe("a-b-c"));
  it("trims leading/trailing hyphens", () => expect(slugify("-Hello-")).toBe("hello"));
  it("caps at 63 chars", () => expect(slugify("x".repeat(80)).length).toBe(63));
  it("returns '' for empty/symbol-only", () => expect(slugify("@@@")).toBe(""));
});
```

- [ ] **Step 2 — Run, verify it fails.** Run: `npm run test` (in `apps/web`). Expected: FAIL (`slugify` not found).

- [ ] **Step 3 — Implement `src/lib/slugify.ts`:**

```ts
// name → URL-safe slug matching backend ^[a-z0-9][a-z0-9-]{0,62}$
export function slugify(name: string): string {
  return name
    .normalize("NFD").replace(/[̀-ͯ]/g, "") // strip diacritics
    .replace(/đ/g, "d").replace(/Đ/g, "D")
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-") // non-alnum → hyphen
    .replace(/-+/g, "-")          // collapse
    .replace(/^-+|-+$/g, "")      // trim
    .slice(0, 63)
    .replace(/-+$/g, "");          // re-trim if slice cut mid-hyphen
}
```

- [ ] **Step 4 — Run, verify pass.** Run: `npm run test`. Expected: PASS (6 tests).

- [ ] **Step 5 — Extract `DiffList`.** Read `CompareTab.tsx`; move its `Change[]` rendering into `src/components/DiffList.tsx` with contract `export function DiffList({ changes }: { changes: Change[] }): JSX.Element`. Render each: `path` · badge(`added|removed|changed`) · `left → right` (omit side that's empty). Update `CompareTab.tsx` to import and use `<DiffList changes={...} />`.

- [ ] **Step 6 — Add hooks** in `src/api/hooks.ts`, following the existing pattern:

```ts
export function useUpdateTenant(slug: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (b: { name?: string; status?: string }) =>
      request<Tenant>("PATCH", `/tenants/${slug}`, b),
    onSuccess: () => { qc.invalidateQueries({ queryKey: keys.tenants }); qc.invalidateQueries({ queryKey: keys.tenant(slug) }); },
  });
}
export function useProcess(slug: string) {
  return useMutation({
    mutationFn: (claim: Claim) => request<ClaimDecision>("POST", `/tenants/${slug}/process`, claim),
  });
}
```
(Verify exact `keys.*` names and `request` signature against the file before writing.)

- [ ] **Step 7 — Typecheck + commit.** Run: `npm run typecheck` (expect clean). `git add` the shared files. `git commit -m "feat(web): shared slugify util, DiffList component, update/process hooks"`.

---

## Task 1 — Slug auto-gen (create form)

**Files:** Modify `src/pages/TenantList.tsx` (create form). Test: `tests/smoke.spec.ts` (or new `tests/create-tenant.spec.ts`).

- [ ] **Step 1 — Playwright test (red).** Add a test: navigate to tenant list, type Name "Acme Health Co"; assert the slug field/preview shows `acme-health-co` without manual typing; submit; assert new tenant row appears.
- [ ] **Step 2 — Run, expect fail.** Run: `npx playwright test tests/<file> -g "auto slug"`. Expected: FAIL.
- [ ] **Step 3 — Implement.** In the create form: remove manual slug `useState` entry point; derive `const derivedSlug = slugify(name)`; show read-only slug text + an "Edit" toggle that reveals an input (override). Submit posts `slug: override || derivedSlug`. On 409 (`ApiError`, `er.fields`/message) show field error. Block submit when `derivedSlug === "" && !override`.
- [ ] **Step 4 — Run, expect pass.** Run: `npx playwright test tests/<file>`. Expected: PASS.
- [ ] **Step 5 — Commit.** `git commit -m "feat(web): auto-generate tenant slug from name (PLAN §1)"`.

---

## Task 2 — Version Details diff panel (VersionsTab)

**Files:** Modify `src/pages/tabs/VersionsTab.tsx`. Test: `tests/versions.spec.ts`.

- [ ] **Step 1 — Test (red).** Seeded tenant with ≥2 versions: click version row vN; assert a "Version Details" panel renders with at least one change entry; assert Rollback button present in panel.
- [ ] **Step 2 — Run, expect fail.**
- [ ] **Step 3 — Implement.** Add `selected` state (versionNumber). On row click set selected. Right panel: if `selected > 1` call `useDiff(\`${slug}@${selected-1}\`, \`${slug}@${selected}\`)` and render `<DiffList changes={data.changes} />`; if `selected === 1` show "initial version" note. Move Rollback button into the panel (keep table compact). Use existing `useDiff` (verify its arg shape: it takes `left,right` strings → `/diff?left=&right=`).
- [ ] **Step 4 — Run, expect pass.**
- [ ] **Step 5 — Commit.** `git commit -m "feat(web): version-to-version diff panel in VersionsTab (PLAN §2)"`.

---

## Task 3 — Tenant edit / archive + row actions (TenantList)

**Files:** Modify `src/pages/TenantList.tsx`. Test: `tests/tenant-actions.spec.ts`.

- [ ] **Step 1 — Test (red).** From list: click Archive on a tenant → its status badge becomes `archived`; click Activate → back to `active`. Click Clone → create wizard/form opens prefilled with that tenant as source.
- [ ] **Step 2 — Run, expect fail.**
- [ ] **Step 3 — Implement.** Per row add actions: **Edit** (inline name edit or small modal → `useUpdateTenant(slug).mutate({name})`), **Archive/Activate** (`useUpdateTenant(slug).mutate({status: next})`), **Clone** (opens create with `cloneFrom=slug` preset — integrates with Task 4 wizard; until then preset the existing form). Status badge reflects invalidated query.
- [ ] **Step 4 — Run, expect pass.**
- [ ] **Step 5 — Commit.** `git commit -m "feat(web): tenant edit/archive + row actions (PLAN §3)"`.

---

## Task 4 — Create Tenant wizard + success modal

**Files:** Create `src/components/CreateTenantWizard.tsx`; modify `src/pages/TenantList.tsx` (replace inline form). Test: `tests/create-wizard.spec.ts`.

- [ ] **Step 1 — Test (red).** Open "Create Tenant"; Step 1 enter Name (slug auto), choose source = default; Next → Step 2 shows a Tenant Summary (name, slug, source, claim-type/tier preview); click Create → success modal "Tenant Created Successfully" with "Go to Tenant"; clicking it navigates to `/t/<slug>`.
- [ ] **Step 2 — Run, expect fail.**
- [ ] **Step 3 — Implement.** `CreateTenantWizard`: step state `1|2`. Step 1 fields: Name (→`slugify`), source select (`default` | clone tenant). Step 2: summary; for clone, fetch source active config (`useConfig`/existing hook) to preview enabled claim types + tier count; for default, summarize default. Create via `useCreateTenant`. On success show modal with `Link to={\`/t/${slug}\`}`. Replace the inline form usage in `TenantList`. Reuse `slugify`. YAGNI: no full editor in wizard.
- [ ] **Step 4 — Run, expect pass.**
- [ ] **Step 5 — Commit.** `git commit -m "feat(web): create-tenant wizard with review + success modal (PLAN §4)"`.

---

## Task 5 — Runtime Logs / Processing Results (live-only)

**Files:** Create `src/pages/Runtime.tsx`; modify router/nav in `src/App.tsx`. Test: `tests/runtime.spec.ts`.

- [ ] **Step 1 — Test (red).** Navigate to `/runtime`; select seeded tenant `safeguard`; enter claim (type OUTPATIENT, amount under auto-approve threshold, submittedAt); Run → result panel shows approval outcome (e.g. `auto_approved`) and a non-empty trace; the run appears as a row in the session log table.
- [ ] **Step 2 — Run, expect fail.**
- [ ] **Step 3 — Implement.** `Runtime` page: tenant `<select>` (from `useTenants`), claim form (type, amount, submittedAt default today, dynamic custom-field inputs from selected tenant config), Run button → `useProcess(slug).mutate(claim)`. Render `ClaimDecision`: accepted/rejectionReasons, approval (`auto_approved` | routed tier/role/committee), requiredDocuments, notifications (event→channels), slaDeadline/slaDays, escalation, customFieldValidation, and `trace[]`. Keep results in `useState<ProcessedRow[]>` → render a logs table + click row → Log Detail panel. No persistence. Register route `/runtime` + nav link in `App.tsx` (follow existing nav pattern).
- [ ] **Step 4 — Run, expect pass.**
- [ ] **Step 5 — Commit.** `git commit -m "feat(web): runtime process-a-claim screen with decision + trace (PLAN §5)"`.

---

## Task 6 — Logo branding URL + preview

**Files:** Create `src/schemaform/widgets/LogoUrlWidget.tsx`; modify widget registry + branding schema usage. Test: extend `tests/smoke.spec.ts`.

- [ ] **Step 1 — Test (red).** In a tenant's Config → Branding, set logo URL field to a known URL; assert an `<img>` preview with that `src` renders; click Remove → field clears and preview disappears.
- [ ] **Step 2 — Run, expect fail.**
- [ ] **Step 3 — Implement.** Read `src/schemaform/` to find how widgets map to fields (by `widget` descriptor or field key). Add `LogoUrlWidget`: text input bound to value + `<img src={value} onError=hide>` preview + Remove button (sets ""). Register it for `branding.logoUrl` (via widget name in schema UI descriptor, or key match). Follow existing widget signature.
- [ ] **Step 4 — Run, expect pass.**
- [ ] **Step 5 — Commit.** `git commit -m "feat(web): logo URL field with live preview (PLAN §6)"`.

---

## Final verification

- [ ] Run full Playwright suite: `npx playwright test` (all green).
- [ ] Run `npm run typecheck` (clean) and `npm run lint`.
- [ ] Run backend Go tests `go test ./...` from `engineer-challenge/` (unchanged, green).
- [ ] Manual smoke via run skill: create tenant (wizard) → configure → process a claim → compare versions.

## Self-review notes
- Spec coverage: §1→T1, §2→T2, §3→T3, §4→T4, §5→T5, §6→T6, shared→T0. All covered.
- Hook/type names (`useUpdateTenant`, `useProcess`, `Claim`, `ClaimDecision`, `Change`, `keys.*`) to be confirmed against `src/api/*` at execution.
- Components whose bodies depend on unread files (CompareTab, App router, schemaform) specify exact path + contract + "read existing pattern first" — no fabricated internals.
