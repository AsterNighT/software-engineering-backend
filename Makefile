
.PHONY: all build run check fmt lint test clean setup-tools swag start
start: swag build run

build:
	go build -o build/main ./cmd/_rename_this_/main.go

run:
	build/main

check: fmt lint test

fmt:
	go fmt ./...

lint: 
	golangci-lint run

test:
	go test ./...

swag:
	swag init -g pkg/router/router.go

setup-tools:
	go get github.com/swaggo/swag/cmd/swag
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0