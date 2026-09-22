import { CATALOG_SHIPS } from "./catalog-known.js";
import type { Inventory, RepoRecord } from "./schema.js";

export function writeReport(inv: Inventory): string {
  const ok = inv.repos.filter((r) => r.access === "ok");
  const blocked = inv.repos.filter((r) => r.access !== "ok");
  const stale = ok.filter((r) => r.riskFlags.includes("stale-12mo") || r.archived);
  const active = ok.filter((r) => !stale.includes(r));
  const live = CATALOG_SHIPS.filter((s) => s.liveClient);
  const concept = CATALOG_SHIPS.filter((s) => !s.liveClient);
  const lines: string[] = [];
  lines.push(`# Inventory report`);
  lines.push("");
  lines.push(`Generated ${inv.generatedAt}. Control plane: \`${inv.controlPlaneRepo}\`.`);
  lines.push(`Owners scanned: ${inv.owners.join(", ")}.`);
  lines.push("");
  lines.push(`| Bucket | Count |`);
  lines.push(`| --- | ---: |`);
  lines.push(`| Reachable GitHub repos | ${ok.length} |`);
  lines.push(`| Denied / limited owners | ${blocked.length} |`);
  lines.push(`| Catalog ships (not scanned) | ${CATALOG_SHIPS.length} |`);
  lines.push(`| Live client (do not touch) | ${live.length} |`);
  lines.push(`| Dead or archivable (scanned) | ${stale.length} |`);
  lines.push("");
  lines.push(`## Active — manager can already see`);
  lines.push("");
  if (!active.length) lines.push("_None._");
  for (const r of active) lines.push(row(r));
  lines.push("");
  lines.push(`## Internal products (from this scan)`);
  lines.push("");
  for (const r of ok.filter((x) => x.org.toLowerCase() === "cashtro")) lines.push(row(r));
  if (!ok.some((x) => x.org.toLowerCase() === "cashtro")) lines.push("_None._");
  lines.push("");
  lines.push(`## Named in kernel catalog — not GitHub-scanned`);
  lines.push("");
  lines.push("These come from `internal/catalog`. Recon did **not** open their repos.");
  lines.push("");
  for (const s of CATALOG_SHIPS) {
    const flag = s.liveClient ? "LIVE CLIENT — do not touch" : "safe to onboard later";
    lines.push(`- **${s.name}** · ${s.stage} · ${s.stack} · ${s.sector} · ${flag}`);
  }
  lines.push("");
  lines.push(`## Dead or archivable`);
  lines.push("");
  if (!stale.length) lines.push("_None in the reachable set._");
  for (const r of stale) lines.push(row(r));
  lines.push("");
  lines.push(`## Needs access`);
  lines.push("");
  if (!blocked.length) lines.push("_All configured owners returned repos._");
  for (const r of blocked) {
    lines.push(`- \`${r.fullName}\` — access **${r.access}**. See \`docs/ACCESS_REQUIRED.md\`.`);
  }
  lines.push("");
  lines.push(`Expected behind that wall: ~48 Evolu-Jeunes private repos, including`);
  lines.push(`${live.map((s) => s.name).join(", ")}, plus ${concept.map((s) => s.name).join(", ")}.`);
  lines.push("");
  lines.push(`## What the manager can already control`);
  lines.push("");
  lines.push(
    ok.length
      ? `Only \`cashtro/cashtro\`: Go kernel, 14 agentics, this control-plane API. Client fleet is dark until Option A.`
      : `Nothing yet.`,
  );
  lines.push("");
  return lines.join("\n");
}

function row(r: RepoRecord): string {
  return `- **${r.fullName}** · ${r.visibility} · ${r.framework || r.primaryLanguage || "unknown"} · CI ${r.ciStatus || "none"} · PRs ${r.openPrs} · risks: ${r.riskFlags.join(", ") || "none"}`;
}

export function writeGraph(inv: Inventory): string {
  const ok = inv.repos.filter((r) => r.access === "ok");
  const live = CATALOG_SHIPS.filter((s) => s.liveClient);
  const later = CATALOG_SHIPS.filter((s) => !s.liveClient);
  const lines = [
    "# Inventory graph",
    "",
    "Cashtro (control plane) and Evolu-Jeunes (delivery org) mapped together.",
    "Live client ships are catalog-named only. Do not touch production.",
    "Canonical desk map: `GET /map` · JSON `GET /api/map` · `docs/FLEET_MAP.md`.",
    "",
    "```mermaid",
    "flowchart TB",
    "  castro[Castro]",
    "  subgraph cashtro_home[cashtro GitHub user]",
  ];
  for (const r of ok) {
    const id = "repo_" + r.fullName.replace(/[^a-zA-Z0-9]/g, "_");
    lines.push(`    ${id}[${r.fullName}]`);
  }
  if (!ok.length) lines.push("    cashtro_empty[no reachable repos]");
  lines.push("  end");
  lines.push("  subgraph evolu_home[Evolu-Jeunes GitHub org]");
  lines.push("    evolu_member[github.com/evoluJeunes]");
  lines.push("    evolu_site[evolujeunes.ca / LIVE do not touch]");
  lines.push("    evolu_dark[~48 private repos / access limited]");
  lines.push("    subgraph live[LIVE CLIENT — catalog named only]");
  for (const s of live) {
    lines.push(`      cat_${s.name.replace(/[^a-zA-Z0-9]/g, "_")}[${s.name} / ${s.stage}]`);
  }
  lines.push("    end");
  lines.push("    subgraph later[Safe to onboard later]");
  for (const s of later) {
    lines.push(`      cat_${s.name.replace(/[^a-zA-Z0-9]/g, "_")}[${s.name} / ${s.stage}]`);
  }
  lines.push("    end");
  lines.push("  end");
  lines.push("  castro -->|owns| cashtro_home");
  lines.push("  castro -->|owns| evolu_home");
  const manager = ok.find((r) => r.fullName === "cashtro/cashtro");
  if (manager) {
    lines.push("  repo_cashtro_cashtro -->|control plane / recon| evolu_home");
    lines.push("  repo_cashtro_cashtro -->|catalog seed| live");
    lines.push("  repo_cashtro_cashtro -->|catalog seed| later");
  }
  lines.push("  evolu_home -.-> evolu_dark");
  lines.push("```", "");
  return lines.join("\n");
}
