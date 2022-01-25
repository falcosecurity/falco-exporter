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
