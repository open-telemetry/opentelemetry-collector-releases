package internal

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

var (
	baseArchs = []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}
	winArchs  = []string{"amd64", "arm64"}

	// otelcol (core) distro
	otelColDist = newDistributionBuilder(CoreDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullDistBuildConfig{
				targetOS:   []string{"darwin", "linux", "windows"},
				targetArch: baseArchs,
			},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", baseArchs, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)
	}).WithDefaultArchives().WithDefaultNfpms().WithDefaultMSIConfig().Build()

	// otlp distro
	otlpDist = newDistributionBuilder(OTLPDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullDistBuildConfig{
				targetOS:   []string{"darwin", "linux", "windows"},
				targetArch: []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"},
			},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)
	}).WithDefaultArchives().WithDefaultNfpms().WithDefaultMSIConfig().Build()

	// contrib distro
	contribDist = newDistributionBuilder(ContribDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&preBuiltBuildConfig{
				targetOS:   []string{"darwin", "linux", "windows"},
				targetArch: []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"},
				preBuilt: config.PreBuiltOptions{
					Path: "artifacts/otelcol-contrib_{{ .Target }}" +
						"/otelcol-contrib{{- if eq .Os \"windows\" }}.exe{{ end }}",
				},
			},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)
	}).WithDefaultArchives().WithDefaultNfpms().WithDefaultMSIConfig().Build()

	// contrib build-only project
	contribBuildOnlyDist = newDistributionBuilder(ContribDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullDistBuildConfig{
				targetOS:   []string{"darwin", "linux", "windows"},
				targetArch: []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"},
			},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}, containerImageOptions{armVersion: "7"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)
		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", baseArchs, containerImageOptions{}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)
	}).WithDefaultArchives().WithDefaultNfpms().WithDefaultMSIConfig().Build()

	// k8s distro
	k8sArchs = []string{"amd64", "arm64", "ppc64le", "s390x"}
	k8sDist  = newDistributionBuilder(K8sDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []buildConfig{
			&fullDistBuildConfig{
				targetOS:   []string{"linux"},
				targetArch: k8sArchs,
			},
		}
		d.containerImages = slices.Concat(
			newContainerImages(d.name, "linux", k8sArchs, containerImageOptions{}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImages(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
		)

		d.containerImageManifests = slices.Concat(
			newContainerImageManifests(d.name, "linux", k8sArchs, containerImageOptions{}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2019"}),
			newContainerImageManifests(d.name, "windows", winArchs, containerImageOptions{winVersion: "2022"}),
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
	b.dist.archives = newArchives(b.dist.name)
	return b
}

func (b *distributionBuilder) WithDefaultNfpms() *distributionBuilder {
	b.dist.nfpms = newNfpms(b.dist.name)
	return b
}

func (b *distributionBuilder) WithDefaultMSIConfig() *distributionBuilder {
	b.dist.msiConfig = newMSIConfig(b.dist.name)
	return b
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
}

func newContainerImageManifests(dist, os string, archs []string, opts containerImageOptions) []config.DockerManifest {
	imageNames := []string{DockerHub, GHCR}
	tags := []string{`{{ .Version }}`, "latest"}

	if os == "windows" {
		for i, tag := range tags {
			tags[i] = fmt.Sprintf("%s-%s-%s", tag, os, opts.winVersion)
		}
	}
	var r []config.DockerManifest
	for _, imageName := range imageNames {
		for _, tag := range tags {
			r = append(r, osDockerManifest(imageName, tag, dist, os, archs))
		}
	}
	return r
}

type containerImageOptions struct {
	armVersion string
	winVersion string
}

// There are lots of complications around this function.
// Should receive target OS and target arch. CGO is disabled so can cross compile.
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
	if _, ok := DefaultConfigDists[dist]; ok {
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
	if _, ok := DefaultConfigDists[dist]; ok {
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

func newArchives(dist string) []config.Archive {
	return []config.Archive{
		{
			ID:           dist,
			NameTemplate: "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}{{ if .Mips }}_{{ .Mips }}{{ end }}",
			Builds:       []string{dist},
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
		Signs:           Sign(),
		DockerSigns:     DockerSigns(),
		SBOMs:           SBOM(),
		Version:         2,
		Monorepo: config.Monorepo{
			TagPrefix: "v",
		},
		Partial: config.Partial{By: "target"},
	}
}

type fullDistBuildConfig struct {
	targetOS   []string
	targetArch []string
}

func (c *fullDistBuildConfig) Build(dist string) config.Build {
	return config.Build{
		// ID:     dist + "-" + c.targetOS,
		ID:     dist,
		Dir:    "_build",
		Binary: dist,
		BuildDetails: config.BuildDetails{
			Env:     []string{"CGO_ENABLED=0"},
			Flags:   []string{"-trimpath"},
			Ldflags: []string{"-s", "-w"},
		},
		// Goos:   []string{c.targetOS},
		Goos:   c.targetOS,
		Goarch: c.targetArch,
		Goarm:  ArmVersions(dist),
		Ignore: IgnoreBuildCombinations(dist),
	}
}

type preBuiltBuildConfig struct {
	targetOS   []string
	targetArch []string
	preBuilt   config.PreBuiltOptions
}

func (c *preBuiltBuildConfig) Build(dist string) config.Build {
	return config.Build{
		// ID:     dist + "-" + c.targetOS,
		ID:       dist,
		Builder:  "prebuilt",
		PreBuilt: c.preBuilt,
		Dir:      "_build",
		Binary:   dist,
		// Goos:   []string{c.targetOS},
		Goos:   c.targetOS,
		Goarch: c.targetArch,
		Goarm:  ArmVersions(dist),
		Ignore: IgnoreBuildCombinations(dist),
	}
}

func dockerImageWithOS(dist, os, arch string, opts containerImageOptions) config.Docker {
	dockerArchName := osArchName(os, arch, opts)
	imageTemplates := make([]string, 0)
	for _, prefix := range ImagePrefixes {
		dockerArchTag := strings.ReplaceAll(dockerArchName, "/", "")
		// if os == "windows" {
		// 	dockerArchTag = fmt.Sprintf("%s-%s-%s", os, opts.winVersion, dockerArchTag)
		// }
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
	if _, ok := DefaultConfigDists[dist]; ok {
		files = append(files, "config.yaml")
	}
	imageConfig := config.Docker{
		ImageTemplates: imageTemplates,
		Dockerfile:     "Dockerfile",
		Use:            "buildx",
		BuildFlagTemplates: []string{
			"--pull",
			fmt.Sprintf("--platform=linux/%s", dockerArchName),
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
	if arch == "arm" {
		imageConfig.Goarm = opts.armVersion
	}
	return imageConfig
}

func osArchName(os, arch string, opts containerImageOptions) string {
	armVersion := opts.armVersion
	winVersion := opts.winVersion

	switch os {
	case "linux":
		switch arch {
		case ArmArch:
			return fmt.Sprintf("%s/v%s", arch, armVersion)
		}
	case "windows":
		return fmt.Sprintf("%s-%s-%s", os, winVersion, arch)
	}
	return arch
}

func osDockerManifest(prefix, version, dist, os string, archs []string) config.DockerManifest {
	var imageTemplates []string
	for _, arch := range archs {
		switch arch {
		case ArmArch:
			for _, armVers := range ArmVersions(dist) {
				dockerArchTag := strings.ReplaceAll(archName(arch, armVers), "/", "")
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
