# ClaimConfig — Multi-Tenant Claims Configuration Platform

A multi-tenant platform for **configuring** how insurance claims are handled per tenant, and for **running** claims through a config-driven decision engine. Each tenant owns a versioned configuration spanning branding, claim types, approval routing, notifications, SLAs and custom fields. The same configuration that the admin UI edits is what the engine executes — there is no second source of truth.

> **Module:** `claimsplatform` · **Backend:** Go 1.25 + Gin + GORM/Postgres · **Frontend:** React 19 + Ant Design 6 + Vite

---

## Table of contents

- [What it does](#what-it-does)
- [Architecture](#architecture)
- [Configuration dimensions](#configuration-dimensions)
- [Project structure](#project-structure)
- [Getting started](#getting-started)
- [API reference](#api-reference)
- [Demo walkthrough (UI)](#demo-walkthrough-ui)
- [Testing](#testing)
- [Conventions & design notes](#conventions--design-notes)

---

## What it does

- **Onboard tenants** with a server-authoritative slug derived from the display name (diacritics, including Vietnamese, are stripped).
- **Configure each tenant** across six independent *dimensions*, validated as a whole before anything is saved.
- **Version every config change** — edits become drafts, drafts get published, and any previous version can be rolled back to.
- **Process claims** through a pure decision engine driven entirely by the active config: accept/reject, required documents, approval routing (auto-approve / tiered / committee), SLA deadlines, escalation and notification fan-out — each step explained in a human-readable trace.
- **Preview** a claim against any version or an unsaved draft without persisting anything.
- **Compare** any two configs (two versions of one tenant, or two different tenants) as a field-level changelog.

Three sample tenants — **SafeGuard**, **HealthFirst** and **GovHealth** — are seeded on startup so the app is useful immediately.

---

## Architecture

A clean, dependency-inverted layering: the domain has no outward dependencies; use cases orchestrate; HTTP and persistence are interchangeable adapters behind a port.

```mermaid
flowchart TD
    subgraph Frontend["Web (React + antd)"]
      UI["Tenants · Runtime · Compare"]
    end
    subgraph API["httpapi (Gin)"]
      H["thin handlers + error mapping"]
    end
    subgraph App["usecase (application services)"]
      S["Service: tenants / config / claims"]
    end
    subgraph Domain["domain (entities + ports)"]
      E["engine (pure)"]
      R["dimensions registry"]
      P["ConfigurationRepository (port)"]
    end
    subgraph Infra["adapters"]
      DB[("configrepo · GORM/Postgres")]
      MEM[("configrepo/memory · tests")]
    end

    UI -->|/api| H --> S --> E
    S --> P
    P -.implemented by.-> DB
    P -.implemented by.-> MEM
    E --> R
```

**The dimension plugin pattern.** Each configuration aspect is a self-contained package under `internal/dimensions/*` implementing one interface:

```go
type Dimension interface {
    Key() string
    Evaluate(cfg ConfigDocument, claim Claim, dec *ClaimDecision) // runtime logic
    Validate(cfg ConfigDocument) []FieldError                     // config-time rules
    DefaultConfig() any                                           // starter values
    UISchema() []FieldDescriptor                                  // widget metadata for the UI
}
```

Dimensions self-register via `init()` into `internal/registry`. The engine, the config validator and the JSON-Schema endpoint all iterate that registry, so adding a new dimension means adding one package — no central switch to edit.

**Claim processing pipeline** (`internal/engine`): ordering is explicit, not import-order-dependent.

```mermaid
flowchart LR
    A["claimTypes gate"] -->|rejected| Z["return decision"]
    A -->|accepted| B["approval · sla · customFields · branding"]
    B --> C["notifications (last — reacts to approval outcome)"]
    C --> Z
```

`claimTypes` runs first as a gate (a disabled type short-circuits to rejection). Core dimensions then run independently. `notifications` runs last because it reads the approval outcome to decide which events fire.

---

## Configuration dimensions

| Key | What it controls | Runtime effect |
|---|---|---|
| `branding` | Display name, logo, colors, support email | None (presentation only); validated for required name + `#RRGGBB` colors |
| `claimTypes` | Which of `OUTPATIENT / INPATIENT / DENTAL / MATERNITY / OPTICAL` are enabled and their required documents | **Gate** — disabled type ⇒ rejected; sets `requiredDocuments` |
| `approval` | `autoApproveThreshold`, model (`tiered` or `committee`), tiers / committee | Auto-approve under threshold, else route to a tier or committee |
| `notifications` | Channels (`email`/`sms`/`webhook`) and per-event channel map | Emits `claim_submitted` + `claim_auto_approved`/`claim_routed` |
| `sla` | `defaultDays`, per-claim-type overrides, escalation | Computes `slaDeadline` from `submittedAt`; attaches escalation |
| `customFields` | Tenant-defined fields (type, required, options, regex/min/max) | Validates the claim's `customFields` against the schema |

---

## Project structure

```
engineer-challenge/
├── cmd/server/             # main(): migrate → seed → serve (DATABASE_URL, PORT)
├── migrations/             # embedded SQL; golang-migrate runs them on boot
├── internal/
│   ├── domain/             # entities, value objects, errors, repository port (no deps)
│   ├── engine/             # pure ProcessClaim(cfg, claim) -> decision
│   ├── registry/           # Dimension interface + registration
│   ├── dimensions/         # branding, claimtypes, approval, notifications, sla, customfields
│   ├── usecase/            # application services (tenants, config, claims) + defaults
│   ├── configrepo/         # GORM/Postgres adapter (+ memory adapter + contract tests)
│   ├── configschema/       # JSON Schema per dimension (struct = source of truth)
│   ├── configvalidation/   # aggregates every dimension's Validate
│   ├── comparison/         # config diff via r3labs/diff
│   ├── httpapi/            # Gin router, thin handlers, error→status mapping, request logger
│   ├── tenantcontext/      # slug-resolving middleware
│   ├── seed/               # SafeGuard / HealthFirst / GovHealth sample tenants
│   └── slug/               # display name → URL-safe slug (Unicode/Vietnamese aware)
└── apps/web/               # React 19 + antd 6 + Vite admin UI
    └── src/
        ├── api/            # typed client + TanStack Query hooks
        ├── schemaform/     # AJV-driven dynamic form (renders the config schema)
        ├── pages/          # TenantList, TenantDetail (tabs), Runtime, Compare
        └── components/     # CreateTenantWizard, ClaimForm, DecisionView, DiffView
```

---

## Getting started

### Prerequisites

- **Go** ≥ 1.25
- **Node.js** ≥ 20 (for the web app)
- **PostgreSQL** ≥ 14

### 1. Start PostgreSQL

Any Postgres instance works. For a throwaway local one:

```bash
docker run --name claims-pg -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=claims -p 5432:5432 -d postgres:16
```

### 2. Run the API

The server runs migrations and seeds the three sample tenants automatically on startup.

```bash
export DATABASE_URL='postgres://postgres:postgres@localhost:5432/claims?sslmode=disable'
export PORT=8080            # optional, defaults to 8080
go run ./cmd/server
# → {"level":"INFO","msg":"api listening","port":"8080"}
```

| Env var | Required | Default | Notes |
|---|---|---|---|
| `DATABASE_URL` | ✅ | — | Postgres URL (used by both the migrator and GORM). Never logged. |
| `PORT` | — | `8080` | HTTP listen port |
| `GIN_MODE` | — | `release` | Set `debug` for Gin's verbose route logs |
| `LOG_FORMAT` | — | text | Set `json` for machine-parseable logs |

### 3. Run the web app

```bash
cd apps/web
npm install
npm run dev        # http://localhost:5173 — proxies /api → http://localhost:8080
```

Open <http://localhost:5173>: **Tenants** (list / create / clone / configure), **Runtime** (process a claim), **Compare** (diff two configs).

---

## API reference

Base path: `/api`. Errors share one envelope: `{"error":{"code","message",...}}`.

| Method | Path | Description |
|---|---|---|
| `GET` | `/tenants` | List tenants |
| `POST` | `/tenants` | Create a tenant (`{name, config?}`; `config` omitted ⇒ default) |
| `GET` | `/config-schema` | Per-dimension JSON Schema + UI widget metadata |
| `GET` | `/config-default` | The starter config a new tenant gets |
| `GET` | `/diff?left=<ref>&right=<ref>` | Field-level diff of two configs |
| `GET` | `/tenants/:slug` | Get one tenant |
| `PATCH` | `/tenants/:slug` | Update name / status |
| `GET` | `/tenants/:slug/config` | Active (published) config |
| `GET` | `/tenants/:slug/versions` | List all versions |
| `GET` | `/tenants/:slug/versions/:n` | Get version `n` |
| `POST` | `/tenants/:slug/versions` | Save a new **draft** (`{config, note}`) |
| `POST` | `/tenants/:slug/versions/:n/publish` | Publish version `n` |
| `POST` | `/tenants/:slug/rollback` | Roll back (`{targetVersion}`) |
| `POST` | `/tenants/:slug/process` | Run a claim against the **active** config |
| `POST` | `/tenants/:slug/preview` | Run a claim against a version or inline config (no writes) |

**Config refs** (for `/diff`): `slug` = the tenant's active config, `slug@n` = a specific version (e.g. `acme-health@2`).

---

## Demo walkthrough (UI)

A 5-minute click-through of the admin UI: **create a tenant → configure it → process claims → compare configs**. Make sure the API (`:8080`) and web app (`:5173`) are running, then open <http://localhost:5173>. The left sidebar has three sections: **Tenants**, **Runtime**, **Compare**.

### 1 · Create a tenant

1. Open **Tenants** (the default landing page) and click **Create Tenant** (top-right).
2. The **Add New Tenant** wizard opens with a stepper:
   - **Basics** — type a **Tenant name** (e.g. `Acme Health`). Under **Start from**, choose *Default configuration* (or *Clone from: …* to copy an existing tenant's config). Click **Next**. *(The URL slug, e.g. `acme-health`, is generated from the name by the server.)*
   - **Branding → Claim Types → Approval Rules → Notifications → SLA → Custom Fields** — each step is a form for that dimension. Adjust anything you like (e.g. on **Claim Types** toggle which types are enabled; on **Approval Rules** set the model and auto-approve threshold), or just click **Next** to accept the defaults.
   - **Review & Create** — check the summary (enabled claim types, approval model, SLA, custom-field count) and click **Create Tenant**.
3. On success you'll see **“Tenant Created Successfully!”** — click **Go to Tenant** to open its detail page.

> If the server rejects a field, the wizard jumps back to the step that owns it and highlights the error — fix it and continue.

### 2 · Configure it

You're now on the tenant page (`/t/acme-health`), which has four tabs: **Config · Preview · Versions · Compare**.

1. On the **Config** tab the whole configuration is shown as cards (one per dimension). Edit some fields, for example:
   - **Approval → Auto-approve threshold** → `100000`
   - **Claim Types → DENTAL** → enable it and add required documents (e.g. `receipt`, `dental_chart`)
   - **Branding → Display name** → `Acme Health Plan` (the page re-themes to the tenant's brand colors)
2. Click **Save draft + Publish**. This saves a new config **version** and publishes it in one step, so it immediately becomes the active config. Invalid values are highlighted inline and block publishing.
3. Open the **Versions** tab to see the history: every version with its status, author and note, with the live one tagged **active**. Click a version to see a field-level **diff vs the previous version**; for any non-active version, **Rollback to this version** restores it (rollback creates a *new* published version rather than deleting history).

### 3 · Process claims

1. Click **Runtime** in the sidebar (*Runtime — process a claim*).
2. In the **Claim input** card, pick your tenant from **Select a tenant…**.
3. Fill the claim form and click **Process claim**:
   - **Claim type** (e.g. `OUTPATIENT`), **Amount**, **Submitted at**, and optional **Custom fields (JSON)** (pre-filled with sample keys like `employeeId`, `policyNumber`).
4. The **Processing result** card shows the decision: accepted/rejected, required documents, approval outcome (**auto-approved** or **routed →** tier/committee), SLA deadline, notifications fired, custom-field validation, and a step-by-step **trace** explaining each dimension.
5. Each run is added to the **Processing log** table — click a row to re-inspect it. Try a few claims against your tenant to see the rules:
   - a **small** amount → **auto-approved** (≤ threshold),
   - a **large** amount → **routed** for approval,
   - a **disabled** claim type → **rejected**.

> The tenant's **Preview** tab does the same thing but lets you run a claim against a *specific version* or an unsaved draft **without persisting** — handy for “what would v1 have decided?” before you publish.

### 4 · Compare configs

- **Across tenants:** click **Compare** in the sidebar, choose a **left tenant** and a **right tenant** (try `safeguard` vs `healthfirst`), then click **Diff**. The changelog lists every differing field as *added / removed / changed* with its left and right values.
- **Across versions of one tenant:** open the tenant → **Compare** tab → pick two versions (e.g. `acme-health@1` vs `acme-health@2`) → **Diff** to see exactly what your edit in step 2 changed.

---

## Testing

**Backend**

```bash
go vet ./...
go test ./...          # unit tests run anywhere; configrepo integration tests use
                       # Testcontainers and require a running Docker daemon
```

The persistence layer is covered by a shared **contract test** (`configrepo/repotest`) run against both the in-memory and the Postgres (Testcontainers) implementations, so both adapters obey the same behavior.

**Frontend**

```bash
cd apps/web
npm run lint
npm run typecheck
npm run build
npx playwright test    # smoke test (apps/web/tests/smoke.spec.ts)
```

---

## Conventions & design notes

- **Config is the single source of truth.** The struct definitions drive the JSON Schema (`invopop/jsonschema`), the UI form, the validator and the engine — they cannot drift apart.
- **Library-first.** Diffing uses `r3labs/diff`; schema generation uses `invopop/jsonschema`; migrations use `golang-migrate` (embedded via `embed.FS`). No hand-rolled equivalents.
- **Thin handlers.** All HTTP→status mapping lives in one place (`httpapi.fail`); domain errors (`ValidationError`, `ErrTenantNotFound`, `ErrSlugTaken`, …) map to `422/404/409`.
- **Slugs are server-authoritative** and Unicode-aware: `"Phòng khám Đa khoa"` → `"phong-kham-da-khoa"`. Collisions get `-2`, `-3`, … suffixes.
- **Security defaults:** proxy headers are untrusted by default, request bodies are capped at 1 MiB, and **secrets are never logged** (the `DATABASE_URL` is deliberately kept out of all log lines).
```
