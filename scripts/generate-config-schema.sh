#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# Bundles each distribution's config JSON Schema via cmd/configschemagen,
# from the manifest.yaml and the already-generated OCB build directory
# (distributions/<name>/_build, produced by `make generate-sources`).

REPO_DIR="$( cd "$(dirname "$( dirname "${BASH_SOURCE[0]}" )")" &> /dev/null && pwd )"

while getopts d:g: flag
do
    case "${flag}" in
        d) distributions=${OPTARG};;
        g) GO=${OPTARG};;
        *) exit 1;;
    esac
done

[[ -n "$GO" ]] || GO='go'

if [[ -z $distributions ]]; then
    echo "List of distributions to generate schemas for not provided. Use '-d' to specify the names of the distributions. Ex.:"
    echo "$0 -d otelcol"
    exit 1
fi

for distribution in $(echo "$distributions" | tr "," "\n")
do
    manifest="${REPO_DIR}/distributions/${distribution}/manifest.yaml"
    build_dir="${REPO_DIR}/distributions/${distribution}/_build"
    out="${build_dir}/config.schema.json"

    if [[ ! -f "$manifest" ]]; then
        echo "Skipping ${distribution}: no manifest.yaml"
        continue
    fi
    if [[ ! -f "${build_dir}/go.mod" ]]; then
        echo "Skipping ${distribution}: no generated build dir at ${build_dir} (run 'make generate-sources' first)"
        continue
    fi

    echo "Generating config schema for: $distribution"
    (cd "$REPO_DIR" && "$GO" run ./cmd/configschemagen -manifest "$manifest" -build-dir "$build_dir" -o "$out")
done
