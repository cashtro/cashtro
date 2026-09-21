import assert from "node:assert/strict";
import { test } from "node:test";
import { Inventory, mergeInventory, repoKey, type RepoRecord } from "./schema.js";

function stub(name: string, scannedAt: string): RepoRecord {
  return {
    name,
    org: "cashtro",
    fullName: `cashtro/${name}`,
    visibility: "public",
    defaultBranch: "main",
    lastCommitAt: scannedAt,
    archived: false,
    access: "ok",
    primaryLanguage: "Go",
    framework: "Go",
    packageManager: "go",
    runtimeVersion: "go1.22",
    deployTarget: ["github-actions"],
    ciStatus: "success",
    agentConfig: [],
    requiredEnvVars: [],
    entryCommands: ["make:test"],
    openPrs: 0,
    openIssues: 0,
    riskFlags: [],
    scannedAt,
  };
}

test("merge is idempotent by org/name", () => {
  const a = stub("cashtro", "2026-09-21T00:00:00Z");
  const b = { ...stub("cashtro", "2026-09-21T12:00:00Z"), ciStatus: "failure" };
  const once = mergeInventory([], [a]);
  const twice = mergeInventory(once, [b]);
  const thrice = mergeInventory(twice, [b]);
  assert.equal(twice.length, 1);
  assert.equal(thrice.length, 1);
  assert.equal(repoKey(twice[0]), "cashtro/cashtro");
  assert.equal(twice[0].ciStatus, "failure");
  assert.deepEqual(twice, thrice);
});

test("inventory schema rejects a bad row", () => {
  const parsed = Inventory.safeParse({
    schemaVersion: "1",
    generatedAt: "now",
    controlPlaneRepo: "cashtro/cashtro",
    owners: ["cashtro"],
    repos: [{ name: "x" }],
  });
  assert.equal(parsed.success, false);
});
