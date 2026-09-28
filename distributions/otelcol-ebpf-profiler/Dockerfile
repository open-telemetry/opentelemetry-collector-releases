FROM alpine:3.24@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6 as certs
RUN apk --update add ca-certificates


FROM scratch

# By default, requires to run as root or with privileges to see all processes.
ARG USER_UID=0
ARG USER_GID=0
USER ${USER_UID}:${USER_GID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
EXPOSE 4317 4318 55679
