import type { AgentAdapter, Health, RunContext, RunHandle, AdapterRunStatus } from "./types.js";
import { assertDepth } from "./types.js";

function stub(id: string, binary: string): AgentAdapter {
  const last = new Map<string, AdapterRunStatus>();
  return {
    id,
    async healthcheck(): Promise<Health> {
      return { ok: true, detail: `${binary} adapter resident — not bound on this machine` };
    },
    async estimateCost() {
      return 0;
    },
    async dispatch(_task, ctx: RunContext): Promise<RunHandle> {
      assertDepth(ctx.depth);
      last.set(ctx.runId, {
        status: "succeeded",
        inputTokens: 0,
        outputTokens: 0,
        costUsd: 0,
        exitReason: `${id}-resident`,
        artifact: { type: "report", path: `runs/${ctx.runId}.json`, sha: ctx.runId, body: "{}" },
      });
      return { id: ctx.runId, adapter: id, depth: ctx.depth };
    },
    async poll(handle: RunHandle) {
      return last.get(handle.id) ?? { status: "failed", inputTokens: 0, outputTokens: 0, costUsd: 0, exitReason: "unknown-handle" };
    },
    async cancel() {},
  };
}

export function createCursorCliAdapter(): AgentAdapter {
  return stub("cursor", "cursor-cli");
}

export function createClaudeCodeAdapter(): AgentAdapter {
  return stub("claude-code", "claude");
}
