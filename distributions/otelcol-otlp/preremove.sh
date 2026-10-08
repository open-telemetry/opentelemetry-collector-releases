#!/bin/sh

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

if [ "$1" != "1" ]; then
    if command -v systemctl >/dev/null 2>&1; then
        debsystemctl=$(command -v deb-systemd-invoke || echo systemctl)
        "$debsystemctl" stop otelcol-otlp.service
        systemctl disable otelcol-otlp.service
    fi
fi
