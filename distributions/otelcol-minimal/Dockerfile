FROM alpine:3.19 as certs
RUN apk --update add ca-certificates

FROM scratch

ARG USER_UID=10001
USER ${USER_UID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-minimal /otelcol-minimal
ENTRYPOINT ["/otelcol-minimal"]
CMD ["--config", "/etc/otelcol-minimal/config.yaml"]
EXPOSE 4317 55678 55679
