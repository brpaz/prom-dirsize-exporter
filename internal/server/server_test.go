package server_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/brpaz/prom-dirsize-exporter/internal/server"
	"github.com/brpaz/prom-dirsize-exporter/internal/testutil"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewMetricsServer(t *testing.T) {
	t.Parallel()

	srv := server.NewMetricsServer(
		server.WithPort(3000),
		server.WithPath("/metrics"),
		server.WithLogger(zap.NewNop()),
	)

	assert.NotNil(t, srv)
}

// TestMetricsServerStart tests the Start method of MetricsServer.
func TestMetricsServerStart_ReturnsSuccess(t *testing.T) {
	t.Parallel()

	logger := zap.NewNop()
	port, err := testutil.GetFreePort()
	if err != nil {
		t.Fatalf("Error getting free port: %s", err)
	}

	srv := server.NewMetricsServer(
		server.WithPort(port),
		server.WithLogger(logger),
	)

	// Start the MetricsServer in a separate goroutine.
	go func() {
		err := srv.Start()
		assert.NoError(t, err, "Expected no error when starting the server")
	}()

	t.Cleanup(func() {
		_ = srv.Stop()
	})

	// Wait for a short time to allow the server to start.
	time.Sleep(100 * time.Millisecond)

	// Send a test HTTP request to the server.
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check if the response status code is OK.
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(respBody), "Prometheus Directory Size Exporter is up and running")

	// Stop the server
	err = srv.Stop()
	assert.NoError(t, err)
}

// TestMetricsServerStartWithInvalidPort tests the Start method of MetricsServer with an invalid port.
func TestMetricsServerStart_WithInvalidPort_ReturnsError(t *testing.T) {
	t.Parallel()

	srv := server.NewMetricsServer(
		server.WithLogger(zap.NewNop()),
		server.WithPort(999999),
	)

	t.Cleanup(func() {
		_ = srv.Stop()
	})

	errCh := make(chan error, 1)
	timeout := 1 * time.Second
	go func() {
		err := srv.Start()
		if err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "listen tcp: address 999999: invalid port")
	case <-time.After(timeout):
		t.Fatal("Timed out waiting for error")
	}
}

func TestMetricsServer_ServesMetricsEndpoint(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name        string
		metricsPath string
	}{
		{
			name:        "With Custom metrics path",
			metricsPath: "/custom-metrics-path",
		},
		{
			name:        "With Default metrics path",
			metricsPath: "",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			metricsPath := scenario.metricsPath
			t.Parallel()
			logger := zap.NewNop()
			port, err := testutil.GetFreePort()
			if err != nil {
				t.Fatalf("Error getting free port: %s", err)
			}

			var srv *server.MetricsServer
			var metricsEndpoint string
			if metricsPath != "" {
				srv = server.NewMetricsServer(
					server.WithLogger(logger),
					server.WithPort(port),
					server.WithPath(metricsPath),
				)

				metricsEndpoint = fmt.Sprintf("http://localhost:%d%s", port, metricsPath)
			} else {
				srv = server.NewMetricsServer(
					server.WithPort(port),
					server.WithLogger(logger),
				)
				metricsEndpoint = fmt.Sprintf("http://localhost:%d/metrics", port)
			}

			// Start the MetricsServer in a separate goroutine.
			go func() {
				err := srv.Start()
				assert.NoError(t, err, "Expected no error when starting the server")
			}()

			t.Cleanup(func() {
				_ = srv.Stop()
			})

			// Wait for a short time to allow the server to start.
			time.Sleep(100 * time.Millisecond)

			// Send a test HTTP request to the server.
			resp, err := http.Get(metricsEndpoint)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check if the response status code is OK.
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Contains(t, string(respBody), "go_gc_duration_seconds")
		})
	}
}
