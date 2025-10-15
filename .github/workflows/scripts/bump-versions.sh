#!/bin/bash
set -e
# This script reads next versions, and updates the version in the specified files.
# It will infer the next semantic version (e.g. v0.110.0 -> v0.111.0 or v1.16.0 -> v1.17.0)
# based on the version(s) read in.

# List of files to update
manifest_files=(
  "distributions/otelcol-contrib/manifest.yaml"
  "distributions/otelcol/manifest.yaml"
  "distributions/otelcol-k8s/manifest.yaml"
  "distributions/otelcol-otlp/manifest.yaml"
  "distributions/otelcol-ebpf-profiler/manifest.yaml"
)

# Function to display usage
usage() {
  echo "Usage: $0 [--commit] [--pull-request]"
  echo
  echo "  --commit: Commit the changes to a new branch"
  echo "  --pull-request: Push the changes to the repo and create a draft PR (requires --commit)"
  exit 1
}

# Function to validate semantic version and strip leading 'v'
validate_and_strip_version() {
  local var_name=$1
  local version=${!var_name}
  # Strip leading 'v' if present
  version=${version#v}
  if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Invalid version: $version. Must be a semantic version (e.g., 1.2.3)."
    exit 1
  fi
  eval "$var_name='$version'"
}
commit_changes=false
create_pr=false
# Parse named arguments
next_beta_core=$(awk '/^.*go\.opentelemetry\.io\/collector\/.* v0/ {print $4; exit}' distributions/otelcol/manifest.yaml)
next_beta_contrib=$(awk '/^.*github\.com\/open-telemetry\/opentelemetry-collector-contrib\/.* v0/ {print $4; exit}' distributions/otelcol-contrib/manifest.yaml)
next_stable_core=$(awk '/^.*go\.opentelemetry\.io\/collector\/.* v1/ {print $4; exit}' distributions/otelcol/manifest.yaml)

while [[ "$#" -gt 0 ]]; do
  case $1 in
    --commit) commit_changes=true ;;
    --pull-request) create_pr=true ;;
    *) echo "Unknown parameter passed: $1"; usage ;;
  esac
  shift
done

# Check if --pull-request is passed without --commit
if [ "$create_pr" = true ] && [ "$commit_changes" = false ]; then
  echo "--pull-request requires --commit"
  usage
fi

# Validate and strip versions
if [ -n "$next_beta_core" ]; then
  validate_and_strip_version next_beta_core
fi
if [ -n "$next_beta_contrib" ]; then
  validate_and_strip_version next_beta_contrib
fi
if [ -n "$next_stable_core" ]; then
  validate_and_strip_version next_stable_core
fi

# Function to compare two semantic versions and return the maximum
max_version() {
  local version1=$1
  local version2=$2

  # Strip leading 'v' if present
  version1=${version1#v}
  version2=${version2#v}

  # Split versions into components
  IFS='.' read -r -a ver1 <<< "$version1"
  IFS='.' read -r -a ver2 <<< "$version2"

  # Compare major, minor, and patch versions
  for i in {0..2}; do
    if [[ ${ver1[i]} -gt ${ver2[i]} ]]; then
      echo "$version1"
      return
    elif [[ ${ver1[i]} -lt ${ver2[i]} ]]; then
      echo "$version2"
      return
    fi
  done

  # If versions are equal, return either
  echo "$version1"
}

# Determine the maximum of next_beta_core and next_beta_contrib
next_distribution_version=$(max_version "$next_beta_core" "$next_beta_contrib")
validate_and_strip_version next_distribution_version

# Update versions in each manifest file
echo "Making the following updates:"
echo "  - core beta module set to $next_beta_core"
echo "  - core stable module set to $next_stable_core"
echo "  - contrib beta module set to $next_beta_contrib"
echo "  - distribution version to $next_distribution_version"
for file in "${manifest_files[@]}"; do
  if [ -f "$file" ]; then
    sed "s/version: .*/version: $next_distribution_version/" "$file" > "$file.tmp"
    mv "$file.tmp" "$file"
  else
    echo "File $file does not exist"
  fi
done

echo "Version update completed."

# Make a new changelog update
make chlog-update VERSION="v$next_distribution_version"

# Commit changes and draft PR
if [ "$commit_changes" = false ]; then
  echo "Changes not committed and PR not created."
  exit 0
fi

commit_changes() {
  local next_version=$1
  shift 1
  local branch_name="update-version-${next_version}"

  git checkout -b "$branch_name"
  git add .
  git commit -m "Update version to $next_version"
  git push -u origin "$branch_name"
}

create_pr() {
  local next_version=$1
  shift 1
  local branch_name="update-version-${next_version}"

  gh pr create --title "[chore] Prepare release $next_version" \
    --body "This PR updates the version to $next_version" \
    --base main --head "$branch_name"
}

# TODO: Once Collector 1.0 is released, we can consider removing the
# beta version check for commit and PR creation
if [ -n "$next_beta_core" ]; then
  if [ "$commit_changes" = true ]; then
    commit_changes "$next_beta_core"
  fi
  if [ "$create_pr" = true ]; then
    create_pr "$next_beta_core"
  fi
elif [ -n "$next_beta_contrib" ]; then
  if [ "$commit_changes" = true ]; then
    commit_changes "$next_beta_contrib"
  fi
  if [ "$create_pr" = true ]; then
    create_pr "$next_beta_contrib"
  fi
else
  if [ "$commit_changes" = true ]; then
    commit_changes "$next_stable_core"
  fi
  if [ "$create_pr" = true ]; then
    create_pr "$next_stable_core"
  fi
fi

echo "Changes committed and PR created:"
gh pr view
