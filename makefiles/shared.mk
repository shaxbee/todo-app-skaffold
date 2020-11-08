ifndef _include_shared_mk
_include_shared_mk := 1

OS ?= $(shell uname -s | tr [:upper:] [:lower:])

_bullet := $(shell printf "\033[34;1m▶\033[0m")

.PHONY: help clean clean-bin build generate deploy git-dirty git-hooks

all: deps generate format lint test build

help: ## Help
	@cat $(sort $(MAKEFILE_LIST)) | grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' | sort

clean: clean-bin ## Clean targets

deps: ## Download dependencies

vendor: ## Vendor dependencies

generate: ## Generate code

format: ## Format code

lint: ## Lint code

test: ## Run tests

test-coverage: ## Run tests with coverage

integration-test: ## Run integration tests

build: ## Build all targets

bootstrap: ## Bootstrap

deploy: ## Deploy

clean-bin: ## Clean installed tools
	$(info $(_bullet) Cleaning <bin>)
	rm -rf bin/

git-dirty: ## Check for uncommited changes
	$(info $(_bullet) Checking for uncommited changes)
	git status --porcelain
	git diff --quiet --exit-code

git-hooks: ## Configure git hooks
	$(info $(_bullet) Configuring git hooks)
	git config core.hooksPath .githooks

endif