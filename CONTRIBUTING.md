# Contributing to OpenTelemetry Collector Distributions

## Introduction

Welcome to the OpenTelemetry Collector Distributions repository! 🎉

Thank you for considering contributing to this project. Whether you're adding new distributions, improving existing ones, or enhancing documentation, your efforts are invaluable to the community.

This repository assembles various OpenTelemetry Collector distributions, such as the "core" and "contrib" versions. Each distribution contains:

- Binaries for multiple platforms and architectures.
- Multi-architecture container images.
- Packages for Linux distributions (RPM, deb), Windows (msi) and macOS (brew).

For more details about each distribution, please refer to their respective directories within the repository.

---

## Prerequisites

Before you begin, ensure you have the following tools installed:

- **Go** – [Install Go](https://golang.org/doc/install)
- **Docker** – [Install Docker](https://docs.docker.com/get-docker/)
- **Make** – [Install Make](https://www.gnu.org/software/make/)

Additional Notes:

- Ensure your Go environment is set up correctly, with `$GOPATH` and `$PATH` configured.
- Docker is essential for building and testing container images.
- Familiarity with the OpenTelemetry Collector and its components will be beneficial.

---

## Repository Structure

Understanding the repository's structure will help you navigate and contribute effectively:

- **`distributions/`**: Contains directories for each distribution (e.g., `otelcol`, `otelcol-contrib`).
  - Each distribution directory includes:
    - `Dockerfile`: Defines how to build the container image for the distribution.
    - `manifest.yaml`: Used with the OpenTelemetry Collector Builder (`ocb`) to generate the sources for the distribution.

- **`scripts/`**: Contains scripts used by the main `Makefile` to automate various tasks.

---

## Workflow

### Pull Request Guidelines

- Fork the repository and create a new branch.
- Ensure your code adheres to the project's coding standards.
- Run all tests and ensure they pass before submitting.
- Link relevant issues in the PR description.

---

## Building Distributions

To build a specific distribution:

```bash
cd distributions/<distribution-name>
ocb --config manifest.yaml
```

To build all distributions:

```bash
make build
```

If you're only interested in generating the sources for the distributions:

```bash
make generate
```

---

## Generating goreleaser Configuration

`goreleaser` plays a significant role in producing the final artifacts. Given that the final configuration file for this tool would be extensive and would cause relatively big changes for each new distribution, a `Makefile` target exists to generate the `.goreleaser.yaml` for the repository.

To generate the `.goreleaser.yaml` file:

```bash
make generate-goreleaser
```

After generating the configuration, you can test the `goreleaser` build process with:

```bash
make goreleaser-verify
```

---

## Building Multi-Architecture Docker Images

`goreleaser` will build Docker images for various architectures, including `x86_64`, `386`, `arm`, `arm64`, and `ppc64le`. The build process involves executing `RUN` steps on the target architecture, which means the system you run it on needs support for emulating foreign architectures.

To set up the environment for building multi-architecture images:

```bash
# Install QEMU and related packages
sudo apt-get install qemu binfmt-support qemu-user-static

# Enable support for QEMU within Docker
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
```

---

## Adding Support for New Platforms or Architectures

When introducing a new collector distribution image or binary for a different platform or architecture, consider the following:

1. **Continuous Integration**:
   - Add the new platform or architecture to the CI test matrix for both the core and contrib repositories to ensure they can be compiled with the new combination. Failing to do so may cause the release to fail due to compilation issues on those uncovered platforms.

2. **goreleaser Configuration**:
   - In the `goreleaser/configure.go` file, add the new platform or architecture.
   - Regenerate the `.goreleaser.yaml` file (see the "Generating goreleaser Configuration" section above).

3. **GitHub Actions**:
   - In the `.github/workflows/ci-goreleaser.yaml` file, under the "Setup QEMU" action, add the new platform and architecture.
   - In the `.github/workflows/release.yaml` file, under the "Setup QEMU" action, add the new platform and architecture.

---

## Further Help

Need assistance? Join our community:

- **Slack Discussions**: [OpenTelemetry](https://cloud-native.slack.com/archives/CJFCJHG4Q)
- **Issues**: If you encounter a bug or have a feature request, [open an issue](https://github.com/open-telemetry/opentelemetry-collector-releases/issues)

---

Thank you for contributing! 🚀
