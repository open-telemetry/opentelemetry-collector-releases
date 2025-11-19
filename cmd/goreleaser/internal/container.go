// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"slices"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

const (
	armArchitecture = "arm"
	dockerHubRepo   = "otel"
	ghcrRepo        = "ghcr.io/open-telemetry/opentelemetry-collector-releases"
)

var (
	imageRepositories = []string{dockerHubRepo, ghcrRepo}
)

// containerImageOptions contains options for container image configuration.
type containerImageOptions struct {
	armVersion    string
	winVersion    string
	binaryRelease bool
}

func (o *containerImageOptions) version() string {
	if o.armVersion != "" {
		return o.armVersion
	}
	return o.winVersion
}

type osArchInfo struct {
	os, arch, version string
}

func (o *osArchInfo) buildPlatform() string {
	switch o.os {
	case "linux":
		switch o.arch {
		case armArchitecture:
			return fmt.Sprintf("linux/arm/v%s", o.version)
		}
	case "windows":
		return fmt.Sprintf("windows/%s", o.arch)
	}
	return fmt.Sprintf("linux/%s", o.arch)
}

func (o *osArchInfo) imageTag() string {
	switch o.os {
	case "linux":
		switch o.arch {
		case armArchitecture:
			return fmt.Sprintf("armv%s", o.version)
		}
	case "windows":
		return fmt.Sprintf("windows-%s-%s", o.version, o.arch)
	}
	return o.arch
}

// newContainerImages creates container image configurations.
func newContainerImages(dist string, targetOS string, targetArchs []string, opts containerImageOptions) []config.Docker {
	var images []config.Docker
	for _, targetArch := range targetArchs {
		images = append(images, buildDockerImageWithOS(dist, targetOS, targetArch, opts))
	}
	return images
}

// newContainerImageManifests creates container image manifest configurations.
func newContainerImageManifests(dist, os string, archs []string, opts containerImageOptions) []config.DockerManifest {
	tags := []string{`{{ .Version }}`, "{{ .Env.CONTAINER_IMAGE_EPHEMERAL_TAG }}"}
	if os == "windows" {
		for i, tag := range tags {
			tags[i] = fmt.Sprintf("%s-%s-%s", tag, os, opts.winVersion)
		}
	}

	var r []config.DockerManifest
	for _, imageRepo := range imageRepositories {
		for _, tag := range tags {
			r = append(r, buildOSDockerManifest(imageRepo, tag, dist, os, archs, opts))
		}
	}
	return r
}

func buildDockerImageWithOS(dist, os, arch string, opts containerImageOptions) config.Docker {
	osArch := osArchInfo{os: os, arch: arch, version: opts.version()}
	var imageTemplates []string
	for _, prefix := range imageRepositories {
		imageTemplates = append(
			imageTemplates,
			fmt.Sprintf("%s/%s:{{ .Version }}-%s", prefix, imageName(dist, opts), osArch.imageTag()),
			fmt.Sprintf("%s/%s:{{ .Env.CONTAINER_IMAGE_EPHEMERAL_TAG }}-%s", prefix, imageName(dist, opts), osArch.imageTag()),
		)
	}

	label := func(name, template string) string {
		return fmt.Sprintf("--label=org.opencontainers.image.%s={{%s}}", name, template)
	}
	imageConfig := config.Docker{
		ImageTemplates: imageTemplates,
		Dockerfile:     "Dockerfile",
		Use:            "buildx",
		BuildFlagTemplates: []string{
			"--pull",
			fmt.Sprintf("--platform=%s", osArch.buildPlatform()),
			label("created", ".Date"),
			label("name", ".ProjectName"),
			label("revision", ".FullCommit"),
			label("version", ".Version"),
			label("source", ".GitURL"),
			"--label=org.opencontainers.image.licenses=Apache-2.0",
		},
		Goos:   os,
		Goarch: arch,
	}
	if arch == armArchitecture {
		imageConfig.Goarm = opts.armVersion
	}
	if os == "windows" {
		imageConfig.BuildFlagTemplates = slices.Insert(
			imageConfig.BuildFlagTemplates, 1,
			fmt.Sprintf("--build-arg=WIN_VERSION=%s", opts.winVersion),
		)
		imageConfig.Dockerfile = "Windows.dockerfile"
		imageConfig.Use = "docker"
		imageConfig.SkipBuild = "{{ not (eq .Runtime.Goos \"windows\") }}"
		imageConfig.SkipPush = "{{ not (eq .Runtime.Goos \"windows\") }}"
	}
	return imageConfig
}

func buildOSDockerManifest(prefix, version, dist, os string, archs []string, opts containerImageOptions) config.DockerManifest {
	var imageTemplates []string
	for _, arch := range archs {
		switch arch {
		case armArchitecture:
			for _, armVers := range armVersions(dist) {
				dockerArchTag := (&osArchInfo{os: os, arch: arch, version: armVers}).imageTag()
				imageTemplates = append(
					imageTemplates,
					fmt.Sprintf("%s/%s:%s-%s", prefix, imageName(dist, opts), version, dockerArchTag),
				)
			}
		default:
			imageTemplates = append(
				imageTemplates,
				fmt.Sprintf("%s/%s:%s-%s", prefix, imageName(dist, opts), version, arch),
			)
		}
	}

	manifest := config.DockerManifest{
		NameTemplate:   fmt.Sprintf("%s/%s:%s", prefix, imageName(dist, opts), version),
		ImageTemplates: imageTemplates,
	}
	if os == "windows" {
		manifest.SkipPush = "{{ not (eq .Runtime.Goos \"windows\") }}"
	}
	return manifest
}
