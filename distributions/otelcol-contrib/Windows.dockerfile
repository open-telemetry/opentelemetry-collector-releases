# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:5ba1caa6da010f2baf95c4fd7605c4eecbeec5ca156602c70d4281b2b3c98360
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol-contrib.exe ./otelcol-contrib.exe
COPY config.yaml ./config.yaml

ENTRYPOINT ["otelcol-contrib.exe"]
CMD ["--config", "config.yaml"]
EXPOSE 4317 4318 55679
