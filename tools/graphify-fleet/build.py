#!/usr/bin/env python3
"""Build the Graphify-facing agent × project work map.

Reads Cursor cloud-agent dumps plus the kernel catalog, then writes:

- inventory/agent-project-map.json
- inventory/fleet/*.md (wikilinks Graphify and humans can follow)
- internal/fleet/map.go (AST-extractable Worked() edges)
"""

from __future__ import annotations

import json
import re
from collections import defaultdict
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
DUMPS = [
    Path("/home/ubuntu/.cursor/projects/workspace/agent-tools/60402102-cf6d-46fe-82be-438ee91d71ae.txt"),
    Path("/home/ubuntu/.cursor/projects/workspace/agent-tools/d39e0e97-9736-437d-9d7d-609d6b4c4048.txt"),
]

KERNEL_AGENTS = [
    {"id": "init", "name": "Init", "role": "kernel", "mode": "live", "work": "Boots the OS and publishes the manifesto."},
    {"id": "delivery", "name": "Delivery", "role": "ship", "mode": "live", "work": "Owns the idea → concept → production line."},
    {"id": "router", "name": "Router", "role": "model", "mode": "resident", "work": "Optional OpenRouter model bus."},
    {"id": "research", "name": "Research", "role": "library", "mode": "live", "work": "Ingests sourced findings into notes."},
    {"id": "explorer", "name": "Explorer", "role": "search", "mode": "live", "work": "Searches processes, ships, and notes."},
    {"id": "operator", "name": "Operator", "role": "computer-use", "mode": "resident", "work": "Browser and desktop for shippers."},
    {"id": "reviewer", "name": "Reviewer", "role": "qa", "mode": "resident", "work": "Reads walkthrough artifacts before a ship is done."},
    {"id": "architect", "name": "Architect", "role": "design", "mode": "resident", "work": "Shapes apps and workflows before they hit the line."},
    {"id": "deploy", "name": "Deploy", "role": "release", "mode": "resident", "work": "CI, preview, and production promotion."},
    {"id": "security", "name": "Security", "role": "guard", "mode": "resident", "work": "Triage CVE and SAST findings."},
    {"id": "memory", "name": "Memory", "role": "recall", "mode": "live", "work": "Episodic store and recall."},
    {"id": "comms", "name": "Comms", "role": "signal", "mode": "live", "work": "Outbound only after a human confirm."},
    {"id": "planner", "name": "Planner", "role": "backlog", "mode": "live", "work": "Turns a goal into idea-stage ships."},
    {"id": "investigator", "name": "Investigator", "role": "incident", "mode": "resident", "work": "Traces a failing check to blast radius."},
]

PROJECTS = [
    {"slug": "cashtro", "name": "cashtro/cashtro", "kind": "control-plane", "stage": "active", "access": "ok", "notes": "Manager repo. Go kernel + TypeScript control plane."},
    {"slug": "evolu-jeunes", "name": "Evolu-Jeunes/*", "kind": "org", "stage": "denied", "access": "limited", "notes": "48 claimed private repos. Token cannot see them."},
    {"slug": "btk-avocats", "name": "BTK Avocats", "kind": "product", "stage": "production", "access": "catalog-only", "notes": "LIVE CLIENT. Law firm WordPress platform."},
    {"slug": "md-clinic", "name": "MD Clinic", "kind": "product", "stage": "production", "access": "catalog-only", "notes": "LIVE CLIENT. Medical clinic WordPress/ACF."},
    {"slug": "solution-hypotheque-qc", "name": "Solution Hypothèque QC", "kind": "product", "stage": "production", "access": "catalog-only", "notes": "LIVE CLIENT. Mortgage services."},
    {"slug": "educonnexion", "name": "Éduconnexion", "kind": "product", "stage": "production", "access": "catalog-only", "notes": "LIVE CLIENT. Education platform."},
    {"slug": "proximity", "name": "Proximity", "kind": "product", "stage": "production", "access": "catalog-only", "notes": "LIVE CLIENT. Agency Next.js + Azure."},
    {"slug": "scanapp", "name": "ScanApp", "kind": "product", "stage": "concept", "access": "catalog-only", "notes": "Scan, CRM, and AI bots. Safe to onboard later."},
    {"slug": "cashtro-catalog", "name": "Cashtro delivery catalog", "kind": "product", "stage": "idea", "access": "ok", "notes": "Stdlib Go board for idea → concept → production."},
    {"slug": "axel", "name": "Axel brand AI bot", "kind": "idea", "stage": "idea", "access": "named-only", "notes": "Named by a cloud agent. Not in the kernel catalog."},
    {"slug": "happier", "name": "Happier non-profit", "kind": "idea", "stage": "idea", "access": "named-only", "notes": "Named by a cloud agent. Not in the kernel catalog."},
    {"slug": "wealth", "name": "Wealth / trading agents", "kind": "idea", "stage": "idea", "access": "named-only", "notes": "Named by wealth and trading cloud agents."},
    {"slug": "n8n-desk", "name": "Internal n8n desk", "kind": "internal", "stage": "concept", "access": "named-only", "notes": "n8n workflow and desk browser-test work."},
]

CATALOG_SHIPS = [
    "btk-avocats",
    "md-clinic",
    "solution-hypotheque-qc",
    "educonnexion",
    "proximity",
    "scanapp",
    "cashtro-catalog",
]


def ident(s: str) -> str:
    s = s.replace("É", "E").replace("é", "e").replace("à", "a").replace("ô", "o")
    s = re.sub(r"[^A-Za-z0-9]+", " ", s).strip()
    parts = s.split()
    if not parts:
        return "Unnamed"
    out = "".join(p[:1].upper() + p[1:] for p in parts)
    if out[0].isdigit():
        out = "N" + out
    return out


def slug(s: str) -> str:
    s = s.replace("É", "E").replace("é", "e").replace("à", "a").replace("ô", "o")
    s = re.sub(r"[^a-z0-9]+", "-", s.lower()).strip("-")
    return s or "unnamed"


def load_cloud_agents() -> list[dict]:
    agents: list[dict] = []
    seen: set[str] = set()
    for path in DUMPS:
        if not path.exists():
            continue
        text = path.read_text()
        start = text.find("{")
        if start < 0:
            continue
        data = json.loads(text[start:])
        for a in data.get("agents", []):
            bc = a.get("bcId")
            if not bc or bc in seen:
                continue
            seen.add(bc)
            agents.append(a)
    return agents


def classify_projects(name: str) -> list[str]:
    n = name.lower()
    hits: list[str] = ["cashtro"]
    if "scan" in n:
        hits.append("scanapp")
    if "proximity" in n or "teams conversation" in n:
        hits.append("proximity")
    if "axel" in n:
        hits.append("axel")
    if "happier" in n:
        hits.append("happier")
    if "wealth" in n or "trading" in n or "kimmy" in n or "kimi" in n:
        hits.append("wealth")
    if "n8n" in n:
        hits.append("n8n-desk")
    if "voltron" in n or "desk" in n or "agentic" in n or "ai os" in n or "saas ai" in n:
        hits.append("cashtro-catalog")
    if "graphiphy" in n or "graphify" in n or "grapiphy" in n:
        hits.extend(CATALOG_SHIPS)
        hits.append("evolu-jeunes")
        hits.append("axel")
        hits.append("happier")
        hits.append("wealth")
        hits.append("n8n-desk")
    # unique preserve order
    out: list[str] = []
    for h in hits:
        if h not in out:
            out.append(h)
    return out


def kernel_projects(agent_id: str) -> list[str]:
    if agent_id == "delivery":
        return ["cashtro", "cashtro-catalog", *CATALOG_SHIPS]
    if agent_id == "planner":
        return ["cashtro", "cashtro-catalog"]
    if agent_id == "explorer":
        return ["cashtro", *CATALOG_SHIPS]
    if agent_id == "research":
        return ["cashtro"]
    return ["cashtro"]


def parse_dumps(agents: list[dict]) -> dict[str, dict]:
    groups: dict[str, dict] = {}
    for a in agents:
        name = a.get("name") or "Unnamed"
        key = slug(name)
        g = groups.setdefault(
            key,
            {
                "id": key,
                "name": name,
                "kind": "cloud",
                "source": a.get("source"),
                "runs": 0,
                "statuses": defaultdict(int),
                "sources": defaultdict(int),
                "branches": set(),
                "bcIds": [],
                "urls": [],
                "projects": classify_projects(name),
                "work": name,
            },
        )
        g["runs"] += 1
        g["statuses"][a.get("status") or "UNKNOWN"] += 1
        g["sources"][a.get("source") or "unknown"] += 1
        if a.get("branchName"):
            g["branches"].add(a["branchName"])
        g["bcIds"].append(a["bcId"])
        if a.get("url"):
            g["urls"].append(a["url"])
        # prefer a more specific source label if any run is user-facing
        if a.get("source") in ("web", "desktop", "automations") and g["source"] == "internal":
            g["source"] = a.get("source")
    return groups


def go_string(s: str) -> str:
    return json.dumps(s, ensure_ascii=False)


def write_go(projects: list[dict], kernel: list[dict], cloud: list[dict], edges: list[dict]) -> str:
    lines = [
        "// Code generated by tools/graphify-fleet/build.py. DO NOT EDIT.",
        "//",
        "// Package fleet is the Graphify-facing map of every agentic and every",
        "// project they have worked. `Worked(agent, project)` calls are the",
        "// extractable edges. Rebuild with: python3 tools/graphify-fleet/build.py",
        "package fleet",
        "",
        "// Project is a delivery ship, control-plane repo, or named idea.",
        "type Project struct {",
        "\tSlug   string",
        "\tName   string",
        "\tKind   string",
        "\tStage  string",
        "\tAccess string",
        "}",
        "",
        "// Agent is a kernel process or a Cursor cloud run family.",
        "type Agent struct {",
        "\tID   string",
        "\tName string",
        "\tKind string",
        "\tRole string",
        "}",
        "",
        "// Work is one agent × project edge.",
        "type Work struct {",
        "\tAgent   Agent",
        "\tProject Project",
        "}",
        "",
        "// Edges is every Worked() call recorded at init.",
        "var Edges []Work",
        "",
        "// Worked records that agent ran work against project.",
        "// Graphify extracts each call as an edge.",
        "func Worked(agent Agent, project Project) {",
        "\tEdges = append(Edges, Work{Agent: agent, Project: project})",
        "}",
        "",
    ]
    for p in projects:
        fn = ident(p["slug"]) + "Project"
        lines.append(f"// {p['name']}")
        lines.append(f"func {fn}() Project {{")
        lines.append(
            f"\treturn Project{{Slug: {go_string(p['slug'])}, Name: {go_string(p['name'])}, Kind: {go_string(p['kind'])}, Stage: {go_string(p['stage'])}, Access: {go_string(p['access'])}}}"
        )
        lines.append("}")
        lines.append("")
    lines.append("// Kernel agentics registered in internal/agents.")
    for a in kernel:
        fn = ident(a["id"]) + "Agent"
        lines.append(f"func {fn}() Agent {{")
        lines.append(
            f"\treturn Agent{{ID: {go_string(a['id'])}, Name: {go_string(a['name'])}, Kind: \"kernel\", Role: {go_string(a['role'])}}}"
        )
        lines.append("}")
        lines.append("")
    lines.append("// Cloud run families observed on cashtro/cashtro.")
    used: set[str] = set()
    for a in cloud:
        name = ident(a["id"]) + "Agent"
        if name in used:
            name = ident(a["id"] + "Cloud") + "Agent"
        used.add(name)
        a["_go"] = name
        role = a.get("source") or "cloud"
        lines.append(f"func {name}() Agent {{")
        lines.append(
            f"\treturn Agent{{ID: {go_string(a['id'])}, Name: {go_string(a['name'])}, Kind: \"cloud\", Role: {go_string(role)}}}"
        )
        lines.append("}")
        lines.append("")
    lines.append("func init() {")
    proj_ident = {p["slug"]: ident(p["slug"]) + "Project" for p in projects}
    kern_ident = {a["id"]: ident(a["id"]) + "Agent" for a in kernel}
    for e in edges:
        if e["agentKind"] == "kernel":
            ag = kern_ident[e["agent"]]
        else:
            ag = next(c["_go"] for c in cloud if c["id"] == e["agent"])
        pr = proj_ident[e["project"]]
        lines.append(f"\tWorked({ag}(), {pr}()) // {e['work']}")
    lines.append("}")
    lines.append("")
    return "\n".join(lines)


def write_md(title: str, body: str) -> str:
    return f"# {title}\n\n{body.rstrip()}\n"


def main() -> None:
    now = datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")
    raw = load_cloud_agents()
    groups = parse_dumps(raw)
    cloud = sorted(groups.values(), key=lambda g: (-g["runs"], g["name"].lower()))
    for g in cloud:
        g["statuses"] = dict(g["statuses"])
        g["sources"] = dict(g["sources"])
        g["branches"] = sorted(g["branches"])

    edges: list[dict] = []
    for a in KERNEL_AGENTS:
        for p in kernel_projects(a["id"]):
            edges.append({"agent": a["id"], "agentKind": "kernel", "project": p, "work": a["work"], "runs": 1})
    for g in cloud:
        for p in g["projects"]:
            edges.append(
                {
                    "agent": g["id"],
                    "agentKind": "cloud",
                    "project": p,
                    "work": g["name"],
                    "runs": g["runs"],
                }
            )

    payload = {
        "schemaVersion": "1",
        "generatedAt": now,
        "source": "cursor-cloud list-cloud-agents + kernel catalog + named cloud work",
        "reachableRepo": "cashtro/cashtro",
        "counts": {
            "projects": len(PROJECTS),
            "kernelAgents": len(KERNEL_AGENTS),
            "cloudFamilies": len(cloud),
            "cloudRuns": len(raw),
            "workedEdges": len(edges),
        },
        "projects": PROJECTS,
        "kernelAgents": KERNEL_AGENTS,
        "cloudAgents": [
            {
                "id": g["id"],
                "name": g["name"],
                "kind": "cloud",
                "source": g["source"],
                "runs": g["runs"],
                "statuses": g["statuses"],
                "sources": g["sources"],
                "branches": g["branches"],
                "projects": g["projects"],
                "bcIds": g["bcIds"],
                "url": g["urls"][0] if g["urls"] else None,
            }
            for g in cloud
        ],
        "worked": edges,
    }

    (ROOT / "inventory" / "agent-project-map.json").write_text(json.dumps(payload, indent=2) + "\n")

    fleet = ROOT / "inventory" / "fleet"
    for sub in ("projects", "kernel", "cloud"):
        (fleet / sub).mkdir(parents=True, exist_ok=True)

    by_project: dict[str, list[dict]] = defaultdict(list)
    for e in edges:
        by_project[e["project"]].append(e)

    (fleet / "README.md").write_text(
        write_md(
            "Fleet map — agents × work × projects",
            f"""Generated {now}. Graphify extracts `internal/fleet` (`Worked` calls) and this folder.

Reachable GitHub repo: [[cashtro]]. Catalog ships are named only until Option A unlocks [[evolu-jeunes]].

| Layer | Count |
| --- | ---: |
| Projects | {len(PROJECTS)} |
| Kernel agentics | {len(KERNEL_AGENTS)} |
| Cloud run families | {len(cloud)} |
| Cloud runs observed | {len(raw)} |
| Worked edges | {len(edges)} |

## Projects

{chr(10).join(f"- [[{p['slug']}]] — {p['name']} ({p['stage']}, {p['access']})" for p in PROJECTS)}

## Kernel

{chr(10).join(f"- [[{a['id']}]] — {a['name']} · {a['role']} · {a['mode']}" for a in KERNEL_AGENTS)}

## Cloud families

{chr(10).join(f"- [[{g['id']}]] — {g['name']} · {g['runs']} run(s) · {', '.join('[['+p+']]' for p in g['projects'])}" for g in cloud)}
""",
        )
    )

    for p in PROJECTS:
        workers = by_project[p["slug"]]
        kernel_w = [e for e in workers if e["agentKind"] == "kernel"]
        cloud_w = [e for e in workers if e["agentKind"] == "cloud"]
        (fleet / "projects" / f"{p['slug']}.md").write_text(
            write_md(
                p["name"],
                f"""Slug: `{p['slug']}`
Kind: {p['kind']} · Stage: {p['stage']} · Access: {p['access']}

{p['notes']}

Index: [[README]]

## Kernel agents that worked this project

{chr(10).join(f"- [[{e['agent']}]] — {e['work']}" for e in kernel_w) or "_None._"}

## Cloud agents that worked this project

{chr(10).join(f"- [[{e['agent']}]] — {e['work']} ({e['runs']} run(s))" for e in cloud_w) or "_None._"}
""",
            )
        )

    for a in KERNEL_AGENTS:
        projs = kernel_projects(a["id"])
        (fleet / "kernel" / f"{a['id']}.md").write_text(
            write_md(
                a["name"],
                f"""Kernel agentic `{a['id']}` · role `{a['role']}` · mode `{a['mode']}`.

{a['work']}

Index: [[README]]

## Projects this agent worked

{chr(10).join(f"- [[{p}]]" for p in projs)}
""",
            )
        )

    for g in cloud:
        (fleet / "cloud" / f"{g['id']}.md").write_text(
            write_md(
                g["name"],
                f"""Cloud run family `{g['id']}` · source `{g['source']}` · {g['runs']} run(s).

Statuses: {", ".join(f"{k}={v}" for k, v in g["statuses"].items())}

Index: [[README]]

## Projects this agent worked

{chr(10).join(f"- [[{p}]]" for p in g["projects"])}

## Runs

{chr(10).join(f"- `{bc}`" for bc in g["bcIds"][:20])}
{("…" if len(g["bcIds"]) > 20 else "")}
""",
            )
        )

    go_path = ROOT / "internal" / "fleet" / "map.go"
    go_path.parent.mkdir(parents=True, exist_ok=True)
    go_path.write_text(write_go(PROJECTS, KERNEL_AGENTS, cloud, edges))

    print(f"projects={len(PROJECTS)} kernel={len(KERNEL_AGENTS)} cloud={len(cloud)} runs={len(raw)} edges={len(edges)}")


if __name__ == "__main__":
    main()
