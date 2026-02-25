FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659 as certs
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
