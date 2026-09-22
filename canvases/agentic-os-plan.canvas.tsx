export default function AgentOSPlan() {
  const cashtro = [
    { name: "cashtro/cashtro", meta: "control plane · reachable", tone: "ok" },
    { name: "Cashtro OS kernel", meta: "14 agentics on :8080", tone: "ok" },
    { name: "Control plane API", meta: "recon · registry · $2/run", tone: "ok" },
  ];
  const evolu = [
    { name: "Evolu-Jeunes org", meta: "0 public repos · ~48 private", tone: "hot" },
    { name: "evoluJeunes member", meta: "github.com/evoluJeunes", tone: "ok" },
    { name: "evolujeunes.ca", meta: "LIVE nonprofit · do not touch", tone: "hot" },
    { name: "BTK · MD Clinic · Hypothèque · Éduconnexion · Proximity", meta: "production catalog · do not touch", tone: "hot" },
    { name: "ScanApp · Cashtro catalog", meta: "safe to onboard later", tone: "gold" },
  ];
  const card = (item: { name: string; meta: string; tone: string }) => (
    <article
      key={item.name}
      style={{
        border: "1px solid #262b3a",
        background: "#151925",
        padding: 12,
        marginBottom: 8,
      }}
    >
      <div
        style={{
          fontFamily: "ui-monospace, monospace",
          fontSize: 10,
          letterSpacing: "0.1em",
          textTransform: "uppercase",
          color: item.tone === "hot" ? "#ff4d3d" : item.tone === "gold" ? "#f0c14b" : "#3dffa6",
        }}
      >
        {item.tone === "hot" ? "dark / live" : item.tone === "gold" ? "later" : "reachable"}
      </div>
      <h3 style={{ margin: "6px 0 4px", fontSize: 16 }}>{item.name}</h3>
      <p style={{ margin: 0, color: "#8b93a7", fontSize: 12 }}>{item.meta}</p>
    </article>
  );
  return (
    <div style={{ background: "#07080c", color: "#eef1f7", minHeight: "100%", padding: 24, fontFamily: "system-ui, sans-serif" }}>
      <p style={{ fontSize: 11, letterSpacing: "0.16em", textTransform: "uppercase", color: "#8b93a7" }}>
        Admiral fleet · visual plan
      </p>
      <h1 style={{ fontSize: 42, letterSpacing: "-0.04em", margin: "8px 0 12px" }}>
        Cashtro × Evolu-Jeunes
      </h1>
      <p style={{ maxWidth: 640, color: "#8b93a7", lineHeight: 1.5 }}>
        Castro owns both GitHub homes. Cashtro is the public manager (this repo). Evolu-Jeunes is
        the private delivery org. One fleet. Live clients stay untouched until Option A.
      </p>
      <div
        style={{
          margin: "20px 0",
          border: "1px solid #ff4d3d",
          padding: 16,
          display: "flex",
          justifyContent: "space-between",
          gap: 16,
          flexWrap: "wrap",
        }}
      >
        <strong style={{ fontSize: 22 }}>Castro</strong>
        <span style={{ color: "#8b93a7" }}>owns → cashtro user · owns → Evolu-Jeunes org</span>
      </div>
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16 }}>
        <section style={{ background: "#0e1118", border: "1px solid #262b3a", padding: 16 }}>
          <h2 style={{ marginTop: 0 }}>Cashtro</h2>
          {cashtro.map(card)}
        </section>
        <section style={{ background: "#0e1118", border: "1px solid #262b3a", padding: 16 }}>
          <h2 style={{ marginTop: 0 }}>Evolu-Jeunes</h2>
          {evolu.map(card)}
        </section>
      </div>
    </div>
  );
}
