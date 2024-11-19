GO ?= go
GORELEASER ?= goreleaser

OTELCOL_BUILDER_VERSION ?= 0.113.0
OTELCOL_BUILDER_DIR ?= ${HOME}/bin
OTELCOL_BUILDER ?= ${OTELCOL_BUILDER_DIR}/ocb

DISTRIBUTIONS ?= "axoflow-otel-collector"

ci: check build
check: ensure-goreleaser-up-to-date

build: go ocb
	@./scripts/build.sh -d "${DISTRIBUTIONS}" -b ${OTELCOL_BUILDER}

generate: generate-sources generate-goreleaser

generate-goreleaser: go
	chmod +x ./scripts/generate-goreleaser.sh
	@./scripts/generate-goreleaser.sh -d "${DISTRIBUTIONS}" -g ${GO}

generate-sources: go ocb
	@./scripts/build.sh -d "${DISTRIBUTIONS}" -s true -b ${OTELCOL_BUILDER}

goreleaser-verify: goreleaser
	@${GORELEASER} release --snapshot --clean

ensure-goreleaser-up-to-date: generate-goreleaser
	@git diff -s --exit-code distributions/*/.goreleaser.yaml || (echo "Check failed: The goreleaser templates have changed but the .goreleaser.yamls haven't. Run 'make generate-goreleaser' and update your PR." && exit 1)

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

REMOTE?=git@github.com:axoflow/axoflow-otel-collector-releases.git
.PHONY: push-tags
push-tags:
	@[ "${TAG}" ] || ( echo ">> env var TAG is not set"; exit 1 )
	@echo "Adding tag ${TAG}"
	@git tag -a ${TAG} -s -m "Version ${TAG}"
	@echo "Pushing tag ${TAG}"
	@git push ${REMOTE} ${TAG}
