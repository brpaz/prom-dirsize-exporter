// Package collector defines the main Prometheus collector for directory size metrics.
package collector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

const (
	CollectorNamespace = "directory"
	CollectorName      = "size_bytes"
)

// DirectoryCollector collects directory size metrics
type DirectoryCollector struct {
	logger      *zap.Logger
	directories []string
	mutex       sync.Mutex
	metricsMap  map[string]prometheus.Gauge
}

// DirectoryCollectorOption represents an option to customize DirectoryCollector behavior
type DirectoryCollectorOption func(*DirectoryCollector)

// WithDirectories sets the directories to monitor
func WithDirectories(dirs []string) DirectoryCollectorOption {
	return func(c *DirectoryCollector) {
		c.directories = dirs
	}
}

// WithLogger sets the logger of the DirectoryCollector
func WithLogger(logger *zap.Logger) DirectoryCollectorOption {
	return func(c *DirectoryCollector) {
		c.logger = logger
	}
}

// NewDirectoryCollector creates a new DirectoryCollector with the provided options
func NewDirectoryCollector(opts ...DirectoryCollectorOption) *DirectoryCollector {
	collector := &DirectoryCollector{
		logger:     zap.NewNop(),
		metricsMap: make(map[string]prometheus.Gauge),
	}

	// Apply options
	for _, opt := range opts {
		opt(collector)
	}

	return collector
}

// Directories returns the directories being monitored
func (c *DirectoryCollector) Directories() []string {
	return c.directories
}

// Describe implements the prometheus.Collector interface.
func (c *DirectoryCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metricsMap {
		metric.Describe(ch)
	}
}

// Collect implements the prometheus.Collector interface.
func (c *DirectoryCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	c.logger.Info("start collector", zap.String("directories", strings.Join(c.directories, ",")))

	for _, dir := range c.directories {

		// before processing, check if the directory exists
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			c.logger.Error("directory does not exist", zap.String("directory", dir))
			continue
		}

		wg.Add(1)
		go func(directory string) {
			c.logger.Info("collecting directory size", zap.String("directory", directory))
			defer wg.Done()

			size, err := c.getDirectorySize(directory)
			if err != nil {
				c.logger.Error("error getting directory size", zap.String("directory", directory), zap.Error(err))
				return
			}

			c.logger.Info("directory size collected", zap.String("directory", directory), zap.Int64("size", size))

			c.updateMetric(directory, size, ch)
		}(dir)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// updateMetric updates or creates a new metric for the given directory.
func (c *DirectoryCollector) updateMetric(directory string, size int64, ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	metric, ok := c.metricsMap[directory]
	if !ok {
		metric = prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace:   CollectorNamespace,
			Name:        CollectorName,
			Help:        "Size of the directory in bytes.",
			ConstLabels: prometheus.Labels{"name": filepath.Base(directory), "path": directory},
		})
		c.metricsMap[directory] = metric
	}
	metric.Set(float64(size))
	ch <- metric
}

// getDirectorySize calculates the total size of a directory using the "du" command
func (c *DirectoryCollector) getDirectorySize(path string) (int64, error) {
	cmd := exec.Command("du", "-sb", path)

	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0o644)
	if err != nil {
		c.logger.Error("error opening /dev/null", zap.Error(err))
	}
	defer devNull.Close()

	cmd.Stderr = devNull

	// TODO Ignoring the error is not pretty, but du always return error exit status,
	// even for not "fatal" errors like permission denied
	output, _ := cmd.Output()

	sizeStr := strings.Fields(string(output))[0]
	var size int64
	_, err = fmt.Sscanf(sizeStr, "%d", &size)
	if err != nil {
		return 0, fmt.Errorf("error parsing size: %w", err)
	}

	return size, nil
}
