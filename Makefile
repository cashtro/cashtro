.PHONY: test vet run build always-on ultron fleet-on

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/cashtro

ultron:
	go run ./cmd/ultron

build:
	mkdir -p bin
	go build -o bin/cashtro ./cmd/cashtro
	go build -o bin/ultron ./cmd/ultron

always-on:
	./scripts/always-on.sh

ultron-on:
	./scripts/ultron-on.sh

fleet-on:
	./scripts/fleet-on.sh
