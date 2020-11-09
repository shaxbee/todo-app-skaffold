include makefiles/shared.mk

include makefiles/go.mk
include makefiles/sqlc.mk
include makefiles/openapi.mk
include makefiles/kind.mk
include makefiles/skaffold.mk

.PHONY: all deploy run debug build-goose build-todo-service bootstrap-deployment

build: build-goose build-todo-service

build-goose: ## Build goose
	$(info $(_bullet) Building <goose>)
	$(GO) build -o bin/goose ./cmd/goose

build-todo-service: ## Build todo-service
	$(info $(_bullet) Building <todo-service>) 
	$(GO) build -o bin/todo-service ./services/todo

bootstrap: bootstrap-deployment

bootstrap-deployment: ## Bootstrap deployment
	$(info $(_bullet) Bootstrap <deployment>)
	kubectl apply --context $(BOOTSTRAP_CONTEXT) -k ops/bootstrap/overlays/dev

deploy: deploy-skaffold

run: run-skaffold ## Run

dev: dev-skaffold ## Run in development mode

debug: debug-skaffold ## Run in debug mode