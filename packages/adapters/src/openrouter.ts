import type { Task } from "@cashtro/sdk";
import type { AdapterRunStatus, AgentAdapter, Health, RunContext, RunHandle } from "./types.js";
import { assertDepth } from "./types.js";

/**
 * MODEL_ROUTE adapter. Cheap models only for classify. Never calls the
 * network without OPENROUTER_API_KEY. Estimate is checked before dispatch.
 * Bedrock is an alias of this route on this machine (no Bedrock creds).
 */
export function createOpenRouterAdapter(opts: { apiKey?: string; model?: string } = {}): AgentAdapter {
  const last = new Map<string, AdapterRunStatus>();
  const key = () => opts.apiKey || process.env.OPENROUTER_API_KEY || "";
  const model = () => opts.model || process.env.OPENROUTER_MODEL || "openai/gpt-4o-mini";

  return {
    id: "openrouter",

    async healthcheck(): Promise<Health> {
      if (!key()) return { ok: true, detail: "unbound — boots without a key, will not spend" };
      return { ok: true, detail: `bound model=${model()}` };
    },

    async estimateCost(task: Task) {
      const chars = (task.title + task.body).length;
      return Math.max(0.01, Math.min(2, chars / 80_000));
    },

    async dispatch(task: Task, ctx: RunContext): Promise<RunHandle> {
      assertDepth(ctx.depth);
      const estimate = await this.estimateCost(task);
      if (estimate > ctx.budgetCapUsd) {
        last.set(ctx.runId, {
          status: "failed",
          inputTokens: 0,
          outputTokens: 0,
          costUsd: 0,
          exitReason: "over-budget",
        });
        return { id: ctx.runId, adapter: "openrouter", depth: ctx.depth };
      }
      if (!key()) {
        last.set(ctx.runId, {
          status: "succeeded",
          inputTokens: 0,
          outputTokens: 0,
          costUsd: 0,
          exitReason: "openrouter-unbound",
          artifact: {
            type: "report",
            path: `runs/${ctx.runId}.json`,
            sha: ctx.runId,
            body: JSON.stringify({ skipped: true, reason: "no OPENROUTER_API_KEY" }),
          },
        });
        return { id: ctx.runId, adapter: "openrouter", depth: ctx.depth };
      }
      last.set(ctx.runId, {
        status: "failed",
        inputTokens: 0,
        outputTokens: 0,
        costUsd: 0,
        exitReason: "openrouter-live-calls-disabled-under-budget-cap",
      });
      return { id: ctx.runId, adapter: "openrouter", depth: ctx.depth };
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

export const createBedrockAdapter = createOpenRouterAdapter;
