import assert from "node:assert/strict";
import { test } from "node:test";
import { createN8nAdapter, resolveSpecialist, webhookFromTask, EMBEDDED_N8N_FLEET } from "./n8n.js";
import type { Task } from "@cashtro/sdk";

const task: Task = {
  id: "t1",
  projectId: null,
  title: "os pulse",
  body: "",
  source: "test",
  priority: "p2",
  status: "todo",
  dependsOn: [],
  assigneeAgentId: null,
  idempotencyKey: "idem-n8n",
};

test("fleet is 58 specialists on 14 kernel seats plus Steel", () => {
  assert.equal(EMBEDDED_N8N_FLEET.length, 58);
  assert.equal(new Set(EMBEDDED_N8N_FLEET.map((s) => s.seat)).size, 15);
  assert.equal(new Set(EMBEDDED_N8N_FLEET.map((s) => s.id)).size, 58);
  assert.ok(EMBEDDED_N8N_FLEET.some((s) => s.seat === "contrarian"));
  assert.ok(EMBEDDED_N8N_FLEET.some((s) => s.webhook === "improve-init"));
});

test("resolveSpecialist is wide: webhook, id, then longest name", () => {
  assert.equal(resolveSpecialist({ title: "os pulse", body: "" })?.webhook, "os-pulse");
  assert.equal(resolveSpecialist({ title: "webhook:ship-list", body: "" })?.id, "n8n-03");
  assert.equal(resolveSpecialist({ title: "ping n8n-07", body: "" })?.webhook, "model-status");
  assert.equal(resolveSpecialist({ title: "backlog scanapp", body: "" })?.webhook, "backlog-scanapp");
  assert.equal(resolveSpecialist({ title: "webhook:improve-delivery", body: "" })?.id, "n8n-42");
  assert.equal(resolveSpecialist({ title: "steel contradict", body: "" })?.seat, "contrarian");
  assert.notEqual(webhookFromTask({ title: "control plane", body: "" }), "plan-control-plane");
});

test("unbound n8n falls back to Voltron http", async () => {
  const adapter = createN8nAdapter({
    baseUrl: "",
    fallback: {
      id: "http",
      healthcheck: async () => ({ ok: true, detail: "mock" }),
      estimateCost: async () => 0,
      dispatch: async (_t, ctx) => ({ id: ctx.runId, adapter: "http", depth: ctx.depth }),
      poll: async () => ({
        status: "succeeded",
        inputTokens: 0,
        outputTokens: 0,
        costUsd: 0,
        exitReason: "http:init.os.about",
      }),
      cancel: async () => {},
    },
  });
  const h = await adapter.dispatch(task, { runId: "r1", depth: 1, budgetCapUsd: 2, kernelUrl: "http://voltron.test" });
  const st = await adapter.poll(h);
  assert.equal(st.status, "succeeded");
  assert.match(st.exitReason, /n8n-fallback:http:init.os.about/);
});

test("bound n8n posts the webhook and stays at $0", async () => {
  const adapter = createN8nAdapter({
    baseUrl: "http://n8n.test",
    fetchImpl: async (input, init) => {
      assert.equal(String(input), "http://n8n.test/webhook/os-pulse");
      assert.equal(init?.method, "POST");
      return new Response(JSON.stringify({ ok: true }), { status: 200 });
    },
  });
  const h = await adapter.dispatch(task, { runId: "r2", depth: 2, budgetCapUsd: 2, kernelUrl: "http://voltron.test" });
  const st = await adapter.poll(h);
  assert.equal(st.status, "succeeded");
  assert.equal(st.costUsd, 0);
  assert.equal(st.exitReason, "n8n:os-pulse");
});

test("bound n8n that errors falls back to Voltron", async () => {
  const adapter = createN8nAdapter({
    baseUrl: "http://n8n.test",
    fetchImpl: async () => new Response("down", { status: 502 }),
    fallback: {
      id: "http",
      healthcheck: async () => ({ ok: true, detail: "mock" }),
      estimateCost: async () => 0,
      dispatch: async (_t, ctx) => ({ id: ctx.runId, adapter: "http", depth: ctx.depth }),
      poll: async () => ({
        status: "succeeded",
        inputTokens: 0,
        outputTokens: 0,
        costUsd: 0,
        exitReason: "http:init.os.about",
      }),
      cancel: async () => {},
    },
  });
  const h = await adapter.dispatch(task, { runId: "r3", depth: 1, budgetCapUsd: 2, kernelUrl: "http://voltron.test" });
  const st = await adapter.poll(h);
  assert.equal(st.status, "succeeded");
  assert.match(st.exitReason, /n8n-fallback:os-pulse:http:init.os.about/);
});
