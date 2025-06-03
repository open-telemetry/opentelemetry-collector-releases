FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 as certs
RUN apk --update add ca-certificates

FROM scratch

ARG USER_UID=10001
ARG USER_GID=10001
USER ${USER_UID}:${USER_GID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-k8s /otelcol-k8s
ENTRYPOINT ["/otelcol-k8s"]
# `4137` and `4318`: OTLP
# `55678`: OpenCensus
# `55679`: zpages
# `6831`, `14268`, and `14250`: Jaeger
# `9411`: Zipkin
EXPOSE 4317 4318 55678 55679 6831 14268 14250 9411
