FROM alpine:3.24@sha256:a2d49ea686c2adfe3c992e47dc3b5e7fa6e6b5055609400dc2acaeb241c829f4 as certs
RUN apk --update add ca-certificates

FROM scratch

ARG USER_UID=10001
ARG USER_GID=10001
USER ${USER_UID}:${USER_GID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-k8s /otelcol-k8s
ENTRYPOINT ["/otelcol-k8s"]
# `4137` and `4318`: OTLP
# `55679`: zpages
# `6831`, `14268`, and `14250`: Jaeger
# `9411`: Zipkin
EXPOSE 4317 4318 55679 6831 14268 14250 9411
