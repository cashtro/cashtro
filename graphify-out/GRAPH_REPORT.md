# Graph Report - workspace  (2026-09-22)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 592 nodes · 1298 edges · 28 communities (25 shown, 3 thin omitted)
- Extraction: 98% EXTRACTED · 2% INFERRED · 0% AMBIGUOUS · INFERRED: 21 edges (avg confidence: 0.84)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `67ca60e3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Community 0
- Community 1
- Community 2
- Community 3
- Community 4
- Community 5
- Community 6
- Community 7
- Community 8
- Community 9
- Community 10
- Community 11
- Community 12
- Community 13
- Community 14
- Community 15
- Community 16
- Community 17
- Community 18
- Community 19
- Community 20
- Community 21
- Community 22
- Community 23
- Community 24
- Community 25
- Community 27

## God Nodes (most connected - your core abstractions)
1. `init()` - 129 edges
2. `Agent` - 117 edges
3. `Kernel` - 42 edges
4. `api` - 28 edges
5. `writeJSON()` - 24 edges
6. `Call` - 20 edges
7. `Catalog` - 18 edges
8. `Project` - 16 edges
9. `Result` - 15 edges
10. `Spec` - 12 edges

## Surprising Connections (you probably didn't know these)
- `run()` --calls--> `New()`  [EXTRACTED]
  cmd/cashtro/main.go → internal/server/server.go
- `TestCardUnbound()` --calls--> `Card()`  [INFERRED]
  internal/model/openrouter_test.go → internal/model/openrouter.go
- `withDb()` --calls--> `buildApp()`  [EXTRACTED]
  apps/api/src/app.test.ts → apps/api/src/app.ts
- `Boot()` --calls--> `New()`  [EXTRACTED]
  internal/agents/builtin.go → internal/catalog/catalog.go
- `Boot()` --calls--> `FromEnv()`  [EXTRACTED]
  internal/agents/builtin.go → internal/model/openrouter.go

## Import Cycles
- None detected.

## Communities (28 total, 3 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.06
Nodes (50): deliveryAgent, researchAgent, residentAgent, go_pkg_context, go_pkg_encoding_json, go_pkg_errors, go_pkg_github_com_cashtro_cashtro_internal_catalog, go_pkg_github_com_cashtro_cashtro_internal_kernel (+42 more)

### Community 1 - "Community 1"
Cohesion: 0.10
Nodes (31): main(), run(), go_pkg_flag, go_pkg_github_com_cashtro_cashtro_internal_agents, go_pkg_github_com_cashtro_cashtro_internal_server, go_pkg_log, go_pkg_net, go_pkg_net_http (+23 more)

### Community 2 - "Community 2"
Cohesion: 0.14
Nodes (13): Agent, AgentRuntime, CreateTask, Project, ProjectKind, ProjectStatus, Run, RunStatus (+5 more)

### Community 3 - "Community 3"
Cohesion: 0.10
Nodes (24): ref_node_child_process, ref_node_timers, ref_node_util, CATALOG_SHIPS, Config, findRoot(), loadJSON(), main() (+16 more)

### Community 4 - "Community 4"
Cohesion: 0.11
Nodes (25): CreateShip, Option, Profile, Ship, Stage, go_pkg_sync, go_pkg_unicode, sync.RWMutex (+17 more)

### Community 5 - "Community 5"
Cohesion: 0.05
Nodes (116): Agent, AgentLoopArchitecturesAgent(), AgentToolingTrustBoundaryAgent(), AgentToolingTrustReviewAgent(), AiOsDashboardAgent(), ApiRpcPrivilegeReviewAgent(), ApiRpcPrivilegeReviewerAgent(), ArchitectAgent() (+108 more)

### Community 6 - "Community 6"
Cohesion: 0.24
Nodes (6): net/http.Request, net/http.ResponseWriter, decodeJSON(), writeError(), writeJSON(), api

### Community 7 - "Community 7"
Cohesion: 0.10
Nodes (20): bin, recon, dependencies, zod, devDependencies, tsx, @types/node, typescript (+12 more)

### Community 8 - "Community 8"
Cohesion: 0.10
Nodes (19): dependencies, zod, devDependencies, tsx, @types/node, typescript, tsx, @types/node (+11 more)

### Community 9 - "Community 9"
Cohesion: 0.16
Nodes (9): go_pkg_fmt, go_pkg_sort, go_pkg_strings, clip(), Kernel, Confirm, Fact, Mail (+1 more)

### Community 10 - "Community 10"
Cohesion: 0.18
Nodes (16): collections, datetime, json, pathlib, re, classify_projects(), go_string(), ident() (+8 more)

### Community 11 - "Community 11"
Cohesion: 0.21
Nodes (13): routerAgent, go_pkg_bytes, net/http.Client, Bound(), Card(), FromEnv(), Client, chatAPIRequest (+5 more)

### Community 12 - "Community 12"
Cohesion: 0.12
Nodes (16): Project, Work, AxelProject(), BtkAvocatsProject(), CashtroCatalogProject(), CashtroProject(), EduconnexionProject(), EvoluJeunesProject() (+8 more)

### Community 13 - "Community 13"
Cohesion: 0.17
Nodes (12): AppOpts, buildApp(), fastify, FastifyRequest, jsonArr(), parseKeys(), port, prisma (+4 more)

### Community 14 - "Community 14"
Cohesion: 0.15
Nodes (12): tsx, @types/node, typescript, zod, name, prisma, seed, private (+4 more)

### Community 15 - "Community 15"
Cohesion: 0.15
Nodes (12): compilerOptions, declaration, esModuleInterop, forceConsistentCasingInFileNames, module, moduleResolution, outDir, resolveJsonModule (+4 more)

### Community 16 - "Community 16"
Cohesion: 0.24
Nodes (8): KERNEL_AGENTS, prisma, seed(), exec, withDb(), ref_node_fs, ref_node_os, ref_node_path

### Community 17 - "Community 17"
Cohesion: 0.17
Nodes (11): name, packageManager, private, scripts, api, api:dev, db:migrate, db:seed (+3 more)

### Community 18 - "Community 18"
Cohesion: 0.49
Nodes (9): build(), cmd_start(), cmd_status(), health_ok(), log(), run_loop(), voltron.sh script, supervisor_running() (+1 more)

### Community 20 - "Community 20"
Cohesion: 0.22
Nodes (8): compilerOptions, module, moduleResolution, outDir, rootDir, extends, include, ../../tsconfig.base.json

### Community 21 - "Community 21"
Cohesion: 0.29
Nodes (7): dependencies, @cashtro/sdk, fastify, @fastify/swagger, @fastify/swagger-ui, @prisma/client, zod

### Community 22 - "Community 22"
Cohesion: 0.29
Nodes (7): scripts, dev, prisma:migrate, prisma:seed, start, test, typecheck

### Community 23 - "Community 23"
Cohesion: 0.29
Nodes (6): compilerOptions, outDir, rootDir, extends, include, ../../tsconfig.base.json

### Community 24 - "Community 24"
Cohesion: 0.29
Nodes (6): compilerOptions, outDir, rootDir, extends, include, ../../tsconfig.base.json

### Community 25 - "Community 25"
Cohesion: 0.40
Nodes (5): devDependencies, prisma, tsx, @types/node, typescript

## Knowledge Gaps
- **117 isolated node(s):** `name`, `version`, `private`, `type`, `dev` (+112 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 151 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Kernel` connect `Community 0` to `Community 1`, `Community 4`, `Community 6`, `Community 9`, `Community 11`?**
  _High betweenness centrality (0.046) - this node is a cross-community bridge._
- **Why does `api` connect `Community 6` to `Community 0`, `Community 4`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **Why does `init()` connect `Community 5` to `Community 12`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `name`, `version`, `private` to the rest of the system?**
  _117 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `Community 0` be split into smaller, more focused modules?**
  _Cohesion score 0.05938037865748709 - nodes in this community are weakly interconnected._
- **Should `Community 1` be split into smaller, more focused modules?**
  _Cohesion score 0.10252100840336134 - nodes in this community are weakly interconnected._
- **Should `Community 2` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._