package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	goflag "flag"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

	var probesAddr string
	pflag.StringVar(&probesAddr, "probes-listen-address", ":19376", "address on which to expose readiness/liveness probes endpoints")

	var serverCARootFile string
	pflag.StringVar(&serverCARootFile, "server-ca", "", "CA root file path for metrics https server")

	var serverCertFile string
	pflag.StringVar(&serverCertFile, "server-cert", "", "cert file path for metrics https server")

	var serverKeyFile string
	pflag.StringVar(&serverKeyFile, "server-key", "", "key file path for metrics https server")

	var timeout time.Duration
	pflag.DurationVar(&timeout, "timeout", time.Minute*2, "timeout for initial gRPC connection")

	config := &client.Config{
		DialOptions: []grpc.DialOption{
			// Instruct `client.NewForConfig` to wait until the underlying connection is up,
			// the dialer will use the default gRPC backoff if needed.
			grpc.WithBlock(),
		},
	}
	pflag.StringVar(&config.UnixSocketPath, "client-socket", "unix:///run/falco/falco.sock", "unix socket path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.Hostname, "client-hostname", "", "hostname for connecting to a Falco gRPC server, if set, takes precedence over --client-socket")
	pflag.Uint16Var(&config.Port, "client-port", 5060, "port for connecting to a Falco gRPC server")
	pflag.StringVar(&config.CertFile, "client-cert", "/etc/falco/certs/client.crt", "cert file path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.KeyFile, "client-key", "/etc/falco/certs/client.key", "key file path for connecting to a Falco gRPC server")
	pflag.StringVar(&config.CARootFile, "client-ca", "/etc/falco/certs/ca.crt", "CA root file path for connecting to a Falco gRPC server")
	pflag.BoolVar(&config.GRPCAuth, "grpc-auth", true, "Whether or not falco-exporter authenticates itself to the gRPC server")


	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	pflag.Parse()

	go serveMetrics(addr, serverCARootFile, serverCertFile, serverKeyFile)
	probeMux := enableProbes(probesAddr)

	if config.Hostname != "" {
		config.UnixSocketPath = ""
		log.Printf("connecting to gRPC server at %s:%d (timeout %s)", config.Hostname, config.Port, timeout)
	} else {
		if !strings.HasPrefix(config.UnixSocketPath, "unix://") {
			config.UnixSocketPath = "unix://" + config.UnixSocketPath
		}
		log.Printf("connecting to gRPC server at %s (timeout %s)", config.UnixSocketPath, timeout)
	}

	// main context
	ctx := withSignals(context.Background())

	// cancel the pending connection after timeout is reached
	dialerCtx, cancelTimeout := context.WithTimeout(ctx, timeout)
	c, err := client.NewForConfig(dialerCtx, config)
	if err != nil {
		log.Fatalf("gRPC: %v\n", err)
	}
	defer c.Close()

	log.Println("connected to gRPC server, subscribing events stream")

	oc, err := c.Outputs()
	if err != nil {
		log.Fatalf("gRPC: %v\n", err)
	}

	fsc, err := oc.Sub(ctx)
	if err != nil {
		log.Fatalf("gRPC: %v\n", err)
	}

	cancelTimeout()
	enableReadiness(probeMux)
	log.Println("ready")

	if err := exporter.Watch(ctx, fsc, time.Second); err != nil {
		log.Fatalf("gRPC: %v\n", err)
	} else {
		log.Println("gRPC stream closed")
	}
}

func serveMetrics(addr string, caFile string, cert string, key string) {
	// Configure mTLS if applies
	var mTLS = false
	var tlsConfig *tls.Config = nil
	if caFile != "" {
		// Load CA cert
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Fatalf("mTLS: %v\n", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		// Create the TLS Config with the CA pool and enable Client certificate validation
		tlsConfig = &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		log.Println("TLS configured successfully")
		mTLS = true
	}

	http.Handle("/metrics", promhttp.Handler())
	if mTLS {
		server := &http.Server{
			Addr:      addr,
			TLSConfig: tlsConfig,
		}
		log.Printf("listening on https://%s/metrics\n", addr)
		if err := server.ListenAndServeTLS(cert, key); err != nil {
			log.Fatalf("TLS server: %v", err)
		}
	} else {
		log.Printf("listening on http://%s/metrics\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("server: %v", err)
		}
	}
}

func enableProbes(probesAddr string) (probeMux *http.ServeMux) {
	// probes are served in a different ServeMux since a possible mTLS config may be used on main server (/metrics)
	probeMux = http.NewServeMux()
	probeMux.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	go func() {
		if err := http.ListenAndServe(probesAddr, probeMux); err != nil {
			log.Fatalf("healthz server: %v", err)
		}
	}()
	return
}

func enableReadiness(probeMux *http.ServeMux) {
	probeMux.HandleFunc("/readiness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
}
