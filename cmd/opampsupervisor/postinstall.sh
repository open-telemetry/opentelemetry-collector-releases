#!/bin/sh

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

if command -v systemctl >/dev/null 2>&1; then
    if [ -d /run/systemd/system ]; then
        systemctl daemon-reload
    fi
    systemctl enable opampsupervisor.service
    if [ -f /etc/opampsupervisor/config.yaml ]; then
        if [ -d /run/systemd/system ]; then
            systemctl restart opampsupervisor.service
        fi
    fi
fi
