# escape=`
ARG WIN_VERSION=2019
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}

COPY otelcol-k8s.exe ./otelcol-k8s.exe

ENV NO_WINDOWS_SERVICE=1

ENTRYPOINT ["otelcol-k8s.exe"]
# `4137` and `4318`: OTLP
# `55678`: OpenCensus
# `55679`: zpages
# `6831`, `14268`, and `14250`: Jaeger
# `9411`: Zipkin
EXPOSE 4317 4318 55678 55679 6831 14268 14250 9411
