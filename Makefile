GO=$(shell which go)
OTELCOL_BUILDER_VERSION ?= 0.35.0
OTELCOL_BUILDER_DIR ?= ~/bin
OTELCOL_BUILDER ?= ${OTELCOL_BUILDER_DIR}/opentelemetry-collector-builder

YQ_VERSION ?= 4.11.1
YQ_DIR ?= ${OTELCOL_BUILDER_DIR}
YQ ?= ${YQ_DIR}/yq

DISTRIBUTIONS ?= "otelcol"

ci: check build
check: ensure-goreleaser-up-to-date

build: otelcol-builder
	@./scripts/build.sh -d "${DISTRIBUTIONS}" -b ${OTELCOL_BUILDER} -g ${GO}

generate: generate-sources generate-goreleaser

generate-goreleaser: yq
	@./scripts/generate-goreleaser-config.sh -d "${DISTRIBUTIONS}" -y "${YQ}"

generate-sources: otelcol-builder
	@./scripts/build.sh -d "${DISTRIBUTIONS}" -s true -b ${OTELCOL_BUILDER} -g ${GO}

goreleaser-verify:
	@goreleaser release --snapshot --rm-dist

ensure-goreleaser-up-to-date: generate-goreleaser
	@git diff -s --exit-code .goreleaser.yaml || (echo "Build failed: The goreleaser templates have changed but the .goreleaser.yaml hasn't. Run 'make generate-goreleaser' and update your PR." && exit 1)

otelcol-builder:
ifeq (, $(shell which opentelemetry-collector-builder))
	@{ \
	set -e ;\
	echo Installing opentelemetry-collector-builder at $(OTELCOL_BUILDER_DIR);\
	mkdir -p $(OTELCOL_BUILDER_DIR) ;\
	curl -sLo $(OTELCOL_BUILDER) https://github.com/open-telemetry/opentelemetry-collector-builder/releases/download/v$(OTELCOL_BUILDER_VERSION)/opentelemetry-collector-builder_$(OTELCOL_BUILDER_VERSION)_linux_amd64 ;\
	chmod +x $(OTELCOL_BUILDER) ;\
	}
else
OTELCOL_BUILDER=$(shell which opentelemetry-collector-builder)
endif

yq:
ifeq (, $(shell which yq))
	@{ \
	set -e ;\
	echo Installing yq at $(YQ_DIR);\
	mkdir -p $(YQ_DIR) ;\
	curl -sLo $(YQ) https://github.com/mikefarah/yq/releases/download/v$(YQ_VERSION)/yq_linux_amd64 ;\
	chmod +x $(YQ) ;\
	}
else
YQ=$(shell which yq)
endif
