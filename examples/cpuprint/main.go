/*
Copyright 2024 The K8sHorizMetrics Authors.

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

package main

import (
	"log"
	"path/filepath"
	"time"

	"github.com/jthomperoo/k8shorizmetrics/v4"
	"github.com/jthomperoo/k8shorizmetrics/v4/metricsclient"
	"github.com/jthomperoo/k8shorizmetrics/v4/podsclient"
	v2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	cpuInitializationPeriodSeconds = 300
	initialReadinessDelaySeconds   = 30
	namespace                      = "default"
)

var podMatchSelector = labels.SelectorFromSet(labels.Set{
	"run": "php-apache",
})

func main() {
	clusterConfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		log.Fatalf("Fail to create out-of-cluster Kubernetes config: %s", err)
	}

	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		log.Fatalf("Fail to set up Kubernetes clientset: %s", err)
	}

	metricsclient := metricsclient.NewClient(clusterConfig, clientset.Discovery())
	podsclient := &podsclient.OnDemandPodLister{
		Clientset: clientset,
	}
	cpuInitializationPeriod := time.Duration(cpuInitializationPeriodSeconds) * time.Second
	initialReadinessDelay := time.Duration(initialReadinessDelaySeconds) * time.Second

	// Set up the metric gatherer, needs to be able to query metrics and pods with the clients provided, along with
	// config options
	gather := k8shorizmetrics.NewGatherer(metricsclient, podsclient, cpuInitializationPeriod, initialReadinessDelay)

	// This is the metric spec, this targets the CPU resource metric, gathering utilization values
	// Equivalent to the following YAML:
	// metrics:
	// - type: Resource
	//   resource:
	// 	   name: cpu
	// 	   target:
	// 	     type: Utilization
	spec := v2.MetricSpec{
		Type: v2.ResourceMetricSourceType,
		Resource: &v2.ResourceMetricSource{
			Name: corev1.ResourceCPU,
			Target: v2.MetricTarget{
				Type: v2.UtilizationMetricType,
			},
		},
	}

	// Loop infinitely, wait 5 seconds between each loop
	for {
		time.Sleep(5 * time.Second)

		// Gather the metrics using the specs, targeting the namespace and pod selector defined above
		metric, err := gather.GatherSingleMetric(spec, namespace, podMatchSelector)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println("CPU metrics:")

		for pod, podmetric := range metric.Resource.PodMetricsInfo {
			actualCPU := podmetric.Value
			requestedCPU := metric.Resource.Requests[pod]
			log.Printf("Pod: %s, CPU usage: %dm (%0.2f%% of requested)\n", pod, actualCPU, float64(actualCPU)/float64(requestedCPU)*100.0)
		}

		log.Println("----------")
	}
}
