module go.opentelemetry.io/collector/internal/tools

go 1.24.0

require go.opentelemetry.io/build-tools/chloggen v0.29.0

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/cobra v1.10.1 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v0.57.1 // Release failed, use v0.57.2
	v0.57.0 // Release failed, use v0.57.2
)
