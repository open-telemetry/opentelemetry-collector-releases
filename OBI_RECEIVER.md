# OBI Receiver Integration

This document explains how OBI (eBPF-based zero-code instrumentation) is integrated into the OpenTelemetry Collector distributions.

## Overview

The OBI receiver is included in the `otelcol-contrib` distribution as an external component. It is maintained in its own repository at [opentelemetry-ebpf-instrumentation](https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation) and included here as a git submodule.

## Architecture

### Why a Submodule?

OBI requires eBPF code generation (converting C code to Go) before the module can be built. The generated files are:
- **Not committed** to OBI's git repository (policy decision)
- **Platform-specific** (Linux kernel versions)
- **Large** (~10MB of binary objects and generated Go code)

Therefore, we need:
1. The OBI source code (git submodule at `internal/obi`)
2. A generation step before building (`make generate-obi`)
3. A replace directive to point the builder to the local submodule

### Integration Pattern

The `otelcol-contrib` distribution manifest includes:

```yaml
receivers:
  - gomod: go.opentelemetry.io/obi v0.5.0
    import: go.opentelemetry.io/obi/collector

replaces:
  - go.opentelemetry.io/obi => ../../internal/obi
```

This uses the OpenTelemetry Collector Builder's `import:` field to directly reference the OBI collector package without needing a wrapper module.

## Building with OBI

### Prerequisites

- Git with submodule support
- Docker (for eBPF code generation)
- Standard Go build tools

### First Time Setup

```bash
# Clone with submodules
git clone --recurse-submodules https://github.com/open-telemetry/opentelemetry-collector-releases.git

# Or initialize submodules if already cloned
git submodule update --init --recursive

# Generate OBI eBPF artifacts (requires Docker, takes 2-5 minutes)
make generate-obi
```

### Building Distributions

Once OBI artifacts are generated:

```bash
# Build all distributions
make build

# Or build specific distribution
make build DISTRIBUTIONS="otelcol-contrib"
```

The build process will automatically check that OBI artifacts exist before building.

## CI/CD Integration

### Required Workflow Changes

All CI workflows that build distributions must:

1. **Initialize submodules:**
   ```yaml
   - uses: actions/checkout@v6
     with:
       submodules: 'recursive'
   ```

2. **Generate OBI artifacts (Linux only):**
   ```yaml
   - name: Generate OBI eBPF artifacts
     if: runner.os == 'Linux'
     run: make generate-obi
   ```

3. **Build as normal:**
   ```yaml
   - name: Build distributions
     run: make build
   ```

### Build Time Impact

- **Without caching:** Adds ~2-5 minutes to Linux builds
- **With caching:** Near-zero impact after first build
- **Submodule checkout:** Adds ~10-20 seconds

### Artifact Caching

To optimize CI builds, cache generated artifacts:

```yaml
- name: Cache OBI artifacts
  id: cache-obi
  uses: actions/cache@v3
  with:
    path: |
      internal/obi/**/*_bpfel.go
      internal/obi/**/*_bpfeb.go
      internal/obi/**/*.o
    key: obi-artifacts-${{ hashFiles('internal/obi/.git/HEAD') }}

- name: Generate OBI eBPF artifacts
  if: runner.os == 'Linux' && steps.cache-obi.outputs.cache-hit != 'true'
  run: make generate-obi
```

## Maintenance

### Updating OBI Version

1. **Update the submodule to a new commit:**
   ```bash
   cd internal/obi
   git fetch
   git checkout v0.6.0  # or specific commit
   cd ../..
   git add internal/obi
   ```

2. **Update the manifest version:**
   ```bash
   # Edit distributions/otelcol-contrib/manifest.yaml
   # Change: gomod: go.opentelemetry.io/obi v0.5.0
   # To:     gomod: go.opentelemetry.io/obi v0.6.0
   ```

3. **Regenerate artifacts:**
   ```bash
   make generate-obi
   ```

4. **Test the build:**
   ```bash
   make build DISTRIBUTIONS="otelcol-contrib"
   ```

5. **Commit the changes:**
   ```bash
   git commit -m "Update OBI receiver to v0.6.0"
   ```

### Troubleshooting

#### Error: "OBI submodule not initialized"

**Solution:**
```bash
git submodule update --init --recursive
```

#### Error: "OBI eBPF artifacts not generated"

**Solution:**
```bash
make generate-obi
```

If generation fails, check:
- Docker is installed and running
- Docker can pull images
- You have sufficient disk space (generation needs ~500MB temp space)

#### Slow Builds

**Solution:** Implement CI artifact caching (see above)

## Platform Support

- ✅ **Linux (x86_64, arm64):** Full support with eBPF instrumentation
- ⚠️ **Other platforms:** Receiver included but runs as no-op (OBI uses build tags)

The OBI code itself handles platform detection:
- `collector/factory_linux.go` - Full implementation
- `collector/factory_others.go` - No-op stub

## Why Not in Collector-Contrib?

The OBI receiver is an **external component** maintained in its own repository, not a **donated component**. According to collector-contrib documentation:

> "you can just host it in your own repository as a Go module"

The collector-releases repository is the appropriate place for:
- Packaging external components into distributions
- Managing build-time dependencies (like submodules)
- Distribution-specific configurations

The collector-contrib repository is for:
- Components donated to the project
- Code maintained by the OpenTelemetry community
- Components that follow contrib's lifecycle

## Future Improvements

### Option 1: Pre-Generated Artifacts (Recommended)

If OBI publishes pre-generated artifacts as GitHub release assets:

1. Remove the submodule
2. Add a download step: `scripts/download-obi-artifacts.sh`
3. Remove the replace directive (use published module)

This would eliminate all generation complexity.

### Option 2: OBI Commits Artifacts

If OBI changes policy and commits generated files to git:

1. Remove the submodule
2. Remove generation targets from Makefile
3. Remove the replace directive
4. Reference OBI as a normal Go dependency

This would be the simplest solution but requires OBI team approval.

## References

- **OBI Repository:** https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation
- **OBI Collector Example:** https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation/tree/main/examples/otel-collector
- **Collector Builder Docs:** https://opentelemetry.io/docs/collector/custom-collector/

## Questions?

For questions about:
- **OBI functionality:** Open an issue in the OBI repository
- **Build integration:** Open an issue in this repository
- **Distribution inclusion:** Discuss in Collector SIG meetings
