import { api } from "./lib";

type Project = { slug: string; status: string; org: string; kind: string };
type Fleet = {
  seats: number;
  kernelSeats?: number;
  specialists: number;
  depthLimit: number;
  n8n: string;
  fabric?: string;
  agents?: Array<{ name: string; workers: Array<{ name: string; tool: string }> }>;
};

export default async function FleetPage() {
  const projects = (await api<Project[]>("/projects")) ?? [];
  const costs = (await api<{ totalUsd: number }>("/costs")) ?? { totalUsd: 0 };
  const fleet = (await api<Fleet>("/fleet")) ?? { seats: 14, specialists: 0, depthLimit: 3, n8n: "unknown", agents: [] };
  const specialists = (fleet.agents ?? []).flatMap((a) =>
    a.workers.filter((w) => w.tool.startsWith("n8n:")).map((w) => ({ seat: a.name, ...w })),
  );
  return (
    <section>
      <h1>Fleet</h1>
      <p>
        Cost this window: ${costs.totalUsd.toFixed(4)} · {fleet.specialists} n8n specialists · {fleet.kernelSeats ?? 14} kernel
        + Steel · depth ≤ {fleet.depthLimit} · {fleet.n8n}
      </p>
      <ul style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(140px, 1fr))", gap: 8, padding: 0, listStyle: "none" }}>
        {specialists.map((s) => (
          <li key={s.name} style={{ border: "1px solid #2a2a2c", borderRadius: 8, padding: 8 }}>
            <strong>{s.name}</strong>
            <div style={{ color: "#9a958a", fontSize: 12 }}>
              {s.seat} · {s.tool.replace(/^n8n:/, "")}
            </div>
          </li>
        ))}
      </ul>
      <ul>
        {projects.map((p) => (
          <li key={p.slug}>
            {p.slug} · {p.status} · {p.org}
          </li>
        ))}
      </ul>
    </section>
  );
}
