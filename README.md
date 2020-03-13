# falco-exporter
> Prometheus Metrics Exporter for Falco output events

Status: **Under development**

## Prerequisites

Before using **falco-exporter**, you need [Falco installed](https://falco.org/docs/installation/) and running with the [gRPC Output](https://falco.org/docs/grpc/) enabled. The Falco gRPC server works only with mutual TLS by design. Therefore, you also need valid [certificate files](https://falco.org/docs/grpc/#certificates) to configure **falco-exporter** properly.


## Usage

### Run it manually
```
make
./falco-exporter
```
http://localhost:9376/metrics

```
Usage of ./falco-exporter:
      --client-ca string         CA root file path for connecting to a Falco gRPC server (default "/etc/falco/certs/ca.crt")
      --client-cert string       cert file path for connecting to a Falco gRPC server (default "/etc/falco/certs/client.crt")
      --client-hostname string   hostname for connecting to a Falco gRPC server (default "localhost")
      --client-key string        key file path for connecting to a Falco gRPC server (default "/etc/falco/certs/client.key")
      --client-port uint16       port for connecting to a Falco gRPC server (default 5060)
      --listen-address string    address on which to expose the Prometheus metrics (default ":9376")
```

### Deploy in Kubernetes

Using the [provided Helm chart](deploy/helm/falco-exporter/) is the easiest way to deploy **falco-exporter**.

To install the chart with the release name `falco-exporter` and default [configuration values](deploy/helm/falco-exporter/values.yaml):
```
helm install falco-exporter \
      --set-file certs.ca.crt=/path/to/ca.crt,certs.client.key=/path/to/client.key,certs.client.crt=/path/to/client.crt \
      ./deploy/helm/falco-exporter
```

The command deploys **falco-exporter** as Daemon Set on your the Kubernetes cluster. If a [Prometheus installation](https://github.com/helm/charts/tree/master/stable/prometheus) is running within your cluster, metrics provided by **falco-exporter** will be automatically discovered.

Alternatively, it is possible to deploy **falco-exporter** without using Helm. Templates for manual installation are [here](deploy/k8s/falco-exporter).