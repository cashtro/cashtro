import { readFile } from "node:fs/promises";
import path from "node:path";
import { PrismaClient } from "@prisma/client";

const prisma = new PrismaClient();

const KERNEL_AGENTS = [
  { name: "init", role: "kernel", runtime: "http" },
  { name: "delivery", role: "ship", runtime: "http" },
  { name: "router", role: "model", runtime: "http" },
  { name: "research", role: "library", runtime: "http" },
  { name: "explorer", role: "search", runtime: "http" },
  { name: "operator", role: "computer-use", runtime: "http" },
  { name: "reviewer", role: "qa", runtime: "http" },
  { name: "architect", role: "design", runtime: "http" },
  { name: "deploy", role: "release", runtime: "http" },
  { name: "security", role: "guard", runtime: "http" },
  { name: "memory", role: "recall", runtime: "http" },
  { name: "comms", role: "signal", runtime: "http" },
  { name: "planner", role: "backlog", runtime: "http" },
  { name: "investigator", role: "incident", runtime: "http" },
  { name: "manager", role: "instinct", runtime: "http" },
];

async function readJSON<T>(rel: string, fallback: T): Promise<T> {
  try {
    return JSON.parse(await readFile(path.resolve(process.cwd(), rel), "utf8")) as T;
  } catch {
    try {
      return JSON.parse(await readFile(path.resolve(process.cwd(), "../..", rel), "utf8")) as T;
    } catch {
      return fallback;
    }
  }
}

export async function seed(client: PrismaClient = prisma) {
  await client.controlState.upsert({
    where: { id: "global" },
    update: {},
    create: { id: "global", paused: false },
  });

  const inventory = await readJSON<{ repos: Array<{ name: string; org: string; fullName: string; access: string; framework: string | null; deployTarget: string[]; visibility: string }> }>(
    "inventory/repos.json",
    { repos: [] },
  );
  const backlog = await readJSON<{ tasks: Array<{ id: string; title: string; body: string; source: string; priority: string; status: string; projectSlug: string; dependsOn: string[] }> }>(
    "state/backlog.json",
    { tasks: [] },
  );

  for (const repo of inventory.repos) {
    const slug = repo.access === "ok" ? repo.name : `${repo.org}-denied`.toLowerCase();
    const project = await client.project.upsert({
      where: { slug },
      update: {
        repoUrl: repo.access === "ok" ? `https://github.com/${repo.fullName}` : null,
        status: repo.access === "ok" ? "active" : "denied",
        deployTarget: (repo.deployTarget || []).join(",") || null,
        tags: JSON.stringify([repo.visibility, repo.framework || "unknown"]),
      },
      create: {
        slug,
        repoUrl: repo.access === "ok" ? `https://github.com/${repo.fullName}` : null,
        org: repo.org,
        kind: repo.name === "cashtro" ? "control-plane" : "internal",
        status: repo.access === "ok" ? "active" : "denied",
        deployTarget: (repo.deployTarget || []).join(",") || null,
        healthUrl: repo.name === "cashtro" ? "http://127.0.0.1:8080/health" : null,
        tags: JSON.stringify([repo.visibility, repo.framework || "unknown"]),
      },
    });
    if (repo.name === "cashtro" && repo.access === "ok") {
      await client.capability.upsert({
        where: { id: "cap-kernel-invoke" },
        update: {},
        create: {
          id: "cap-kernel-invoke",
          projectId: project.id,
          name: "kernel.invoke",
          description: "POST /api/agents/:id/invoke on the Voltron kernel",
          invokeSpec: "POST http://127.0.0.1:8080/api/agents/{id}/invoke",
        },
      });
    }
  }

  const cashtro = await client.project.findUnique({ where: { slug: "cashtro" } });
  if (cashtro) {
    for (const spec of KERNEL_AGENTS) {
      const existing = await client.agent.findFirst({ where: { name: spec.name, projectId: cashtro.id } });
      const agent =
        existing ??
        (await client.agent.create({
          data: {
            projectId: cashtro.id,
            name: spec.name,
            role: spec.role,
            runtime: spec.runtime,
            modelRoute: "openrouter",
            maxCostPerRun: 2,
            enabled: true,
          },
        }));
      await client.secretRef.upsert({
        where: { id: `mcp-${spec.name}` },
        update: {},
        create: {
          id: `mcp-${spec.name}`,
          scope: `agent:${spec.name}`,
          provider: "env",
          key: `MCP_TOKEN_${spec.name.toUpperCase()}`,
        },
      });
      if (!existing) {
        await client.worker.create({
          data: { agentId: agent.id, name: `${spec.name}-worker`, tool: "kernel.invoke", concurrency: 1, timeoutSec: 120 },
        });
      }
    }
  }

  for (const t of backlog.tasks) {
    const project = await client.project.findUnique({ where: { slug: t.projectSlug } });
    await client.task.upsert({
      where: { idempotencyKey: t.id },
      update: { title: t.title, body: t.body, status: t.status === "doing" ? "doing" : t.status === "done" ? "done" : t.status === "blocked" ? "blocked" : "todo" },
      create: {
        title: t.title,
        body: t.body,
        source: t.source,
        priority: t.priority,
        status: t.status === "doing" ? "doing" : t.status === "done" ? "done" : t.status === "blocked" ? "blocked" : "todo",
        dependsOn: JSON.stringify(t.dependsOn || []),
        projectId: project?.id,
        idempotencyKey: t.id,
      },
    });
  }
}

if (import.meta.url === `file://${process.argv[1]}` || process.argv[1]?.endsWith("seed.ts")) {
  seed()
    .then(() => prisma.$disconnect())
    .catch(async (err) => {
      console.error(err);
      await prisma.$disconnect();
      process.exit(1);
    });
}
