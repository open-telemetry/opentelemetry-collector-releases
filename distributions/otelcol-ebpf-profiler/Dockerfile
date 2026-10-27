FROM alpine:3.24@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a as certs
RUN apk --update add 'ca-certificates>=20261020-r0'


FROM scratch

# By default, requires to run as root or with privileges to see all processes.
ARG USER_UID=0
ARG USER_GID=0
USER ${USER_UID}:${USER_GID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-ebpf-profiler /otelcol-ebpf-profiler
ENTRYPOINT ["/otelcol-ebpf-profiler"]
EXPOSE 4317 4318 55679
