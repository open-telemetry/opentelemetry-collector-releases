# OpenTelemetry Collector Core Distro

This distribution contains all the components from the [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector) repository and a small selection of components tied to open source projects from the [OpenTelemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) repository.

This distribution is considered "classic" and is no longer accepting new components outside of components from the Core repo.

## Components

The full list of components is available in the [manifest](manifest.yaml)

### Rules for Component Inclusion

Since Core is a "classic" distribution its components are strictly limited to what currently exists in its [manifest](manifest.yaml) and any future components in Core.
No other components from Contrib should be added.
