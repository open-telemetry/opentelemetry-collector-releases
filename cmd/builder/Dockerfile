FROM golang:1.27-alpine3.23@sha256:0908ac9b9319e09d7c238aabe914e0395c51d63c4e3d0ae8c554fda9158a5769
RUN apk --update add ca-certificates

ARG SERVICE_NAME=ocb

RUN addgroup --gid 10001 --system ${SERVICE_NAME} && \
    adduser --ingroup ${SERVICE_NAME} --shell /bin/false \
    --disabled-password --uid 10001 ${SERVICE_NAME}

USER ${SERVICE_NAME}
WORKDIR /home/${SERVICE_NAME}

COPY --chmod=755 ocb /usr/local/bin/ocb
ENTRYPOINT [ "ocb" ]
