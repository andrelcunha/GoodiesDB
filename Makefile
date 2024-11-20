run: build
	@./bin/goredisclone

build:
	@go build -o bin/goredisclone ./cmd/redis-server