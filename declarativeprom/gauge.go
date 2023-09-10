package declarativeprom

import "github.com/prometheus/client_golang/prometheus"

type Gauge struct{}

// GetCollectorInitializer implementation of Gauge. Uses prometheus.GaugeVec for multi-label capability.
func (c Gauge) GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector {
	return func() prometheus.Collector {
		return prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: name, Help: help},
			labels,
		)
	}
}

// SetGauge sets the declarative Gauge metric struct with the received float64 value.
func SetGauge(recordedMetric interface{}, value float64) {
	extractedMetric := extractToMappedMetric(recordedMetric)

	g := getOrRegisterMetric(extractedMetric)
	gVec := g.(*prometheus.GaugeVec)

	gVec.With(extractedMetric.PrometheusLabels).Set(value)
}
