# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:ecd58cb94454ed143639e44d28abf240d4f1c0f1566574a58d7f179027fad80d
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol.exe ./otelcol.exe
COPY config.yaml ./config.yaml

ENTRYPOINT ["otelcol.exe"]
CMD ["--config", "config.yaml"]
EXPOSE 4317 4318 55679
