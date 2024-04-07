build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags=release -ldflags="-s -w" -o bin/manager ./cmd

run:
	go run ./cmd

integration-test:
	go test -v ./tests/integration/...
