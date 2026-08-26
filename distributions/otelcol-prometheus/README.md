# OpenTelemetry Collector Prometheus Distro

This distribution runs selected Prometheus exporters in-process as OpenTelemetry Collector metrics receivers. It preserves the exporters' familiar metric schemas while allowing the metrics to use Collector processors and OTLP or Prometheus-compatible exporters.

Exporter-backed receivers are compatibility alternatives to overlapping native OpenTelemetry receivers, not replacements for them.

## Recommendation

This distribution is experimental and is not recommended for production. Its bridge-backed receivers are currently at Development stability.

For production, build a custom Collector with the OpenTelemetry Collector Builder and include only the components required by your environment.

## Configuration

This distribution does not install a default configuration because each embedded exporter requires environment-specific targets, credentials, or files. Pass a configuration with `--config`.

The following example embeds `stackdriver_exporter` and sends its metrics over OTLP:

```yaml
extensions:
  health_check:

receivers:
  stackdriver_exporter:
    scrape_interval: 60s
    exporter_config:
      project_ids: ["my-project"]
      metrics_prefixes:
        - compute.googleapis.com/instance

processors:
  memory_limiter:
    check_interval: 1s
    limit_mib: 512
  batch:

exporters:
  otlp:
    endpoint: otel-backend.example.com:4317

service:
  extensions: [health_check]
  pipelines:
    metrics:
      receivers: [stackdriver_exporter]
      processors: [memory_limiter, batch]
      exporters: [otlp]
```

## Components

The full component list is available in the [manifest](manifest.yaml).

### Rules for Component Inclusion

Components included in this distribution are mainly [Prometheus exporters](https://prometheus.io/docs/instrumenting/exporters/) from the [Prometheus GitHub organization](https://github.com/prometheus) that have implemented the [Collector Receiver interfaces](https://pkg.go.dev/go.opentelemetry.io/collector/receiver). Exporters who implement said interface, but are not part of the Prometheus GitHub organization MAY be included, but require approval from the [Prometheus Interoperability SIG](https://github.com/open-telemetry/community/blob/main/sigs.md#prometheus-interoperability).

Components from [OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector) and [OpenTelemetry Collector Contrib](https://github.com/open-telemetry/opentelemetry-collector-contrib) can be added if they provide a better user experience for users of the Prometheus distribution.
