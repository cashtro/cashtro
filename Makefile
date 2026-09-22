.PHONY: test vet run build voltron voltron-status

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/cashtro

build:
	go build -o bin/cashtro ./cmd/cashtro

# Voltron = Cashtro OS under a restart-forever supervisor.
voltron:
	./scripts/voltron.sh start

voltron-status:
	./scripts/voltron.sh status

# Graphify tree of Cashtro agentics and the Evolu-Jeunes fleet.
graphify:
	python3 -m graphify extract ./internal/graphifytree --code-only --no-cluster --out .
	python3 -m graphify cluster-only . --no-label --no-viz
	python3 -m graphify tree --graph graphify-out/graph.json --output graphify-out/GRAPH_TREE.html --label "Cashtro × Evolu-Jeunes"
	python3 -c 'from pathlib import Path; p=Path("graphify-out/GRAPH_TREE.html"); t=p.read_text(); n="window.expandAll = () => { expandBranch(rootNode); updateTree(rootNode); };"; p.write_text(t.replace(n, n+"\n    expandAll();", 1) if "\n    expandAll();" not in t else t)'
	python3 -m graphify export svg --graph graphify-out/graph.json
	python3 -m graphify export html --graph graphify-out/graph.json
	python3 -m graphify export callflow-html --graph graphify-out/graph.json --output graphify-out/GRAPH_CALLFLOW.html
