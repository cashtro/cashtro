import type { RepoRecord } from "./schema.js";
import type { GhRepo, GitHub } from "./github.js";

const PROBE_FILES = [
  "package.json",
  "pnpm-lock.yaml",
  "yarn.lock",
  "package-lock.json",
  "composer.json",
  "go.mod",
  "requirements.txt",
  "pyproject.toml",
  "next.config.js",
  "next.config.mjs",
  "next.config.ts",
  "Dockerfile",
  "vercel.json",
  "azure-pipelines.yml",
  "Makefile",
  ".nvmrc",
  "AGENTS.md",
  "CLAUDE.md",
  ".mcp.json",
  "agent.manifest.json",
];

function pickEnvKeys(text: string): string[] {
  const found = new Set<string>();
  const re = /process\.env\.([A-Z][A-Z0-9_]+)|os\.Getenv\("([A-Z][A-Z0-9_]+)"\)|\$\{?([A-Z][A-Z0-9_]+)\}?/g;
  let m: RegExpExecArray | null;
  while ((m = re.exec(text))) {
    const key = m[1] || m[2] || m[3];
    if (key && key.length > 2 && !["HTTP", "HTTPS", "GET", "POST"].includes(key)) found.add(key);
  }
  return [...found].sort();
}

export async function probeRepo(gh: GitHub, repo: GhRepo, now: string): Promise<RepoRecord> {
  const files = new Map<string, string>();
  for (const path of PROBE_FILES) {
    const body = await gh.getFile(repo.full_name, path);
    if (body) files.set(path, body);
  }
  const workflows = await gh.listDir(repo.full_name, ".github/workflows");
  const wp = await gh.listDir(repo.full_name, "wp-content");
  const srcHints = [files.get("go.mod"), files.get("package.json"), files.get("Makefile")].filter(Boolean).join("\n");

  let framework: string | null = null;
  let packageManager: string | null = null;
  let runtimeVersion: string | null = null;
  const entryCommands: string[] = [];
  const deployTarget: string[] = [];
  const agentConfig: string[] = [];

  if (files.has("go.mod")) {
    framework = "Go";
    packageManager = "go";
    const go = files.get("go.mod")!.match(/^go\s+([0-9.]+)/m);
    if (go) runtimeVersion = `go${go[1]}`;
  }
  if (files.has("composer.json") || wp.length) {
    framework = framework ? `${framework}+WordPress` : "WordPress";
    packageManager = packageManager || "composer";
  }
  if (files.has("package.json")) {
    packageManager = files.has("pnpm-lock.yaml")
      ? "pnpm"
      : files.has("yarn.lock")
        ? "yarn"
        : files.has("package-lock.json")
          ? "npm"
          : "npm";
    try {
      const pkg = JSON.parse(files.get("package.json")!);
      if (pkg.engines?.node) runtimeVersion = runtimeVersion || `node${pkg.engines.node}`;
      if (files.has("next.config.js") || files.has("next.config.mjs") || files.has("next.config.ts") || pkg.dependencies?.next) {
        framework = "Next.js";
      } else if (pkg.dependencies?.react) {
        framework = "React";
      } else if (!framework) {
        framework = "Node";
      }
      for (const [name] of Object.entries(pkg.scripts || {})) entryCommands.push(`npm:${name}`);
    } catch {
      /* ignore broken package.json */
    }
  }
  if (files.has(".nvmrc")) runtimeVersion = runtimeVersion || `node${files.get(".nvmrc")!.trim()}`;
  if (files.has("requirements.txt") || files.has("pyproject.toml")) {
    framework = framework ? `${framework}+Python` : "Python";
  }
  if (files.has("Makefile")) {
    for (const line of files.get("Makefile")!.split("\n")) {
      const t = line.match(/^([a-zA-Z0-9_-]+):/);
      if (t) entryCommands.push(`make:${t[1]}`);
    }
  }

  if (files.has("Dockerfile")) deployTarget.push("docker");
  if (files.has("vercel.json")) deployTarget.push("vercel");
  if (files.has("azure-pipelines.yml")) deployTarget.push("azure");
  if (workflows.some((w) => /azure|deploy/i.test(w))) deployTarget.push("github-actions");
  if (workflows.length && !deployTarget.includes("github-actions")) deployTarget.push("github-actions");

  for (const name of ["AGENTS.md", "CLAUDE.md", ".mcp.json", "agent.manifest.json"]) {
    if (files.has(name)) agentConfig.push(name);
  }

  const ciStatus = workflows.length ? await gh.latestWorkflow(repo.full_name) : null;
  const counts = await gh.openCounts(repo.full_name);
  const envVars = pickEnvKeys(srcHints);

  const riskFlags: string[] = [];
  if (!workflows.length) riskFlags.push("no-ci");
  if (!files.has("go.mod") && !files.has("package.json") && !files.has("composer.json")) {
    /* unknown stack is not automatically no-tests */
  }
  const last = repo.pushed_at || repo.updated_at;
  if (last && Date.now() - Date.parse(last) > 365 * 24 * 3600 * 1000) riskFlags.push("stale-12mo");
  if (repo.archived) riskFlags.push("archived");
  if (envVars.some((k) => /SECRET|TOKEN|PASSWORD|PRIVATE_KEY/i.test(k)) && files.has("package.json")) {
    /* referenced, not committed — flag only if a secret-looking value is in a probed file */
  }
  for (const [path, body] of files) {
    if (/(AKIA[0-9A-Z]{16})|(sk-or-[A-Za-z0-9]{20,})|(ghp_[A-Za-z0-9]{20,})/.test(body)) {
      riskFlags.push(`possible-secret:${path}`);
    }
  }

  return {
    name: repo.name,
    org: repo.owner.login,
    fullName: repo.full_name,
    visibility: repo.private ? "private" : "public",
    defaultBranch: repo.default_branch ?? null,
    lastCommitAt: last ?? null,
    archived: !!repo.archived,
    access: "ok",
    primaryLanguage: repo.language,
    framework,
    packageManager,
    runtimeVersion,
    deployTarget: [...new Set(deployTarget)],
    ciStatus,
    agentConfig,
    requiredEnvVars: envVars,
    entryCommands: [...new Set(entryCommands)],
    openPrs: counts.prs,
    openIssues: counts.issues,
    riskFlags: [...new Set(riskFlags)],
    scannedAt: now,
  };
}

export function deniedOwner(owner: string, access: "denied" | "limited", now: string): RepoRecord {
  return {
    name: "*",
    org: owner,
    fullName: `${owner}/*`,
    visibility: "unknown",
    defaultBranch: null,
    lastCommitAt: null,
    archived: false,
    access,
    primaryLanguage: null,
    framework: null,
    packageManager: null,
    runtimeVersion: null,
    deployTarget: [],
    ciStatus: null,
    agentConfig: [],
    requiredEnvVars: [],
    entryCommands: [],
    openPrs: 0,
    openIssues: 0,
    riskFlags: ["access-denied"],
    scannedAt: now,
  };
}
