#!/bin/bash
REPO_DIR="$( cd "$(dirname $( dirname "${BASH_SOURCE[0]}" ))" &> /dev/null && pwd )"
BUILDER=$(which ocb)
GO=$(which go)

# default values
skipcompilation=false

while getopts d:s:b:g: flag
do
    case "${flag}" in
        d) distributions=${OPTARG};;
        s) skipcompilation=${OPTARG};;
        b) BUILDER=${OPTARG};;
        g) GO=${OPTARG};;
    esac
done

if [ -z $distributions ]; then
    echo "List of distributions to build not provided. Use '-d' to specify the names of the distributions to build. Ex.:"
    echo "$0 -d otelcol"
    exit 1
fi

if [ "$skipcompilation" = true ]; then
    echo "Skipping the compilation, we'll only generate the sources."
fi

echo "Distributions to build: $distributions";

for distribution in $(echo $distributions | tr "," "\n")
do
    pushd "${REPO_DIR}/distributions/${distribution}" > /dev/null
    mkdir -p _build
    echo Building: ${distribution}
    echo Using Builder: ${BUILDER}
    echo Using Go: ${GO}
    ${BUILDER} --skip-compilation=${skipcompilation} --go ${GO} --config manifest.yaml > _build/build.log 2>&1
    if [ $? != 0 ]; then
        echo "❌ ERROR: failed to build the distribution '${distribution}'."
        echo "🪵 Build logs for '${distribution}'"
        echo "----------------------"
        cat _build/build.log
        echo "----------------------"
        exit 1
    else
        echo "✅ SUCCESS: distribution '${distribution}' built."
    fi

    popd > /dev/null
done
