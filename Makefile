GO ?= go
GORELEASER ?= goreleaser

# SRC_ROOT is the top of the source tree.
SRC_ROOT := $(shell git rev-parse --show-toplevel)

# renovate: datasource=github-releases depName=OCB packageName=open-telemetry/opentelemetry-collector
OTELCOL_BUILDER_VERSION ?= 0.139.0

OTELCOL_BUILDER_DIR ?= ${HOME}/bin
OTELCOL_BUILDER ?= ${OTELCOL_BUILDER_DIR}/ocb

GOCMD?= go
TOOLS_MOD_DIR   := $(SRC_ROOT)/internal/tools
TOOLS_BIN_DIR   := $(SRC_ROOT)/.tools
TOOLS_MOD_REGEX := "\s+_\s+\".*\""
TOOLS_PKG_NAMES := $(shell grep -E $(TOOLS_MOD_REGEX) < $(TOOLS_MOD_DIR)/tools.go | tr -d " _\"" | grep -vE '/v[0-9]+$$')
TOOLS_BIN_NAMES := $(addprefix $(TOOLS_BIN_DIR)/, $(notdir $(shell echo $(TOOLS_PKG_NAMES))))
CHLOGGEN        := $(TOOLS_BIN_DIR)/chloggen
CHLOGGEN_CONFIG := .chloggen/config.yaml

DISTRIBUTIONS ?= "otelcol,otelcol-contrib,otelcol-k8s,otelcol-otlp,otelcol-ebpf-profiler"
BINARIES ?= "builder,opampsupervisor"

ci: check build
check: ensure-goreleaser-up-to-date validate-components

build: go ocb
	@./scripts/build.sh -d "${DISTRIBUTIONS}" -b ${OTELCOL_BUILDER}

generate: generate-sources generate-goreleaser

generate-goreleaser: go
	@./scripts/generate-goreleaser.sh -d "${DISTRIBUTIONS}" -b "${BINARIES}" -g ${GO}

generate-sources: go ocb generate-msi
	@./scripts/build.sh -d "${DISTRIBUTIONS}" -s true -b ${OTELCOL_BUILDER}

generate-msi: go ocb
	$(GO) run cmd/msi-generator/main.go -d "${DISTRIBUTIONS}"

goreleaser-verify: goreleaser
	@${GORELEASER} release --snapshot --clean

ensure-goreleaser-up-to-date: generate-goreleaser
	@git diff -s --exit-code distributions/*/.goreleaser.yaml || (echo "Check failed: The goreleaser templates have changed but the .goreleaser.yamls haven't. Run 'make generate-goreleaser' and update your PR." && exit 1)

validate-components:
	@./scripts/validate-components.sh

.PHONY: ocb
ocb:
ifeq (, $(shell command -v ocb 2>/dev/null))
	@{ \
	[ ! -x '$(OTELCOL_BUILDER)' ] || exit 0; \
	set -e ;\
	os=$$(uname | tr A-Z a-z) ;\
	machine=$$(uname -m) ;\
	[ "$${machine}" != x86 ] || machine=386 ;\
	[ "$${machine}" != x86_64 ] || machine=amd64 ;\
	echo "Installing ocb ($${os}/$${machine}) at $(OTELCOL_BUILDER_DIR)";\
	mkdir -p $(OTELCOL_BUILDER_DIR) ;\
	CGO_ENABLED=0 go install -trimpath -ldflags="-s -w" go.opentelemetry.io/collector/cmd/builder@v$(OTELCOL_BUILDER_VERSION) ;\
	mv $$(go env GOPATH)/bin/builder $(OTELCOL_BUILDER) ;\
	}
else
OTELCOL_BUILDER=$(shell command -v ocb)
endif

.PHONY: go
go:
	@{ \
		if ! command -v '$(GO)' >/dev/null 2>/dev/null; then \
			echo >&2 '$(GO) command not found. Please install golang. https://go.dev/doc/install'; \
			exit 1; \
		fi \
	}

.PHONY: goreleaser
goreleaser:
	@{ \
		if ! command -v '$(GORELEASER)' >/dev/null 2>/dev/null; then \
			echo >&2 '$(GORELEASER) command not found. Please install goreleaser. https://goreleaser.com/install/'; \
			exit 1; \
		fi \
	}

REMOTE?=git@github.com:open-telemetry/opentelemetry-collector-releases.git
.PHONY: push-tags
push-tags:
	@[ "${TAG}" ] || ( echo ">> env var TAG is not set"; exit 1 )
	@echo "Adding tag ${TAG}"
	@git tag -a ${TAG} -s -m "Version ${TAG}"
	@echo "Pushing tag ${TAG}"
	@git push ${REMOTE} ${TAG}

# Used for debug only
REMOTE?=git@github.com:open-telemetry/opentelemetry-collector-releases.git
.PHONY: delete-tags
delete-tags:
	@[ "${TAG}" ] || ( echo ">> env var TAG is not set"; exit 1 )
	@echo "Deleting local tag ${TAG}"
	@if [ -n "$$(git tag -l ${TAG})" ]; then \
		git tag -d ${TAG}; \
	fi
	@if [ -n "$$(git tag -l cmd/builder/${TAG})" ]; then \
		git tag -d cmd/builder/${TAG}; \
	fi
	@echo "Deleting remote tag ${TAG}"
	@git push ${REMOTE} :refs/tags/${TAG}
	@git push ${REMOTE} :refs/tags/cmd/builder/${TAG}

# Used for debug only
REMOTE?=git@github.com:open-telemetry/opentelemetry-collector-releases.git
.PHONY: repeat-tags
repeat-tags: delete-tags push-tags

.PHONY: install-tools
install-tools: $(TOOLS_BIN_NAMES)

$(TOOLS_BIN_DIR):
	mkdir -p $@

$(TOOLS_BIN_NAMES): $(TOOLS_BIN_DIR) $(TOOLS_MOD_DIR)/go.mod
	cd $(TOOLS_MOD_DIR) && $(GOCMD) build -o $@ -trimpath $(filter %/$(notdir $@),$(TOOLS_PKG_NAMES))

FILENAME?=$(shell git branch --show-current)
.PHONY: chlog-new
chlog-new: $(CHLOGGEN)
	$(CHLOGGEN) new --config $(CHLOGGEN_CONFIG) --filename $(FILENAME)

.PHONY: chlog-validate
chlog-validate: $(CHLOGGEN)
	$(CHLOGGEN) validate --config $(CHLOGGEN_CONFIG)

.PHONY: chlog-preview
chlog-preview: $(CHLOGGEN)
	$(CHLOGGEN) update --config $(CHLOGGEN_CONFIG) --dry

.PHONY: chlog-update
chlog-update: $(CHLOGGEN)
	$(CHLOGGEN) update --config $(CHLOGGEN_CONFIG) --version $(VERSION)
