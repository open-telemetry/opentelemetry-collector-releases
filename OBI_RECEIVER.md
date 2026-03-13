# OBI Receiver Integration

This document explains how OBI (eBPF-based zero-code instrumentation) is integrated into OpenTelemetry Collector distributions, why the integration is structured this way, and how to build and maintain it.

## Overview

The OBI receiver is included in the `otelcol-contrib` distribution as an external component. It is maintained in its own repository at [opentelemetry-ebpf-instrumentation](https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation).

To avoid git submodules, this repository downloads the OBI release tarball (`obi-vX.Y.Z-source-generated.tar.gz`) during builds. This tarball is published alongside each OBI release and already contains all pre-generated BPF source files.

## Architecture

Build preparation does the following:

1. Reads the OBI version from `distributions/otelcol-contrib/manifest.yaml`
2. Downloads `obi-vX.Y.Z-source-generated.tar.gz` from the OBI GitHub release
3. Verifies the SHA256 checksum against the `SHA256SUMS` release asset
4. Extracts the tarball to `internal/obi-src` (untracked)
5. Creates a version-keyed stamp file (`internal/obi-src/.obi-vX.Y.Z`) so subsequent builds skip the download

No BPF toolchain (Docker, clang, bpf2go) is required. The pre-generated files are included in the `source-generated` tarball.

### Integration Pattern

The `otelcol-contrib` distribution manifest includes:

```yaml
receivers:
  - gomod: go.opentelemetry.io/obi v0.5.0
    import: go.opentelemetry.io/obi/collector

replaces:
  - go.opentelemetry.io/obi => ../../../internal/obi-src
```

This keeps the pinned module version in the manifest while using a locally prepared source tree for builds.

## Building with OBI

### Prerequisites

- Standard Go build tools (no Docker required)

### Build Commands

```bash
# Build all distributions (includes OBI preparation for otelcol-contrib)
make build

# Build only otelcol-contrib
make build DISTRIBUTIONS="otelcol-contrib"
```

`make build` and `make generate` call `scripts/prepare-obi.sh` automatically when `otelcol-contrib` is selected.

## CI/CD Integration

CI jobs that build `otelcol-contrib` should run the `.github/actions/fetch-obi` composite action **before** any `make` target that triggers source generation. The action:

1. Checks the Actions cache for the OBI version (cache key: `obi-source-vX.Y.Z`)
2. On a cache miss: downloads and verifies the tarball, extracts it, and writes the stamp file
3. On a cache hit: restores from cache (stamp file already present → `prepare-obi.sh` exits immediately)

This means no Docker or network access is needed during the actual `make` build step.

```yaml
- name: Fetch OBI source
  if: inputs.distribution == 'otelcol-contrib'
  uses: ./.github/actions/fetch-obi
```

## Maintenance

### Updating the OBI Version

1. Update `distributions/otelcol-contrib/manifest.yaml`:
   - `gomod: go.opentelemetry.io/obi vX.Y.Z`
2. Build contrib to verify:
   ```bash
   make build DISTRIBUTIONS="otelcol-contrib"
   ```
3. Commit both the manifest change and any `go.sum` updates.

No repository submodule updates are required. The CI cache is invalidated automatically because the cache key is version-keyed.

## Platform Support

- Linux amd64/arm64: full OBI support (eBPF receiver active)
- All other platforms (Linux ppc64le/s390x/riscv64, Windows, macOS): OBI's non-eBPF stubs are compiled; the receiver registers but does not activate eBPF collection

The `go.opentelemetry.io/obi/collector` package uses build tags (`//go:build linux && (amd64 || arm64)`) to select the appropriate implementation. No cross-compilation guards are needed in this repository.

## References

- OBI repository: https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation
- OBI collector example: https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation/tree/main/examples/otel-collector
- Collector builder docs: https://opentelemetry.io/docs/collector/custom-collector/

