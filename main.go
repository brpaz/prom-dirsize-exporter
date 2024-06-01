// Package main provides the entry point for the application.
// It is responsible for initializing the global dependencies like logging and executing the root command.
package main

import (
	"fmt"
	"os"

	"github.com/brpaz/prom-dirsize-exporter/cmd"
	"go.uber.org/zap"
)

func main() {
	logger, err := createLogger()
	if err != nil {
		panic(fmt.Errorf("error creating logger: %w", err))
	}

	defer func() {
		_ = logger.Sync()
	}()

	if err := cmd.NewRootCmd(logger).Execute(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// createLogger creates a new zap logger based on the application envrionment
func createLogger() (*zap.Logger, error) {
	appEnv := os.Getenv("APP_ENV")

	if appEnv == "dev" {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}
