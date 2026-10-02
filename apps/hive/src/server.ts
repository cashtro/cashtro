import { readFileSync } from "node:fs";
import http, { type IncomingMessage, type ServerResponse } from "node:http";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { Hive, NotFound, type HiveEvent } from "./hive.js";

const here = path.dirname(fileURLToPath(import.meta.url));
const UI = readFileSync(path.join(here, "ui.html"), "utf8");

type Params = Record<string, string>;
type Handler = (ctx: { req: IncomingMessage; res: ServerResponse; params: Params; body: any }) => unknown;

class HttpError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
  }
}

export function createHiveServer(hive: Hive, opts: { token?: string } = {}): http.Server {
  const token = opts.token ?? process.env.HIVE_TOKEN ?? "";
  const routes: Array<{ method: string; pattern: RegExp; keys: string[]; handler: Handler }> = [];
  const route = (method: string, spec: string, handler: Handler) => {
    const keys: string[] = [];
    const pattern = new RegExp(
      "^" + spec.replace(/:(\w+)/g, (_, k) => (keys.push(k), "([^/]+)")) + "/?$",
    );
    routes.push({ method, pattern, keys, handler });
  };

  route("GET", "/health", () => ({
    ok: true,
    swarms: hive.listSwarms().length,
    connectors: hive.listConnectors().length,
  }));

  route("GET", "/api/swarms", () =>
    hive.listSwarms().map((s) => ({
      id: s.id,
      name: s.name,
      question: s.question,
      options: s.options,
      status: s.status,
      participants: s.participants.length,
      rooms: s.rooms.length,
      decision: s.decision.status,
      winner: s.decision.winner === null ? null : s.options[s.decision.winner],
      createdAt: s.createdAt,
    })),
  );
  route("POST", "/api/swarms", ({ res, body }) => {
    res.statusCode = 201;
    return hive.createSwarm({
      name: str(body.name),
      question: str(body.question),
      options: Array.isArray(body.options) ? body.options.map(String) : str(body.options).split(/\n|,/),
      roomSize: body.roomSize === undefined ? undefined : Number(body.roomSize),
      surrogateConnectorId: body.surrogateConnectorId || undefined,
    });
  });
  route("GET", "/api/swarms/:id", ({ params }) => hive.getSwarm(params.id));
  route("PATCH", "/api/swarms/:id", ({ params, body }) =>
    hive.updateSwarm(params.id, {
      surrogateConnectorId: body.surrogateConnectorId,
      roomSize: body.roomSize === undefined ? undefined : Number(body.roomSize),
    }),
  );
  route("DELETE", "/api/swarms/:id", ({ params }) => (hive.deleteSwarm(params.id), { ok: true }));
  route("POST", "/api/swarms/:id/start", ({ params }) => hive.setStatus(params.id, "live"));
  route("POST", "/api/swarms/:id/pause", ({ params }) => hive.setStatus(params.id, "draft"));
  route("POST", "/api/swarms/:id/close", ({ params }) => hive.setStatus(params.id, "closed"));
  route("GET", "/api/swarms/:id/summary", ({ params }) => hive.summary(params.id));
  route("GET", "/api/swarms/:id/export", ({ params, res }) => {
    const s = hive.getSwarm(params.id);
    res.setHeader("content-disposition", `attachment; filename="swarm-${s.id.slice(0, 8)}.json"`);
    return { swarm: s, summary: hive.summary(params.id) };
  });

  route("POST", "/api/swarms/:id/participants", ({ params, body, res }) => {
    const kind = body.kind === "agent" ? "agent" : "human";
    res.statusCode = 201;
    return hive.join(params.id, { name: str(body.name), kind, connectorId: body.connectorId });
  });
  route("DELETE", "/api/swarms/:id/participants/:pid", ({ params }) => (hive.leave(params.id, params.pid), { ok: true }));
  route("POST", "/api/swarms/:id/messages", ({ params, body, res }) => {
    res.statusCode = 201;
    return hive.say(params.id, str(body.participantId), str(body.text));
  });

  route("POST", "/api/swarms/:id/decision/start", ({ params }) => hive.startDecision(params.id).then((s) => s.decision));
  route("POST", "/api/swarms/:id/decision/pull", ({ params, body }) => {
    const target = body.target === null || body.target === undefined ? null : Number(body.target);
    hive.pull(params.id, str(body.participantId), target, body.conviction === undefined ? 1 : Number(body.conviction));
    return { ok: true };
  });

  route("GET", "/api/connectors", () => hive.listConnectors());
  route("POST", "/api/connectors", ({ body, res }) => {
    res.statusCode = 201;
    return hive.addConnector(body);
  });
  route("DELETE", "/api/connectors/:id", ({ params }) => (hive.deleteConnector(params.id), { ok: true }));
  route("POST", "/api/connectors/:id/test", ({ params }) => hive.testConnector(params.id));

  return http.createServer(async (req, res) => {
    const url = new URL(req.url || "/", "http://hive");
    try {
      if (req.method === "GET" && (url.pathname === "/" || url.pathname === "/index.html")) {
        res.writeHead(200, { "content-type": "text/html; charset=utf-8" });
        return res.end(UI);
      }
      const stream = url.pathname.match(/^\/api\/swarms\/([^/]+)\/stream$/);
      if (req.method === "GET" && stream) return sse(hive, stream[1], req, res);

      if (token && req.method !== "GET") {
        const auth = req.headers.authorization || "";
        if (auth !== `Bearer ${token}`) throw new HttpError(401, "unauthorized");
      }
      for (const r of routes) {
        if (r.method !== req.method) continue;
        const m = url.pathname.match(r.pattern);
        if (!m) continue;
        const params: Params = {};
        r.keys.forEach((k, i) => (params[k] = decodeURIComponent(m[i + 1])));
        const body = req.method === "GET" || req.method === "DELETE" ? {} : await readJson(req);
        const out = await r.handler({ req, res, params, body });
        return json(res, res.statusCode || 200, out);
      }
      throw new HttpError(404, "not found");
    } catch (err) {
      const status = err instanceof HttpError ? err.status : err instanceof NotFound ? 404 : 400;
      return json(res, status, { error: (err as Error).message });
    }
  });
}

function sse(hive: Hive, id: string, req: IncomingMessage, res: ServerResponse): void {
  const swarm = hive.getSwarm(id);
  res.writeHead(200, {
    "content-type": "text/event-stream",
    "cache-control": "no-cache",
    connection: "keep-alive",
  });
  const send = (ev: HiveEvent) => res.write(`data: ${JSON.stringify(ev)}\n\n`);
  send({ type: "swarm", swarm });
  hive.events.on(id, send);
  const ping = setInterval(() => res.write(": ping\n\n"), 15000);
  req.on("close", () => {
    clearInterval(ping);
    hive.events.off(id, send);
  });
}

function str(v: unknown): string {
  return typeof v === "string" ? v : v === undefined || v === null ? "" : String(v);
}

function json(res: ServerResponse, status: number, body: unknown): void {
  res.writeHead(status, { "content-type": "application/json" });
  res.end(JSON.stringify(body));
}

async function readJson(req: IncomingMessage): Promise<any> {
  let raw = "";
  for await (const chunk of req) {
    raw += chunk;
    if (raw.length > 1_000_000) throw new HttpError(413, "body too large");
  }
  if (!raw) return {};
  try {
    return JSON.parse(raw);
  } catch {
    throw new HttpError(400, "invalid json");
  }
}

const isMain = process.argv[1] && path.resolve(process.argv[1]) === fileURLToPath(import.meta.url);
if (isMain) {
  const dataFile = process.env.HIVE_DATA || path.resolve(here, "..", "data", "hive.json");
  const hive = new Hive({ dataFile });
  hive.start();
  const port = Number(process.env.HIVE_PORT || 8790);
  const host = process.env.HIVE_HOST || "127.0.0.1";
  createHiveServer(hive).listen(port, host, () => {
    console.log(`Hive · collective superintelligence on http://${host}:${port}`);
  });
  const shutdown = () => {
    hive.stop();
    process.exit(0);
  };
  process.on("SIGINT", shutdown);
  process.on("SIGTERM", shutdown);
}
