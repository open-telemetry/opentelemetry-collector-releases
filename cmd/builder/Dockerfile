FROM golang:1.27-alpine3.23@sha256:f70850ef6ce341e8b6427f677049145d63c0b7ce002d4d7e3d464828b1b96dcd
RUN apk --update add 'ca-certificates>=20261020-r0'

ARG SERVICE_NAME=ocb

RUN addgroup --gid 10001 --system ${SERVICE_NAME} && \
    adduser --ingroup ${SERVICE_NAME} --shell /bin/false \
    --disabled-password --uid 10001 ${SERVICE_NAME}

USER ${SERVICE_NAME}
WORKDIR /home/${SERVICE_NAME}

COPY --chmod=755 ocb /usr/local/bin/ocb
ENTRYPOINT [ "ocb" ]
