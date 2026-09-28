#!/bin/sh

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

if command -v systemctl >/dev/null 2>&1; then
    if [ -d /run/systemd/system ]; then
        systemctl daemon-reload
    fi
    systemctl enable otelcol-prometheus.service
    if [ -f /etc/otelcol-prometheus/config.yaml ]; then
        if [ -d /run/systemd/system ]; then
            systemctl restart otelcol-prometheus.service
        fi
    else
        echo "Make sure to configure otelcol-prometheus by creating /etc/otelcol-prometheus/config.yaml"
    fi
fi
