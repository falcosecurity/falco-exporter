# falco-exporter
> Prometheus Metrics Exporter for Falco output events

[![Release](https://img.shields.io/github/release/falcosecurity/falco-exporter.svg?style=flat-square)](https://github.com/falcosecurity/falco-exporter/releases/latest)
[![License](https://img.shields.io/github/license/falcosecurity/falco-exporter?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/falcosecurity/falco-exporter?style=flat-square)](https://goreportcard.com/report/github.com/falcosecurity/falco-exporter)
[![Docker pulls](https://img.shields.io/docker/pulls/falcosecurity/falco-exporter?style=flat-square)](https://hub.docker.com/r/falcosecurity/falco-exporter)

## Prerequisites

- Before using **falco-exporter**, you need [Falco installed](https://falco.org/docs/installation/) and running with the [gRPC Output](https://falco.org/docs/grpc/) enabled (over Unix socket by default).
- Since `v0.2.0` the minimum required version of Falco is `0.24.0`.


## Usage

### Run it manually

```shell
make
./falco-exporter
```
Then check the metrics endpoint at http://localhost:9376/metrics

Command line usage:
```
$ ./falco-exporter --help
Usage of ./falco-exporter:
      --client-ca string         CA root file path for connecting to a Falco gRPC server (default "/etc/falco/certs/ca.crt")
      --client-cert string       cert file path for connecting to a Falco gRPC server (default "/etc/falco/certs/client.crt")
      --client-hostname string   hostname for connecting to a Falco gRPC server, if set, takes precedence over --client-socket
      --client-key string        key file path for connecting to a Falco gRPC server (default "/etc/falco/certs/client.key")
      --client-port uint16       port for connecting to a Falco gRPC server (default 5060)
      --client-socket string     unix socket path for connecting to a Falco gRPC server (default "unix:///var/run/falco.sock")
      --listen-address string    address on which to expose the Prometheus metrics (default ":9376")
      --timeout duration         timeout for initial gRPC connection (default 2m0s)
```

### Deploy in Kubernetes

Using the [provided Helm chart](deploy/helm/falco-exporter/) is the easiest way to deploy **falco-exporter**.

To install the chart with the release name `falco-exporter` and default [configuration values](deploy/helm/falco-exporter/values.yaml):
```shell
helm install falco-exporter ./deploy/helm/falco-exporter
```

The command deploys **falco-exporter** as Daemon Set on your the Kubernetes cluster. If a [Prometheus installation](https://github.com/helm/charts/tree/master/stable/prometheus) is running within your cluster, metrics provided by **falco-exporter** will be automatically discovered.

Alternatively, it is possible to deploy **falco-exporter** without using Helm. Templates for manual installation are [here](deploy/k8s/falco-exporter).

### Grafana

The **Falco dashboard** can be imported into Grafana by copy-paste the provided [grafana/dashboard.json](grafana/dashboard.json) or by [getting it from the Grafana Dashboards](https://grafana.com/grafana/dashboards/11914) website.

You can find detailed Grafana importing instructions [here](https://grafana.com/docs/reference/export_import/).

![Falco dashboard](https://github.com/falcosecurity/falco-exporter/raw/master/grafana/preview.png)

## Connection options

**falco-exporter** uses gRPC over Unix socket by default. 

You may change this behavior by setting `--client-hostname`. Note that the Falco gRPC server over the network works only with mutual TLS by design. Therefore, when `--client-hostname` is set  you also need valid [certificate files](https://falco.org/docs/grpc/#certificates) to configure **falco-exporter** properly (see the *Command line usage* above).