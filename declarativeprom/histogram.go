package declarativeprom

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Histogram struct{}

// GetCollectorInitializer implementation of Histogram. Uses prometheus.HistogramVec for multi-label capability.
func (h Histogram) GetCollectorInitializer(name, help string, labels []string) func() prometheus.Collector {
	return func() prometheus.Collector {
		return prometheus.NewHistogramVec(
			prometheus.HistogramOpts{Name: name, Help: help},
			labels,
		)
	}
}

// RecordHistogram records the declarative Histogram metric struct with the received float64 value.
func RecordHistogram(recordedMetric interface{}, value float64) {
	extractedMetric := extractToMappedMetric(recordedMetric)

	h := getOrRegisterMetric(extractedMetric)
	hVec := h.(*prometheus.HistogramVec)

	hVec.With(extractedMetric.PrometheusLabels).Observe(value)
}

// NewTimer returns a *prometheus.Timer that can be used for timing metrics for the received declarative Histogram metric struct.
func NewTimer(recordedMetric interface{}) *prometheus.Timer {
	extractedMetric := extractToMappedMetric(recordedMetric)

	h := getOrRegisterMetric(extractedMetric)
	hVec := h.(*prometheus.HistogramVec)

	return prometheus.NewTimer(hVec.With(extractedMetric.PrometheusLabels))
}
