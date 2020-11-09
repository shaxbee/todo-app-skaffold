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

.PHONY: clean clean-kind bootstrap-kind

clean: clean-kind

clean-bin: clean-kind

clean-kind: $(KIND) # Delete cluster
	$(info $(_bullet) Cleaning <kind>)
	$(KIND) delete cluster --name $(KIND_CLUSTER_NAME)
	docker rm --force kind-registry &>/dev/null || exit 0

bootstrap: bootstrap-kind

bootstrap-kind: $(KIND) ## Bootstrap cluster
	$(info $(_bullet) Bootstraping <kind>)
	$(env | grep KIND) scripts/bootstrap-kind

endif