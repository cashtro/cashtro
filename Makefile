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
