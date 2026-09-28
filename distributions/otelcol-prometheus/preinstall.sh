#!/bin/sh

# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0

getent passwd otelcol-prometheus >/dev/null || useradd --system --user-group --no-create-home --shell /sbin/nologin otelcol-prometheus
