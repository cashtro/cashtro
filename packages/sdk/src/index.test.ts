import assert from "node:assert/strict";
import { test } from "node:test";
import { CreateTask } from "./index.js";

test("create task requires an idempotency key", () => {
  const bad = CreateTask.safeParse({ title: "x" });
  assert.equal(bad.success, false);
  const ok = CreateTask.parse({ title: "scan", idempotencyKey: "idem-1234" });
  assert.equal(ok.body, "");
});
