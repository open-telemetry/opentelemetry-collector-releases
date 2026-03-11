FROM golang:1.25-alpine3.23@sha256:7993d12d3da2533bc2be7b2048c2d33f5872d52275c52f710dcf31a02e78b9b1
RUN apk --update add ca-certificates

ARG SERVICE_NAME=ocb

RUN addgroup --gid 10001 --system ${SERVICE_NAME} && \
    adduser --ingroup ${SERVICE_NAME} --shell /bin/false \
    --disabled-password --uid 10001 ${SERVICE_NAME}

USER ${SERVICE_NAME}
WORKDIR /home/${SERVICE_NAME}

COPY --chmod=755 ocb /usr/local/bin/ocb
ENTRYPOINT [ "ocb" ]
