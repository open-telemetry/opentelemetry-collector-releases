FROM golang:1.24-alpine3.20@sha256:9f98e9893fbc798c710f3432baa1e0ac6127799127c3101d2c263c3a954f0abe
RUN apk --update add ca-certificates

ARG SERVICE_NAME=ocb

RUN addgroup --gid 10001 --system ${SERVICE_NAME} && \
    adduser --ingroup ${SERVICE_NAME} --shell /bin/false \
    --disabled-password --uid 10001 ${SERVICE_NAME}

USER ${SERVICE_NAME}
WORKDIR /home/${SERVICE_NAME}

COPY --chmod=755 ocb /usr/local/bin/ocb
ENTRYPOINT [ "ocb" ]