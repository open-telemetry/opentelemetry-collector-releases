# escape=`
ARG WIN_VERSION=2019
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}

COPY otelcol-contrib.exe ./otelcol-contrib.exe
COPY config.yaml ./config.yaml

ENV NO_WINDOWS_SERVICE=1

ENTRYPOINT ["otelcol-contrib.exe"]
CMD ["--config", "config.yaml"]
EXPOSE 4317 4318 55678 55679
