FROM golang:1.26.4@sha256:d184d9be4c13614e28498d632eeaaac704d662f18ad357e1df74a44424236cea

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
