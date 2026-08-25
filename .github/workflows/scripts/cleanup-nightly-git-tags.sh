#!/bin/bash
set -euo pipefail
# Delete git tags belonging to nightly runs older than the cutoff.
#
# Environment:
#   MAX_NIGHTLY_RUNS  number of nightly runs to delete per invocation
#   DRY_RUN           "true" to log what would be deleted without deleting it

CUTOFF=$(date -u -d '2 weeks ago' +%s)
git fetch --tags

# Age comes from the timestamp in the tag name, not from the commit it points at:
# nightly tags are cut from whatever main HEAD happens to be, so a quiet period
# would make fresh tags look weeks old. Group by the shared suffix, then take the
# MAX_NIGHTLY_RUNS oldest runs (`sort -u` on `nightly.<YYYYMMDDHHMM>` is
# chronological).
mapfile -t RUNS_TO_DELETE < <(
  git tag -l \
    | grep -oE 'nightly\.[0-9]{12}$' \
    | sort -u \
    | while read -r RUN; do
        STAMP=${RUN#nightly.}
        if ! RUN_DATE=$(date -u -d "${STAMP:0:8} ${STAMP:8:2}:${STAMP:10:2}" +%s 2>/dev/null); then
          continue
        fi
        if [ "$RUN_DATE" -lt "$CUTOFF" ]; then
          echo "$RUN"
        fi
      done \
    | head -n "$MAX_NIGHTLY_RUNS"
)

if [ ${#RUNS_TO_DELETE[@]} -eq 0 ]; then
  echo "No nightly git tags older than the cutoff"
  exit 0
fi

# Anchor each run timestamp to the end of the tag name so the dots are literal and
# a run can never match a longer, unrelated suffix.
PATTERN=$(printf '%s\n' "${RUNS_TO_DELETE[@]}" | sed 's/\./\\./g' | paste -sd'|' -)
mapfile -t TAGS_TO_DELETE < <(git tag -l | grep -E "(${PATTERN})\$")

if [ ${#TAGS_TO_DELETE[@]} -eq 0 ]; then
  echo "No git tags matched the selected runs"
  exit 0
fi

if [ "$DRY_RUN" = "true" ]; then
  printf 'DRY RUN: would delete git tag: %s\n' "${TAGS_TO_DELETE[@]}"
  exit 0
fi

printf 'Deleting git tag: %s\n' "${TAGS_TO_DELETE[@]}"
git push origin --delete "${TAGS_TO_DELETE[@]}"
