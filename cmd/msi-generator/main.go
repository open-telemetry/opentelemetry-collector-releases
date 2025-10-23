// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

const (
	coreDistro       = "otelcol"
	contribDistro    = "otelcol-contrib"
	otlpDistro       = "otelcol-otlp"
	templateFilename = "cmd/msi-generator/windows-installer.wxs.tmpl"
	finalFilename    = "windows-installer.wxs"
	distroFolder     = "distributions"
)

var (
	distFlag = flag.String("d", "", "Collector distributions to build")
)

func main() {
	flag.Parse()

	if len(*distFlag) == 0 {
		log.Fatal("no distribution to template")
	}
	distros := strings.Split(*distFlag, ",")

	for _, distro := range distros {
		log.Println("Templating MSI installer for distribution: " + distro)
		TemplateDist(distro)
	}
}

func TemplateDist(dist string) {
	switch dist {
	case coreDistro, contribDistro:
		templateDist(dist, true)
	case otlpDistro:
		templateDist(dist, false)
	default:
		log.Println("Unknown distribution: " + dist)
	}
}

func templateDist(dist string, addConfig bool) {
	// Parse the base template
	baseTemplate, err := template.New("base").Delims("<<", ">>").ParseFiles(templateFilename)
	if err != nil {
		panic(err)
	}

	// Data for the base template
	data := map[string]interface{}{
		"AddConfig": addConfig,
	}

	// Execute the base template to generate a new template
	var generatedTemplateContent bytes.Buffer
	err = baseTemplate.ExecuteTemplate(&generatedTemplateContent, "base", data)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s/%s", distroFolder, dist, finalFilename), generatedTemplateContent.Bytes(), 0644)
	if err != nil {
		panic(err)
	}
}
