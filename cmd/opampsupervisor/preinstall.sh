#!/bin/sh

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

getent passwd opampsupervisor >/dev/null || useradd --system --user-group --no-create-home --shell /sbin/nologin opampsupervisor
