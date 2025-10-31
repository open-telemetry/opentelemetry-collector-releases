// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import "slices"

var (
	ebpfProfilerArchs = []string{"amd64", "arm64"}
)

var (
	// ebpf-profiler distro
	ebpfProfilerDist = newDistributionBuilder(ebpfProfilerDistro).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&fullBuildConfig{TargetOS: "linux", TargetArch: ebpfProfilerArchs, BuildDir: defaultBuildDir},
		}
		d.ContainerImages = slices.Concat(
			newContainerImages(d.Name, "linux", ebpfProfilerArchs, containerImageOptions{}),
		)
		d.ContainerImageManifests = slices.Concat(
			newContainerImageManifests(d.Name, "linux", ebpfProfilerArchs, containerImageOptions{}),
		)
		d.Env = append(d.Env, "TARGET_ARCH={{ .Runtime.Goarch }}")
		d.LdFlags = "-extldflags=-static"
		d.GoTags = "osusergo,netgo"
	}).withDefaultArchives().
		withDefaultChecksum().
		withDefaultSigns().
		withDefaultDockerSigns().
		withDefaultSBOMs().
		withDefaultMonorepo().
		withDefaultEnv().
		withDefaultPartial().
		withDefaultRelease().
		withNightlyConfig().
		withDefaultSnapshot().
		build()
)
