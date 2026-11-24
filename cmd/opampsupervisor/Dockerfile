FROM alpine:3.24@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf
RUN apk --update add 'ca-certificates>=20261020-r0'

ARG SERVICE_NAME=opampsupervisor

RUN addgroup --gid 10001 --system ${SERVICE_NAME} && \
    adduser --ingroup ${SERVICE_NAME} --shell /bin/false \
    --disabled-password --uid 10001 ${SERVICE_NAME}

USER ${SERVICE_NAME}
WORKDIR /home/${SERVICE_NAME}

COPY --chmod=755 opampsupervisor /usr/local/bin/opampsupervisor
ENTRYPOINT [ "opampsupervisor" ]
