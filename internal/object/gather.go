/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Modifications Copyright 2022 The K8sHorizMetrics Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.

Modified to split up evaluations and metric gathering to work with the
Custom Pod Autoscaler framework.
Original source:
https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/podautoscaler/horizontal.go
https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/podautoscaler/replica_calculator.go
*/

package object

import (
	"fmt"

	"github.com/jthomperoo/k8shorizmetrics/v3/internal/podutil"
	"github.com/jthomperoo/k8shorizmetrics/v3/metrics/object"
	"github.com/jthomperoo/k8shorizmetrics/v3/metrics/value"
	metricsclient "github.com/jthomperoo/k8shorizmetrics/v3/metricsclient"
	autoscaling "k8s.io/api/autoscaling/v2"
	"k8s.io/apimachinery/pkg/labels"
)

// Gather (Object) provides functionality for retrieving metrics for object metric specs.
type Gather struct {
	MetricsClient   metricsclient.Client
	PodReadyCounter podutil.PodReadyCounter
}

// Gather retrieves an object metric
func (c *Gather) Gather(metricName string, namespace string, objectRef *autoscaling.CrossVersionObjectReference, podSelector labels.Selector, metricSelector labels.Selector) (*object.Metric, error) {
	// Get metrics
	utilization, timestamp, err := c.MetricsClient.GetObjectMetric(metricName, namespace, objectRef, metricSelector)
	if err != nil {
		return nil, fmt.Errorf("unable to get metric %s: %s on %s %s: %w", metricName, objectRef.Kind, namespace, objectRef.Name, err)
	}

	// Calculate number of ready pods
	readyPodCount, err := c.PodReadyCounter.GetReadyPodsCount(namespace, podSelector)
	if err != nil {
		return nil, fmt.Errorf("unable to calculate ready pods: %w", err)
	}

	return &object.Metric{
		Current: value.MetricValue{
			Value: &utilization,
		},
		ReadyPodCount: &readyPodCount,
		Timestamp:     timestamp,
	}, nil
}

// GatherPerPod retrieves an object per pod metric
func (c *Gather) GatherPerPod(metricName string, namespace string, objectRef *autoscaling.CrossVersionObjectReference, metricSelector labels.Selector) (*object.Metric, error) {
	// Get metrics
	utilization, timestamp, err := c.MetricsClient.GetObjectMetric(metricName, namespace, objectRef, metricSelector)
	if err != nil {
		return nil, fmt.Errorf("unable to get metric %s: %s on %s %s/%w", metricName, objectRef.Kind, namespace, objectRef.Name, err)
	}

	return &object.Metric{
		Current: value.MetricValue{
			AverageValue: &utilization,
		},
		Timestamp: timestamp,
	}, nil
}
