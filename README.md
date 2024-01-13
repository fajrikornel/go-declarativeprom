# DeclarativeProm - Prometheus client wrapper for Go, declaratively define and manage metrics

DeclarativeProm is a lightweight Prometheus client wrapper for Go, which enables declarative definition and managing of Prometheus metrics.

## Installation

```shell
go get github.com/fajrikornel/go-declarativeprom
```

## Why use this library?

DeclarativeProm provides a simple interface for developers to instrument their applications using Prometheus. It supports:

- Using simple struct types as a declarative definition of Prometheus metrics.
- No need to manage the lifecycle of metrics; registering metrics, passing around collectors, etc.

Using a declarative approach to metrics will make app instrumentation a lot simpler. Below usage section will illustrate how.

Define the app's prometheus metrics declaratively:
```go
import dprom "github.com/fajrikornel/go-declarativeprom/declarativeprom"

type ApiResponseTime struct {
    dprom.Histogram `help:"Histogram that tracks API response times from the server"`
    Path            string
    Method          string
}

type ApiHit struct {
    dprom.Counter `help:"Counter that tracks API hits to the server"`
    Path          string
    Method        string
    ResponseCode  int
}

type SomeQueueSize struct {
    dprom.Gauge `help:"Gauge that tracks the size of some queue"`
    QueueName   string
}
```

Whenever we want to do something with the metric, just use the provided functions with the corresponding struct:
```go
dprom.IncrementCounter( // use the IncrementCounter function with the declared metric struct
    ApiHit{ // specify the labels for this particular incrementation
        Path:   "/test",
        Method: "GET",
        ResponseCode: 200,
    })
```

Other lifecycle management is done by the library. No need to manage registering the metrics, passing the prometheus collectors around in a variable, etc. Once imported, this library registers and stores the metrics in an application-wide manner.

So essentially we can declare a metric with a struct and interact with it everywhere. This is for developers who does not need more control in registering and managing the metrics' lifecycle, and just want to have a simple interface in instrumenting the app.

## Features

### Declarative definition of metrics

A struct definition can be embedded with one of the following available `MetricType`:
```go
dprom.Counter
dprom.Histogram
dprom.Gauge
```

Embedding of one of these structs will signify that the embedded struct is a prometheus metric.

These rules apply:
- The struct name, transformed into snake case, will be the metric name.
- The fields of that struct will become the labels to that metric. Field name, transformed into snake case, will be the label name.
- The embedded `MetricType` field can be tagged with `help` to provide the help message of the corresponding metric.

### Instrumentation functions according to metric types

According to the embedded `MetricType` inside the declared metric struct, those metrics can be instrumented with the following functions.

Passing an invalid struct will result in a panic.

#### Counter functions

The `IncrementCounter` function can be used to increment `Counter` metrics.

```go
// declare the metric
type ApiHit struct {
    dprom.Counter `help:"Counter that tracks API hits to the server"`
    Path          string
    Method        string
    ResponseCode  int
}

// increment the api_hit metric with some labels
dprom.IncrementCounter(
	ApiHit{
        Path: "/some/path",
        Method: "GET",
        ResponseCode: 200,
    })
```

#### Histogram functions

The `RecordHistogram` and `NewTimer` functions can be used to record `Histogram` metrics.

```go
// declare the metric
type ApiResponseTime struct {
    dprom.Histogram `help:"Histogram that tracks API response times from the server"`
    Path            string
    Method          string
}

// record a float64 value to the api_response_time metric with some labels
dprom.RecordHistogram(
    ApiResponseTime{
        Path: "/some/path", 
        Method: "GET", 
        ResponseCode: 200,
    }, 1.23)

// NewTimer returns a *prometheus.Timer which can be used to time a process
t := dprom.NewTimer(
    ApiResponseTime{
        Path: "/some/other/path",
        Method: "GET",
        ResponseCode: 200,
    })

// do something...

t.ObserveDuration() // record duration of process since invocation NewTimer() until invocation of ObserveDuration()
```

#### Gauge functions

The `SetGauge` function can be used to set the value of `Gauge` metrics.

```go
// declare the metric
type SomeQueueSize struct {
    dprom.Gauge `help:"Gauge that tracks the size of some queue"`
    QueueName   string
}

// set the api_hit gauge metric with some labels to some float64 value
dprom.SetGauge(
    SomeQueueSize{
        QueueName: "some_queue_name",
    }, 1.23)
```

### Setting and getting registerers and gatherers

The `SetGatherer`, `SetRegisterer`, `GetGatherer`, and `GetRegisterer` functions can be used to interact with the prom client's gatherers and registerers. By default, this library will use prometheus' DefaultRegisterer.

```go
// create custom registry
registry := prometheus.NewRegistry()

// set the library to use the new registry for its registerer and gatherer
dprom.SetRegisterer(registry)
dprom.SetGatherer(registry)

//use the GetGatherer and GetRegisterer functions for promhttp.HandlerFor function
http.Handle("/metrics",
    promhttp.HandlerFor(dprom.GetGatherer(),
        promhttp.HandlerOpts{Registry: dprom.GetRegisterer()},
    ))
```

The `SetGatherer` and `SetRegisterer` functions should be used before interacting with any metrics.

## Example application

Example provided in the `/examples` folder of this repository.
