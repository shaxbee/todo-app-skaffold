ifndef _include_kind_mk
_include_kind_mk := 1

include makefiles/shared.mk

KIND := bin/kind
KIND_VERSION ?= 0.9.0
KIND_CLUSTER_NAME ?= local
KIND_KUBERNETES_VERSION ?= 1.17.11
KIND_HOST_PORT ?= 80

BOOTSTRAP_CONTEXT := kind-$(KIND_CLUSTER_NAME)

$(KIND):
	$(info $(_bullet) Installing <kind>)
	@mkdir -p bin
	curl -sSfL https://kind.sigs.k8s.io/dl/v$(KIND_VERSION)/kind-$(OS)-amd64 -o $(KIND)
	chmod u+x $(KIND)

clean: clean-kind

clean-bin: clean-kind

bootstrap: bootstrap-kind

.PHONY: clean-kind bootstrap-kind

clean-kind bootstrap-kind: export KIND := $(KIND) 
clean-kind bootstrap-kind: export KIND_CLUSTER_NAME := $(KIND_CLUSTER_NAME)
clean-kind bootstrap-kind: export KIND_KUBERNETES_VERSION := $(KIND_KUBERNETES_VERSION)
clean-kind bootstrap-kind: export KIND_HOST_PORT := $(KIND_HOST_PORT)

clean-kind: $(KIND) # Delete cluster
	$(info $(_bullet) Cleaning <kind>)
	scripts/clean-kind

bootstrap-kind: $(KIND) ## Bootstrap cluster
	$(info $(_bullet) Bootstraping <kind>)
	scripts/bootstrap-kind

endif