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
	ArmArch          = "arm"
	CoreDistro       = "otelcol"
	ContribDistro    = "otelcol-contrib"
	K8sDistro        = "otelcol-k8s"
	OTLPDistro       = "otelcol-otlp"
	DockerHub        = "otel"
	GHCR             = "ghcr.io/open-telemetry/opentelemetry-collector-releases"
	BinaryNamePrefix = "otelcol"
	ImageNamePrefix  = "opentelemetry-collector"
)

var (
	baseArchs         = []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}
	winArchs          = []string{"386", "amd64", "arm64", "ppc64le"}
	winContainerArchs = []string{"amd64", "arm64"}
	darwinArchs       = []string{"amd64", "arm64"}

	imageRepos         = []string{DockerHub, GHCR}
	defaultConfigDists = map[string]bool{CoreDistro: true, ContribDistro: true}

	// otelcol (core) distro
	otelColDist = newDistributionBuilder(CoreDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullDistBuildConfig{targetOS: "linux", targetArch: baseArchs, armVersion: []string{"7"}},
			&fullDistBuildConfig{targetOS: "darwin", targetArch: darwinArchs},
			&fullDistBuildConfig{targetOS: "windows", targetArch: winArchs},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", baseArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
			newContainerImageManifests(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
	}).WithPackagingDefaults().Build()

	// otlp distro
	otlpDist = newDistributionBuilder(OTLPDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullDistBuildConfig{targetOS: "linux", targetArch: baseArchs, armVersion: []string{"7"}},
			&fullDistBuildConfig{targetOS: "darwin", targetArch: darwinArchs},
			&fullDistBuildConfig{targetOS: "windows", targetArch: winArchs},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
			newContainerImageManifests(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
	}).WithPackagingDefaults().Build()

	// contrib distro
	contribDist = newDistributionBuilder(ContribDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&preBuiltBuildConfig{
				targetOS:   "linux",
				targetArch: baseArchs,
				preBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib_{{ .Target }}/otelcol-contrib",
				},
			},
			&preBuiltBuildConfig{
				targetOS:   "darwin",
				targetArch: darwinArchs,
				preBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib_{{ .Target }}/otelcol-contrib",
				},
			},
			&preBuiltBuildConfig{
				targetOS:   "windows",
				targetArch: winArchs,
				preBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib_{{ .Target }}/otelcol-contrib.exe",
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
			newContainerImageManifests(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
	}).WithPackagingDefaults().Build()

	// contrib build-only project
	contribBuildOnlyDist = newDistributionBuilder(ContribDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullDistBuildConfig{targetOS: "linux", targetArch: baseArchs, armVersion: []string{"7"}},
			&fullDistBuildConfig{targetOS: "darwin", targetArch: darwinArchs},
			&fullDistBuildConfig{targetOS: "windows", targetArch: winArchs},
		}
	}).Build()

	// k8s distro
	k8sArchs = []string{"amd64", "arm64", "ppc64le", "s390x"}
	k8sDist  = newDistributionBuilder(K8sDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullDistBuildConfig{targetOS: "linux", targetArch: k8sArchs},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", k8sArchs, containerImageOptions{}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)

		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", k8sArchs, containerImageOptions{}),
			newContainerImageManifests(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winContainerArchs, containerImageOptions{winVersion: "2022"}),
		)
	}).WithDefaultArchives().Build()
)

func BuildDist(dist string, onlyBuild bool) config.Project {
	switch dist {
	case CoreDistro:
		return otelColDist.BuildProject()
	case OTLPDistro:
		return otlpDist.BuildProject()
	case K8sDistro:
		return k8sDist.BuildProject()
	case ContribDistro:
		if onlyBuild {
			return contribBuildOnlyDist.BuildProject()
		}
		return contribDist.BuildProject()
	default:
		panic("Unknown distribution")
	}
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
		d.archives = newArchives(d.name, builds)
	})
	return b
}

func (b *distributionBuilder) WithDefaultNfpms() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.nfpms = newNfpms(d.name)
	})
	return b
}

func (b *distributionBuilder) WithDefaultMSIConfig() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.msiConfig = newMSIConfig(d.name)
	})
	return b
}

func (b *distributionBuilder) WithDefaultSigns() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.signs = Sign()
	})
	return b
}

func (b *distributionBuilder) WithDefaultDockerSigns() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.dockerSigns = DockerSigns()
	})
	return b
}

func (b *distributionBuilder) WithDefaultSBOMs() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.sboms = SBOM()
	})
	return b
}

func (b *distributionBuilder) WithPackagingDefaults() *distributionBuilder {
	return b.WithDefaultArchives().
		WithDefaultNfpms().
		WithDefaultMSIConfig().
		WithDefaultSigns().
		WithDefaultDockerSigns().
		WithDefaultSBOMs()
}

func (b *distributionBuilder) WithConfigFunc(configFunc func(*distribution)) *distributionBuilder {
	b.configFuncs = append(b.configFuncs, configFunc)
	return b
}

func (b *distributionBuilder) Build() *distribution {
	for _, configFunc := range b.configFuncs {
		configFunc(b.dist)
	}
	return b.dist
}

type buildConfig interface {
	Build(dist string) config.Build
	OS() string
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

func newNfpms(dist string) []config.NFPM {
	nfpmContents := config.NFPMContents{
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
	if _, ok := defaultConfigDists[dist]; ok {
		nfpmContents = append(nfpmContents, &config.NFPMContent{
			Source:      "config.yaml",
			Destination: path.Join("/etc", dist, "config.yaml"),
			Type:        "config|noreplace",
		})
	}
	return []config.NFPM{
		{
			ID:          dist,
			Builds:      []string{dist},
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

func newMSIConfig(dist string) []config.MSI {
	files := []string{"opentelemetry.ico"}
	if _, ok := defaultConfigDists[dist]; ok {
		files = append(files, "config.yaml")
	}
	return []config.MSI{
		{
			ID:    dist,
			Name:  fmt.Sprintf("%s_{{ .Version }}_{{ .Os }}_{{ .MsiArch }}", dist),
			WXS:   "windows-installer.wxs",
			Files: files,
		},
	}
}

func newArchives(dist string, builds []string) []config.Archive {
	return []config.Archive{
		{
			ID:           dist,
			NameTemplate: "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}{{ if .Mips }}_{{ .Mips }}{{ end }}",
			Builds:       builds,
		},
	}
}

func (d *distribution) BuildProject() config.Project {
	builds := make([]config.Build, 0, len(d.buildConfigs))
	for _, buildConfig := range d.buildConfigs {
		builds = append(builds, buildConfig.Build(d.name))
	}

	return config.Project{
		ProjectName: "opentelemetry-collector-releases",
		Checksum: config.Checksum{
			NameTemplate: fmt.Sprintf("{{ .ProjectName }}_%v_checksums.txt", d.name),
		},
		Env:             []string{"COSIGN_YES=true"},
		Builds:          builds,
		Archives:        d.archives,
		MSI:             d.msiConfig,
		NFPMs:           d.nfpms,
		Dockers:         d.containerImages,
		DockerManifests: d.containerImageManifests,
		Signs:           d.signs,
		DockerSigns:     d.dockerSigns,
		SBOMs:           d.sboms,
		Version:         2,
		Monorepo: config.Monorepo{
			TagPrefix: "v",
		},
		Partial: config.Partial{By: "target"},
	}
}

type fullDistBuildConfig struct {
	targetOS   string
	targetArch []string
	armVersion []string
}

func (c *fullDistBuildConfig) Build(dist string) config.Build {
	buildConfig := config.Build{
		ID:     dist + "-" + c.targetOS,
		Dir:    "_build",
		Binary: dist,
		BuildDetails: config.BuildDetails{
			Env:     []string{"CGO_ENABLED=0"},
			Flags:   []string{"-trimpath"},
			Ldflags: []string{"-s", "-w"},
		},
		Goos:   []string{c.targetOS},
		Goarch: c.targetArch,
		Goarm:  c.armVersion,
	}
	return buildConfig
}

func (c *fullDistBuildConfig) OS() string {
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
		Goarm:    ArmVersions(dist),
	}
}

func (c *preBuiltBuildConfig) OS() string {
	return c.targetOS
}

func dockerImageWithOS(dist, os, arch string, opts containerImageOptions) config.Docker {
	osArch := osArch{os: os, arch: arch, version: opts.version()}
	imageTemplates := make([]string, 0)
	for _, prefix := range imageRepos {
		dockerArchTag := strings.ReplaceAll(osArch.imageTag(), "/", "")
		imageTemplates = append(
			imageTemplates,
			fmt.Sprintf("%s/%s:{{ .Version }}-%s", prefix, imageName(dist), dockerArchTag),
			fmt.Sprintf("%s/%s:latest-%s", prefix, imageName(dist), dockerArchTag),
		)
	}

	label := func(name, template string) string {
		return fmt.Sprintf("--label=org.opencontainers.image.%s={{%s}}", name, template)
	}
	files := make([]string, 0)
	if _, ok := defaultConfigDists[dist]; ok {
		files = append(files, "config.yaml")
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
		Files:  files,
		Goos:   os,
		Goarch: arch,
	}
	if arch == ArmArch {
		imageConfig.Goarm = opts.armVersion
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
		case ArmArch:
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
		case ArmArch:
			return fmt.Sprintf("arm/v%s", o.version)
		}
	case "windows":
		return fmt.Sprintf("windows-%s-%s", o.version, o.arch)
	}
	return o.arch
}

func osDockerManifest(prefix, version, dist, os string, archs []string) config.DockerManifest {
	var imageTemplates []string
	for _, arch := range archs {
		switch arch {
		case ArmArch:
			for _, armVers := range ArmVersions(dist) {
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

	return config.DockerManifest{
		NameTemplate:   fmt.Sprintf("%s/%s:%s", prefix, imageName(dist), version),
		ImageTemplates: imageTemplates,
	}
}

func ArmVersions(dist string) []string {
	if dist == K8sDistro {
		return nil
	}
	return []string{"7"}
}

// imageName translates a distribution name to a container image name.
func imageName(dist string) string {
	return strings.Replace(dist, BinaryNamePrefix, ImageNamePrefix, 1)
}

// archName translates architecture to docker platform names.
func archName(arch, armVersion string) string {
	switch arch {
	case ArmArch:
		return fmt.Sprintf("%s/v%s", arch, armVersion)
	default:
		return arch
	}
}

func Sign() []config.Sign {
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

func DockerSigns() []config.Sign {
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

func SBOM() []config.SBOM {
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
