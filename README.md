# OpenTelemetry Collector distributions

This repository assembles OpenTelemetry Collector distributions, such as the "core" distribution, or "contrib". It may contain non-official distributions, focused on specific use-cases, such as the load-balancer.

Each distribution contains:

- Binaries for a multitude of platforms and architectures
- Multi-arch container images (x86_64 and arm64)
- Packages to be used with Linux distributions (apk, RPM, deb), Mac OS (brew)

More details about each individual distribution can be seen in its own readme files.

Current list of distributions:

- [OpenTelemetry Collector (also known as "otelcol")](./distributions/otelcol)
- [OpenTelemetry Collector Contrib (also known as "otelcol-contrib")](./distributions/otelcol-contrib)
