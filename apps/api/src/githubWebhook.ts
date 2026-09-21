export type IngestedTask = {
  title: string;
  body: string;
  source: string;
  idempotencyKey: string;
};

export function ingestGithubWebhook(body: unknown): IngestedTask | null {
  if (!body || typeof body !== "object") return null;
  const ev = body as Record<string, unknown>;
  const repo = (ev.repository as { full_name?: string } | undefined)?.full_name || "unknown";
  if (repo.toLowerCase().includes("evolu") && repo.toLowerCase() !== "unknown") {
    // still allow — we don't invent names, we echo what GitHub sent
  }

  const action = String(ev.action || "");
  const pr = ev.pull_request as { number?: number; title?: string; html_url?: string } | undefined;
  if (pr?.number) {
    return {
      title: `GitHub PR #${pr.number} ${action || "event"}: ${pr.title || ""}`.trim(),
      body: `${pr.html_url || ""}\nrepo=${repo}`,
      source: "github",
      idempotencyKey: `gh-pr-${repo}-${pr.number}-${action || "event"}`.slice(0, 80),
    };
  }

  const run = ev.workflow_run as { id?: number; name?: string; conclusion?: string; html_url?: string } | undefined;
  if (run?.id) {
    return {
      title: `CI ${run.conclusion || action || "run"}: ${run.name || "workflow"}`,
      body: `${run.html_url || ""}\nrepo=${repo}`,
      source: "github",
      idempotencyKey: `gh-ci-${repo}-${run.id}`.slice(0, 80),
    };
  }

  const pusher = ev.pusher as { name?: string } | undefined;
  const after = typeof ev.after === "string" ? ev.after : "";
  if (after || pusher) {
    return {
      title: `push on ${repo}`,
      body: `sha=${after}\nby=${pusher?.name || "unknown"}`,
      source: "github",
      idempotencyKey: `gh-push-${repo}-${after || "head"}`.slice(0, 80),
    };
  }

  if (action) {
    return {
      title: `github ${action} on ${repo}`,
      body: JSON.stringify({ action, repo }),
      source: "github",
      idempotencyKey: `gh-misc-${repo}-${action}`.slice(0, 80),
    };
  }
  return null;
}
