FROM golang:1.25.1@sha256:0caf875670e0ec9ebe7f4a9f4cf02add9d06ffccb055cc1066c83270c237dfb9

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
