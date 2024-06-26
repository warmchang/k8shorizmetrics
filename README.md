[![Build](https://github.com/jthomperoo/k8shorizmetrics/workflows/main/badge.svg)](https://github.com/jthomperoo/k8shorizmetrics/actions)
[![go.dev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/jthomperoo/k8shorizmetrics/v4)
[![Go Report
Card](https://goreportcard.com/badge/github.com/jthomperoo/k8shorizmetrics/v4)](https://goreportcard.com/report/github.com/jthomperoo/k8shorizmetrics/v4)
[![License](https://img.shields.io/:license-apache-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)

# k8shorizmetrics

`k8shorizmetrics` is a library that provides the internal workings of the Kubernetes Horizontal Pod Autoscaler (HPA)
wrapped up in a simple API. The project allows querying metrics just as the HPA does, and also running the calculations
to work out the target replica count that the HPA does.

## Install

```bash
go get -u github.com/jthomperoo/k8shorizmetrics/v4@v4.0.0
```

## Features

- Simple API, based directly on the code from the HPA, but detangled for ease of use.
- Dependent only on versioned and public Kubernetes Golang modules, allows easy install without replace directives.
- Splits the HPA into two parts, metric gathering and evaluation, only use what you need.
- Allows insights into how the HPA makes decisions.
- Supports scaling to and from 0.

## Quick Start

The following is a simple program that can run inside a Kubernetes cluster that gets the CPU resource metrics for
pods with the label `run: php-apache`.

```go
package main

import (
	"log"
	"time"

	"github.com/jthomperoo/k8shorizmetrics/v4"
	"github.com/jthomperoo/k8shorizmetrics/v4/metricsclient"
	"github.com/jthomperoo/k8shorizmetrics/v4/podsclient"
	"k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Kubernetes API setup
	clusterConfig, _ := rest.InClusterConfig()
	clientset, _ := kubernetes.NewForConfig(clusterConfig)
	// Metrics and pods clients setup
	metricsclient := metricsclient.NewClient(clusterConfig, clientset.Discovery())
	podsclient := &podsclient.OnDemandPodLister{Clientset: clientset}
	// HPA configuration options
	cpuInitializationPeriod := time.Duration(300) * time.Second
	initialReadinessDelay := time.Duration(30) * time.Second

	// Setup gatherer
	gather := k8shorizmetrics.NewGatherer(metricsclient, podsclient, cpuInitializationPeriod, initialReadinessDelay)

	// Target resource values
	namespace := "default"
	podSelector := labels.SelectorFromSet(labels.Set{
		"run": "php-apache",
	})

	// Metric spec to gather, CPU resource utilization
	spec := v2.MetricSpec{
		Type: v2.ResourceMetricSourceType,
		Resource: &v2.ResourceMetricSource{
			Name: corev1.ResourceCPU,
			Target: v2.MetricTarget{
				Type: v2.UtilizationMetricType,
			},
		},
	}

	metric, _ := gather.GatherSingleMetric(spec, namespace, podSelector)

	for pod, podmetric := range metric.Resource.PodMetricsInfo {
		actualCPU := podmetric.Value
		requestedCPU := metric.Resource.Requests[pod]
		log.Printf("Pod: %s, CPU usage: %dm (%0.2f%% of requested)\n", pod, actualCPU, float64(actualCPU)/float64(requestedCPU)*100.0)
	}
}
```

## Documentation

See the [Go doc](https://pkg.go.dev/github.com/jthomperoo/k8shorizmetrics/v4).

## Migration

This section explains how to migrate between versions of the library.

### From v1 to v2

There are two changes you need to make to migrate from `v1` to `v2`:

1. Switch from using `k8s.io/api/autoscaling/v2beta2` to `k8s.io/api/autoscaling/v2`.
2. Switch from using `github.com/jthomperoo/k8shorizmetrics` to `github.com/jthomperoo/k8shorizmetrics/v2`.

### From v2 to v3

The breaking changes introduced by `v3` are:

- Gather now returns the `GathererMultiMetricError` error type if any of the metrics fail to gather. This error is
returned for partial errors, meaning some metrics gathered successfully and others did not. If this partial error
occurs the `GathererMultiMetricError` error will have the `Partial` property set to `true`. This can be checked for
using `errors.As`.
- Evaluate now returns the `EvaluatorMultiMetricError` error type if any of the metrics fail to
evaluate. This error is returned for partial errors, meaning some metrics evaluted successfully and others did not.
If this partial error occurs the `EvaluatorMultiMetricError` error will have the `Partial` property set to `true`. This
can be checked for using `errors.As`.

To update to `v3` you will need to update all references in your code that refer to
`github.com/jthomperoo/k8shorizmetrics/v2` to use `github.com/jthomperoo/k8shorizmetrics/v3`.

If you want the behaviour to stay the same and to swallow partial errors you can use code like this:

```go
metrics, err := gather.Gather(specs, namespace, podMatchSelector)
if err != nil {
	gatherErr := &k8shorizmetrics.GathererMultiMetricError{}
	if !errors.As(err, &gatherErr) {
		log.Fatal(err)
	}

	if !gatherErr.Partial {
		log.Fatal(err)
	}

	// Not a partial error, just continue as normal
}
```

You can use similar code for the `Evaluate` method of the `Evaluater`.

### From v3 to v4

To update to `v4` you will need to update all references in your code that refer to
`github.com/jthomperoo/k8shorizmetrics/v3` to use `github.com/jthomperoo/k8shorizmetrics/v4`.

The only behaviour change is around serialisation into JSON. Fields are now serialised using camel case rather
than snake case to match the Kubernetes conventions.

If you are relying on JSON serialised values you need to use camel case now. For example the Resource Metric field
`PodMetricsInfo` is now serialised as `podMetricsInfo` ratherthan `pod_metrics_info`.

## Examples

See the [examples directory](./examples/) for some examples, [cpuprint](./examples/cpuprint/) is a good start.

## Developing and Contributing

See the [contribution guidelines](CONTRIBUTING.md) and [code of conduct](CODE_OF_CONDUCT.md).
