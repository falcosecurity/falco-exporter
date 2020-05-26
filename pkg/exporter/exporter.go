package exporter

import (
	"context"
	"fmt"

	"github.com/falcosecurity/client-go/pkg/api/output"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

// A RecvFunc waits for subscribed events and forwards to metric counters.
type RecvFunc func() error

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

// Subscribe to a ServiceClient to receive a stream of Falco output events.
// If success, it returns a RecvFunc that can be used.
// To stop the subscription, use Cancel() on the context provided.
func Subscribe(ctx context.Context, outputClient output.ServiceClient, opts ...grpc.CallOption) (RecvFunc, error) {

	// Keepalive true means that the client will wait indefinitely for new events to come
	fcs, err := outputClient.Subscribe(ctx, &output.Request{Keepalive: true}, opts...)
	if err != nil {
		return nil, err
	}

	return func() error {
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

			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
		}
	}, nil
}
