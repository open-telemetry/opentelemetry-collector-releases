# escape=`
ARG WIN_VERSION=2019
ARG WIN_VERSION_SHA=sha256:ecd58cb94454ed143639e44d28abf240d4f1c0f1566574a58d7f179027fad80d
FROM mcr.microsoft.com/windows/nanoserver:ltsc${WIN_VERSION}@${WIN_VERSION_SHA}

COPY otelcol-k8s.exe ./otelcol-k8s.exe

ENTRYPOINT ["otelcol-k8s.exe"]
# `4137` and `4318`: OTLP
# `55679`: zpages
# `6831`, `14268`, and `14250`: Jaeger
# `9411`: Zipkin
EXPOSE 4317 4318 55679 6831 14268 14250 9411
