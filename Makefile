.PHONY: test vet run build always-on

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/cashtro

build:
	mkdir -p bin
	go build -o bin/cashtro ./cmd/cashtro

always-on:
	./scripts/always-on.sh
