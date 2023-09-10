package declarativeprom

import "github.com/prometheus/client_golang/prometheus"

type Gauge struct{}

func (c Gauge) GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector {
	return func() prometheus.Collector {
		return prometheus.NewGaugeVec(
			prometheus.GaugeOpts{Name: name, Help: help},
			labels,
		)
	}
}

func SetGauge(recordedMetric interface{}, value float64) {
	extractedMetric := extractToMappedMetric(recordedMetric)

	g := getOrRegisterMetric(extractedMetric)
	gVec := g.(*prometheus.GaugeVec)

	gVec.With(extractedMetric.PrometheusLabels).Set(value)
}
