# escape=`
ARG WIN_VERSION=2019
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}

COPY otelcol-otlp.exe ./otelcol-otlp.exe

ENV NO_WINDOWS_SERVICE=1

ENTRYPOINT ["otelcol-otlp.exe"]
EXPOSE 4317 4318
