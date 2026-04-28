#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# This script creates and pushes the specified TAG to REMOTE

set -euo pipefail

if ! command -v yq &> /dev/null; then
    echo "This script requires 'yq'. Please install and try again."
    exit 1
fi

if [ -z "${TAG:-}" ]; then
    echo "TAG must be set (e.g. TAG=v0.100.0)"
    exit 1
fi
if [[ ! "${TAG}" =~ ^v.* ]]; then
    echo "TAG must start with lowercase 'v' (e.g. v0.100.0)"
    exit 1
fi

REMOTE="${REMOTE:-git@github.com:open-telemetry/opentelemetry-collector-releases.git}"
VALIDATE="${VALIDATE:-true}"
VERSION="${TAG#v}"

if [ "${VALIDATE}" = "true" ]; then
    for dir in distributions/*/; do
        manifest="${dir}manifest.yaml"
        if [ -f "${manifest}" ]; then
            dist_version=$(yq '.dist.version' "${manifest}")
            if [ "${dist_version}" != "${VERSION}" ]; then
                echo "Version mismatch in ${manifest}: dist.version is set to '${dist_version}', expected '${VERSION}'"
                echo "Please ensure the '[chore] Prepare release ${VERSION}' PR has been merged."
                echo "If this version mismatch is expected, please re-run with the extra argument 'VALIDATE=false'"
                exit 1
            fi
        fi
    done
fi

echo "Adding tag ${TAG}"
git tag -a "${TAG}" -s -m "Version ${TAG}"
echo "Pushing tag ${TAG}"
git push "${REMOTE}" "${TAG}"
