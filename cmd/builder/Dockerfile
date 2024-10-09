FROM golang:1.23-alpine3.20
RUN apk --update add ca-certificates

ARG SERVICE_NAME=ocb

RUN addgroup --gid 1001 -S ${SERVICE_NAME} && \
    adduser -G ${SERVICE_NAME} --shell /bin/false --disabled-password \
    -H --uid 1001 ${SERVICE_NAME} && \
    chown ${SERVICE_NAME}:${SERVICE_NAME} /usr/local/go/bin

USER ${SERVICE_NAME}
WORKDIR /home/${SERVICE_NAME}

COPY --chmod=755 ocb /usr/local/bin/ocb
ENTRYPOINT [ "ocb" ]