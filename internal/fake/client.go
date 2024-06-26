/*
Copyright 2022 The K8sHorizMetrics Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	"time"

	"github.com/jthomperoo/k8shorizmetrics/v4/metrics/podmetrics"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// MetricsClient (fake) provides a way to insert functionality into a metricsclient
type MetricsClient struct {
	GetResourceMetricReactor func(resource corev1.ResourceName, namespace string, selector labels.Selector) (podmetrics.MetricsInfo, time.Time, error)
	GetRawMetricReactor      func(metricName string, namespace string, selector labels.Selector, metricSelector labels.Selector) (podmetrics.MetricsInfo, time.Time, error)
	GetObjectMetricReactor   func(metricName string, namespace string, objectRef *autoscalingv2.CrossVersionObjectReference, metricSelector labels.Selector) (int64, time.Time, error)
	GetExternalMetricReactor func(metricName string, namespace string, selector labels.Selector) ([]int64, time.Time, error)
}

// GetResourceMetric calls the fake metricsclient function
func (f *MetricsClient) GetResourceMetric(resource corev1.ResourceName, namespace string, selector labels.Selector) (podmetrics.MetricsInfo, time.Time, error) {
	return f.GetResourceMetricReactor(resource, namespace, selector)
}

// GetRawMetric calls the fake metricsclient function
func (f *MetricsClient) GetRawMetric(metricName string, namespace string, selector labels.Selector, metricSelector labels.Selector) (podmetrics.MetricsInfo, time.Time, error) {
	return f.GetRawMetricReactor(metricName, namespace, selector, metricSelector)
}

// GetObjectMetric calls the fake metricsclient function
func (f *MetricsClient) GetObjectMetric(metricName string, namespace string, objectRef *autoscalingv2.CrossVersionObjectReference, metricSelector labels.Selector) (int64, time.Time, error) {
	return f.GetObjectMetricReactor(metricName, namespace, objectRef, metricSelector)
}

// GetExternalMetric calls the fake metricsclient function
func (f *MetricsClient) GetExternalMetric(metricName string, namespace string, selector labels.Selector) ([]int64, time.Time, error) {
	return f.GetExternalMetricReactor(metricName, namespace, selector)
}
