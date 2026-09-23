// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import "slices"

var (
	// otelcol (core) distro
	otelColDist = newDistributionBuilder(coreDistro).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&fullBuildConfig{TargetOS: "aix", TargetArch: aixArchs, BuildDir: defaultBuildDir},
			&fullBuildConfig{TargetOS: "linux", TargetArch: baseArchs, BuildDir: defaultBuildDir, ArmVersion: []string{"7"}, Ppc64Version: []string{"power8"}},
			&fullBuildConfig{TargetOS: "darwin", TargetArch: darwinArchs, BuildDir: defaultBuildDir},
			&fullBuildConfig{TargetOS: "solaris", TargetArch: solarisArchs, BuildDir: defaultBuildDir},
			&fullBuildConfig{TargetOS: "windows", TargetArch: winArchs, BuildDir: defaultBuildDir},
		}
		d.ContainerImages = slices.Concat(
			newContainerImages(d.Name, "linux", baseArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.Name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019", winVersionSHA: "sha256:217694d5470363a77b86c42d439090b7466cd959ad5a15221c556a4707481305"}),
			newContainerImages(d.Name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022", winVersionSHA: "sha256:e4ce8c20390c3785c3cbeef15c579d186b3599d37525c596590cf4508e38d3ff"}),
		)
		d.ContainerImageManifests = slices.Concat(
			newContainerImageManifests(d.Name, "linux", baseArchs, containerImageOptions{}),
		)
	}).withPackagingDefaults().withDefaultConfigIncluded().withVarLibDir("otel", "otel").build()
)
