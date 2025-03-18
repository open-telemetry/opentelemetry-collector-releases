# escape=`
ARG WIN_VERSION=2019
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}

COPY otelcol-otlp.exe ./otelcol-otlp.exe

ENTRYPOINT ["otelcol-otlp.exe"]
EXPOSE 4317 4318
