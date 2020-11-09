ifndef _include_kind_mk
_include_kind_mk := 1

include makefiles/shared.mk

KIND := bin/kind
KIND_VERSION ?= 0.9.0
KIND_CLUSTER_NAME ?= local

KUBERNETES_VERSION ?= 1.17.11
BOOTSTRAP_CONTEXT := kind-$(KIND_CLUSTER_NAME)

$(KIND):
	$(info $(_bullet) Installing <kind>)
	@mkdir -p bin
	curl -sSfL https://kind.sigs.k8s.io/dl/v$(KIND_VERSION)/kind-$(OS)-amd64 -o $(KIND)
	chmod u+x $(KIND)

.PHONY: clean clean-kind bootstrap-kind

clean: clean-kind

clean-bin: clean-kind

clean-kind: # Delete cluster
	$(info $(_bullet) Cleaning <kind>)
	! test -f $(KIND) || \
	$(KIND) delete cluster --name $(KIND_CLUSTER_NAME)
	! docker top kind-registry &>/dev/null || \
	docker rm --force kind-registry

bootstrap: bootstrap-kind

bootstrap-kind: $(KIND) ## Bootstrap Kubernetes in Docker
	$(info $(_bullet) Bootstraping <kind>)
	docker top kind-registry &>/dev/null || \
	docker run \
		--detach \
		--restart always \
		--publish 5000:5000 \
		--name kind-registry \
		registry:2
	$(KIND) get clusters | grep -q $(KIND_CLUSTER_NAME) || \
	$(KIND) create cluster \
		--name $(KIND_CLUSTER_NAME) \
		--config kind.yaml \
		--image kindest/node:v$(KUBERNETES_VERSION) \
		--wait 1m
	kubectl apply --context $(BOOTSTRAP_CONTEXT) -k ops/bootstrap/overlays/dev
	docker network connect kind kind-registry &>/dev/null || true

endif