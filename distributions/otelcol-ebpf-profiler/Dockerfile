FROM golang:1.24.5@sha256:fdcd2e5a34587bd5426c90e1531fd5ba448c89bb738df0f33860dfc69439a1f5

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
