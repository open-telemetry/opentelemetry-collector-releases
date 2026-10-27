# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:08d84cf1d9d6edd5f117bd955cc0e9d2a6db4b7f89d51c64518a7fe2372fb8d8
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol-k8s.exe ./otelcol-k8s.exe

ENTRYPOINT ["otelcol-k8s.exe"]
# `4137` and `4318`: OTLP
# `55679`: zpages
# `6831`, `14268`, and `14250`: Jaeger
# `9411`: Zipkin
EXPOSE 4317 4318 55679 6831 14268 14250 9411
