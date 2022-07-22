//go:build releaser

// Deprecated: will be removed in next minor version.
package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/open-telemetry/opentelemetry-collector-releases/internal/goreleaser"
)

var (
	// Deprecated: will be removed in next minor version.
	ImagePrefixes = goreleaser.ImagePrefixes
	// Deprecated: will be removed in next minor version.
	Architectures = goreleaser.Architectures

	distsFlag = flag.String("d", "", "Collector distributions(s) to build, comma-separated")
)

func main() {
	flag.Parse()

	log.Print("DEPRECATED: will be removed in next minor version")

	if len(*distsFlag) == 0 {
		log.Fatal("no distributions to build")
	}
	dists := strings.Split(*distsFlag, ",")

	project := Generate(ImagePrefixes, dists)

	if err := yaml.NewEncoder(os.Stdout).Encode(&project); err != nil {
		log.Fatal(err)
	}
}

// Deprecated: will be removed in next minor version.
var Generate = goreleaser.Generate

// Deprecated: will be removed in next minor version.
var Builds = goreleaser.Builds

// Deprecated: will be removed in next minor version.
var Build = goreleaser.Build

// Deprecated: will be removed in next minor version.
var Archives = goreleaser.Archives

// Deprecated: will be removed in next minor version.
var Archive = goreleaser.Archive

// Deprecated: will be removed in next minor version.
var Packages = goreleaser.Packages

// Deprecated: will be removed in next minor version.
var Package = goreleaser.Package

// Deprecated: will be removed in next minor version.
var DockerImages = goreleaser.DockerImages

// Deprecated: will be removed in next minor version.
var DockerImage = goreleaser.DockerImage

// Deprecated: will be removed in next minor version.
var DockerManifests = goreleaser.DockerManifests

// Deprecated: will be removed in next minor version.
var DockerManifest = goreleaser.DockerManifest
