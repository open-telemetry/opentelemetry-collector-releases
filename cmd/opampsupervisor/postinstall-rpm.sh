#!/bin/sh

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload
    systemctl try-restart opampsupervisor.service
fi
