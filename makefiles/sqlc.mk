ifndef include_sqlc_mk
_include_sqlc_mk := 1

include makefiles/shared.mk

SQLC := bin/sqlc
SQLC_VERSION ?= 1.5.0

$(SQLC):
	$(info $(_bullet) Installing <sqlc>)
	@mkdir -p bin
	curl -sSfL https://github.com/kyleconroy/sqlc/releases/download/v$(SQLC_VERSION)/sqlc-v$(SQLC_VERSION)-$(OS)-amd64.tar.gz | tar -C bin -xz
	chmod u+x $(SQLC)

.PHONY: generate generate-sqlc

generate: generate-sqlc

generate-sqlc: $(SQLC) ## Generate SQLC code
	$(info $(_bullet) Generating <sqlc>)
	$(SQLC) generate

endif