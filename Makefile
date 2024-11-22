VERSION := $(shell git describe --tags --always)

run: build
	@./bin/goredisclone

build:
	@go build -ldflags "-X main.version=$(VERSION)" -o bin/goredisclone ./cmd/redis-server