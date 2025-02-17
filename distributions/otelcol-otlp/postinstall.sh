#!/bin/sh

# Copyright The OpenTelemetry Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

if command -v systemctl >/dev/null 2>&1; then
    systemctl daemon-reload
    systemctl enable otelcol-otlp.service
    if [ -f /etc/otelcol-otlp/config.yaml ]; then
        systemctl restart otelcol-otlp.service
    else
        echo "Make sure to configure otelcol-otlp by creating /etc/otelcol-otlp/config.yaml"
    fi
fi
