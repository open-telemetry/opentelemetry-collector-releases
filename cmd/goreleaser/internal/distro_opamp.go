// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"path"
	"slices"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

const (
	opampReleaseHeader = "### Release of OpAMP supervisor artifacts"
)

var (
	// OpAMP Supervisor binary
	opampDist = newDistributionBuilder(opampBinary).withConfigFunc(func(d *distribution) {
		d.BuildConfigs = []buildConfig{
			&fullBuildConfig{TargetOS: "linux", TargetArch: opAmpArchs, BinaryName: "opampsupervisor"},
			&fullBuildConfig{TargetOS: "darwin", TargetArch: darwinArchs, BinaryName: "opampsupervisor"},
			&fullBuildConfig{TargetOS: "windows", TargetArch: []string{"amd64"}, BinaryName: "opampsupervisor"},
		}
		d.ContainerImages = slices.Concat(
			newContainerImages(d.Name, "linux", opAmpArchs, containerImageOptions{binaryRelease: true}),
		)
		d.ContainerImageManifests = slices.Concat(
			newContainerImageManifests(d.Name, "linux", opAmpArchs, containerImageOptions{binaryRelease: true}),
		)
		d.LdFlags = "-s -w -X github.com/open-telemetry/opentelemetry-collector-contrib/cmd/opampsupervisor/internal.version={{ .Version }}"
	}).withBinaryPackagingDefaults().
		withBinaryMonorepo(".contrib/cmd/opampsupervisor").
		withDefaultBinaryRelease(opampReleaseHeader).
		withDefaultNfpms().
		// This is required because of some non-obvious path/workdir handling in
		// Github Actions specific to the binaries CI.
		withConfigFunc(func(d *distribution) {
			d.Nfpms[0].Contents = append(d.Nfpms[0].Contents, config.NFPMContent{
				Source:      "config.example.yaml",
				Destination: path.Join("/etc", d.Name, "config.example.yaml"),
				Type:        "config|noreplace",
			})
			for i, content := range d.Nfpms[0].Contents {
				content.Source = path.Join("cmd", d.Name, content.Source)
				d.Nfpms[0].Contents[i] = content
			}
			d.Nfpms[0].Scripts.PreInstall = path.Join("cmd", d.Name, d.Nfpms[0].Scripts.PreInstall)
			d.Nfpms[0].Scripts.PostInstall = path.Join("cmd", d.Name, d.Nfpms[0].Scripts.PostInstall)
			d.Nfpms[0].Scripts.PreRemove = path.Join("cmd", d.Name, d.Nfpms[0].Scripts.PreRemove)
		}).
		withNightlyConfig().
		build()
)
