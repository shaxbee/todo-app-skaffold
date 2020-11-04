OS := $(shell uname -s | tr [:upper:] [:lower:])

SQLC := bin/sqlc
SQLC_VERSION ?= 1.5.0

GOLANGCILINT := bin/golangci-lint
GOLANGCILINT_VERSION ?= 1.31.0

OPENAPIGENERATORCLI := scripts/openapi-generator-cli
OPENAPIGENERATORCLI_VERSION :=

KIND := bin/kind
KIND_VERSION ?= 0.9.0
KIND_CLUSTER_NAME ?= local

SKAFFOLD := bin/skaffold
SKAFFOLD_VERSION := 1.16.0

KUBERNETES_VERSION ?= 1.17.11
KUBERNETES_CONTEXT := kind-$(KIND_CLUSTER_NAME)

bullet := $(shell printf "\033[34;1m▶\033[0m")

all: generate lint test

$(SQLC): ; $(info $(bullet) Installing <sqlc>)
	@mkdir -p bin
	curl -sSfL https://github.com/kyleconroy/sqlc/releases/download/v$(SQLC_VERSION)/sqlc-v$(SQLC_VERSION)-$(OS)-amd64.tar.gz | tar -C bin -xvz
	chmod u+x $(SQLC)

$(GOLANGCILINT): ; $(info $(bullet) Installing <golangci-lint>)
	@mkdir -p bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b bin v$(GOLANGCILINT_VERSION)

$(KIND): ; $(info $(bullet) Installing <kind>)
	@mkdir -p bin
	curl -sSfL https://kind.sigs.k8s.io/dl/v$(KIND_VERSION)/kind-$(OS)-amd64 -o $(KIND)
	chmod u+x $(KIND)

$(SKAFFOLD): ; $(info $(bullet) Installing <skaffold>)
	curl -sSfL https://storage.googleapis.com/skaffold/releases/v$(SKAFFOLD_VERSION)/skaffold-$(OS)-amd64 -o $(SKAFFOLD)
	chmod u+x $(SKAFFOLD)

clean: clean-kind clean-bin

clean-bin: ; $(info $(bullet) Cleaning <bin>)
	rm -rf bin/

clean-kind: $(KIND) ; $(info $(bullet) Cleaning <kind>)
	$(KIND) delete cluster --name $(KIND_CLUSTER_NAME) || exit 0

generate: generate-sqlc generate-openapi

generate-sqlc: $(SQLC) ; $(info $(bullet) Generating <sqlc>)
	$(SQLC) generate

generate-openapi: $(SQLC) ; $(info $(bullet) Generating <openapi>)
	$(OPENAPIGENERATORCLI) generate \
		--input-spec api/todo.yaml \
		--output pkg/api/todo \
		--generator-name go-experimental \
		--package-name=todo \
		--additional-properties withGoCodegenComment \
		--import-mappings=uuid.UUID=github.com/google/uuid --type-mappings=UUID=uuid.UUID \

lint: $(GOLANGCILINT) ; $(info $(bullet) Running linter)
	$(GOLANGCILINT) run ./...

test: ; $(info $(bullet) Running tests)
	go test ./...

test-coverage: ; $(info $(bullet) Running tests with coverage)
	go test -cover ./...

bootstrap-kind: $(KIND); $(info $(bullet) Bootstrap <kind>)
	$(KIND) create cluster \
		--name $(KIND_CLUSTER_NAME) \
		--image kindest/node:v$(KUBERNETES_VERSION) \
		--wait 1m

bootstrap-deploy: $(SKAFFOLD); $(info $(bullet) Bootstrap <deploy>)
	$(SKAFFOLD) build -q | $(SKAFFOLD) deploy -p bootstrap --kube-context=$(KUBERNETES_CONTEXT) --build-artifacts -

bootstrap: bootstrap-kind bootstrap-deploy
