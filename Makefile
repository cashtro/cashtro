.PHONY: test vet run build voltron voltron-status graphify graphify-map

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

# Rebuild the agent × project Worked() map, then extract a local graph.
graphify-map:
	python3 tools/graphify-fleet/build.py
	gofmt -w internal/fleet/map.go

graphify: graphify-map
	graphify extract . --code-only
	graphify cluster-only . --no-label
