# falco-exporter
> Prometheus Metrics Exporter for Falco output events

[![Release](https://img.shields.io/github/release/falcosecurity/falco-exporter.svg?style=flat-square)](https://github.com/falcosecurity/falco-exporter/releases/latest)
[![License](https://img.shields.io/github/license/falcosecurity/falco-exporter?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/falcosecurity/falco-exporter?style=flat-square)](https://goreportcard.com/report/github.com/falcosecurity/falco-exporter)
[![Docker pulls](https://img.shields.io/docker/pulls/falcosecurity/falco-exporter?style=flat-square)](https://hub.docker.com/r/falcosecurity/falco-exporter)

## Prerequisites

- Before using **falco-exporter**, you need [Falco installed](https://falco.org/docs/getting-started/installation/) and running with the [gRPC Output](https://falco.org/docs/grpc/) enabled (over Unix socket by default).
- Since **falco-exporter** `v0.3.0`: 
  - the minimum required version of Falco is `0.24.0`
  - if using Helm, the minimum required version of the [Falco Chart](https://github.com/falcosecurity/charts/tree/master/falco) is `v1.2.0`
- Since **falco-exporter** `v0.8.0`:
  - the default Unix socket path is `/run/falco/falco.sock` to be compatible with Falco 0.33.0 and later (in previous version it defaulted to `/var/run/falco.sock`)

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
      --client-ca string               CA root file path for connecting to a Falco gRPC server (default "/etc/falco/certs/ca.crt")
      --client-cert string             cert file path for connecting to a Falco gRPC server (default "/etc/falco/certs/client.crt")
      --client-hostname string         hostname for connecting to a Falco gRPC server, if set, takes precedence over --client-socket
      --client-key string              key file path for connecting to a Falco gRPC server (default "/etc/falco/certs/client.key")
      --client-port uint16             port for connecting to a Falco gRPC server (default 5060)
      --client-socket string           unix socket path for connecting to a Falco gRPC server (default "unix:///run/falco/falco.sock")
      --listen-address string          address on which to expose the Prometheus metrics (default ":9376")
      --probes-listen-address string   address on which to expose readiness/liveness probes endpoints (default ":19376")
      --server-ca string               CA root file path for metrics https server
      --server-cert string             cert file path for metrics https server
      --server-key string              key file path for metrics https server
      --timeout duration               timeout for initial gRPC connection (default 2m0s)
```

### Run with Docker

To run **falco-exporter** in a container using Docker:

```shell
docker run -v /path/to/falco.sock:/var/run/falco.sock falcosecurity/falco-exporter
```

### Deploy in Kubernetes

### Using Helm

Using the [falco-exporter Helm Chart](https://github.com/falcosecurity/charts/tree/master/falco-exporter) is the easiest way to deploy **falco-exporter**.

Before installing the chart, add the `falcosecurity` charts repository:

```shell
helm repo add falcosecurity https://falcosecurity.github.io/charts
helm repo update
```

Finally, to install the chart with the release name `falco-exporter` and default [configuration values](https://github.com/falcosecurity/charts/blob/master/falco-exporter/values.yaml):

```shell
helm install falco-exporter falcosecurity/falco-exporter
```

The full documentation of the Helm Chart is [here](https://github.com/falcosecurity/charts/tree/master/falco-exporter).

### Using resource templates

Alternatively, it is possible to deploy **falco-exporter** without using Helm. Templates for manual installation are [here](deploy/k8s/falco-exporter).

### Grafana

The **Falco dashboard** can be imported into Grafana by copy-paste the provided [grafana/dashboard.json](grafana/dashboard.json) or by [getting it from the Grafana Dashboards](https://grafana.com/grafana/dashboards/11914) website.

You can find detailed Grafana importing instructions [here](https://grafana.com/docs/reference/export_import/).

![Falco dashboard](https://github.com/falcosecurity/falco-exporter/raw/master/grafana/preview.png)

## Event priority

Falco events have a priority value, as defined [here](https://github.com/falcosecurity/falco/blob/b76420fe471f8af220d742543637b5aae02ee556/userspace/engine/falco_common.h#L82-L89).
The exported metrics will include a `priority` label that uses a numeric index. The meaning of these indices is reported in the following table.

| ID  | Priority      |
| --- | ------------- |
| 7   | debug         |
| 6   | informational |
| 5   | notice        |
| 4   | warning       |
| 3   | error         |
| 2   | critical      |
| 1   | alert         |
| 0   | emergency     |

## Connection options

**falco-exporter** uses gRPC over a Unix socket by default. 

You may change this behavior by setting `--client-hostname`. Note that the Falco gRPC server over the network works only with mutual TLS by design. Therefore, when `--client-hostname` is set  you also need valid [certificate files](https://falco.org/docs/grpc/#certificates) to configure **falco-exporter** properly (see the *Command line usage* above).
