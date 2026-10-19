#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# A distribution is released by several parallel base-release.yaml invocations
# (linux, windows, aix, ...) that all run `goreleaser continue --merge` over the
# same .goreleaser.yaml. Its docker_manifests only reference linux per-arch images,
# so exactly one invocation -- the one whose prepare jobs push those images -- may
# create them. That invocation declares `owns_docker_manifests: true`; the rest skip
# the manifests via the skip_push guard in cmd/goreleaser/internal/container.go.
#
# Getting this wrong fails quietly: with no owner the manifests are silently never
# published, and with more than one the extra invocations race and fail with
# "no such manifest". This script makes both cases a build error instead.

set -euo pipefail

WORKFLOW_DIR=".github/workflows"
BASE_WORKFLOW="$/.github/workflows/base-release.yaml"
DISTRIBUTIONS_DIR="distributions"

if ! command -v yq &> /dev/null; then
    echo "This script requires 'yq'. Please install it and try again."
    exit 1
fi

errors=""
checked=0

for workflow in "${WORKFLOW_DIR}"/release-*.yaml; do
  # Distributions released by this workflow through the reusable base workflow.
  distributions="$(
    yq -r ".jobs[] | select(.uses == \"${BASE_WORKFLOW}\") | .with.distribution" \
      "$workflow" | sort -u
  )"

  [[ -z "$distributions" ]] && continue

  while IFS= read -r distribution; do
    goreleaser_config="${DISTRIBUTIONS_DIR}/${distribution}/.goreleaser.yaml"

    if [[ ! -f "$goreleaser_config" ]]; then
      errors="${errors}\n${workflow}: distribution '${distribution}' has no ${goreleaser_config}"
      continue
    fi

    # Distributions without docker_manifests have no manifests to own.
    if [[ "$(yq -r '.docker_manifests // "" | length' "$goreleaser_config")" == "0" ]]; then
      expected_owners=0
    else
      expected_owners=1
    fi

    owners="$(
      yq -r ".jobs | to_entries | .[]
             | select(.value.uses == \"${BASE_WORKFLOW}\")
             | select(.value.with.distribution == \"${distribution}\")
             | select(.value.with.owns_docker_manifests == true)
             | .key" "$workflow"
    )"
    owner_count="$(grep -c . <<< "$owners" || true)"
    checked=$((checked + 1))

    if [[ "$owner_count" -ne "$expected_owners" ]]; then
      if [[ "$expected_owners" -eq 0 ]]; then
        errors="${errors}\n${workflow}: '${distribution}' defines no docker_manifests, but ${owner_count} job(s) set owns_docker_manifests: ${owners//$'\n'/, }"
      else
        errors="${errors}\n${workflow}: '${distribution}' needs exactly 1 job with owns_docker_manifests: true, found ${owner_count}${owners:+ (${owners//$'\n'/, })}"
      fi
      continue
    fi

    # The owner must be the invocation that actually pushes the linux images the
    # manifests reference.
    if [[ "$expected_owners" -eq 1 ]]; then
      owner_goos="$(yq -r ".jobs.\"${owners}\".with.goos" "$workflow")"
      if [[ "$owner_goos" != *linux* ]]; then
        errors="${errors}\n${workflow}: job '${owners}' owns the manifests for '${distribution}' but its goos (${owner_goos}) excludes linux, so it never pushes the images they reference"
      fi
    fi
  done <<< "$distributions"
done

if [[ -n "$errors" ]]; then
  echo "Docker manifest ownership check failed:"
  printf '%b\n' "$errors"
  exit 1
fi

echo "Docker manifest ownership is valid (${checked} distribution/workflow pairs checked)."
