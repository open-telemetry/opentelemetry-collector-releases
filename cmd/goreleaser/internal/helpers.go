// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"strings"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

func BuildDistribution(dist string, onlyBuild bool) config.Project {
	switch dist {
	case coreDistro:
		return otelColDist.buildProject()
	case otlpDistro:
		return otlpDist.buildProject()
	case k8sDistro:
		return k8sDist.buildProject()
	case ebpfProfilerDistro:
		return ebpfProfilerDist.buildProject()
	case contribDistro:
		if onlyBuild {
			return contribBuildOnlyDist.buildProject()
		}
		return contribDist.buildProject()
	case ocbBinary:
		return ocbDist.buildProject()
	case opampBinary:
		return opampDist.buildProject()
	default:
		panic("Unknown distribution")
	}
}

func armVersions(dist string) []string {
	if dist == k8sDistro {
		return nil
	}
	return []string{"7"}
}

// imageName translates a distribution name to a container image name.
func imageName(dist string, opts containerImageOptions) string {
	if opts.binaryRelease {
		return imageNamePrefix + "-" + dist
	}
	return strings.Replace(dist, binaryNamePrefix, imageNamePrefix, 1)
}
