import type { AgentAdapter } from "./types.js";
import { createHttpAdapter } from "./http.js";
import { createOpenRouterAdapter, createBedrockAdapter } from "./openrouter.js";
import { createClaudeCodeAdapter, createCursorCliAdapter } from "./cli.js";
import { createLocalAdapter } from "./local.js";
import { createN8nAdapter } from "./n8n.js";

export function defaultAdapters(): Record<string, AgentAdapter> {
  const openrouter = createOpenRouterAdapter();
  const http = createHttpAdapter();
  return {
    http,
    local: createLocalAdapter(),
    n8n: createN8nAdapter({ fallback: http }),
    openrouter,
    bedrock: createBedrockAdapter(),
    cursor: createCursorCliAdapter(),
    "claude-code": createClaudeCodeAdapter(),
  };
}

/** Cheapest capable adapter. HTTP (Voltron) first — $0. Model route only if asked. */
export function pickAdapter(runtime: string, adapters: Record<string, AgentAdapter> = defaultAdapters()): AgentAdapter {
  if (runtime === "n8n" && adapters.n8n) return adapters.n8n;
  if (runtime === "http" && adapters.http) return adapters.http;
  if (runtime === "bedrock" && adapters.bedrock) return adapters.bedrock;
  if (runtime === "cursor" && adapters.cursor) return adapters.cursor;
  if (runtime === "claude-code" && adapters["claude-code"]) return adapters["claude-code"];
  if (adapters[runtime]) return adapters[runtime];
  return adapters.http;
}
