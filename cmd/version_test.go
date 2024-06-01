package cmd_test

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"

	"github.com/brpaz/prom-dirsize-exporter/cmd"
	"github.com/brpaz/prom-dirsize-exporter/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	mockVersion   = "v0.1.0"
	mockGitCommit = "123456"
	mockBuildDate = "2021-09-26T22:00:00Z"
)

func TestVersionCmd(t *testing.T) {
	t.Parallel()

	version.BuildDate = mockBuildDate
	version.GitCommit = mockGitCommit
	version.Version = mockVersion

	versionCmd := cmd.NewVersionCmd()

	// Capture the output
	buf := new(bytes.Buffer)

	versionCmd.SetOut(buf)
	assert.NotNil(t, versionCmd)

	err := versionCmd.RunE(versionCmd, []string{})
	require.NoError(t, err)

	expectedOutput := fmt.Sprintf(
		"Version: %s\nGit commit: %s\nBuild date: %s\nGo version: %s\n",
		mockVersion, mockGitCommit, mockBuildDate, runtime.Version(),
	)

	assert.Equal(t, expectedOutput, buf.String())
}
