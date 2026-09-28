// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"slices"
	"testing"
)

func TestPrometheusDistribution(t *testing.T) {
	project := BuildDistribution("otelcol-prometheus", false)

	wantBuilds := map[string]struct {
		goos   string
		goarch []string
	}{
		"otelcol-prometheus-linux": {
			goos:   "linux",
			goarch: []string{"amd64", "arm64"},
		},
		"otelcol-prometheus-darwin": {
			goos:   "darwin",
			goarch: []string{"arm64"},
		},
		"otelcol-prometheus-windows": {
			goos:   "windows",
			goarch: []string{"amd64"},
		},
	}

	if len(project.Builds) != len(wantBuilds) {
		t.Fatalf("got %d builds, want %d", len(project.Builds), len(wantBuilds))
	}
	for _, build := range project.Builds {
		want, ok := wantBuilds[build.ID]
		if !ok {
			t.Errorf("unexpected build %q", build.ID)
			continue
		}
		if len(build.Goos) != 1 || build.Goos[0] != want.goos {
			t.Errorf("build %q GOOS = %v, want [%s]", build.ID, build.Goos, want.goos)
		}
		if !slices.Equal(build.Goarch, want.goarch) {
			t.Errorf("build %q GOARCH = %v, want %v", build.ID, build.Goarch, want.goarch)
		}
	}

	if len(project.NFPMs) == 0 {
		t.Error("Prometheus distribution must produce deb and RPM packages")
	}
	if len(project.MSI) == 0 {
		t.Error("Prometheus distribution must produce an MSI package")
	}
	if len(project.Dockers) != 2 {
		t.Errorf("got %d container image builds, want 2 linux image builds", len(project.Dockers))
	}
	if len(project.DockerManifests) == 0 {
		t.Error("Prometheus distribution must produce multi-architecture image manifests")
	}
}
