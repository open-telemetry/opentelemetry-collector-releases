FROM golang:1.25.1@sha256:d6bdb04c3109d29739c19566092e63b469d4cdf193df206fadf286aa4a549ba6

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
