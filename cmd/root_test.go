package cmd_test

import (
	"testing"

	"github.com/brpaz/prom-dirsize-exporter/cmd"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRootCmd(t *testing.T) {
	rootCmd := cmd.NewRootCmd(zap.NewNop())

	assert.NotNil(t, rootCmd)
	assert.Equal(t, "prom-dirsize-exporter", rootCmd.Use)
}
