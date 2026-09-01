// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"slices"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

var (
	// contrib distro
	contribDist = newDistributionBuilder(contribDistro).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&preBuiltBuildConfig{
				TargetOS:   "aix",
				TargetArch: aixArchs,
				PreBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib-aix_{{ .Target }}/otelcol-contrib",
				},
			},
			&preBuiltBuildConfig{
				TargetOS:   "linux",
				TargetArch: baseArchs,
				PreBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib-linux_{{ .Target }}/otelcol-contrib",
				},
			},
			&preBuiltBuildConfig{
				TargetOS:   "darwin",
				TargetArch: darwinArchs,
				PreBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib-darwin_{{ .Target }}/otelcol-contrib",
				},
			},
			&preBuiltBuildConfig{
				TargetOS:   "windows",
				TargetArch: winArchs,
				PreBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib-windows_{{ .Target }}/otelcol-contrib.exe",
				},
			},
		}
		d.ContainerImages = slices.Concat(
			newContainerImages(d.Name, "linux", baseArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.Name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019", winVersionSHA: "sha256:217694d5470363a77b86c42d439090b7466cd959ad5a15221c556a4707481305"}),
			newContainerImages(d.Name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022", winVersionSHA: "sha256:e4ce8c20390c3785c3cbeef15c579d186b3599d37525c596590cf4508e38d3ff"}),
		)
		d.ContainerImageManifests = slices.Concat(
			newContainerImageManifests(d.Name, "linux", baseArchs, containerImageOptions{}),
		)
	}).withPackagingDefaults().withDefaultConfigIncluded().withVarLibDir("otelcol-contrib", "otelcol-contrib").build()

	// contrib build-only project
	contribBuildOnlyDist = newDistributionBuilder(contribDistro).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&fullBuildConfig{TargetOS: "aix", TargetArch: aixArchs, BuildDir: defaultBuildDir},
			&fullBuildConfig{TargetOS: "linux", TargetArch: baseArchs, BuildDir: defaultBuildDir, ArmVersion: []string{"7"}},
			&fullBuildConfig{TargetOS: "darwin", TargetArch: darwinArchs, BuildDir: defaultBuildDir},
			&fullBuildConfig{TargetOS: "windows", TargetArch: winArchs, BuildDir: defaultBuildDir},
		}
	}).withBinArchive().
		withDefaultMonorepo().
		withDefaultEnv().
		withDefaultPartial().
		withDefaultRelease().
		withNightlyConfig().
		build()
)
