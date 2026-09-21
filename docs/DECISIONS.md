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
