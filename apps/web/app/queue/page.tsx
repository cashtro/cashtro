import { api } from "../lib";

type Task = { id: string; title: string; status: string; priority: string };

const API = process.env.CONTROL_PLANE_URL || "http://127.0.0.1:8787";

export default async function Queue() {
  const tasks = (await api<Task[]>("/tasks")) ?? [];
  return (
    <section>
      <h1>Queue</h1>
      <p style={{ color: "#9a958a", fontSize: 14 }}>
        Avant chaque coup : une question, la position, le coup le plus court, la réponse adverse.
      </p>
      <ul>
        {tasks.map((t) => {
          const can = t.status === "todo" || t.status === "doing" || t.status === "blocked";
          return (
            <li key={t.id}>
              {t.title} · {t.status} · {t.priority}
              {can ? (
                <>
                  {" "}
                  ·{" "}
                  <a href={`${API}/ui/coup/${t.id}`} style={{ color: "inherit" }}>
                    Coup
                  </a>
                </>
              ) : null}
            </li>
          );
        })}
      </ul>
    </section>
  );
}
