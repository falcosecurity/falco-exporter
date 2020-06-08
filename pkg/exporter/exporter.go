package exporter

import (
	"context"
	"fmt"

	"github.com/falcosecurity/client-go/pkg/api/output"
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
		},
	)
)

func init() {
	prometheus.MustRegister(eventsCounter)
}

// Subscribe to a ServiceClient to receive a stream of falco output events.
func Subscribe(ctx context.Context, outputClient output.ServiceClient) error {
	// Keepalive true means that the client will wait indefinitely for new events to come
	// Use keepalive false if you only want to receive the accumulated events and stop
	fcs, err := outputClient.Subscribe(ctx, &output.Request{Keepalive: true})
	if err != nil {
		return err
	}

	for {
		res, err := fcs.Recv()
		if err != nil {
			return err
		}

		labels := prometheus.Labels{
			"rule":         res.Rule,
			"priority":     fmt.Sprintf("%d", res.Priority),
			"hostname":     res.Hostname,
			"source":       res.Source.String(),
			"k8s_ns_name":  "",
			"k8s_pod_name": "",
		}

		// Ensure OutputFields are enabled
		if res.OutputFields != nil {
			labels["k8s_ns_name"] = res.OutputFields["k8s.ns.name"]
			labels["k8s_pod_name"] = res.OutputFields["k8s.pod.name"]
		}

		eventsCounter.With(labels).Inc()
	}
}
