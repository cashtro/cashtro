# DECISIONS

## 2026-09-21 — Fill the blank mission header

- **Choice:** `CONTROL_PLANE_REPO=cashtro/cashtro`, orgs=`Evolu-Jeunes`,
  users=`cashtro`, `MODEL_ROUTE=OpenRouter`, `HARD_BUDGET_PER_RUN=$2`,
  `DEPLOY_TARGET=local only`, Prisma SQLite→Postgres.
- **Alternatives:** Treat Evolu-Jeunes as the manager; Bedrock-direct;
  Azure deploy in Phase 0; $0 budget.
- **Why:** This repo already calls itself the control plane. OpenRouter
  is the only model bus in the kernel. Local-only avoids live client
  sites. $2 is a hard cap we can enforce before any paid call. Castro
  left the blanks empty; autonomy says decide.

## 2026-09-21 — Do not rewrite the Go kernel

- **Choice:** Add a TypeScript workspace beside `cmd/` and `internal/`.
- **Alternatives:** Port the kernel to Node; new empty repo.
- **Why:** Mission: do not rewrite working code. The kernel is the
  Voltron runtime. The control plane is the fleet front door.

## 2026-09-21 — Queue without paid Redis

- **Choice:** SQLite `Job` table as the local queue. BullMQ only when
  `REDIS_URL` is set. Tests use SQLite.
- **Alternatives:** Require Redis in Phase 3; paid Upstash.
- **Why:** Redis is not on this machine. Paid Redis is a paid
  third-party service (stop-and-ask). Same dispatch API either way.

## 2026-09-21 — Branch name

- **Choice:** `cursor/control-plane-b1d0` off `main`, one PR for
  Phase 0–3.
- **Alternatives:** `feat/control-plane/0` per phase.
- **Why:** This Cloud Agent requires `cursor/<name>-b1d0`. Phases 0–3
  are one vertical slice; 4–7 wait for Castro to read REPORT.md.

## 2026-09-21 — Catalog ships in REPORT, not as scanned repos

- **Choice:** Recon REPORT lists kernel catalog ships (BTK, MD Clinic,
  Hypothèque, Éduconnexion, Proximity, ScanApp, delivery catalog) as
  "named, not GitHub-scanned."
- **Alternatives:** Invent Evolu-Jeunes repo slugs; omit catalog names
  until PAT lands.
- **Why:** Mission says do not guess missing repos. Castro still needs
  a phone-readable list of what he owns. The catalog is first-party
  seed data, not a GitHub listing.

## 2026-09-21 — Seed catalog ships as projects, ScanApp first

- **Choice:** `state/catalog-ships.json` is the source. Seed upserts
  every named ship. ScanApp gets draft capabilities and attached tasks.
  Live clients stay `production` with no repoUrl until recon.
- **Alternatives:** Wait for PAT before any ScanApp rows; invent
  `Evolu-Jeunes/scanapp`.
- **Why:** Castro said there was no ScanApp work in the plane. The
  backlog task was orphaned because slug `scanapp` did not exist. Do
  not invent a GitHub path.

## 2026-09-21 — HTTP adapter is the first real runtime

- **Choice:** Phase 4 ships `http` (Voltron kernel), `openrouter` /
  `bedrock` alias (no live spend without a key), `cursor` and
  `claude-code` as resident stubs. Depth hard-stop is 3 in code.
- **Alternatives:** Require Redis/BullMQ; call OpenRouter for the
  first E2E; wait for Evolu-Jeunes.
- **Why:** Voltron is already on :8080. OpenRouter spend needs a key
  and would burn the $2 cap. Redis is still not on this box.

## 2026-09-21 — Four screens on Fastify `/ui` plus Next.js `apps/web`

- **Choice:** Phone pane is Fastify `/ui/{fleet,queue,run,recon}`
  against the same DB. `apps/web` is the Next.js shell.
- **Alternatives:** Next only; wait for Phase 7.
- **Why:** Castro said keep working. The Fastify pane loads without a
  second process. Next talks to the same API.

## 2026-09-21 — Bind residents as local live verbs

- **Choice:** Six seats execute in-process. No browser farm, no prod
  deploy, no OpenRouter spend. Operator only hits 127.0.0.1.
- **Alternatives:** Leave “ready to bind”; wire real Playwright.
- **Why:** Castro said go attack. Binding a computer-use farm is a
  paid/scope jump. Local verbs unblock the desk tonight.

## 2026-09-21 — MCP reads files, not a second API

- **Choice:** `tools/mcp` reads inventory/backlog/catalog off disk.
- **Alternatives:** MCP that proxies Fastify; skip MCP until Evolu access.
- **Why:** Mission interop. Disk works when the API is down. No guessed repos.
