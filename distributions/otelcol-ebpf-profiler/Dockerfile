FROM golang:1.25.5@sha256:0f406d34b7cb7255d0700af02ec28a2c88f1e00701055f4c282aa4c3ec0b3245

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
