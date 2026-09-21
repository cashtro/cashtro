import type { Task } from "@cashtro/sdk";

export const MAX_DEPTH = 3;

export type Health = { ok: boolean; detail: string };

export type RunHandle = { id: string; adapter: string; depth: number };

export type AdapterRunStatus = {
  status: "queued" | "running" | "succeeded" | "failed" | "cancelled";
  inputTokens: number;
  outputTokens: number;
  costUsd: number;
  exitReason: string;
  artifact?: { type: string; path: string; sha: string; body?: string };
};

export type RunContext = {
  runId: string;
  depth: number;
  budgetCapUsd: number;
  kernelUrl: string;
};

export interface AgentAdapter {
  id: string;
  healthcheck(): Promise<Health>;
  dispatch(task: Task, ctx: RunContext): Promise<RunHandle>;
  poll(handle: RunHandle): Promise<AdapterRunStatus>;
  cancel(handle: RunHandle): Promise<void>;
  estimateCost(task: Task): Promise<number>;
}

export function assertDepth(depth: number): void {
  if (!Number.isFinite(depth) || depth < 1) {
    throw new DepthError("depth must be >= 1");
  }
  if (depth > MAX_DEPTH) {
    throw new DepthError(`depth ${depth} exceeds hard stop ${MAX_DEPTH}`);
  }
}

export class DepthError extends Error {
  readonly code = "depth-limit";
  constructor(message: string) {
    super(message);
    this.name = "DepthError";
  }
}
