#!/usr/bin/env node
import { mkdir, readFile, writeFile, access } from "node:fs/promises";
import path from "node:path";
import { GitHub, resolveToken } from "./github.js";
import { deniedOwner, probeRepo } from "./probe.js";
import { writeGraph, writeReport } from "./report.js";
import { CursorState, Inventory, mergeInventory, type RepoRecord } from "./schema.js";

async function findRoot(start: string): Promise<string> {
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

type Config = {
  controlPlaneRepo: string;
  orgs: string[];
  users: string[];
};

async function loadJSON<T>(file: string, fallback: T): Promise<T> {
  try {
    return JSON.parse(await readFile(file, "utf8")) as T;
  } catch {
    return fallback;
  }
}

async function main() {
  const root = await findRoot(path.resolve(process.cwd()));
  const cfg = (await loadJSON<Config>(path.join(root, "control-plane.config.json"), {
    controlPlaneRepo: "cashtro/cashtro",
    orgs: ["Evolu-Jeunes"],
    users: ["cashtro"],
  })) as Config;

  const token = await resolveToken();
  const gh = new GitHub(token);
  const now = new Date().toISOString();
  const cursorPath = path.join(root, "state/recon-cursor.json");
  const invPath = path.join(root, "inventory/repos.json");
  const cursor = CursorState.parse(
    await loadJSON(cursorPath, { completed: [], updatedAt: now }),
  );
  const prevInv = Inventory.safeParse(await loadJSON(invPath, null));
  const prev = prevInv.success ? prevInv.data.repos : [];

  const next: RepoRecord[] = [];
  const owners = [...cfg.users.map((u) => ({ owner: u, kind: "user" as const })), ...cfg.orgs.map((o) => ({ owner: o, kind: "org" as const }))];

  for (const { owner, kind } of owners) {
    const listed = await gh.listOwnerRepos(owner, kind);
    if (listed.access !== "ok" || listed.repos.length === 0) {
      next.push(deniedOwner(owner, listed.access === "ok" ? "limited" : listed.access, now));
      continue;
    }
    for (const repo of listed.repos) {
      const key = repo.full_name.toLowerCase();
      process.stderr.write(`recon ${repo.full_name}\n`);
      const probed = await probeRepo(gh, repo, now);
      next.push(probed);
      if (!cursor.completed.includes(key)) cursor.completed.push(key);
    }
  }

  const merged = mergeInventory(prev, next);
  const inventory: Inventory = {
    schemaVersion: "1",
    generatedAt: now,
    controlPlaneRepo: cfg.controlPlaneRepo,
    owners: [...cfg.users, ...cfg.orgs],
    repos: merged,
  };
  Inventory.parse(inventory);

  await mkdir(path.join(root, "inventory"), { recursive: true });
  await mkdir(path.join(root, "state"), { recursive: true });
  await writeFile(invPath, JSON.stringify(inventory, null, 2) + "\n");
  await writeFile(path.join(root, "inventory/REPORT.md"), writeReport(inventory));
  await writeFile(path.join(root, "inventory/GRAPH.md"), writeGraph(inventory));
  cursor.updatedAt = now;
  await writeFile(cursorPath, JSON.stringify(cursor, null, 2) + "\n");
  process.stdout.write(`recon ok · ${merged.filter((r) => r.access === "ok").length} repos · ${merged.length} rows\n`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
