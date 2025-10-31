// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import "slices"

const (
	ocbReleaseHeader = "### Images and binaries for collector distributions here: https://github.com/open-telemetry/opentelemetry-collector-releases/releases/tag/{{ .Tag }}"
)

var (
	// OCB binary
	ocbDist = newDistributionBuilder(ocbBinary).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&fullBuildConfig{TargetOS: "linux", TargetArch: ocbArchs, BinaryName: "ocb"},
			&fullBuildConfig{TargetOS: "darwin", TargetArch: darwinArchs, BinaryName: "ocb"},
			&fullBuildConfig{TargetOS: "windows", TargetArch: []string{"amd64"}, BinaryName: "ocb"},
		}
		d.ContainerImages = slices.Concat(
			newContainerImages(d.Name, "linux", ocbArchs, containerImageOptions{binaryRelease: true}),
		)
		d.ContainerImageManifests = slices.Concat(
			newContainerImageManifests(d.Name, "linux", ocbArchs, containerImageOptions{binaryRelease: true}),
		)
		d.LdFlags = "-s -w -X go.opentelemetry.io/collector/cmd/builder/internal.version={{ .Version }}"
	}).withBinaryPackagingDefaults().
		withBinaryMonorepo(".core/cmd/builder").
		withDefaultBinaryRelease(ocbReleaseHeader).
		withNightlyConfig().
		build()
)
