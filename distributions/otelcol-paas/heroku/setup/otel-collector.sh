#!/bin/bash

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

if [ "$DYNOTYPE" == "run" ]; then
    exit 0
fi

# Set configuration file
export CONFIG_DIR="$HOME/.otel"

if [[ -z "$CONFIG" ]]; then
    export CONFIG="${$CONFIG_DIR/config.yaml}"
fi

# Set log file
if [[ -z "$LOG_FILE" ]]; then
    export LOG_FILE=/dev/stdout
else
    mkdir -p $(dirname $LOG_FILE)
fi

$CONFIG_DIR/otelcol-paas-linux --config $CONFIG > $LOG_FILE 2>&1 &