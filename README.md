# OpenTelemetry Collector distributions

> :warning: **Important note:** Git tags in this repository may change at any time to fix any issues found during a release. They are only meant to trigger Github releases and should not be relied upon.

This repository assembles OpenTelemetry Collector distributions, such as the "core" distribution, or "contrib".

Each distribution contains:

- Binaries for a multitude of platforms and architectures
- Multi-arch container images
- Packages to be used with Linux distributions (RPM, deb), Mac OS (brew) for the above-mentioned architectures

More details about each individual distribution can be seen in its own readme files.

Current list of distributions:

- [OpenTelemetry Collector (also known as "otelcol")](./distributions/otelcol)
- [OpenTelemetry Collector Contrib (also known as "otelcol-contrib")](./distributions/otelcol-contrib)
- [OpenTelemetry Collector for Kubernetes (also known as "otelcol-k8s")](./distributions/otelcol-k8s)
