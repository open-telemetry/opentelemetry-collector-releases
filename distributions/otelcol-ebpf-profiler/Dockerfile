FROM golang:1.25.3@sha256:0afe9b553123ec160b8282154b7156241ebf43c9033ea2bf9528c6c58c17e89d

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
