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

// Package resource contains models for resource metrics as returned by the K8s metrics APIs.
package resource

import (
	"time"

	"github.com/jthomperoo/k8shorizmetrics/v4/metrics/podmetrics"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Metric (Resource) is a resource metric known to Kubernetes, as specified in requests and limits, describing each pod
// in the current scale target (e.g. CPU or memory).  Such metrics are built in to Kubernetes, and have special scaling
// options on top of those available to normal per-pod metrics (the "pods" source).
type Metric struct {
	PodMetricsInfo podmetrics.MetricsInfo `json:"podMetricsInfo"`
	Requests       map[string]int64       `json:"requests"`
	ReadyPodCount  int64                  `json:"readyPodCount"`
	IgnoredPods    sets.String            `json:"ignoredPods"`
	MissingPods    sets.String            `json:"missingPods"`
	TotalPods      int                    `json:"totalPods"`
	Timestamp      time.Time              `json:"timestamp,omitempty"`
}
