

# Kubernetes Daemon Sets templates for falco-exporter

The [templates directory](./templates) gives you the required YAML files to stand up **falco-exporter** on Kubernetes as Deamon Set. 
This will result in a **falco-exporter** Pod being deployed to each node.

## Configuration


### Certificates
The Daemon Set relies on a Kubernetes Secret to store the [certificates](https://falco.org/docs/grpc/#certificates). 

Before deploying **falco-exporter** you have to add these certificates to `secret-certs.yaml`:

```
echo "  ca.crt: `cat /path/to/ca.crt | base64 -w0`" >> secret-certs.yaml
echo "  client.crt: `cat /path/to/client.crt | base64 -w0`" >> secret-certs.yaml
echo "  client.key: `cat /path/to/client.key | base64 -w0`" >> secret-certs.yaml
```


### Falco gRCP server

The default configurations for connecting to a Falco gRPC server are:

- hostname: `falco-grpc.default.svc.cluster.local`
- port: `5060`

If needed, please modify them in `daemonset.yaml` according to your installation.

## Deploying

```
kubectl create \
    -f daemonset.yaml \
    -f secret-certs.yaml \
    -f serviceaccount.yaml \
    -f service.yaml
```