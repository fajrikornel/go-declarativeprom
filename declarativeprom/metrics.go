package declarativeprom

import (
	"github.com/prometheus/client_golang/prometheus"
)

var registerer *prometheus.Registerer
var gatherer *prometheus.Gatherer
var metricsMap map[string]prometheus.Collector

type MetricType interface {
	GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector
}

type mappedMetric struct {
	MetricName               string
	Help                     string
	PrometheusLabels         prometheus.Labels
	CollectorInitializerFunc func() prometheus.Collector
}

func init() {
	registerer = &prometheus.DefaultRegisterer
	gatherer = &prometheus.DefaultGatherer
	metricsMap = make(map[string]prometheus.Collector)
}

func SetRegisterer(customRegisterer prometheus.Registerer) {
	registerer = &customRegisterer
}

func SetGatherer(customGatherer prometheus.Gatherer) {
	gatherer = &customGatherer
}

func GetRegisterer() prometheus.Registerer {
	return *registerer
}

func GetGatherer() prometheus.Gatherer {
	return *gatherer
}
