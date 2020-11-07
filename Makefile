include makefiles/shared.mk

include makefiles/go.mk
include makefiles/sqlc.mk
include makefiles/openapi.mk
include makefiles/kind.mk
include makefiles/skaffold.mk

.PHONY: all help clean generate build-goose build-todo-service build bootstrap

all: generate build

build: build-goose build-todo-service

build-goose: ## Build goose
	$(info $(_bullet) Building <goose>)
	$(GO) build -o bin/goose ./cmd/goose

build-todo-service: ## Build todo-service
	$(info $(_bullet) Building <todo-service>) 
	$(GO) build -o bin/todo-service ./services/todo