# OpenTelemetry Collector distributions

> :warning: **Important note:** Git tags in this repository may change at any time to fix any issues found during a release. They are only meant to trigger Github releases and should not be relied upon.

This repository assembles OpenTelemetry Collector distributions, such as the "core" distribution, or "contrib". It may contain non-official distributions, focused on specific use-cases, such as the load-balancer.

Each distribution contains:

- Binaries for a multitude of platforms and architectures (at least linux_amd64, linux_arm64, windows_amd64 and darwin_arm64)
- Multi-arch container images (at least amd64 and arm64)
- Packages to be used with Linux distributions (apk, RPM, deb), Mac OS (brew) for the above-mentioned architectures

More details about each individual distribution can be seen in its own readme files.

Current list of distributions:

- [OpenTelemetry Collector (also known as "otelcol")](./distributions/otelcol)
- [OpenTelemetry Collector Contrib (also known as "otelcol-contrib")](./distributions/otelcol-contrib)
