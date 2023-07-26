# Criteria for Supported Distributions

Support for a distribution implies all the duties covered by the Approver and Maintainer role requirements. In addition, support means that the Collector SIG is the owner and maintainer of the binaries/images of the different Collector distributions and is responsible for the pipeline that produces those artifacts.

Distributions supported by the Collector SIG should fulfill the following criteria:

1. Support at least one distribution that is recommended for production which includes support for Prometheus, Jaeger, Zipkin, and OpenCensus.
2. Serve a specific purpose and those purposes should have minimal overlap.
3. Meet general needs and not be too niche.
4. Only target the needs of the OpenTelemetry project.
5. May be focused on development and proof of concept use cases.  The distribution should clearly indicate whether the Collector SIG recommends the distribution be used in Production environments.
6. Must only include components from the `opentelemetry-collector` and `opentelemetry-collector-contrib` repositories.
7. Have a clearly defined list of criteria for which components are included.
8. Must include the following assets except where the specific purpose of the distribution is naturally associated with a subset of these assets. In such cases, it should be clearly stated which assets are skipped and why.  Additional asset may be included if the distro desires:
    - Binaries for linux_amd64, linux_arm64, windows_amd64 and darwin_arm64
    - linux_amd64 and linux_arm64 container images
    - Packages to be used with Linux distributions (apk, RPM, deb), Mac OS (brew) for each distributed binary.


