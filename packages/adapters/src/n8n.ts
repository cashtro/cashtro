import { readFileSync } from "node:fs";
import path from "node:path";
import type { Task } from "@cashtro/sdk";
import type { AdapterRunStatus, AgentAdapter, Health, RunContext, RunHandle } from "./types.js";
import { assertDepth } from "./types.js";
import { createHttpAdapter } from "./http.js";

export type N8nSpecialist = {
  id: string;
  seat: string;
  name?: string;
  capability?: string;
  webhook: string;
};

export type N8nOpts = {
  baseUrl?: string;
  apiKey?: string;
  fetchImpl?: typeof fetch;
  fallback?: AgentAdapter;
  fleet?: N8nSpecialist[];
};

/** Width 58: 40 verbs + 14 self-improve + 4 Steel. Used when state/n8n-fleet.json is not on disk. */
export const EMBEDDED_N8N_FLEET: N8nSpecialist[] = [
  { id: "n8n-01", seat: "init", name: "os.pulse", capability: "os.about", webhook: "os-pulse" },
  { id: "n8n-02", seat: "init", name: "os.manifest", capability: "os.about", webhook: "os-manifest" },
  { id: "n8n-03", seat: "delivery", name: "ship.list", capability: "delivery.list", webhook: "ship-list" },
  { id: "n8n-04", seat: "delivery", name: "ship.create.idea", capability: "delivery.create", webhook: "ship-create" },
  { id: "n8n-05", seat: "delivery", name: "ship.advance", capability: "delivery.advance", webhook: "ship-advance" },
  { id: "n8n-06", seat: "delivery", name: "ship.profile", capability: "delivery.profile", webhook: "ship-profile" },
  { id: "n8n-07", seat: "router", name: "model.status", capability: "model.status", webhook: "model-status" },
  { id: "n8n-08", seat: "router", name: "model.chat.cheap", capability: "model.chat", webhook: "model-chat-cheap" },
  { id: "n8n-09", seat: "research", name: "research.list", capability: "research.list", webhook: "research-list" },
  { id: "n8n-10", seat: "research", name: "research.ingest", capability: "research.ingest", webhook: "research-ingest" },
  { id: "n8n-11", seat: "research", name: "note.write", capability: "note.write", webhook: "note-write" },
  { id: "n8n-12", seat: "explorer", name: "search.agents", capability: "explorer.search", webhook: "search-agents" },
  { id: "n8n-13", seat: "explorer", name: "search.ships", capability: "explorer.search", webhook: "search-ships" },
  { id: "n8n-14", seat: "explorer", name: "search.notes", capability: "explorer.search", webhook: "search-notes" },
  { id: "n8n-15", seat: "operator", name: "browse.local", capability: "operator.browse", webhook: "browse-local" },
  { id: "n8n-16", seat: "operator", name: "browse.api", capability: "operator.browse", webhook: "browse-api" },
  { id: "n8n-17", seat: "operator", name: "browse.health", capability: "operator.browse", webhook: "browse-health" },
  { id: "n8n-18", seat: "reviewer", name: "watch.artifacts", capability: "reviewer.watch", webhook: "watch-artifacts" },
  { id: "n8n-19", seat: "reviewer", name: "watch.report", capability: "reviewer.watch", webhook: "watch-report" },
  { id: "n8n-20", seat: "reviewer", name: "watch.ui", capability: "reviewer.watch", webhook: "watch-ui" },
  { id: "n8n-21", seat: "architect", name: "plan.control-plane", capability: "architect.plan", webhook: "plan-control-plane" },
  { id: "n8n-22", seat: "architect", name: "plan.scanapp", capability: "architect.plan", webhook: "plan-scanapp" },
  { id: "n8n-23", seat: "architect", name: "plan.refuse-prod", capability: "architect.plan", webhook: "plan-refuse-prod" },
  { id: "n8n-24", seat: "deploy", name: "release.local", capability: "deploy.release", webhook: "release-local" },
  { id: "n8n-25", seat: "deploy", name: "release.check", capability: "deploy.release", webhook: "release-check" },
  { id: "n8n-26", seat: "deploy", name: "release.refuse-prod", capability: "deploy.release", webhook: "release-refuse-prod" },
  { id: "n8n-27", seat: "security", name: "triage.access", capability: "security.triage", webhook: "triage-access" },
  { id: "n8n-28", seat: "security", name: "triage.secrets", capability: "security.triage", webhook: "triage-secrets" },
  { id: "n8n-29", seat: "security", name: "triage.flags", capability: "security.triage", webhook: "triage-flags" },
  { id: "n8n-30", seat: "memory", name: "memory.store", capability: "memory.store", webhook: "memory-store" },
  { id: "n8n-31", seat: "memory", name: "memory.recall", capability: "memory.recall", webhook: "memory-recall" },
  { id: "n8n-32", seat: "memory", name: "memory.digest", capability: "memory.recall", webhook: "memory-digest" },
  { id: "n8n-33", seat: "comms", name: "comms.send", capability: "comms.send", webhook: "comms-send" },
  { id: "n8n-34", seat: "comms", name: "comms.pending", capability: "comms.pending", webhook: "comms-pending" },
  { id: "n8n-35", seat: "comms", name: "comms.allow", capability: "comms.allow", webhook: "comms-allow" },
  { id: "n8n-36", seat: "planner", name: "backlog.now", capability: "planner.backlog", webhook: "backlog-now" },
  { id: "n8n-37", seat: "planner", name: "backlog.scanapp", capability: "planner.backlog", webhook: "backlog-scanapp" },
  { id: "n8n-38", seat: "planner", name: "backlog.evolu", capability: "planner.backlog", webhook: "backlog-evolu" },
  { id: "n8n-39", seat: "investigator", name: "trace.events", capability: "investigator.trace", webhook: "trace-events" },
  { id: "n8n-40", seat: "investigator", name: "trace.health", capability: "investigator.trace", webhook: "trace-health" },
  { id: "n8n-41", seat: "init", name: "init.self-improve", capability: "init.improve", webhook: "improve-init" },
  { id: "n8n-42", seat: "delivery", name: "delivery.self-improve", capability: "delivery.improve", webhook: "improve-delivery" },
  { id: "n8n-43", seat: "router", name: "router.self-improve", capability: "router.improve", webhook: "improve-router" },
  { id: "n8n-44", seat: "research", name: "research.self-improve", capability: "research.improve", webhook: "improve-research" },
  { id: "n8n-45", seat: "explorer", name: "explorer.self-improve", capability: "explorer.improve", webhook: "improve-explorer" },
  { id: "n8n-46", seat: "operator", name: "operator.self-improve", capability: "operator.improve", webhook: "improve-operator" },
  { id: "n8n-47", seat: "reviewer", name: "reviewer.self-improve", capability: "reviewer.improve", webhook: "improve-reviewer" },
  { id: "n8n-48", seat: "architect", name: "architect.self-improve", capability: "architect.improve", webhook: "improve-architect" },
  { id: "n8n-49", seat: "deploy", name: "deploy.self-improve", capability: "deploy.improve", webhook: "improve-deploy" },
  { id: "n8n-50", seat: "security", name: "security.self-improve", capability: "security.improve", webhook: "improve-security" },
  { id: "n8n-51", seat: "memory", name: "memory.self-improve", capability: "memory.improve", webhook: "improve-memory" },
  { id: "n8n-52", seat: "comms", name: "comms.self-improve", capability: "comms.improve", webhook: "improve-comms" },
  { id: "n8n-53", seat: "planner", name: "planner.self-improve", capability: "planner.improve", webhook: "improve-planner" },
  { id: "n8n-54", seat: "investigator", name: "investigator.self-improve", capability: "investigator.improve", webhook: "improve-investigator" },
  { id: "n8n-55", seat: "contrarian", name: "steel.contradict", capability: "contrarian.contradict", webhook: "steel-contradict" },
  { id: "n8n-56", seat: "contrarian", name: "steel.test", capability: "contrarian.test", webhook: "steel-test" },
  { id: "n8n-57", seat: "contrarian", name: "steel.optimize", capability: "contrarian.optimize", webhook: "steel-optimize" },
  { id: "n8n-58", seat: "contrarian", name: "steel.challenge", capability: "contrarian.challenge", webhook: "steel-challenge" },
];

export function loadN8nFleet(start = process.cwd()): N8nSpecialist[] {
  const paths = [
    path.resolve(start, "state/n8n-fleet.json"),
    path.resolve(start, "../../state/n8n-fleet.json"),
    path.resolve(start, "../../../state/n8n-fleet.json"),
  ];
  for (const file of paths) {
    try {
      const raw = JSON.parse(readFileSync(file, "utf8")) as { specialists?: N8nSpecialist[] };
      if (raw.specialists?.length >= 40) return raw.specialists;
    } catch {
      // try next
    }
  }
  return EMBEDDED_N8N_FLEET;
}

/**
 * Fabric router. Wide match, never a substring of "plane" as "plan".
 * Explicit webhook: / n8n-NN win. Else longest webhook/name/capability hit.
 */
export function resolveSpecialist(task: { title: string; body: string }, fleet: N8nSpecialist[] = loadN8nFleet()): N8nSpecialist | undefined {
  if (!fleet.length) return undefined;
  const blob = `${task.title} ${task.body}`.toLowerCase();
  const hyphen = blob.replace(/[._]+/g, "-").replace(/\s+/g, "-");
  const hook = blob.match(/webhook:([a-z0-9-]+)/);
  if (hook?.[1]) return fleet.find((s) => s.webhook === hook[1]) ?? fleet[0];
  const id = blob.match(/\bn8n-(\d{2})\b/);
  if (id) return fleet.find((s) => s.id === `n8n-${id[1]}`) ?? fleet[0];

  let best: N8nSpecialist | undefined;
  let score = 0;
  for (const spec of fleet) {
    const keys = [spec.webhook, spec.name, spec.capability]
      .filter(Boolean)
      .map((k) => String(k).toLowerCase().replace(/[._]+/g, "-"));
    for (const key of keys) {
      if (key.length < 6) continue;
      if ((hyphen.includes(key) || blob.includes(key.replace(/-/g, " "))) && key.length > score) {
        best = spec;
        score = key.length;
      }
    }
  }
  return best ?? fleet[0];
}

export function webhookFromTask(task: { title: string; body: string }, fleet?: N8nSpecialist[]): string {
  return resolveSpecialist(task, fleet ?? loadN8nFleet())?.webhook ?? "os-pulse";
}

/**
 * n8n fabric. Same seats, wider verbs. Unbound (no N8N_BASE_URL) falls
 * back to the Voltron HTTP adapter so the fleet still runs tonight.
 * A bound webhook that errors also falls back — fabric, not a hard cut.
 */
export function createN8nAdapter(opts: N8nOpts = {}): AgentAdapter {
  const fetchImpl = opts.fetchImpl ?? fetch;
  const fallback = opts.fallback ?? createHttpAdapter({ fetchImpl: opts.fetchImpl });
  const fleet = opts.fleet ?? loadN8nFleet();
  const last = new Map<string, AdapterRunStatus>();

  function base() {
    return (opts.baseUrl || process.env.N8N_BASE_URL || "").replace(/\/$/, "");
  }

  async function viaFallback(task: Task, ctx: RunContext, reason: string): Promise<RunHandle> {
    const handle = await fallback.dispatch(task, ctx);
    const status = await fallback.poll(handle);
    last.set(ctx.runId, {
      ...status,
      exitReason: `${reason}:${status.exitReason}`,
    });
    return { id: ctx.runId, adapter: "n8n", depth: ctx.depth };
  }

  return {
    id: "n8n",

    async healthcheck(): Promise<Health> {
      const url = base();
      if (!url) return { ok: true, detail: "n8n unbound — falling back to Voltron http" };
      try {
        const res = await fetchImpl(`${url}/healthz`);
        return { ok: res.ok, detail: `n8n ${res.status}` };
      } catch (err) {
        return { ok: false, detail: err instanceof Error ? err.message : "n8n unreachable" };
      }
    },

    async estimateCost() {
      return 0;
    },

    async dispatch(task: Task, ctx: RunContext): Promise<RunHandle> {
      assertDepth(ctx.depth);
      const url = base();
      const webhook = webhookFromTask(task, fleet);
      if (url) {
        try {
          const res = await fetchImpl(`${url}/webhook/${encodeURIComponent(webhook)}`, {
            method: "POST",
            headers: {
              "content-type": "application/json",
              ...(opts.apiKey || process.env.N8N_API_KEY ? { "x-n8n-api-key": opts.apiKey || process.env.N8N_API_KEY || "" } : {}),
            },
            body: JSON.stringify({ task, runId: ctx.runId, depth: ctx.depth, webhook, fabric: "n8n" }),
          });
          const text = await res.text();
          if (res.ok) {
            last.set(ctx.runId, {
              status: "succeeded",
              inputTokens: 0,
              outputTokens: 0,
              costUsd: 0,
              exitReason: `n8n:${webhook}`,
              artifact: { type: "report", path: `runs/${ctx.runId}.json`, sha: ctx.runId, body: text },
            });
            return { id: ctx.runId, adapter: "n8n", depth: ctx.depth };
          }
        } catch {
          // bound n8n missed — Voltron still has the seat
        }
        return viaFallback(task, ctx, `n8n-fallback:${webhook}`);
      }
      return viaFallback(task, ctx, "n8n-fallback");
    },

    async poll(handle: RunHandle) {
      return (
        last.get(handle.id) ?? {
          status: "failed",
          inputTokens: 0,
          outputTokens: 0,
          costUsd: 0,
          exitReason: "unknown-handle",
        }
      );
    },

    async cancel() {},
  };
}
