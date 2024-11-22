VERSION := $(shell git describe --tags --always)

run: build
	@./bin/goodiesDb

build:
	@go build -ldflags "-X main.version=$(VERSION)" -o bin/goodiesDb ./cmd/goodiesdb-server