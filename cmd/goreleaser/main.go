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

package main

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/open-telemetry/opentelemetry-collector-releases/cmd/goreleaser/internal"
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
	project := internal.BuildDist(*distFlag, *contribBuildOrRestFlag)

	e := yaml.NewEncoder(os.Stdout)
	e.SetIndent(2)
	if err := e.Encode(&project); err != nil {
		log.Fatal(err)
	}
}
