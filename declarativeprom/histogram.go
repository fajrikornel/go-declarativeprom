package declarativeprom

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Histogram struct{}

func (h Histogram) GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector {
	return func() prometheus.Collector {
		return prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: name, Help: help},
			labels,
		)
	}
}

func RecordHistogram(recordedMetric interface{}, value float64) {
	extractedMetric := extractToMappedMetric(recordedMetric)

	h := getOrRegisterMetric(extractedMetric)
	hVec := h.(*prometheus.HistogramVec)

	hVec.With(extractedMetric.PrometheusLabels).Observe(value)
}

func NewTimer(recordedMetric interface{}) *prometheus.Timer {
	extractedMetric := extractToMappedMetric(recordedMetric)

	h := getOrRegisterMetric(extractedMetric)
	hVec := h.(*prometheus.HistogramVec)

	return prometheus.NewTimer(hVec.With(extractedMetric.PrometheusLabels))
}
