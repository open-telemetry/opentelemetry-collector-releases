#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# This script verifies that all components declared in manifest.yaml files are
# defined in the builder-config.yaml from the opentelemetry-collector-contrib
# repository, or ebpf-profiler, ensuring they were built and tested successfully.

set -euo pipefail

BUILDER_CONFIG_URL="https://raw.githubusercontent.com/open-telemetry/opentelemetry-collector-contrib/main/cmd/otelcontribcol/builder-config.yaml"
MANIFEST_DIR="distributions"
CURL_MAX_RETRIES=3
CURL_RETRY_DELAY_SECONDS=1

# Ensure required tools are available
if ! command -v curl &> /dev/null || ! command -v yq &> /dev/null; then
    echo "This script requires 'curl' and 'yq'. Please install them and try again."
    exit 1
fi

CURL_FETCH_ARGS=(
  --fail
  --show-error
  --location
  --silent
  --retry "${CURL_MAX_RETRIES}"
  --retry-delay "${CURL_RETRY_DELAY_SECONDS}"
  --retry-all-errors
)

# Fetch and parse valid components from builder-config.yaml
echo "Fetching builder-config.yaml..."
builder_config_tmp="$(mktemp)"
trap 'rm -f "$builder_config_tmp"' EXIT
if ! http_status="$(
  curl "${CURL_FETCH_ARGS[@]}" \
    --write-out '%{http_code}' \
    --output "$builder_config_tmp" \
    "$BUILDER_CONFIG_URL"
)"; then
  echo "Failed to fetch builder-config.yaml via curl." >&2
  exit 1
fi
if [[ "$http_status" != "200" ]]; then
  echo "Failed to fetch builder-config.yaml: HTTP ${http_status}" >&2
  exit 1
fi

valid_components="$(
  yq -r '
      (
        .extensions[]?.gomod,
        .receivers[]?.gomod,
        .connectors[]?.gomod,
        .processors[]?.gomod,
        .exporters[]?.gomod,
        .providers[]?.gomod
      )
    ' "$builder_config_tmp" \
    | awk '{print $1}' \
    | sort -u
)"
valid_components+='
go.opentelemetry.io/ebpf-profiler'

if [[ -z "$valid_components" ]]; then
  echo "Error: No valid 'gomod' entries found in builder-config.yaml!"
  exit 1
fi

echo "Verifying all manifest.yaml files in '${MANIFEST_DIR}'..."

# We accumulate invalid components here as a multi-line string
invalid_components=""

# Use process substitution to avoid subshell issues
while IFS= read -r manifest_file; do
  echo "Checking $manifest_file"

  # Extract and trim components from the local manifest.yaml
  manifest_components="$(
    yq -r '
      (
        .extensions[]?.gomod,
        .receivers[]?.gomod,
        .connectors[]?.gomod,
        .processors[]?.gomod,
        .exporters[]?.gomod,
        .providers[]?.gomod
      )
    ' "$manifest_file" \
    | awk '{print $1}' \
    | sort -u
  )"

  # Compare each manifest component against the valid list
  while IFS= read -r component; do
    if ! printf '%s\n' "$valid_components" | grep -qxF "$component"; then
      invalid_components="${invalid_components}\n${component}"
    fi
  done <<< "$manifest_components"

done < <(find "$MANIFEST_DIR" -type f -name "manifest.yaml")

if [[ -n "$invalid_components" ]]; then
  echo
  echo "The following components MUST be listed in"
  echo "https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/cmd/otelcontribcol/builder-config.yaml"
  echo "to ensure that they can be built:"
  printf '%b\n' "$invalid_components" | sort -u
  echo
  exit 1
else
  echo "All manifest.yaml components are valid!"
fi
