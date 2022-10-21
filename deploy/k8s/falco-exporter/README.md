

# Kubernetes Daemon Sets templates for falco-exporter

The [templates directory](./templates) gives you the required YAML files to stand up **falco-exporter** on Kubernetes as Deamon Set. 
This will result in a **falco-exporter** Pod being deployed to each node.

## Configuration

The default configurations for connecting to a Falco gRPC server over Unix socket are:

- client-socket: `unix:///run/falco/falco.sock`
- timeout: `2m`
- listen-address: `0.0.0.0:9376`

If needed, please modify them in `daemonset.yaml` according to your installation.

## Deploying

```
kubectl create \
    -f ./templates/daemonset.yaml \
    -f ./templates/serviceaccount.yaml \
    -f ./templates/service.yaml
```