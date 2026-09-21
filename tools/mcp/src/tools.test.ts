import assert from "node:assert/strict";
import path from "node:path";
import { test } from "node:test";
import { callTool, findRoot } from "./tools.js";

test("MCP tools read inventory and refuse to invent Evolu-Jeunes repos", async () => {
  const root = await findRoot(path.resolve(process.cwd()));
  const inv = (await callTool(root, "list_inventory")) as { repos: Array<{ fullName: string; access: string }> };
  assert.ok(inv.repos.some((r) => r.fullName === "cashtro/cashtro" && r.access === "ok"));
  const access = (await callTool(root, "get_access")) as { denied: string[] };
  assert.ok(access.denied.some((n) => n.startsWith("Evolu-Jeunes")));
  const catalog = (await callTool(root, "list_catalog")) as { ships: Array<{ slug: string }> };
  assert.ok(catalog.ships.some((s) => s.slug === "scanapp"));
  const backlog = (await callTool(root, "list_backlog", { status: "blocked" })) as { tasks: Array<{ id: string }> };
  assert.ok(backlog.tasks.some((t) => t.id === "evolu-access" || t.id === "scanapp-find-repo"));
  const n8n = (await callTool(root, "list_n8n_fleet")) as { specialists: Array<{ id: string; seat: string }> };
  assert.equal(n8n.specialists.length, 40);
  assert.equal(new Set(n8n.specialists.map((s) => s.seat)).size, 14);
});
