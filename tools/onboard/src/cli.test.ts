import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import path from "node:path";
import { test } from "node:test";
import { readFile } from "node:fs/promises";

test("onboard drafts cashtro handshake without inventing Evolu-Jeunes repos", () => {
  const cli = path.resolve(import.meta.dirname, "cli.ts");
  const root = path.resolve(import.meta.dirname, "../../..");
  const denied = spawnSync("pnpm", ["exec", "tsx", cli, "Evolu-Jeunes/scanapp"], { cwd: root, encoding: "utf8" });
  assert.notEqual(denied.status, 0);
  const ok = spawnSync("pnpm", ["exec", "tsx", cli, "cashtro/cashtro"], { cwd: root, encoding: "utf8" });
  assert.equal(ok.status, 0, ok.stderr);
  assert.match(ok.stdout, /cashtro\/cashtro/);
});

test("draft manifest is schemaVersion 1", async () => {
  const raw = await readFile(path.resolve(import.meta.dirname, "../../../inventory/drafts/cashtro/agent.manifest.json"), "utf8").catch(
    async () => await readFile(path.resolve(import.meta.dirname, "../../../agent.manifest.json"), "utf8"),
  );
  const json = JSON.parse(raw) as { schemaVersion: string; project: { slug: string } };
  assert.equal(json.schemaVersion, "1");
  assert.equal(json.project.slug, "cashtro");
});
