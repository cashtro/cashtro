# Graph Report - workspace  (2026-09-22)

## Corpus Check
- cluster-only mode — file stats not available

## Summary
- 42 nodes · 75 edges · 5 communities
- Extraction: 100% EXTRACTED · 0% INFERRED · 0% AMBIGUOUS
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `26dd180b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- Community 0
- Community 1
- Community 2
- Community 3
- Community 4

## God Nodes (most connected - your core abstractions)
1. `Kernel()` - 16 edges
2. `OwnsEvoluJeunes()` - 11 edges
3. `Fleet()` - 9 edges
4. `OwnsCashtro()` - 6 edges
5. `Home` - 4 edges
6. `Map()` - 4 edges
7. `Delivery()` - 3 edges
8. `Explorer()` - 3 edges
9. `Planner()` - 3 edges
10. `Repo()` - 2 edges

## Surprising Connections (you probably didn't know these)
- `OwnsCashtro()` --calls--> `Repo()`  [EXTRACTED]
  castro/owns.go → cashtro/agents.go
- `OwnsCashtro()` --calls--> `Kernel()`  [EXTRACTED]
  castro/owns.go → cashtro/agents.go
- `OwnsCashtro()` --calls--> `ControlPlane()`  [EXTRACTED]
  castro/owns.go → cashtro/agents.go
- `OwnsEvoluJeunes()` --calls--> `Fleet()`  [EXTRACTED]
  castro/owns.go → evolujeunes/org.go
- `OwnsEvoluJeunes()` --calls--> `Member()`  [EXTRACTED]
  castro/owns.go → evolujeunes/org.go

## Import Cycles
- None detected.

## Communities (5 total, 0 thin omitted)

### Community 0 - "Community 0"
Cohesion: 0.29
Nodes (12): Architect(), Comms(), Deploy(), Init(), Investigator(), Kernel(), Memory(), Operator() (+4 more)

### Community 1 - "Community 1"
Cohesion: 0.42
Nodes (8): BTKAvocats(), CashtroDeliveryCatalog(), Educonnexion(), Fleet(), MDClinic(), Proximity(), ScanApp(), SolutionHypothequeQC()

### Community 2 - "Community 2"
Cohesion: 0.36
Nodes (7): ControlPlane(), Repo(), Home, Map(), OwnsCashtro(), go_pkg_github_com_cashtro_cashtro_internal_graphifytree_cashtro, go_pkg_github_com_cashtro_cashtro_internal_graphifytree_evolujeunes

### Community 3 - "Community 3"
Cohesion: 0.25
Nodes (8): Delivery(), Explorer(), Planner(), OwnsEvoluJeunes(), Member(), Org(), PrivateRepos(), Site()

### Community 4 - "Community 4"
Cohesion: 0.50
Nodes (3): TestMapHasBothHomesAndEveryAgentic(), go_pkg_testing, testing.T

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `OwnsEvoluJeunes()` connect `Community 3` to `Community 1`, `Community 2`?**
  _High betweenness centrality (0.409) - this node is a cross-community bridge._
- **Why does `Kernel()` connect `Community 0` to `Community 2`, `Community 3`?**
  _High betweenness centrality (0.234) - this node is a cross-community bridge._
- **Why does `Fleet()` connect `Community 1` to `Community 3`?**
  _High betweenness centrality (0.232) - this node is a cross-community bridge._