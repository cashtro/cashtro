import assert from "node:assert/strict";
import { test } from "node:test";
import { scanText } from "./cli.js";

test("flags a github pat and ignores clean text", () => {
  assert.ok(scanText("x", "ghp_abcdefghijklmnopqrstuv").length >= 1);
  assert.equal(scanText("x", "just a backlog note about tokens as names").length, 0);
});
