import { randomUUID } from "node:crypto";
import type { DecisionState, Message, MessageRole, Participant, Room, Swarm } from "./types.js";

// Pure swarm mechanics. No I/O, no timers: the Hive drives these.

export const DECISION = {
  radius: 1,
  captureRadius: 0.14,
  damping: 0.9,
  gain: 0.0015,
  maxTicks: 600, // 60s at 100ms
  reconsiderEvery: 8,
  // An agent switches to where the swarm is heading if it likes that option
  // at least this close to its own favourite (the "best acceptable" move).
  switchMargin: 0.18,
};

const STOP = new Set(
  "a an and are as at be but by for from has have i i'm in is it it's its of on or so that the this to was we with you they our not if can will would should could just think also more very what about which there their them than then into out up do does".split(
    " ",
  ),
);

export function now(): string {
  return new Date().toISOString();
}

export function newDecision(optionCount: number): DecisionState {
  return {
    status: "idle",
    x: 0,
    y: 0,
    vx: 0,
    vy: 0,
    ticks: 0,
    maxTicks: DECISION.maxTicks,
    winner: null,
    support: new Array(optionCount).fill(0),
    trace: [],
  };
}

export function createSwarm(input: {
  name: string;
  question: string;
  options: string[];
  roomSize?: number;
  surrogateConnectorId?: string;
}): Swarm {
  const options = input.options.map((o) => o.trim()).filter(Boolean);
  if (!input.name.trim()) throw new Error("name is required");
  if (!input.question.trim()) throw new Error("question is required");
  if (options.length < 2 || options.length > 8) throw new Error("give 2 to 8 options");
  const roomSize = Math.max(2, Math.min(12, Math.round(input.roomSize ?? 5)));
  return {
    id: randomUUID(),
    name: input.name.trim(),
    question: input.question.trim(),
    options,
    roomSize,
    status: "draft",
    surrogateConnectorId: input.surrogateConnectorId,
    participants: [],
    rooms: [],
    messages: [],
    decision: newDecision(options.length),
    createdAt: now(),
  };
}

// Thinkscape splits a large group into small rooms where real conversation
// works. New members fill the emptiest room; a new room opens when all are full.
export function assignRoom(swarm: Swarm): Room {
  const counts = new Map(swarm.rooms.map((r) => [r.id, 0]));
  for (const p of swarm.participants) counts.set(p.roomId, (counts.get(p.roomId) ?? 0) + 1);
  let best: Room | undefined;
  for (const room of swarm.rooms) {
    const n = counts.get(room.id) ?? 0;
    if (n < swarm.roomSize && (!best || n < (counts.get(best.id) ?? 0))) best = room;
  }
  if (best) return best;
  const room: Room = { id: randomUUID(), name: `Room ${swarm.rooms.length + 1}`, relayedUpTo: 0 };
  swarm.rooms.push(room);
  return room;
}

export function addParticipant(
  swarm: Swarm,
  input: { name: string; kind: Participant["kind"]; connectorId?: string },
): Participant {
  if (swarm.status === "closed") throw new Error("swarm is closed");
  if (!input.name.trim()) throw new Error("name is required");
  const room = assignRoom(swarm);
  const p: Participant = {
    id: randomUUID(),
    name: input.name.trim(),
    kind: input.kind,
    connectorId: input.connectorId,
    roomId: room.id,
    target: null,
    conviction: 0,
    joinedAt: now(),
  };
  swarm.participants.push(p);
  return p;
}

export function removeParticipant(swarm: Swarm, participantId: string): boolean {
  const before = swarm.participants.length;
  swarm.participants = swarm.participants.filter((p) => p.id !== participantId);
  return swarm.participants.length !== before;
}

export function words(text: string): string[] {
  return text
    .toLowerCase()
    .replace(/[^\p{L}\p{N}\s'-]/gu, " ")
    .split(/\s+/)
    .filter((w) => w.length > 2 && !STOP.has(w));
}

// Which options a message argues for: an option counts when all of its
// significant words appear in the text.
export function detectMentions(text: string, options: string[]): number[] {
  const bag = new Set(words(text));
  const out: number[] = [];
  options.forEach((opt, i) => {
    const ws = words(opt);
    if (ws.length > 0 && ws.every((w) => bag.has(w))) out.push(i);
  });
  return out;
}

export function postMessage(
  swarm: Swarm,
  input: { roomId: string; authorId: string; authorName: string; role: MessageRole; text: string },
): Message {
  const text = input.text.trim().slice(0, 2000);
  if (!text) throw new Error("text is required");
  if (!swarm.rooms.some((r) => r.id === input.roomId)) throw new Error("unknown room");
  const msg: Message = {
    id: randomUUID(),
    roomId: input.roomId,
    authorId: input.authorId,
    authorName: input.authorName,
    role: input.role,
    text,
    mentions: detectMentions(text, swarm.options),
    at: now(),
  };
  swarm.messages.push(msg);
  return msg;
}

export function roomMessages(swarm: Swarm, roomId: string): Message[] {
  return swarm.messages.filter((m) => m.roomId === roomId);
}

// Thinkscape links rooms in a ring: each room's surrogate carries its
// insight into the next room, so ideas propagate across the whole group.
export function nextRoom(swarm: Swarm, roomId: string): Room | undefined {
  if (swarm.rooms.length < 2) return undefined;
  const i = swarm.rooms.findIndex((r) => r.id === roomId);
  return swarm.rooms[(i + 1) % swarm.rooms.length];
}

// Member messages this room produced since its surrogate last relayed.
export function unrelayed(swarm: Swarm, room: Room): Message[] {
  const own = roomMessages(swarm, room.id).filter((m) => m.role === "member");
  return own.slice(room.relayedUpTo);
}

export function markRelayed(swarm: Swarm, room: Room): void {
  room.relayedUpTo = roomMessages(swarm, room.id).filter((m) => m.role === "member").length;
}

// Offline surrogate: the room's leading option plus the message that best
// carries the room's most repeated ideas.
export function extractiveInsight(swarm: Swarm, room: Room, msgs: Message[]): string {
  const tally = new Array(swarm.options.length).fill(0);
  for (const m of msgs) for (const i of m.mentions) tally[i]++;
  const freq = new Map<string, number>();
  for (const m of msgs) for (const w of new Set(words(m.text))) freq.set(w, (freq.get(w) ?? 0) + 1);
  let quote = msgs[0];
  let best = -1;
  for (const m of msgs) {
    const ws = new Set(words(m.text));
    let score = 0;
    for (const w of ws) score += (freq.get(w) ?? 0) - 1;
    score = score / Math.sqrt(ws.size || 1);
    if (score > best) {
      best = score;
      quote = m;
    }
  }
  const lead = argmax(tally);
  const people = new Set(msgs.map((m) => m.authorId)).size;
  const head =
    lead >= 0 && tally[lead] > 0
      ? `${room.name} leans "${swarm.options[lead]}" (${tally[lead]} of ${msgs.length} messages, ${people} voices).`
      : `${room.name} has no clear lean yet (${people} voices).`;
  return `${head} Strongest point: "${clip(quote.text, 220)}" (${quote.authorName})`;
}

export function clip(text: string, n: number): string {
  return text.length <= n ? text : text.slice(0, n - 1).trimEnd() + "…";
}

export function argmax(xs: number[]): number {
  let best = -1;
  for (let i = 0; i < xs.length; i++) if (best < 0 || xs[i] > xs[best]) best = i;
  return best;
}

// Conversation-level collective answer: option mentions across every room,
// with each room counting once per distinct voice (no single loud member wins).
export function conversationTally(swarm: Swarm): number[] {
  const tally = new Array(swarm.options.length).fill(0);
  const seen = new Set<string>();
  for (const m of swarm.messages) {
    if (m.role !== "member") continue;
    for (const i of m.mentions) {
      const key = `${m.authorId}:${i}`;
      if (seen.has(key)) continue;
      seen.add(key);
      tally[i]++;
    }
  }
  return tally;
}

// --- Swarm decision: the puck --------------------------------------------

export function optionPoint(i: number, n: number): [number, number] {
  const a = (i / n) * Math.PI * 2 - Math.PI / 2;
  return [Math.cos(a) * DECISION.radius, Math.sin(a) * DECISION.radius];
}

export function startDecision(swarm: Swarm): void {
  if (swarm.participants.length === 0) throw new Error("no participants");
  swarm.decision = newDecision(swarm.options.length);
  swarm.decision.status = "running";
  swarm.decision.startedAt = now();
  for (const p of swarm.participants) {
    if (p.kind === "agent" && p.prefs) {
      p.target = argmax(p.prefs);
      const sorted = [...p.prefs].sort((a, b) => b - a);
      p.conviction = clamp(0.35 + (sorted[0] - (sorted[1] ?? 0)) * 1.5, 0.2, 1);
    } else {
      p.target = null;
      p.conviction = 0;
    }
  }
}

export function setPull(swarm: Swarm, participantId: string, target: number | null, conviction = 1): void {
  const p = swarm.participants.find((x) => x.id === participantId);
  if (!p) throw new Error("unknown participant");
  if (target !== null && (target < 0 || target >= swarm.options.length || !Number.isInteger(target)))
    throw new Error("bad option");
  p.target = target;
  p.conviction = target === null ? 0 : clamp(conviction, 0, 1);
}

// The option the puck is currently moving toward, if it is moving at all.
export function heading(d: DecisionState, n: number): number {
  const speed = Math.hypot(d.vx, d.vy);
  if (speed < 1e-4) return -1;
  let best = -1;
  let bestCos = 0.5;
  for (let i = 0; i < n; i++) {
    const [ox, oy] = optionPoint(i, n);
    const dx = ox - d.x;
    const dy = oy - d.y;
    const cos = (dx * d.vx + dy * d.vy) / (Math.hypot(dx, dy) * speed || 1);
    if (cos > bestCos) {
      bestCos = cos;
      best = i;
    }
  }
  return best;
}

// One 100ms step. Every member's pull is a vector toward their option; the
// puck follows the sum. Agents watch where it heads and, like people in a
// real swarm, swing to that option when they find it acceptable.
export function decisionTick(swarm: Swarm): void {
  const d = swarm.decision;
  if (d.status !== "running") return;
  const n = swarm.options.length;
  d.ticks++;

  if (d.ticks % DECISION.reconsiderEvery === 0) {
    const h = heading(d, n);
    if (h >= 0) {
      for (const p of swarm.participants) {
        if (p.kind !== "agent" || !p.prefs || p.target === null || p.target === h) continue;
        if (p.prefs[h] >= p.prefs[p.target] - DECISION.switchMargin) p.target = h;
      }
    }
  }

  let fx = 0;
  let fy = 0;
  const support = new Array(n).fill(0);
  let total = 0;
  for (const p of swarm.participants) {
    if (p.target === null || p.conviction <= 0) continue;
    const [ox, oy] = optionPoint(p.target, n);
    const dx = ox - d.x;
    const dy = oy - d.y;
    const len = Math.hypot(dx, dy) || 1;
    fx += (dx / len) * p.conviction;
    fy += (dy / len) * p.conviction;
    support[p.target] += p.conviction;
    total += p.conviction;
  }
  const crowd = Math.max(1, swarm.participants.length);
  d.vx = d.vx * DECISION.damping + (fx / crowd) * DECISION.gain;
  d.vy = d.vy * DECISION.damping + (fy / crowd) * DECISION.gain;
  d.x += d.vx;
  d.y += d.vy;
  const r = Math.hypot(d.x, d.y);
  if (r > DECISION.radius) {
    d.x = (d.x / r) * DECISION.radius;
    d.y = (d.y / r) * DECISION.radius;
  }
  d.support = support.map((s) => (total > 0 ? s / total : 0));
  if (d.ticks % 2 === 0) d.trace.push([round(d.x), round(d.y)]);

  for (let i = 0; i < n; i++) {
    const [ox, oy] = optionPoint(i, n);
    if (Math.hypot(ox - d.x, oy - d.y) <= DECISION.captureRadius) return finish(d, i);
  }
  if (d.ticks >= d.maxTicks) finish(d, nearest(d, n));
}

function finish(d: DecisionState, winner: number): void {
  d.status = "done";
  d.winner = winner;
  d.vx = 0;
  d.vy = 0;
  d.endedAt = now();
}

function nearest(d: DecisionState, n: number): number {
  let best = 0;
  let bestDist = Infinity;
  for (let i = 0; i < n; i++) {
    const [ox, oy] = optionPoint(i, n);
    const dist = Math.hypot(ox - d.x, oy - d.y);
    if (dist < bestDist) {
      bestDist = dist;
      best = i;
    }
  }
  return best;
}

function clamp(v: number, lo: number, hi: number): number {
  return Math.max(lo, Math.min(hi, Number.isFinite(v) ? v : lo));
}

function round(v: number): number {
  return Math.round(v * 1000) / 1000;
}
