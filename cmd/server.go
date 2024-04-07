package cmd

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/brpaz/prom-dirsize-exporter/internal/collector"
	"github.com/brpaz/prom-dirsize-exporter/internal/server"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	serveFlagMetricsPort   = "metricsPort"
	serveFlagMetricsPath   = "metricsPath"
	serviceFlagDirectories = "directories"
)

// setFlagsFromEnv sets the command flags from environment variables.
// The environment variables take precedence over any defined flag.
func setFlagsFromEnv(cmd *cobra.Command) error {
	if os.Getenv("METRICS_PORT") != "" {
		return cmd.Flags().Set(serveFlagMetricsPort, os.Getenv("METRICS_PORT"))
	}

	if os.Getenv("METRICS_PATH") != "" {
		return cmd.Flags().Set(serveFlagMetricsPath, os.Getenv("METRICS_PATH"))
	}

	if os.Getenv("DIRECTORIES") != "" {
		return cmd.Flags().Set(serviceFlagDirectories, os.Getenv("DIRECTORIES"))
	}

	return nil
}

// NewServeCmd returns a new instance of the serve command that will start the metrics http server
func NewServeCmd(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Short:   "Starts the metrics http server",
		Example: `prom-dirsize-exporter serve --metricsPort 8080 --metricsPath /metrics --directories /var/log;/var/tmp`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setFlagsFromEnv(cmd); err != nil {
				return errors.Join(errors.New("error setting command flags from environment variables"), err)
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dirsList, err := cmd.Flags().GetString(serviceFlagDirectories)
			if err != nil {
				return errors.Join(errors.New("error reading directories flag"), err)
			}

			metricsPath, err := cmd.Flags().GetString(serveFlagMetricsPath)
			if err != nil {
				return errors.Join(errors.New("error reading metricsPath flag"), err)
			}

			metricsPort, err := cmd.Flags().GetInt(serveFlagMetricsPort)
			if err != nil {
				return errors.Join(errors.New("error reading metricsPort flag"), err)
			}

			directoriesToMonitor := make([]string, 0)
			directoriesToMonitor = append(directoriesToMonitor, filepath.SplitList(dirsList)...)
			return runServer(logger, directoriesToMonitor, metricsPort, metricsPath)
		},
	}

	cmd.PersistentFlags().IntP("metricsPort", "p", server.DefaultMetricsPort, "the port where the metrics server will listen")
	cmd.PersistentFlags().StringP("metricsPath", "m", server.DefaultMetricsPath, "the path where the metrics will be exposed")
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
