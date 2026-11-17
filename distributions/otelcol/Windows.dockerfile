# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:a0d53c656749f02ad8490478f18f51e99b4a3fba23b1dd198083e94ed1620ab9
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol.exe ./otelcol.exe
COPY config.yaml ./config.yaml

ENTRYPOINT ["otelcol.exe"]
CMD ["--config", "config.yaml"]
EXPOSE 4317 4318 55679
