// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"slices"
	"testing"
)

func TestK8sDistributionWindowsArchitectures(t *testing.T) {
	project := BuildDistribution("otelcol-k8s", false)

	// The Windows binaries are built for every Windows architecture, as for otelcol and otelcol-contrib.
	var windowsBuilds int
	for _, build := range project.Builds {
		if !slices.Equal(build.Goos, []string{"windows"}) {
			continue
		}
		windowsBuilds++
		if !slices.Equal(build.Goarch, winArchs) {
			t.Errorf("build %q GOARCH = %v, want %v", build.ID, build.Goarch, winArchs)
		}
	}
	if windowsBuilds != 1 {
		t.Fatalf("got %d Windows builds, want 1", windowsBuilds)
	}

	// The Windows container images keep their own architecture list.
	var windowsImages int
	for _, image := range project.Dockers {
		if image.Goos != "windows" {
			continue
		}
		windowsImages++
		if !slices.Contains(winContainerArchs, image.Goarch) {
			t.Errorf("Windows container image GOARCH = %q, want one of %v", image.Goarch, winContainerArchs)
		}
	}
	if windowsImages == 0 {
		t.Error("otelcol-k8s must still produce Windows container images")
	}
}
