// Fake viewer runtime: in-memory db with live snapshots + a sample.json that answers per prompt kind.
(() => {
  const store = {}; const subs = {};
  const snap = (col) => ({ docs: Object.entries(store[col] || {}).map(([id, d]) => ({ id, exists: true, data: () => d })) });
  const emit = (col) => (subs[col] || []).forEach((f) => setTimeout(() => f(snap(col)), 0));
  let n = 0;
  const db = {
    collection: (col) => ({
      onSnapshot: (f) => { (subs[col] ||= []).push(f); setTimeout(() => f(snap(col)), 0); return () => {}; },
      add: async (d) => { const id = 'id' + (++n); (store[col] ||= {})[id] = d; emit(col); return { id }; },
    }),
    doc: (path) => { const [col, id] = path.split('/'); return {
      onSnapshot: (f) => { const fire = () => f({ exists: !!(store[col] && store[col][id]), data: () => (store[col] || {})[id] }); (subs[col] ||= []).push(fire); setTimeout(fire, 0); return () => {}; },
      update: async (p) => { const merge = (a, b) => { for (const [k, v] of Object.entries(b)) { a[k] = v && typeof v === 'object' && !Array.isArray(v) && a[k] && typeof a[k] === 'object' ? merge({ ...a[k] }, v) : v; } return a; }; store[col][id] = merge({ ...(store[col][id] || {}) }, JSON.parse(JSON.stringify(p))); emit(col); },
      delete: async () => { delete store[col][id]; emit(col); },
    }; },
  };
  window.__calls = [];
  const sampleFn = async (prompt, o) => { const t = prompt.includes('Ordre du patron') ? 'Page d accueil rédigée : Bienvenue à la clinique…' : 'Proposition rédigée pour Clinique Test'; await new Promise((r) => setTimeout(r, 300)); o.onText && o.onText({ text: t.slice(0, 10), delta: t.slice(0, 10) }); await new Promise((r) => setTimeout(r, 300)); o.onText && o.onText({ text: t, delta: t.slice(10) }); return { text: t }; };
  const sample = Object.assign(sampleFn, { json: async (prompt, o) => {
    window.__calls.push(o.modelTier);
    await new Promise((r) => setTimeout(r, 30));
    o.onText && o.onText({ text: '{...}', delta: '{...}' });
    if (prompt.includes('Boussole')) return 'Proposition pour Clinique Test';
    if (prompt.includes('brouillon indépendant')) return { diagnostic: 'Diag <b>x</b>', proposition: 'Proposer un portail de réservation', etapes: ['a', 'b'], risques: ['r'], livrable: { titre: 'Script vidéo', contenu: 'Texte du livrable' }, confiance: 70 };
    if (prompt.includes('agent relais')) return { idee_forte: 'Portail de réservation bilingue', soutien: '2/2', desaccord: '' };
    if (prompt.includes('Critique les brouillons')) return { position: 'Position révisée', changements: 'ok', grille: { impact: 4, cout: 3, delai: 4, risque: 3, conformite: 5 } };
    return { titre: 'Plan final', resume: 'Résumé', diagnostic: 'D', solution: 'S', etapes: [{ semaine: 'S1', action: 'Audit', equipe: 'Infra' }], budget: { mise_en_place: '5 000 $', mensuel: '690 $', hypotheses: 'h' }, risques: [{ risque: 'R', parade: 'P' }], livrables: [{ equipe: 'Contenu', titre: 'Script', contenu: '...' }], points_a_valider: ['Prix'] };
  } });
  window.__drive = []; const mcp = { callTool: async (s, t, i) => { window.__drive.push([s, t, i.title, i.textContent.length]); return { payload: { id: 'f1', webViewLink: 'https://docs.google.com/document/d/f1' } }; } }; window.claude = { use: async (name) => (name === 'db' ? db : name === 'sample' ? sample : name === 'mcp' ? mcp : null) };
  window.__seed = async () => {
    await db.collection('clients').add({ name: 'Clinique Test', sector: 'Santé', stage: 'actif', mrr: 1490, need: 'Appels manqués' });
    await db.collection('clients').add({ name: 'Cabinet Test', sector: 'Droit', stage: 'prospect', mrr: 0, need: 'Dossiers' });
    await db.collection('tasks').add({ title: 'Remplacer les clés', owner: 'toi', status: 'à faire', area: 'site', order: 1 });
    await db.collection('offers').add({ code: 'A', name: 'Diagnostic', forWho: 'PME', setup: '1 500 $', monthly: '—' });
    const T = [['infra','Alex','Architecte cloud'],['infra','Sam','Sécurité'],['contenu','Léa','Vidéo IA'],['web','Max','Automatisation'],['ventes','Zoé','Ventes']];
    for (const [team, name, role] of T) await db.collection('agents').add({ team, name, role, method: 'Étape 1\nÉtape 2', engine: name === 'Sam' ? 'maison' : 'claude', active: true });
  };
})();
