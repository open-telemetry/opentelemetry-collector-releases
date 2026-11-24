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
- [OpenTelemetry Collector Prometheus (also known as "otelcol-prometheus")](./distributions/otelcol-prometheus)
- [OpenTelemetry Collector eBPF Profiler (also known as "otelcol-ebpf-profiler")](./distributions/otelcol-ebpf-profiler)

## Linux systemd services

The packaged `otelcol`, `otelcol-contrib`, `otelcol-otlp` and
`otelcol-prometheus` services enable these protections by default:

- `NoNewPrivileges=true` prevents the process and its children from gaining
  privileges through setuid/setgid programs or file capabilities.
- `LockPersonality=true`, `RestrictRealtime=true` and `RestrictSUIDSGID=true`
  prevent ABI personality changes, realtime scheduling and creation of
  setuid/setgid files.
- `ProtectControlGroups=true`, `ProtectKernelModules=true`,
  `ProtectKernelTunables=true` and `ProtectHostname=true` prevent changes to
  the host's cgroups, kernel modules, kernel tunables and hostname.

The complete set requires systemd 242 or newer. Older systemd releases ignore
unknown directives, so some protections will not apply. Availability also depends
on kernel support and the service manager's ability to create namespaces and apply
seccomp filters.

These settings preserve access to application logs, devices and networking for
receivers, and allow writes to configured storage directories. Stronger
restrictions depend on the components and paths used by each deployment. For
example, hiding devices can break hardware monitoring, and making the entire
filesystem read-only can prevent persistent queues from writing their state.

Use a systemd drop-in to customize the service instead of editing the packaged
unit. If a receiver needs a privileged port, grant the capability through systemd
rather than relying on a file capability, which `NoNewPrivileges` prevents from
adding privileges at execution:

```sh
sudo systemctl edit otelcol-contrib
```

```ini
[Service]
AmbientCapabilities=CAP_NET_BIND_SERVICE
```

After saving, restart the service with `sudo systemctl restart otelcol-contrib`.
Other workloads that require a restricted operation can override the relevant
setting in a drop-in after assessing the required permissions. The OpAMP
Supervisor and distributions without a packaged systemd service are outside
this change.

## Community

This repository is part of the Collector SIG. Check out the [Community section](https://github.com/open-telemetry/opentelemetry-collector?tab=readme-ov-file#community) on the main Collector repository to see how to get involved.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

### Maintainers

- [Alex Boten](https://github.com/codeboten), Grafana Labs
- [Andrzej Stencel](https://github.com/andrzej-stencel), Elastic
- [Antoine Toulme](https://github.com/atoulme), Splunk
- [Bogdan Drutu](https://github.com/bogdandrutu), Snowflake
- [Damien Mathieu](https://github.com/dmathieu), Elastic
- [Dmitrii Anoshin](https://github.com/dmitryax), Splunk
- [Evan Bradley](https://github.com/evan-bradley), Dynatrace
- [Moritz Wiesinger](https://github.com/mowies), Dynatrace
- [Pablo Baeyens](https://github.com/mx-psi), DataDog
- [Sean Marciniak](https://github.com/MovieStoreGuy), Splunk
- [Tyler Helmuth](https://github.com/TylerHelmuth), Grafana Labs
- [Yang Song](https://github.com/songy23), DataDog

For more information about the maintainer role, see the [community repository](https://github.com/open-telemetry/community/blob/main/guides/contributor/membership.md#maintainer).

### Approvers

- [Christos Markou](https://github.com/ChrsMark), Elastic
- [Curtis Robert](https://github.com/crobert-1), Splunk
- [Douglas Camata](https://github.com/douglascamata), Coralogix
- [Sam DeHaan](https://github.com/dehaansa), Grafana Labs

For more information about the approver role, see the [community repository](https://github.com/open-telemetry/community/blob/main/guides/contributor/membership.md#approver).

### Emeritus

- [Anthony Mirabella](https://github.com/Aneurysm9), Approver
- [Bryan Aguilar](https://github.com/bryan-aguilar), Approver
- [David Ashpole](https://github.com/dashpole), Approver
- [Daniel Jaglowski](https://github.com/djaglowski), Maintainer
- [Edmo Vamerlatti Costa](https://github.com/edmocosta), Maintainer
- [John L. Peterson (Jack)](https://github.com/jackgopack4), Approver
- [Juraci Paixão Kröhling](https://github.com/jpkrohling), Maintainer
- [Matt Wear](https://github.com/mwear), Approver
- [Przemek Maciolek](https://github.com/pmm-sumo), Approver
- [Ruslan Kovalov](https://github.com/kovrus), Approver
- [Tigran Najaryan](https://github.com/tigrannajaryan), Maintainer
- [Ziqi Zhao](https://github.com/fatsheep9146), Approver

For more information about the emeritus role, see the [community repository](https://github.com/open-telemetry/community/blob/main/guides/contributor/membership.md#emeritus-maintainerapprovertriager).
