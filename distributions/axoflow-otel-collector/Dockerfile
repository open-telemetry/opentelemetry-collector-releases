FROM alpine:3.16 as certs
RUN apk --update add ca-certificates

FROM scratch

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 axoflow-otel-collector /axoflow-otel-collector
COPY configs/axoflow-otel-collector.yaml /etc/axoflow-otel-collector/config.yaml
ENTRYPOINT ["/axoflow-otel-collector"]
CMD ["--config", "/etc/axoflow-otel-collector/config.yaml"]
EXPOSE 4317 55678 55679
