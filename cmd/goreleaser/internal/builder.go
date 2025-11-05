// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"path"

	"github.com/goreleaser/goreleaser-pro/v2/pkg/config"
)

const containerEphemeralTag = "CONTAINER_IMAGE_EPHEMERAL_TAG={{ if .IsNightly }}nightly{{ else }}latest{{ end }}"

// distributionBuilder is used to build distribution configurations.
type distributionBuilder struct {
	dist        *distribution
	configFuncs []func(*distribution)
}

// buildConfig is the interface for build configurations.
type buildConfig interface {
	Build(dist string) config.Build
	OS() string
}

// fullBuildConfig represents a full build configuration.
type fullBuildConfig struct {
	TargetOS     string
	TargetArch   []string
	ArmVersion   []string
	Ppc64Version []string
	BinaryName   string
	BuildDir     string
}

func (c *fullBuildConfig) Build(dist string) config.Build {
	buildConfiguration := config.Build{
		ID:     dist + "-" + c.TargetOS,
		Binary: dist,
		BuildDetails: config.BuildDetails{
			Flags:   []string{"{{ .Env.BUILD_FLAGS }}"},
			Ldflags: []string{"{{ .Env.LD_FLAGS }}"},
		},
		Goos:    []string{c.TargetOS},
		Goarch:  c.TargetArch,
		Goarm:   c.ArmVersion,
		Goppc64: c.Ppc64Version,
	}

	if c.BinaryName != "" {
		buildConfiguration.Binary = c.BinaryName
	}

	if c.BuildDir != "" {
		buildConfiguration.Dir = c.BuildDir
	}

	return buildConfiguration
}

func (c *fullBuildConfig) OS() string {
	return c.TargetOS
}

// preBuiltBuildConfig represents a pre-built build configuration.
type preBuiltBuildConfig struct {
	TargetOS   string
	TargetArch []string
	PreBuilt   config.PreBuiltOptions
}

func (c *preBuiltBuildConfig) Build(dist string) config.Build {
	return config.Build{
		ID:       dist + "-" + c.TargetOS,
		Builder:  "prebuilt",
		PreBuilt: c.PreBuilt,
		Dir:      "_build",
		Binary:   dist,
		Goos:     []string{c.TargetOS},
		Goarch:   c.TargetArch,
		Goarm:    armVersions(dist),
		Goppc64:  []string{"power8"},
	}
}

func (c *preBuiltBuildConfig) OS() string {
	return c.TargetOS
}

// distribution represents a collector distribution configuration.
type distribution struct {
	// Name of the distribution (i.e. otelcol, otelcol-contrib, k8s)
	Name string

	BuildConfigs            []buildConfig
	Archives                []config.Archive
	MsiConfig               []config.MSI
	Nfpms                   []config.NFPM
	ContainerImages         []config.Docker
	ContainerImageManifests []config.DockerManifest
	Signs                   []config.Sign
	DockerSigns             []config.Sign
	Sboms                   []config.SBOM
	Nightly                 config.Nightly
	Checksum                config.Checksum
	Partial                 config.Partial
	Monorepo                config.Monorepo
	Release                 config.Release
	Snapshot                config.Snapshot
	Changelog               config.Changelog
	Env                     []string
	EnableCgo               bool
	LdFlags                 string
	GoTags                  string
}

// buildProject builds the goreleaser project configuration from the distribution.
func (d *distribution) buildProject() config.Project {
	builds := make([]config.Build, 0, len(d.BuildConfigs))
	for _, buildConfig := range d.BuildConfigs {
		builds = append(builds, buildConfig.Build(d.Name))
	}

	return config.Project{
		ProjectName:     projectName,
		Release:         d.Release,
		Checksum:        d.Checksum,
		Env:             d.Env,
		Builds:          builds,
		Archives:        d.Archives,
		MSI:             d.MsiConfig,
		NFPMs:           d.Nfpms,
		Dockers:         d.ContainerImages,
		DockerManifests: d.ContainerImageManifests,
		Signs:           d.Signs,
		DockerSigns:     d.DockerSigns,
		SBOMs:           d.Sboms,
		Nightly:         d.Nightly,
		Version:         2,
		Monorepo:        d.Monorepo,
		Partial:         d.Partial,
		Snapshot:        d.Snapshot,
		Changelog:       d.Changelog,
	}
}

// newDistributionBuilder creates a new distribution builder.
func newDistributionBuilder(name string) *distributionBuilder {
	return &distributionBuilder{dist: &distribution{Name: name}}
}

func (b *distributionBuilder) withDefaultArchives() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		builds := make([]string, 0, len(d.BuildConfigs))
		for _, build := range d.BuildConfigs {
			builds = append(builds, fmt.Sprintf("%s-%s", d.Name, build.OS()))
		}
		d.Archives = b.newArchives(d.Name, builds)
	})
	return b
}

func (b *distributionBuilder) withBinArchive() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.Archives = append(d.Archives, config.Archive{
			Formats: []string{"binary"},
		})
	})
	return b
}

func (b *distributionBuilder) newArchives(dist string, ids []string) []config.Archive {
	return []config.Archive{
		{
			ID:           dist,
			NameTemplate: "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}{{ if .Mips }}_{{ .Mips }}{{ end }}",
			IDs:          ids,
		},
	}
}

func (b *distributionBuilder) withDefaultNfpms() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.Nfpms = b.newNfpms(d.Name)
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
			IDs:         []string{dist + "-linux"},
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

func (b *distributionBuilder) withDefaultMSIConfig() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.MsiConfig = b.newMSIConfig(d.Name)
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

func (b *distributionBuilder) withDefaultSigns() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.Signs = b.signs()
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

func (b *distributionBuilder) withDefaultDockerSigns() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.DockerSigns = b.newDockerSigns()
	})
	return b
}

func (b *distributionBuilder) newDockerSigns() []config.Sign {
	condition := ""
	switch b.dist.Name {
	case ocbBinary, opampBinary:
		condition = "$SKIP_SIGNS != 'true'"
	}
	return []config.Sign{
		{
			If:        condition,
			Artifacts: "all",
			Args: []string{
				"sign",
				"${artifact}",
			},
		},
	}
}

func (b *distributionBuilder) withNightlyConfig() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.Nightly = b.newNightly()
	})
	return b
}

func (b *distributionBuilder) newNightly() config.Nightly {
	return config.Nightly{
		VersionTemplate:   "{{ incpatch .Version}}-nightly.{{ .ShortCommit }}",
		TagName:           fmt.Sprintf("nightly-%s", b.dist.Name),
		PublishRelease:    false,
		KeepSingleRelease: true,
	}
}

func (b *distributionBuilder) withDefaultSBOMs() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.Sboms = b.newSboms()
	})
	return b
}

func (b *distributionBuilder) newSboms() []config.SBOM {
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

func (b *distributionBuilder) withDefaultChecksum() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.Checksum = config.Checksum{
			NameTemplate: fmt.Sprintf("{{ .ProjectName }}_%v{{ if eq .Runtime.Goos \"windows\" }}_{{ .Runtime.Goos }}{{ end }}_checksums.txt", d.Name),
		}
	})
	return b
}

func (b *distributionBuilder) withDefaultBinaryChecksum() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.Checksum = config.Checksum{
			NameTemplate: "checksums.txt",
		}
	})
	return b
}

func (b *distributionBuilder) withDefaultSnapshot() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		d.Snapshot = config.Snapshot{
			VersionTemplate: "{{ incpatch .Version }}-next",
		}
	})
	return b
}

func (b *distributionBuilder) withDefaultMonorepo() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.Monorepo = config.Monorepo{
			TagPrefix: "v",
		}
	})
	return b
}

func (b *distributionBuilder) withBinaryMonorepo(dir string) *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.Monorepo = b.newBinaryMonorepo(dir)
	})
	return b
}

func (b *distributionBuilder) newBinaryMonorepo(dir string) config.Monorepo {
	return config.Monorepo{
		TagPrefix: fmt.Sprintf("cmd/%s/", b.dist.Name),
		Dir:       dir,
	}
}

func (b *distributionBuilder) withDefaultEnv() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		ldFlags := "-s -w"
		if b.dist.LdFlags != "" {
			ldFlags = b.dist.LdFlags
		}

		env := []string{
			"COSIGN_YES=true",
			"LD_FLAGS=" + ldFlags,
			"BUILD_FLAGS=-trimpath",
			containerEphemeralTag,
			"GOPROXY=https://proxy.golang.org,direct",
		}
		if b.dist.GoTags != "" {
			env = append(env, "GO_TAGS="+b.dist.GoTags)
		}
		if !b.dist.EnableCgo {
			env = append(env, "CGO_ENABLED=0")
		}
		b.dist.Env = append(b.dist.Env, env...)
	})
	return b
}

func (b *distributionBuilder) withDefaultPartial() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.Partial = config.Partial{
			By: "target",
		}
	})
	return b
}

func (b *distributionBuilder) withDefaultRelease() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.Release = config.Release{
			ReplaceExistingArtifacts: true,
		}
	})
	return b
}

func (b *distributionBuilder) withDefaultBinaryRelease(header string) *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		b.dist.Release = b.newBinaryRelease(header)
	})
	return b
}

func (b *distributionBuilder) newBinaryRelease(header string) config.Release {
	return config.Release{
		MakeLatest: "false",
		Header: config.IncludedMarkdown{
			Content: header,
		},
		GitHub: config.Repo{
			Owner: "open-telemetry",
			Name:  "opentelemetry-collector-releases",
		},
	}
}

func (b *distributionBuilder) withPackagingDefaults() *distributionBuilder {
	return b.withDefaultArchives().
		withDefaultSnapshot().
		withDefaultChecksum().
		withDefaultMonorepo().
		withDefaultEnv().
		withDefaultNfpms().
		withDefaultMSIConfig().
		withDefaultSigns().
		withDefaultDockerSigns().
		withDefaultSBOMs().
		withDefaultPartial().
		withDefaultRelease().
		withNightlyConfig()
}

func (b *distributionBuilder) withBinaryPackagingDefaults() *distributionBuilder {
	b.dist.Changelog = config.Changelog{
		Disable: "true",
	}
	return b.withBinArchive().
		withDefaultSnapshot().
		withDefaultChecksum().
		withDefaultEnv().
		withDefaultSigns().
		withDefaultDockerSigns().
		withDefaultSBOMs().
		withDefaultBinaryChecksum()
}

// withConfigFunc adds a configuration function to the builder.
func (b *distributionBuilder) withConfigFunc(configFunc func(*distribution)) *distributionBuilder {
	b.configFuncs = append(b.configFuncs, configFunc)
	return b
}

func (b *distributionBuilder) withDefaultConfigIncluded() *distributionBuilder {
	b.configFuncs = append(b.configFuncs, func(d *distribution) {
		for i, container := range d.ContainerImages {
			container.Files = append(container.Files, "config.yaml")
			d.ContainerImages[i] = container
		}

		for i, nfpm := range d.Nfpms {
			nfpm.Contents = append(nfpm.Contents, config.NFPMContent{
				Source:      "config.yaml",
				Destination: path.Join("/etc", d.Name, "config.yaml"),
				Type:        "config|noreplace",
			})
			d.Nfpms[i] = nfpm
		}

		for i := range d.MsiConfig {
			d.MsiConfig[i].Files = append(d.MsiConfig[i].Files, "config.yaml")
		}
	})
	return b
}

// Build constructs the final distribution.
func (b *distributionBuilder) build() *distribution {
	for _, configFunc := range b.configFuncs {
		configFunc(b.dist)
	}
	return b.dist
}
