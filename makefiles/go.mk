ifndef _include_go_mk
_include_go_mk = 1

include makefiles/shared.mk
include makefiles/gobin.mk

GO ?= go
FORMAT_FILES ?= .

GOFUMPT := bin/gofumpt

GOLANGCILINT := bin/golangci-lint
GOLANGCILINT_VERSION ?= 1.31.0
GOLANGCILINT_CONCURRENCY ?= 16

$(GOFUMPT): $(GOBIN)
	$(info $(_bullet) Installing <gofumpt>)
	@mkdir -p bin
	GOBIN=bin $(GOBIN) mvdan.cc/gofumpt

$(GOLANGCILINT): $(GOBIN)
	$(info $(_bullet) Installing <golangci-lint>)
	@mkdir -p bin
	GOBIN=bin $(GOBIN) github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCILINT_VERSION)

.PHONY: all deps format lint test test-coverage integration-test

all: format lint test-coverage integration-test

deps: ## Download dependencies
	$(info $(_bullet) Downloading dependencies)
	$(GO) mod download

format: $(GOFUMPT) ## Format code
	$(info $(_bullet) Formatting code)
	$(GOFUMPT) -w $(FORMAT_FILES)

lint: $(GOLANGCILINT) ## Lint code
	$(info $(_bullet) Running linter) 
	$(GOLANGCILINT) run --concurrency $(GOLANGCILINT_CONCURRENCY) ./...

test: ## Run tests
	$(info $(_bullet) Running tests)
	$(GO) test ./...
	
test-coverage: ## Run tests with coverage
	$(info $(_bullet) Running tests with coverage) 
	$(GO) test -cover ./...

integration-test: ## Run integration tests
	$(info $(_bullet) Running integration tests) 
	$(GO) test -tags integration -count 1 ./...

endif