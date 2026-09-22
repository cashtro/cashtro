import assert from "node:assert/strict";
import { test } from "node:test";
import { challenge, createContrarianAdapter } from "./contrarian.js";
import type { Task } from "@cashtro/sdk";

test("steel refuses live clients and prefers the shorter option", () => {
  const live = challenge({ title: "touch BTK Avocats site", body: "" });
  assert.equal(live.liveClient, true);
  assert.match(live.moreOptimized, /Do not touch/);
  const local = challenge({
    title: "dispatch pane",
    body: "",
    shortest: "four beats then play",
    opponent: "one-click dispatch skips context",
  });
  assert.equal(local.liveClient, false);
  assert.match(local.contradiction, /one-click/);
});

test("contrarian adapter is $0 and records the challenge", async () => {
  const adapter = createContrarianAdapter();
  const task: Task = {
    id: "t-steel",
    projectId: null,
    title: "improve init",
    body: "opponent: manifesto is a slogan",
    source: "test",
    priority: "p2",
    status: "todo",
    dependsOn: [],
    assigneeAgentId: null,
    idempotencyKey: "idem-steel",
  };
  const h = await adapter.dispatch(task, { runId: "r-steel", depth: 1, budgetCapUsd: 2, kernelUrl: "http://127.0.0.1:8080" });
  const st = await adapter.poll(h);
  assert.equal(st.status, "succeeded");
  assert.equal(st.costUsd, 0);
  assert.equal(st.exitReason, "steel:contradict");
  assert.ok(st.artifact?.body?.includes("contradiction"));
});
