ifndef include_skaffold_mk
_include_skaffold_mk := 1

include makefiles/shared.mk
include makefiles/kubectl.mk

SKAFFOLD := bin/skaffold
SKAFFOLD_VERSION ?= 1.16.0

$(SKAFFOLD): $(KUBECTL)
	$(info $(_bullet) Installing <skaffold>)
	@mkdir -p bin
	curl -sSfL https://storage.googleapis.com/skaffold/releases/v$(SKAFFOLD_VERSION)/skaffold-$(OS)-amd64 -o $(SKAFFOLD)
	chmod u+x $(SKAFFOLD)

deploy: deploy-skaffold

.PHONY: clean-skaffold build-skaffold deploy-skaffold run-skaffold deploy-skaffold dev-skaffold

clean-skaffold: $(SKAFFOLD) ## Clean Skaffold
	$(info $(_bullet) Cleaning <skaffold>)
	! kubectl config current-context &>/dev/null || \
	PATH=bin:$(PATH) $(SKAFFOLD) delete

build-skaffold: $(SKAFFOLD) ## Build artifacts with Skaffold
	$(info $(_bullet) Building artifacts with <skaffold>)
	PATH=bin:$(PATH) $(SKAFFOLD) build

deploy-skaffold: $(SKAFFOLD) build-skaffold ## Deploy artifacts with Skaffold
	$(info $(_bullet) Deploying with <skaffold>)
	$(SKAFFOLD) build -q | $(SKAFFOLD) deploy --force --build-artifacts -

run-skaffold: $(SKAFFOLD) ## Run with Skaffold
	$(info $(_bullet) Running stack with <skaffold>)
	$(SKAFFOLD) run --force

dev-skaffold: $(SKAFFOLD) ## Run in development mode with Skaffold
	$(info $(_bullet) Running stack in development mode with <skaffold>)
	$(SKAFFOLD) dev --force --port-forward

debug-skaffold: $(SKAFFOLD) ## Run in debugging mode with Skaffold
	$(info $(_bullet) Running stack in debugging mode with <skaffold>)
	$(SKAFFOLD) debug --force --port-forward

clean-skaffold build-skaffold deploy-skaffold run-skaffold dev-skaffold debug-skaffold: export PATH := $(shell pwd)/bin:$(PATH)

endif
