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
			"event_pod_name",
			"event_namespace",
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

	var namespace string
	var podName string

	for {
		res, err := fcs.Recv()
		if err != nil {
			return err
		}

		namespace = ""
		podName = ""

		//Ensure OutputFields are enabled
		if res.OutputFields != nil{
			ns,ok := res.OutputFields["k8s.ns.name"]
			if ok{
				namespace = ns
			}
			pn,ok := res.OutputFields["k8s.pod.name"]
			if ok{
				podName = pn
			}
		}

		eventsCounter.With(prometheus.Labels{
			"rule":     res.Rule,
			"priority": fmt.Sprintf("%d", res.Priority),
			"hostname": res.Hostname,
			"event_pod_name": podName,
			"event_namespace": namespace,
		}).Inc()
	}
}
