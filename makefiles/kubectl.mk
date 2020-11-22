ifndef _include_kubectl_mk
_include_kubectl_mk := 1

include makefiles/shared.mk

KUBECTL := bin/kubectl
KUBECTL_VERSION ?= 1.17.11

$(KUBECTL):
	$(info $(_bullet) Installing <kubectl>)
	@mkdir -p bin
	curl -sSfL https://storage.googleapis.com/kubernetes-release/release/v$(KUBECTL_VERSION)/bin/$(OS)/amd64/kubectl -o $(KUBECTL)
	chmod u+x $(KUBECTL)

endif

