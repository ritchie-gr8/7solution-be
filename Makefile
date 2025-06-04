SHELL=cmd.exe
API_BINARY=apiApp

## up: starts all containers in the background without forcing build
up:
		@echo Starting Docker images...
		docker-compose up -d
		@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_api
		@echo Stopping docker images (if running...)
		docker-compose down
		@echo Building (when required) and starting docker images...
		docker-compose up --build -d
		@echo Docker images built and started!

## down: stop docker compose
down:
		@echo Stopping docker compose...
		docker-compose down
		@echo Done!

## build_api: builds the api binary as a linux executable
build_api:
		@echo Building API binary...
		set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${API_BINARY} ./cmd/server
		@echo Done!

## start_api: starts the API service locally
start_api:
		@echo Starting API service...
		go run ./cmd/server
		@echo Done!

## test: run all tests
test:
		@echo Running tests...
		go test -v ./...
		@echo Done!
