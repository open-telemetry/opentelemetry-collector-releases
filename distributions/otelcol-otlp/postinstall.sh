#!/bin/sh

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

if command -v systemctl >/dev/null 2>&1; then
    if [ -d /run/systemd/system ]; then
        systemctl daemon-reload
    fi
    systemctl enable otelcol-otlp.service
    if [ -f /etc/otelcol-otlp/config.yaml ]; then
        if [ -d /run/systemd/system ]; then
            systemctl restart otelcol-otlp.service
        fi
    else
        echo "Make sure to configure otelcol-otlp by creating /etc/otelcol-otlp/config.yaml"
    fi
fi
