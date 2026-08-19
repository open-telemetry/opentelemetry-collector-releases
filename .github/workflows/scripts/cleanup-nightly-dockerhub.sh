#!/bin/bash
set -euo pipefail
# Delete Docker Hub tags belonging to nightly runs older than the cutoff.
#
# Environment:
#   DISTRIBUTIONS     whitespace-separated list of image names
#   MAX_NIGHTLY_RUNS  number of nightly runs to delete per invocation
#   DRY_RUN           "true" to log what would be deleted without deleting it
#   DOCKERHUB_ORG     Docker Hub org owning the repositories
#   DOCKER_USERNAME   Docker Hub user
#   DOCKER_TOKEN      Docker Hub access token

CUTOFF=$(date -u -d '2 weeks ago' +%Y-%m-%dT%H:%M:%SZ)

# Obtain a Docker Hub JWT
HUB_TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "{\"username\": \"${DOCKER_USERNAME}\", \"password\": \"${DOCKER_TOKEN}\"}" \
  "https://hub.docker.com/v2/users/login/" | jq -r '.token')

if [ -z "$HUB_TOKEN" ] || [ "$HUB_TOKEN" = "null" ]; then
  echo "Failed to authenticate with Docker Hub"
  exit 1
fi

for IMAGE in $DISTRIBUTIONS; do
  echo "Processing Docker Hub image: ${DOCKERHUB_ORG}/${IMAGE}"

  # Collect all old nightly tags across pages, then group into runs and cap.
  OLD_TAGS="[]"
  PAGE=1
  while true; do
    RESPONSE=$(curl -s \
      -H "Authorization: Bearer ${HUB_TOKEN}" \
      "https://hub.docker.com/v2/repositories/${DOCKERHUB_ORG}/${IMAGE}/tags/?page_size=100&page=${PAGE}")

    COUNT=$(echo "$RESPONSE" | jq '.results | length')
    if [ "$COUNT" -eq 0 ]; then
      break
    fi

    BATCH=$(echo "$RESPONSE" | jq --arg cutoff "$CUTOFF" '[.results[] | select(
      (.name | test("^v[0-9]+\\.[0-9]+\\.[0-9]+-nightly\\.[0-9]+")) and
      (.last_updated < $cutoff)
    ) | {
      name: .name,
      last_updated: .last_updated,
      run: (.name | capture("nightly\\.(?<ts>[0-9]+)") | .ts)
    }]')
    OLD_TAGS=$(jq -n --argjson a "$OLD_TAGS" --argjson b "$BATCH" '$a + $b')

    NEXT=$(echo "$RESPONSE" | jq -r '.next')
    if [ "$NEXT" = "null" ] || [ -z "$NEXT" ]; then
      break
    fi
    PAGE=$((PAGE + 1))
  done

  while IFS=$'\t' read -r TAG_NAME UPDATED; do
    if [ "$DRY_RUN" = "true" ]; then
      echo "  DRY RUN: would delete Docker Hub tag: ${DOCKERHUB_ORG}/${IMAGE}:${TAG_NAME} (updated: $UPDATED)"
      continue
    fi

    echo "  Deleting Docker Hub tag: ${DOCKERHUB_ORG}/${IMAGE}:${TAG_NAME} (updated: $UPDATED)"
    STATUS=$(curl -s -o /dev/null -w '%{http_code}' -X DELETE \
      -H "Authorization: Bearer ${HUB_TOKEN}" \
      "https://hub.docker.com/v2/repositories/${DOCKERHUB_ORG}/${IMAGE}/tags/${TAG_NAME}/")
    if [ "$STATUS" != "204" ] && [ "$STATUS" != "202" ]; then
      echo "    Failed to delete ${DOCKERHUB_ORG}/${IMAGE}:${TAG_NAME} (HTTP $STATUS)" >&2
      exit 1
    fi
  done < <(echo "$OLD_TAGS" | jq -r --argjson max "$MAX_NIGHTLY_RUNS" '
    (map(.run) | unique | .[:$max]) as $runs
    | map(select(.run | IN($runs[])))
    | sort_by(.run)[]
    | [.name, .last_updated] | @tsv')
done
