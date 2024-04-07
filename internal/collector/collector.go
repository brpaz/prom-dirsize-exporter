package collector

import (
	"fmt"
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

	for _, dir := range c.directories {
		wg.Add(1)
		go func(directory string) {
			c.logger.Info("Collecting directory size", zap.String("directory", directory))
			defer wg.Done()

			size, err := c.getDirectorySize(directory)
			if err != nil {
				c.logger.Error("Error getting directory size", zap.String("directory", directory), zap.Error(err))
				return
			}

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

// getDirectorySize calculates the total size of a directory
func (c *DirectoryCollector) getDirectorySize(path string) (int64, error) {
	cmd := exec.Command("du", "-s", path)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	sizeStr := strings.Fields(string(output))[0]
	var size int64
	_, err = fmt.Sscanf(sizeStr, "%d", &size)
	if err != nil {
		return 0, err
	}

	return size * 1024, nil // du returns size in 1K blocks
}
