package declarativeprom

import "github.com/prometheus/client_golang/prometheus"

type Counter struct{}

func (c Counter) GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector {
	return func() prometheus.Collector {
		return prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: name, Help: help},
			labels,
		)
	}
}

func IncrementCounter(recordedMetric interface{}) {
	extractedMetric := extractToMappedMetric(recordedMetric)

	c := getOrRegisterMetric(extractedMetric)
	cVec := c.(*prometheus.CounterVec)

	cVec.With(extractedMetric.PrometheusLabels).Inc()
}
