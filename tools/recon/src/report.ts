import type { Inventory, RepoRecord } from "./schema.js";

export function writeReport(inv: Inventory): string {
  const ok = inv.repos.filter((r) => r.access === "ok");
  const blocked = inv.repos.filter((r) => r.access !== "ok");
  const stale = ok.filter((r) => r.riskFlags.includes("stale-12mo") || r.archived);
  const active = ok.filter((r) => !stale.includes(r));
  const lines: string[] = [];
  lines.push(`# Inventory report`);
  lines.push("");
  lines.push(`Generated ${inv.generatedAt}. Control plane: \`${inv.controlPlaneRepo}\`.`);
  lines.push(`Owners scanned: ${inv.owners.join(", ")}.`);
  lines.push("");
  lines.push(`| Bucket | Count |`);
  lines.push(`| --- | ---: |`);
  lines.push(`| Reachable repos | ${ok.length} |`);
  lines.push(`| Denied / limited owners | ${blocked.length} |`);
  lines.push(`| Active | ${active.length} |`);
  lines.push(`| Dead or archivable | ${stale.length} |`);
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
  lines.push(`## What the manager can already control`);
  lines.push("");
  lines.push(
    ok.length
      ? `Only the reachable public surface. Client fleet at Evolu-Jeunes is not in this file until Option A lands.`
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
  const lines = ["# Inventory graph", "", "```mermaid", "flowchart LR", "  manager[cashtro/cashtro]"];
  for (const r of ok) {
    const id = r.name.replace(/[^a-zA-Z0-9]/g, "_");
    lines.push(`  ${id}[${r.fullName}]`);
    lines.push(`  manager --> ${id}`);
  }
  for (const r of inv.repos.filter((x) => x.access !== "ok")) {
    const id = r.org.replace(/[^a-zA-Z0-9]/g, "_") + "_dark";
    lines.push(`  ${id}[${r.org} / denied]`);
    lines.push(`  manager -.-> ${id}`);
  }
  const go = ok.filter((r) => (r.framework || "").includes("Go"));
  const ts = ok.filter((r) => /Next|React|Node/.test(r.framework || ""));
  if (go.length > 1) {
    lines.push(`  subgraph go_stack[Go]`);
    for (const r of go) lines.push(`    ${r.name.replace(/[^a-zA-Z0-9]/g, "_")}`);
    lines.push(`  end`);
  }
  if (ts.length > 1) {
    lines.push(`  subgraph ts_stack[TypeScript]`);
    for (const r of ts) lines.push(`    ${r.name.replace(/[^a-zA-Z0-9]/g, "_")}`);
    lines.push(`  end`);
  }
  lines.push("```", "");
  return lines.join("\n");
}
