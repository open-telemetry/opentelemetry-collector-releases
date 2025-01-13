# escape=`
ARG WIN_VERSION=2019
FROM mcr.microsoft.com/windows/servercore:ltsc${WIN_BASE_IMAGE}

SHELL ["powershell", "-Command", "$ErrorActionPreference = 'Stop'; $ProgressPreference = 'SilentlyContinue';"]


COPY --chmod=755 otelcol-otlp ./otelcol-otlp.exe
COPY config.yaml ./config.yaml

ENV NO_WINDOWS_SERVICE=1

ENTRYPOINT ["otelcol-otlp.exe"]
CMD ["--config", "config.yaml"]
EXPOSE 13133 14250 14268 4317 6060 8888 9411 9443 9080
