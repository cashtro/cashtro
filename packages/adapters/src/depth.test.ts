import assert from "node:assert/strict";
import { test } from "node:test";
import { assertDepth, DepthError, MAX_DEPTH } from "./types.js";
import { createHttpAdapter } from "./http.js";
import type { Task } from "@cashtro/sdk";

const task: Task = {
  id: "t1",
  projectId: null,
  title: "os about",
  body: "",
  source: "test",
  priority: "p2",
  status: "todo",
  dependsOn: [],
  assigneeAgentId: null,
  idempotencyKey: "idem-depth",
};

test("depth hard stop is 3", () => {
  assert.equal(MAX_DEPTH, 3);
  assertDepth(1);
  assertDepth(3);
  assert.throws(() => assertDepth(4), DepthError);
});

test("http adapter refuses depth 4 before any fetch", async () => {
  let fetched = false;
  const adapter = createHttpAdapter({
    kernelUrl: "http://127.0.0.1:9",
    fetchImpl: async () => {
      fetched = true;
      return new Response("no");
    },
  });
  await assert.rejects(
    () => adapter.dispatch(task, { runId: "r1", depth: 4, budgetCapUsd: 2, kernelUrl: "http://127.0.0.1:9" }),
    DepthError,
  );
  assert.equal(fetched, false);
});
