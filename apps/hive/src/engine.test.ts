import assert from "node:assert/strict";
import { test } from "node:test";
import * as E from "./engine.js";

const base = () => E.createSwarm({ name: "Roadmap", question: "What should we build next?", options: ["Mobile app", "Client portal", "AI assistant"], roomSize: 3 });

test("createSwarm validates its input", () => {
  assert.throws(() => E.createSwarm({ name: "", question: "q", options: ["a", "b"] }), /name/);
  assert.throws(() => E.createSwarm({ name: "n", question: "q", options: ["only"] }), /2 to 8/);
  assert.equal(E.createSwarm({ name: "n", question: "q", options: ["a", "b"], roomSize: 99 }).roomSize, 12);
});

test("members fill small rooms and a new room opens when all are full", () => {
  const s = base();
  for (let i = 0; i < 7; i++) E.addParticipant(s, { name: `p${i}`, kind: "human" });
  assert.equal(s.rooms.length, 3);
  const sizes = s.rooms.map((r) => s.participants.filter((p) => p.roomId === r.id).length);
  assert.deepEqual(sizes, [3, 3, 1]);
});

test("rooms form a ring so insights travel through every room", () => {
  const s = base();
  for (let i = 0; i < 9; i++) E.addParticipant(s, { name: `p${i}`, kind: "human" });
  const [a, b, c] = s.rooms;
  assert.equal(E.nextRoom(s, a.id), b);
  assert.equal(E.nextRoom(s, c.id), a);
});

test("messages detect which options they argue for", () => {
  assert.deepEqual(E.detectMentions("I think the client portal beats a mobile app", ["Mobile app", "Client portal", "AI assistant"]), [0, 1]);
  assert.deepEqual(E.detectMentions("nothing relevant", ["Mobile app", "Client portal"]), []);
});

test("extractive surrogate reports the room's lean and relays only fresh messages", () => {
  const s = base();
  const p1 = E.addParticipant(s, { name: "Ana", kind: "human" });
  const p2 = E.addParticipant(s, { name: "Ben", kind: "human" });
  const room = s.rooms[0];
  E.postMessage(s, { roomId: room.id, authorId: p1.id, authorName: "Ana", role: "member", text: "The client portal saves support time" });
  E.postMessage(s, { roomId: room.id, authorId: p2.id, authorName: "Ben", role: "member", text: "Agreed, client portal first, support time matters" });
  const fresh = E.unrelayed(s, room);
  assert.equal(fresh.length, 2);
  const insight = E.extractiveInsight(s, room, fresh);
  assert.match(insight, /leans "Client portal"/);
  E.markRelayed(s, room);
  assert.equal(E.unrelayed(s, room).length, 0);
});

test("conversation tally counts each voice once per option", () => {
  const s = base();
  const p = E.addParticipant(s, { name: "Ana", kind: "human" });
  const r = s.rooms[0].id;
  E.postMessage(s, { roomId: r, authorId: p.id, authorName: "Ana", role: "member", text: "AI assistant" });
  E.postMessage(s, { roomId: r, authorId: p.id, authorName: "Ana", role: "member", text: "AI assistant again!" });
  assert.deepEqual(E.conversationTally(s), [0, 0, 1]);
});

test("the puck converges on the option the crowd pulls toward", () => {
  const s = base();
  const ps = [0, 1, 2, 3, 4].map((i) => E.addParticipant(s, { name: `p${i}`, kind: "human" }));
  E.startDecision(s);
  ps.forEach((p, i) => E.setPull(s, p.id, i < 3 ? 2 : 0));
  while (s.decision.status === "running") E.decisionTick(s);
  assert.equal(s.decision.winner, 2);
  assert.ok(s.decision.ticks < s.decision.maxTicks);
});

test("agents swing to an acceptable option the swarm is heading toward", () => {
  const s = base();
  const humans = [0, 1, 2].map((i) => E.addParticipant(s, { name: `h${i}`, kind: "human" }));
  const agent = E.addParticipant(s, { name: "bot", kind: "agent", connectorId: "x" });
  agent.prefs = [0.8, 0.7, 0.1]; // favourite: Mobile app, Client portal is acceptable
  E.startDecision(s);
  assert.equal(agent.target, 0);
  for (const h of humans) E.setPull(s, h.id, 1);
  while (s.decision.status === "running") E.decisionTick(s);
  assert.equal(s.decision.winner, 1);
  assert.equal(agent.target, 1);
});

test("setPull rejects options that do not exist", () => {
  const s = base();
  const p = E.addParticipant(s, { name: "Ana", kind: "human" });
  E.startDecision(s);
  assert.throws(() => E.setPull(s, p.id, 7), /bad option/);
});
