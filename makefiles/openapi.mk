ifndef _include_openapi_mk
_include_openapi_mk := 1
_openapi_mk_path := $(lastword $(MAKEFILE_LIST))

include makefiles/shared.mk

OPENAPI_SPEC ?= api/api.yaml
OPENAPI_OUTPUT ?= pkg/api
OPENAPI_PACKAGE_NAME ?= api

OPENAPIGENERATORCLI := $(dir $(_openapi_mk_path))scripts/openapi-generator-cli
OPENAPIGENERATORCLI_VERSION ?= 4.3.1

lint: lint-openapi

generate: generate-openapi

.PHONY: lint-openapi generate-openapi

lint-openapi generate-openapi: export OPENAPIGENERATORCLI_VERSION := $(OPENAPIGENERATORCLI_VERSION)

lint-openapi: ## List OpenAPI spec
	$(info $(_bullet) Linting <openapi>)
	$(OPENAPIGENERATORCLI) validate --input-spec $(OPENAPI_SPEC)

generate-openapi: ## Generate OpenAPI code
	$(info $(_bullet) Generating <openapi>)
	$(OPENAPIGENERATORCLI) generate \
		--input-spec $(OPENAPI_SPEC) \
		--output $(OPENAPI_OUTPUT) \
		--generator-name go-experimental \
		--package-name=$(OPENAPI_PACKAGE_NAME) \
		--additional-properties withGoCodegenComment \
		--import-mappings=uuid.UUID=github.com/google/uuid --type-mappings=UUID=uuid.UUID

endif