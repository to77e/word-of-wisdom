LOCAL_BIN:=$(CURDIR)/bin
GOLANGCI_BIN:=$(LOCAL_BIN)/golangci-lint



.PHONY: run
run:
	go run -race cmd/app/main.go

.PHONY: build
build:
	go build -o bin/word-of-wisdom cmd/app/main.go

.PHONY: generate
generate:
	go generate ./...



.PHONY: test_unit
test_unit:
	go test -v -timeout 5s -count 1 -race -run Unit ./...

.PHONY: test_unit_multi
test_unit_multi:
	go test -v -timeout 5s -count 30 -race -run Unit ./...

.PHONY: test_integration
test_integration:
	go test -v -timeout 10s -count 1 -race -run Integration ./...

.PHONY: test
test:
	go test -v -count 1 -race ./...



.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(GOLANGCI_BIN)),)
	$(info Downloading golangci-lint)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif

.PHONY: lint
lint: install-lint
	$(GOLANGCI_BIN) run --config=.golangci.yaml ./...



.PHONY: deps
deps:
	$(info Install dependencies...)
	go mod tidy
	go mod download



.PHONY: env
env:
	docker compose -f docker-compose.yaml up -d --force-recreate

.PHONY: env_down
env_down:
	docker compose -f docker-compose.yaml down -v --remove-orphans