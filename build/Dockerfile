FROM golang:1.14-buster AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 0

WORKDIR /opt/falco-exporter

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make

FROM alpine:3.11
WORKDIR /opt/falco-exporter

COPY --from=builder /opt/falco-exporter/falco-exporter /usr/bin/falco-exporter

EXPOSE 9376/tcp

ENTRYPOINT ["/usr/bin/falco-exporter"]
