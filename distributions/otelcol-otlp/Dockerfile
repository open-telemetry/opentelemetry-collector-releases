FROM alpine:3.24@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf as certs
RUN apk --update add 'ca-certificates>=20261020-r0'

FROM scratch

ARG USER_UID=10001
ARG USER_GID=10001
USER ${USER_UID}:${USER_GID}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --chmod=755 otelcol-otlp /otelcol-otlp
ENTRYPOINT ["/otelcol-otlp"]
EXPOSE 4317 4318
