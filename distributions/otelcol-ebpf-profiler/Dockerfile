FROM golang:1.25.5@sha256:b6ba5234ed128185b0b81070813c77dd5c973aec7703f4646805f11377627408

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
