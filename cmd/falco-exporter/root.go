package main

import (
	"context"
	goflag "flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/falcosecurity/client-go/pkg/client"
	"github.com/falcosecurity/falco-exporter/pkg/exporter"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/ssgreg/repeat"
)

func main() {

	addr := ""
	pflag.StringVar(&addr, "listen-address", ":9376", "address on which to expose the Prometheus metrics")

	config := &client.Config{}
	pflag.StringVar(&config.Hostname, "client-hostname", "localhost", "hostname for connecting to a Falco gRPC server")
	pflag.Uint16Var(&config.Port, "client-port", 5060, "port for connecting to a Falco gRPC server")
	pflag.StringVar(&config.CertFile, "client-cert", "/etc/falco/certs/client.crt", "cert file path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.KeyFile, "client-key", "/etc/falco/certs/client.key", "key file path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.CARootFile, "client-ca", "/etc/falco/certs/ca.crt", "CA root file path for connecting to a Falco gRPC server")

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

	g := run.Group{}
	g.Add(
		func() error {
			op := func(c int) error {
				err := exporter.Subscribe(ctx, outputClient)
				if err != nil {
					return repeat.HintTemporary(
						fmt.Errorf("server subscription (attempt %d) unavailable: %w", c, err),
					)
				}
				return nil
			}
			return repeat.Repeat(
				repeat.FnWithCounter(op),
				repeat.LimitMaxTries(30),
				repeat.WithDelay(
					repeat.FullJitterBackoff(10*time.Second).Set(),
					repeat.SetContext(ctx),
				),
			)
		},
		func(err error) {
			log.Println("exiting due to repeated error: ", err)
		},
	)

	http.Handle("/metrics", promhttp.Handler())
	ln, _ := net.Listen("tcp", addr)
	g.Add(
		func() error {
			return http.Serve(ln, nil)
		},
		func(error) {
			ln.Close()
		},
	)

	// Exit with first error from the run group.
	log.Fatalln(g.Run())
}
