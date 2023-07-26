# OpenTelemetry Collector distributions

> :warning: **Important note:** Git tags in this repository may change at any time to fix any issues found during a release. They are only meant to trigger Github releases and should not be relied upon.

This repository assembles OpenTelemetry Collector distributions, such as the "core" distribution, or "contrib". It may contain non-official distributions, focused on specific use-cases, such as the load-balancer.

Each distribution contains:

- Binaries for a multitude of platforms and architectures (at least linux_amd64, linux_arm64, windows_amd64 and darwin_arm64)
- Multi-arch container images (at least amd64 and arm64)
- Packages to be used with Linux distributions (apk, RPM, deb), Mac OS (brew) for the above mentioned architectures

More details about each individual distribution can be seen in its own readme files.

Current list of distributions:

- [OpenTelemetry Collector (also known as "otelcol")](./distributions/otelcol)
- [OpenTelemetry Collector Contrib (also known as "otelcol-contrib")](./distributions/otelcol-contrib)

## Criteria for Distributions

- To honor commitments made by OpenTelemetry to other Open Source projects the Collector SIG should support at least 1 distribution that includes support for Prometheus, Jaeger, Zipkin, and OpenCensus.
- Distributions supported by the Collector SIG should serve a specific purpose and those purposes should have minimal overlap.
- Distributions supported by the Collector SIG should meet general needs and not be too niche.
- Distributions supported by the Collector SIG should only target the needs of the OpenTelemetry project.
- Distributions supported by the Collector SIG are not required to be production ready and may be focused on development and proof of concept use cases.  The distribution should clearly indicate whether the Collector SIG considers it to be production ready.
- Distributions supported by the Collector SIG must only include components from the `opentelemetry-collector` and `opentelemetry-collector-contrib` repositories.
- Distributions supported by the Collector SIG should have a clearly defined list of criteria for which components are included.
- Distributions supported by the Collector SIG must include the following assets except where the specific purpose of the distribution is naturally associated with a subset of these assets. In such cases, it should be clearly stated which assets are skipped and why.  Additional asset may be included if the distro desires:
    - Binaries for linux_amd64, linux_arm64, windows_amd64 and darwin_arm64
    - linux_amd64 and linux_arm64 container images
    - Packages to be used with Linux distributions (apk, RPM, deb), Mac OS (brew) for each distributed binary.

Support for a distribution implies all the duties covered by the Approver and Maintainer role requirements. In addition, support means that the Collector SIG is the owner and maintainer of the binaries/images of the different Collector distributions and is responsible for the pipeline that produces those artifacts.
