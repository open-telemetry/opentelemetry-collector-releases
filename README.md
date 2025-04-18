# OpenTelemetry Collector distributions

> :warning: **Important note:** Git tags in this repository may change at any time to fix any issues found during a release. They are only meant to trigger Github releases and should not be relied upon.

This repository assembles OpenTelemetry Collector distributions, such as the "core" distribution, or "contrib".

Each distribution contains:

- Binaries for a multitude of platforms and architectures
- Multi-arch container images
- Packages to be used with Linux distributions (RPM, deb), Mac OS (brew) for the above-mentioned architectures

More details about each individual distribution can be seen in its own readme files.

Current list of distributions:

- [OpenTelemetry Collector (also known as "otelcol")](./distributions/otelcol)
- [OpenTelemetry Collector Contrib (also known as "otelcol-contrib")](./distributions/otelcol-contrib)
- [OpenTelemetry Collector for Kubernetes (also known as "otelcol-k8s")](./distributions/otelcol-k8s)
- [OpenTelemetry Collector OTLP (also known as "otelcol-otlp")](./distributions/otelcol-otlp)
- [OpenTelemetry Collector eBPF Profiler (also known as "otelcol-ebpf-profiler")](./distributions/otelcol-ebpf-profiler)

## Community

This repository is part of the Collector SIG. Check out the [Community section](https://github.com/open-telemetry/opentelemetry-collector?tab=readme-ov-file#community) on the main Collector repository to see how to get involved.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

Approvers ([@open-telemetry/collector-releases-approvers](https://github.com/orgs/open-telemetry/teams/collector-releases-approvers)):

- [Antoine Toulme](https://github.com/atoulme), Splunk
- [Christos Markou](https://github.com/ChrsMark), Elastic
- [Curtis Robert](https://github.com/crobert-1), Splunk
- [David Ashpole](https://github.com/dashpole), Google
- [John L. Peterson (Jack)](https://github.com/jackgopack4), Datadog
- [Matt Wear](https://github.com/mwear), Lightstep
- [Moritz Wiesinger](https://github.com/mowies), Dynatrace
- [Ziqi Zhao](https://github.com/fatsheep9146), Alibaba

Emeritus Approvers:

- [Anthony Mirabella](https://github.com/Aneurysm9)
- [Bryan Aguilar](https://github.com/bryan-aguilar)
- [Przemek Maciolek](https://github.com/pmm-sumo)
- [Ruslan Kovalov](https://github.com/kovrus)

Maintainers ([@open-telemetry/collector-contrib-maintainers](https://github.com/orgs/open-telemetry/teams/collector-contrib-maintainers)):

- [Alex Boten](https://github.com/codeboten), Honeycomb
- [Andrzej Stencel](https://github.com/andrzej-stencel), Elastic
- [Bogdan Drutu](https://github.com/bogdandrutu), Snowflake
- [Daniel Jaglowski](https://github.com/djaglowski), observIQ
- [Dmitrii Anoshin](https://github.com/dmitryax), Splunk
- [Evan Bradley](https://github.com/evan-bradley), Dynatrace
- [Juraci Paixão Kröhling](https://github.com/jpkrohling), Grafana Labs
- [Pablo Baeyens](https://github.com/mx-psi), DataDog
- [Sean Marciniak](https://github.com/MovieStoreGuy), Splunk
- [Tyler Helmuth](https://github.com/TylerHelmuth), Honeycomb
- [Yang Song](https://github.com/songy23), DataDog

Emeritus Maintainers

- [Tigran Najaryan](https://github.com/tigrannajaryan)

Learn more about roles in the [community repository](https://github.com/open-telemetry/community/blob/main/guides/contributor/membership.md).
