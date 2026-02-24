# OBI Receiver Integration

This document explains how OBI (eBPF-based zero-code instrumentation) is integrated into OpenTelemetry Collector distributions, why the integration is structured this way, and how to build and maintain it.

## Overview

The OBI receiver is included in the `otelcol-contrib` distribution as an external component. It is maintained in its own repository at [opentelemetry-ebpf-instrumentation](https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation).

To avoid git submodules, this repository prepares OBI from a tagged upstream source archive during builds.

## Architecture

OBI requires eBPF code generation (converting C code to Go) before Linux builds can compile OBI's eBPF-enabled paths. Generated files are not present in the upstream tag, so build preparation does the following:

1. Reads the OBI version from `distributions/otelcol-contrib/manifest.yaml`
2. Downloads the matching tagged source archive into `internal/obi-src` (untracked)
3. Runs `make docker-generate` in `internal/obi-src` when generated files are missing on Linux

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

- Docker (required on Linux when generated OBI artifacts are missing)
- Standard Go build tools

### Build Commands

```bash
# Build all distributions (includes OBI preparation for otelcol-contrib)
make build

# Build only otelcol-contrib
make build DISTRIBUTIONS="otelcol-contrib"
```

`make build` and `make generate` call `scripts/prepare-obi.sh` automatically when `otelcol-contrib` is selected.

## CI/CD Integration

No submodule checkout is required.

Linux CI jobs that build `otelcol-contrib` must still have Docker available for first-time OBI artifact generation. Caching `internal/obi-src` can reduce repeated generation costs.

## Maintenance

### Updating the OBI Version

1. Update `distributions/otelcol-contrib/manifest.yaml`:
   - `gomod: go.opentelemetry.io/obi vX.Y.Z`
2. Build contrib:
   ```bash
   make build DISTRIBUTIONS="otelcol-contrib"
   ```
3. Verify OBI preparation and generation complete successfully.

No repository submodule updates are required.

## Platform Support

- Linux: full OBI support (requires generated eBPF artifacts)
- Non-Linux: OBI's non-Linux code paths are used; Linux eBPF generation is skipped

## References

- OBI repository: https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation
- OBI collector example: https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation/tree/main/examples/otel-collector
- Collector builder docs: https://opentelemetry.io/docs/collector/custom-collector/
