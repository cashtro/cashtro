import { z } from "zod";

export const Access = z.enum(["ok", "denied", "limited"]);

export const RepoRecord = z.object({
  name: z.string(),
  org: z.string(),
  fullName: z.string(),
  visibility: z.enum(["public", "private", "internal", "unknown"]),
  defaultBranch: z.string().nullable(),
  lastCommitAt: z.string().nullable(),
  archived: z.boolean(),
  access: Access,
  primaryLanguage: z.string().nullable(),
  framework: z.string().nullable(),
  packageManager: z.string().nullable(),
  runtimeVersion: z.string().nullable(),
  deployTarget: z.array(z.string()),
  ciStatus: z.string().nullable(),
  agentConfig: z.array(z.string()),
  requiredEnvVars: z.array(z.string()),
  entryCommands: z.array(z.string()),
  openPrs: z.number().int().nonnegative(),
  openIssues: z.number().int().nonnegative(),
  riskFlags: z.array(z.string()),
  scannedAt: z.string(),
});

export type RepoRecord = z.infer<typeof RepoRecord>;

export const Inventory = z.object({
  schemaVersion: z.literal("1"),
  generatedAt: z.string(),
  controlPlaneRepo: z.string(),
  owners: z.array(z.string()),
  repos: z.array(RepoRecord),
});

export type Inventory = z.infer<typeof Inventory>;

export const CursorState = z.object({
  completed: z.array(z.string()),
  updatedAt: z.string(),
});

export type CursorState = z.infer<typeof CursorState>;

export function repoKey(r: Pick<RepoRecord, "org" | "name">): string {
  return `${r.org}/${r.name}`.toLowerCase();
}

export function mergeInventory(prev: RepoRecord[], next: RepoRecord[]): RepoRecord[] {
  const map = new Map<string, RepoRecord>();
  for (const r of prev) map.set(repoKey(r), r);
  for (const r of next) map.set(repoKey(r), r);
  return [...map.values()].sort((a, b) => a.fullName.localeCompare(b.fullName));
}
