package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

const (
	DefaultMetricsPort = 8080
	DefaultMetricsPath = "/metrics"

	// Define constants for signal types
	sigInt  = syscall.SIGINT
	sigTerm = syscall.SIGTERM
)

type MetricsServer struct {
	logger      *zap.Logger
	httpServer  *http.Server
	port        int
	metricsPath string
}

// MetricsServerOption is a function that configures a MetricsServer
type MetricsServerOption func(*MetricsServer)

// WithLogger sets the logger of the MetricsServer
func WithLogger(logger *zap.Logger) MetricsServerOption {
	return func(c *MetricsServer) {
		c.logger = logger
	}
}

// WithPort sets the port of the MetricsServer
func WithPort(port int) MetricsServerOption {
	return func(c *MetricsServer) {
		c.port = port
	}
}

// WithPath sets the metrics path of the MetricsServer
func WithPath(metricsPath string) MetricsServerOption {
	return func(c *MetricsServer) {
		c.metricsPath = metricsPath
	}
}

// NewMetricsServer creates a new MetricsServer with the provided options.
// It uses golang http.Server to create a new server instance to expose the prometheus metrics.
func NewMetricsServer(opts ...MetricsServerOption) *MetricsServer {
	srv := &MetricsServer{
		port:        DefaultMetricsPort,
		metricsPath: DefaultMetricsPath,
		logger:      zap.NewNop(),
	}

	for _, opt := range opts {
		opt(srv)
	}

	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", srv.port),
		Handler: initRoutes(srv.metricsPath),
	}

	return srv
}

func (s *MetricsServer) Start() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, sigInt, sigTerm)

	errSrvStart := make(chan error, 1)
	go func() {
		s.logger.Info("Starting server on port", zap.Int("port", s.port))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errSrvStart <- err
		}
	}()

	// Block until a signal is received or the server fails to start
	select {
	case <-stop:
		err := s.Stop()
		return err
	case err := <-errSrvStart:
		s.logger.Error("Failed to start HTTP server", zap.Error(err))
		return err
	}
}

// Stop stops the MetricsServer gracefully
func (s *MetricsServer) Stop() error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Info("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}
