FROM golang:1.25.5@sha256:0ece421d4bb2525b7c0b4cad5791d52be38edf4807582407525ca353a429eccc

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
