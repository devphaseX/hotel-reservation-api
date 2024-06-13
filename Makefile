dev: build run

build:
	@go build -o ./bin/api
.PHONY:build

run: 
	@./bin/api --listenAddress :5000
.PHONY:run

test: 
	@go test -v ./...
.PHONY:test