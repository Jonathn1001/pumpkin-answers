# Tenant Admin UI — Design-Parity Build

**Date:** 2026-06-14
**Status:** Approved (design)
**Branch:** `feat/engineer-challenge`

## Goal

Build out the remaining admin-UI features so the frontend reaches parity with
`.claude/UI_DESIGN.png` (7 mockup screens) for the multi-tenant claims
configuration platform (AI Challenge 15). Every new piece of UI must be wired to
a real, already-tested backend endpoint — no dead UI.

## Key finding

Five of the six gaps are **frontend-only**. The exception is slug generation:
deriving a URL-safe slug from a tenant name is a **domain concern and lives in the
backend** (see Decisions) so every client gets identical, server-authoritative
slugs. The Go backend already exposes (and tests) every other endpoint required:

| Need | Endpoint | Notes |
|---|---|---|
| Update/archive tenant | `PATCH /api/tenants/:slug` | handler+route exist; FE hook was dropped |
| Run a claim | `POST /api/tenants/:slug/process` | returns full `ClaimDecision` incl. `trace[]` |
| Version-to-version diff | `GET /api/diff?left=slug@N-1&right=slug@N` | `parseRef` supports `slug@version` |
| Duplicate slug guard | — | `domain.ErrSlugTaken` → HTTP 409 |
| Versions / rollback | `GET /versions`, `POST /rollback` | already wired in `VersionsTab` |

No new backend logic is planned. If any backend change becomes necessary, it is
done test-first (TDD) and the existing Go test suite must stay green.

## Decisions (locked)

- **Runtime Logs screen:** live-only. Process a sample claim → render the result;
  keep processed claims in session state as a logs table. No DB persistence, no
  migration. (Matches challenge `processClaim` scope.)
- **Branding logo:** URL field + live image preview + Remove. No file upload, no
  object storage.
- **Slug:** generated **server-side** from the tenant name. `slug` is not part of
  the create payload at all — there is no client override (BE or FE).
  `POST /tenants` takes `name` (+ optional `cloneFrom`). A derived slug that
  collides gets a numeric suffix (`acme`, `acme-2`, …) so onboarding always
  succeeds. No client-side slug preview (avoids FE/BE logic drift) — the FE shows
  the slug returned in the create response.

## Shared building blocks (built first)

- `internal/slug` (Go) — `Make(name)` derives a slug matching backend
  `slugPattern` `^[a-z0-9][a-z0-9-]{0,62}$`: strips diacritics (incl. Vietnamese
  via NFD), collapses non-alphanumerics to `-`, trims, caps 63. Used by
  `CreateTenant` with numeric-suffix collision handling.
- Reuse the existing `DiffView` component (renders `Change[]`) for the version
  diff panel — no new component needed.
- Hooks: restore `useUpdateTenant` (PATCH), add `useProcess` (POST `/process`).

## Features (build order)

### 1. Slug generation (backend-owned)
Backend: `CreateTenant(name, cloneFrom)` derives the slug via `slug.Make(name)`
with numeric-suffix uniqueness; `slug` is removed from the request entirely (no
override). Name with no usable chars → validation error on `name`. Frontend: the
create form posts `name` only; the derived slug appears in the response /
refreshed list. TDD via Go tests (`internal/slug`, `internal/usecase`,
`internal/httpapi`).

### 2. Version Details / version-to-version diff (`VersionsTab`, mockup #5)
Selecting a version opens a right-hand "Version Details (vN)" panel listing
changes vs the previous version via `useDiff('slug@N-1','slug@N')`, rendered with
`DiffList` (path · added/removed/changed · before→after). Rollback button lives in
the panel. Edge: v1 has no predecessor → show "initial version" (diff vs default
config).

### 3. Tenant edit/archive + row actions (`TenantList`, mockup #1)
Per-row actions: **Edit** (rename), **Archive/Activate** (toggle status via
`useUpdateTenant` → PATCH), **Clone** (open create prefilled with `cloneFrom`).
Status badge updates immediately (query invalidation).

### 4. Create Tenant wizard (mockup #7, replaces inline form)
Two steps: (1) Basics — Name → auto slug, source config (default | clone tenant
X). (2) Review & Create — Tenant Summary (name/slug/source + preview of resulting
config: claim types, approval tiers read from source). On success → modal
"Tenant Created Successfully" → "Go to Tenant". The wizard does **not** embed the
full editor; detailed config happens in the editor tabs afterward (zero-code
onboarding flow). Reuses `slugify`.

### 5. Runtime Logs / Processing Results (mockup #6, live-only)
New `Runtime` page (+ route + nav entry). Pick tenant + enter a sample claim
(type, amount, submittedAt, custom fields) → `useProcess` → `/process` → render the
full `ClaimDecision`: accepted/rejection reasons, approval (auto_approved or routed
tier+role / committee), required documents, notifications (event→channels), SLA
deadline+days, escalation, custom-field validation, and `trace[]` per-dimension
explanations. Processed claims accumulate in session state as a logs table + Log
Detail panel. No persistence.

### 6. Logo branding URL + preview (mockup #3)
SchemaForm widget for `branding.logoUrl`: URL input + live image preview + Remove
(clears URL). Existing color swatches kept.

## Verification ("đảm bảo tính năng")

- Extend the existing Playwright smoke: slug auto-gen; version-diff shows changes;
  archive flips status; create-wizard → success modal → go to tenant; runtime
  process renders a decision + trace.
- Backend Go tests remain green (no BE changes expected).
- Each feature is committed and verified independently, in the order above
  (shared → 1 → 2 → 3 → 4 → 5 → 6).

## Out of scope (YAGNI)

- DB persistence of processed claims / runtime logs.
- File upload / object storage for logos.
- Embedding the full config editor inside the create wizard.
- Any backend feature work beyond wiring (unless a gap forces a test-first change).
