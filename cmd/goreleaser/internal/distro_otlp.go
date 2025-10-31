// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import "slices"

var (
	// otlp distro
	otlpDist = newDistributionBuilder(otlpDistro).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&fullBuildConfig{TargetOS: "linux", TargetArch: baseArchs, BuildDir: defaultBuildDir, ArmVersion: []string{"7"}, Ppc64Version: []string{"power8"}},
			&fullBuildConfig{TargetOS: "darwin", TargetArch: darwinArchs, BuildDir: defaultBuildDir},
			&fullBuildConfig{TargetOS: "windows", TargetArch: winArchs, BuildDir: defaultBuildDir},
		}
		d.ContainerImages = slices.Concat(
			newContainerImages(d.Name, "linux", baseArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.Name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.Name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.ContainerImageManifests = slices.Concat(
			newContainerImageManifests(d.Name, "linux", baseArchs, containerImageOptions{}),
		)
	}).withPackagingDefaults().build()
)
