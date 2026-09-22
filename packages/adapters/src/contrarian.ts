import type { Task } from "@cashtro/sdk";
import type { AdapterRunStatus, AgentAdapter, RunContext, RunHandle } from "./types.js";
import { assertDepth } from "./types.js";

export type SteelChallenge = {
  contradiction: string;
  moreOptimized: string;
  tests: string[];
  liveClient: boolean;
};

const LIVE = ["btk", "md clinic", "hypothe", "educonnexion", "proximity"];

export function challenge(input: { title: string; body: string; opponent?: string; shortest?: string }): SteelChallenge {
  const blob = `${input.title} ${input.body} ${input.opponent ?? ""}`.toLowerCase();
  const liveClient = LIVE.some((t) => blob.includes(t));
  if (liveClient) {
    return {
      contradiction: "Live client is in scope. The optimized option is to stop.",
      moreOptimized: "Do not touch live client production. Catalog only.",
      tests: ["refuse-prod", "no-client-url"],
      liveClient: true,
    };
  }
  const shortest = (input.shortest || "").trim();
  const opponent = (input.opponent || "").trim();
  const contradiction =
    opponent ||
    "The first position is usually a dispatch. The cheaper move is to fill the missing context, then a smaller verb.";
  const moreOptimized = shortest
    ? `Keep the shortest move (${shortest}) only if the contradiction still leaves it standing.`
    : "Ask the specific questions, then one local verb, depth 1, $0.";
  return {
    contradiction,
    moreOptimized,
    tests: ["coup-complete", "depth<=3", "budget<=2", "not-live-client"],
    liveClient: false,
  };
}

/** Local Steel department. $0. No kernel process. Contradicts toward the cheaper option. */
export function createContrarianAdapter(): AgentAdapter {
  const last = new Map<string, AdapterRunStatus>();
  return {
    id: "contrarian",
    async healthcheck() {
      return { ok: true, detail: "steel local — contradict toward the cheaper option" };
    },
    async estimateCost() {
      return 0;
    },
    async dispatch(task: Task, ctx: RunContext): Promise<RunHandle> {
      assertDepth(ctx.depth);
      const steel = challenge({ title: task.title, body: task.body });
      const body = JSON.stringify({ department: "contrarian", crew: "Steel", ...steel, task: task.id });
      last.set(ctx.runId, {
        status: steel.liveClient ? "failed" : "succeeded",
        inputTokens: 0,
        outputTokens: 0,
        costUsd: 0,
        exitReason: steel.liveClient ? "steel:refuse-live-client" : "steel:contradict",
        artifact: { type: "report", path: `runs/${ctx.runId}.json`, sha: ctx.runId, body },
      });
      return { id: ctx.runId, adapter: "contrarian", depth: ctx.depth };
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
