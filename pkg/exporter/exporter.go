// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
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

package exporter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/falcosecurity/client-go/pkg/api/outputs"
	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	eventsCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "falco_events",
		},
		[]string{
			"rule",
			"priority",
			"hostname",
			"source",
			"k8s_ns_name",
			"k8s_pod_name",
			"tags",
		},
	)
)

func init() {
	prometheus.MustRegister(eventsCounter)
}

// Forward processes a single *outputs.Response and forwards the event to the metric counters.
func Forward(res *outputs.Response) error {
	labels := prometheus.Labels{
		"rule":         res.Rule,
		"priority":     fmt.Sprintf("%d", res.Priority),
		"hostname":     res.Hostname,
		"source":       res.Source,
		"k8s_ns_name":  "",
		"k8s_pod_name": "",
		"tags":         fmt.Sprintf(",%s,", strings.Join(res.Tags, ",")),
	}

	// Ensure OutputFields are enabled
	if res.OutputFields != nil {
		labels["k8s_ns_name"] = res.OutputFields["k8s.ns.name"]
		labels["k8s_pod_name"] = res.OutputFields["k8s.pod.name"]
	}

	eventsCounter.With(labels).Inc()
	return nil
}

// Watch allows to watch and process a stream of *outputs.Response from a given outputs.Service_SubClient.
// The timeout parameter specifies the frequency of the watch operation.
//
// It waits until the stream is closed, the context is cancelled or an error occured.
func Watch(ctx context.Context, sc outputs.Service_SubClient, timeout time.Duration) error {
	return client.OutputsWatch(ctx, sc, Forward, timeout)
}
