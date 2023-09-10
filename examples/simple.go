package main

import (
	dprom "github.com/fajrikornel/go-declarativeprom/declarativeprom"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	// create custom registry
	registry := prometheus.NewRegistry()

	// set the library to use the new registry for its registerer and gatherer
	dprom.SetRegisterer(registry)
	dprom.SetGatherer(registry)

	go invokeSomeMetrics()

	//use the GetGatherer and GetRegisterer functions for promhttp.HandlerFor function
	http.Handle("/metrics",
		promhttp.HandlerFor(dprom.GetGatherer(),
			promhttp.HandlerOpts{Registry: dprom.GetRegisterer()},
		))

	log.Println("inspect me at localhost:8080/metrics")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//
// Defining the metrics declaratively:
// embed the metric type, provide the help description,
// and add the labels as the struct fields
//

type NumberOfHits struct {
	dprom.Counter  `help:"Counter that tracks how many function calls were made"`
	Method         string
	SomeOtherLabel string
}

type HitDuration struct {
	dprom.Histogram `help:"Histogram that tracks duration of the function calls"`
	Method          string
	SomeOtherLabel  float64
}

//
// example functions hello() and hi() that uses the metrics
//

func hello() {
	log.Println("hello() invoked along with metrics")

	// use NewTimer function with the declared HitDuration metric struct, specifying the label
	timer := dprom.NewTimer(HitDuration{
		Method:         "hello()",
		SomeOtherLabel: 123.123,
	})
	defer timer.ObserveDuration() // use ObserveDuration like any other prometheus.Timer variable

	time.Sleep(250 * time.Millisecond)

	// increment the declared NumberOfHits metric and specifying the label
	dprom.IncrementCounter(NumberOfHits{
		Method:         "hello()",
		SomeOtherLabel: "someOtherLabelHello",
	})
}

func hi() {
	log.Println("hi() invoked along with metrics")

	// do the same as hello() function, except we can specify different labels according to use case
	timer := dprom.NewTimer(HitDuration{
		Method:         "hi()",
		SomeOtherLabel: 321.321,
	})
	defer timer.ObserveDuration()

	time.Sleep(750 * time.Millisecond)

	// do the same as hello() function, except we can specify different labels according to use case
	dprom.IncrementCounter(NumberOfHits{
		Method:         "hi()",
		SomeOtherLabel: "someOtherLabelHi",
	})
}

// function that randomly invokes hello() and hi() functions for metric testing
func invokeSomeMetrics() {
	for true {
		i := rand.Intn(2)
		if i == 0 {
			hi()
		} else {
			hello()
		}
	}
}
