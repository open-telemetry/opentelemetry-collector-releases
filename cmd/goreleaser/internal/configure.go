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

	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/nfpm/v2/files"
)

const ArmArch = "arm"

var (
	ImagePrefixes = []string{"ghcr.io/axoflow/axoflow-otel-collector"}
	Architectures = []string{"amd64", "arm64"}
	ArmVersions   = []string{}
)

func Generate(imagePrefixes []string, dists []string) config.Project {
	return config.Project{
		ProjectName: "axoflow-otel-collector-releases",
		Checksum: config.Checksum{
			NameTemplate: "{{ .ProjectName }}_checksums.txt",
		},

		Builds:          Builds(dists),
		Archives:        Archives(dists),
		Dockers:         DockerImages(imagePrefixes, dists),
		DockerManifests: DockerManifests(imagePrefixes, dists),
	}
}

func Builds(dists []string) (r []config.Build) {
	for _, dist := range dists {
		r = append(r, Build(dist))
	}
	return
}

// Build configures a goreleaser build.
// https://goreleaser.com/customization/build/
func Build(dist string) config.Build {
	return config.Build{
		ID:     dist,
		Dir:    path.Join("distributions", dist, "_build"),
		Binary: dist,
		BuildDetails: config.BuildDetails{
			Env:     []string{"CGO_ENABLED=0"},
			Flags:   []string{"-trimpath"},
			Ldflags: []string{"-s", "-w"},
		},
		Goos:   []string{"linux"},
		Goarch: Architectures,
		Goarm:  ArmVersions,
	}
}

func Archives(dists []string) (r []config.Archive) {
	for _, dist := range dists {
		r = append(r, Archive(dist))
	}
	return
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

func Packages(dists []string) (r []config.NFPM) {
	for _, dist := range dists {
		r = append(r, Package(dist))
	}
	return
}

// Package configures goreleaser to build a system package.
// https://goreleaser.com/customization/nfpm/
func Package(dist string) config.NFPM {
	return config.NFPM{
		ID:      dist,
		Builds:  []string{dist},
		Formats: []string{"apk", "deb", "rpm"},

		License:     "Apache 2.0",
		Description: fmt.Sprintf("OpenTelemetry Collector - %s", dist),
		Maintainer:  "The OpenTelemetry Collector maintainers <cncf-opentelemetry-maintainers@lists.cncf.io>",

		NFPMOverridables: config.NFPMOverridables{
			PackageName: dist,
			Scripts: config.NFPMScripts{
				PreInstall:  path.Join("distributions", dist, "preinstall.sh"),
				PostInstall: path.Join("distributions", dist, "postinstall.sh"),
				PreRemove:   path.Join("distributions", dist, "preremove.sh"),
			},
			Contents: files.Contents{
				{
					Source:      path.Join("distributions", dist, fmt.Sprintf("%s.service", dist)),
					Destination: path.Join("/lib", "systemd", "system", fmt.Sprintf("%s.service", dist)),
				},
				{
					Source:      path.Join("distributions", dist, fmt.Sprintf("%s.conf", dist)),
					Destination: path.Join("/etc", dist, fmt.Sprintf("%s.conf", dist)),
					Type:        "config|noreplace",
				},
				{
					Source:      path.Join("configs", fmt.Sprintf("%s.yaml", dist)),
					Destination: path.Join("/etc", dist, "config.yaml"),
					Type:        "config",
				},
			},
		},
	}
}

func DockerImages(imagePrefixes, dists []string) (r []config.Docker) {
	for _, dist := range dists {
		for _, arch := range Architectures {
			switch arch {
			case ArmArch:
				for _, vers := range ArmVersions {
					r = append(r, DockerImage(imagePrefixes, dist, arch, vers))
				}
			default:
				r = append(r, DockerImage(imagePrefixes, dist, arch, ""))
			}
		}
	}
	return
}

// DockerImage configures goreleaser to build a container image.
// https://goreleaser.com/customization/docker/
func DockerImage(imagePrefixes []string, dist, arch, armVersion string) config.Docker {
	dockerArchName := archName(arch, armVersion)
	var imageTemplates []string
	for _, prefix := range imagePrefixes {
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

	return config.Docker{
		ImageTemplates: imageTemplates,
		Dockerfile:     path.Join("distributions", dist, "Dockerfile"),

		Use: "buildx",
		BuildFlagTemplates: []string{
			"--pull",
			fmt.Sprintf("--platform=linux/%s", dockerArchName),
			label("created", ".Date"),
			label("name", ".ProjectName"),
			label("revision", ".FullCommit"),
			label("version", ".Version"),
			label("source", ".GitURL"),
		},
		Files:  []string{path.Join("configs", fmt.Sprintf("%s.yaml", dist))},
		Goos:   "linux",
		Goarch: arch,
		Goarm:  armVersion,
	}
}

func DockerManifests(imagePrefixes, dists []string) (r []config.DockerManifest) {
	for _, dist := range dists {
		for _, prefix := range imagePrefixes {
			r = append(r, DockerManifest(prefix, `{{ .Version }}`, dist))
			r = append(r, DockerManifest(prefix, "latest", dist))
		}
	}
	return
}

// DockerManifest configures goreleaser to build a multi-arch container image manifest.
// https://goreleaser.com/customization/docker_manifest/
func DockerManifest(prefix, version, dist string) config.DockerManifest {
	var imageTemplates []string
	for _, arch := range Architectures {
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
	return strings.Replace(dist, "otelcol", "axoflow-otel-collector", 1)
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
