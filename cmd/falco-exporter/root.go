package main

import (
	"context"
	goflag "flag"
	"log"
	"net/http"
	"time"

	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/falcosecurity/falco-exporter/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

func main() {

	var addr string
	pflag.StringVar(&addr, "listen-address", ":9376", "address on which to expose the Prometheus metrics")

	var timeout time.Duration
	pflag.DurationVar(&timeout, "timeout", time.Minute*2, "timeout for initial gRPC connection")

	config := &client.Config{
		DialOptions: []grpc.DialOption{
			// Instruct `client.NewForConfig` to wait until the underlying connection is up,
			// the dialer will use the default gRPC backoff if needed.
			grpc.WithBlock(),
		},
	}
	pflag.StringVar(&config.Hostname, "client-hostname", "localhost", "hostname for connecting to a Falco gRPC server")
	pflag.Uint16Var(&config.Port, "client-port", 5060, "port for connecting to a Falco gRPC server")
	pflag.StringVar(&config.CertFile, "client-cert", "/etc/falco/certs/client.crt", "cert file path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.KeyFile, "client-key", "/etc/falco/certs/client.key", "key file path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.CARootFile, "client-ca", "/etc/falco/certs/ca.crt", "CA root file path for connecting to a Falco gRPC server")

	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()

	go serveMetrics(addr)

	log.Printf("connecting to gRPC server %s:%d", config.Hostname, config.Port)
	// cancel the pending connection after timeout is reached
	dialerCtx, cancelTimeout := context.WithTimeout(context.Background(), timeout)
	c, err := client.NewForConfig(dialerCtx, config)
	if err != nil {
		log.Fatalf("gRPC: %v\n", err)
	}
	defer c.Close()
	outputClient, err := c.Output()
	if err != nil {
		log.Fatalf("gRPC: %v\n", err)
	}

	log.Println("connected to gRPC server, subscribing events stream")

	recv, err := exporter.Subscribe(context.Background(), outputClient)
	if err != nil {
		log.Fatalf("gRPC: %v\n", err)
	}

	cancelTimeout()
	enableReadiness()
	log.Println("ready, receiving events")
	if err := recv(); err != nil {
		log.Fatalf("gRPC: %v\n", err)
	}
}

func serveMetrics(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("listening on %s/metrics\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server: %v", err)
	}
}

func enableReadiness() {
	http.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
