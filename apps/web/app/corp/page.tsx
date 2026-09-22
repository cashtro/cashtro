import { api } from "../lib";

const API = process.env.CONTROL_PLANE_URL || "http://127.0.0.1:8787";

type Dept = {
  id: string;
  crew: string;
  role: string;
  layer: string;
  mission: string;
  gap: string;
  ask: string[];
  selfImprove: string;
  maxPush: string;
};

export default async function CorpPage() {
  const corp = (await api<{ principle: string; kernelSeats: number; departments: number; seats: Dept[] }>("/corp")) ?? {
    principle: "",
    kernelSeats: 14,
    departments: 0,
    seats: [],
  };
  return (
    <section>
      <h1>Corporation</h1>
      <p style={{ color: "#9a958a" }}>{corp.principle}</p>
      <p>
        {corp.kernelSeats} kernel · {corp.departments} départements ·{" "}
        <a href={`${API}/ui/corp`} style={{ color: "inherit" }}>
          pane
        </a>
      </p>
      <ul style={{ listStyle: "none", padding: 0 }}>
        {corp.seats.map((d) => (
          <li key={d.id} style={{ border: "1px solid #2a2a2c", borderRadius: 12, padding: 12, marginBottom: 10 }}>
            <strong>{d.crew}</strong> · {d.id} · {d.role} · {d.layer}
            <p>{d.mission}</p>
            <p style={{ color: "#9a958a", fontSize: 14 }}>Lacune : {d.gap}</p>
            <p style={{ color: "#9a958a", fontSize: 14 }}>Auto-améliore : {d.selfImprove}</p>
            <ul>
              {d.ask.map((q) => (
                <li key={q}>{q}</li>
              ))}
            </ul>
          </li>
        ))}
      </ul>
    </section>
  );
}
