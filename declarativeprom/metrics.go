package declarativeprom

import (
	"github.com/prometheus/client_golang/prometheus"
)

// prometheus.Registerer and prometheus.Gatherer that will be used by this library.
var registerer *prometheus.Registerer
var gatherer *prometheus.Gatherer

// metricsMap is a map of the metric name to prometheus.Collector that is registered.
var metricsMap map[string]prometheus.Collector

// MetricType is the interface which different types of metrics (Counter, Histogram, Gauge)
// needs to implement.
//
// GetCollectorInitializer should take the metrics' name, help, and label names, and return
// a function that when invoked will return the corresponding prometheus.Collector according to the metric type.
type MetricType interface {
	GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector
}

// mappedMetric is an internal struct that will be used to map user-defined declarative structs to
// its corresponding prometheus metric type.
type mappedMetric struct {
	MetricName               string
	Help                     string
	PrometheusLabels         prometheus.Labels
	CollectorInitializerFunc func() prometheus.Collector
}

// init sets the defaults.
func init() {
	registerer = &prometheus.DefaultRegisterer
	gatherer = &prometheus.DefaultGatherer
	metricsMap = make(map[string]prometheus.Collector)
}

// SetRegisterer sets the prometheus.Registerer that will be used by this library.
func SetRegisterer(customRegisterer prometheus.Registerer) {
	registerer = &customRegisterer
}

// SetGatherer sets the prometheus.Gatherer that will be used by this library.
func SetGatherer(customGatherer prometheus.Gatherer) {
	gatherer = &customGatherer
}

// GetRegisterer returns the prometheus.Registerer used by this library.
func GetRegisterer() prometheus.Registerer {
	return *registerer
}

// GetGatherer returns the prometheus.Gatherer used by this library.
func GetGatherer() prometheus.Gatherer {
	return *gatherer
}
