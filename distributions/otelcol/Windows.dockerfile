# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:217694d5470363a77b86c42d439090b7466cd959ad5a15221c556a4707481305
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol.exe ./otelcol.exe
COPY config.yaml ./config.yaml

ENTRYPOINT ["otelcol.exe"]
CMD ["--config", "config.yaml"]
EXPOSE 4317 4318 55679
