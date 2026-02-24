#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

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

obi_module_version="$(
  awk '/- gomod: go\.opentelemetry\.io\/obi / {print $NF; exit}' "${MANIFEST}"
)"

if [[ -z "${obi_module_version}" ]]; then
  echo "ERROR: failed to resolve OBI version from ${MANIFEST}" >&2
  exit 1
fi

obi_archive_url="https://github.com/open-telemetry/opentelemetry-ebpf-instrumentation/archive/refs/tags/${obi_module_version}.tar.gz"
version_file="${OBI_DIR}/.obi-version"
prepared_version=""
if [[ -f "${version_file}" ]]; then
  prepared_version="$(cat "${version_file}")"
fi

if [[ "${prepared_version}" != "${obi_module_version}" ]]; then
  tmpdir="$(mktemp -d)"
  archive="${tmpdir}/obi-source.tar.gz"
  trap 'rm -rf "${tmpdir}"' EXIT

  echo "Preparing OBI source ${obi_module_version} from ${obi_archive_url}"
  curl --fail --show-error --location --retry 3 --retry-delay 1 \
    --output "${archive}" "${obi_archive_url}"
  tar -xzf "${archive}" -C "${tmpdir}"

  extracted_dir="$(find "${tmpdir}" -mindepth 1 -maxdepth 1 -type d | head -n1)"
  if [[ -z "${extracted_dir}" ]]; then
    echo "ERROR: failed to unpack OBI archive ${obi_archive_url}" >&2
    exit 1
  fi

  rm -rf "${OBI_DIR}"
  mv "${extracted_dir}" "${OBI_DIR}"
  echo "${obi_module_version}" > "${version_file}"
fi

# Linux builds compile OBI's eBPF-enabled paths and require generated files.
if [[ "$(uname -s)" != "Linux" ]]; then
  exit 0
fi

if ! find "${OBI_DIR}" -name "*_bpfel.go" | grep -q .; then
  if ! command -v docker > /dev/null 2>&1; then
    echo "ERROR: docker is required to generate OBI eBPF artifacts on Linux." >&2
    exit 1
  fi
  echo "Generating OBI eBPF artifacts (this may take a few minutes)..."
  make -C "${OBI_DIR}" docker-generate
fi
