package main

import (
	"context"
	goflag "flag"
	"log"
	"net/http"

	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/falcosecurity/falco-exporter/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
)

func main() {

	addr := ""
	pflag.StringVar(&addr, "listen-address", ":9376", "address on which to expose the Prometheus metrics")

	config := &client.Config{}
	pflag.StringVar(&config.Hostname, "client-hostname", "localhost", "hostname for connecting to a Falco gRPC server")
	pflag.Uint16Var(&config.Port, "client-port", 5060, "port for connecting to a Falco gRPC server")
	pflag.StringVar(&config.CertFile, "client-cert", "/tmp/client.crt", "cert file path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.KeyFile, "client-key", "/tmp/client.key", "key file path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.CARootFile, "client-ca", "/tmp/ca.crt", "CA root file path for connecting to a Falco gRPC server")

	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()

	c, err := client.NewForConfig(config)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	defer c.Close()
	outputClient, err := c.Output()
	if err != nil {
		log.Fatalf("unable to obtain an output client: %v", err)
	}

	ctx := context.Background()
	go func() {
		if err := exporter.Subscribe(ctx, outputClient); err != nil {
			if err != nil {
				log.Fatalf("exporter error: %v", err)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	if err = http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("%v", err)
	}
}
