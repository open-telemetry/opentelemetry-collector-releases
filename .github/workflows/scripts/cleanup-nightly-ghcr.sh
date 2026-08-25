#!/bin/bash
set -euo pipefail
# Delete GHCR container versions belonging to nightly runs older than the cutoff.
#
# Environment:
#   DISTRIBUTIONS     whitespace-separated list of image names
#   MAX_NIGHTLY_RUNS  number of nightly runs to delete per invocation
#   DRY_RUN           "true" to log what would be deleted without deleting it
#   ORG               GitHub org owning the packages
#   GHCR_REPO         repository segment of the package name
#   GH_TOKEN          token used by `gh api`

CUTOFF=$(date -u -d '2 weeks ago' +%Y-%m-%dT%H:%M:%SZ)

for IMAGE in $DISTRIBUTIONS; do
  # Images are published under ghcr.io/<org>/<repo>/<image>, so the package name
  # includes the repo segment and its slash has to be percent-encoded for the API.
  PACKAGE="${GHCR_REPO}%2F${IMAGE}"
  echo "Processing GHCR package: ${GHCR_REPO}/${IMAGE}"

  # Collect all old nightly versions across pages, then group into runs and cap.
  OLD_VERSIONS="[]"
  PAGE=1
  while true; do
    VERSIONS=$(gh api \
      -H "Accept: application/vnd.github+json" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      "/orgs/${ORG}/packages/container/${PACKAGE}/versions?per_page=100&page=${PAGE}")

    COUNT=$(echo "$VERSIONS" | jq 'length')
    if [ "$COUNT" -eq 0 ]; then
      break
    fi

    # `all` rather than `any`: a version that also carries a release tag such as
    # `latest` is left alone. `run` records which nightly run produced the version.
    BATCH=$(echo "$VERSIONS" | jq --arg cutoff "$CUTOFF" '[.[] | select(
      (.metadata.container.tags | length > 0) and
      (.metadata.container.tags | all(test("^v[0-9]+\\.[0-9]+\\.[0-9]+-nightly\\.[0-9]+"))) and
      (.updated_at < $cutoff)
    ) | {
      id: .id,
      tags: .metadata.container.tags,
      updated_at: .updated_at,
      run: (.metadata.container.tags[0] | capture("nightly\\.(?<ts>[0-9]+)") | .ts)
    }]')
    OLD_VERSIONS=$(jq -n --argjson a "$OLD_VERSIONS" --argjson b "$BATCH" '$a + $b')

    PAGE=$((PAGE + 1))
  done

  while IFS=$'\t' read -r VERSION_ID TAGS UPDATED; do
    if [ "$DRY_RUN" = "true" ]; then
      echo "  DRY RUN: would delete GHCR version $VERSION_ID (tags: $TAGS, updated: $UPDATED)"
      continue
    fi

    echo "  Deleting GHCR version $VERSION_ID (tags: $TAGS, updated: $UPDATED)"
    gh api \
      --method DELETE \
      -H "Accept: application/vnd.github+json" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      "/orgs/${ORG}/packages/container/${PACKAGE}/versions/${VERSION_ID}"
  done < <(echo "$OLD_VERSIONS" | jq -r --argjson max "$MAX_NIGHTLY_RUNS" '
    (map(.run) | unique | .[:$max]) as $runs
    | map(select(.run | IN($runs[])))
    | sort_by(.run)[]
    | [.id, (.tags | join(", ")), .updated_at] | @tsv')
done
