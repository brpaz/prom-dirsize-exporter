package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewRootCmd returns a new instance of the root command for the application
func NewRootCmd(logger *zap.Logger) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "prom-dirsize-exporter",
		Short: "Prometheus directory size exporter",
		Long: `Prometheus directory size exporter is a tool that exports the size of directories to Prometheus.
			See https://github.com/brpaz/prom-dirsize-exporter for more information.
		`,
	}

	// Reggister subcommands
	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(NewServeCmd(logger))

	return rootCmd
}
