import type { Task } from "@cashtro/sdk";
import type { AdapterRunStatus, AgentAdapter, Health, RunContext, RunHandle } from "./types.js";
import { assertDepth } from "./types.js";

export type HttpAdapterOpts = {
  kernelUrl?: string;
  fetchImpl?: typeof fetch;
};

type Last = AdapterRunStatus & { handleId: string };

/** Voltron / Cashtro OS kernel over HTTP. Zero model spend. */
export function createHttpAdapter(opts: HttpAdapterOpts = {}): AgentAdapter {
  const fetchImpl = opts.fetchImpl ?? fetch;
  const last = new Map<string, Last>();

  async function kernelUrl(ctx?: RunContext): Promise<string> {
    return (ctx?.kernelUrl || opts.kernelUrl || process.env.VOLTRON_URL || "http://127.0.0.1:8080").replace(/\/$/, "");
  }

  return {
    id: "http",

    async healthcheck(): Promise<Health> {
      const base = await kernelUrl();
      try {
        const res = await fetchImpl(`${base}/health`);
        if (!res.ok) return { ok: false, detail: `http ${res.status}` };
        const body = (await res.json()) as { running?: number; status?: string };
        return { ok: body.status === "ok", detail: `${base} running=${body.running ?? "?"}` };
      } catch (err) {
        return { ok: false, detail: err instanceof Error ? err.message : "unreachable" };
      }
    },

    async estimateCost() {
      return 0;
    },

    async dispatch(task: Task, ctx: RunContext): Promise<RunHandle> {
      assertDepth(ctx.depth);
      const base = await kernelUrl(ctx);
      const agent = task.assigneeAgentId ? undefined : undefined;
      const agentName = pickKernelAgent(task);
      const capability = pickCapability(task, agentName);
      const res = await fetchImpl(`${base}/api/agents/${encodeURIComponent(agentName)}/invoke`, {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ capability, payload: { title: task.title, body: task.body, taskId: task.id } }),
      });
      const text = await res.text();
      let parsed: { ok?: boolean; message?: string } = {};
      try {
        parsed = JSON.parse(text) as { ok?: boolean; message?: string };
      } catch {
        parsed = { ok: false, message: text.slice(0, 400) };
      }
      const handle: RunHandle = { id: ctx.runId, adapter: "http", depth: ctx.depth };
      const ok = res.ok && parsed.ok !== false;
      last.set(handle.id, {
        handleId: handle.id,
        status: ok ? "succeeded" : "failed",
        inputTokens: 0,
        outputTokens: 0,
        costUsd: 0,
        exitReason: ok ? `http:${agentName}.${capability}` : `http-error:${res.status}`,
        artifact: {
          type: "report",
          path: `runs/${ctx.runId}.json`,
          sha: ctx.runId,
          body: text,
        },
      });
      void agent;
      return handle;
    },

    async poll(handle: RunHandle): Promise<AdapterRunStatus> {
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

    async cancel(handle: RunHandle) {
      const prev = last.get(handle.id);
      if (prev && prev.status === "running") {
        last.set(handle.id, { ...prev, status: "cancelled", exitReason: "cancelled" });
      }
    },
  };
}

function pickKernelAgent(task: Task): string {
  const blob = `${task.title} ${task.body}`.toLowerCase();
  if (/\bexplor/.test(blob)) return "explorer";
  if (/\bdeliver|\bship\b/.test(blob)) return "delivery";
  if (/\bresearch\b/.test(blob)) return "research";
  if (/\bplanner\b|\bbacklog\b/.test(blob)) return "planner";
  if (/\bmemory\b/.test(blob)) return "memory";
  return "init";
}

function pickCapability(task: Task, agent: string): string {
  if (agent === "explorer") return "explorer.search";
  if (agent === "delivery") return "delivery.list";
  if (agent === "research") return "research.list";
  if (agent === "planner") return "planner.backlog";
  if (agent === "memory") return "memory.recall";
  return "os.about";
}
