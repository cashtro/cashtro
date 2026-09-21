import Fastify, { type FastifyInstance, type FastifyRequest } from "fastify";
import swagger from "@fastify/swagger";
import swaggerUi from "@fastify/swagger-ui";
import { PrismaClient } from "@prisma/client";
import { CreateTask } from "@cashtro/sdk";
import type { AgentAdapter } from "@cashtro/adapters";
import { drainQueue, orchestrate } from "./orchestrate.js";
import { ingestGithubWebhook } from "./githubWebhook.js";
import { renderPane } from "./pane.js";

export type AppOpts = {
  prisma: PrismaClient;
  apiKeys?: Array<{ id: string; secret: string; scope: string }>;
  budgetCap?: number;
  kernelUrl?: string;
  fetchImpl?: typeof fetch;
  adapters?: Record<string, AgentAdapter>;
};

function parseKeys(raw: string | undefined): Array<{ id: string; secret: string; scope: string }> {
  return (raw || "local:change-me:admin")
    .split(",")
    .map((row) => row.trim())
    .filter(Boolean)
    .map((row) => {
      const [id, secret, scope] = row.split(":");
      return { id: id || "anon", secret: secret || "", scope: scope || "read" };
    });
}

function jsonArr(value: string | null | undefined): string[] {
  try {
    const v = JSON.parse(value || "[]");
    return Array.isArray(v) ? v.map(String) : [];
  } catch {
    return [];
  }
}

export async function buildApp(opts: AppOpts): Promise<FastifyInstance> {
  const prisma = opts.prisma;
  const keys = opts.apiKeys ?? parseKeys(process.env.CONTROL_PLANE_API_KEYS);
  const budgetCap = opts.budgetCap ?? Number(process.env.HARD_BUDGET_PER_RUN_USD || 2);

  const app = Fastify({ logger: false });
  await app.register(swagger, {
    openapi: {
      openapi: "3.1.0",
      info: { title: "Cashtro Control Plane", version: "0.3.0" },
    },
  });
  await app.register(swaggerUi, { routePrefix: "/docs" });

  app.decorateRequest("actor", "");
  app.addHook("preHandler", async (req, reply) => {
    if (
      req.url === "/docs" ||
      req.url.startsWith("/docs/") ||
      req.url.startsWith("/documentation") ||
      req.url === "/metrics" ||
      req.url === "/health" ||
      req.url === "/ui" ||
      req.url.startsWith("/ui/")
    ) {
      req.actor = "public";
      return;
    }
    const header = req.headers.authorization || "";
    const token = header.startsWith("Bearer ") ? header.slice(7) : "";
    const match = keys.find((k) => k.secret && k.secret === token);
    if (!match) {
      return reply.code(401).send({ error: "unauthorized" });
    }
    if (req.method !== "GET" && match.scope === "read") {
      return reply.code(403).send({ error: "read-only key" });
    }
    req.actor = match.id;
  });

  async function emit(actor: string, type: string, payload: unknown) {
    await prisma.event.create({ data: { actor, type, payload: JSON.stringify(payload) } });
  }

  app.get("/health", async () => ({ status: "ok", service: "control-plane" }));
  app.get("/ui", async (req, reply) => {
    reply.header("content-type", "text/html; charset=utf-8");
    return renderPane(prisma, "fleet");
  });
  app.get("/ui/:screen", async (req, reply) => {
    const { screen } = req.params as { screen: string };
    reply.header("content-type", "text/html; charset=utf-8");
    return renderPane(prisma, screen);
  });

  app.post("/ui/act/dispatch/:id", async (req, reply) => {
    const { id } = req.params as { id: string };
    const result = await orchestrate({
      prisma,
      taskId: id,
      actor: req.actor || "pane",
      budgetCap,
      kernelUrl: opts.kernelUrl,
      fetchImpl: opts.fetchImpl,
      adapters: opts.adapters,
    });
    if (!result.ok) return reply.code(result.status).send(result);
    return reply.redirect("/ui/run");
  });

  app.post("/ui/act/pause", async (req, reply) => {
    await prisma.controlState.upsert({
      where: { id: "global" },
      update: { paused: true },
      create: { id: "global", paused: true },
    });
    await drainQueue(prisma, req.actor || "pane");
    return reply.redirect("/ui/fleet");
  });

  app.get("/projects", async (req) => {
    const q = req.query as { kind?: string; status?: string; tag?: string };
    const rows = await prisma.project.findMany();
    return rows
      .filter((p) => (!q.kind || p.kind === q.kind) && (!q.status || p.status === q.status) && (!q.tag || jsonArr(p.tags).includes(q.tag)))
      .map((p) => ({ ...p, tags: jsonArr(p.tags) }));
  });

  app.get("/projects/:id", async (req, reply) => {
    const { id } = req.params as { id: string };
    const project = await prisma.project.findFirst({
      where: { OR: [{ id }, { slug: id }] },
      include: { capabilities: true, agents: true },
    });
    if (!project) return reply.code(404).send({ error: "not found" });
    return { ...project, tags: jsonArr(project.tags) };
  });

  app.post("/projects/:id/sync", async (req, reply) => {
    const { id } = req.params as { id: string };
    const project = await prisma.project.findFirst({ where: { OR: [{ id }, { slug: id }] } });
    if (!project) return reply.code(404).send({ error: "not found" });
    await emit(req.actor, "project.sync", { id: project.id, slug: project.slug });
    return { ok: true, queued: "recon", slug: project.slug };
  });

  app.get("/agents", async () => prisma.agent.findMany({ include: { workers: true } }));

  app.get("/fleet", async () => {
    const agents = await prisma.agent.findMany({ include: { workers: true } });
    const specialists = agents.flatMap((a) =>
      a.workers
        .filter((w) => w.tool.startsWith("n8n:"))
        .map((w) => ({ id: w.name, seat: a.name, webhook: w.tool.slice(4), tool: w.tool })),
    );
    const bySeat = Object.fromEntries(
      agents.map((a) => [a.name, specialists.filter((s) => s.seat === a.name).map((s) => s.id)]),
    );
    return {
      architecture: "wide-not-deep",
      fabric: "n8n-inside-voltron",
      seats: agents.length,
      specialists: specialists.length,
      depthLimit: 3,
      n8n: process.env.N8N_BASE_URL ? "bound" : "unbound-fallback-voltron",
      bySeat,
      agents: agents.map((a) => ({
        name: a.name,
        runtime: a.runtime,
        workers: a.workers.map((w) => ({ name: w.name, tool: w.tool })),
      })),
    };
  });

  app.post("/tasks", async (req, reply) => {
    const parsed = CreateTask.safeParse(req.body);
    if (!parsed.success) return reply.code(400).send({ error: parsed.error.flatten() });
    const existing = await prisma.task.findUnique({ where: { idempotencyKey: parsed.data.idempotencyKey } });
    if (existing) return existing;
    const created = await prisma.task.create({
      data: {
        title: parsed.data.title,
        body: parsed.data.body,
        source: parsed.data.source,
        priority: parsed.data.priority,
        projectId: parsed.data.projectId,
        assigneeAgentId: parsed.data.assigneeAgentId,
        dependsOn: JSON.stringify(parsed.data.dependsOn),
        idempotencyKey: parsed.data.idempotencyKey,
      },
    });
    await emit(req.actor, "task.create", { id: created.id });
    return reply.code(201).send(created);
  });

  app.get("/tasks", async (req) => {
    const q = req.query as { status?: string; project?: string };
    return prisma.task.findMany({
      where: {
        ...(q.status ? { status: q.status } : {}),
        ...(q.project ? { project: { is: { OR: [{ id: q.project }, { slug: q.project }] } } } : {}),
      },
    });
  });

  app.post("/tasks/:id/dispatch", async (req, reply) => {
    const { id } = req.params as { id: string };
    const body = (req.body || {}) as { depth?: number };
    const result = await orchestrate({
      prisma,
      taskId: id,
      actor: req.actor,
      budgetCap,
      depth: body.depth,
      kernelUrl: opts.kernelUrl,
      fetchImpl: opts.fetchImpl,
      adapters: opts.adapters,
    });
    if (!result.ok) return reply.code(result.status).send(result);
    return result.run;
  });

  app.post("/tasks/:id/cancel", async (req, reply) => {
    const { id } = req.params as { id: string };
    const task = await prisma.task.findUnique({ where: { id } });
    if (!task) return reply.code(404).send({ error: "not found" });
    const updated = await prisma.task.update({ where: { id }, data: { status: "cancelled" } });
    await prisma.run.updateMany({ where: { taskId: id, status: { in: ["queued", "running"] } }, data: { status: "cancelled", endedAt: new Date(), exitReason: "cancelled" } });
    await emit(req.actor, "task.cancel", { id });
    return updated;
  });

  app.post("/control/pause", async (req) => {
    const state = await prisma.controlState.upsert({
      where: { id: "global" },
      update: { paused: true },
      create: { id: "global", paused: true },
    });
    const drained = await drainQueue(prisma, req.actor);
    await emit(req.actor, "control.pause", { paused: true, drained });
    return { ...state, drained };
  });

  app.get("/control", async () => {
    const state = await prisma.controlState.findUnique({ where: { id: "global" } });
    const queued = await prisma.job.count({ where: { status: "queued" } });
    return { paused: state?.paused ?? false, queued };
  });

  app.get("/costs", async () => {
    const since = new Date(Date.now() - 7 * 24 * 3600 * 1000);
    const runs = await prisma.run.findMany({ where: { startedAt: { gte: since } } });
    const total = runs.reduce((n, r) => n + r.costUsd, 0);
    const byAgent: Record<string, number> = {};
    for (const r of runs) byAgent[r.agentId] = (byAgent[r.agentId] || 0) + r.costUsd;
    return { windowDays: 7, runs: runs.length, totalUsd: total, byAgent, hardCapPerRun: budgetCap };
  });

  app.post("/control/resume", async (req) => {
    const state = await prisma.controlState.upsert({
      where: { id: "global" },
      update: { paused: false },
      create: { id: "global", paused: false },
    });
    await emit(req.actor, "control.resume", { paused: false });
    return state;
  });

  app.get("/runs/:id", async (req, reply) => {
    const { id } = req.params as { id: string };
    const run = await prisma.run.findUnique({ where: { id }, include: { artifacts: true } });
    if (!run) return reply.code(404).send({ error: "not found" });
    return run;
  });

  app.get("/runs/:id/stream", async (req, reply) => {
    const { id } = req.params as { id: string };
    const run = await prisma.run.findUnique({ where: { id } });
    if (!run) return reply.code(404).send({ error: "not found" });
    reply.header("Content-Type", "text/event-stream");
    return `data: ${JSON.stringify({ id: run.id, status: run.status, costUsd: run.costUsd })}\n\n`;
  });

  app.get("/metrics", async () => {
    const [runs, cost, paused] = await Promise.all([
      prisma.run.count(),
      prisma.run.aggregate({ _sum: { costUsd: true } }),
      prisma.controlState.findUnique({ where: { id: "global" } }),
    ]);
    const pause = paused?.paused ? 1 : 0;
    return [
      `# HELP cashtro_runs_total Completed and in-flight runs`,
      `# TYPE cashtro_runs_total counter`,
      `cashtro_runs_total ${runs}`,
      `# HELP cashtro_cost_usd_total Spend recorded on runs`,
      `# TYPE cashtro_cost_usd_total counter`,
      `cashtro_cost_usd_total ${cost._sum.costUsd ?? 0}`,
      `# HELP cashtro_paused Kill switch`,
      `# TYPE cashtro_paused gauge`,
      `cashtro_paused ${pause}`,
      "",
    ].join("\n");
  });

  app.post("/webhooks/github", async (req) => {
    await emit(req.actor || "github", "webhook.github", req.body ?? {});
    const ingested = ingestGithubWebhook(req.body);
    if (!ingested) return { ok: true, task: null };
    const existing = await prisma.task.findUnique({ where: { idempotencyKey: ingested.idempotencyKey } });
    if (existing) return { ok: true, task: existing, deduped: true };
    const task = await prisma.task.create({
      data: {
        title: ingested.title,
        body: ingested.body,
        source: ingested.source,
        priority: "p1",
        idempotencyKey: ingested.idempotencyKey,
      },
    });
    await emit(req.actor || "github", "task.create", { id: task.id, source: "github" });
    return { ok: true, task };
  });

  app.post("/webhooks/:agent/callback", async (req) => {
    const { agent } = req.params as { agent: string };
    await emit(req.actor || agent, "webhook.agent", { agent, body: req.body ?? {} });
    return { ok: true };
  });

  return app;
}

declare module "fastify" {
  interface FastifyRequest {
    actor: string;
  }
}

void (0 as unknown as FastifyRequest);
