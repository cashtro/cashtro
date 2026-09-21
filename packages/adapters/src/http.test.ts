import assert from "node:assert/strict";
import { test } from "node:test";
import { createHttpAdapter } from "./http.js";
import type { Task } from "@cashtro/sdk";

const task: Task = {
  id: "t1",
  projectId: null,
  title: "ping kernel",
  body: "",
  source: "test",
  priority: "p2",
  status: "todo",
  dependsOn: [],
  assigneeAgentId: null,
  idempotencyKey: "idem-http",
};

test("http adapter invokes kernel and records a $0 artifact", async () => {
  const adapter = createHttpAdapter({
    kernelUrl: "http://voltron.test",
    fetchImpl: async (input) => {
      const url = String(input);
      if (url.endsWith("/health")) {
        return Response.json({ status: "ok", running: 14 });
      }
      assert.match(url, /\/api\/agents\/init\/invoke$/);
      return Response.json({ ok: true, message: "Cashtro OS", data: { version: "0.2.0" } });
    },
  });
  const health = await adapter.healthcheck();
  assert.equal(health.ok, true);
  assert.equal(await adapter.estimateCost(task), 0);
  const handle = await adapter.dispatch(task, {
    runId: "run-1",
    depth: 1,
    budgetCapUsd: 2,
    kernelUrl: "http://voltron.test",
  });
  const status = await adapter.poll(handle);
  assert.equal(status.status, "succeeded");
  assert.equal(status.costUsd, 0);
  assert.ok(status.artifact?.body?.includes("Cashtro OS"));
});
