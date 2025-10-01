APP=server

.PHONY: run build test tidy

run:
	@echo "PORT=$${PORT:-8080}"
	PORT=$${PORT:-8080} go run ./cmd/server

build:
	GOOS=$${GOOS:-$(shell go env GOOS)} GOARCH=$${GOARCH:-$(shell go env GOARCH)} \
	go build -o ./bin/$(APP) ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy
