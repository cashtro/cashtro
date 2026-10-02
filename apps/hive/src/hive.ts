import { EventEmitter } from "node:events";
import { existsSync, mkdirSync, readFileSync, renameSync, writeFileSync } from "node:fs";
import path from "node:path";
import { createConnector, liveMind, turnContext, type Mind } from "./connectors.js";
import * as E from "./engine.js";
import type { Connector, HiveState, Message, Participant, Swarm } from "./types.js";

export interface HiveOpts {
  dataFile?: string;
  mind?: Mind;
  // How many turns an agent takes in a row before it waits for a human.
  agentTurns?: number;
  // Relay surrogate insights every N conversation ticks.
  relayEvery?: number;
  // Agent messages a swarm allows between two human messages.
  agentBudget?: number;
}

export type HiveEvent =
  | { type: "swarm"; swarm: Swarm }
  | { type: "message"; message: Message }
  | { type: "decision"; decision: Swarm["decision"]; pulls: Array<Pick<Participant, "id" | "target" | "conviction">> }
  | { type: "deleted"; id: string };

// The Hive hosts every swarm and connector, drives the conversation and
// decision loops, and fans out live events.
export class Hive {
  readonly events = new EventEmitter();
  private state: HiveState = { swarms: [], connectors: [] };
  private mind: Mind;
  private agentTurns: number;
  private relayEvery: number;
  private agentBudget: number;
  private quiet = new Map<string, number>();
  private turns = new Map<string, number>();
  private busy = new Set<string>();
  private convoCount = 0;
  private saveTimer?: NodeJS.Timeout;
  private loops: NodeJS.Timeout[] = [];

  constructor(private opts: HiveOpts = {}) {
    this.mind = opts.mind ?? liveMind;
    this.agentTurns = opts.agentTurns ?? 6;
    this.relayEvery = opts.relayEvery ?? 2;
    this.agentBudget = opts.agentBudget ?? Number(process.env.HIVE_AGENT_BUDGET || 40);
    this.events.setMaxListeners(0);
    if (opts.dataFile && existsSync(opts.dataFile)) {
      this.state = JSON.parse(readFileSync(opts.dataFile, "utf8")) as HiveState;
    }
    if (this.state.connectors.length === 0) this.seedConnectors();
  }

  private seedConnectors(): void {
    const starters = [
      { name: "Ada (analyst)", persona: "a data-driven analyst who trusts evidence, cost and risk numbers" },
      { name: "Marco (builder)", persona: "a pragmatic builder who wants speed, shipping and simple solutions" },
      { name: "Lina (customer voice)", persona: "a customer advocate who cares about people, trust and quality" },
    ];
    for (const s of starters) this.state.connectors.push(createConnector({ ...s, kind: "persona" }));
  }

  // --- persistence -------------------------------------------------------

  private touch(): void {
    if (!this.opts.dataFile) return;
    clearTimeout(this.saveTimer);
    this.saveTimer = setTimeout(() => this.flush(), 400);
  }

  flush(): void {
    const file = this.opts.dataFile;
    if (!file) return;
    mkdirSync(path.dirname(file), { recursive: true });
    const tmp = `${file}.tmp`;
    writeFileSync(tmp, JSON.stringify(this.state));
    renameSync(tmp, file);
  }

  private emit(swarmId: string, ev: HiveEvent): void {
    this.events.emit(swarmId, ev);
    this.events.emit("*", { swarmId, ...ev });
  }

  private changed(swarm: Swarm): void {
    this.emit(swarm.id, { type: "swarm", swarm });
    this.touch();
  }

  // --- swarms ------------------------------------------------------------

  listSwarms(): Swarm[] {
    return this.state.swarms;
  }

  getSwarm(id: string): Swarm {
    const s = this.state.swarms.find((x) => x.id === id);
    if (!s) throw new NotFound("swarm not found");
    return s;
  }

  createSwarm(input: Parameters<typeof E.createSwarm>[0]): Swarm {
    if (input.surrogateConnectorId) this.getConnector(input.surrogateConnectorId);
    const s = E.createSwarm(input);
    this.state.swarms.unshift(s);
    this.touch();
    return s;
  }

  updateSwarm(id: string, input: { surrogateConnectorId?: string | null; roomSize?: number }): Swarm {
    const s = this.getSwarm(id);
    if (input.surrogateConnectorId !== undefined) {
      if (input.surrogateConnectorId) this.getConnector(input.surrogateConnectorId);
      s.surrogateConnectorId = input.surrogateConnectorId || undefined;
    }
    if (input.roomSize !== undefined) s.roomSize = Math.max(2, Math.min(12, Math.round(input.roomSize)));
    this.changed(s);
    return s;
  }

  deleteSwarm(id: string): void {
    this.getSwarm(id);
    this.state.swarms = this.state.swarms.filter((s) => s.id !== id);
    this.emit(id, { type: "deleted", id });
    this.touch();
  }

  setStatus(id: string, status: Swarm["status"]): Swarm {
    const s = this.getSwarm(id);
    s.status = status;
    if (status === "closed" && s.decision.status === "running") s.decision.status = "done";
    if (status === "live" && s.rooms.length > 0 && !s.messages.some((m) => m.role === "system")) {
      for (const r of s.rooms) {
        E.postMessage(s, {
          roomId: r.id,
          authorId: "hive",
          authorName: "Hive",
          role: "system",
          text: `Welcome to ${r.name}. Question: ${s.question} Options: ${s.options.join(" · ")}`,
        });
      }
    }
    this.changed(s);
    return s;
  }

  join(id: string, input: { name: string; kind: Participant["kind"]; connectorId?: string }): Participant {
    const s = this.getSwarm(id);
    let name = input.name;
    if (input.kind === "agent") {
      const c = this.getConnector(input.connectorId ?? "");
      name = name?.trim() || c.name;
    }
    const roomsBefore = s.rooms.length;
    const p = E.addParticipant(s, { name, kind: input.kind, connectorId: input.connectorId });
    if (s.status === "live" && s.rooms.length > roomsBefore) {
      const r = s.rooms[s.rooms.length - 1];
      E.postMessage(s, {
        roomId: r.id,
        authorId: "hive",
        authorName: "Hive",
        role: "system",
        text: `Welcome to ${r.name}. Question: ${s.question} Options: ${s.options.join(" · ")}`,
      });
    }
    this.changed(s);
    return p;
  }

  leave(id: string, participantId: string): void {
    const s = this.getSwarm(id);
    if (!E.removeParticipant(s, participantId)) throw new NotFound("participant not found");
    this.changed(s);
  }

  say(id: string, participantId: string, text: string): Message {
    const s = this.getSwarm(id);
    if (s.status !== "live") throw new Error("swarm is not live");
    const p = s.participants.find((x) => x.id === participantId);
    if (!p) throw new NotFound("participant not found");
    const msg = E.postMessage(s, { roomId: p.roomId, authorId: p.id, authorName: p.name, role: "member", text });
    if (p.kind === "human") {
      this.quiet.delete(s.id);
      for (const a of s.participants) if (a.roomId === p.roomId) this.turns.delete(a.id);
    }
    this.emit(s.id, { type: "message", message: msg });
    this.touch();
    return msg;
  }

  summary(id: string) {
    const s = this.getSwarm(id);
    const tally = E.conversationTally(s);
    const lead = E.argmax(tally);
    return {
      question: s.question,
      options: s.options,
      conversation: { tally, leading: lead >= 0 && tally[lead] > 0 ? s.options[lead] : null },
      decision: {
        status: s.decision.status,
        winner: s.decision.winner === null ? null : s.options[s.decision.winner],
        support: s.decision.support,
        seconds: s.decision.ticks / 10,
      },
      rooms: s.rooms.map((r) => ({
        name: r.name,
        members: s.participants.filter((p) => p.roomId === r.id).map((p) => p.name),
        messages: E.roomMessages(s, r.id).length,
      })),
      insights: s.messages.filter((m) => m.role === "surrogate").map((m) => ({ to: roomName(s, m.roomId), text: m.text })),
    };
  }

  // --- connectors --------------------------------------------------------

  listConnectors(): Connector[] {
    return this.state.connectors;
  }

  getConnector(id: string): Connector {
    const c = this.state.connectors.find((x) => x.id === id);
    if (!c) throw new NotFound("connector not found");
    return c;
  }

  addConnector(input: Parameters<typeof createConnector>[0]): Connector {
    const c = createConnector(input);
    this.state.connectors.push(c);
    this.touch();
    return c;
  }

  deleteConnector(id: string): void {
    this.getConnector(id);
    const used = this.state.swarms.some(
      (s) => s.surrogateConnectorId === id || s.participants.some((p) => p.connectorId === id),
    );
    if (used) throw new Error("connector is in use by a swarm");
    this.state.connectors = this.state.connectors.filter((c) => c.id !== id);
    this.touch();
  }

  async testConnector(id: string): Promise<{ ok: boolean; reply?: string; error?: string }> {
    const c = this.getConnector(id);
    try {
      const reply = await this.mind.speak(c, {
        question: "Which is better for a first test?",
        options: ["Short answer", "Long answer"],
        roomName: "Test room",
        selfName: c.name,
        recent: [],
      });
      return { ok: Boolean(reply), reply };
    } catch (err) {
      return { ok: false, error: (err as Error).message };
    }
  }

  // --- conversation loop -------------------------------------------------

  // One conversation step for every live swarm: at most one agent speaks per
  // room, and every few steps each room's surrogate relays to the next room.
  async convoTick(): Promise<void> {
    this.convoCount++;
    const jobs: Promise<void>[] = [];
    for (const s of this.state.swarms) {
      if (s.status !== "live") continue;
      for (const room of s.rooms) {
        jobs.push(this.agentTurn(s, room.id));
        if (this.convoCount % this.relayEvery === 0) jobs.push(this.relay(s, room.id));
      }
    }
    await Promise.all(jobs);
  }

  private async agentTurn(s: Swarm, roomId: string): Promise<void> {
    const key = `turn:${s.id}:${roomId}`;
    if (this.busy.has(key) || (this.quiet.get(s.id) ?? 0) >= this.agentBudget) return;
    const room = s.rooms.find((r) => r.id === roomId);
    if (!room) return;
    const msgs = E.roomMessages(s, roomId);
    const last = msgs[msgs.length - 1];
    const agents = s.participants.filter(
      (p) => p.kind === "agent" && p.roomId === roomId && (this.turns.get(p.id) ?? 0) < this.agentTurns,
    );
    // Rotate: the agent who spoke least recently goes next, never twice in a row.
    const lastSpoke = (p: Participant) => {
      for (let i = msgs.length - 1; i >= 0; i--) if (msgs[i].authorId === p.id) return i;
      return -1;
    };
    const next = agents.filter((p) => p.id !== last?.authorId).sort((a, b) => lastSpoke(a) - lastSpoke(b))[0];
    if (!next || !next.connectorId) return;
    const c = this.state.connectors.find((x) => x.id === next.connectorId);
    if (!c) return;
    this.busy.add(key);
    try {
      const text = await this.mind.speak(c, turnContext(s.question, s.options, room.name, next.name, msgs));
      if (!text || s.status !== "live" || !s.participants.includes(next)) return;
      this.turns.set(next.id, (this.turns.get(next.id) ?? 0) + 1);
      this.quiet.set(s.id, (this.quiet.get(s.id) ?? 0) + 1);
      const msg = E.postMessage(s, { roomId, authorId: next.id, authorName: next.name, role: "member", text });
      this.emit(s.id, { type: "message", message: msg });
      this.touch();
    } catch (err) {
      this.turns.set(next.id, this.agentTurns);
      this.system(s, roomId, `${next.name} could not answer (${(err as Error).message}).`);
    } finally {
      this.busy.delete(key);
    }
  }

  private async relay(s: Swarm, roomId: string): Promise<void> {
    const key = `relay:${s.id}:${roomId}`;
    if (this.busy.has(key)) return;
    const room = s.rooms.find((r) => r.id === roomId);
    const to = room && E.nextRoom(s, roomId);
    if (!room || !to) return;
    const fresh = E.unrelayed(s, room);
    if (fresh.length < 2) return;
    this.busy.add(key);
    try {
      let text = "";
      const c = s.surrogateConnectorId && this.state.connectors.find((x) => x.id === s.surrogateConnectorId);
      if (c && c.kind !== "persona") {
        try {
          text = await this.mind.summarize(c, turnContext(s.question, s.options, room.name, "Surrogate", fresh, 20));
        } catch {
          text = "";
        }
      }
      if (!text) text = E.extractiveInsight(s, room, fresh);
      E.markRelayed(s, room);
      const msg = E.postMessage(s, {
        roomId: to.id,
        authorId: `surrogate:${room.id}`,
        authorName: `Surrogate · ${room.name}`,
        role: "surrogate",
        text,
      });
      // A new idea from another room earns each agent there one reply.
      for (const p of s.participants) {
        if (p.roomId === to.id) this.turns.set(p.id, Math.min(this.turns.get(p.id) ?? 0, this.agentTurns - 1));
      }
      this.emit(s.id, { type: "message", message: msg });
      this.touch();
    } finally {
      this.busy.delete(key);
    }
  }

  private system(s: Swarm, roomId: string, text: string): void {
    const msg = E.postMessage(s, { roomId, authorId: "hive", authorName: "Hive", role: "system", text });
    this.emit(s.id, { type: "message", message: msg });
  }

  // --- decision loop -----------------------------------------------------

  async startDecision(id: string): Promise<Swarm> {
    const s = this.getSwarm(id);
    if (s.status === "closed") throw new Error("swarm is closed");
    await Promise.all(
      s.participants.map(async (p) => {
        if (p.kind !== "agent" || !p.connectorId) return;
        const c = this.state.connectors.find((x) => x.id === p.connectorId);
        if (!c) return;
        try {
          p.prefs = await this.mind.rank(c, s.question, s.options);
        } catch {
          p.prefs = undefined;
        }
      }),
    );
    E.startDecision(s);
    this.changed(s);
    return s;
  }

  pull(id: string, participantId: string, target: number | null, conviction?: number): void {
    const s = this.getSwarm(id);
    if (s.decision.status !== "running") throw new Error("no decision is running");
    E.setPull(s, participantId, target, conviction);
  }

  decisionTick(): void {
    for (const s of this.state.swarms) {
      if (s.decision.status !== "running") continue;
      E.decisionTick(s);
      const done = (s.decision.status as Swarm["decision"]["status"]) === "done";
      const pulls = s.participants.map((p) => ({ id: p.id, target: p.target, conviction: p.conviction }));
      this.emit(s.id, { type: "decision", decision: s.decision, pulls });
      if (done) {
        this.system(s, s.rooms[0].id, `The swarm converged on "${s.options[s.decision.winner!]}" in ${(s.decision.ticks / 10).toFixed(1)}s.`);
        this.changed(s);
      }
    }
  }

  start(convoMs = Number(process.env.HIVE_CONVO_MS || 4000)): void {
    this.loops.push(setInterval(() => this.decisionTick(), 100));
    this.loops.push(setInterval(() => void this.convoTick(), convoMs));
  }

  stop(): void {
    for (const t of this.loops) clearInterval(t);
    this.loops = [];
    clearTimeout(this.saveTimer);
    this.flush();
  }
}

export class NotFound extends Error {}

function roomName(s: Swarm, roomId: string): string {
  return s.rooms.find((r) => r.id === roomId)?.name ?? "?";
}
