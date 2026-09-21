import assert from "node:assert/strict";
import { test } from "node:test";
import { ingestGithubWebhook } from "./githubWebhook.js";

test("github PR opened becomes an idempotent task", () => {
  const a = ingestGithubWebhook({
    action: "opened",
    repository: { full_name: "cashtro/cashtro" },
    pull_request: { number: 14, title: "control plane", html_url: "https://github.com/cashtro/cashtro/pull/14" },
  });
  const b = ingestGithubWebhook({
    action: "opened",
    repository: { full_name: "cashtro/cashtro" },
    pull_request: { number: 14, title: "control plane", html_url: "https://github.com/cashtro/cashtro/pull/14" },
  });
  assert.ok(a);
  assert.equal(a?.idempotencyKey, b?.idempotencyKey);
  assert.match(a?.title || "", /PR #14/);
});

test("CI conclusion becomes a task", () => {
  const t = ingestGithubWebhook({
    action: "completed",
    repository: { full_name: "cashtro/cashtro" },
    workflow_run: { id: 99, name: "go", conclusion: "success", html_url: "https://github.com/cashtro/cashtro/actions" },
  });
  assert.ok(t);
  assert.match(t?.title || "", /success/);
});
