FROM golang:1.24.5@sha256:14fd8a55e59a560704e5fc44970b301d00d344e45d6b914dda228e09f359a088

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
