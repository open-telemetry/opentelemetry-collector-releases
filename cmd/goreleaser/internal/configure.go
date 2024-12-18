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
	ImagePrefixes      = []string{DockerHub, GHCR}
	Architectures      = []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}
	DefaultConfigDists = map[string]bool{CoreDistro: true, ContribDistro: true}
	MSIWindowsDists    = map[string]bool{CoreDistro: true, ContribDistro: true, OTLPDistro: true}
	K8sDockerSkipArchs = map[string]bool{"arm": true, "386": true}
	K8sGoos            = []string{"linux"}
	K8sArchs           = []string{"amd64", "arm64", "ppc64le", "s390x"}
)

func GenerateContribBuildOnly(dist string, buildOrRest bool) config.Project {
	return config.Project{
		ProjectName: "opentelemetry-collector-releases",
		Builds:      Builds(dist, buildOrRest),
		Version:     2,
		Monorepo: config.Monorepo{
			TagPrefix: "v",
		},
	}
}

func Generate(dist string, buildOrRest bool) config.Project {
	return config.Project{
		ProjectName: "opentelemetry-collector-releases",
		Checksum: config.Checksum{
			NameTemplate: fmt.Sprintf("{{ .ProjectName }}_%v_checksums.txt", dist),
		},
		Env:             []string{"COSIGN_YES=true"},
		Builds:          Builds(dist, buildOrRest),
		Archives:        Archives(dist),
		MSI:             WinPackages(dist),
		NFPMs:           Packages(dist),
		Dockers:         DockerImages(dist),
		DockerManifests: DockerManifests(dist),
		Signs:           Sign(),
		DockerSigns:     DockerSigns(),
		SBOMs:           SBOM(),
		Version:         2,
		Monorepo: config.Monorepo{
			TagPrefix: "v",
		},
	}
}

func Builds(dist string, buildOrRest bool) []config.Build {
	return []config.Build{
		Build(dist, buildOrRest),
	}
}

// Build configures a goreleaser build.
// https://goreleaser.com/customization/build/
func Build(dist string, buildOrRest bool) config.Build {
	goos := []string{"darwin", "linux", "windows"}
	archs := Architectures

	if dist == ContribDistro && !buildOrRest {
		// only return build config for contrib build file
		return config.Build{
			ID:      dist,
			Builder: "prebuilt",
			PreBuilt: config.PreBuiltOptions{
				Path: "artifacts/otelcol-contrib_{{ .Target }}" +
					"/otelcol-contrib{{- if eq .Os \"windows\" }}.exe{{ end }}",
			},
			Goos:   goos,
			Goarch: archs,
			Goarm:  ArmVersions(dist),
			Dir:    "_build",
			Binary: dist,
			Ignore: IgnoreBuildCombinations(dist),
		}
	}

	if dist == K8sDistro {
		goos = K8sGoos
		archs = K8sArchs
	}

	return config.Build{
		ID:     dist,
		Dir:    "_build",
		Binary: dist,
		BuildDetails: config.BuildDetails{
			Env:     []string{"CGO_ENABLED=0"},
			Flags:   []string{"-trimpath"},
			Ldflags: []string{"-s", "-w"},
		},
		Goos:   goos,
		Goarch: archs,
		Goarm:  ArmVersions(dist),
		Ignore: IgnoreBuildCombinations(dist),
	}
}

func IgnoreBuildCombinations(dist string) []config.IgnoredBuild {
	if dist == K8sDistro {
		return nil
	}
	return []config.IgnoredBuild{
		{Goos: "darwin", Goarch: "386"},
		{Goos: "darwin", Goarch: "arm"},
		{Goos: "darwin", Goarch: "s390x"},
		{Goos: "windows", Goarch: "arm"},
		{Goos: "windows", Goarch: "arm64"},
		{Goos: "windows", Goarch: "s390x"},
	}
}

func ArmVersions(dist string) []string {
	if dist == K8sDistro {
		return nil
	}
	return []string{"7"}
}

func Archives(dist string) []config.Archive {
	return []config.Archive{
		Archive(dist),
	}
}

// Archive configures a goreleaser archive (tarball).
// https://goreleaser.com/customization/archive/
func Archive(dist string) config.Archive {
	return config.Archive{
		ID:           dist,
		NameTemplate: "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}{{ if .Mips }}_{{ .Mips }}{{ end }}",
		Builds:       []string{dist},
	}
}

func WinPackages(dist string) []config.MSI {
	if _, ok := MSIWindowsDists[dist]; !ok {
		return nil
	}
	return []config.MSI{
		WinPackage(dist),
	}
}

// Package configures goreleaser to build a Windows MSI package.
// https://goreleaser.com/customization/msi/
func WinPackage(dist string) config.MSI {
	files := []string{"opentelemetry.ico"}
	if _, ok := DefaultConfigDists[dist]; ok {
		files = append(files, "config.yaml")
	}
	return config.MSI{
		ID:    dist,
		Name:  fmt.Sprintf("%s_{{ .Version }}_{{ .Os }}_{{ .MsiArch }}", dist),
		WXS:   "windows-installer.wxs",
		Files: files,
	}
}

func Packages(dist string) []config.NFPM {
	if dist == K8sDistro {
		return nil
	}
	return []config.NFPM{
		Package(dist),
	}
}

// Package configures goreleaser to build a system package.
// https://goreleaser.com/customization/nfpm/
func Package(dist string) config.NFPM {
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
	return config.NFPM{
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
	}
}

func DockerImages(dist string) []config.Docker {
	var r []config.Docker
	for _, arch := range Architectures {
		if dist == K8sDistro && K8sDockerSkipArchs[arch] {
			continue
		}
		switch arch {
		case ArmArch:
			for _, vers := range ArmVersions(dist) {
				r = append(r, DockerImage(dist, arch, vers))
			}
		default:
			r = append(r, DockerImage(dist, arch, ""))
		}
	}
	return r
}

// DockerImage configures goreleaser to build a container image.
// https://goreleaser.com/customization/docker/
func DockerImage(dist, arch, armVersion string) config.Docker {
	dockerArchName := archName(arch, armVersion)
	imageTemplates := make([]string, 0)
	for _, prefix := range ImagePrefixes {
		dockerArchTag := strings.ReplaceAll(dockerArchName, "/", "")
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
	return config.Docker{
		ImageTemplates: imageTemplates,
		Dockerfile:     "Dockerfile",

		Use: "buildx",
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
		Goos:   "linux",
		Goarch: arch,
		Goarm:  armVersion,
	}
}

func DockerManifests(dist string) []config.DockerManifest {
	r := make([]config.DockerManifest, 0)
	for _, prefix := range ImagePrefixes {
		r = append(r, DockerManifest(prefix, `{{ .Version }}`, dist))
		r = append(r, DockerManifest(prefix, "latest", dist))
	}
	return r
}

// DockerManifest configures goreleaser to build a multi-arch container image manifest.
// https://goreleaser.com/customization/docker_manifest/
func DockerManifest(prefix, version, dist string) config.DockerManifest {
	var imageTemplates []string
	for _, arch := range Architectures {
		if dist == K8sDistro {
			if _, ok := K8sDockerSkipArchs[arch]; ok {
				continue
			}
		}
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
