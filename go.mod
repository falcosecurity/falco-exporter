module github.com/falcosecurity/falco-exporter

go 1.13

require (
	// todo(leogr): update version once client-go with bidi support has been released
	github.com/falcosecurity/client-go v0.1.1-0.20200609153459-3b6f8eb9e49d
	github.com/prometheus/client_golang v1.1.0
	github.com/spf13/pflag v1.0.5
	google.golang.org/grpc v1.28.0
)
