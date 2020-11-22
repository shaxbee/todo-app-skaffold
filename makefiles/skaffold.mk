ifndef include_skaffold_mk
_include_skaffold_mk := 1

include makefiles/shared.mk
include makefiles/kubectl.mk

SKAFFOLD := bin/skaffold
SKAFFOLD_VERSION ?= 1.16.0

$(SKAFFOLD):
	$(info $(_bullet) Installing <skaffold>)
	@mkdir -p bin
	curl -sSfL https://storage.googleapis.com/skaffold/releases/v$(SKAFFOLD_VERSION)/skaffold-$(OS)-amd64 -o $(SKAFFOLD)
	chmod u+x $(SKAFFOLD)

deploy: deploy-skaffold

.PHONY: clean-skaffold build-skaffold deploy-skaffold run-skaffold dev-skaffold debug-skaffold

clean-skaffold build-skaffold deploy-skaffold run-skaffold dev-skaffold debug-skaffold: $(SKAFFOLD) $(KUBECTL)
clean-skaffold build-skaffold deploy-skaffold run-skaffold dev-skaffold debug-skaffold: export PATH := $(shell pwd)/bin:$(PATH)

clean-skaffold: ## Clean Skaffold
	$(info $(_bullet) Cleaning <skaffold>)
	! kubectl config current-context &>/dev/null || \
	$(SKAFFOLD) delete

build-skaffold: ## Build artifacts with Skaffold
	$(info $(_bullet) Building artifacts with <skaffold>)
	$(SKAFFOLD) build

deploy-skaffold: build-skaffold ## Deploy artifacts with Skaffold
	$(info $(_bullet) Deploying with <skaffold>)
	$(SKAFFOLD) build -q | $(SKAFFOLD) deploy --force --build-artifacts -

run-skaffold: ## Run with Skaffold
	$(info $(_bullet) Running stack with <skaffold>)
	$(SKAFFOLD) run --force

dev-skaffold: ## Run in development mode with Skaffold
	$(info $(_bullet) Running stack in development mode with <skaffold>)
	$(SKAFFOLD) dev --force --port-forward

debug-skaffold: ## Run in debugging mode with Skaffold
	$(info $(_bullet) Running stack in debugging mode with <skaffold>)
	$(SKAFFOLD) debug --force --port-forward

endif
