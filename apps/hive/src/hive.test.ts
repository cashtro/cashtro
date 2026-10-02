import assert from "node:assert/strict";
import http from "node:http";
import type { AddressInfo } from "node:net";
import { mkdtemp, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import path from "node:path";
import { test } from "node:test";
import { Hive } from "./hive.js";
import { createHiveServer } from "./server.js";
import { parseScores, personaRank } from "./connectors.js";

async function listen(server: http.Server): Promise<string> {
  await new Promise<void>((r) => server.listen(0, "127.0.0.1", r));
  return `http://127.0.0.1:${(server.address() as AddressInfo).port}`;
}

async function call(base: string, method: string, url: string, body?: unknown, token?: string) {
  const res = await fetch(base + url, {
    method,
    headers: { "content-type": "application/json", ...(token ? { authorization: `Bearer ${token}` } : {}) },
    body: body ? JSON.stringify(body) : undefined,
  });
  return { status: res.status, body: (await res.json()) as any };
}

test("host a swarm end to end: create, connect, converse, relay, decide", async () => {
  // A webhook agent that always argues for the client portal.
  const hook = http.createServer((req, res) => {
    let raw = "";
    req.on("data", (c) => (raw += c));
    req.on("end", () => {
      const body = JSON.parse(raw);
      res.setHeader("content-type", "application/json");
      if (body.type === "rank") return res.end(JSON.stringify({ scores: [0.1, 0.9, 0.2] }));
      res.end(JSON.stringify({ text: "Client portal first: it is what clients ask for." }));
    });
  });
  const hookUrl = await listen(hook);
  const hive = new Hive({ agentTurns: 2, relayEvery: 1 });
  const server = createHiveServer(hive, { token: "" });
  const base = await listen(server);
  try {
    const health = await call(base, "GET", "/health");
    assert.equal(health.body.ok, true);
    assert.equal(health.body.connectors, 3, "starter personas are seeded");

    const created = await call(base, "POST", "/api/swarms", {
      name: "Roadmap",
      question: "What should we build next?",
      options: "Mobile app\nClient portal\nAI assistant",
      roomSize: 2,
    });
    assert.equal(created.status, 201);
    const id = created.body.id;

    const conn = await call(base, "POST", "/api/connectors", { name: "Portal bot", kind: "webhook", url: hookUrl + "/hook" });
    assert.equal(conn.status, 201);
    const bad = await call(base, "POST", "/api/connectors", { name: "x", kind: "webhook", url: "ftp://nope" });
    assert.equal(bad.status, 400);

    const ana = (await call(base, "POST", `/api/swarms/${id}/participants`, { name: "Ana", kind: "human" })).body;
    await call(base, "POST", `/api/swarms/${id}/participants`, { kind: "agent", connectorId: conn.body.id });
    const persona = hive.listConnectors()[0];
    await call(base, "POST", `/api/swarms/${id}/participants`, { kind: "agent", connectorId: persona.id });
    await call(base, "POST", `/api/swarms/${id}/participants`, { kind: "agent", connectorId: persona.id, name: "Ada 2" });

    let swarm = (await call(base, "GET", `/api/swarms/${id}`)).body;
    assert.equal(swarm.rooms.length, 2);

    const early = await call(base, "POST", `/api/swarms/${id}/messages`, { participantId: ana.id, text: "hi" });
    assert.equal(early.status, 400, "cannot talk before the swarm is live");

    await call(base, "POST", `/api/swarms/${id}/start`);
    const said = await call(base, "POST", `/api/swarms/${id}/messages`, { participantId: ana.id, text: "I like the mobile app idea" });
    assert.equal(said.status, 201);
    assert.deepEqual(said.body.mentions, [0]);

    for (let i = 0; i < 6; i++) await hive.convoTick();
    const live = hive.getSwarm(id);
    const bot = live.messages.filter((m) => m.authorName === "Portal bot");
    // Two turns of its own, plus one reply to the insight relayed from Room 2.
    assert.equal(bot.length, 3, "agent takes turns up to its limit, then waits for a human");
    assert.ok(live.messages.some((m) => m.role === "surrogate"), "surrogates relay insights between rooms");

    // Agents alone cannot chatter forever: the swarm budget stops them.
    const capped = new Hive({ agentTurns: 100, relayEvery: 1, agentBudget: 5 });
    const cs = capped.createSwarm({ name: "c", question: "q", options: ["a", "b"] });
    const pc = capped.listConnectors()[0].id;
    capped.join(cs.id, { name: "x", kind: "agent", connectorId: pc });
    capped.join(cs.id, { name: "y", kind: "agent", connectorId: pc });
    capped.setStatus(cs.id, "live");
    for (let i = 0; i < 20; i++) await capped.convoTick();
    assert.equal(capped.getSwarm(cs.id).messages.filter((m) => m.role === "member").length, 5);

    const summary = (await call(base, "GET", `/api/swarms/${id}/summary`)).body;
    assert.equal(summary.rooms.length, 2);
    assert.ok(summary.conversation.tally[1] >= 1);

    // Ana pulls hard toward the AI assistant; the agents bring their own pulls.
    await call(base, "POST", `/api/swarms/${id}/decision/start`);
    const pull = await call(base, "POST", `/api/swarms/${id}/decision/pull`, { participantId: ana.id, target: 1, conviction: 1 });
    assert.equal(pull.status, 200);
    for (let i = 0; i < 700 && hive.getSwarm(id).decision.status === "running"; i++) hive.decisionTick();
    const d = hive.getSwarm(id).decision;
    assert.equal(d.status, "done");
    assert.equal(d.winner, 1);

    const inUse = await call(base, "DELETE", `/api/connectors/${conn.body.id}`);
    assert.equal(inUse.status, 400);
    assert.equal((await call(base, "DELETE", `/api/swarms/${id}`)).status, 200);
    assert.equal((await call(base, "GET", `/api/swarms/${id}`)).status, 404);
  } finally {
    server.close();
    hook.close();
  }
});

test("a HIVE_TOKEN guards every write", async () => {
  const server = createHiveServer(new Hive(), { token: "s3cret" });
  const base = await listen(server);
  try {
    const body = { name: "n", question: "q", options: ["a", "b"] };
    assert.equal((await call(base, "POST", "/api/swarms", body)).status, 401);
    assert.equal((await call(base, "POST", "/api/swarms", body, "s3cret")).status, 201);
    assert.equal((await call(base, "GET", "/api/swarms")).status, 200);
  } finally {
    server.close();
  }
});

test("state survives a restart", async () => {
  const dir = await mkdtemp(path.join(tmpdir(), "hive-"));
  const file = path.join(dir, "hive.json");
  try {
    const a = new Hive({ dataFile: file });
    const s = a.createSwarm({ name: "Keep", question: "q", options: ["a", "b"] });
    a.join(s.id, { name: "Ana", kind: "human" });
    a.stop();
    const b = new Hive({ dataFile: file });
    assert.equal(b.getSwarm(s.id).participants[0].name, "Ana");
    b.stop();
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
});

test("model scores are parsed defensively", () => {
  assert.deepEqual(parseScores("Sure: [0.2, 0.9]", 2), [0.2, 0.9]);
  assert.deepEqual(parseScores("[2, 8]", 2), [0.25, 1]);
  assert.equal(parseScores("no idea", 2), null);
  assert.equal(parseScores("[1]", 2), null);
  const c = { id: "c", name: "n", kind: "persona" as const, persona: "loves speed", createdAt: "" };
  assert.deepEqual(personaRank(c, "q", ["a", "b"]), personaRank(c, "q", ["a", "b"]), "personas are deterministic");
});
