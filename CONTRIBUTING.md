# Contribution guidelines

This repository contains a set of resources that ultimately results in OpenTelemetry Collector distributions. This document contains information needed to start contributing to this repository, including how to add new distributions.

## Understanding this repository

### Distribution directory

Each distribution has its own directory at the root of this repository, such as `opentelemetry-collector` or `opentelemetry-collector-loadbalancer`. Within each one of those, you'll find at least two files:

- `Dockerfile`, determining how to build the container image for this distribution
- `manifest.yaml`, which is used with the [opentelemetry-collector-builder](https://github.com/open-telemetry/opentelemetry-collector-builder) to generate the sources for the distribution.

Within each distribution, you are expected to be able to build it using the builder, like:

    $ opentelemetry-collector-builder --config manifest.yaml

You can build all distributions by running:

    $ make build

If you only interested in generating the sources for the distributions, use:

    $ make generate

### Distribution configurations

Due to an incompatibility between `goreleaser` and how this directory is structured, the default configuration files to be included in the container images should be placed under [./configs](./configs) instead of within the distribution's main directory.

### Scripts

The main `Makefile` is mostly a wrapper around scripts under the [./scripts](./scripts) directory.

### goreleaser

[goreleaser](https://goreleaser.com) plays a big role in producing the final artifacts. Given that the final configuration file for this tool would be huge and would cause relatively big changes for each new distribution, a `Makefile` target exists to generate the `.goreleaser.yaml` for the repository. The `Makefile` contains a list of distributions (`DISTRIBUTIONS`) that is passed down to the script, which will generate snippets based on the templates from `./scripts/goreleaser-templates/`. Adding a new distribution is then only a matter of adding the distribution's directory, plus adding it to the Makefile. Adding a new snippet is only a matter of adding a new template.

Once there's a change either to the templates or to the list of distributions, a new `.goreleaser` file can be generated with:

    $ make generate-goreleaser

