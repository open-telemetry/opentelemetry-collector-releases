This file contains the default configuration files for the container images. Due to a problem with goreleaser, these cannot be in the a directory with the same name as the binary for that image. For intance, if the binary is `opentelemetry-collector`, the container image cannot include anything from a directory named `opentelemetry-collector`, which is exactly what we have.

This problem is solvable in a better way, but for now, storing the configs in this directory is the simplest solution that works.
