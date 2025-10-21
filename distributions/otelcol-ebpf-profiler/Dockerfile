FROM golang:1.25.3@sha256:ffa2e570108dd80c155d6ea9447b2410d0ed739e8cc9e256d6bd5d818c7a03e2

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
