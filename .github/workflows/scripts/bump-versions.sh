#!/bin/bash
set -e
# This script reads current versions and takes optional next versions, and updates the
# version in the specified files. If next version is not provided, it will infer the 
# next semantic version (e.g. v0.110.0 -> v0.111.0 or v1.16.0 -> v1.17.0) based on the 
# current version(s) read in.

# List of files to update
manifest_files=(
  "distributions/otelcol-contrib/manifest.yaml"
  "distributions/otelcol/manifest.yaml"
  "distributions/otelcol-k8s/manifest.yaml"
  "distributions/otelcol-otlp/manifest.yaml"
)

# Function to display usage
usage() {
  echo "Usage: $0 [--commit] [--pull-request] [--next-beta-core <next-beta-core>] [--next-beta-contrib <next-beta-contrib>] [--next-stable <next-stable>]"
  echo "  --next-beta-core: Next beta version of the core component (e.g., v0.111.0)"
  echo "  --next-beta-contrib: Next beta version of the contrib component (e.g., v0.111.0)"
  echo "  --next-stable: Next stable version of the core component (e.g., v1.17.0)"
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
current_beta_core=$(awk '/^.*go\.opentelemetry\.io\/collector\/.* v0/ {print $4; exit}' distributions/otelcol/manifest.yaml)
current_beta_contrib=$(awk '/^.*github\.com\/open-telemetry\/opentelemetry-collector-contrib\/.* v0/ {print $4; exit}' distributions/otelcol-contrib/manifest.yaml)
current_stable=$(awk '/^.*go\.opentelemetry\.io\/collector\/.* v1/ {print $4; exit}' distributions/otelcol/manifest.yaml)
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --next-beta-core) next_beta_core="$2"; shift ;;
    --next-beta-contrib) next_beta_contrib="$2"; shift ;;
    --next-stable) next_stable="$2"; shift ;;
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
if [ -n "$current_beta_core" ]; then
  validate_and_strip_version current_beta_core
fi
if [ -n "$current_beta_contrib" ]; then
  validate_and_strip_version current_beta_contrib
fi
if [ -n "$current_stable" ]; then
  validate_and_strip_version current_stable
fi
if [ -n "$next_beta_core" ]; then
  validate_and_strip_version next_beta_core
fi
if [ -n "$next_beta_contrib" ]; then
  validate_and_strip_version next_beta_contrib
fi
if [ -n "$next_stable" ]; then
  validate_and_strip_version next_stable
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
if  [ -n "$current_beta_core" ] && [ -z "$next_beta_core" ]; then
  next_beta_core=$(bump_version "$current_beta_core")
fi
if  [ -n "$current_beta_contrib" ] && [ -z "$next_beta_contrib" ]; then
  next_beta_contrib=$(bump_version "$current_beta_contrib")
fi

# Determine the maximum of next_beta_core and next_beta_contrib
next_distribution_version=$(max_version "$next_beta_core" "$next_beta_contrib")
validate_and_strip_version next_distribution_version

# Infer the next stable version if current_stable provided and next version not supplied
if [ -n "$current_stable" ] && [ -z "$next_stable" ]; then
  next_stable=$(bump_version "$current_stable")
fi

# add escape characters to the current versions to work with sed
escaped_current_beta_core=${current_beta_core//./\\.}
escaped_current_beta_contrib=${current_beta_contrib//./\\.}
escaped_current_stable=${current_stable//./\\.}

# Determine the OS and set the sed command accordingly
if [[ "$OSTYPE" == "darwin"* ]]; then
  # macOS
  sed_command="sed -i ''"
else
  # Linux
  sed_command="sed -i"
fi

# Update versions in each manifest file
echo "Updating core beta version from $current_beta_core to $next_beta_core,"
echo "core stable version from $current_stable to $next_stable,"
echo "contrib beta version from $current_beta_contrib to $next_beta_contrib,"
echo "and distribution version to $next_distribution_version"
for file in "${manifest_files[@]}"; do
  if [ -f "$file" ]; then
    $sed_command "s/\(^.*go\.opentelemetry\.io\/collector\/.*\) v$escaped_current_beta_core/\1 v$next_beta_core/" "$file"
    $sed_command "s/\(^.*github\.com\/open-telemetry\/opentelemetry-collector-contrib\/.*\) v$escaped_current_beta_contrib/\1 v$next_beta_contrib/" "$file"
    $sed_command "s/\(^.*go\.opentelemetry\.io\/collector\/.*\) v$escaped_current_stable/\1 v$next_stable/" "$file"
    $sed_command "s/version: .*/version: $next_distribution_version/" "$file"
  else
    echo "File $file does not exist"
  fi
done

# Update Makefile OCB version
$sed_command "s/OTELCOL_BUILDER_VERSION ?= $escaped_current_beta_core/OTELCOL_BUILDER_VERSION ?= $next_beta_core/" Makefile

echo "Version update completed."

# Make a new changelog update
make chlog-update VERSION="v$next_distribution_version"

# Commit changes and draft PR
if [ "$commit_changes" = false ]; then
  echo "Changes not committed and PR not created."
  exit 0
fi

git config --global user.name "github-actions[bot]"
git config --global user.email "github-actions[bot]@users.noreply.github.com"

commit_changes() {
  local current_version=$1
  local next_version=$2
  shift 2
  local branch_name="update-version-${next_version}"

  git checkout -b "$branch_name"
  git add .
  git commit -m "Update version from $current_version to $next_version"
  git push -u origin "$branch_name"
}

create_pr() {
  local current_version=$1
  local next_version=$2
  shift 2
  local branch_name="update-version-${next_version}"

  gh pr create --title "[chore] Prepare release $next_version" \
    --body "This PR updates the version from $current_version to $next_version" \
    --base main --head "$branch_name" --draft  
}

# TODO: Once Collector 1.0 is released, we can consider removing the
# beta version check for commit and PR creation
if [ -n "$current_beta_core" ]; then
  if [ "$commit_changes" = true ]; then
    commit_changes "$current_beta_core" "$next_beta_core"
  fi
  if [ "$create_pr" = true ]; then
    create_pr "$current_beta_core" "$next_beta_core"
  fi
elif [ -n "$current_beta_contrib" ]; then
  if [ "$commit_changes" = true ]; then
    commit_changes "$current_beta_contrib" "$next_beta_contrib"
  fi
  if [ "$create_pr" = true ]; then
    create_pr "$current_beta_contrib" "$next_beta_contrib"
  fi
else
  if [ "$commit_changes" = true ]; then
    commit_changes "$current_stable" "$next_stable"
  fi
  if [ "$create_pr" = true ]; then
    create_pr "$current_stable" "$next_stable"
  fi
fi

echo "Changes committed and PR created:"
gh pr view
