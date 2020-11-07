ifndef _include_openapi_mk
_include_openapi_mk := 1

include makefiles/shared.mk

OPENAPI_SPEC ?= api/api.yaml
OPENAPI_OUTPUT ?= pkg/api
OPENAPI_PACKAGE_NAME ?= api

OPENAPIGENERATORCLI := scripts/openapi-generator-cli
OPENAPIGENERATORCLI_VERSION ?= 4.3.1

.PHONY: generate generate-openapi

generate: generate-openapi

generate-openapi: ## Generate OpenAPI code
	$(info $(_bullet) Generating <openapi>)
	OPENAPIGENERATORCLI_VERSION=$(OPENAPIGENERATORCLI_VERSION) $(OPENAPIGENERATORCLI) generate \
		--input-spec $(OPENAPI_SPEC) \
		--output $(OPENAPI_OUTPUT) \
		--generator-name go-experimental \
		--package-name=$(OPENAPI_PACKAGE_NAME) \
		--additional-properties withGoCodegenComment \
		--import-mappings=uuid.UUID=github.com/google/uuid --type-mappings=UUID=uuid.UUID

endif