// v6 checks: command palette, voice (no API, fake API, refused API), filters, themes, screenshots.
const { chromium } = require('/opt/node-tools/node_modules/playwright'); const fs = require('fs');
const html = fs.readFileSync('centre-pandora.html', 'utf8'), mock = fs.readFileSync('mock2.js', 'utf8');
const page = (pre = '') => `<!doctype html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><script>${pre}</script><script>${mock}</script></head><body>${html}</body></html>`;
const NO_SR = 'delete window.webkitSpeechRecognition; delete window.SpeechRecognition; window.webkitSpeechRecognition = undefined;';
const FAKE_SR = `window.SpeechRecognition = class { start() { setTimeout(() => { this.onresult && this.onresult({ resultIndex: 0, results: [Object.assign([{ transcript: 'bonjour les agents' }], { isFinal: true })] }); setTimeout(() => this.onend && this.onend(), 50); }, 80); } stop() { this.onend && this.onend(); } };`;
const DENY_SR = `window.SpeechRecognition = class { start() { setTimeout(() => this.onerror && this.onerror({ error: 'not-allowed' }), 30); } stop() {} };`;
const out = {}; const errs = [];
async function open(b, pre, vp = { width: 1280, height: 900 }, extra = {}) {
  const p = await b.newPage({ viewport: vp, ...extra }); p.on('pageerror', (e) => errs.push(e.message));
  await p.setContent(page(pre)); await p.evaluate(() => window.__seed()); await p.waitForTimeout(500);
  return p;
}
(async () => {
  const b = await chromium.launch();
  // 1. Guide shows on first visit, then the palette: Ctrl+K, "console", Enter.
  let p = await open(b, NO_SR);
  out.guideShown = await p.isVisible('#guide'); await p.click('#g-next'); await p.click('#g-next'); out.guideLast = await p.textContent('#g-next'); await p.click('#g-next');
  out.guideClosed = !(await p.isVisible('#guide'));
  await p.keyboard.press('Control+k'); out.paletteOpen = await p.isVisible('#cmdk');
  await p.keyboard.type('console'); out.firstHit = await p.textContent('#k-list .k-it[aria-selected="true"]');
  await p.keyboard.press('Enter'); await p.waitForTimeout(300);
  out.paletteSpace = await p.evaluate(() => space); out.navCurrent = await p.getAttribute('nav.spaces a[aria-current="page"]', 'href'); out.paletteClosed = !(await p.isVisible('#cmdk'));
  // Escape closes; agent order command focuses the order box.
  await p.click('#open-k'); await p.keyboard.type('ordre atlas'); await p.keyboard.press('Enter'); await p.waitForTimeout(400);
  out.orderToAtlas = await p.evaluate(() => [document.activeElement.id, document.getElementById('o-who').selectedOptions[0].textContent]);
  await p.keyboard.press('Control+k'); await p.keyboard.press('Escape'); out.escCloses = !(await p.isVisible('#cmdk'));
  // 2. Filters in the Console.
  await p.fill('#f-cq', 'blender'); await p.waitForTimeout(100);
  out.filterConsole = await p.evaluate(() => [document.querySelectorAll('[data-agent-card="c"]:not([hidden])').length, document.getElementById('f-ccount').textContent, netFilter && netFilter.size]);
  await p.fill('#f-cq', ''); await p.selectOption('#f-cteam', 'contra'); out.filterTeam = await p.$$eval('[data-agent-card="c"]:not([hidden])', (n) => n.length);
  // 3. Voice with no API: help panel, then the field gets focus.
  await p.selectOption('#f-cteam', ''); await p.fill('#o-text', 'Déjà écrit');
  await p.click('[data-mic="o-text"]'); out.helpNoApi = await p.isVisible('#michelp');
  await p.click('#mh-ok'); out.focusAfterHelp = await p.evaluate(() => document.activeElement.id);
  // Canvas hover card.
  await p.evaluate(() => document.querySelector('.net canvas').scrollIntoView({ block: 'center' })); await p.waitForTimeout(200); const node = await p.evaluate(() => netNodes[0]); const box = await (await p.$('.net canvas')).boundingBox();
  await p.mouse.move(box.x + node.x, box.y + node.y); await p.waitForTimeout(150);
  out.hoverTip = await p.evaluate(() => !netTip.hidden && netTip.textContent.slice(0, 60));
  await p.screenshot({ path: 'v6-tip.png', clip: { x: box.x, y: box.y, width: box.width, height: box.height } });
  await p.close();
  // 4. Fake API: the dictated text is added after what is already there.
  p = await open(b, FAKE_SR); await p.click('#g-skip'); await p.click('nav a[href="#console"]:visible'); await p.waitForTimeout(200);
  await p.fill('#o-text', 'Déjà écrit'); await p.click('[data-mic="o-text"]'); out.recClass = await p.getAttribute('[data-mic="o-text"]', 'class'); await p.waitForTimeout(400);
  out.dictated = await p.inputValue('#o-text'); out.recStopped = await p.getAttribute('[data-mic="o-text"]', 'class');
  await p.close();
  // 5. Refused API: help panel, remembered for the next click.
  p = await open(b, DENY_SR); await p.click('#g-skip'); await p.click('nav a[href="#console"]:visible'); await p.waitForTimeout(200);
  await p.click('[data-mic="c-topic"]'); await p.waitForTimeout(150); out.helpDenied = await p.isVisible('#michelp'); await p.keyboard.press('Escape');
  out.rememberedSession = await p.evaluate(() => micBlocked);
  await p.close();
  // 6. Screenshots: both themes, desktop and phone, Accueil / Console / Agents.
  for (const theme of ['light', 'dark']) for (const [w, h, tag] of [[1280, 900, 'd'], [390, 844, 'm']]) {
    p = await open(b, NO_SR, { width: w, height: h }, { colorScheme: theme }); await p.click('#g-skip');
    for (const sp of ['accueil', 'console', 'agents']) {
      await p.evaluate((s) => go(s), sp); await p.waitForTimeout(700);
      out[`scrollW-${theme}-${tag}-${sp}`] = await p.evaluate(() => document.documentElement.scrollWidth);
      await p.screenshot({ path: `v6-${sp}-${theme}-${tag}.png`, fullPage: sp !== 'console' || tag === 'm' ? false : false });
      if (sp === 'console') { await (await p.$('.net')).screenshot({ path: `v6-net-${theme}-${tag}.png` }); }
    }
    if (tag === 'd') { await p.evaluate(() => go('accueil')); await p.waitForTimeout(500); await p.screenshot({ path: `v6-accueil-${theme}-full.png`, fullPage: true }); await p.evaluate(() => go('agents')); await p.waitForTimeout(400); await p.screenshot({ path: `v6-agents-${theme}-full.png`, fullPage: true }); }
    await p.close();
  }
  // 7. Theme toggle button and reduced motion.
  p = await open(b, NO_SR, { width: 1280, height: 900 }, { reducedMotion: 'reduce' }); await p.click('#g-skip');
  await p.click('#theme-t'); out.themeAfterToggle = await p.evaluate(() => document.documentElement.dataset.theme);
  await p.evaluate(() => go('console')); await p.waitForTimeout(300); out.loopStoppedReduced = await p.evaluate(() => netLooping === false);
  await p.screenshot({ path: 'v6-console-toggle-dark.png' });
  await p.close();
  console.log(JSON.stringify(out, null, 1)); console.log('errors', JSON.stringify(errs)); await b.close();
})();
