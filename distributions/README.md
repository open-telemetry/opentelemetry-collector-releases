# Criteria for Supported Distributions

Support for a distribution implies all the duties covered by the Approver and Maintainer role requirements. In addition, support means that the Collector SIG is the owner and maintainer of the binaries/images of the different Collector distributions and is responsible for the pipeline that produces those artifacts.

Distributions supported by the Collector SIG should fulfill the following criteria:

1. To honor commitments made by OpenTelemetry to other Open Source projects the Collector SIG should support at least 1 distribution that includes support for Prometheus, Jaeger, Zipkin, and OpenCensus.
2. Distributions supported by the Collector SIG should serve a specific purpose and those purposes should have minimal overlap.
3. Distributions supported by the Collector SIG should meet general needs and not be too niche.
4. Distributions supported by the Collector SIG should only target the needs of the OpenTelemetry project.
5. Distributions supported by the Collector SIG are not required to be production ready and may be focused on development and proof of concept use cases.  The distribution should clearly indicate whether the Collector SIG considers it to be production ready.
6. Distributions supported by the Collector SIG must only include components from the `opentelemetry-collector` and `opentelemetry-collector-contrib` repositories.
7. Distributions supported by the Collector SIG should have a clearly defined list of criteria for which components are included.
8. Distributions supported by the Collector SIG must include the following assets except where the specific purpose of the distribution is naturally associated with a subset of these assets. In such cases, it should be clearly stated which assets are skipped and why.  Additional asset may be included if the distro desires:
    - Binaries for linux_amd64, linux_arm64, windows_amd64 and darwin_arm64
    - linux_amd64 and linux_arm64 container images
    - Packages to be used with Linux distributions (apk, RPM, deb), Mac OS (brew) for each distributed binary.


