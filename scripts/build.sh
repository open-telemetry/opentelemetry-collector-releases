#!/bin/bash

REPO_DIR="$( cd "$(dirname "$( dirname "${BASH_SOURCE[0]}" )")" &> /dev/null && pwd )"
BUILDER=''
GO=''
GOPROXY='direct'

# default values
skipcompilation=false

while getopts d:s:b:g:l: flag
do
    case "${flag}" in
        d) distributions=${OPTARG};;
        s) skipcompilation=${OPTARG};;
        b) BUILDER=${OPTARG};;
        g) GO=${OPTARG};;
        l) latest=${OPTARG};;
        *) exit 1;;
    esac
done

[[ -n "$BUILDER" ]] || BUILDER='ocb'
[[ -n "$GO" ]] || GO='go'

if [[ -z $distributions ]]; then
    echo "List of distributions to build not provided. Use '-d' to specify the names of the distributions to build. Ex.:"
    echo "$0 -d otelcol"
    exit 1
fi

if [[ "$skipcompilation" = true ]]; then
    echo "Skipping the compilation, we'll only generate the sources."
fi

if [[ "$latest" = true ]]; then
    echo "Using latest commits instead of pinned versions."
fi

echo "Distributions to build: $distributions";

for distribution in $(echo "$distributions" | tr "," "\n")
do
    pushd "${REPO_DIR}/distributions/${distribution}" > /dev/null || exit
    mkdir -p _build

    echo "Building: $distribution"
    echo "Using Builder: $(command -v "$BUILDER")"
    echo "Using Go: $(command -v "$GO")"

    if [[ "$latest" = true ]]; then
        echo "Using latest main versions for all components."
        sed -i 's/\(gomod: github.com\/open-telemetry\/opentelemetry-collector-contrib.*\?\) v[0-9]\.[0-9]\+\.[0-9]\+/\1 main/' manifest.yaml
        sed -i 's/\(gomod: go\.opentelemetry\.io\/collector.*\?\) v[0-9]\.[0-9]\+\.[0-9]\+/\1 main/' manifest.yaml
    fi

    if "$BUILDER" --skip-compilation="${skipcompilation}" --skip-strict-versioning --go "$GO" --verbose --config manifest.yaml > _build/build.log 2>&1; then
        echo "✅ SUCCESS: distribution '${distribution}' built."
    else
        echo "❌ ERROR: failed to build the distribution '${distribution}'."
        echo "🪵 Build logs for '${distribution}'"
        echo "----------------------"
        cat _build/build.log
        echo "----------------------"
        exit 1
    fi

    popd > /dev/null || exit
done
