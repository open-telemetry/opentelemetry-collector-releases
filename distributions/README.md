# Criteria for Supported Distributions

Support for a distribution implies all the duties covered by the Approver and Maintainer role requirements. In addition, support means that the Collector SIG is the owner and maintainer of the binaries/images of the different Collector distributions and is responsible for the pipeline that produces those artifacts.

The Collector SIG will support at least one distribution that is recommended for production which includes support for Prometheus, Jaeger, Zipkin, and OpenCensus.

Distributions supported by the Collector SIG should fulfill the following criteria:

1. Serve a specific purpose that has minimal overlap with the purpose of any other distribution.
2. Should meet general needs and be desired by many users.
3. Should not be specific to any vendor.
4. May be focused on development or proof of concept use cases.  The distribution should clearly indicate whether the Collector SIG recommends the distribution be used in production environments.
5. Must only include components from the `opentelemetry-collector` and `opentelemetry-collector-contrib` repositories.
    - Components that are marked [Unmaintained](https://github.com/open-telemetry/opentelemetry-collector#unmaintained) will be kept in any distributions for six months. After six months of being unmaintained the component will be removed from the distributions.
6. Have a clearly defined list of criteria for which components are included.
7. Must include the following assets except where the specific purpose of the distribution is naturally associated with a subset of these assets. In such cases, it should be clearly stated which assets are skipped and why.  Additional assets may be included if the maintainers agree:
    - Binaries for linux_amd64, linux_arm64, windows_amd64 and darwin_arm64
    - linux_amd64 and linux_arm64 container images
    - Packages to be used with Linux distributions (RPM, deb), macOS (brew), Windows (MSI) for each distributed binary.


