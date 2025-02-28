# OpenTelemetry Collector Gateway Distro

This distribution contains the components from both the [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector) repository and the [OpenTelemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) repository. This distribution includes open source and vendor supported components.

This distribution is intended to receive telemetry signal in various protocols and export data with OTLP. It doesn't contain any active scrapers. When idle, the distribution does not use any resources.
The debug and file exporters are also present to help debug.
The distribution accepts any extensions, processors, connectors and configuration providers of [Alpha stability](https://github.com/open-telemetry/opentelemetry-collector#alpha) or higher.

## Components

The full list of components is available in the [manifest](manifest.yaml)

### Rules for Component Inclusion

- Include all extensions at [Alpha stability](https://github.com/open-telemetry/opentelemetry-collector#alpha) or higher and pipeline components that have at least 1 signal at [Alpha stability](https://github.com/open-telemetry/opentelemetry-collector#alpha) or higher.
