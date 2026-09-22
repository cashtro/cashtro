# DECISIONS

## 2026-09-22 — Vapi talk works without a pre-created assistant

- **Choice:** `POST /api/vapi/talk` hits live Vapi whenever `VAPI_API_KEY`
  is set. The kernel creates the Cashtro Teal assistant on boot if none
  exists, or sends a transient assistant on `/chat`. `Start voice` on the
  desk uses the mic locally and loads the Vapi HTML widget when
  `VAPI_PUBLIC_KEY` is set. `.env` is loaded at process start.
- **Why:** Castro asked to integrate Vapi and make it work. The old path
  stayed local unless both a key and an assistant ID were already set.
- **Stop:** Do not buy Vapi numbers or fire paid outbound without allow.

## 2026-09-22 — L'Inquisiteur contradicts, blocks, and steps up every department

- **Choice:** Live kernel agentic `inquisitor`. Ten lines + five brains +
  the corporation are departments. Each has a `smarter` loop that fills
  its own gaps. `Spec.Rules` divorces every agentic. Mutating inquisitor
  verbs park specific questions on the ask-gate before they act. Dissent
  can `block` a department; `vapi.call` refuses a blocked line until that
  department self-improves. Boot ticks every department once and does
  **not** freeze the desk.
- **Why:** Castro said the chain was not specialized enough. Push every
  department to the maximum so it can close its own lacunes. Always
  understand before targeting an action. Always contradict and test
  toward the most optimized option.
- **Stop:** Do not auto-block the whole corporation at boot. Do not buy
  Vapi numbers or fire paid calls.

## 2026-09-22 — Evolu-Jeunes/Teal is a fresh voice repo

- **Choice:** Put the Vapi/Teal voice brain in a fresh Evolu-Jeunes
  product at `teal/` (`github.com/Evolu-Jeunes/Teal`). Cashtro OS stays
  the manager. Teal listens on `:8090`. Publish with
  `scripts/publish-evolu-jeunes-teal.sh` once a PAT can create org
  repos (this Cloud Agent token returns 403 on Evolu-Jeunes writes).
- **Why:** Castro asked to put this in another repo, a fresh one, in
  Evolu-Jeunes.
- **Stop:** Do not buy Vapi numbers. Do not treat Teal as a second OS
  that replaces Cashtro.

## 2026-09-22 — Vapi is a kernel voice lane across every bridge

- **Choice:** Absorb a live `vapi` agentic into Cashtro OS. Talk
  (`POST /api/vapi/talk`), outbound (`POST /api/vapi/call` then human
  allow → `vapi.fire`), voice switch (`POST /api/vapi/voice`), and
  line context (`POST /api/vapi/bridge`) run on the kernel. Vapi is
  listed on **every** ecosystem `Links()` hop so Proximity, Scan App,
  Panda, NFT/Giant, école, marketing, Empire, trading, and Pandora
  share the same voice script. Teams scrape stays on `teal`.
- **Why:** Castro asked to put Vappy in function with the agentic
  project: automate calls, talk with him, change voice, build it in
  the repo, connect it to infrastructure not just Teams, and run it
  across bridge projects.
- **Stop:** Do not buy Vapi numbers or fire live paid calls until
  Castro confirms. `$2/run` still applies. Local talk works without
  `VAPI_API_KEY`. Free Vapi numbers cannot outbound.

## 2026-09-22 — Teal brain is Teams-only, voice is Vapi, Cursor is a card

- **Choice:** Absorb a live `teal` agentic into Cashtro OS. Ingest
  surface is Microsoft Teams only. Voice is Vapi (the "Vappy"
  assistant), now also a dedicated kernel lane. Cursor Cloud Agents
  launch from a Teams Adaptive Card after a human confirm. New intel
  registers back into this kernel.
- **Why:** Castro asked to connect the Teams AI Teal voice AI (Vapi) to the
  Pandora brainstorming brain, scrape Teams, manage interns, listen to calls,
  and keep getting smarter. Teams MCP and Zapier Teams have no tenant bind
  here, so Graph scrape waits on `TEAMS_TOKEN`. The public Vapi link we can
  stand up without that bind is https://vapi.ai/ (dashboard
  https://dashboard.vapi.ai/). Lock the exact assistant with `VAPI_SHARE_URL`.
- **Alternatives:** A second web dock; Slack; paid Vapi signup from this agent.
- **Stop:** Do not buy Vapi numbers or post to a live Teams webhook until Castro
  confirms. Outbound stays on the human gate.

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
