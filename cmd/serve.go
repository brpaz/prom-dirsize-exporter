package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/brpaz/prom-dirsize-exporter/internal/collector"
	"github.com/brpaz/prom-dirsize-exporter/internal/server"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	serveFlagMetricsPort   = "metrics-port"
	serveFlagMetricsPath   = "metrics-path"
	serviceFlagDirectories = "directories"
)

// SetFlagsFromEnv sets the command flags from environment variables.
// The environment variables take precedence over any defined flag.
func SetFlagsFromEnv(cmd *cobra.Command) error {
	if !cmd.Flags().Changed(serveFlagMetricsPort) && os.Getenv("METRICS_PORT") != "" {
		if err := cmd.Flags().Set(serveFlagMetricsPort, os.Getenv("METRICS_PORT")); err != nil {
			return err
		}
	}

	if !cmd.Flags().Changed(serveFlagMetricsPath) && os.Getenv("METRICS_PATH") != "" {
		if err := cmd.Flags().Set(serveFlagMetricsPath, os.Getenv("METRICS_PATH")); err != nil {
			return err
		}
	}

	if !cmd.Flags().Changed("directories") && os.Getenv("DIRECTORIES") != "" {
		if err := cmd.Flags().Set(serviceFlagDirectories, os.Getenv("DIRECTORIES")); err != nil {
			return err
		}
	}

	return nil
}

// NewServeCmd returns a new instance of the serve command that will start the metrics http server
func NewServeCmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "Starts the prometheus exporter",
		Example: `prom-dirsize-exporter serve --metricsPort 8080 --metricsPath /metrics --directories /var/log:/var/tmp`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := SetFlagsFromEnv(cmd); err != nil {
				return fmt.Errorf("error setting flags from environment variables: %w", err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dirsList, err := cmd.Flags().GetString(serviceFlagDirectories)
			if err != nil {
				return fmt.Errorf("error reading directories flag: %w", err)
			}

			metricsPath, err := cmd.Flags().GetString(serveFlagMetricsPath)
			if err != nil {
				return fmt.Errorf("error reading metricsPath flag: %w", err)
			}

			metricsPort, err := cmd.Flags().GetInt(serveFlagMetricsPort)
			if err != nil {
				return fmt.Errorf("error reading metricsPort flag: %w", err)
			}

			directoriesToMonitor := make([]string, 0)
			directoriesToMonitor = append(directoriesToMonitor, filepath.SplitList(dirsList)...)
			return runServer(logger, directoriesToMonitor, metricsPort, metricsPath)
		},
	}

	cmd.PersistentFlags().IntP("metrics-port", "p", server.DefaultMetricsPort, "the port where the metrics server will listen")
	cmd.PersistentFlags().StringP("metrics-path", "m", server.DefaultMetricsPath, "the path where the metrics will be exposed")
	cmd.PersistentFlags().StringP("directories", "d", "", "a colon separated list of directories to monitor")

	return cmd
}

func runServer(logger *zap.Logger, directoriesToMonitor []string, metricsPort int, metricsPath string) error {
	// Initialize and register collector
	dirsizeCollector := collector.NewDirectoryCollector(
		collector.WithLogger(logger),
		collector.WithDirectories(directoriesToMonitor),
	)

	prometheus.MustRegister(dirsizeCollector)

	// Create metrics server
	metricsServer := server.NewMetricsServer(
		server.WithLogger(logger),
		server.WithPort(metricsPort),
		server.WithPath(metricsPath),
	)

	return metricsServer.Start()
}
