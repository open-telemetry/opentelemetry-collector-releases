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
	DockerHub        = "johnpeterson785"
	GHCR             = "ghcr.io/jackgopack4/opentelemetry-collector-releases"
	BinaryNamePrefix = "otelcol"
	ImageNamePrefix  = "opentelemetry-collector"
)

var (
	ImagePrefixes      = []string{DockerHub, GHCR}
	Architectures      = []string{"386", "amd64", "arm", "arm64", "ppc64le", "s390x"}
	ArmVersions        = []string{"7"}
	DefaultConfigDists = map[string]bool{CoreDistro: true, ContribDistro: true}
	MSIWindowsDists    = map[string]bool{CoreDistro: true, ContribDistro: true, OTLPDistro: true}
	K8sDockerSkipArchs = map[string]bool{"arm": true, "386": true}
	K8sGoos            = []string{"linux"}
	K8sArchs           = []string{"amd64", "arm64", "ppc64le", "s390x"}
	AlwaysIgnored      = map[config.IgnoredBuild]bool{
		{Goos: "darwin", Goarch: "386"}:    true,
		{Goos: "darwin", Goarch: "arm"}:    true,
		{Goos: "darwin", Goarch: "s390x"}:  true,
		{Goos: "windows", Goarch: "arm"}:   true,
		{Goos: "windows", Goarch: "arm64"}: true,
		{Goos: "windows", Goarch: "s390x"}: true,
	}
)

// Copied from go/src/internal/platform/supported.go, see:
// https://cs.opensource.google/go/go/+/d7fcb5cf80953f1d63246f1ae9defa60c5ce2d76:src/internal/platform/supported.go;l=222
func InternalLinkPIESupported(goos, goarch string) bool {
	switch goos + "/" + goarch {
	case "android/arm64",
		"darwin/amd64", "darwin/arm64",
		"linux/amd64", "linux/arm64", "linux/ppc64le",
		"windows/386", "windows/amd64", "windows/arm", "windows/arm64":
		return true
	}
	return false
}

func Generate(dist string) config.Project {
	return config.Project{
		ProjectName: "opentelemetry-collector-releases",
		Checksum: config.Checksum{
			NameTemplate: fmt.Sprintf("{{ .ProjectName }}_%v_checksums.txt", dist),
		},
		Env:             []string{"COSIGN_YES=true"},
		Builds:          Builds(dist),
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

func Builds(dist string) []config.Build {
	return []config.Build{
		Build(dist, true),
		Build(dist, false),
	}
}

func generateIgnored(goos, archs []string, pie bool) []config.IgnoredBuild {
	ignored := make([]config.IgnoredBuild, 0)
	var build config.IgnoredBuild
	for _, goos := range goos {
		for _, arch := range archs {
			build = config.IgnoredBuild{
				Goos:   goos,
				Goarch: arch,
			}
			if _, ok := AlwaysIgnored[build]; ok || !pie && InternalLinkPIESupported(goos, arch) || pie && !InternalLinkPIESupported(goos, arch) {
				ignored = append(ignored, build)
			}
		}
	}
	return ignored
}

// Build configures a goreleaser build.
// https://goreleaser.com/customization/build/
func Build(dist string, pie bool) config.Build {
	var goos []string
	var archs []string
	var ignore []config.IgnoredBuild
	var armVersions []string
	id := dist
	ldflags := []string{"-s", "-w"}
	if pie {
		ldflags = append(ldflags, "-buildmode=pie")
		id = id + "-pie"
	}
	if dist == K8sDistro {
		goos = K8sGoos
		archs = K8sArchs
		armVersions = make([]string, 0)
	} else {
		goos = []string{"darwin", "linux", "windows"}
		archs = Architectures
		armVersions = ArmVersions
	}
	ignore = generateIgnored(goos, archs, pie)
	return config.Build{
		ID:     id,
		Dir:    "_build",
		Binary: dist,
		BuildDetails: config.BuildDetails{
			Env:     []string{"CGO_ENABLED=0"},
			Flags:   []string{"-trimpath"},
			Ldflags: ldflags,
		},
		Goos:   goos,
		Goarch: archs,
		Goarm:  armVersions,
		Ignore: ignore,
	}
}

func Archives(dist string) (r []config.Archive) {
	return []config.Archive{
		Archive(dist, true),
		Archive(dist, false),
	}
}

// Archive configures a goreleaser archive (tarball).
// https://goreleaser.com/customization/archive/
func Archive(dist string, pie bool) config.Archive {
	id := dist
	build := dist
	if pie {
		id = id + "-pie"
		build = build + "-pie"
	}
	return config.Archive{
		ID:           id,
		NameTemplate: "{{ .Binary }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}{{ if .Arm }}v{{ .Arm }}{{ end }}{{ if .Mips }}_{{ .Mips }}{{ end }}",
		Builds:       []string{build},
	}
}

func WinPackages(dist string) []config.MSI {
	if _, ok := MSIWindowsDists[dist]; !ok {
		return []config.MSI{}
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

func Packages(dist string) (r []config.NFPM) {
	if dist == K8sDistro {
		return []config.NFPM{}
	}
	return []config.NFPM{
		Package(dist, true),
		Package(dist, false),
	}
}

// Package configures goreleaser to build a system package.
// https://goreleaser.com/customization/nfpm/
func Package(dist string, pie bool) config.NFPM {
	id := dist
	build := dist
	if pie {
		id = id + "-pie"
		build = build + "-pie"
	}
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
		ID:      id,
		Builds:  []string{build},
		Formats: []string{"deb", "rpm"},

		License:     "Apache 2.0",
		Description: fmt.Sprintf("OpenTelemetry Collector - %s", dist),
		Maintainer:  "The OpenTelemetry Collector maintainers <cncf-opentelemetry-maintainers@lists.cncf.io>",
		Overrides: map[string]config.NFPMOverridables{
			"rpm": {
				Dependencies: []string{
					"/bin/sh",
				},
			},
		},

		NFPMOverridables: config.NFPMOverridables{
			PackageName: id,
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
	r := make([]config.Docker, 0)
	for _, arch := range Architectures {
		if dist == K8sDistro {
			if _, ok := K8sDockerSkipArchs[arch]; ok {
				continue
			}
		}
		switch arch {
		case ArmArch:
			for _, vers := range ArmVersions {
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
			for _, armVers := range ArmVersions {
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
