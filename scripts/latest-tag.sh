VERSION=''

while getopts v: flag
do
    case "${flag}" in
        v) VERSION=${OPTARG};;
    esac
done

[[ -n "$VERSION" ]] || echo "Unknown version. Use ./scripts/latest-tag.sh -v VERSION" || exit 1

# Push latest tag for otel/opentelemetry-collector-contrib
docker pull otel/opentelemetry-collector-contrib:${VERSION}-386
docker pull otel/opentelemetry-collector-contrib:${VERSION}-amd64
docker pull otel/opentelemetry-collector-contrib:${VERSION}-arm64
docker pull otel/opentelemetry-collector-contrib:${VERSION}-ppc64le
docker tag otel/opentelemetry-collector-contrib:${VERSION}-386 otel/opentelemetry-collector-contrib:latest-386
docker tag otel/opentelemetry-collector-contrib:${VERSION}-amd64 otel/opentelemetry-collector-contrib:latest-amd64
docker tag otel/opentelemetry-collector-contrib:${VERSION}-arm64 otel/opentelemetry-collector-contrib:latest-arm64
docker tag otel/opentelemetry-collector-contrib:${VERSION}-ppc64le otel/opentelemetry-collector-contrib:latest-ppc64le
docker push otel/opentelemetry-collector-contrib:latest-386
docker push otel/opentelemetry-collector-contrib:latest-amd64
docker push otel/opentelemetry-collector-contrib:latest-arm64
docker push otel/opentelemetry-collector-contrib:latest-ppc64le
docker manifest create otel/opentelemetry-collector-contrib:latest --amend otel/opentelemetry-collector-contrib:latest-386 --amend otel/opentelemetry-collector-contrib:latest-amd64 --amend otel/opentelemetry-collector-contrib:latest-arm64 --amend otel/opentelemetry-collector-contrib:latest-ppc64le
docker manifest push docker.io/otel/opentelemetry-collector-contrib:latest

# Push latest tag for otel/opentelemetry-collector
docker pull otel/opentelemetry-collector:${VERSION}-386
docker pull otel/opentelemetry-collector:${VERSION}-amd64
docker pull otel/opentelemetry-collector:${VERSION}-arm64
docker pull otel/opentelemetry-collector:${VERSION}-ppc64le
docker tag otel/opentelemetry-collector:${VERSION}-386 otel/opentelemetry-collector:latest-386
docker tag otel/opentelemetry-collector:${VERSION}-amd64 otel/opentelemetry-collector:latest-amd64
docker tag otel/opentelemetry-collector:${VERSION}-arm64 otel/opentelemetry-collector:latest-arm64
docker tag otel/opentelemetry-collector:${VERSION}-ppc64le otel/opentelemetry-collector:latest-ppc64le
docker push otel/opentelemetry-collector:latest-386
docker push otel/opentelemetry-collector:latest-amd64
docker push otel/opentelemetry-collector:latest-arm64
docker push otel/opentelemetry-collector:latest-ppc64le
docker manifest create otel/opentelemetry-collector:latest --amend otel/opentelemetry-collector:latest-386 --amend otel/opentelemetry-collector:latest-amd64 --amend otel/opentelemetry-collector:latest-arm64 --amend otel/opentelemetry-collector:latest-ppc64le
docker manifest push docker.io/otel/opentelemetry-collector:latest
