// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import "slices"

var (
	// k8s distro
	k8sDist = newDistributionBuilder(k8sDistro).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&fullBuildConfig{TargetOS: "linux", TargetArch: k8sArchs, BuildDir: defaultBuildDir, Ppc64Version: []string{"power8"}},
			&fullBuildConfig{TargetOS: "windows", TargetArch: winContainerArchs, BuildDir: defaultBuildDir},
		}
		d.ContainerImages = slices.Concat(
			newContainerImages(d.Name, "linux", k8sArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.Name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019", winVersionSHA: "sha256:217694d5470363a77b86c42d439090b7466cd959ad5a15221c556a4707481305"}),
			newContainerImages(d.Name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022", winVersionSHA: "sha256:e4ce8c20390c3785c3cbeef15c579d186b3599d37525c596590cf4508e38d3ff"}),
		)
		d.ContainerImageManifests = slices.Concat(
			newContainerImageManifests(d.Name, "linux", k8sArchs, containerImageOptions{}),
		)
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
