package orgmap

import "strings"

// PageHTML is the desk map. Self-contained — no CDN.
func PageHTML() []byte {
	g := Joint()
	var cashtro, evolu strings.Builder
	for _, n := range g.Nodes {
		if n.Org == "cashtro" || n.ID == "cashtro-user" {
			cashtro.WriteString(card(n))
		}
		if n.Org == EvoluOrg || n.ID == "evolu-org" || n.ID == "evolu-member" || n.ID == "evolu-site" || n.ID == "evolu-dark" || n.Kind == "catalog-ship" {
			evolu.WriteString(card(n))
		}
	}
	html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Cashtro × Evolu-Jeunes · fleet map</title>
  <link rel="icon" href="/favicon.svg" type="image/svg+xml">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=IBM+Plex+Mono:wght@400;500&family=Sora:wght@400;500;600&family=Syne:wght@600;700;800&display=swap" rel="stylesheet">
  <style>
    :root {
      --void: #07080c; --panel: #0e1118; --lift: #151925; --line: #262b3a;
      --ink: #eef1f7; --muted: #8b93a7; --hot: #ff4d3d; --gold: #f0c14b;
      --live: #3dffa6; --idea: #6ea8ff;
    }
    * { box-sizing: border-box; }
    html, body { margin: 0; background: var(--void); color: var(--ink); font-family: Sora, system-ui, sans-serif; }
    body { min-height: 100vh; background: linear-gradient(180deg, rgba(255,77,61,0.06), transparent 28%), var(--void); }
    a { color: var(--hot); }
    .wrap { width: min(1440px, calc(100% - 32px)); margin: 0 auto; }
    header.top { display: flex; justify-content: space-between; gap: 24px; align-items: flex-end; padding: 28px 0 18px; border-bottom: 1px solid var(--line); }
    .word { font-family: Syne, sans-serif; font-weight: 800; font-size: clamp(36px, 5vw, 64px); line-height: 0.86; letter-spacing: -0.04em; margin: 0; }
    .word em { font-style: normal; color: var(--hot); }
    .sub { font-family: "IBM Plex Mono", ui-monospace, monospace; font-size: 11px; letter-spacing: 0.16em; text-transform: uppercase; color: var(--muted); margin: 8px 0 0; }
    .lede { max-width: 40rem; color: var(--muted); line-height: 1.5; margin: 10px 0 0; font-size: 14px; }
    .chip { border: 1px solid var(--line); background: var(--panel); padding: 6px 10px; font-family: "IBM Plex Mono", ui-monospace, monospace; font-size: 11px; letter-spacing: 0.04em; text-transform: uppercase; color: inherit; text-decoration: none; display: inline-block; }
    .chip.live { color: var(--live); border-color: rgba(61,255,166,0.35); }
    .chip.hot { color: var(--hot); border-color: rgba(255,77,61,0.4); }
    .stats { display: flex; flex-wrap: wrap; gap: 8px; justify-content: flex-end; max-width: 520px; }
    .castro { margin: 20px 0 0; border: 1px solid var(--hot); background: var(--panel); padding: 16px 18px; display: flex; justify-content: space-between; gap: 16px; flex-wrap: wrap; align-items: center; }
    .castro b { font-family: Syne, sans-serif; font-size: 22px; }
    .arrow { display: grid; grid-template-columns: 1fr auto 1fr; gap: 12px; align-items: center; margin: 14px 0; font-family: "IBM Plex Mono", ui-monospace, monospace; font-size: 11px; letter-spacing: 0.12em; text-transform: uppercase; color: var(--muted); }
    .arrow span { height: 1px; background: var(--line); }
    .islands { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; padding: 0 0 28px; }
    .island { background: var(--panel); border: 1px solid var(--line); padding: 16px; min-width: 0; }
    .island h2 { margin: 0 0 4px; font-family: Syne, sans-serif; font-size: 28px; letter-spacing: -0.03em; }
    .island .tag { font-family: "IBM Plex Mono", ui-monospace, monospace; font-size: 11px; letter-spacing: 0.14em; text-transform: uppercase; color: var(--muted); margin: 0 0 14px; }
    .card { border: 1px solid var(--line); background: var(--lift); padding: 12px; margin-bottom: 8px; }
    .card h3 { margin: 0 0 4px; font-family: Syne, sans-serif; font-size: 16px; }
    .card p { margin: 0; color: var(--muted); font-size: 12px; line-height: 1.4; }
    .card a { font-size: 12px; }
    .badge { font-family: "IBM Plex Mono", ui-monospace, monospace; font-size: 10px; letter-spacing: 0.1em; text-transform: uppercase; margin-bottom: 6px; color: var(--gold); }
    .badge.ok { color: var(--live); }
    .badge.live-client { color: var(--hot); }
    .badge.limited, .badge.catalog-only { color: var(--gold); }
    pre { background: var(--void); border: 1px solid var(--line); padding: 14px; overflow: auto; font-family: "IBM Plex Mono", ui-monospace, monospace; font-size: 11px; line-height: 1.45; color: var(--muted); }
    footer { border-top: 1px solid var(--line); padding: 16px 0 32px; color: var(--muted); font-size: 12px; display: flex; justify-content: space-between; gap: 10px; flex-wrap: wrap; font-family: "IBM Plex Mono", ui-monospace, monospace; letter-spacing: 0.06em; text-transform: uppercase; }
    @media (max-width: 900px) { .islands { grid-template-columns: 1fr; } header.top { flex-direction: column; align-items: flex-start; } .stats { justify-content: flex-start; } }
  </style>
</head>
<body>
  <div class="wrap">
    <header class="top">
      <div>
        <p class="word">CASHT<em>R</em>O × EVOLU</p>
        <p class="sub"><a href="/">Voltron desk</a> · joint fleet map · FR/EN</p>
        <p class="lede">` + esc(g.Summary) + `</p>
      </div>
      <div class="stats">
        <a class="chip" href="https://github.com/cashtro/cashtro">cashtro/cashtro</a>
        <a class="chip" href="https://github.com/Evolu-Jeunes">Evolu-Jeunes</a>
        <a class="chip live" href="https://evolujeunes.ca">evolujeunes.ca</a>
        <a class="chip hot" href="/api/map">GET /api/map</a>
      </div>
    </header>
    <section class="castro">
      <div><b>Castro</b><p class="lede" style="margin:6px 0 0">Owns both GitHub homes. Cashtro is the manager. Evolu-Jeunes is the delivery org.</p></div>
      <div class="stats">
        <span class="chip live">control plane reachable</span>
        <span class="chip hot">client fleet dark</span>
      </div>
    </section>
    <div class="arrow"><span></span>mapped together<span></span></div>
    <main class="islands">
      <section class="island" id="cashtro">
        <h2>Cashtro</h2>
        <p class="tag">GitHub user · public control plane</p>
        ` + cashtro.String() + `
      </section>
      <section class="island" id="evolu-jeunes">
        <h2>Evolu-Jeunes</h2>
        <p class="tag">GitHub org · private delivery fleet</p>
        ` + evolu.String() + `
      </section>
    </main>
    <section class="island" style="margin-bottom:28px">
      <h2 style="font-size:18px">Mermaid</h2>
      <p class="tag">Same graph as docs/FLEET_MAP.md</p>
      <pre id="mermaid">` + esc(Mermaid()) + `</pre>
    </section>
    <footer>
      <span>alejandro@proximityagency.ca</span>
      <span>Live clients stay untouched · Option A for private trees</span>
    </footer>
  </div>
</body>
</html>
`
	return []byte(html)
}

func card(n Node) string {
	badge := n.Access
	url := ""
	if n.URL != "" {
		url = `<div><a href="` + esc(n.URL) + `">` + esc(n.URL) + `</a></div>`
	}
	return `<article class="card" data-id="` + esc(n.ID) + `">
  <div class="badge ` + esc(n.Access) + `">` + esc(badge) + ` · ` + esc(n.Kind) + `</div>
  <h3>` + esc(n.Name) + `</h3>
  <p>` + esc(n.Notes) + `</p>` + url + `
</article>`
}
