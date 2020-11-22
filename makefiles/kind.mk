ifndef _include_kind_mk
_include_kind_mk := 1
_kind_mk_path := $(dir $(lastword $(MAKEFILE_LIST)))

include makefiles/shared.mk
include makefiles/kubectl.mk

KIND := bin/kind
KIND_VERSION ?= 0.9.0
KIND_CLUSTER_NAME ?= local
KIND_K8S_VERSION ?= 1.17.11
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

clean-kind bootstrap-kind: export PATH := bin:$(PATH)
clean-kind bootstrap-kind: export CLUSTER_NAME := $(KIND_CLUSTER_NAME)
clean-kind bootstrap-kind: export K8S_VERSION := $(KIND_K8S_VERSION)
clean-kind bootstrap-kind: export HOST_PORT := $(KIND_HOST_PORT)

clean-kind: $(KIND) # Delete cluster
	$(info $(_bullet) Cleaning <kind>)
	$(dir $(_kind_mk_path))scripts/clean-kind

bootstrap-kind: $(KUBECTL) $(KIND)
	$(info $(_bullet) Bootstraping <kind>)
	$(dir $(_kind_mk_path))scripts/bootstrap-kind

endif
