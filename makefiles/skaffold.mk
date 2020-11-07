ifndef include_skaffold_mk
_include_skaffold_mk := 1

include makefiles/shared.mk

SKAFFOLD := bin/skaffold
SKAFFOLD_VERSION ?= 1.16.0

$(SKAFFOLD):
	$(info $(_bullet) Installing <skaffold>)
	@mkdir -p bin
	curl -sSfL https://storage.googleapis.com/skaffold/releases/v$(SKAFFOLD_VERSION)/skaffold-$(OS)-amd64 -o $(SKAFFOLD)
	chmod u+x $(SKAFFOLD)

.PHONY: clean-skaffold build-skaffold deploy-skaffold dev-skaffold

clean-skaffold: $(SKAFFOLD) ## Clean Skaffold stack
	$(info $(_bullet) Cleaning <skaffold>)
	! kubectl config current-context &>/dev/null || \
	$(SKAFFOLD) delete

build-skaffold: $(SKAFFOLD) ## Build artifacts with Skaffold
	$(info $(_bullet) Building artifacts with <skaffold>)
	$(SKAFFOLD) build

deploy: deploy-skaffold

deploy-skaffold: $(SKAFFOLD) build-skaffold ## Deploy artifacts with Skaffold
	$(info $(_bullet) Deploying with <skaffold>)
	$(SKAFFOLD) build -q | $(SKAFFOLD) deploy --force --build-artifacts -

run-skaffold: $(SKAFFOLD)
	$(info $(_bullet) Run stack with <skaffold>)	
	$(SKAFFOLD) run --force --port-forward --tail

dev-skaffold: $(SKAFFOLD) ## Run stack in development mode with Skaffold
	$(info $(_bullet) Run stack in development mode with <skaffold>)
	$(SKAFFOLD) dev --force --port-forward

endif
