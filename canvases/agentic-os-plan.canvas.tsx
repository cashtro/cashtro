export default function AgenticOsPlan() {
  return (
    <div style={{ padding: 24, fontFamily: "Georgia, serif", maxWidth: 920 }}>
      <h1 style={{ fontSize: 28, marginBottom: 8 }}>Ultron over Cashtro</h1>
      <p style={{ opacity: 0.8, marginBottom: 24 }}>
        Separated control plane. Companies run from one IDE. Cashtro OS stays the
        subordinate kernel.
      </p>
      <pre style={{ whiteSpace: "pre-wrap", lineHeight: 1.5, fontSize: 14 }}>
{`ULTRON IDE  :9090
  auth + RBAC (owner/admin/operator/client/viewer)
  companies: Giant · ScanApp · Proximity · Empire · …
  fleet: agents + workers per company
       │
       ▼ bridge (HTTP)
CASHTRO OS  :8080
  kernel · 14 agentics · delivery · research · disk`}
      </pre>
    </div>
  );
}
