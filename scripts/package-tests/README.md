# Build and test deb/rpm packages

## Prerequisites

Tools:

- [Go](https://go.dev/)
- [GoReleaser](https://goreleaser.com/)
- [Podman](https://podman.io/)
- make

## How to build and test

To build the Collector Linux packages, a few steps are required:

- Run `make generate` to (re-)generate sources and GoReleaser files
- Go to the distribution folder that you want to build (under the `distributions` folder)
- Run `goreleaser release --snapshot --clean --skip sbom,sign,archive,docker`
    - This will build the necessary release assets with all architectures and packaging types into the `dist` folder inside your
      current folder. (We can skip many parts of the release build that we don't need for running the package tests locally)
    - We use GoReleaser Pro only features in CI. If you want to run this locally, and you run into `unmarshal` errors, 
    you may have to remove the parts that goreleaser complains about or use a pro license.
- Go back to the root of the repo
- To start the package tests,
  run: `./scripts/package-tests/package-tests.sh ./distributions/<otelcol|otelcol-contrib>/dist/<otelcol|otelcol-contrib>_*-SNAPSHOT-*_linux_amd64.<deb|rpm> <otelcol|otelcol-contrib>`
