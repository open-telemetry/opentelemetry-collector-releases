# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:08d84cf1d9d6edd5f117bd955cc0e9d2a6db4b7f89d51c64518a7fe2372fb8d8
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol.exe ./otelcol.exe
COPY config.yaml ./config.yaml

ENTRYPOINT ["otelcol.exe"]
CMD ["--config", "config.yaml"]
EXPOSE 4317 4318 55679
