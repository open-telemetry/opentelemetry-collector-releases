# Build and test deb/rpm/apk packages

## Prerequisites

Tools:

- [Go](https://go.dev/)
- [GoReleaser](https://goreleaser.com/)
- [Podman](https://podman.io/)
- make

## How to build and test

To build the collector linux packages, a few steps are required:

- Run `make generate` to (re-)generate sources and goreleaser files
- Go to the distribution folder that you want to build (under the `distributions` folder)
- Run `goreleaser release --snapshot --clean --skip sbom,sign,archive,docker`
    - This will build a full release with all architectures and packaging types into the `dist` folder inside your
      current folder. (We can skip many parts of the release build that we don't need)
    - If you run into `unmarshal` errors, start to remove the parts that goreleaser complains about. This likely happens
      due to a missing goreleaser pro license and therefore feature that you can't use.
- Go back to the root of the repo
- To start the package tests,
  run: `./scripts/package-tests/package-tests.sh ./distributions/<otelcol|otelcol-contrib>/dist/<otelcol|otelcol-contrib>_*-SNAPSHOT-*_linux_amd64.<deb|rpm> <otelcol|otelcol-contrib>`

