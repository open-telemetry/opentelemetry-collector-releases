FROM alpine:3.21 as certs
RUN apk --update add ca-certificates

FROM ghcr.io/open-telemetry/opentelemetry-collector-releases/opentelemetry-collector-opampsupervisor:0.126.0 AS supervisor
RUN mkdir -p /etc/otelcol-contrib/supervisor-data

FROM scratch

ARG USER_UID=10001
ARG USER_GID=10001
USER ${USER_UID}:${USER_GID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-contrib /otelcol-contrib
COPY --from=supervisor --chmod=755 /usr/local/bin/opampsupervisor /opampsupervisor
COPY --from=supervisor --chmod=644 --chown=10001:10001 /etc/otelcol-contrib/supervisor-data /etc/otelcol-contrib/supervisor-data
COPY config.yaml /etc/otelcol-contrib/config.yaml

ENTRYPOINT ["/opampsupervisor"]
EXPOSE 4317 4318 55678 55679
