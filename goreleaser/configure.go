//go:build releaser

package main

// This file is a script which generates the .goreleaser.yaml file for all
// supported OpenTelemetry Collector distributions.
//
// Run it with `make generate-goreleaser`.

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goreleaser/goreleaser/pkg/config"
	"github.com/goreleaser/nfpm/v2/files"
	yaml "gopkg.in/yaml.v2"
)

var (
	ImagePrefixes = []string{"otel"}
	Architectures = []string{"386", "amd64", "arm64"}

	distsFlag = flag.String("d", "", "Collector distributions(s) to build, comma-separated")
)

func main() {
	flag.Parse()

	if len(*distsFlag) == 0 {
		log.Fatal("no distributions to build")
	}
	dists := strings.Split(*distsFlag, ",")

	tree := os.DirFS(".")
	imageNames, err := readImageNames(tree)
	if err != nil {
		log.Fatal(err)
	}

	project := Generate(ImagePrefixes, imageNames, dists)

	if err := yaml.NewEncoder(os.Stdout).Encode(&project); err != nil {
		log.Fatal(err)
	}
}

func Generate(imagePrefixes []string, imageNames map[string]string, dists []string) config.Project {
	return config.Project{
		ProjectName: "opentelemetry-collector-releases",
		Checksum: config.Checksum{
			NameTemplate: "{{ .ProjectName }}_checksums.txt",
		},

		Builds:          Builds(dists),
		Archives:        Archives(dists),
		NFPMs:           Packages(dists),
		Dockers:         DockerImages(imagePrefixes, dists, imageNames),
		DockerManifests: DockerManifests(imagePrefixes, dists, imageNames),
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
		ID:      dist,
		Dir:     path.Join("distributions", dist, "_build"),
		Binary:  dist,
		Env:     []string{"CGO_ENABLED=0"},
		Flags:   []string{"-trimpath"},
		Ldflags: []string{"-s", "-w"},

		Goos:   []string{"darwin", "linux", "windows"},
		Goarch: Architectures,
		Ignore: []config.IgnoredBuild{
			{Goos: "windows", Goarch: "arm64"},
		},
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

func DockerImages(imagePrefixes, dists []string, imageNames map[string]string) (r []config.Docker) {
	for _, dist := range dists {
		for _, arch := range Architectures {
			r = append(r, DockerImage(imagePrefixes, imageNames[dist], dist, arch))
		}
	}
	return
}

// DockerImage configures goreleaser to build a container image.
// https://goreleaser.com/customization/docker/
func DockerImage(imagePrefixes []string, imageName, dist, arch string) config.Docker {
	var imageTemplates []string
	for _, prefix := range imagePrefixes {
		imageTemplates = append(
			imageTemplates,
			fmt.Sprintf("%s/%s:{{ .Version }}-%s", prefix, imageName, arch),
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
			fmt.Sprintf("--platform=linux/%s", arch),
			label("created", ".Date"),
			label("name", ".ProjectName"),
			label("revision", ".FullCommit"),
			label("version", ".Version"),
			label("source", ".GitURL"),
		},
		Files:  []string{path.Join("configs", fmt.Sprintf("%s.yaml", dist))},
		Goos:   "linux",
		Goarch: arch,
	}
}

func DockerManifests(imagePrefixes, dists []string, imageNames map[string]string) (r []config.DockerManifest) {
	for _, dist := range dists {
		r = append(r, DockerManifest(imagePrefixes, imageNames[dist])...)
	}
	return
}

// DockerManifest configures goreleaser to build a multi-arch container image manifest.
// https://goreleaser.com/customization/docker_manifest/
func DockerManifest(imagePrefixes []string, imageName string) (manifests []config.DockerManifest) {
	for _, prefix := range imagePrefixes {
		var imageTemplates []string
		for _, arch := range Architectures {
			imageTemplates = append(
				imageTemplates,
				fmt.Sprintf("%s/%s:{{ .Version }}-%s", prefix, imageName, arch),
			)
		}

		manifests = append(manifests, config.DockerManifest{
			NameTemplate:   fmt.Sprintf("%s/%s:{{ .Version }}", prefix, imageName),
			ImageTemplates: imageTemplates,
		})
	}
	return
}

// readImageNames returns a mapping from distribution name to the container
// image name for that distribution.
func readImageNames(tree fs.FS) (map[string]string, error) {
	entries, err := fs.ReadDir(tree, "distributions")
	if err != nil {
		return nil, err
	}

	names := make(map[string]string, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		imageName, err := fs.ReadFile(tree, filepath.Join("distributions", name, "image-name"))
		if err != nil {
			return nil, fmt.Errorf("read image name for %s: %w", name, err)
		}

		names[name] = string(imageName)
	}

	return names, nil
}
