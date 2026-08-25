#!/usr/bin/env bash
#
# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# Determine which ref of the collector core/contrib repository a binary release
# should be built against, and write it to $GITHUB_OUTPUT as COLLECTOR_REF.
#
# The ref is read from the distribution manifests rather than from the tag pushed
# to this repository. Those are two independent things: a releases-repo-only
# re-tag (e.g. v0.115.1, v0.124.1) has no counterpart in core or contrib, so
# reusing this repository's tag as the dependency ref makes the checkout fail
# with "A branch or tag with the name '<tag>' could not be found". The manifests
# always name a core/contrib version that was actually released.
#
# Environment:
#   BINARY         which binary is being released ('builder' or 'opampsupervisor')
#   GITHUB_REF     ref that triggered the release
#   GITHUB_OUTPUT  file the resulting COLLECTOR_REF is appended to

set -euo pipefail

# Nightly builds track the tip of core/contrib instead of a released version.
if [[ "${GITHUB_REF}" == *-nightly* ]]; then
  echo "Nightly release, building against main"
  echo "COLLECTOR_REF=main" >> "${GITHUB_OUTPUT}"
  exit 0
fi

# The patterns are used as awk dynamic regexes, so literal dots are written as
# [.] rather than \. to keep them clear of awk's string escape processing.
case "${BINARY}" in
  builder)
    manifest="distributions/otelcol/manifest.yaml"
    module_pattern='go[.]opentelemetry[.]io/collector/'
    ;;
  opampsupervisor)
    manifest="distributions/otelcol-contrib/manifest.yaml"
    module_pattern='github[.]com/open-telemetry/opentelemetry-collector-contrib/'
    ;;
  *)
    echo "Unknown binary '${BINARY}', cannot determine collector ref" >&2
    exit 1
    ;;
esac

if [ ! -f "${manifest}" ]; then
  echo "Manifest ${manifest} not found" >&2
  exit 1
fi

# Match the first beta (v0) module of the dependency and take its version, e.g.
#   - gomod: go.opentelemetry.io/collector/receiver/nopreceiver v0.159.0
# yields v0.159.0. This is the same extraction bump-versions.sh uses.
collector_ref=$(awk -v pattern="${module_pattern}.* v0" \
  '$0 ~ pattern {print $4; exit}' "${manifest}")

# An empty ref would make the parent pipeline silently fall back to the default
# branch during the checkout step, so a release would be built against main
# without anyone noticing.
if [ -z "${collector_ref}" ]; then
  echo "Could not determine collector ref for '${BINARY}' from ${manifest}" >&2
  exit 1
fi

echo "Building ${BINARY} against collector ref ${collector_ref} (from ${manifest})"
echo "COLLECTOR_REF=${collector_ref}" >> "${GITHUB_OUTPUT}"
