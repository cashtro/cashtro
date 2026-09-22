import type { PrismaClient } from "@prisma/client";
import { readFile } from "node:fs/promises";
import path from "node:path";
import { loadCorporation, probesFor } from "./corporation.js";

export async function renderCoup(
  prisma: PrismaClient,
  taskId: string,
  missing = false,
): Promise<{ status: number; html: string }> {
  const task = await prisma.task.findUnique({ where: { id: taskId } });
  if (!task) {
    return { status: 404, html: chrome({ lang: "fr", title: "Coup", page: "coup", body: "<p>Tâche introuvable.</p>" }) };
  }
  const can = task.status === "todo" || task.status === "doing" || task.status === "blocked";
  const blob = `${task.title} ${task.body}`;
  const probes = probesFor(blob);
  const probeList = probes.map((q) => `<li>${esc(q)}</li>`).join("");
  const alert = missing
    ? `<p class="alert">Questions spécifiques d’abord (avec un ?), puis position, coup court, Steel. Ensuite seulement on joue.</p>`
    : "";
  const form = can
    ? `<form method="post" action="/ui/act/dispatch/${esc(task.id)}">
        <p class="mut">Comprendre avant de positionner. Ces questions visent le trou de contexte, pas le slogan.</p>
        <ul class="mut">${probeList}</ul>
        <label for="question">Questions spécifiques</label>
        <textarea id="question" name="question" required rows="3" placeholder="Réponds aux questions ci-dessus. Il faut un ? — on comble le manque avant d’agir."></textarea>
        <label for="position">Position</label>
        <textarea id="position" name="position" required rows="3" placeholder="Seulement après les questions. Où on est. Ce qui est déjà vrai."></textarea>
        <label for="shortest">Coup le plus court</label>
        <textarea id="shortest" name="shortest" required rows="2" placeholder="Le plus petit geste qui avance."></textarea>
        <label for="opponent">Réponse adverse (Steel)</label>
        <textarea id="opponent" name="opponent" required rows="2" placeholder="Steel contredit. Quelle option plus optimisée survit ?"></textarea>
        <p style="margin-top:16px"><button type="submit">Jouer</button>
        <a href="/ui/queue" style="margin-left:8px;color:inherit">Retour file</a></p>
      </form>`
    : `<p class="mut">Cette tâche n’est plus jouable (${esc(task.status)}).</p>
       <p><a href="/ui/queue" style="color:inherit">Retour file</a></p>`;
  const body = `${alert}
    <p class="mut">Avant chaque coup : des questions spécifiques pour combler le contexte, la position, le coup le plus court, Steel. Ensuite seulement on joue.</p>
    <div class="card"><strong>${esc(task.title)}</strong>
      <div class="mut">${esc(task.status)} · ${esc(task.priority)} · ${esc(task.source)}</div>
      <div class="mut">${esc(task.body || "sans corps")}</div>
    </div>
    ${form}`;
  return { status: 200, html: chrome({ lang: "fr", title: "Coup", page: "coup", body }) };
}

export async function renderPane(prisma: PrismaClient, screen: string): Promise<string> {
  const projects = await prisma.project.findMany();
  const tasks = await prisma.task.findMany({ orderBy: { createdAt: "desc" }, take: 40 });
  const runs = await prisma.run.findMany({ orderBy: { startedAt: "desc" }, take: 20, include: { artifacts: true } });
  const cost = await prisma.run.aggregate({ _sum: { costUsd: true } });
  const paused = await prisma.controlState.findUnique({ where: { id: "global" } });
  const specialistRows = await prisma.worker.findMany({
    where: { tool: { startsWith: "n8n:" } },
    include: { agent: true },
    orderBy: { name: "asc" },
  });
  const specialists = specialistRows.length;
  let report = "_recon report missing_";
  try {
    report = await readFile(path.resolve(process.cwd(), "../../inventory/REPORT.md"), "utf8");
  } catch {
    try {
      report = await readFile(path.resolve(process.cwd(), "inventory/REPORT.md"), "utf8");
    } catch {
      report = "_inventory/REPORT.md not found_";
    }
  }

  const page = ["fleet", "queue", "run", "recon", "corp"].includes(screen) ? screen : "fleet";
  const body =
    page === "fleet"
      ? fleet(projects, cost._sum.costUsd ?? 0, paused?.paused ?? false, specialists, specialistRows)
      : page === "queue"
        ? queue(tasks)
        : page === "run"
          ? runView(runs)
          : page === "corp"
            ? corpView()
            : recon(report);

  return chrome({ lang: "en", title: `Cashtro · ${page}`, page, body });
}

function chrome(opts: { lang: string; title: string; page: string; body: string }): string {
  const { lang, title, page, body } = opts;
  return `<!doctype html>
<html lang="${esc(lang)}">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width,initial-scale=1"/>
  <title>${esc(title)}</title>
  <style>
    :root { color-scheme: dark; --bg:#0b0b0c; --ink:#f4f1ea; --mut:#9a958a; --line:#2a2a2c; --red:#c42b2b; --ok:#3d9a5b; }
    * { box-sizing: border-box; }
    body { margin:0; font: 16px/1.45 ui-sans-serif, system-ui, sans-serif; background:var(--bg); color:var(--ink); }
    header { position:sticky; top:0; background:var(--bg); border-bottom:1px solid var(--line); padding:12px 16px; }
    h1 { font-size:1.1rem; margin:0 0 8px; }
    nav { display:flex; gap:8px; flex-wrap:wrap; }
    nav a { color:var(--ink); text-decoration:none; padding:6px 10px; border:1px solid var(--line); border-radius:999px; font-size:.85rem; }
    nav a[aria-current="page"] { background:var(--red); border-color:var(--red); }
    button { background:var(--red); color:#fff; border:0; border-radius:8px; padding:6px 10px; font:inherit; }
    a.btn { display:inline-block; background:var(--red); color:#fff; text-decoration:none; border-radius:8px; padding:6px 10px; font-size:.85rem; }
    main { padding:16px; max-width:720px; }
    .card { border:1px solid var(--line); border-radius:12px; padding:12px; margin:0 0 10px; }
    .grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(140px,1fr)); gap:8px; margin:12px 0; }
    .chip { border:1px solid var(--line); border-radius:8px; padding:8px; font-size:.78rem; }
    .chip strong { display:block; }
    .dot { display:inline-block; width:.6rem; height:.6rem; border-radius:50%; background:var(--mut); margin-right:6px; }
    .dot.ok { background:var(--ok); }
    .dot.bad { background:var(--red); }
    .mut { color:var(--mut); font-size:.85rem; }
    pre { white-space:pre-wrap; word-break:break-word; font-size:.8rem; }
    table { width:100%; border-collapse:collapse; font-size:.9rem; }
    td, th { text-align:left; padding:6px 4px; border-bottom:1px solid var(--line); }
    label { display:block; font-size:.85rem; margin:12px 0 4px; }
    textarea { width:100%; background:#161618; color:var(--ink); border:1px solid var(--line); border-radius:8px; padding:8px; font:inherit; }
    .alert { border:1px solid var(--red); padding:8px 12px; border-radius:8px; margin:0 0 12px; }
  </style>
</head>
<body>
  <header>
    <h1>Cashtro control plane</h1>
    <nav>
      <a href="/ui/fleet" ${page === "fleet" ? `aria-current="page"` : ""}>Fleet</a>
      <a href="/ui/corp" ${page === "corp" ? `aria-current="page"` : ""}>Corp</a>
      <a href="/ui/queue" ${page === "queue" || page === "coup" ? `aria-current="page"` : ""}>Queue</a>
      <a href="/ui/run" ${page === "run" ? `aria-current="page"` : ""}>Run</a>
      <a href="/ui/recon" ${page === "recon" ? `aria-current="page"` : ""}>Recon</a>
    </nav>
  </header>
  <main>${body}</main>
</body>
</html>`;
}

function fleet(
  projects: Array<{ slug: string; status: string; org: string; kind: string }>,
  cost: number,
  paused: boolean,
  specialists = 0,
  specialistRows: Array<{ name: string; tool: string; agent: { name: string } }> = [],
) {
  const rows = projects
    .map((p) => {
      const ok = p.status === "active" || p.status === "concept";
      return `<div class="card"><span class="dot ${ok ? "ok" : "bad"}"></span><strong>${esc(p.slug)}</strong>
        <div class="mut">${esc(p.org)} · ${esc(p.kind)} · ${esc(p.status)}</div></div>`;
    })
    .join("");
  const chips = specialistRows
    .map(
      (s) =>
        `<div class="chip"><strong>${esc(s.name)}</strong><span class="mut">${esc(s.agent.name)} · ${esc(s.tool.replace(/^n8n:/, ""))}</span></div>`,
    )
    .join("");
  return `<p class="mut">Cost recorded: $${cost.toFixed(4)} · kill switch ${paused ? "ON" : "off"} · n8n specialists ${specialists} · 14 kernel seats + Steel · depth ≤ 3 · fabric inside Voltron, not a second OS</p>
    <form method="post" action="/ui/act/pause"><button type="submit">Pause everything</button></form>
    <h2 style="font-size:1rem">${specialists} specialists · 14 kernel + Steel</h2>
    <div class="grid">${chips || "<p class='mut'>Reseed to load the n8n fabric.</p>"}</div>
    ${rows || "<p>No projects.</p>"}`;
}

function queue(tasks: Array<{ id: string; title: string; status: string; priority: string; source: string }>) {
  if (!tasks.length) return "<p>Queue empty.</p>";
  return `<table><thead><tr><th>Task</th><th>Status</th><th></th></tr></thead><tbody>${tasks
    .map((t) => {
      const can = t.status === "todo" || t.status === "doing" || t.status === "blocked";
      const btn = can ? `<a class="btn" href="/ui/coup/${esc(t.id)}">Coup</a>` : "";
      return `<tr><td>${esc(t.title)}</td><td>${esc(t.status)}</td><td>${btn}</td></tr>`;
    })
    .join("")}</tbody></table>`;
}

function runView(runs: Array<{ id: string; status: string; costUsd: number; exitReason: string | null; artifacts: Array<{ path: string }> }>) {
  if (!runs.length) return "<p>No runs yet. Dispatch a task.</p>";
  return runs
    .map(
      (r) => `<div class="card"><strong>${esc(r.status)}</strong> · $${r.costUsd.toFixed(4)}
        <div class="mut">${esc(r.id)} · ${esc(r.exitReason || "")}</div>
        <div class="mut">${r.artifacts.map((a) => esc(a.path)).join(" · ") || "no artifacts"}</div></div>`,
    )
    .join("");
}

function corpView() {
  const corp = loadCorporation();
  const cards = corp.seats
    .map((d) => {
      const steel = d.id === "contrarian";
      return `<div class="card"${steel ? ` style="border-color:var(--red)"` : ""}>
        <strong>${esc(d.crew)}</strong> · ${esc(d.id)}
        <div class="mut">${esc(d.role)} · ${esc(d.layer)}${d.runtime ? ` · ${esc(d.runtime)}` : ""}</div>
        <p>${esc(d.mission)}</p>
        <p class="mut">Lacune : ${esc(d.gap)}</p>
        <p class="mut">Auto-améliore : ${esc(d.selfImprove)}</p>
        <p class="mut">Max : ${esc(d.maxPush)}</p>
        <ul class="mut">${d.ask.map((q) => `<li>${esc(q)}</li>`).join("")}</ul>
      </div>`;
    })
    .join("");
  return `<p class="mut">${esc(corp.principle)}</p>
    <p class="mut">${corp.kernelSeats} kernel · ${corp.departments} départements · Steel contredit vers l’option plus courte.</p>
    ${cards}`;
}

function recon(md: string) {
  return `<p class="mut">inventory/REPORT.md</p><pre>${esc(md)}</pre>`;
}

function esc(s: string) {
  return s.replace(/[&<>"]/g, (c) => ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;" })[c] || c);
}
