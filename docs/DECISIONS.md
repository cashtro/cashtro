# DECISIONS

## 2026-09-22 — Graphify lives in the kernel HUD

- **Choice:** `internal/graphify` builds a Graphify-shaped map from brains,
  lines, ecosystem links, processes, ships, and the real boot→invoke path.
  The desk is a game HUD (`/api/live` every 1.2s). Giant is the world hub.
- **Why:** MCP Graphify is not bound in this environment. Castro asked to
  watch the code unfold and the ecosystem progress live, Giant included.
  Mermaid (`docs/MAP.md`) stays the readable schema next to the live map.

## 2026-09-22 — Use both lanes, and Mermaid too

- **Choice:** `model.chat` route `both` calls Ollama (internal) and, on
  OpenRouter, Kimi K3 plus GLM 5.3 at reasoning `max`. Maps stay in
  Graphify and in Mermaid (`docs/MAP.md`, `inventory/GRAPH.md`).
- **Why:** Internal work stays on the local daemon. The two top external
  models still answer the same prompt. Mermaid is the readable schema
  next to the Graphify map.

## 2026-09-22 — Model route is Kimi K3 and GLM 5.3 max

- **Choice:** Primary `moonshotai/kimi-k3`. Also `z-ai/glm-5.3` with
  reasoning effort `max`. Override with `OPENROUTER_MODEL` and
  `OPENROUTER_ALSO_MODEL`.
- **Why:** The swarm thinks through two top models. Kimi K3 is the
  default brain. GLM top runs at max reasoning when that route is called.

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
