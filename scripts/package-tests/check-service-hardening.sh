#!/bin/bash

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

service="${1:?usage: check-service-hardening.sh SERVICE}"
for property in NoNewPrivileges LockPersonality RestrictRealtime RestrictSUIDSGID \
    ProtectControlGroups ProtectKernelModules ProtectKernelTunables ProtectHostname; do
    value=$(systemctl show "$service" --property="$property" --value)
    if [[ "$value" != "yes" ]]; then
        echo "$service: expected $property=yes, got $value" >&2
        exit 1
    fi
done

# Check the running process as well as the loaded unit configuration.
pid=$(systemctl show "$service" --property=MainPID --value)
if [[ -z "$pid" || "$pid" == "0" ]]; then
    echo "$service: no running main process" >&2
    exit 1
fi
awk '/^NoNewPrivs:/ { found = 1; if ($2 != 1) exit 1 } END { if (!found) exit 1 }' "/proc/$pid/status"
