#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

podman_cp() {
    local container="$1"
    local src="$2"
    local dest="$3"
    local dest_dir
    dest_dir="$( dirname "$dest" )"

    echo "Copying $src to $container:$dest ..."
    podman exec "$container" mkdir -p "$dest_dir"
    podman cp "$src" "$container":"$dest"
}

install_pkg() {
    local container="$1"
    local pkg_path="$2"
    local pkg_base
    pkg_base=$( basename "$pkg_path" )

    echo "Installing $pkg_base ..."
    podman_cp "$container" "$pkg_path" /tmp/"$pkg_base"
    if [[ "${pkg_base##*.}" = "deb" ]]; then
        podman exec "$container" dpkg -i /tmp/"$pkg_base"
    else
        podman exec "$container" rpm -ivh /tmp/"$pkg_base"
    fi
}

uninstall_pkg() {
    local container="$1"
    local pkg_type="$2"
    local pkg_name="$3"

    echo "Uninstalling $pkg_name ..."
    if [[ "$pkg_type" = "deb" ]]; then
        podman exec "$container" dpkg -r "$pkg_name"
    else
        podman exec "$container" rpm -e "$pkg_name"
    fi
}
