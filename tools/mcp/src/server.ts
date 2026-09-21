#!/usr/bin/env node
import { createInterface } from "node:readline";
import { callTool, findRoot, TOOLS } from "./tools.js";

const root = await findRoot(process.cwd());

const rl = createInterface({ input: process.stdin, crlfDelay: Infinity });
for await (const line of rl) {
  if (!line.trim()) continue;
  let msg: { jsonrpc?: string; id?: unknown; method?: string; params?: { name?: string; arguments?: Record<string, unknown> } };
  try {
    msg = JSON.parse(line) as typeof msg;
  } catch {
    continue;
  }
  const id = msg.id;
  if (msg.method === "initialize") {
    write({ jsonrpc: "2.0", id, result: { protocolVersion: "2024-11-05", capabilities: { tools: {} }, serverInfo: { name: "cashtro-registry", version: "0.1.0" } } });
    continue;
  }
  if (msg.method === "tools/list") {
    write({ jsonrpc: "2.0", id, result: { tools: TOOLS } });
    continue;
  }
  if (msg.method === "tools/call") {
    try {
      const name = msg.params?.name || "";
      const data = await callTool(root, name, msg.params?.arguments || {});
      write({ jsonrpc: "2.0", id, result: { content: [{ type: "text", text: JSON.stringify(data, null, 2) }] } });
    } catch (err) {
      write({ jsonrpc: "2.0", id, error: { code: -32000, message: err instanceof Error ? err.message : "tool failed" } });
    }
    continue;
  }
  if (msg.method?.startsWith("notifications/")) continue;
  write({ jsonrpc: "2.0", id, error: { code: -32601, message: `unknown method ${msg.method}` } });
}

function write(obj: unknown) {
  process.stdout.write(JSON.stringify(obj) + "\n");
}
