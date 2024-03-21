FROM alpine:3.16 as certs
RUN apk --update add ca-certificates

FROM scratch

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-k8s /otelcol-k8s
COPY config.yaml /etc/otelcol-k8s/config.yaml
ENTRYPOINT ["/otelcol-k8s"]
CMD ["--config", "/etc/otelcol-k8s/config.yaml"]
EXPOSE 4317 55678 55679
