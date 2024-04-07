package collector_test

import (
	"testing"
	"time"

	"github.com/brpaz/prom-dirsize-exporter/internal/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewDirectoryCollector(t *testing.T) {
	c := collector.NewDirectoryCollector(
		collector.WithDirectories([]string{"/tmp"}),
	)

	assert.NotNil(t, c)
	assert.Equal(t, []string{"/tmp"}, c.Directories())
}

func TestDirectoryCollector_Collect(t *testing.T) {
	c := collector.NewDirectoryCollector(
		collector.WithDirectories([]string{"./testdata/example_directory"}),
	)

	// Create a new Prometheus registry
	// TODO is this the best way to test this?
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	registry.MustRegister(c)

	ch := make(chan prometheus.Metric, 1)
	defer close(ch)

	go func() {
		c.Collect(ch)
	}()

	select {
	case metric := <-ch:
		assert.NotNil(t, metric)
		assert.Implements(t, (*prometheus.Gauge)(nil), metric)

		metrics, _ := registry.Gather()
		assert.Equal(t, 1, len(metrics))
		assert.Equal(t, "directory_size_bytes", metrics[0].GetName())
		assert.Equal(t, float64(2105344), metrics[0].Metric[0].Gauge.GetValue())
	case timeout := <-time.After(1 * time.Second):
		t.Fatalf("Timed out waiting for metric to be collected. %v", timeout)
	}
}

func TestDirectoryCollector_Collect_WithNonExistingDirectory(t *testing.T) {
	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore)

	c := collector.NewDirectoryCollector(
		collector.WithLogger(logger),
		collector.WithDirectories([]string{"/tmp/some-non-existing-dir"}),
	)

	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	registry.MustRegister(c)

	ch := make(chan prometheus.Metric, 1)
	defer close(ch)

	go func() {
		c.Collect(ch)
	}()

	time.Sleep(1 * time.Second)
	assert.Len(t, ch, 0)

	// Check if the error was logged
	errorLog := observedLogs.FilterMessage("Error getting directory size").All()
	assert.Equal(t, 1, len(errorLog))
}
