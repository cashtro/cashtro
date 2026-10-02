// "Activer la salle": every agent of the room lights up at once, at most PARALLEL call Claude together,
// all finish, and a rate_limited error is retried. Uses mock2 (37 agents) plus 9 test agents in Marketing.
const { chromium } = require('/opt/node-tools/node_modules/playwright'); const fs = require('fs');
const html = fs.readFileSync('centre-pandora.html', 'utf8'), mock = fs.readFileSync('mock2.js', 'utf8');
(async () => {
  const b = await chromium.launch(); const errs = []; const out = {};
  for (const [w, h, tag] of [[1280, 900, 'd'], [390, 844, 'm']]) {
    const p = await b.newPage({ viewport: { width: w, height: h } }); p.on('pageerror', (e) => errs.push(e.message));
    await p.setContent(`<!doctype html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><script>${mock}</script></head><body>${html}</body></html>`);
    await p.evaluate(() => window.__seed()); await p.waitForTimeout(400); await p.click('#g-skip');
    if (tag === 'd') {
      // 9 marketing agents to stress a big room.
      await p.evaluate(async () => { const db = await claude.use('db'); for (const n of ['Écho2', 'Signal2', 'Vitrine', 'Lien', 'Boost', 'Clic', 'Ancre', 'Ruche', 'Phare']) await db.collection('agents').add({ name: n, team: 'marketing', role: 'Spécialiste ' + n, method: 'Étape 1\nÉtape 2', engine: 'claude', active: true, status: 'libre' }); });
      await p.waitForTimeout(300);
      // One rate_limited error on the very first call: must be retried (after 2 s), not fail.
      await p.evaluate(() => { const real = sample; let first = true; window.__fly = 0; window.__flyPeak = 0; window.__retried = 0; sample = Object.assign(async (...a) => { if (first) { first = false; window.__retried++; throw { code: 'rate_limited' }; } window.__fly++; window.__flyPeak = Math.max(window.__flyPeak, window.__fly); try { return await real(...a); } finally { window.__fly--; } }, { json: real.json }); window.__peak = { busy: 0, run: 0 }; setInterval(() => { window.__peak.busy = Math.max(window.__peak.busy, S.agents.filter((x) => x.status === 'travail').length); window.__peak.run = Math.max(window.__peak.run, Object.keys(liveOut).filter((k) => !k.startsWith('d:')).length); }, 30); });
      await p.click('nav a[href="#console"]:visible'); await p.click('#tab-salles'); await p.waitForTimeout(200);
      await p.click('[data-group="c"][data-team="marketing"] >> text=Activer la salle');
      out.panelOpen = await p.isVisible('#activate'); out.panelCount = await p.textContent('#x-count');
      await p.fill('#x-text', 'Chacun propose 3 idées de campagne.'); await p.click('#x-go');
      await p.waitForTimeout(250);
      out.litAtOnce = await p.evaluate(() => S.agents.filter((x) => x.team === 'marketing' && x.status === 'travail').length);
      out.badge = await p.textContent('.live-badge'); out.tab = await p.evaluate(() => ctab);
      await p.screenshot({ path: 'v7-activate-live-d.png' });
      await p.locator('.net canvas').screenshot({ path: 'v7-net-busy-d.png' });
      await p.waitForFunction(() => S.orders.filter((o) => o.team === 'marketing' && o.status === 'fait').length === 9, null, { timeout: 30000 });
      await p.waitForTimeout(400);
      out.done = await p.evaluate(() => S.orders.filter((o) => o.team === 'marketing' && o.status === 'fait').length);
      out.errors = await p.evaluate(() => S.orders.filter((o) => o.status === 'erreur').length);
      out.peak = await p.evaluate(() => ({ ...window.__peak, claudeCallsAtOnce: window.__flyPeak, rateLimitedRetried: window.__retried })); out.busyAfter = await p.evaluate(() => S.agents.filter((x) => x.status === 'travail').length);
      out.toast = await p.evaluate(() => [...document.querySelectorAll('.toast')].map((t) => t.textContent).join(' | '));
      // Second run, no rate limit: exactly PARALLEL (6) Claude calls at the same time.
      await p.evaluate(() => { window.__flyPeak = 0; });
      await p.click('#tab-salles'); await p.click('[data-group="c"][data-team="marketing"] >> text=Activer la salle'); await p.fill('#x-text', 'Deuxième mission.'); await p.click('#x-go');
      await p.waitForFunction(() => S.orders.filter((o) => o.team === 'marketing' && o.status === 'fait').length === 18, null, { timeout: 30000 });
      out.secondRunCallsAtOnce = await p.evaluate(() => window.__flyPeak);
      await p.click('#tab-salles'); await p.screenshot({ path: 'v7-salles-d.png', fullPage: true });
      for (const t of ['ordres', 'conseil', 'journal']) { await p.click('#tab-' + t); await p.waitForTimeout(150); await p.screenshot({ path: `v7-${t}-d.png`, fullPage: true }); }
      await p.click('nav a[href="#accueil"]:visible'); await p.waitForTimeout(400); await p.screenshot({ path: 'v7-accueil-d.png', fullPage: true });
      await p.evaluate(() => go('console')); await p.click('#tab-reseau'); await p.waitForTimeout(400); await p.screenshot({ path: 'v7-console-d.png', fullPage: true });
    } else {
      await p.click('nav a[href="#console"]:visible'); await p.waitForTimeout(300);
      out.ctaMobile = await p.isVisible('.cta-float .btn'); out.scrollW = await p.evaluate(() => document.documentElement.scrollWidth);
      await p.screenshot({ path: 'v7-console-m.png' }); await p.screenshot({ path: 'v7-console-m-full.png', fullPage: true });
      await p.click('#tab-salles'); await p.waitForTimeout(200); await p.screenshot({ path: 'v7-salles-m.png', fullPage: false });
      await p.click('.cta-float .btn'); await p.waitForTimeout(200); await p.screenshot({ path: 'v7-activate-m.png' }); await p.keyboard.press('Escape');
      await p.click('nav a[href="#accueil"]:visible'); await p.waitForTimeout(300); await p.screenshot({ path: 'v7-accueil-m-full.png', fullPage: true });
      out.scrollWHome = await p.evaluate(() => document.documentElement.scrollWidth);
    }
    await p.close();
  }
  console.log(JSON.stringify(out, null, 1)); console.log('errors', JSON.stringify(errs)); await b.close();
})();
