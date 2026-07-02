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

# Verify checksum against the exact entry in the upstream SHA256SUMS release asset.
echo "Verifying OBI ${OBI_VERSION} tarball checksum..."
if ! expected_checksum="$(
  curl --fail --show-error --location --silent --retry 3 --retry-delay 1 \
    "${BASE_URL}/SHA256SUMS" \
    | awk -v filename="${TARBALL}" '
        {
          checksum_filename = $2
          sub(/^\*/, "", checksum_filename)
          if (checksum_filename == filename) {
            print $1
            found = 1
            exit
          }
        }
        END { if (!found) exit 1 }
      '
)"; then
  rm -f "${TARBALL_CACHE}"
  echo "ERROR: ${TARBALL} not found in SHA256SUMS." >&2
  exit 1
fi

if [[ "$(uname -s)" == "Linux" ]]; then
  actual_checksum="$(sha256sum "${TARBALL_CACHE}" | awk '{print $1}')"
else
  actual_checksum="$(shasum -a 256 "${TARBALL_CACHE}" | awk '{print $1}')"
fi

if [[ "${actual_checksum}" != "${expected_checksum}" ]]; then
  rm -f "${TARBALL_CACHE}"
  echo "ERROR: checksum verification failed for ${TARBALL}." >&2
  exit 1
fi
echo "SHA256 verified: ${TARBALL}"

# Extract to OBI_DIR.
rm -rf "${OBI_DIR}"
mkdir -p "${OBI_DIR}"
tar xzf "${TARBALL_CACHE}" --strip-components=1 -C "${OBI_DIR}"

# Version-keyed stamp file: signals that this OBI version is ready and lets
# subsequent Make invocations (and the CI composite action) skip the download.
touch "${OBI_STAMP}"
echo "OBI ${OBI_VERSION} source prepared at ${OBI_DIR}"
