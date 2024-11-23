VERSION := $(shell git describe --tags --always)

run: build
	@./bin/goodiesdb-server

build:
	@go build -ldflags "-X main.version=$(VERSION)" -o bin/goodiesdb-server ./cmd/goodiesdb-server