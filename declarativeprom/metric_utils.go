package declarativeprom

import (
	"fmt"
	"github.com/gobeam/stringy"
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
)

const Help = "help"

func extractToMappedMetric(recordedMetric interface{}) mappedMetric {
	metricName := stringy.New(reflect.TypeOf(recordedMetric).Name()).SnakeCase().ToLower()
	var metricType MetricType
	var help string
	labels := make([]string, 0)
	prometheusLabels := prometheus.Labels{}

	recordedMetricReflect := reflect.ValueOf(recordedMetric)

	if !recordedMetricReflect.Type().Implements(reflect.TypeOf((*MetricType)(nil)).Elem()) {
		panic("METRIC_DOES_NOT_IMPLEMENT_METRIC_TYPE")
	}

	for i := 0; i < recordedMetricReflect.NumField(); i++ {
		fieldReflect := recordedMetricReflect.Type().Field(i)
		fieldNameReflect := fieldReflect.Name
		fieldTypeReflect := fieldReflect.Type
		fieldValueReflect := recordedMetricReflect.Field(i).Interface()

		if fieldTypeReflect.Implements(reflect.TypeOf((*MetricType)(nil)).Elem()) {
			metricType = fieldValueReflect.(MetricType)
			help = fieldReflect.Tag.Get(Help)

		} else {
			labelName := stringy.New(fieldNameReflect).SnakeCase().ToLower()
			labels = append(labels, labelName)
			prometheusLabels[labelName] = fmt.Sprintf("%v", recordedMetricReflect.Field(i))

		}
	}

	return mappedMetric{
		MetricName:               metricName,
		Help:                     help,
		PrometheusLabels:         prometheusLabels,
		CollectorInitializerFunc: metricType.GetCollectorInitializer(metricName, help, labels),
	}
}

func getOrRegisterMetric(extractedMetric mappedMetric) prometheus.Collector {
	registeredMetric := metricsMap[extractedMetric.MetricName]
	if registeredMetric == nil {
		collector := extractedMetric.CollectorInitializerFunc()

		err := (*registerer).Register(collector)
		if err != nil {
			panic(err)
		}

		metricsMap[extractedMetric.MetricName] = collector
		registeredMetric = collector
	}

	return registeredMetric
}
