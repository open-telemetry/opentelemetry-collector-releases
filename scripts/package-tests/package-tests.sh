#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

set -euov pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
REPO_DIR="$( cd "$SCRIPT_DIR/../../../../" && pwd )"
export REPO_DIR
PKG_PATH="${1:-}"
DISTRO="${2:-}"

SERVICE_NAME=$DISTRO
PROCESS_NAME=$DISTRO

# shellcheck source=scripts/package-tests/common.sh
source "$SCRIPT_DIR"/common.sh

if [[ -z "$PKG_PATH" ]]; then
    echo "usage: ${BASH_SOURCE[0]} DEB_OR_RPM_PATH" >&2
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
image_name="axoflow-otel-collector-$pkg_type-test"
container_name="$image_name"
container_exec="podman exec $container_name"

trap 'podman rm -fv $container_name >/dev/null 2>&1 || true' EXIT

podman build -t "$image_name" -f "$SCRIPT_DIR/Dockerfile.test.$pkg_type" "$SCRIPT_DIR"
podman rm -fv "$container_name" >/dev/null 2>&1 || true

# test install
podman run --name "$container_name" -d "$image_name"

# ensure that the system is up and running by checking if systemctl is running
$container_exec systemctl is-system-running --quiet --wait
install_pkg "$container_name" "$PKG_PATH"

# ensure service has started and still running after 5 seconds
sleep 5

echo "Checking $SERVICE_NAME service status ..."
$container_exec systemctl --no-pager status "$SERVICE_NAME"

echo "Checking $PROCESS_NAME process ..."
if [ "$DISTRO" = "axoflow-otel-collector" ]; then
  $container_exec pgrep -a -u axoflow-otel-collector -f "$PROCESS_NAME"
fi

# test uninstall
echo
uninstall_pkg "$container_name" "$pkg_type" "$DISTRO"

echo "Checking $SERVICE_NAME service status after uninstall ..."
if $container_exec systemctl --no-pager status "$SERVICE_NAME"; then
    echo "$SERVICE_NAME service still running after uninstall" >&2
    exit 1
fi
echo "$SERVICE_NAME service successfully stopped after uninstall"

echo "Checking $SERVICE_NAME service existence after uninstall ..."
if $container_exec systemctl list-unit-files --all | grep "$SERVICE_NAME"; then
    echo "$SERVICE_NAME service still exists after uninstall" >&2
    exit 1
fi
echo "$SERVICE_NAME service successfully removed after uninstall"

echo "Checking $PROCESS_NAME process after uninstall ..."
if $container_exec pgrep "$PROCESS_NAME"; then
    echo "$PROCESS_NAME process still running after uninstall"
    exit 1
fi
echo "$PROCESS_NAME process successfully killed after uninstall"
