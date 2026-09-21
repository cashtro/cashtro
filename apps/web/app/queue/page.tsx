import { api } from "../lib";

type Task = { id: string; title: string; status: string; priority: string };

export default async function Queue() {
  const tasks = (await api<Task[]>("/tasks")) ?? [];
  return (
    <section>
      <h1>Queue</h1>
      <ul>
        {tasks.map((t) => (
          <li key={t.id}>
            {t.title} · {t.status} · {t.priority}
          </li>
        ))}
      </ul>
    </section>
  );
}
