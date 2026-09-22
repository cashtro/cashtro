import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import { promisify } from "node:util";
import { test } from "node:test";
import { PrismaClient } from "@prisma/client";
import { createLocalAdapter } from "@cashtro/adapters";
import { seed } from "../prisma/seed.js";
import { buildApp } from "./app.js";

const exec = promisify(execFile);

async function withDb() {
  const dir = await mkdtemp(path.join(tmpdir(), "cp-"));
  const db = path.join(dir, "test.db");
  process.env.DATABASE_URL = `file:${db}`;
  const apiRoot = path.resolve(import.meta.dirname, "..");
  await exec("pnpm", ["exec", "prisma", "migrate", "deploy"], {
    cwd: apiRoot,
    env: { ...process.env, DATABASE_URL: `file:${db}` },
  });
  const prisma = new PrismaClient();
  await seed(prisma);
  const app = await buildApp({
    prisma,
    apiKeys: [{ id: "test", secret: "test-key", scope: "admin" }],
    budgetCap: 2,
  });
  return {
    app,
    prisma,
    async close() {
      await app.close();
      await prisma.$disconnect();
      await rm(dir, { recursive: true, force: true });
    },
  };
}

function auth(json?: unknown) {
  return {
    headers: { authorization: "Bearer test-key", "content-type": "application/json" },
    payload: JSON.stringify(json ?? {}),
  };
}

test("registry seed + every mutating route writes an event", async (t) => {
  const ctx = await withDb();
  t.after(() => ctx.close());

  const health = await ctx.app.inject({ method: "GET", url: "/health" });
  assert.equal(health.statusCode, 200);

  const docs = await ctx.app.inject({ method: "GET", url: "/docs" });
  assert.ok(docs.statusCode === 200 || docs.statusCode === 302, `docs ${docs.statusCode}`);

  const spec = await ctx.app.inject({ method: "GET", url: "/docs/json" });
  assert.equal(spec.statusCode, 200);
  const openapi = spec.json();
  assert.equal(openapi.openapi.startsWith("3."), true);

  const projects = await ctx.app.inject({ method: "GET", url: "/projects", ...auth() });
  assert.equal(projects.statusCode, 200);
  const list = projects.json();
  assert.ok(list.some((p: { slug: string }) => p.slug === "cashtro"));
  const scanapp = list.find((p: { slug: string; status: string }) => p.slug === "scanapp");
  assert.ok(scanapp, "scanapp concept project must be seeded");
  assert.equal(scanapp.status, "concept");

  const scanDetail = await ctx.app.inject({ method: "GET", url: "/projects/scanapp", ...auth() });
  assert.equal(scanDetail.statusCode, 200);
  const scanBody = scanDetail.json();
  assert.ok(scanBody.capabilities.some((c: { name: string }) => c.name === "scan.ingest"));
  assert.ok(scanBody.capabilities.some((c: { name: string }) => c.name === "crm.upsert"));
  assert.ok(scanBody.capabilities.some((c: { name: string }) => c.name === "bot.reply"));

  const scanTasks = await ctx.app.inject({ method: "GET", url: "/tasks?project=scanapp", ...auth() });
  assert.equal(scanTasks.statusCode, 200);
  assert.ok(scanTasks.json().length >= 3, "scanapp must have concept work attached");

  const detail = await ctx.app.inject({ method: "GET", url: "/projects/cashtro", ...auth() });
  assert.equal(detail.statusCode, 200);
  assert.ok(detail.json().agents.length >= 14);

  const agents = await ctx.app.inject({ method: "GET", url: "/agents", ...auth() });
  assert.equal(agents.statusCode, 200);
  assert.ok(agents.json().length >= 14);

  const fabric = await ctx.app.inject({ method: "GET", url: "/fleet", ...auth() });
  assert.equal(fabric.statusCode, 200);
  assert.equal(fabric.json().specialists, 58);
  assert.equal(fabric.json().kernelSeats, 14);
  assert.equal(fabric.json().departments, 15);
  assert.equal(fabric.json().selfImprove, true);
  assert.equal(fabric.json().contrarian, "steel-local");
  assert.equal(fabric.json().architecture, "wide-not-deep");
  assert.match(String(fabric.json().n8n), /unbound/);

  const corp = await ctx.app.inject({ method: "GET", url: "/corp", ...auth() });
  assert.equal(corp.statusCode, 200);
  assert.equal(corp.json().seats.length, 15);
  assert.ok(corp.json().seats.some((s: { crew: string }) => s.crew === "Steel"));

  const steelTask = await ctx.app.inject({
    method: "POST",
    url: "/tasks",
    ...auth({ title: "steel contradict the cheaper option", idempotencyKey: "idem-steel-0001" }),
  });
  const steelRun = await ctx.app.inject({
    method: "POST",
    url: `/tasks/${steelTask.json().id}/dispatch`,
    ...auth(),
  });
  assert.equal(steelRun.statusCode, 200);
  assert.equal(steelRun.json().status, "succeeded");
  assert.equal(steelRun.json().exitReason, "steel:contradict");
  assert.equal(steelRun.json().costUsd, 0);

  const created = await ctx.app.inject({
    method: "POST",
    url: "/tasks",
    ...auth({ title: "ping kernel", idempotencyKey: "idem-ping-0001" }),
  });
  assert.equal(created.statusCode, 201);
  const again = await ctx.app.inject({
    method: "POST",
    url: "/tasks",
    ...auth({ title: "ping kernel", idempotencyKey: "idem-ping-0001" }),
  });
  assert.equal(again.json().id, created.json().id);

  const dispatched = await ctx.app.inject({ method: "POST", url: `/tasks/${created.json().id}/dispatch`, ...auth() });
  assert.equal(dispatched.statusCode, 200);
  assert.equal(dispatched.json().status, "succeeded");
  assert.ok(dispatched.json().costUsd >= 0);
  const run = await ctx.app.inject({ method: "GET", url: `/runs/${dispatched.json().id}`, ...auth() });
  assert.equal(run.statusCode, 200);
  assert.ok(run.json().artifacts.length >= 1);

  const stream = await ctx.app.inject({ method: "GET", url: `/runs/${dispatched.json().id}/stream`, ...auth() });
  assert.match(stream.body, /data:/);

  const metrics = await ctx.app.inject({ method: "GET", url: "/metrics" });
  assert.match(metrics.body, /cashtro_cost_usd_total/);

  const gh = await ctx.app.inject({
    method: "POST",
    url: "/webhooks/github",
    ...auth({
      action: "opened",
      repository: { full_name: "cashtro/cashtro" },
      pull_request: { number: 14, title: "control plane", html_url: "https://example.test/14" },
    }),
  });
  assert.equal(gh.json().ok, true);
  assert.ok(gh.json().task?.id);
  const againGh = await ctx.app.inject({
    method: "POST",
    url: "/webhooks/github",
    ...auth({
      action: "opened",
      repository: { full_name: "cashtro/cashtro" },
      pull_request: { number: 14, title: "control plane", html_url: "https://example.test/14" },
    }),
  });
  assert.equal(againGh.json().deduped, true);
  assert.equal(againGh.json().task.id, gh.json().task.id);
  const cb = await ctx.app.inject({ method: "POST", url: "/webhooks/init/callback", ...auth({ done: true }) });
  assert.equal(cb.json().ok, true);

  const events = await ctx.prisma.event.count();
  assert.ok(events >= 4);
});

test("pause refuses dispatch", async (t) => {
  const ctx = await withDb();
  t.after(() => ctx.close());
  const task = await ctx.app.inject({
    method: "POST",
    url: "/tasks",
    ...auth({ title: "paused", idempotencyKey: "idem-pause-0001" }),
  });
  const pause = await ctx.app.inject({ method: "POST", url: "/control/pause", ...auth() });
  assert.equal(pause.json().paused, true);
  assert.ok(typeof pause.json().drained === "number");
  const dispatch = await ctx.app.inject({ method: "POST", url: `/tasks/${task.json().id}/dispatch`, ...auth() });
  assert.equal(dispatch.statusCode, 409);
  await ctx.app.inject({ method: "POST", url: `/tasks/${task.json().id}/cancel`, ...auth() });
  const cancelled = await ctx.app.inject({ method: "GET", url: `/tasks?status=cancelled`, ...auth() });
  assert.ok(cancelled.json().some((x: { id: string }) => x.id === task.json().id));
});

test("budget cap blocks an over-budget dispatch", async (t) => {
  const dir = await mkdtemp(path.join(tmpdir(), "cp-"));
  t.after(async () => rm(dir, { recursive: true, force: true }));
  const db = path.join(dir, "test.db");
  process.env.DATABASE_URL = `file:${db}`;
  const apiRoot = path.resolve(import.meta.dirname, "..");
  await exec("pnpm", ["exec", "prisma", "migrate", "deploy"], {
    cwd: apiRoot,
    env: { ...process.env, DATABASE_URL: `file:${db}` },
  });
  const prisma = new PrismaClient();
  t.after(() => prisma.$disconnect());
  await seed(prisma);
  const app = await buildApp({
    prisma,
    apiKeys: [{ id: "test", secret: "test-key", scope: "admin" }],
    budgetCap: 0.001,
    adapters: { http: createLocalAdapter(1), local: createLocalAdapter(1) },
  });
  t.after(() => app.close());
  const task = await app.inject({
    method: "POST",
    url: "/tasks",
    ...auth({ title: "expensive", idempotencyKey: "idem-budget-0001" }),
  });
  const dispatch = await app.inject({ method: "POST", url: `/tasks/${task.json().id}/dispatch`, ...auth() });
  assert.equal(dispatch.statusCode, 409);
  assert.equal(dispatch.json().error, "over-budget");
});

test("depth 4 is refused and /ui + /costs render", async (t) => {
  const ctx = await withDb();
  t.after(() => ctx.close());
  const task = await ctx.app.inject({
    method: "POST",
    url: "/tasks",
    ...auth({ title: "too deep", idempotencyKey: "idem-depth-0001" }),
  });
  const deep = await ctx.app.inject({
    method: "POST",
    url: `/tasks/${task.json().id}/dispatch`,
    ...auth({ depth: 4 }),
  });
  assert.equal(deep.statusCode, 409);
  assert.equal(deep.json().error, "depth-limit");

  const fleet = await ctx.app.inject({ method: "GET", url: "/ui/fleet" });
  assert.equal(fleet.statusCode, 200);
  assert.match(fleet.body, /scanapp/i);
  assert.match(fleet.body, /n8n specialists 58/);
  assert.match(fleet.body, /n8n-40/);
  assert.match(fleet.body, /Steel/);
  const corpUi = await ctx.app.inject({ method: "GET", url: "/ui/corp" });
  assert.equal(corpUi.statusCode, 200);
  assert.match(corpUi.body, /Steel/);
  assert.match(corpUi.body, /Lacune/);
  const costs = await ctx.app.inject({ method: "GET", url: "/costs", ...auth() });
  assert.equal(costs.statusCode, 200);
  assert.equal(typeof costs.json().totalUsd, "number");

  const queue = await ctx.app.inject({ method: "GET", url: "/ui/queue" });
  assert.equal(queue.statusCode, 200);
  assert.match(queue.body, /Coup/);
  assert.doesNotMatch(queue.body, />Dispatch</);

  const ready = await ctx.app.inject({
    method: "POST",
    url: "/tasks",
    ...auth({ title: "pane dispatch", idempotencyKey: "idem-pane-0001" }),
  });
  const coupPage = await ctx.app.inject({ method: "GET", url: `/ui/coup/${ready.json().id}` });
  assert.equal(coupPage.statusCode, 200);
  assert.match(coupPage.body, /Questions spécifiques/);
  assert.match(coupPage.body, /Position/);
  assert.match(coupPage.body, /Coup le plus court/);
  assert.match(coupPage.body, /Réponse adverse \(Steel\)/);
  assert.match(coupPage.body, /Jouer/);

  const blocked = await ctx.app.inject({
    method: "POST",
    url: `/ui/act/dispatch/${ready.json().id}`,
    headers: { "content-type": "application/x-www-form-urlencoded" },
    payload: "",
  });
  assert.equal(blocked.statusCode, 302);
  assert.equal(blocked.headers.location, `/ui/coup/${ready.json().id}?missing=1`);

  const pane = await ctx.app.inject({
    method: "POST",
    url: `/ui/act/dispatch/${ready.json().id}`,
    headers: { "content-type": "application/x-www-form-urlencoded" },
    payload:
      "question=faut-il+jouer+%3F&position=queue+ouverte&shortest=quatre+champs&opponent=un+clic+de+moins",
  });
  assert.ok(pane.statusCode === 302 || pane.statusCode === 200, `pane ${pane.statusCode} ${pane.body}`);
  const coupEvents = await ctx.prisma.event.findMany({ where: { type: "coup.before" } });
  assert.ok(coupEvents.some((e) => e.payload.includes(ready.json().id)));
});
