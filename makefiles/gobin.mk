ifndef _include_gobin_mk
_include_gobin_mk := 1

include makefiles/shared.mk

GOBIN := bin/gobin
GOBIN_VERSION := 0.0.14

$(GOBIN):
	$(info $(_bullet) Installing <gobin>)
	@mkdir -p bin
	curl -sSfL https://github.com/myitcv/gobin/releases/download/v$(GOBIN_VERSION)/$(OS)-amd64 -o $(GOBIN)
	chmod u+x $(GOBIN)

endif