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

test("real Voltron http dispatch records run + artifact + $0", async (t) => {
  let health: Response;
  try {
    health = await fetch("http://127.0.0.1:8080/health");
  } catch {
    t.skip("Voltron not on :8080");
    return;
  }
  if (!health.ok) {
    t.skip("Voltron health not ok");
    return;
  }

  const dir = await mkdtemp(path.join(tmpdir(), "cp-e2e-"));
  t.after(() => rm(dir, { recursive: true, force: true }));
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
    budgetCap: 2,
    kernelUrl: "http://127.0.0.1:8080",
  });
  t.after(() => app.close());

  const created = await app.inject({
    method: "POST",
    url: "/tasks",
    headers: { authorization: "Bearer test-key", "content-type": "application/json" },
    payload: JSON.stringify({ title: "os about", idempotencyKey: "idem-voltron-e2e-0001" }),
  });
  assert.equal(created.statusCode, 201);
  const dispatched = await app.inject({
    method: "POST",
    url: `/tasks/${created.json().id}/dispatch`,
    headers: { authorization: "Bearer test-key", "content-type": "application/json" },
    payload: "{}",
  });
  assert.equal(dispatched.statusCode, 200, dispatched.body);
  const run = dispatched.json();
  assert.equal(run.status, "succeeded");
  assert.equal(run.costUsd, 0);
  assert.match(String(run.exitReason), /http:init\.os\.about|local-queue/);
  const detail = await app.inject({
    method: "GET",
    url: `/runs/${run.id}`,
    headers: { authorization: "Bearer test-key" },
  });
  assert.ok(detail.json().artifacts.length >= 1);
});
