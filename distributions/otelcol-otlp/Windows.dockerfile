# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:a0d53c656749f02ad8490478f18f51e99b4a3fba23b1dd198083e94ed1620ab9
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol-otlp.exe ./otelcol-otlp.exe

ENTRYPOINT ["otelcol-otlp.exe"]
EXPOSE 4317 4318
