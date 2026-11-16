# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:14564f500b1e714e66652f7d0a101ea8917b4d79b8ac878e7f6546ff087a0fc8
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol.exe ./otelcol.exe
COPY config.yaml ./config.yaml

ENTRYPOINT ["otelcol.exe"]
CMD ["--config", "config.yaml"]
EXPOSE 4317 4318 55679
