// A connector is how a new mind joins the hive: an offline persona, a local
// Ollama model, an OpenRouter model, or any HTTP agent behind a webhook.
export type ConnectorKind = "persona" | "ollama" | "openrouter" | "webhook";

export interface Connector {
  id: string;
  name: string;
  kind: ConnectorKind;
  persona: string;
  model?: string;
  url?: string;
  createdAt: string;
}

export type ParticipantKind = "human" | "agent";

export interface Participant {
  id: string;
  name: string;
  kind: ParticipantKind;
  connectorId?: string;
  roomId: string;
  // Swarm decision: which option this member pulls toward, and how hard (0..1).
  target: number | null;
  conviction: number;
  // Agent preference per option (0..1), filled when a decision starts.
  prefs?: number[];
  joinedAt: string;
}

export type MessageRole = "member" | "surrogate" | "system";

export interface Message {
  id: string;
  roomId: string;
  authorId: string;
  authorName: string;
  role: MessageRole;
  text: string;
  // Options this message argues for, detected from its text.
  mentions: number[];
  at: string;
}

export interface Room {
  id: string;
  name: string;
  // Index into messages of the last message the surrogate already relayed.
  relayedUpTo: number;
}

export interface DecisionState {
  status: "idle" | "running" | "done";
  x: number;
  y: number;
  vx: number;
  vy: number;
  ticks: number;
  maxTicks: number;
  winner: number | null;
  // Share of total pull on each option at the last tick.
  support: number[];
  trace: Array<[number, number]>;
  startedAt?: string;
  endedAt?: string;
}

export type SwarmStatus = "draft" | "live" | "closed";

export interface Swarm {
  id: string;
  name: string;
  question: string;
  options: string[];
  roomSize: number;
  status: SwarmStatus;
  surrogateConnectorId?: string;
  participants: Participant[];
  rooms: Room[];
  messages: Message[];
  decision: DecisionState;
  createdAt: string;
}

export interface HiveState {
  swarms: Swarm[];
  connectors: Connector[];
}
