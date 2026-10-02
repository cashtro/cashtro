# Cashtro Hive

A platform for hosting **collective superintelligence** sessions. It follows the
ideas in Louis Rosenberg's talk (Unanimous AI): use AI to connect people
rather than replace them, so a group can solve problems its members can't solve
alone.

Hive puts two of those mechanisms in one place:

| Mechanism | Inspired by | What Hive does |
| --- | --- | --- |
| **Conversational swarm** | Thinkscape | Splits a large group into small rooms (default 5). Each room gets an AI **surrogate** that relays the room's strongest point and leaning option to the next room, so ideas travel through the whole group in a ring. |
| **Swarm decision** | Swarm AI | Every member pulls a shared puck toward the option they prefer, as hard as their conviction. The group converges in real time. Agents watch where the puck is heading and switch to it when it's an option they can accept. |

## Host, manage, use, create, connect

- **Create** swarms: a question, 2–8 options and a room size.
- **Host and manage** them: go live, pause, close, export the full transcript as JSON, or delete.
- **Use** them: join as a human, talk in your room, and pull the puck.
- **Connect** new minds as members or as the room surrogate:
  - `persona`: offline and deterministic. It's free and always available, and three starter personas are pre-loaded.
  - `ollama`: a local model (`OLLAMA_HOST`, `OLLAMA_MODEL`, default `llama3.2`).
  - `openrouter`: a cloud model (`OPENROUTER_API_KEY`, `OPENROUTER_MODEL`).
  - `webhook`: any HTTP agent, such as a Cashtro OS agent or your own bot.

### Webhook contract

Hive sends `POST <url>` with JSON:

```json
{ "type": "speak" | "rank" | "summarize",
  "persona": "…", "question": "…", "options": ["…"],
  "roomName": "Room 2", "selfName": "Nora",
  "recent": [{ "author": "Ana", "role": "member", "text": "…" }] }
```

Answer `{ "text": "…" }` for `speak` and `summarize`, or
`{ "scores": [0.2, 0.9, …] }` (one number per option) for `rank`.

## Run

```bash
pnpm install
pnpm hive                 # http://127.0.0.1:8790
pnpm --filter @cashtro/hive test
```

| Env | Default | Meaning |
| --- | --- | --- |
| `HIVE_PORT` / `HIVE_HOST` | `8790` / `127.0.0.1` | Bind address |
| `HIVE_DATA` | `apps/hive/data/hive.json` | Where swarms and connectors are saved |
| `HIVE_TOKEN` | unset | When set, every write needs `Authorization: Bearer <token>` |
| `HIVE_CONVO_MS` | `4000` | How often agents speak and surrogates relay |
| `HIVE_AGENT_BUDGET` | `40` | Agent messages allowed per swarm between two human messages, which caps model spend |
| `HIVE_MODEL_TIMEOUT_MS` | `25000` | Per-call timeout for model and webhook agents |

## API

| Method | Path | What it does |
| --- | --- | --- |
| `GET` | `/health` | Liveness and counts |
| `GET` / `POST` | `/api/swarms` | List or create swarms |
| `GET` / `PATCH` / `DELETE` | `/api/swarms/{id}` | Read, set surrogate or room size, delete |
| `POST` | `/api/swarms/{id}/start` · `/pause` · `/close` | Lifecycle |
| `GET` | `/api/swarms/{id}/stream` | Live events (SSE) |
| `GET` | `/api/swarms/{id}/summary` · `/export` | Collective answer and insights; full JSON export |
| `POST` / `DELETE` | `/api/swarms/{id}/participants[/{pid}]` | Join a human or an agent (`{kind, name, connectorId}`); leave |
| `POST` | `/api/swarms/{id}/messages` | `{participantId, text}` into the member's room |
| `POST` | `/api/swarms/{id}/decision/start` · `/decision/pull` | Start the puck; `{participantId, target, conviction}` |
| `GET` / `POST` | `/api/connectors` | List or connect minds |
| `DELETE` | `/api/connectors/{id}` | Disconnect (refused while a swarm uses it) |
| `POST` | `/api/connectors/{id}/test` | Ask the mind a test question |

## Layout

- `src/engine.ts`: pure swarm mechanics (rooms, mention detection, surrogate insight, puck physics)
- `src/connectors.ts`: persona, Ollama, OpenRouter and webhook minds
- `src/hive.ts`: hosts the swarms and runs the conversation and decision loops, persistence and events
- `src/server.ts`: HTTP API, SSE and the web UI (`src/ui.html`). No runtime dependencies.
