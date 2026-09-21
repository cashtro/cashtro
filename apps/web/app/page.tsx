import { api } from "./lib";

type Project = { slug: string; status: string; org: string; kind: string };

export default async function Fleet() {
  const projects = (await api<Project[]>("/projects")) ?? [];
  const costs = (await api<{ totalUsd: number }>("/costs")) ?? { totalUsd: 0 };
  return (
    <section>
      <h1>Fleet</h1>
      <p>Cost this window: ${costs.totalUsd.toFixed(4)}</p>
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
