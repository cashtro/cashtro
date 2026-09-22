import { readFile, access } from "node:fs/promises";
import path from "node:path";

export const TOOLS = [
  {
    name: "list_inventory",
    description: "Repos recon has actually seen. Denied owners stay denied.",
    inputSchema: { type: "object", properties: {} },
  },
  {
    name: "list_catalog",
    description: "Kernel catalog ships. Not GitHub-scanned.",
    inputSchema: { type: "object", properties: {} },
  },
  {
    name: "list_backlog",
    description: "Normalized control-plane backlog.",
    inputSchema: { type: "object", properties: { status: { type: "string" } } },
  },
  {
    name: "list_n8n_fleet",
    description: "n8n specialists on 14 kernel seats plus Steel. Wide, not deep. Each seat self-improves.",
    inputSchema: { type: "object", properties: {} },
  },
  {
    name: "list_corporation",
    description: "Cashtro corporation: 14 kernel departments + Steel. Gaps, asks, self-improve.",
    inputSchema: { type: "object", properties: {} },
  },
  {
    name: "get_access",
    description: "What this agent can and cannot see on GitHub.",
    inputSchema: { type: "object", properties: {} },
  },
] as const;

export async function findRoot(start: string): Promise<string> {
  let dir = start;
  for (;;) {
    try {
      await access(path.join(dir, "control-plane.config.json"));
      return dir;
    } catch {
      const parent = path.dirname(dir);
      if (parent === dir) return start;
      dir = parent;
    }
  }
}

export async function callTool(root: string, name: string, args: Record<string, unknown> = {}): Promise<unknown> {
  if (name === "list_inventory") {
    return JSON.parse(await readFile(path.join(root, "inventory/repos.json"), "utf8"));
  }
  if (name === "list_catalog") {
    return JSON.parse(await readFile(path.join(root, "state/catalog-ships.json"), "utf8"));
  }
  if (name === "list_backlog") {
    const raw = JSON.parse(await readFile(path.join(root, "state/backlog.json"), "utf8")) as {
      tasks: Array<{ status: string }>;
    };
    const status = typeof args.status === "string" ? args.status : "";
    return { ...raw, tasks: status ? raw.tasks.filter((t) => t.status === status) : raw.tasks };
  }
  if (name === "list_n8n_fleet") {
    return JSON.parse(await readFile(path.join(root, "state/n8n-fleet.json"), "utf8"));
  }
  if (name === "list_corporation") {
    return JSON.parse(await readFile(path.join(root, "state/corporation.json"), "utf8"));
  }
  if (name === "get_access") {
    const inv = JSON.parse(await readFile(path.join(root, "inventory/repos.json"), "utf8")) as {
      repos: Array<{ access: string; fullName: string }>;
    };
    return {
      reachable: inv.repos.filter((r) => r.access === "ok").map((r) => r.fullName),
      denied: inv.repos.filter((r) => r.access !== "ok").map((r) => r.fullName),
      note: "Connecting Cursor desktop ≠ this Cloud Agent token. Evolu-Jeunes needs the Cursor GitHub App on the org.",
    };
  }
  throw new Error(`unknown tool ${name}`);
}
