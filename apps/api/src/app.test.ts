import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import { promisify } from "node:util";
import { test } from "node:test";
import { PrismaClient } from "@prisma/client";
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

  const detail = await ctx.app.inject({ method: "GET", url: "/projects/cashtro", ...auth() });
  assert.equal(detail.statusCode, 200);
  assert.ok(detail.json().agents.length >= 14);

  const agents = await ctx.app.inject({ method: "GET", url: "/agents", ...auth() });
  assert.equal(agents.statusCode, 200);
  assert.ok(agents.json().length >= 14);

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
  assert.ok(dispatched.json().costUsd > 0);

  const run = await ctx.app.inject({ method: "GET", url: `/runs/${dispatched.json().id}`, ...auth() });
  assert.equal(run.statusCode, 200);
  assert.ok(run.json().artifacts.length >= 1);

  const stream = await ctx.app.inject({ method: "GET", url: `/runs/${dispatched.json().id}/stream`, ...auth() });
  assert.match(stream.body, /data:/);

  const metrics = await ctx.app.inject({ method: "GET", url: "/metrics" });
  assert.match(metrics.body, /cashtro_cost_usd_total/);

  const gh = await ctx.app.inject({ method: "POST", url: "/webhooks/github", ...auth({ action: "opened" }) });
  assert.equal(gh.json().ok, true);
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
