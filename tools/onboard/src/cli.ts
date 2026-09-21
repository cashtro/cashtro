#!/usr/bin/env node
import { mkdir, readFile, writeFile, access } from "node:fs/promises";
import path from "node:path";

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

type Repo = { name: string; org: string; fullName: string; access: string; framework?: string | null };

function draftManifest(repo: Repo) {
  return {
    schemaVersion: "1",
    project: { slug: repo.name, kind: repo.name === "cashtro" ? "control-plane" : "product", owner: repo.org },
    capabilities: [
      { name: "test", cmd: "pnpm test", timeoutSec: 300 },
      { name: "build", cmd: "pnpm build", timeoutSec: 600 },
    ],
    agents: [{ name: `${repo.name}-operator`, runtime: "http", promptFile: "AGENTS.md" }],
    env: ["DATABASE_URL"],
    health: { url: "", expect: 200 },
  };
}

function draftAgents(repo: Repo) {
  return `# AGENTS — ${repo.fullName}

Drafted by \`pnpm onboard\`. Review before merge.

- Owner: ${repo.org}
- Do not touch live client production from this handshake.
- Publish \`agent.manifest.json\` so the control plane can register capabilities.
`;
}

async function main() {
  const root = await findRoot(path.resolve(process.cwd()));
  const want = process.argv[2] || "cashtro/cashtro";
  const inv = JSON.parse(await readFile(path.join(root, "inventory/repos.json"), "utf8")) as { repos: Repo[] };
  const repo = inv.repos.find((r) => r.fullName.toLowerCase() === want.toLowerCase() || r.name === want);
  if (!repo || repo.access !== "ok") {
    process.stderr.write(`onboard: ${want} is not a reachable inventory row\n`);
    process.exit(2);
  }
  const dest = path.join(root, "inventory/drafts", repo.name);
  await mkdir(dest, { recursive: true });
  const manifest = draftManifest(repo);
  await writeFile(path.join(dest, "agent.manifest.json"), JSON.stringify(manifest, null, 2) + "\n");
  await writeFile(path.join(dest, "AGENTS.md"), draftAgents(repo));
  if (repo.name === "cashtro") {
    await writeFile(path.join(root, "agent.manifest.json"), JSON.stringify(manifest, null, 2) + "\n");
    // AGENTS.md at root is the reviewed copy; do not overwrite blindly if it exists.
  }
  process.stdout.write(`onboard drafted ${repo.fullName} → inventory/drafts/${repo.name}\n`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
