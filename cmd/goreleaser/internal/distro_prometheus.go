// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import "slices"

var (
	prometheusDist = newDistributionBuilder(prometheusDistro).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&fullBuildConfig{TargetOS: "linux", TargetArch: []string{"amd64", "arm64"}, BuildDir: defaultBuildDir},
			&fullBuildConfig{TargetOS: "darwin", TargetArch: []string{"arm64"}, BuildDir: defaultBuildDir},
			&fullBuildConfig{TargetOS: "windows", TargetArch: []string{"amd64"}, BuildDir: defaultBuildDir},
		}
		d.ContainerImages = newContainerImages(
			d.Name,
			"linux",
			[]string{"amd64", "arm64"},
			containerImageOptions{},
		)
		d.ContainerImageManifests = slices.Concat(
			newContainerImageManifests(
				d.Name,
				"linux",
				[]string{"amd64", "arm64"},
				containerImageOptions{},
			),
		)
	}).withPackagingDefaults().
		withVarLibDir(prometheusDistro, prometheusDistro).
		build()
)
