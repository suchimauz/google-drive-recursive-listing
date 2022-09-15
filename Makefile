GOOS ?= windows

.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	go mod download && CGO_ENABLED=0 GOOS=$(GOOS) go build -o ./.bin/app.exe main.go

run:
	go run main.go

environment:
	cp env.dist .env

test:
	go test --short -coverprofile=cover.out -v ./...
	make test.coverage

test.coverage:
	go tool cover -func=cover.out

lint:
	golangci-lint run
