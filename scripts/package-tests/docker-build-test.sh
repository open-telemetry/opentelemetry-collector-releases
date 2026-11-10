#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

set -euov pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PKG_PATH="${1:-}"
DISTRO="${2:-}"

# Global array for trap functions
TRAP_FUNCS=()

# Helper to add functions to the trap list
add_trap_func() {
	TRAP_FUNCS+=("$1")
}

# Wrapper that calls each trap'ed function.
# It uses LIFO order, like Go's `defer`.
run_traps() {
	if [ "${#TRAP_FUNCS[@]}" -gt 0 ]; then
		for ((i = ${#TRAP_FUNCS[@]} - 1; i >= 0; i--)); do
			"${TRAP_FUNCS[i]}"
		done
	fi
}
trap 'run_traps' EXIT


if [[ -z "$PKG_PATH" || -z "$DISTRO" ]]; then
    echo "usage: ${BASH_SOURCE[0]} DEB_OR_RPM_PATH DISTRO" >&2
    exit 1
fi

if [[ ! -f "$PKG_PATH" ]]; then
    echo "$PKG_PATH not found!" >&2
    exit 1
fi


pkg_base="$( basename "$PKG_PATH" )"
pkg_type="${pkg_base##*.}"
if [[ ! "$pkg_type" =~ ^(deb|rpm)$ ]]; then
    echo "$PKG_PATH not supported!" >&2
    exit 1
fi
image_name="otelcontribcol-$pkg_type-build-test"

# Use a temporary directory containing only the package as build context, so
# that the directory the package is in (e.g. the repo root or a dist folder
# with packages for all architectures) isn't sent to the builder.
build_ctx="$( mktemp -d )"
build_ctx_cleanup() {
	rm -rf "$build_ctx"
}
add_trap_func build_ctx_cleanup

podman_cleanup() {
	podman rmi -f "$image_name" >/dev/null 2>&1 || true
}
add_trap_func podman_cleanup

cp "$PKG_PATH" "$build_ctx/"

# test install, reinstall and uninstall during an image build, where systemd is not running
podman build \
    --build-arg PKG="$pkg_base" \
    --build-arg DISTRO="$DISTRO" \
    -t "$image_name" \
    -f "$SCRIPT_DIR/Dockerfile.build-test.$pkg_type" \
    "$build_ctx"

echo "$pkg_base successfully installed, reinstalled and removed during image build"
