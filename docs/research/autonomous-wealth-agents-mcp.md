# Cashtro OS — Autonomous wealth agents & MCP playbook

Research swarm: **20 agents**, 2026-09-21. Constraint: new agentics register into this kernel. They do not become a second product.

Sources: kernel/process table in this repo; MCP spec [2026-07-28](https://modelcontextprotocol.io/specification/2026-07-28); [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk); [official registry](https://modelcontextprotocol.io/registry/about).

---

## Verdict

Cashtro already has the right primitives (processes, verbs, journal, human confirms, delivery line). **Wealth is blocked by missing runtime, not missing ideas:** outbound after `comms.allow` does not send, six residents are stubs, state dies on restart, nothing ticks overnight, and there is no MCP adapter wrapping `POST /api/agents/{id}/invoke`.

Autonomy that makes money is a **ladder with unlocks**, not a chat that never sleeps. Sell **outcomes** (ship packs, care retainers, vertical MCP seats). Meter agent hours internally as COGS.

---

## How to create agents (in this kernel)

There is no plugin loader. Compile-time:

1. Implement `kernel.Agent` (`Spec`, `Boot`, `Invoke`) in `internal/agents/` — copy `plannerInvoke` / `researchInvoke` in `live.go`, or `resident()` in `builtin.go`.
2. Append to `Builtins()` in `internal/agents/builtin.go`.
3. `agents.Boot()` registers and autostarts. HTTP is already generic: `POST /api/agents/{id}/invoke`.

| ID | Mode today | Capabilities |
|---|---|---|
| init | live | `os.about` |
| delivery | live | `delivery.list/create/advance/profile` |
| router | live iff OpenRouter key | `model.status`, `model.chat` |
| research | live | `research.list`, `research.ingest`, `note.write` |
| explorer | live | `explorer.search` |
| memory | live | `memory.store`, `memory.recall` |
| comms | live | `comms.send/pending/allow/deny` (park only) |
| planner | live | `planner.backlog` |
| operator, reviewer, architect, deploy, security, investigator | **resident stubs** | browse, watch, plan, release, triage, trace |

**Smallest new wealth agent:** `rainmaker` with `rainmaker.scan` + `rainmaker.pitch` → `planner.backlog` / `delivery.create` → `comms.send` → human allow. Tests to copy: `internal/agents/builtin_test.go`, `internal/server/server_test.go`.

---

## How to create MCPs

MCP is a **kernel adapter**, not a second OS. Host (Cursor/Claude) → Client → Server. Servers expose **tools** (verbs), **resources** (lists), **prompts** (templates). Spec `2026-07-28` is **stateless JSON-RPC**; sessions/initialize are gone. Long work uses the **Tasks** extension. Mid-tool human input uses elicitation (`input_required`). Sampling is deprecated — Cashtro already has optional `model.chat`.

**Build:** official `github.com/modelcontextprotocol/go-sdk` next to `cmd/cashtro`. One in-process server that calls `Kernel.InvokeCap`. Do not fork a TypeScript OS.

**Surface = kernel names.**

- Always read: `os.about`, `explorer.search`, `delivery.list/profile`, `research.list`, `memory.recall`, `comms.pending`, `model.status`
- Autonomous writes (allowlisted): `research.ingest`, `memory.store`, `planner.backlog`, `delivery.create`
- HITL: `delivery.advance`, `comms.send`
- **Never as agent tools:** `comms.allow`, `comms.deny` (desk HTTP only)
- Resources: `cashtro://os|agents|ships|notes|memory|confirms|events`

**Auth:** local Cursor = stdio + OS-user trust. Remote later = OAuth 2.1. Annotate `readOnlyHint` / `destructiveHint`; still **enforce** confirms in the kernel.

**Distribution (listings are free; cash is hosted seats):** GitHub → [official MCP Registry](https://modelcontextprotocol.io/registry/about) → mcp.so / Smithery → Cursor Marketplace → Claude Connectors. First products: Delivery MCP, Confirm MCP, Research Desk MCP.

---

## How to make them autonomous (without becoming reckless)

### Runtime (steal, don’t replace)

Own the loop in Go. Model only fills the next verb.

| First | Pattern | Why cash |
|---|---|---|
| 1 | Event-driven wake (cron + webhooks) | Night shift starts without a chat |
| 2 | Deterministic tool-calling loop + token/CAD caps | Turns a wake into ingest → backlog → parked outbound |
| 3 | Durable steps + idempotency; confirm = wait-for-event | Survives until the invoice |

v1 host: **cheap VPS + systemd + this binary on `:8080`**. Persist the desk. Watchdog = `nightwatch` agentic, not a sidecar OS. Fly.io later. Do not rewrite onto Workers/Vercel/GHA.

Copy from other stacks into **this** kernel: MCP `tools/list` schemas, LangGraph interrupt/resume bound to confirms, OpenAI handoff vs tool-call, Restate/DO keyed single-writer per `ship_id`, Claude Pre/Post invoke hooks.

### Autonomy ladder (publish this)

| Level | Name | Money |
|---|---|---|
| 0 | Observe | Read only |
| 1 | Draft | Parks confirms; human sends |
| 2 | Reversible | Internal writes; **$0 leaves the org** |
| 3 | Bounded auto | Numeric fences (CAD, tokens, steps, recipients) |
| 4 | Funds / irreversible | Dual-control; initiator ≠ approver |
| 5 | Orchestrate | Fleet inside policy — **not** self-set P&L |

An agent never raises its own level. Kill switch, journal, and confirm store are not allowlisted tools. Default ceiling for revenue agents: **L3**. L4 is a named exception.

### Gaps that block money (ranked)

1. **Allow does not send** — `DecideConfirm` is `"no outbound bind yet"`.
2. **Worker bind** — operator/reviewer/deploy/architect stay `fn == nil`.
3. **No scheduler** — mail is never consumed; no overnight tick.
4. **RAM-only state** — ships/notes/memory/confirms die on restart (journal cap 200).
5. **No MCP package** — REST only.
6. **No billing process** — Stripe should stay ModeResident until a restricted key binds (same as router).
7. **No eval harness** on plan → ship → confirm → deliver.

---

## What to sell (cash, not demos)

Price **eligible tasks + proof**. Walkthrough artifacts are the Definition of Done. Tokens are COGS (`meter.cost`).

### Fastest cash (existing ships)

1. **SKU Ship Pack** — $5k–$15k, 50% deposit, frozen at concept, production = PR + `reviewer.watch`.
2. **Prepaid burst block** — ticket cap, not raw hours.
3. **Production Care** — $2.5k–$8k/mo on BTK / MD Clinic / Hypothèque / Éduconnexion.
4. **Concept sprint** — $2k–$5k idea → concept.
5. **Workflow SLA** — later; SLA the named workflow, not model accuracy.

**WordPress factory SKU (CSF-LIVE):** CAD ~$6,900, 8–12 pages + ACF + FR/EN + 3 forms, human gates on DNS/prod. Official `WordPress/mcp-adapter` + WP-CLI + Playwright QA.

### Recurring verticals (tenant process, not a fork)

| Product | From | Starter CAD/mo | Autonomy |
|---|---|---|---|
| IntakeDesk Legal | BTK | ~$490 | L2; never advise |
| ClinicBook | MD Clinic | ~$590 | L3 in-policy; Law 25 |
| PreQual QC | Hypothèque | ~$690 | L2; no rate advice |
| EnrollConnect | Éduconnexion | ~$390 | L3–L4 on published catalog |

### Skills as SKUs

Mine production ships into `SKILL.md` (weekly + repeatable + time saved). Vertical packs $499–$999; QC agency bundle ~$1,499; run on kernel $149–$399/mo seat. Notes → skill only after ≥3 sourced notes **and** a production twin.

### Revenue / hunter / merch (later binds)

- **`revenue`:** `icp.qualify` → `seller.draft` (CASL consent class) → `pipeline.upsert`. Gmail **draft** after allow; never auto-send. LinkedIn is a CEM.
- **`hunter`:** SEAO / CanadaBuys / public rebuilds. Score `fit × fee × close`. Park idea ships; no auto-bid.
- **E-com:** `merch`, `stock`, `recover` — 15% of **proven incremental** GMV, 10% holdout, confirm is the billing gate.

---

## Swarm contract (existing 14 processes)

Keep InvokeCap (syscall), mail (typed IPC with `ship_id`), journal (audit). Advance is illegal without the artifact:

`plan.ready` (architect → planner) → `work.assigned` (operator **or new `coder`**) → `work.done` → `review.pass` → `release.ready` → `released` → `comms.send` (human allow).

Bind for cash: router key → `architect.plan` → `operator.browse` → `reviewer.watch` → `deploy.release`. Architect never deploys. Operator never mails the client. Comms never advances the board.

**New processes (register here):** `nightwatch`, `billing`, `pricer`, `meter`, `proof`, `care`, `rainmaker`/`hunter`, `coder`, `content`/`seo`. Memory needs a **ledger** (`amount_minor`, `outcome`) plus `memory.week` / `memory.score` — FIFO `Fact` text cannot compound wealth.

---

## 90-day sequence (capability order, CAD ranges not promises)

Do **not** rebuild ScanApp, a second product, a custom LLM, or autonomous mail.

| Window | Bind / sell | Stop if |
|---|---|---|
| W1 | Comms → Gmail drafts; 5 upsell ships; desk audits | — |
| W2 | Research ingest packs (law, clinic) | Do not stand up a vector DB |
| W3 | WordPress MCP + Care Line on 2 live sites | Do not rebuild WP in Next |
| W4 | Stripe invoice + first allowed outbound | **Zero closes → stop binding, only sell audits** |
| W5 | Paid `architect.plan` workshop | No generic app factory |
| W6 | Front Desk intake if a Care client pays | No operator yet |
| W7–8 | `deploy.release` + `reviewer.watch` on one paid ship | No video platform |
| W9–10 | `security.triage` Care add-on; investigator hour-block | ScanApp only if a named buyer |
| W11 | One paid `operator.browse` mandate | No unattended desktop |
| W12 | Freeze new agents; package Cashtro Desk | Repeat Care + Front Desk + one production ship |

Indicative band if closes happen: month-3 cash **$10–25k**, MRR **$4–12k**. Not a forecast.

---

## Do this next (kernel)

1. Persist catalog + notes + confirms + memory (`internal/store`).
2. `nightwatch` ticker + scheduler (`internal/scheduler`).
3. Outbound bind on `comms.allow` (draft-only first).
4. `cmd/cashtro-mcp` wrapping `InvokeCap`.
5. `billing` resident until `STRIPE_RESTRICTED_KEY`; Checkout after confirm.
6. Eval: plan → ship → confirm → (mocked) deliver.

**Line:** *Cashtro earns autonomy. Money never does.*
