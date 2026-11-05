// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"log"
	"os"

	"github.com/open-telemetry/opentelemetry-collector-releases/cmd/goreleaser/internal"
	"go.yaml.in/yaml/v4"
)

var (
	distFlag               = flag.String("d", "", "Collector distributions to build")
	contribBuildOrRestFlag = flag.Bool("generate-build-step", false, "Collector Contrib distribution only - switch between build and package config file - set to true to generate build step, false to generate package step")
)

func main() {
	flag.Parse()

	if len(*distFlag) == 0 {
		log.Fatal("no distribution to build")
	}
	project := internal.BuildDistribution(*distFlag, *contribBuildOrRestFlag)

	os.Stdout.WriteString("# yaml-language-server: $schema=https://goreleaser.com/static/schema-pro.json\n")
	d, err := yaml.NewDumper(os.Stdout, yaml.WithV3Defaults(), yaml.WithIndent(2), yaml.WithLineWidth(-1))
	if err != nil {
		log.Fatal(err)
	}
	if err := d.Dump(&project); err != nil {
		log.Fatal(err)
	}
	if err := d.Close(); err != nil {
		log.Fatal(err)
	}
}
