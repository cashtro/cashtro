import type { AgentAdapter, AdapterRunStatus, RunContext, RunHandle } from "./types.js";
import { assertDepth } from "./types.js";

/** In-process fallback when Voltron is not reachable. Still records a run. */
export function createLocalAdapter(costUsd = 0.01): AgentAdapter {
  const last = new Map<string, AdapterRunStatus>();
  return {
    id: "local",
    async healthcheck() {
      return { ok: true, detail: "local fallback" };
    },
    async estimateCost() {
      return costUsd;
    },
    async dispatch(_task, ctx: RunContext): Promise<RunHandle> {
      assertDepth(ctx.depth);
      last.set(ctx.runId, {
        status: "succeeded",
        inputTokens: 0,
        outputTokens: 0,
        costUsd,
        exitReason: "local-queue",
        artifact: { type: "report", path: `runs/${ctx.runId}.json`, sha: ctx.runId, body: "{}" },
      });
      return { id: ctx.runId, adapter: "local", depth: ctx.depth };
    },
    async poll(handle: RunHandle) {
      return last.get(handle.id) ?? { status: "failed", inputTokens: 0, outputTokens: 0, costUsd: 0, exitReason: "unknown-handle" };
    },
    async cancel() {},
  };
}
