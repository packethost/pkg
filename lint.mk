
# BEGIN: lint-install -makefile lint.mk .
# http://github.com/tinkerbell/lint-install

GOLINT_VERSION ?= v1.42.1


YAMLLINT_VERSION ?= 1.26.3
LINT_OS := $(shell uname)
LINT_ARCH := $(shell uname -m)
LINT_ROOT := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# shellcheck and hadolint lack arm64 native binaries: rely on x86-64 emulation
ifeq ($(LINT_OS),Darwin)
	ifeq ($(LINT_ARCH),arm64)
		LINT_ARCH=x86_64
	endif
endif


GOLINT_CONFIG = $(LINT_ROOT)/.golangci.yml
YAMLLINT_ROOT = out/linters/yamllint-$(YAMLLINT_VERSION)

.PHONY: lint
lint: out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH) $(YAMLLINT_ROOT)/dist/bin/yamllint
	find . -name go.mod -execdir "$(LINT_ROOT)/out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH)" run -c "$(GOLINT_CONFIG)" \;
	PYTHONPATH=$(YAMLLINT_ROOT)/dist $(YAMLLINT_ROOT)/dist/bin/yamllint .

.PHONY: fix
fix: out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH)
	find . -name go.mod -execdir "$(LINT_ROOT)/out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH)" run -c "$(GOLINT_CONFIG)" --fix \;

out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH):
	mkdir -p out/linters
	rm -rf out/linters/golangci-lint-*
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b out/linters $(GOLINT_VERSION)
	mv out/linters/golangci-lint out/linters/golangci-lint-$(GOLINT_VERSION)-$(LINT_ARCH)

$(YAMLLINT_ROOT)/dist/bin/yamllint:
	mkdir -p out/linters
	rm -rf out/linters/yamllint-*
	curl -sSfL https://github.com/adrienverge/yamllint/archive/refs/tags/v$(YAMLLINT_VERSION).tar.gz | tar -C out/linters -zxf -
	cd $(YAMLLINT_ROOT) && pip3 install . -t dist
# END: lint-install -makefile lint.mk .
