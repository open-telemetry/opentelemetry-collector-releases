// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

// Architecture sets shared across distributions.
var (
	baseArchs         = []string{"386", "amd64", "arm", "arm64", "ppc64le", "riscv64", "s390x"}
	winArchs          = []string{"386", "amd64", "arm64"}
	winContainerArchs = []string{"amd64"}
	darwinArchs       = []string{"amd64", "arm64"}
	k8sArchs          = []string{"amd64", "arm64", "ppc64le", "riscv64", "s390x"}
	ocbArchs          = []string{"amd64", "arm64", "ppc64le", "riscv64"}
	opAmpArchs        = []string{"amd64", "arm64", "ppc64le"}
)
