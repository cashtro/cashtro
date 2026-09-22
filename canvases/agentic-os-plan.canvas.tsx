/**
 * Visual plan — Teal corporate brain on Cashtro OS.
 * Teams only. Vapi voice. Cursor cards. Pandora memory.
 */
export const title = "Agentic OS · Teal × Vapi × Cursor";

export const mermaid = `
flowchart TB
  subgraph teams["Microsoft Teams · AI Teal only"]
    scrape["Teams scraper<br/>Graph or ingest"]
    card["Adaptive Card"]
    intern["Intern desk"]
    questions["Questions"]
  end
  subgraph voice["Vapi / Vappy"]
    share["vapi.ai · dashboard.vapi.ai"]
    webhook["/api/vapi/webhook"]
    calls["Listen to calls"]
  end
  subgraph brain["Teal brain · Cashtro OS"]
    memory["Pandora brainstorming"]
    ideas["Generate ideas"]
    spawn["Spawn intel agentics"]
    smarter["Always get smarter"]
  end
  subgraph cursor["Cursor app"]
    agents["cursor.com/agents"]
    gate["Human confirm"]
  end
  scrape --> memory
  webhook --> calls
  calls --> memory
  memory --> ideas
  memory --> questions
  intern --> questions
  questions --> spawn
  ideas --> spawn
  smarter --> spawn
  card --> agents
  agents --> gate
  share --> webhook
  spawn --> brain
`;

export default function AgenticOSPlan() {
  return (
    <section data-canvas="agentic-os-plan">
      <h1>Teal corporate brain</h1>
      <p>
        Microsoft Teams AI Teal is the only conversation surface. Vapi is the
        voice. Cursor launches from the Teams card. Pandora memory stays in
        this kernel.
      </p>
      <pre>{mermaid}</pre>
    </section>
  );
}
