export const title = "Cashtro OS · Giant HUD · Graphify live";

export const mermaid = `
flowchart TB
  player["Castro · Giant Deck"]
  hud["HUD OS :8080<br/>XP · LVL · live ticks"]
  graph["graphify.Build<br/>nodes + edges + flow"]
  giant["GIANT token<br/>world hub"]
  player --> hud --> graph --> giant
  subgraph code["Comment le code se déroule"]
    cmd["cmd/cashtro"] --> boot["agents.Boot"]
    boot --> reg["kernel.Register x15"]
    reg --> kboot["kernel.Boot"]
    kboot --> srv["server.New"]
    srv --> desk["GET / Giant HUD"]
    desk --> inv["kernel.Invoke"]
    inv --> gbuild["manager.graphify"]
    gbuild --> live["GET /api/live"]
    live --> desk
  end
  subgraph brains["Cinq cerveaux"]
    a["Architecte"]
    c["Cartographe · Graphify"]
    f["Forgeron"]
    o["Orfèvre"]
    h["Hustler"]
  end
  graph --> brains
  giant --> nft["NFT + Giant"]
  giant --> trade["Trading bots"]
  giant --> empire["Empire live"]
`;

export default function AgenticOsPlan() {
  return (
    <main style={{ fontFamily: "Syne, Sora, sans-serif", background: "#05060a", color: "#eef1f7", padding: 24 }}>
      <h1 style={{ fontSize: 42, letterSpacing: "-0.04em" }}>VOLTRON · GIANT HUD</h1>
      <p>Live Graphify map in the kernel. Castro watches XP, the five brains, and Giant as the world hub.</p>
      <ol>
        <li>cmd/cashtro boots 15 agentics</li>
        <li>manager.graphify builds the map every tick</li>
        <li>GET /api/live feeds the game desk at 1.2s</li>
        <li>Giant sits at the center of NFT, trading, Empire</li>
      </ol>
    </main>
  );
}
