FROM golang:1.23-alpine3.20
RUN apk --update add ca-certificates

ARG SERVICE_NAME=ocb

RUN addgroup --gid 10001 --system ${SERVICE_NAME} && \
    adduser --ingroup ${SERVICE_NAME} --shell /bin/false \
    --disabled-password --uid 10001 ${SERVICE_NAME}

USER ${SERVICE_NAME}
WORKDIR /home/${SERVICE_NAME}

COPY --chmod=755 ocb /usr/local/bin/ocb
ENTRYPOINT [ "ocb" ]
