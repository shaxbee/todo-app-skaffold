include makefiles/shared.mk

include makefiles/go.mk
include makefiles/sqlc.mk
include makefiles/openapi.mk
include makefiles/docker.mk
include makefiles/kind.mk
include makefiles/skaffold.mk

build: build-goose build-todo-service

bootstrap: bootstrap-deployment

deploy: deploy-skaffold

run: run-skaffold

dev: dev-skaffold

debug: debug-skaffold

.PHONY: build-goose build-todo-service bootstrap-deployment

build-goose: ## Build goose
	$(info $(_bullet) Building <goose>)
	$(GO) build -o bin/goose ./cmd/goose

build-todo-service: ## Build todo-service
	$(info $(_bullet) Building <todo-service>) 
	$(GO) build -o bin/todo-service ./services/todo

bootstrap-deployment: ## Bootstrap deployment
	$(info $(_bullet) Bootstraping <deployment>)
	PATH=bin:$(PATH) kubectl apply --context $(BOOTSTRAP_CONTEXT) -k ops/bootstrap/overlays/local
