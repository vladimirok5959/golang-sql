default: test

clean:
	go clean -testcache ./...

test:
	go test ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

.PHONY: default clean test lint tidy
