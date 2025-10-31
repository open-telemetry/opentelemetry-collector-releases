// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"log"
	"os"

	"github.com/open-telemetry/opentelemetry-collector-releases/cmd/goreleaser/internal"
	"go.yaml.in/yaml/v3"
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

	e := yaml.NewEncoder(os.Stdout)
	e.SetIndent(2)
	if err := e.Encode(&project); err != nil {
		log.Fatal(err)
	}
}
