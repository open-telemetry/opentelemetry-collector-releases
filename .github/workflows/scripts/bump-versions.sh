#!/bin/bash

# This script takes either current beta version, current stable version, or both,
# and optionally next beta version and next stable version, and updates the version
# in the specified files. If next version is not provided, it will infer the next 
# semantic version (e.g. v0.110.0 -> v0.111.0 or v1.16.0 -> v1.17.0) based on the 
# current version(s) passed.

# Function to display usage
usage() {
  echo "Usage: $0 --current_beta <current_beta> [--current_stable <current_stable>] [--next_beta <next_beta>] [--next_stable <next_stable>]"
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

# Parse named arguments
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --current_beta) current_beta="$2"; shift ;;
    --current_stable) current_stable="$2"; shift ;;
    --next_beta) next_beta="$2"; shift ;;
    --next_stable) next_stable="$2"; shift ;;
    *) echo "Unknown parameter passed: $1"; usage ;;
  esac
  shift
done

# Check if at least one of the required arguments is provided
if [ -z "$current_beta" ] && [ -z "$current_stable" ]; then
  usage
fi

# Validate and strip versions
if [ -n "$current_beta" ]; then
  validate_and_strip_version current_beta
fi
if [ -n "$current_stable" ]; then
  validate_and_strip_version current_stable
fi
if [ -n "$next_beta" ]; then
  validate_and_strip_version next_beta
fi
if [ -n "$next_stable" ]; then
  validate_and_strip_version next_stable
fi

# Function to bump the minor version and reset patch version to 0
bump_version() {
  local version=$1
  local major
  major=$(echo "$version" | cut -d. -f1)
  local minor
  minor=$(echo "$version" | cut -d. -f2)
  local new_minor
  new_minor=$((minor + 1))
  echo "$major.$new_minor.0"
}

# Infer the next beta version if not supplied
if  [ -n "$current_beta" ] && [ -z "$next_beta" ]; then
  next_beta=$(bump_version "$current_beta")
fi

# Infer the next stable version if current_stable provided and next version not supplied
if [ -n "$current_stable" ] && [ -z "$next_stable" ]; then
  next_stable=$(bump_version "$current_stable")
fi


# List of files to update
# TODO: Uncomment cmd/builder/builder-config.yaml once PR #671 merged
files=(
  "distributions/otelcol-contrib/manifest.yaml"
  "distributions/otelcol/manifest.yaml"
  "distributions/otelcol-k8s/manifest.yaml"
  "distributions/otelcol-otlp/manifest.yaml"
  "Makefile"
  # "cmd/builder/builder-config-yaml"
)

# Update versions in each file
for file in "${files[@]}"; do
  if [ -f "$file" ]; then
    if [ -n "$current_beta" ]; then
      sed -i.bak "s/$current_beta/$next_beta/g" "$file"
      echo "Updated $file from $current_beta to $next_beta"
    fi
    if [ -n "$current_stable" ]; then
      sed -i.bak "s/$current_stable/$next_stable/g" "$file"
      echo "Updated $file from $current_stable to $next_stable"
    fi
    rm "${file}.bak"
  else
    echo "File $file does not exist"
  fi
done

echo "Version update completed."

Commit changes and draft PR
git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"

# TODO: Once Collector 1.0 is released, we can remove the beta version logic
# for commit and PR creation
if [ -n "$current_beta" ]; then
  branch_name="update-version-${next_beta}"
  git checkout -b "$branch_name"
  git add Makefile \
    distributions/otelcol/manifest.yaml \
    distributions/otelcol-contrib/manifest.yaml \
    distributions/otelcol-k8s/manifest.yaml
  git commit -m "Update version from $current_beta to $next_beta"
  git push -u origin "$branch_name"
  gh pr create --title "[chore] Prepare release $next_beta" \
    --body "This PR updates the version from $current_beta to $next_beta" \
    --base main --head "$branch_name" --draft
else
  branch_name="update-version-${next_stable}"
  git checkout -b "$branch_name"
  git add Makefile \
    distributions/otelcol/manifest.yaml \
    distributions/otelcol-contrib/manifest.yaml \
    distributions/otelcol-k8s/manifest.yaml
  git commit -m "Update version from $current_stable to $next_stable"
  git push -u origin "$branch_name"
  gh pr create --title "[chore] Prepare release $next_stable" \
    --body "This PR updates the version from $current_stable to $next_stable" \
    --base main --head "$branch_name" --draft
fi

echo "Changes committed and PR created."
