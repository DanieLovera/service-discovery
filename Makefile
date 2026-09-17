SHELL := /bin/sh
.DEFAULT_GOAL := help

GO ?= go
BIN_DIR ?= bin
SERVICE ?= registry
DOCKER_REPOSITORY ?= daniel-tpiii
IMAGE ?= $(DOCKER_REPOSITORY)/$(SERVICE):dev
VALID_SERVICES := registry load-balancer mock-service

.PHONY: help fmt fmt-check vet staticcheck test test-race test-integration quality _check-service build build-all run docker-build clean

help:
	@printf '%s\n' \
		'make fmt              Format Go source files' \
		'make fmt-check        Fail if Go source files are not formatted' \
		'make vet              Run go vet' \
		'make staticcheck      Run staticcheck (must be installed)' \
		'make test             Run unit tests' \
		'make test-race        Run tests with the race detector' \
		'make test-integration Run integration test packages when present' \
		'make quality          Run formatting, vet, staticcheck and race tests' \
		'make build-all        Build all process entrypoints' \
		'make run SERVICE=...  Run registry, load-balancer, or mock-service' \
		'make docker-build SERVICE=... Build one service image'

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

_check-service:
	@case " $(VALID_SERVICES) " in \
		*" $(SERVICE) "*) ;; \
		*) echo "Error: SERVICE must be one of: $(VALID_SERVICES) (got '$(SERVICE)')"; exit 1 ;; \
	esac

build: _check-service
	@mkdir -p $(BIN_DIR)
	@CGO_ENABLED=0 $(GO) build -trimpath -o $(BIN_DIR)/$(SERVICE) ./cmd/$(SERVICE)

build-all:
	@mkdir -p $(BIN_DIR)
	@for service in $(VALID_SERVICES); do \
		CGO_ENABLED=0 $(GO) build -trimpath -o $(BIN_DIR)/$$service ./cmd/$$service || exit 1; \
	done

run: _check-service
	@$(GO) run ./cmd/$(SERVICE)

docker-build: _check-service
	@docker build --build-arg SERVICE=$(SERVICE) -t $(IMAGE) .

clean:
	@rm -rf $(BIN_DIR)
