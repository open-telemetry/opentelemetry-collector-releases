FROM alpine:3.21 as certs
RUN apk --update add ca-certificates

FROM scratch

ARG USER_UID=10001
ARG USER_GID=10001
USER ${USER_UID}:${USER_GID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-gateway /otelcol-gateway
COPY config.yaml /etc/otelcol-gateway/config.yaml
ENTRYPOINT ["/otelcol-gateway"]
CMD ["--config", "/etc/otelcol-gateway/config.yaml"]
EXPOSE 4317 4318 55678 55679
