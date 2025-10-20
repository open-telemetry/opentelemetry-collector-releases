// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package msi

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

// Test structure for MSI installation tests
type msiTest struct {
	name                 string
	collectorServiceArgs string
	skipSvcStop          bool
}

func TestMSI(t *testing.T) {
	msiInstallerPath := getInstallerPath(t)

	tests := []msiTest{
		{
			name: "default",
		},
		{
			name:                 "custom",
			collectorServiceArgs: "--config " + quotedIfRequired(getAlternateConfigFile(t)),
			skipSvcStop:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runMsiTest(t, tt, msiInstallerPath)
		})
	}
}

func runMsiTest(t *testing.T, test msiTest, msiInstallerPath string) {
	// Build the MSI installation arguments and include the MSI properties map.
	installLogFile := filepath.Join(os.TempDir(), "install.log")
	args := []string{"/i", msiInstallerPath, "/qn", "/l*v", installLogFile}

	serviceArgs := quotedIfRequired(test.collectorServiceArgs)
	if test.collectorServiceArgs != "" {
		args = append(args, "COLLECTOR_SVC_ARGS="+serviceArgs)
	}

	// Run the MSI installer
	installCmd := exec.Command("msiexec")

	// msiexec is one of the noticeable exceptions about how to format the parameters,
	// see https://pkg.go.dev/os/exec#Command, so we need to join the args manually.
	cmdLine := strings.Join(args, " ")
	installCmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "msiexec " + cmdLine}
	err := installCmd.Run()
	if err != nil {
		logText, _ := os.ReadFile(installLogFile)
		t.Log(string(logText))
	}
	t.Logf("Install command: %s", installCmd.SysProcAttr.CmdLine)
	require.NoError(t, err, "Failed to install the MSI: %v\nArgs: %v", err, args)

	defer func() {
		// Uninstall the MSI
		uninstallCmd := exec.Command("msiexec")
		uninstallCmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: "msiexec /x " + msiInstallerPath + " /qn"}
		err := uninstallCmd.Run()
		t.Logf("Uninstall command: %s", uninstallCmd.SysProcAttr.CmdLine)
		require.NoError(t, err, "Failed to uninstall the MSI: %v", err)
	}()

	// Verify the service
	scm, err := mgr.Connect()
	require.NoError(t, err)
	defer scm.Disconnect()

	collectorSvcName := getServiceName(t)
	service, err := scm.OpenService(collectorSvcName)
	require.NoError(t, err)
	defer service.Close()

	// Wait for the service to reach the running state
	require.Eventually(t, func() bool {
		status, err := service.Query()
		require.NoError(t, err)
		return status.State == svc.Running
	}, 10*time.Second, 500*time.Millisecond, "Failed to start the service")

	if !test.skipSvcStop {
		defer func() {
			_, err = service.Control(svc.Stop)
			require.NoError(t, err)

			require.Eventually(t, func() bool {
				status, err := service.Query()
				require.NoError(t, err)
				return status.State == svc.Stopped
			}, 10*time.Second, 500*time.Millisecond, "Failed to stop the service")
		}()
	}

	assertServiceCommand(t, collectorSvcName, serviceArgs)
}

func assertServiceCommand(t *testing.T, serviceName, collectorServiceArgs string) {
	// Verify the service command
	actualCommand := getServiceCommand(t, serviceName)
	expectedCommand := expectedServiceCommand(t, serviceName, collectorServiceArgs)
	assert.Equal(t, expectedCommand, actualCommand)
}

func getServiceCommand(t *testing.T, serviceName string) string {
	scm, err := mgr.Connect()
	require.NoError(t, err)
	defer scm.Disconnect()

	service, err := scm.OpenService(serviceName)
	require.NoError(t, err)
	defer service.Close()

	config, err := service.Config()
	require.NoError(t, err)

	return config.BinaryPathName
}

func expectedServiceCommand(t *testing.T, serviceName, collectorServiceArgs string) string {
	programFilesDir := os.Getenv("PROGRAMFILES")
	require.NotEmpty(t, programFilesDir, "PROGRAMFILES environment variable is not set")

	collectorDir := filepath.Join(programFilesDir, "OpenTelemetry Collector")
	collectorExe := filepath.Join(collectorDir, serviceName) + ".exe"

	if collectorServiceArgs == "" {
		collectorServiceArgs = "--config " + quotedIfRequired(filepath.Join(collectorDir, "config.yaml"))
	} else {
		// Remove any quotation added for the msiexec command line
		collectorServiceArgs = strings.Trim(collectorServiceArgs, "\"")
		collectorServiceArgs = strings.ReplaceAll(collectorServiceArgs, "\"\"", "\"")
	}

	return quotedIfRequired(collectorExe) + " " + collectorServiceArgs
}

func getServiceName(t *testing.T) string {
	serviceName := os.Getenv("MSI_TEST_COLLECTOR_SERVICE_NAME")
	require.NotEmpty(t, serviceName, "MSI_TEST_COLLECTOR_SERVICE_NAME environment variable is not set")
	return serviceName
}

func getInstallerPath(t *testing.T) string {
	msiInstallerPath := os.Getenv("MSI_TEST_COLLECTOR_PATH")
	require.NotEmpty(t, msiInstallerPath, "MSI_TEST_COLLECTOR_PATH environment variable is not set")
	_, err := os.Stat(msiInstallerPath)
	require.NoError(t, err)
	return msiInstallerPath
}

func getAlternateConfigFile(t *testing.T) string {
	alternateConfigFile := os.Getenv("MSI_TEST_ALTERNATE_CONFIG_FILE")
	require.NotEmpty(t, alternateConfigFile, "MSI_TEST_ALTERNATE_CONFIG_FILE environment variable is not set")
	_, err := os.Stat(alternateConfigFile)
	require.NoError(t, err)
	return alternateConfigFile
}

func quotedIfRequired(s string) string {
	if strings.Contains(s, "\"") || strings.Contains(s, " ") {
		s = strings.ReplaceAll(s, "\"", "\"\"")
		return "\"" + s + "\""
	}
	return s
}
