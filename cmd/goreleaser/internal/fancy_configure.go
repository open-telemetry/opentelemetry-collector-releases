package internal

import (
	"fmt"
	"path"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

var (
	// otelcol (core) distro
	otelColBuildProj = newDistributionBuilder(CoreDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []fullDistBuildConfig{
			{
				targetOS:   []string{"darwin", "linux", "windows"},
				targetArch: []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"},
			},
			// {
			// 	targetOS:   []string{"linux"},
			// 	targetArch: []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"},
			// },
			// {
			// 	targetOS:   []string{"darwin"},
			// 	targetArch: []string{"amd64", "arm64"},
			// },
		}
		d.containerImages = newContainerImages(d.name, "linux", []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}, "7")
		d.containerImageManifests = newContainerImageManifests(d.name, ImagePrefixes, []string{`{{ .Version }}`, "latest"})
	}).WithDefaultArchives().WithDefaultNfpms().WithDefaultMSIConfig().Build()

	// otlp distro
	otelColOTLPBuildProj = newDistributionBuilder(OTLPDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []fullDistBuildConfig{
			{
				targetOS:   []string{"darwin", "linux", "windows"},
				targetArch: []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"},
			},
		}
		d.containerImages = newContainerImages(d.name, "linux", []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}, "7")
		d.containerImageManifests = newContainerImageManifests(d.name, ImagePrefixes, []string{`{{ .Version }}`, "latest"})
	}).WithDefaultArchives().WithDefaultNfpms().WithDefaultMSIConfig().Build()

	// k8s distro
	otelK8sBuildProj = newDistributionBuilder(K8sDistro).WithConfigFunc(func(d *distribution) {
		d.buildConfigs = []fullDistBuildConfig{
			{
				targetOS:   []string{"linux"},
				targetArch: []string{"amd64", "arm64", "ppc64le", "s390x"},
			},
		}
		d.containerImages = newContainerImages(d.name, "linux", []string{"amd64", "arm64", "ppc64le", "s390x"}, "")
		d.containerImageManifests = newContainerImageManifests(d.name, ImagePrefixes, []string{`{{ .Version }}`, "latest"})
	}).WithDefaultArchives().Build()
)

func BuildDist(dist string, buildOrRest bool) config.Project {
	switch dist {
	case CoreDistro:
		return otelColBuildProj.BuildProject(buildOrRest)
	case OTLPDistro:
		return otelColOTLPBuildProj.BuildProject(buildOrRest)
	case K8sDistro:
		return otelK8sBuildProj.BuildProject(buildOrRest)
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

type distribution struct {
	// Name of the distribution (i.e. otelcol, otelcol-contrib, k8s)
	name string

	buildConfigs            []fullDistBuildConfig
	archives                []config.Archive
	msiConfig               []config.MSI
	nfpms                   []config.NFPM
	containerImages         []config.Docker
	containerImageManifests []config.DockerManifest
}

func newContainerImageManifests(dist string, imageNames, tags []string) []config.DockerManifest {
	var r []config.DockerManifest
	for _, imageName := range imageNames {
		for _, tag := range tags {
			r = append(r, DockerManifest(imageName, tag, dist))
		}
	}
	return r
}

// There are lots of complications around this function.
// Should receive target OS and target arch. CGO is disabled so can cross compile.
func newContainerImages(dist string, targetOS string, targetArchs []string, armVersion string) []config.Docker {
	images := []config.Docker{}
	for _, targetArch := range targetArchs {
		if armVersion != "" && targetArch == "arm" {
			images = append(images, DockerImage(dist, targetArch, armVersion))
			continue
		}
		images = append(images, DockerImage(dist, targetArch, ""))
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

func (d *distribution) BuildProject(buildOrRest bool) config.Project {
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
	}
}

type fullDistBuildConfig struct {
	// Target OS (i.e. linux, darwin, windows)
	// targetOS string
	targetOS []string
	// Target architecture (i.e. amd64, arm64)
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
