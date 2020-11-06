OS := $(shell uname -s | tr [:upper:] [:lower:])

GOBIN := bin/gobin
GOBIN_VERSION := 0.0.14

SQLC := bin/sqlc
SQLC_VERSION ?= 1.5.0

OPENAPIGENERATORCLI := scripts/openapi-generator-cli
OPENAPIGENERATORCLI_VERSION ?= 4.3.1

GOLANGCILINT := bin/golangci-lint
GOLANGCILINT_VERSION ?= 1.31.0

GOFUMPT := bin/gofumpt

KIND := bin/kind
KIND_VERSION ?= 0.9.0
KIND_CLUSTER_NAME ?= local

SKAFFOLD := bin/skaffold
SKAFFOLD_VERSION ?= 1.16.0

KUBERNETES_VERSION ?= 1.17.11
KUBERNETES_CONTEXT := kind-$(KIND_CLUSTER_NAME)

bullet := $(shell printf "\033[34;1m▶\033[0m")

all: generate format lint test-coverage integration-test build

$(GOBIN): ; $(info $(bullet) Installing <gobin>)
	@mkdir -p bin
	curl -sSfL https://github.com/myitcv/gobin/releases/download/v$(GOBIN_VERSION)/$(OS)-amd64 -o $(GOBIN)
	chmod u+x $(GOBIN)

$(GOLANGCILINT): $(GOBIN) ; $(info $(bullet) Installing <golangci-lint>)
	@mkdir -p bin
	GOBIN=bin $(GOBIN) github.com/golangci/golangci-lint/cmd/golangci-lint@v$(GOLANGCILINT_VERSION)

$(GOFUMPT): $(GOBIN) ; $(info $(bullet) Installing <gofumpt>)
	@mkdir -p bin
	GOBIN=bin $(GOBIN) mvdan.cc/gofumpt

$(SQLC): ; $(info $(bullet) Installing <sqlc>)
	@mkdir -p bin
	curl -sSfL https://github.com/kyleconroy/sqlc/releases/download/v$(SQLC_VERSION)/sqlc-v$(SQLC_VERSION)-$(OS)-amd64.tar.gz | tar -C bin -xz
	chmod u+x $(SQLC)

$(KIND): ; $(info $(bullet) Installing <kind>)
	@mkdir -p bin
	curl -sSfL https://kind.sigs.k8s.io/dl/v$(KIND_VERSION)/kind-$(OS)-amd64 -o $(KIND)
	chmod u+x $(KIND)

$(SKAFFOLD): ; $(info $(bullet) Installing <skaffold>)
	@mkdir -p bin
	curl -sSfL https://storage.googleapis.com/skaffold/releases/v$(SKAFFOLD_VERSION)/skaffold-$(OS)-amd64 -o $(SKAFFOLD)
	chmod u+x $(SKAFFOLD)

help: ## Help
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clean: clean-kind clean-bin ## Clean targets

clean-bin: ; $(info $(bullet) Cleaning <bin>) ## Clean installed tools
	rm -rf bin/

clean-kind: $(KIND) ; $(info $(bullet) Cleaning <kind>) ## Delete cluster
	$(KIND) delete cluster --name $(KIND_CLUSTER_NAME) || exit 0

generate: generate-sqlc generate-openapi ## Generate code

generate-sqlc: $(SQLC) ; $(info $(bullet) Generating <sqlc>) ## Generate SQLC code
	$(SQLC) generate

generate-openapi: $(SQLC) ; $(info $(bullet) Generating <openapi>) ## Generate OpenAPI code
	OPENAPIGENERATORCLI_VERSION=$(OPENAPIGENERATORCLI_VERSION) $(OPENAPIGENERATORCLI) generate \
		--input-spec api/api.yaml \
		--output pkg/api \
		--generator-name go-experimental \
		--package-name=api \
		--additional-properties withGoCodegenComment \
		--import-mappings=uuid.UUID=github.com/google/uuid --type-mappings=UUID=uuid.UUID

format: $(GOFUMPT) ; $(info $(bullet) Formatting code) ## Format code
	$(GOFUMPT) -w .

lint: $(GOLANGCILINT) ; $(info $(bullet) Running linter) ## Lint code
	$(GOLANGCILINT) run ./...

test: ; $(info $(bullet) Running tests) ## Run tests
	go test ./...

test-coverage: ; $(info $(bullet) Running tests with coverage) ## Run tests with coverage
	go test -cover ./...

integration-test: ; $(info $(bullet) Running integration tests) ## Run integration tests
	go test -tags integration -count 1 ./...

build-goose: ## Build goose
	go build -o bin/goose ./cmd/goose

build-todo-service: ## Build todo-service
	go build -o bin/todo-service ./services/todo

build: build-goose build-todo-service ## Build all targets

bootstrap-kind: $(KIND); $(info $(bullet) Bootstraping <kind>) ## Bootstrap cluster in docker
	$(KIND) get clusters | grep -q $(KIND_CLUSTER_NAME) || \
	$(KIND) create cluster \
		--name $(KIND_CLUSTER_NAME) \
		--image kindest/node:v$(KUBERNETES_VERSION) \
		--wait 1m

bootstrap-deploy: $(SKAFFOLD); $(info $(bullet) Bootstraping <deploy>) ## Bootstrap infrastructure
	$(SKAFFOLD) build -q | $(SKAFFOLD) deploy -p bootstrap --kube-context=$(KUBERNETES_CONTEXT) --build-artifacts -

bootstrap: $(KIND) $(SKAFFOLD) bootstrap-kind bootstrap-deploy ## Bootstrap cluster with infrastructure
