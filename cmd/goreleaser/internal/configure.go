// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

// This file is a script which generates the .goreleaser.yaml file for all
// supported OpenTelemetry Collector distributions.
//
// Run it with `make generate-goreleaser`.

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

const (
	armArch            = "arm"
	coreDistro         = "otelcol"
	contribDistro      = "otelcol-contrib"
	k8sDistro          = "otelcol-k8s"
	otlpDistro         = "otelcol-otlp"
	ebpfProfilerDistro = "otelcol-ebpf-profiler"
	dockerHub          = "otel"
	ghcr               = "ghcr.io/open-telemetry/opentelemetry-collector-releases"
	binaryNamePrefix   = "otelcol"
	imageNamePrefix    = "opentelemetry-collector"
)

var (
	baseArchs         = []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}
	winArchs          = []string{"386", "amd64", "arm64"}
	winContainerArchs = []string{"amd64"}
	darwinArchs       = []string{"amd64", "arm64"}
	k8sArchs          = []string{"amd64", "arm64", "ppc64le", "s390x"}
	ebpfProfilerArchs = []string{"amd64"}

	imageRepos = []string{dockerHub, ghcr}

	// otelcol (core) distro
	otelColDist = newDistributionBuilder(coreDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullBuildConfig{targetOS: "linux", targetArch: baseArchs, armVersion: []string{"7"}, ppc64Version: []string{"power8"}},
			&fullBuildConfig{targetOS: "darwin", targetArch: darwinArchs},
			&fullBuildConfig{targetOS: "windows", targetArch: winArchs},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", baseArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
		)
	}).WithPackagingDefaults().WithDefaultConfigIncluded().Build()

	// otlp distro
	otlpDist = newDistributionBuilder(otlpDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullBuildConfig{targetOS: "linux", targetArch: baseArchs, armVersion: []string{"7"}, ppc64Version: []string{"power8"}},
			&fullBuildConfig{targetOS: "darwin", targetArch: darwinArchs},
			&fullBuildConfig{targetOS: "windows", targetArch: winArchs},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", baseArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
		)
	}).WithPackagingDefaults().Build()

	// contrib distro
	contribDist = newDistributionBuilder(contribDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&preBuiltBuildConfig{
				targetOS:   "linux",
				targetArch: baseArchs,
				preBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib-linux_{{ .Target }}/otelcol-contrib",
				},
			},
			&preBuiltBuildConfig{
				targetOS:   "darwin",
				targetArch: darwinArchs,
				preBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib-darwin_{{ .Target }}/otelcol-contrib",
				},
			},
			&preBuiltBuildConfig{
				targetOS:   "windows",
				targetArch: winArchs,
				preBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib-windows_{{ .Target }}/otelcol-contrib.exe",
				},
			},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", baseArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
		)
	}).WithPackagingDefaults().WithDefaultConfigIncluded().Build()

	// contrib build-only project
	contribBuildOnlyDist = newDistributionBuilder(contribDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullBuildConfig{targetOS: "linux", targetArch: baseArchs, armVersion: []string{"7"}},
			&fullBuildConfig{targetOS: "darwin", targetArch: darwinArchs},
			&fullBuildConfig{targetOS: "windows", targetArch: winArchs},
		}
	}).WithBinArchive().Build()

	// k8s distro
	k8sDist = newDistributionBuilder(k8sDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullBuildConfig{targetOS: "linux", targetArch: k8sArchs, ppc64Version: []string{"power8"}},
			&fullBuildConfig{targetOS: "windows", targetArch: winContainerArchs},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", k8sArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", k8sArchs, containerImageOptions{}),
		)
	}).WithDefaultArchives().WithDefaultChecksum().WithDefaultSigns().WithDefaultDockerSigns().WithDefaultSBOMs().WithNightlyConfig().Build()

	// ebpf-profiler distro
	ebpfProfilerDist = newDistributionBuilder(ebpfProfilerDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullBuildConfig{targetOS: "linux", targetArch: ebpfProfilerArchs},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", ebpfProfilerArchs, containerImageOptions{}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", ebpfProfilerArchs, containerImageOptions{}),
		)
		d.enableCgo = true
		d.ldFlags = "-extldflags=-static"
		d.goTags = "osusergo,netgo"
	}).WithDefaultArchives().WithDefaultChecksum().WithDefaultSigns().WithDefaultDockerSigns().WithDefaultSBOMs().Build()
)

type buildConfig interface {
	Build(dist string) config.Build
	OS() string
}

type distributionBuilder struct {
	dist        *distribution
	configFuncs []func(*distribution)
}

func newDistributionBuilder(name string) *distributionBuilder {
	return &distributionBuilder{dist: &distribution{name: name}}
}

func (b *distributionBuilder) WithDefaultArchives() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		builds := make([]string, 0, len(d.buildConfigs))
		for _, build := range d.buildConfigs {
			builds = append(builds, fmt.Sprintf("%s-%s", d.name, build.OS()))
		}
		d.archives = b.newArchives(d.name, builds)
	})
	return b
}

func (b *distributionBuilder) WithBinArchive() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.archives = append(d.archives, config.Archive{
			Formats: []string{"binary"},
		})
	})
	return b
}

func (b *distributionBuilder) newArchives(dist string, builds []string) []config.Archive {
	return []config.Archive{
		{
			ID:           dist,
			NameTemplate: "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}{{ if .Mips }}_{{ .Mips }}{{ end }}",
			Builds:       builds,
		},
	}
}

func (b *distributionBuilder) WithDefaultNfpms() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.nfpms = b.newNfpms(d.name)
	})
	return b
}

func (b *distributionBuilder) newNfpms(dist string) []config.NFPM {
	nfpmContents := []config.NFPMContent{
		{
			Source:      fmt.Sprintf("%s.service", dist),
			Destination: path.Join("/lib", "systemd", "system", fmt.Sprintf("%s.service", dist)),
		},
		{
			Source:      fmt.Sprintf("%s.conf", dist),
			Destination: path.Join("/etc", dist, fmt.Sprintf("%s.conf", dist)),
			Type:        "config|noreplace",
		},
	}
	return []config.NFPM{
		{
			ID:          dist,
			Builds:      []string{dist + "-linux"},
			Formats:     []string{"deb", "rpm"},
			License:     "Apache 2.0",
			Description: fmt.Sprintf("OpenTelemetry Collector - %s", dist),
			Maintainer:  "The OpenTelemetry Collector maintainers <cncf-opentelemetry-maintainers@lists.cncf.io>",
			Overrides: map[string]config.NFPMOverridables{
				"rpm": {
					Dependencies: []string{"/bin/sh"},
				},
			},
			NFPMOverridables: config.NFPMOverridables{
				PackageName: dist,
				Scripts: config.NFPMScripts{
					PreInstall:  "preinstall.sh",
					PostInstall: "postinstall.sh",
					PreRemove:   "preremove.sh",
				},
				Contents: nfpmContents,
			},
		},
	}
}

func (b *distributionBuilder) WithDefaultMSIConfig() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.msiConfig = b.newMSIConfig(d.name)
	})
	return b
}

func (b *distributionBuilder) newMSIConfig(dist string) []config.MSI {
	files := []string{"opentelemetry.ico"}
	return []config.MSI{
		{
			ID:    dist,
			Name:  fmt.Sprintf("%s_{{ .Version }}_{{ .Os }}_{{ .MsiArch }}", dist),
			WXS:   "windows-installer.wxs",
			Files: files,
		},
	}
}

func (b *distributionBuilder) WithDefaultSigns() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.signs = b.signs()
	})
	return b
}

func (b *distributionBuilder) signs() []config.Sign {
	return []config.Sign{
		{
			Artifacts:   "all",
			Signature:   "${artifact}.sig",
			Certificate: "${artifact}.pem",
			Cmd:         "cosign",
			Args: []string{
				"sign-blob",
				"--output-signature",
				"${artifact}.sig",
				"--output-certificate",
				"${artifact}.pem",
				"${artifact}",
			},
		},
	}
}

func (b *distributionBuilder) WithDefaultDockerSigns() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.dockerSigns = b.dockerSigns()
	})
	return b
}

func (b *distributionBuilder) dockerSigns() []config.Sign {
	return []config.Sign{
		{
			Artifacts: "all",
			Args: []string{
				"sign",
				"${artifact}",
			},
		},
	}
}

func (b *distributionBuilder) WithNightlyConfig() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.nightly = b.nightly()
	})
	return b
}

func (b *distributionBuilder) nightly() config.Nightly {
	return config.Nightly{
		VersionTemplate:   "{{ incpatch .Version}}-nightly.{{ .Now.Format \"200601021504\" }}",
		TagName:           "nightly",
		PublishRelease:    true,
		KeepSingleRelease: false,
	}
}

func (b *distributionBuilder) WithDefaultSBOMs() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.sboms = b.sboms()
	})
	return b
}

func (b *distributionBuilder) sboms() []config.SBOM {
	return []config.SBOM{
		{
			ID:        "archive",
			Artifacts: "archive",
		},
		{
			ID:        "package",
			Artifacts: "package",
		},
	}
}

func (b *distributionBuilder) WithDefaultChecksum() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.checksum = config.Checksum{
			NameTemplate: fmt.Sprintf("{{ .ProjectName }}_%v_checksums.txt", d.name),
		}
	})
	return b
}

func (b *distributionBuilder) WithPackagingDefaults() *distributionBuilder {
	return b.WithDefaultArchives().
		WithDefaultChecksum().
		WithDefaultNfpms().
		WithDefaultMSIConfig().
		WithDefaultSigns().
		WithDefaultDockerSigns().
		WithDefaultSBOMs().
		WithNightlyConfig()
}

func (b *distributionBuilder) WithConfigFunc(configFunc func(*distribution)) *distributionBuilder {
	b.configFuncs = append(b.configFuncs, configFunc)
	return b
}

func (b *distributionBuilder) WithDefaultConfigIncluded() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		for i, container := range d.containerImages {
			container.Files = append(container.Files, "config.yaml")
			d.containerImages[i] = container
		}

		for i, nfpm := range d.nfpms {
			nfpm.Contents = append(nfpm.Contents, config.NFPMContent{
				Source:      "config.yaml",
				Destination: path.Join("/etc", d.name, "config.yaml"),
				Type:        "config|noreplace",
			})
			d.nfpms[i] = nfpm
		}

		for i := range d.msiConfig {
			d.msiConfig[i].Files = append(d.msiConfig[i].Files, "config.yaml")
		}
	})
	return b
}

func (b *distributionBuilder) Build() *distribution {
	for _, configFunc := range b.configFuncs {
		configFunc(b.dist)
	}
	return b.dist
}

type distribution struct {
	// Name of the distribution (i.e. otelcol, otelcol-contrib, k8s)
	name string

	buildConfigs            []buildConfig
	archives                []config.Archive
	msiConfig               []config.MSI
	nfpms                   []config.NFPM
	containerImages         []config.Docker
	containerImageManifests []config.DockerManifest
	signs                   []config.Sign
	dockerSigns             []config.Sign
	sboms                   []config.SBOM
	nightly                 config.Nightly
	checksum                config.Checksum
	enableCgo               bool
	ldFlags                 string
	goTags                  string
}

func (d *distribution) BuildProject() config.Project {
	builds := make([]config.Build, 0, len(d.buildConfigs))
	for _, buildConfig := range d.buildConfigs {
		builds = append(builds, buildConfig.Build(d.name))
	}

	ldFlags := "-s -w"
	if d.ldFlags != "" {
		ldFlags = d.ldFlags
	}

	env := []string{
		"COSIGN_YES=true",
		"LD_FLAGS=" + ldFlags,
		"BUILD_FLAGS=-trimpath",
	}
	if d.goTags != "" {
		env = append(env, "GO_TAGS="+d.goTags)
	}
	if !d.enableCgo {
		env = append(env, "CGO_ENABLED=0")
	}

	return config.Project{
		ProjectName: "opentelemetry-collector-releases",
		Release: config.Release{
			ReplaceExistingArtifacts: true,
		},
		Checksum:        d.checksum,
		Env:             env,
		Builds:          builds,
		Archives:        d.archives,
		MSI:             d.msiConfig,
		NFPMs:           d.nfpms,
		Dockers:         d.containerImages,
		DockerManifests: d.containerImageManifests,
		Signs:           d.signs,
		DockerSigns:     d.dockerSigns,
		SBOMs:           d.sboms,
		Nightly:         d.nightly,
		Version:         2,
		Monorepo: config.Monorepo{
			TagPrefix: "v",
		},
		Partial: config.Partial{By: "target"},
	}
}

func newContainerImageManifests(dist, os string, archs []string, opts containerImageOptions) []config.DockerManifest {
	tags := []string{`{{ .Version }}`, "latest"}
	if os == "windows" {
		for i, tag := range tags {
			tags[i] = fmt.Sprintf("%s-%s-%s", tag, os, opts.winVersion)
		}
	}

	var r []config.DockerManifest
	for _, imageRepo := range imageRepos {
		for _, tag := range tags {
			r = append(r, osDockerManifest(imageRepo, tag, dist, os, archs))
		}
	}
	return r
}

type containerImageOptions struct {
	armVersion string
	winVersion string
}

func (o *containerImageOptions) version() string {
	if o.armVersion != "" {
		return o.armVersion
	}
	return o.winVersion
}

func newContainerImages(dist string, targetOS string, targetArchs []string, opts containerImageOptions) []config.Docker {
	images := []config.Docker{}
	for _, targetArch := range targetArchs {
		images = append(images, dockerImageWithOS(dist, targetOS, targetArch, opts))
	}
	return images
}

type fullBuildConfig struct {
	targetOS     string
	targetArch   []string
	armVersion   []string
	ppc64Version []string
}

func (c *fullBuildConfig) Build(dist string) config.Build {
	buildConfig := config.Build{
		ID:     dist + "-" + c.targetOS,
		Dir:    "_build",
		Binary: dist,
		BuildDetails: config.BuildDetails{
			Flags:   []string{"{{ .Env.BUILD_FLAGS }}"},
			Ldflags: []string{"{{ .Env.LD_FLAGS }}"},
		},
		Goos:    []string{c.targetOS},
		Goarch:  c.targetArch,
		Goarm:   c.armVersion,
		Goppc64: c.ppc64Version,
	}
	return buildConfig
}

func (c *fullBuildConfig) OS() string {
	return c.targetOS
}

type preBuiltBuildConfig struct {
	targetOS   string
	targetArch []string
	preBuilt   config.PreBuiltOptions
}

func (c *preBuiltBuildConfig) Build(dist string) config.Build {
	return config.Build{
		ID:       dist + "-" + c.targetOS,
		Builder:  "prebuilt",
		PreBuilt: c.preBuilt,
		Dir:      "_build",
		Binary:   dist,
		Goos:     []string{c.targetOS},
		Goarch:   c.targetArch,
		Goarm:    armVersions(dist),
		Goppc64:  []string{"power8"},
	}
}

func (c *preBuiltBuildConfig) OS() string {
	return c.targetOS
}

func dockerImageWithOS(dist, os, arch string, opts containerImageOptions) config.Docker {
	osArch := osArch{os: os, arch: arch, version: opts.version()}
	var imageTemplates []string
	for _, prefix := range imageRepos {
		imageTemplates = append(
			imageTemplates,
			fmt.Sprintf("%s/%s:{{ .Version }}-%s", prefix, imageName(dist), osArch.imageTag()),
			fmt.Sprintf("%s/%s:latest-%s", prefix, imageName(dist), osArch.imageTag()),
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
	if arch == armArch {
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

type osArch struct {
	os, arch, version string
}

func (o *osArch) buildPlatform() string {
	switch o.os {
	case "linux":
		switch o.arch {
		case armArch:
			return fmt.Sprintf("linux/arm/v%s", o.version)
		}
	case "windows":
		return fmt.Sprintf("windows/%s", o.arch)
	}
	return fmt.Sprintf("linux/%s", o.arch)
}

func (o *osArch) imageTag() string {
	switch o.os {
	case "linux":
		switch o.arch {
		case armArch:
			return fmt.Sprintf("armv%s", o.version)
		}
	case "windows":
		return fmt.Sprintf("windows-%s-%s", o.version, o.arch)
	}
	return o.arch
}

func BuildDist(dist string, onlyBuild bool) config.Project {
	switch dist {
	case coreDistro:
		return otelColDist.BuildProject()
	case otlpDistro:
		return otlpDist.BuildProject()
	case k8sDistro:
		return k8sDist.BuildProject()
	case ebpfProfilerDistro:
		return ebpfProfilerDist.BuildProject()
	case contribDistro:
		if onlyBuild {
			return contribBuildOnlyDist.BuildProject()
		}
		return contribDist.BuildProject()
	default:
		panic("Unknown distribution")
	}
}

func osDockerManifest(prefix, version, dist, os string, archs []string) config.DockerManifest {
	var imageTemplates []string
	for _, arch := range archs {
		switch arch {
		case armArch:
			for _, armVers := range armVersions(dist) {
				dockerArchTag := (&osArch{os: os, arch: arch, version: armVers}).imageTag()
				imageTemplates = append(
					imageTemplates,
					fmt.Sprintf("%s/%s:%s-%s", prefix, imageName(dist), version, dockerArchTag),
				)
			}
		default:
			imageTemplates = append(
				imageTemplates,
				fmt.Sprintf("%s/%s:%s-%s", prefix, imageName(dist), version, arch),
			)
		}
	}

	manifest := config.DockerManifest{
		NameTemplate:   fmt.Sprintf("%s/%s:%s", prefix, imageName(dist), version),
		ImageTemplates: imageTemplates,
	}
	if os == "windows" {
		manifest.SkipPush = "{{ not (eq .Runtime.Goos \"windows\") }}"
	}
	return manifest
}

func armVersions(dist string) []string {
	if dist == k8sDistro {
		return nil
	}
	return []string{"7"}
}

// imageName translates a distribution name to a container image name.
func imageName(dist string) string {
	return strings.Replace(dist, binaryNamePrefix, imageNamePrefix, 1)
}
