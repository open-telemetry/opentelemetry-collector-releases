FROM golang:1.24.5@sha256:267159cb984d1d034fce6e9db8641bf347f80e5f2e913561ea98c40d5051cb67

COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
