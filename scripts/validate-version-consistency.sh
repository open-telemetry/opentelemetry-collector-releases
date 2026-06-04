#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

# This script validates version consistency across manifest files.
# By default, it only checks that components are internally consistent
# (all core components have the same version, all contrib components have the same version).
#
# Use --check-dist-version to additionally validate that dist.version matches component versions.
# This stricter check is run automatically by the update-version workflow after bumping versions.

set -euo pipefail

MANIFEST_DIR="distributions"
CHECK_DIST_VERSION=false

# Parse arguments
while [[ "$#" -gt 0 ]]; do
    case $1 in
        --check-dist-version) CHECK_DIST_VERSION=true ;;
        *) echo "Unknown parameter: $1"; echo "Usage: $0 [--check-dist-version]"; exit 1 ;;
    esac
    shift
done

check_dependencies() {
    if ! command -v yq &> /dev/null; then
        echo "ERROR: This script requires 'yq'. Please install it and try again."
        exit 1
    fi
}

find_manifest_files() {
    find "$MANIFEST_DIR" -type f -name "manifest.yaml" | sort
}

get_dist_version() {
    local manifest_file="$1"
    yq -r '.dist.version' "$manifest_file"
}

# Check if a module should have its version validated
# Returns 0 (true) if it should be validated, 1 (false) otherwise
should_validate_module_version() {
    local module_path="$1"

    # Validate contrib components
    if [[ "$module_path" == github.com/open-telemetry/opentelemetry-collector-contrib/* ]]; then
        return 0
    fi

    # Validate core collector components, EXCEPT providers (they use different versioning like v1.x)
    if [[ "$module_path" == go.opentelemetry.io/collector/* ]] && \
       [[ "$module_path" != go.opentelemetry.io/collector/confmap/provider/* ]]; then
        return 0
    fi

    # Don't validate:
    # - Core providers (go.opentelemetry.io/collector/confmap/provider/*) - use v1.x versioning
    # - eBPF profiler (go.opentelemetry.io/ebpf-profiler) - has its own versioning
    return 1
}

# Get the module prefix for grouping (core vs contrib)
get_module_prefix() {
    local module_path="$1"

    if [[ "$module_path" == github.com/open-telemetry/opentelemetry-collector-contrib/* ]]; then
        echo "contrib"
    elif [[ "$module_path" == go.opentelemetry.io/collector/* ]]; then
        echo "core"
    else
        echo "other"
    fi
}

# Extract all components from a manifest that should be version-validated
get_validatable_components() {
    local manifest_file="$1"

    yq -r '
      (
        .extensions[]?.gomod,
        .receivers[]?.gomod,
        .connectors[]?.gomod,
        .processors[]?.gomod,
        .exporters[]?.gomod,
        .providers[]?.gomod
      )
    ' "$manifest_file" 2>/dev/null
}

# Check that all distributions have the same dist.version
validate_dist_versions_match() {
    echo "Checking all distributions have the same dist.version..."

    local versions_tmp
    versions_tmp="$(mktemp)"

    while IFS= read -r manifest_file; do
        local version
        version=$(get_dist_version "$manifest_file")
        echo "$version $manifest_file" >> "$versions_tmp"
    done < <(find_manifest_files)

    local unique_versions
    unique_versions=$(awk '{print $1}' "$versions_tmp" | sort -u)
    local version_count
    version_count=$(echo "$unique_versions" | wc -l | tr -d ' ')

    if [[ "$version_count" -gt 1 ]]; then
        echo
        echo "ERROR: Distributions have different dist.version values:"
        echo
        while IFS= read -r line; do
            local version file
            version=$(echo "$line" | awk '{print $1}')
            file=$(echo "$line" | awk '{print $2}')
            echo "  $file: $version"
        done < "$versions_tmp"
        echo
        echo "All distributions must use the same version."
        rm -f "$versions_tmp"
        return 1
    fi

    echo "  All distributions use version: $unique_versions"
    rm -f "$versions_tmp"
    return 0
}

# Check that all components of the same type (core/contrib) have consistent versions
validate_component_internal_consistency() {
    echo
    echo "Checking internal component version consistency..."

    local has_errors=false
    local core_versions_tmp
    local contrib_versions_tmp
    core_versions_tmp="$(mktemp)"
    contrib_versions_tmp="$(mktemp)"

    # Collect all versions by type across all manifests
    while IFS= read -r manifest_file; do
        while IFS= read -r line; do
            [[ -z "$line" ]] && continue

            local module_path version
            module_path=$(echo "$line" | awk '{print $1}')
            version=$(echo "$line" | awk '{print $2}')

            [[ -z "$module_path" || -z "$version" ]] && continue

            if should_validate_module_version "$module_path"; then
                local prefix
                prefix=$(get_module_prefix "$module_path")

                if [[ "$prefix" == "core" ]]; then
                    echo "$version" >> "$core_versions_tmp"
                elif [[ "$prefix" == "contrib" ]]; then
                    echo "$version" >> "$contrib_versions_tmp"
                fi
            fi
        done < <(get_validatable_components "$manifest_file")
    done < <(find_manifest_files)

    # Check core components consistency
    if [[ -s "$core_versions_tmp" ]]; then
        local unique_core_versions
        unique_core_versions=$(sort -u "$core_versions_tmp")
        local core_version_count
        core_version_count=$(echo "$unique_core_versions" | wc -l | tr -d ' ')

        if [[ "$core_version_count" -gt 1 ]]; then
            echo "  ERROR: Core collector components have inconsistent versions:"
            echo "$unique_core_versions" | while read -r ver; do
                echo "    $ver"
            done
            has_errors=true
        else
            echo "  Core components: $unique_core_versions"
        fi
    fi

    # Check contrib components consistency
    if [[ -s "$contrib_versions_tmp" ]]; then
        local unique_contrib_versions
        unique_contrib_versions=$(sort -u "$contrib_versions_tmp")
        local contrib_version_count
        contrib_version_count=$(echo "$unique_contrib_versions" | wc -l | tr -d ' ')

        if [[ "$contrib_version_count" -gt 1 ]]; then
            echo "  ERROR: Contrib components have inconsistent versions:"
            echo "$unique_contrib_versions" | while read -r ver; do
                echo "    $ver"
            done
            has_errors=true
        else
            echo "  Contrib components: $unique_contrib_versions"
        fi
    fi

    rm -f "$core_versions_tmp" "$contrib_versions_tmp"

    if [[ "$has_errors" == "true" ]]; then
        echo
        echo "All core components must use the same version."
        echo "All contrib components must use the same version."
        return 1
    fi

    return 0
}

# Check that components in a manifest match the distribution version
validate_components_match_dist_version() {
    local manifest_file="$1"
    local dist_version="$2"
    local expected_version="v${dist_version}"
    local errors=""

    while IFS= read -r line; do
        [[ -z "$line" ]] && continue

        local module_path version
        module_path=$(echo "$line" | awk '{print $1}')
        version=$(echo "$line" | awk '{print $2}')

        [[ -z "$module_path" || -z "$version" ]] && continue

        if should_validate_module_version "$module_path"; then
            if [[ "$version" != "$expected_version" ]]; then
                errors+="    $module_path: found $version, expected $expected_version\n"
            fi
        fi
    done < <(get_validatable_components "$manifest_file")

    if [[ -n "$errors" ]]; then
        echo "  ERROR in $manifest_file:"
        printf '%b' "$errors"
        return 1
    fi

    return 0
}

# Check all manifests for component version mismatches with dist.version
validate_all_component_versions_match_dist() {
    echo
    echo "Checking components match their distribution version..."

    local has_errors=false

    while IFS= read -r manifest_file; do
        local dist_version
        dist_version=$(get_dist_version "$manifest_file")

        if ! validate_components_match_dist_version "$manifest_file" "$dist_version"; then
            has_errors=true
        else
            echo "  $manifest_file: OK"
        fi
    done < <(find_manifest_files)

    if [[ "$has_errors" == "true" ]]; then
        echo
        echo "Components from opentelemetry-collector-contrib and core collector"
        echo "must use the same version as the distribution (v{dist.version})."
        echo
        echo "Excluded from validation:"
        echo "  - Core providers (go.opentelemetry.io/collector/confmap/provider/*) - use v1.x versioning"
        echo "  - eBPF profiler (go.opentelemetry.io/ebpf-profiler) - has its own versioning"
        return 1
    fi

    return 0
}

main() {
    echo "Validating component version consistency..."
    echo

    check_dependencies

    local has_errors=false

    # Always check that all distributions have the same dist.version
    if ! validate_dist_versions_match; then
        has_errors=true
    fi

    # Always check internal component consistency
    if ! validate_component_internal_consistency; then
        has_errors=true
    fi

    # Only check dist.version matches components if flag is set
    if [[ "$CHECK_DIST_VERSION" == "true" ]]; then
        if ! validate_all_component_versions_match_dist; then
            has_errors=true
        fi
    else
        echo
        echo "Skipping dist.version match validation (use --check-dist-version to enable)"
    fi

    echo
    if [[ "$has_errors" == "true" ]]; then
        echo "Validation FAILED. Please fix the version inconsistencies above."
        exit 1
    else
        echo "All version checks passed!"
    fi
}

main "$@"
