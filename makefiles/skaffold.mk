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

.PHONY: clean-skaffold build-skaffold deploy-skaffold

clean: clean-skaffold

clean-bin: clean-skaffold

clean-skaffold: $(SKAFFOLD)
	$(info $(_bullet) Cleaning <skaffold>)
	! kubectl config current-context 2>/dev/null || \
	$(SKAFFOLD) delete

build: build-skaffold

build-skaffold: $(SKAFFOLD) ## Build artifacts with skaffold
	$(info $(_bullet) Building artifacts with <skaffold>)
	$(SKAFFOLD) build

deploy: deploy-skaffold

deploy-skaffold: $(KUBECTL) $(SKAFFOLD) build-skaffold ## Deploy artifacts with skaffold
	$(info $(_bullet) Deploying with <skaffold>)
	$(SKAFFOLD) build -q | $(SKAFFOLD) deploy --force --build-artifacts -

endif
