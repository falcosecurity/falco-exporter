FROM golang:1.13.4-buster AS builder

ENV GO111MODULE on
ENV CGO_ENABLED 0

WORKDIR /opt/falco-exporter

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make

FROM alpine:3.10
WORKDIR /opt/falco-exporter

RUN addgroup -S falcoexporter && adduser -S falcoexporter -G falcoexporter
USER falcoexporter 

COPY --from=builder /opt/falco-exporter/falco-exporter /usr/bin/falco-exporter

EXPOSE 9376/tcp

ENTRYPOINT ["/usr/bin/falco-exporter"]
