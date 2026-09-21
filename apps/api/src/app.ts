import Fastify, { type FastifyInstance, type FastifyRequest } from "fastify";
import swagger from "@fastify/swagger";
import swaggerUi from "@fastify/swagger-ui";
import { PrismaClient } from "@prisma/client";
import { CreateTask } from "@cashtro/sdk";

export type AppOpts = {
  prisma: PrismaClient;
  apiKeys?: Array<{ id: string; secret: string; scope: string }>;
  budgetCap?: number;
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
    if (req.method === "GET" && (req.url === "/docs" || req.url.startsWith("/docs/") || req.url.startsWith("/documentation") || req.url === "/metrics" || req.url === "/health")) {
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

  app.get("/agents", async () => prisma.agent.findMany());

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
    const control = await prisma.controlState.findUnique({ where: { id: "global" } });
    if (control?.paused) return reply.code(409).send({ error: "paused" });
    const { id } = req.params as { id: string };
    const task = await prisma.task.findUnique({ where: { id } });
    if (!task) return reply.code(404).send({ error: "not found" });
    const agent =
      (task.assigneeAgentId && (await prisma.agent.findUnique({ where: { id: task.assigneeAgentId } }))) ||
      (await prisma.agent.findFirst({ where: { enabled: true } }));
    if (!agent) return reply.code(409).send({ error: "no agent" });
    const estimate = 0.01;
    if (estimate > Math.min(budgetCap, agent.maxCostPerRun)) {
      const blocked = await prisma.run.create({
        data: { taskId: task.id, agentId: agent.id, status: "blocked-budget", costUsd: 0, exitReason: "over-budget" },
      });
      await emit(req.actor, "run.blocked-budget", { runId: blocked.id });
      return reply.code(409).send({ error: "over-budget", run: blocked });
    }
    const run = await prisma.run.create({
      data: {
        taskId: task.id,
        agentId: agent.id,
        status: "succeeded",
        endedAt: new Date(),
        inputTokens: 0,
        outputTokens: 0,
        costUsd: estimate,
        exitReason: "local-queue",
      },
    });
    await prisma.artifact.create({ data: { runId: run.id, type: "report", path: `runs/${run.id}.json`, sha: run.id } });
    await prisma.task.update({ where: { id: task.id }, data: { status: "done", assigneeAgentId: agent.id } });
    await emit(req.actor, "task.dispatch", { taskId: task.id, runId: run.id, costUsd: estimate });
    return run;
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
    await emit(req.actor, "control.pause", { paused: true });
    return state;
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
    return { ok: true };
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
