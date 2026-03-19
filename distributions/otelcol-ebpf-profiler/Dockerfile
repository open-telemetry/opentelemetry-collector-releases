FROM golang:1.26.1@sha256:16e774b791968123d6af5ba4eec19cf91c4208cb1f5849efda5d4ffaf6d1c038

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
