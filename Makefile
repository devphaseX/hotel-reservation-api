dev: build run

build:
	@go build -o ./bin/api

run: 
	@./bin/api --listenAddress :5000

test: 
	@go test -v ./...