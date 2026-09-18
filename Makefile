SHELL := /bin/sh
.DEFAULT_GOAL := help

GO ?= go

.PHONY: help fmt fmt-check vet staticcheck test test-race test-integration quality \
		docker-up docker-up-debug docker-down docker-ps docker-images docker-clean

help:
	@printf '%s\n' \
		'make fmt              	Format Go source files' \
		'make fmt-check        	Fail if Go source files are not formatted' \
		'make vet              	Run go vet' \
		'make staticcheck      	Run staticcheck (must be installed)' \
		'make test             	Run unit tests' \
		'make test-race        	Run tests with the race detector' \
		'make test-integration 	Run integration test packages when present' \
		'make quality          	Run formatting, vet, staticcheck and race tests' \
		'make docker-up        	Build and start the Docker Compose environment' \
		'make docker-up-debug	Build and start the Docker Compose environment with debug tools' \
		'make docker-down      	Stop and remove the Docker Compose environment' \
		'make docker-images     Show Docker Compose service images' \		
		'make docker-ps        	Show Docker Compose service status' \
		'make docker-clean     	Remove Compose containers, local images and volumes'

fmt:
	@$(GO) fmt ./...

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Go files require formatting:"; gofmt -l .; exit 1)

vet:
	@$(GO) vet ./...

staticcheck:
	@command -v staticcheck >/dev/null 2>&1 || { echo "staticcheck is not installed"; exit 1; }
	@staticcheck ./...

test:
	@$(GO) test ./...

test-race:
	@$(GO) test -race ./...

test-integration:
	@if [ -n "$$(find ./test/integration -name '*_test.go' | head -n 1)" ]; then \
		$(GO) test -race ./test/integration/...; \
	else \
		echo "No integration tests implemented yet."; \
	fi

quality: fmt-check vet staticcheck test-race

docker-up:
	@docker compose up --build -d

docker-up-debug:
	@docker compose --profile debug up --build -d

docker-down:
	@docker compose --profile debug down

docker-ps:
	@docker compose --profile debug ps

docker-images:
	@docker compose --profile debug images

docker-clean:
	@docker compose --profile debug down --rmi local -v
