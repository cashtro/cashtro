.PHONY: test vet run build

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/cashtro

build:
	go build -o bin/cashtro ./cmd/cashtro
