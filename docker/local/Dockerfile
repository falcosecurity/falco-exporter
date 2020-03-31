FROM golang:1.14-buster AS builder

ENV CGO_ENABLED 0

ENV CLIENT_CA "/etc/falco/certs/ca.crt"
ENV CLIENT_CERT "/etc/falco/certs/client.crt"
ENV CLIENT_HOSTNAME "localhost"
ENV CLIENT_KEY "/etc/falco/certs/client.key"
ENV CLIENT_PORT "5060"
ENV LISTEN_ADDRESS ":9376"

WORKDIR /opt/falco-exporter

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make

FROM alpine:3.11
WORKDIR /opt/falco-exporter

RUN addgroup -S falcoexporter && adduser -S falcoexporter -G falcoexporter
USER falcoexporter 

COPY --from=builder /opt/falco-exporter/falco-exporter .

EXPOSE 9376/tcp

CMD ["sh", "-c", "./falco-exporter --client-ca ${CLIENT_CA} --client-cert ${CLIENT_CERT} --client-hostname ${CLIENT_HOSTNAME} --client-key ${CLIENT_KEY} --client-port ${CLIENT_PORT} --listen-address ${LISTEN_ADDRESS}"]
