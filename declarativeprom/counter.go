package declarativeprom

import "github.com/prometheus/client_golang/prometheus"

type Counter struct{}

// GetCollectorInitializer implementation of Counter. Uses prometheus.CounterVec for multi-label capability.
func (c Counter) GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector {
	return func() prometheus.Collector {
		return prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: name, Help: help},
			labels,
		)
	}
}

// IncrementCounter receives the declarative Counter metric struct and increments the counter metric.
func IncrementCounter(recordedMetric interface{}) {
	extractedMetric := extractToMappedMetric(recordedMetric)

	c := getOrRegisterMetric(extractedMetric)
	cVec := c.(*prometheus.CounterVec)

	cVec.With(extractedMetric.PrometheusLabels).Inc()
}
