# OpenTelemetry Collector OTLP Distro

This distribution only contains the receiver and exporters for the OpenTelemetry Protocol (OTLP), including both gRPC and HTTP transport.

## Configuration

Unlike the Core and Contrib distributions, this distribution does not provide a default configuration file, and one will need to be created. The location of the config file is specified with the `--config` command line option.

- For the .deb/.rpm systemd service packages, the command line options are set in `/etc/otelcol-otlp/otelcol-otlp.conf`, and the default config path is `/etc/otelcol-otlp/config.yaml`.

- For the Windows installer, the command line options are set during the install process, and the default config path is `%ProgramW6432%\OpenTelemetry Collector\config.yaml`.

- For the Docker images, the command line options are blank by default, and must be specified with a `CMD` directive.
  
  Example: `CMD ["--config", "/etc/otelcol-otlp/config.yaml"]`

## Components

The full list of components is available in the [manifest](manifest.yaml)

### Rules for Component Inclusion

- Only `otlpreceiver`, `otlpexporter`, and `otlphttpexporter` are allowed.
