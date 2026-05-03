run:
	go run ./cmd/reqium

test:
	go test ./...

build:
	go build -o bin/reqium ./cmd/reqium

fmt:
	go fmt ./...
