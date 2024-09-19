# OpenTelemetry Collector OTLP Distro

This distribution only contains the receiver and exporters for the OpenTelemetry Protocol (OTLP), including both gRPC and HTTP transport.

## Usage note

Unlike the Core and Contrib distributions, the deb/rpm/msi installers and the Docker images for this distribution do not provide a default configuration file. One will need to be created in the appropriate location (`/etc/otelcol-otlp/config.yaml` or `%ProgramW6432%\OpenTelemetry Collector\config.yaml` by default).

## Components

The full list of components is available in the [manifest](manifest.yaml)

### Rules for Component Inclusion

- Only `otlpreceiver`, `otlpexporter`, and `otlphttpexporter` are allowed.
