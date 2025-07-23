FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1
RUN apk --update add ca-certificates

ARG SERVICE_NAME=opampsupervisor

RUN addgroup --gid 10001 --system ${SERVICE_NAME} && \
    adduser --ingroup ${SERVICE_NAME} --shell /bin/false \
    --disabled-password --uid 10001 ${SERVICE_NAME}

USER ${SERVICE_NAME}
WORKDIR /home/${SERVICE_NAME}

COPY --chmod=755 opampsupervisor /usr/local/bin/opampsupervisor
ENTRYPOINT [ "opampsupervisor" ]
