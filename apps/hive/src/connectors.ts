import { createHash, randomUUID } from "node:crypto";
import { argmax, clip, now, words } from "./engine.js";
import type { Connector, ConnectorKind, Message } from "./types.js";

// Everything a mind needs to take a turn in a room.
export interface TurnContext {
  question: string;
  options: string[];
  roomName: string;
  selfName: string;
  recent: Array<{ author: string; role: string; text: string }>;
}

export interface Mind {
  speak(c: Connector, ctx: TurnContext): Promise<string>;
  rank(c: Connector, question: string, options: string[]): Promise<number[]>;
  summarize(c: Connector, ctx: TurnContext): Promise<string>;
}

const KINDS: ConnectorKind[] = ["persona", "ollama", "openrouter", "webhook"];

export function createConnector(input: {
  name: string;
  kind: string;
  persona?: string;
  model?: string;
  url?: string;
}): Connector {
  const kind = input.kind as ConnectorKind;
  if (!KINDS.includes(kind)) throw new Error(`kind must be one of ${KINDS.join(", ")}`);
  if (!input.name?.trim()) throw new Error("name is required");
  if (kind === "webhook" && !/^https?:\/\//.test(input.url ?? "")) throw new Error("webhook needs an http(s) url");
  return {
    id: randomUUID(),
    name: input.name.trim(),
    kind,
    persona: (input.persona ?? "").trim() || "a thoughtful generalist who weighs trade-offs",
    model: input.model?.trim() || undefined,
    url: input.url?.trim() || undefined,
    createdAt: now(),
  };
}

// --- Offline persona: deterministic, free, always available --------------

function unit(seed: string): number {
  const h = createHash("sha256").update(seed).digest();
  return h.readUInt32BE(0) / 0xffffffff;
}

export function personaRank(c: Connector, question: string, options: string[]): number[] {
  const traits = new Set(words(c.persona));
  return options.map((opt) => {
    const overlap = words(opt).filter((w) => traits.has(w)).length;
    return Math.min(1, 0.15 + 0.7 * unit(`${c.persona}|${question}|${opt}`) + 0.3 * overlap);
  });
}

export function personaSpeak(c: Connector, ctx: TurnContext): string {
  const prefs = personaRank(c, ctx.question, ctx.options);
  const fav = argmax(prefs);
  const order = prefs.map((p, i) => [p, i] as const).sort((a, b) => b[0] - a[0]);
  const second = order[1]?.[1] ?? fav;
  const last = [...ctx.recent].reverse().find((m) => m.author !== ctx.selfName);
  const lens = clip(c.persona, 90);
  const turn = ctx.recent.filter((m) => m.author === ctx.selfName).length;
  if (!last) {
    return `Speaking as ${lens}: I'd start with "${ctx.options[fav]}". It answers "${clip(ctx.question, 80)}" with the fewest regrets.`;
  }
  const who = last.role === "surrogate" ? "the other room" : last.author;
  const theirs = ctx.options.findIndex((o) => last.text.toLowerCase().includes(o.toLowerCase()));
  if (theirs >= 0 && theirs !== fav) {
    if (prefs[theirs] >= prefs[fav] - 0.15) {
      return `${cap(who)} makes a fair case for "${ctx.options[theirs]}". As ${lens}, I could live with it; it is close to my "${ctx.options[fav]}".`;
    }
    return `I hear ${who} on "${ctx.options[theirs]}", but as ${lens} I still back "${ctx.options[fav]}". What would change your mind?`;
  }
  const lines = [
    `Building on ${who}: "${ctx.options[fav]}" first, with "${ctx.options[second]}" as the fallback.`,
    `If we can't agree, "${ctx.options[second]}" is my acceptable second. The risk I see with the rest is cost and time.`,
    `Our room should tell the others why "${ctx.options[fav]}" holds up from the angle of ${lens}.`,
  ];
  return lines[turn % lines.length];
}

function cap(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

// --- LLM-backed minds -----------------------------------------------------

function systemPrompt(c: Connector, ctx: TurnContext): string {
  return [
    `You are ${ctx.selfName}, a member of a small deliberation room (${ctx.roomName}) inside a larger swarm.`,
    `Persona: ${c.persona}.`,
    `The swarm is answering: "${ctx.question}". Options: ${ctx.options.map((o, i) => `${i + 1}. ${o}`).join("; ")}.`,
    `Reply with one or two short sentences, in the room's language. Name the option you argue for by its exact text.`,
    `React to what others said. Change your mind when an argument is better than yours.`,
  ].join("\n");
}

function transcript(ctx: TurnContext): string {
  if (ctx.recent.length === 0) return "(the room is quiet; open the discussion)";
  return ctx.recent.map((m) => `${m.author}${m.role === "surrogate" ? " [from another room]" : ""}: ${m.text}`).join("\n");
}

async function chat(c: Connector, system: string, user: string): Promise<string> {
  const signal = AbortSignal.timeout(Number(process.env.HIVE_MODEL_TIMEOUT_MS || 25000));
  if (c.kind === "ollama") {
    const base = c.url || process.env.OLLAMA_HOST || "http://127.0.0.1:11434";
    const res = await fetch(`${base.replace(/\/$/, "")}/api/chat`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        model: c.model || process.env.OLLAMA_MODEL || "llama3.2",
        stream: false,
        messages: [
          { role: "system", content: system },
          { role: "user", content: user },
        ],
      }),
      signal,
    });
    if (!res.ok) throw new Error(`ollama ${res.status}`);
    const body = (await res.json()) as { message?: { content?: string } };
    return body.message?.content?.trim() || "";
  }
  if (c.kind === "openrouter") {
    const key = process.env.OPENROUTER_API_KEY;
    if (!key) throw new Error("OPENROUTER_API_KEY is not set");
    const res = await fetch(c.url || "https://openrouter.ai/api/v1/chat/completions", {
      method: "POST",
      headers: { "content-type": "application/json", authorization: `Bearer ${key}` },
      body: JSON.stringify({
        model: c.model || process.env.OPENROUTER_MODEL || "moonshotai/kimi-k3",
        max_tokens: 300,
        messages: [
          { role: "system", content: system },
          { role: "user", content: user },
        ],
      }),
      signal,
    });
    if (!res.ok) throw new Error(`openrouter ${res.status}`);
    const body = (await res.json()) as { choices?: Array<{ message?: { content?: string } }> };
    return body.choices?.[0]?.message?.content?.trim() || "";
  }
  throw new Error(`chat not supported for ${c.kind}`);
}

// Webhook contract: POST JSON {type, ...context}; answer {text} or {scores}.
async function hook(c: Connector, payload: Record<string, unknown>): Promise<Record<string, unknown>> {
  const res = await fetch(c.url!, {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify({ persona: c.persona, ...payload }),
    signal: AbortSignal.timeout(Number(process.env.HIVE_MODEL_TIMEOUT_MS || 25000)),
  });
  if (!res.ok) throw new Error(`webhook ${res.status}`);
  return (await res.json()) as Record<string, unknown>;
}

export function parseScores(raw: string, n: number): number[] | null {
  const match = raw.match(/\[[\s\S]*?\]/);
  if (!match) return null;
  try {
    const arr = JSON.parse(match[0]);
    if (!Array.isArray(arr) || arr.length !== n) return null;
    const nums = arr.map(Number);
    if (nums.some((v) => !Number.isFinite(v))) return null;
    const max = Math.max(...nums, 1e-9);
    return nums.map((v) => Math.max(0, Math.min(1, max > 1 ? v / max : v)));
  } catch {
    return null;
  }
}

export const liveMind: Mind = {
  async speak(c, ctx) {
    if (c.kind === "persona") return personaSpeak(c, ctx);
    if (c.kind === "webhook") {
      const out = await hook(c, { type: "speak", ...ctx });
      return String(out.text ?? "").trim();
    }
    return chat(c, systemPrompt(c, ctx), `Room so far:\n${transcript(ctx)}\n\nYour turn, ${ctx.selfName}.`);
  },
  async rank(c, question, options) {
    if (c.kind === "persona") return personaRank(c, question, options);
    if (c.kind === "webhook") {
      const out = await hook(c, { type: "rank", question, options });
      const scores = Array.isArray(out.scores) ? parseScores(JSON.stringify(out.scores), options.length) : null;
      return scores ?? personaRank(c, question, options);
    }
    const raw = await chat(
      c,
      `Persona: ${c.persona}. Answer only with a JSON array of ${options.length} numbers between 0 and 1.`,
      `Question: ${question}\nRate how much you support each option, in order:\n${options.map((o, i) => `${i + 1}. ${o}`).join("\n")}`,
    );
    return parseScores(raw, options.length) ?? personaRank(c, question, options);
  },
  async summarize(c, ctx) {
    if (c.kind === "webhook") {
      const out = await hook(c, { type: "summarize", ...ctx });
      return String(out.text ?? "").trim();
    }
    return chat(
      c,
      `You are the conversational surrogate of ${ctx.roomName}. In two sentences, pass the room's strongest insight and its leaning option to a neighbouring room. Keep the speakers' language. Do not invent arguments.`,
      transcript(ctx),
    );
  },
};

export function turnContext(
  question: string,
  options: string[],
  roomName: string,
  selfName: string,
  msgs: Message[],
  limit = 12,
): TurnContext {
  return {
    question,
    options,
    roomName,
    selfName,
    recent: msgs.slice(-limit).map((m) => ({ author: m.authorName, role: m.role, text: m.text })),
  };
}
