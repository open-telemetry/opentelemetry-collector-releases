#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# Downloads the OBI (opentelemetry-ebpf-instrumentation) release tarball
# (obi-vX.Y.Z-source-generated.tar.gz), which includes all pre-generated BPF
# source files, verifies its SHA256 checksum, and extracts it to internal/obi-src.
#
# No BPF toolchain (Docker, clang, bpf2go) is required.
#
# The version is read from distributions/otelcol-contrib/manifest.yaml unless
# OBI_VERSION is already set in the environment (e.g. by the CI composite action).

set -euo pipefail

REPO_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." &> /dev/null && pwd )"
DISTRIBUTIONS="${1:-}"
MANIFEST="${REPO_DIR}/distributions/otelcol-contrib/manifest.yaml"
OBI_DIR="${REPO_DIR}/internal/obi-src"

needs_obi=false
if [[ "$DISTRIBUTIONS" == *"otelcol-contrib"* ]]; then
  needs_obi=true
fi

if [[ "${needs_obi}" != "true" ]]; then
  exit 0
fi

if [[ ! -f "${MANIFEST}" ]]; then
  echo "ERROR: OBI manifest not found at ${MANIFEST}" >&2
  exit 1
fi

OBI_VERSION="${OBI_VERSION:-$(
  awk '/- gomod: go\.opentelemetry\.io\/obi / {print $NF; exit}' "${MANIFEST}"
)}"

if [[ -z "${OBI_VERSION}" ]]; then
  echo "ERROR: failed to resolve OBI version from ${MANIFEST}" >&2
  exit 1
fi

TARBALL="obi-${OBI_VERSION}-source-generated.tar.gz"
BASE_URL="https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation/releases/download/${OBI_VERSION}"
TARBALL_CACHE="${REPO_DIR}/.local/${TARBALL}"
OBI_STAMP="${OBI_DIR}/.obi-${OBI_VERSION}"

# Fast path: version-keyed stamp file already exists → nothing to do.
if [[ -f "${OBI_STAMP}" ]]; then
  echo "OBI ${OBI_VERSION} already prepared at ${OBI_DIR}"
  exit 0
fi

echo "Fetching OBI ${OBI_VERSION} source-generated tarball..."

# Download the tarball to .local/ if it is not already cached there.
mkdir -p "${REPO_DIR}/.local"
if [[ ! -f "${TARBALL_CACHE}" ]]; then
  curl --fail --show-error --location --retry 3 --retry-delay 1 \
    --output "${TARBALL_CACHE}" "${BASE_URL}/${TARBALL}"
fi

# Verify checksum against the upstream SHA256SUMS release asset.
echo "Verifying OBI ${OBI_VERSION} tarball checksum..."
if [[ "$(uname -s)" == "Linux" ]]; then
  curl -fsSL "${BASE_URL}/SHA256SUMS" | grep -F "${TARBALL}" \
    | (cd "${REPO_DIR}/.local" && sha256sum --check) \
    || { rm -f "${TARBALL_CACHE}"; echo "ERROR: checksum verification failed." >&2; exit 1; }
else
  curl -fsSL "${BASE_URL}/SHA256SUMS" | grep -F "${TARBALL}" \
    | (cd "${REPO_DIR}/.local" && shasum -a 256 --check) \
    || { rm -f "${TARBALL_CACHE}"; echo "ERROR: checksum verification failed." >&2; exit 1; }
fi

# Extract to OBI_DIR.
rm -rf "${OBI_DIR}"
mkdir -p "${OBI_DIR}"
tar xzf "${TARBALL_CACHE}" --strip-components=1 -C "${OBI_DIR}"

# Version-keyed stamp file: signals that this OBI version is ready and lets
# subsequent Make invocations (and the CI composite action) skip the download.
touch "${OBI_STAMP}"
echo "OBI ${OBI_VERSION} source prepared at ${OBI_DIR}"
