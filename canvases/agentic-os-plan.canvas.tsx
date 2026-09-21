export const title = "Cashtro OS · owned n8n";

export const mermaid = `
flowchart TB
  desk["Desk / :8080"] --> workflow["POST /api/agents/n8n/workflow"]
  desk --> hook["POST /api/n8n/{id}/hook"]
  workflow --> n8n["n8n agentic"]
  hook --> n8n
  n8n --> flow["internal/flow graph"]
  flow --> think["internal/think dual"]
  think --> propose["propose seat<br/>owned · not Moonshot"]
  think --> critique["critique seat<br/>owned · not GLM"]
  flow --> verbs["kernel verbs"]
  verbs --> memory["memory.store"]
  verbs --> note["research.ingest"]
  verbs --> confirm["comms.send → human gate"]
  verbs --> ship["delivery.create"]
  n8n -.->|"never reads"| vendors["N8N_WEBHOOK_URL<br/>MOONSHOT_API_KEY<br/>GLM_API_KEY"]
`;

export default function AgenticOSPlan() {
  return (
    <article data-canvas="agentic-os-plan">
      <h1>Cashtro OS · we own n8n</h1>
      <p>
        Castro asked for <code>/api/agents/n8n/workflow</code> without connecting
        to n8n Cloud, Moonshot, or GLM. The graph, the catch-hook, and the dual
        loop are processes on this kernel.
      </p>
      <section>
        <h2>Fleet</h2>
        <ul>
          <li>kernel — process table, journal, mail, notes, memory, confirms</li>
          <li>delivery — idea → concept → production</li>
          <li>n8n — owned workflow runtime</li>
          <li>think — owned propose / critique</li>
          <li>router — optional OpenRouter, unbound by default</li>
        </ul>
      </section>
      <section>
        <h2>Seeded graphs</h2>
        <ol>
          <li>wealth-dual — manual → dual → memory → note</li>
          <li>intake-hook — webhook → set → if → memory → confirm</li>
          <li>research-line — manual → note → memory</li>
        </ol>
      </section>
    </article>
  );
}
