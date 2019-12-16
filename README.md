# falco-exporter
> Prometheus Metrics Exporter for Falco output events

*Work in progress!!!*

## Usage

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