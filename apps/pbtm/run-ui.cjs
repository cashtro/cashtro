const { chromium } = require('/opt/node-tools/node_modules/playwright');
const fs = require('fs');
(async () => {
  const html = fs.readFileSync('centre-pandora.html', 'utf8');
  const page = `<!doctype html><html><head><meta charset=utf-8><meta name=viewport content="width=device-width,initial-scale=1"><script>${fs.readFileSync('mock-claude.js','utf8')}</script></head><body>${html}</body></html>`;
  const b = await chromium.launch(); const p = await b.newPage({ viewport: { width: 1280, height: 900 } });
  const errs = []; p.on('pageerror', (e) => errs.push(e.message)); p.on('console', (m) => m.type() === 'error' && !/ERR_FAILED|ERR_CERT/.test(m.text()) && errs.push(m.text()));
  
  await p.setContent(page); await p.evaluate(() => window.__seed()); await p.waitForTimeout(400); await p.click('#g-skip', { timeout: 3000 }).catch(() => console.log('no guide'));
  await p.click('nav a[href="#console"]:visible'); await p.waitForTimeout(300);
  await p.click('#tab-salles'); await p.click('.agent >> nth=0'); await p.waitForTimeout(200); await p.click('#tab-ordres');
  await p.fill('#o-text', 'Écris la page d accueil de la clinique'); await p.selectOption('#o-client', { label: 'Clinique Test' });
  await p.click('text=Envoyer l\'ordre');
  await p.waitForTimeout(350); await p.fill('#o-text', 'je tape pendant que ça travaille'); await p.waitForTimeout(250);
  await p.screenshot({ path: 'console-working.png', fullPage: true });
  const typed = await p.inputValue('#o-text');
  await p.waitForTimeout(1200);
  await p.screenshot({ path: 'console-done.png', fullPage: true });
  console.log('typed kept:', typed, '| feed items', await p.$$eval('.feed div', (n) => n.length), '| done orders', await p.locator('.pill.ok:text("fait")').count());
  // team order + swarm logging
  await p.selectOption('#o-who', 't:infra'); await p.fill('#o-text', 'Audit rapide'); await p.click('text=Envoyer l\'ordre'); await p.waitForTimeout(1500);
  console.log('orders total', await p.evaluate(() => document.querySelectorAll('details.deliv').length));
  await p.setViewportSize({ width: 390, height: 800 }); await p.waitForTimeout(200);
  await p.screenshot({ path: 'console-mobile.png', fullPage: false }); console.log('scrollW mobile', await p.evaluate(() => document.documentElement.scrollWidth)); await p.click('nav a[href="#accueil"]:visible'); await p.waitForTimeout(300); await p.setViewportSize({ width: 1280, height: 900 }); await p.screenshot({ path: 'home.png' }); console.log('scrollW', await p.evaluate(() => document.documentElement.scrollWidth), 'errors', JSON.stringify(errs));
  await b.close();
})();
