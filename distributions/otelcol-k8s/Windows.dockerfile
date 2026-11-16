# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:14564f500b1e714e66652f7d0a101ea8917b4d79b8ac878e7f6546ff087a0fc8
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol-k8s.exe ./otelcol-k8s.exe

ENTRYPOINT ["otelcol-k8s.exe"]
# `4137` and `4318`: OTLP
# `55679`: zpages
# `6831`, `14268`, and `14250`: Jaeger
# `9411`: Zipkin
EXPOSE 4317 4318 55679 6831 14268 14250 9411
