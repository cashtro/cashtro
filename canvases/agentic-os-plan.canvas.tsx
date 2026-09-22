/**
 * Visual plan — Vapi voice lane on Cashtro OS.
 * Talk, outbound, voice picker. Every business-line bridge.
 */
export const title = "Agentic OS · Vapi across bridges";

export const mermaid = `
flowchart TB
  subgraph kernel["Cashtro OS kernel"]
    teal["Teal<br/>Teams ingest · intern desk"]
    vapi["Vapi live agentic<br/>talk · call · voice"]
    manager["Manager<br/>ecosystem Links"]
    gate["Human confirm"]
  end
  subgraph voice["Where to change voice"]
    desk["Desk picker"]
    api["POST /api/vapi/voice"]
    dash["dashboard.vapi.ai Assistants"]
    env["VAPI_VOICE_ID"]
    override["per-call assistantOverrides"]
  end
  subgraph lines["Bridge projects"]
    proximity["Proximity"]
    scanapp["Scan App"]
    panda["Panda"]
    nft["NFT + Giant"]
    ecole["École"]
    marketing["Marketing"]
    empire["Empire"]
    trading["Trading"]
    pandora["Pandora / PBTM"]
  end
  subgraph teams["Microsoft Teams · AI Teal"]
    scrape["Graph scrape / ingest"]
    card["Adaptive Card"]
  end
  manager --> vapi
  manager --> teal
  teal --> scrape
  teal --> card
  desk --> vapi
  api --> vapi
  dash --> vapi
  env --> vapi
  override --> vapi
  vapi --> proximity
  vapi --> scanapp
  vapi --> panda
  vapi --> nft
  vapi --> ecole
  vapi --> marketing
  vapi --> empire
  vapi --> trading
  vapi --> pandora
  vapi -->|"POST /call"| gate
  gate -->|"allow → fire"| vapi
`;

export default function AgenticOSPlan() {
  return (
    <section data-canvas="agentic-os-plan">
      <h1>Vapi across every bridge</h1>
      <p>
        Vapi is a kernel voice lane: talk with Castro, park outbound
        calls behind the human gate, switch voices on the desk. Teal
        still ingests Teams. The same script rides Proximity, Scan App,
        Panda, NFT/Giant, école, marketing, Empire, trading, and Pandora.
      </p>
      <pre>{mermaid}</pre>
    </section>
  );
}
