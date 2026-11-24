#!/bin/sh

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

if [ "$1" != "1" ]; then
    if command -v systemctl >/dev/null 2>&1; then
        if [ -d /run/systemd/system ]; then
            systemctl stop opampsupervisor.service
        fi
        systemctl disable opampsupervisor.service
    fi
fi
