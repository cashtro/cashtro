import type { PrismaClient } from "@prisma/client";
import type { AgentAdapter, RunContext } from "@cashtro/adapters";
import { DepthError, MAX_DEPTH, assertDepth, createContrarianAdapter, createHttpAdapter, createN8nAdapter, defaultAdapters, pickAdapter, resolveSpecialist } from "@cashtro/adapters";
import { Task as TaskSchema, type Task } from "@cashtro/sdk";

export type OrchestrateOpts = {
  prisma: PrismaClient;
  taskId: string;
  actor: string;
  budgetCap: number;
  depth?: number;
  kernelUrl?: string;
  fetchImpl?: typeof fetch;
  adapters?: Record<string, AgentAdapter>;
};

export async function orchestrate(opts: OrchestrateOpts) {
  const depth = opts.depth ?? 1;
  try {
    assertDepth(depth);
  } catch (err) {
    if (err instanceof DepthError) {
      return { ok: false as const, status: 409, error: err.code, message: err.message };
    }
    throw err;
  }

  const control = await opts.prisma.controlState.findUnique({ where: { id: "global" } });
  if (control?.paused) {
    return { ok: false as const, status: 409, error: "paused" };
  }

  const row = await opts.prisma.task.findUnique({ where: { id: opts.taskId } });
  if (!row) return { ok: false as const, status: 404, error: "not found" };

  const seats = await opts.prisma.agent.findMany({ where: { enabled: true }, include: { workers: true } });
  const assigned = row.assigneeAgentId ? seats.find((s) => s.id === row.assigneeAgentId) : undefined;
  const fabric = seats.flatMap((seat) =>
    seat.workers
      .filter((w) => w.tool.startsWith("n8n:"))
      .map((w) => ({ id: w.name, seat: seat.name, webhook: w.tool.slice(4), name: w.name })),
  );
  const scoped = assigned ? fabric.filter((s) => s.seat === assigned.name) : fabric;
  const specialist = resolveSpecialist({ title: row.title, body: row.body }, scoped.length ? scoped : fabric);

  const agent =
    assigned ||
    (specialist && seats.find((s) => s.name === specialist.seat)) ||
    seats.find((s) => s.runtime === "http") ||
    seats[0];
  if (!agent) return { ok: false as const, status: 409, error: "no agent" };

  const task = TaskSchema.parse({
    id: row.id,
    projectId: row.projectId,
    title: row.title,
    body: specialist ? `${row.body}\nwebhook:${specialist.webhook}`.trim() : row.body,
    source: row.source,
    priority: row.priority,
    status: row.status,
    dependsOn: safeArr(row.dependsOn),
    assigneeAgentId: row.assigneeAgentId,
    idempotencyKey: row.idempotencyKey,
  });

  const adapters = opts.adapters ?? defaultAdapters();
  if (opts.kernelUrl || opts.fetchImpl) {
    adapters.http = createHttpAdapter({ kernelUrl: opts.kernelUrl, fetchImpl: opts.fetchImpl });
    adapters.n8n = createN8nAdapter({ fallback: adapters.http, fetchImpl: opts.fetchImpl });
  }
  if (!adapters.contrarian) adapters.contrarian = createContrarianAdapter();
  const steelSeat = agent.name === "contrarian" || agent.runtime === "contrarian" || specialist?.seat === "contrarian";
  let adapter = steelSeat
    ? adapters.contrarian
    : specialist && adapters.n8n
      ? adapters.n8n
      : pickAdapter(agent.runtime, adapters);
  if (adapter.id === "http") {
    const health = await adapter.healthcheck();
    if (!health.ok && adapters.local) adapter = adapters.local;
  } else if (adapter.id === "n8n" && adapters.http && adapters.local) {
    const health = await adapters.http.healthcheck();
    if (!health.ok) {
      adapter = createN8nAdapter({ fallback: adapters.local, baseUrl: process.env.N8N_BASE_URL || "" });
    }
  }
  const estimate = await adapter.estimateCost(task);
  const cap = Math.min(opts.budgetCap, agent.maxCostPerRun);
  if (estimate > cap) {
    const blocked = await opts.prisma.run.create({
      data: { taskId: row.id, agentId: agent.id, status: "blocked-budget", costUsd: 0, exitReason: "over-budget" },
    });
    await opts.prisma.event.create({
      data: { actor: opts.actor, type: "run.blocked-budget", payload: JSON.stringify({ runId: blocked.id }) },
    });
    return { ok: false as const, status: 409, error: "over-budget", run: blocked };
  }

  const job = await opts.prisma.job.create({ data: { taskId: row.id, status: "queued", depth } });
  const run = await opts.prisma.run.create({
    data: { taskId: row.id, agentId: agent.id, status: "running", costUsd: 0, exitReason: "dispatching" },
  });
  await opts.prisma.job.update({ where: { id: job.id }, data: { runId: run.id, status: "running" } });

  const ctx: RunContext = {
    runId: run.id,
    depth,
    budgetCapUsd: cap,
    kernelUrl: opts.kernelUrl || process.env.VOLTRON_URL || "http://127.0.0.1:8080",
  };

  let status;
  try {
    const handle = await adapter.dispatch(task, ctx);
    status = await adapter.poll(handle);
  } catch (err) {
    if (err instanceof DepthError) {
      await opts.prisma.run.update({
        where: { id: run.id },
        data: { status: "failed", endedAt: new Date(), exitReason: "depth-limit", costUsd: 0 },
      });
      await opts.prisma.job.update({ where: { id: job.id }, data: { status: "cancelled" } });
      return { ok: false as const, status: 409, error: err.code, message: err.message };
    }
    status = {
      status: "failed" as const,
      inputTokens: 0,
      outputTokens: 0,
      costUsd: 0,
      exitReason: err instanceof Error ? err.message : "adapter-error",
    };
  }

  const finished = await opts.prisma.run.update({
    where: { id: run.id },
    data: {
      status: status.status === "succeeded" ? "succeeded" : status.status === "cancelled" ? "cancelled" : "failed",
      endedAt: new Date(),
      inputTokens: status.inputTokens,
      outputTokens: status.outputTokens,
      costUsd: status.costUsd,
      exitReason: status.exitReason,
    },
    include: { artifacts: true },
  });

  if (status.artifact) {
    await opts.prisma.artifact.create({
      data: {
        runId: run.id,
        type: status.artifact.type,
        path: status.artifact.path,
        sha: status.artifact.sha,
      },
    });
  }

  await opts.prisma.job.update({
    where: { id: job.id },
    data: { status: status.status === "succeeded" ? "done" : status.status === "cancelled" ? "cancelled" : "failed" },
  });
  await opts.prisma.task.update({
    where: { id: row.id },
    data: { status: status.status === "succeeded" ? "done" : row.status, assigneeAgentId: agent.id },
  });
  await opts.prisma.event.create({
    data: {
      actor: opts.actor,
      type: "task.dispatch",
      payload: JSON.stringify({
        taskId: row.id,
        runId: run.id,
        adapter: adapter.id,
        costUsd: status.costUsd,
        depth,
        maxDepth: MAX_DEPTH,
      }),
    },
  });

  const withArt = await opts.prisma.run.findUnique({ where: { id: run.id }, include: { artifacts: true } });
  return { ok: true as const, run: withArt ?? finished, adapter: adapter.id };
}

export async function drainQueue(prisma: PrismaClient, actor: string) {
  const drained = await prisma.job.updateMany({
    where: { status: { in: ["queued"] } },
    data: { status: "cancelled" },
  });
  await prisma.event.create({
    data: { actor, type: "control.drain", payload: JSON.stringify({ cancelled: drained.count }) },
  });
  return drained.count;
}

function safeArr(raw: string): string[] {
  try {
    const v = JSON.parse(raw);
    return Array.isArray(v) ? v.map(String) : [];
  } catch {
    return [];
  }
}

void (0 as unknown as Task);
