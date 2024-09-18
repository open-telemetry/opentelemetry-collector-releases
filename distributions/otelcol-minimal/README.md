# OpenTelemetry Collector Minimal Distro

This distribution only contains the receiver and exporters for the OpenTelemetry Protocol (OTLP), supporting both gRPC and HTTP transport.

## Usage note

This distribution does not provide a default configuration file.

## Components

The full list of components is available in the [manifest](manifest.yaml)

### Rules for Component Inclusion

- Only `otlpreceiver`, `otlpexporter`, and `otlphttpexporter` are allowed.
