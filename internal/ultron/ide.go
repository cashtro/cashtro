package ultron

// ideHTML is the Ultron IDE desk — runs on the server, no Cursor token.
var ideHTML = []byte(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>ULTRON</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=Archivo+Black&family=IBM+Plex+Mono:wght@400;500&family=Manrope:wght@400;600;700&display=swap" rel="stylesheet">
<style>
:root {
  --ink: #0e1412;
  --paper: #e8efe6;
  --panel: #121a17;
  --line: #24322c;
  --mint: #7dffb3;
  --amber: #ffb84d;
  --fog: #9db5a8;
  --danger: #ff6b6b;
  --glow: rgba(125,255,179,.12);
}
* { box-sizing: border-box; }
html, body { margin: 0; min-height: 100%; background: var(--ink); color: var(--paper); font-family: Manrope, sans-serif; }
body {
  background:
    radial-gradient(1200px 600px at 10% -10%, var(--glow), transparent 55%),
    radial-gradient(900px 500px at 100% 0%, rgba(255,184,77,.08), transparent 50%),
    linear-gradient(165deg, #0a100e 0%, #121a17 45%, #0e1412 100%);
  min-height: 100vh;
}
button, input, select { font: inherit; }
button { cursor: pointer; }
code, .mono { font-family: "IBM Plex Mono", monospace; }

#gate {
  min-height: 100vh; display: grid; place-items: center; padding: 2rem;
}
.gate-card {
  width: min(440px, 100%);
  border-top: 3px solid var(--mint);
  padding: 2.5rem 0 1rem;
}
.brand {
  font-family: "Archivo Black", sans-serif;
  font-size: clamp(3rem, 12vw, 5.5rem);
  letter-spacing: -.04em;
  line-height: .9;
  margin: 0 0 .75rem;
  background: linear-gradient(120deg, var(--paper), var(--mint));
  -webkit-background-clip: text; background-clip: text; color: transparent;
  animation: rise .7s ease both;
}
.tag { color: var(--fog); max-width: 28ch; margin: 0 0 2rem; animation: rise .8s ease both; }
.gate-card label { display: block; font-size: .75rem; letter-spacing: .08em; text-transform: uppercase; color: var(--fog); margin: 0 0 .35rem; }
.gate-card input {
  width: 100%; margin: 0 0 1rem; padding: .85rem 1rem;
  background: var(--panel); border: 1px solid var(--line); color: var(--paper); border-radius: 0;
}
.gate-card button, .primary {
  background: var(--mint); color: var(--ink); border: 0; padding: .9rem 1.25rem; font-weight: 700;
  width: 100%; letter-spacing: .02em;
  transition: transform .15s ease, filter .15s ease;
}
.gate-card button:hover, .primary:hover { filter: brightness(1.05); transform: translateY(-1px); }
.err { color: var(--danger); min-height: 1.25rem; font-size: .9rem; }

#ide { display: none; min-height: 100vh; grid-template-columns: 260px 1fr; }
#ide.on { display: grid; animation: fade .35s ease; }
@media (max-width: 900px) {
  #ide.on { grid-template-columns: 1fr; }
  aside { border-bottom: 1px solid var(--line); }
}

aside {
  border-right: 1px solid var(--line);
  background: rgba(18,26,23,.85);
  backdrop-filter: blur(8px);
  padding: 1.25rem 1rem 2rem;
  display: flex; flex-direction: column; gap: 1rem;
}
.aside-brand {
  font-family: "Archivo Black", sans-serif;
  font-size: 1.6rem; letter-spacing: -.03em; margin: 0;
}
.aside-brand span { color: var(--mint); }
.me { font-size: .85rem; color: var(--fog); }
.me strong { color: var(--paper); display: block; }
.nav-label { font-size: .7rem; letter-spacing: .12em; text-transform: uppercase; color: var(--fog); margin: .5rem 0 .25rem; }
.company {
  display: block; width: 100%; text-align: left;
  background: transparent; border: 0; border-left: 2px solid transparent;
  color: var(--paper); padding: .55rem .65rem; margin: 0 0 .2rem;
  transition: background .15s, border-color .15s;
}
.company:hover, .company.active { background: rgba(125,255,179,.08); border-left-color: var(--mint); }
.company small { display: block; color: var(--fog); font-size: .75rem; }
.ghost {
  background: transparent; border: 1px solid var(--line); color: var(--fog);
  padding: .55rem .8rem; width: 100%;
}
.ghost:hover { color: var(--paper); border-color: var(--fog); }

main { padding: 1.5rem clamp(1rem, 3vw, 2.5rem) 3rem; }
.top {
  display: flex; flex-wrap: wrap; gap: 1rem; align-items: end; justify-content: space-between;
  margin-bottom: 1.75rem; padding-bottom: 1rem; border-bottom: 1px solid var(--line);
}
.top h2 {
  font-family: "Archivo Black", sans-serif; font-size: clamp(1.8rem, 4vw, 2.6rem);
  margin: 0; letter-spacing: -.03em;
}
.top p { margin: .35rem 0 0; color: var(--fog); max-width: 48ch; }
.pills { display: flex; flex-wrap: wrap; gap: .5rem; }
.pill {
  font-family: "IBM Plex Mono", monospace; font-size: .75rem;
  padding: .35rem .6rem; border: 1px solid var(--line); color: var(--fog);
}
.pill.ok { color: var(--mint); border-color: rgba(125,255,179,.35); }
.pill.warn { color: var(--amber); border-color: rgba(255,184,77,.35); }

.stage {
  display: grid; gap: 1.5rem;
  grid-template-columns: 1.4fr .9fr;
}
@media (max-width: 980px) { .stage { grid-template-columns: 1fr; } }

.block h3 {
  margin: 0 0 .75rem; font-size: .75rem; letter-spacing: .14em;
  text-transform: uppercase; color: var(--fog); font-weight: 600;
}
.row {
  display: grid; grid-template-columns: 1fr auto; gap: .75rem; align-items: center;
  padding: .85rem 0; border-bottom: 1px solid var(--line);
}
.row:last-child { border-bottom: 0; }
.row strong { display: block; }
.row span { color: var(--fog); font-size: .85rem; }
.row button {
  background: transparent; border: 1px solid var(--mint); color: var(--mint);
  padding: .4rem .7rem; font-size: .8rem; font-weight: 600;
}
.row button:hover { background: rgba(125,255,179,.12); }

.console {
  margin-top: 1.5rem; border-top: 1px solid var(--line); padding-top: 1rem;
}
.console pre {
  margin: 0; padding: 1rem; background: #0a100e; border: 1px solid var(--line);
  overflow: auto; max-height: 280px; font-size: .8rem; color: var(--mint);
}

@keyframes rise { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: none; } }
@keyframes fade { from { opacity: 0; } to { opacity: 1; } }
</style>
</head>
<body>
<section id="gate">
  <div class="gate-card">
    <h1 class="brand">ULTRON</h1>
    <p class="tag">Personal IDE OS. Run every company from one authorized desk. Cashtro is the kernel underneath — no Cursor token on the server.</p>
    <form id="loginForm">
      <label for="email">Email</label>
      <input id="email" name="email" type="email" autocomplete="username" value="alejandro@proximityagency.ca" required>
      <label for="password">Password</label>
      <input id="password" name="password" type="password" autocomplete="current-password" value="ultron-change-me" required>
      <p class="err" id="loginErr"></p>
      <button type="submit">Enter Ultron</button>
    </form>
  </div>
</section>

<div id="ide">
  <aside>
    <p class="aside-brand">ULTRON <span>·</span></p>
    <div class="me" id="meBox"></div>
    <div>
      <div class="nav-label">Companies</div>
      <div id="companyList"></div>
    </div>
    <div style="margin-top:auto; display:grid; gap:.5rem">
      <button class="ghost" type="button" id="btnOsPulse">Pulse Cashtro OS</button>
      <button class="ghost" type="button" id="btnLogout">Sign out</button>
    </div>
  </aside>
  <main>
    <div class="top">
      <div>
        <h2 id="title">Fleet</h2>
        <p id="subtitle">Select a company. Agents and workers appear here. Invokes bridge into Cashtro when bound.</p>
      </div>
      <div class="pills" id="pills"></div>
    </div>
    <div class="stage">
      <div class="block">
        <h3>Agents</h3>
        <div id="agents"></div>
      </div>
      <div class="block">
        <h3>Workers</h3>
        <div id="workers"></div>
        <h3 style="margin-top:1.5rem">Cashtro bridge</h3>
        <div id="bridge"></div>
      </div>
    </div>
    <div class="console"><h3 class="nav-label" style="margin:0 0 .5rem">Console</h3><pre id="out">awaiting command…</pre></div>
  </main>
</div>

<script>
const state = { token: localStorage.getItem('ultron_token') || '', user: null, companies: [], selected: '', about: null };

async function api(path, opts = {}) {
  const headers = Object.assign({ 'Content-Type': 'application/json' }, opts.headers || {});
  if (state.token) headers['Authorization'] = 'Bearer ' + state.token;
  const res = await fetch(path, Object.assign({}, opts, { headers }));
  const text = await res.text();
  let data; try { data = JSON.parse(text); } catch { data = { raw: text }; }
  if (!res.ok) throw new Error((data && data.error) || res.statusText);
  return data;
}

function log(v) {
  document.getElementById('out').textContent = typeof v === 'string' ? v : JSON.stringify(v, null, 2);
}

async function bootAbout() {
  state.about = await api('/api/about');
  const pills = document.getElementById('pills');
  pills.innerHTML = [
    pill(state.about.version, true),
    pill(state.about.companies + ' companies', true),
    pill(state.about.agents + ' agents', true),
    pill(state.about.cashtroOk ? 'cashtro online' : 'cashtro down', state.about.cashtroOk)
  ].join('');
}

function pill(t, ok) {
  return '<span class="pill ' + (ok ? 'ok' : 'warn') + '">' + t + '</span>';
}

async function enter() {
  document.getElementById('gate').style.display = 'none';
  document.getElementById('ide').classList.add('on');
  state.user = await api('/api/me');
  document.getElementById('meBox').innerHTML = '<strong>' + state.user.name + '</strong>' + state.user.role + ' · ' + state.user.email;
  await bootAbout();
  state.companies = await api('/api/companies');
  renderCompanies();
  if (state.companies.length) selectCompany(state.companies[0].id);
}

function renderCompanies() {
  const box = document.getElementById('companyList');
  box.innerHTML = state.companies.map(c => (
    '<button type="button" class="company' + (c.id === state.selected ? ' active' : '') + '" data-id="' + c.id + '">' +
    c.name + '<small>' + c.sector + ' · ' + c.status + '</small></button>'
  )).join('');
  box.querySelectorAll('.company').forEach(btn => btn.addEventListener('click', () => selectCompany(btn.dataset.id)));
}

async function selectCompany(id) {
  state.selected = id;
  renderCompanies();
  const data = await api('/api/companies/' + id);
  const c = data.company;
  document.getElementById('title').textContent = c.name;
  document.getElementById('subtitle').textContent = c.notes || (c.sector + ' · ' + (c.stack || []).join(' · '));
  document.getElementById('agents').innerHTML = (data.agents || []).map(a => (
    '<div class="row"><div><strong>' + a.name + '</strong><span class="mono">' + a.mode + ' · ' + (a.capabilities || []).join(', ') +
    (a.cashtroId ? ' → cashtro:' + a.cashtroId : '') + '</span></div>' +
    '<button type="button" data-agent="' + a.id + '">Invoke</button></div>'
  )).join('') || '<p class="me">No agents bound.</p>';
  document.getElementById('workers').innerHTML = (data.workers || []).map(w => (
    '<div class="row"><div><strong>' + w.name + '</strong><span>' + w.kind + ' · ' + (w.bound ? 'bound' : 'unbound') + '</span></div></div>'
  )).join('') || '<p class="me">No workers.</p>';
  document.querySelectorAll('[data-agent]').forEach(btn => btn.addEventListener('click', () => invoke(btn.dataset.agent)));
  try {
    const os = await api('/api/os/health');
    document.getElementById('bridge').innerHTML = '<div class="row"><div><strong>Cashtro OS</strong><span class="mono">' +
      (os.version || '') + ' · live ' + (os.live ?? '?') + ' · ' + (os.status || '') + '</span></div></div>';
  } catch (e) {
    document.getElementById('bridge').innerHTML = '<p class="err">' + e.message + '</p>';
  }
}

async function invoke(id) {
  try {
    const res = await api('/api/fleet/agents/' + id + '/invoke', { method: 'POST', body: JSON.stringify({ target: state.selected }) });
    log(res);
  } catch (e) { log(String(e.message || e)); }
}

document.getElementById('loginForm').addEventListener('submit', async (e) => {
  e.preventDefault();
  const err = document.getElementById('loginErr');
  err.textContent = '';
  try {
    const data = await api('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email: email.value, password: password.value })
    });
    state.token = data.token;
    localStorage.setItem('ultron_token', data.token);
    await enter();
  } catch (ex) { err.textContent = ex.message; }
});

document.getElementById('btnLogout').addEventListener('click', async () => {
  try { await api('/api/auth/logout', { method: 'POST' }); } catch {}
  localStorage.removeItem('ultron_token');
  location.reload();
});

document.getElementById('btnOsPulse').addEventListener('click', async () => {
  try { log(await api('/api/os/heartbeat', { method: 'POST', body: '{}' })); await bootAbout(); }
  catch (e) { log(String(e.message || e)); }
});

(async () => {
  if (!state.token) return;
  try { await enter(); } catch { localStorage.removeItem('ultron_token'); state.token = ''; }
})();
</script>
</body>
</html>
`)
