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

	// manifestOwnerGuard skips a docker manifest unless the release job declares it
	// owns the per-arch images the manifest references. Kept as a template rather
	// than a plain env lookup because goreleaser renders with missingkey=error, so
	// `.Env.OWNS_DOCKER_MANIFESTS` would fail outright when the variable is unset.
	manifestOwnerGuard = `{{ not (isEnvSet "OWNS_DOCKER_MANIFESTS") }}`
)

var (
	imageRepositories = []string{dockerHubRepo, ghcrRepo}
)

// containerImageOptions contains options for container image configuration.
type containerImageOptions struct {
	armVersion    string
	winVersion    string
	winVersionSHA string
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
			fmt.Sprintf("--build-arg=WIN_VERSION_SHA=%s", opts.winVersionSHA),
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
	switch {
	case os == "windows":
		manifest.SkipPush = "{{ not (eq .Runtime.Goos \"windows\") }}"
	case !opts.binaryRelease:
		// Distributions are released by several parallel base-release.yaml invocations
		// (linux, windows, aix, ...), and each one runs `goreleaser continue --merge`
		// over this same config. Without a guard every invocation would try to create
		// these manifests, but the per-arch images they reference are only pushed by
		// the linux invocation -- so whoever gets there first fails with
		// "no such manifest". Only the invocation that owns those images sets
		// OWNS_DOCKER_MANIFESTS, so exactly one of them does the work.
		//
		// Binary releases (ocb, opampsupervisor) are excluded: they run as a single
		// goreleaser job with no split/merge, so there is nothing to race against and
		// nothing sets the variable for them.
		manifest.SkipPush = manifestOwnerGuard
	}
	return manifest
}
