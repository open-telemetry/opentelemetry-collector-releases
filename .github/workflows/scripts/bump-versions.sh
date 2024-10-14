#!/bin/bash

# This script takes either current beta version, current stable version, or both,
# and optionally next beta version and next stable version, and updates the version
# in the specified files. If next version is not provided, it will infer the next 
# semantic version (e.g. v0.110.0 -> v0.111.0 or v1.16.0 -> v1.17.0) based on the 
# current version(s) passed.

# Function to display usage
usage() {
  echo "Usage: $0 [--commit] [--pull-request] --current-beta-core <current-beta-core> [--current-beta-contrib <current-beta-contrib>] [--current-stable <current-stable>] [--next-beta-core <next-beta-core>] [--next-beta-contrib <next-beta-contrib>] [--next-stable <next-stable>]"
  echo "  --current-beta-core: Current beta version of the core component (e.g., v0.110.0)"
  echo "  --current-beta-contrib: Current beta version of the contrib component (e.g., v0.110.0)"
  echo "  --current-stable: Current stable version of the core component (e.g., v1.16.0)"
  echo "  at least one of the above three arguments is required, or there is no version to update."
  echo
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
while [[ "$#" -gt 0 ]]; do
  case $1 in
    --current-beta-core) current_beta_core="$2"; shift ;;
    --current-beta-contrib) current_beta_contrib="$2"; shift ;;
    --current-stable) current_stable="$2"; shift ;;
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

# Check if at least one of the required arguments is provided
if [ -z "$current_beta_core" ] && [ -z "$current_beta_contrib" ] && [ -z "$current_stable" ]; then
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

# Infer the next stable version if current_stable provided and next version not supplied
if [ -n "$current_stable" ] && [ -z "$next_stable" ]; then
  next_stable=$(bump_version "$current_stable")
fi

# List of files to update
files=(
  "distributions/otelcol-contrib/manifest.yaml"
  "distributions/otelcol/manifest.yaml"
  "distributions/otelcol-k8s/manifest.yaml"
  "distributions/otelcol-otlp/manifest.yaml"
  "Makefile"
)

# Update versions in each file
for file in "${files[@]}"; do
  if [ -f "$file" ]; then
    temp_file=$(mktemp)
    cp "$file" "$temp_file"
    line_number=1
    while IFS= read -r line; do
      if [[ "$line" == *"github.com/open-telemetry/opentelemetry-collector-contrib"* ]]; then
        if [ -n "$current_beta_contrib" ]; then
          sed -i.bak "${line_number}s/$current_beta_contrib/$next_beta_contrib/" "$temp_file"
        fi   
      elif [[ "$line" == *"go.opentelemetry.io/collector"* ]]; then
        if [ -n "$current_beta_core" ]; then
          sed -i.bak "${line_number}s/$current_beta_core/$next_beta_core/" "$temp_file"
        fi
        if [ -n "$current_stable" ]; then
          sed -i.bak "${line_number}s/$current_stable/$next_stable/" "$temp_file"
        fi
      fi
      line_number=$((line_number + 1))
    done < "$file"
    mv "$temp_file" "$file"
    rm "${temp_file}.bak"
  else
    echo "File $file does not exist"
  fi
done

echo "Version update completed."

# Commit changes and draft PR
if [ "$commit_changes" = false ] && [ "$create_pr" = false ]; then
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
  for file in "${files[@]}"; do
    git add "$file"
  done
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

echo "Changes committed and PR created."