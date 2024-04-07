package cmd_test

import (
	"testing"

	"github.com/brpaz/prom-dirsize-exporter/cmd"
	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	versionCmd := cmd.NewVersionCmd()
	assert.NotNil(t, versionCmd)
}
