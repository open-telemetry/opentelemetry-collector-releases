# Axoflow Distribution for OpenTelemetry Collector

This repository assembles OpenTelemetry Collector distributions, such as the "core" distribution, or "contrib". It may contain non-official distributions, focused on specific use-cases, such as the load-balancer.

Each distribution contains:

- Binaries for a multitude of platforms and architectures (at least linux_amd64, linux_arm64, windows_amd64 and darwin_arm64)
- Multi-arch container images (at least amd64 and arm64)

More details about each individual distribution can be seen in its own readme files.

Current list of distributions:

- [Axoflow Distribution for OpenTelemetry Collector (also known as "axoflow-otel-collector")](./distributions/otelcol-contrib)
