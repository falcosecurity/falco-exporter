FROM golang:1.13.4-buster AS builder

ENV CLIENT_CA "/tmp/ca.crt"
ENV CLIENT_CERT "/tmp/client.crt"
ENV CLIENT_HOSTNAME "localhost"
ENV CLIENT_KEY "/tmp/client.key"
ENV CLIENT_PORT "5060"
ENV LISTEN_ADDRESS ":9376"

WORKDIR /opt/falco-exporter

EXPOSE 9376/tcp

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make

FROM debian:bullseye-slim
WORKDIR /opt/falco-exporter

COPY --from=builder /opt/falco-exporter/falco-exporter .

CMD ["sh", "-c", "./falco-exporter --client-ca ${CLIENT_CA} --client-cert ${CLIENT_CERT} --client-hostname ${CLIENT_HOSTNAME} --client-key ${CLIENT_KEY} --client-port ${CLIENT_PORT} --listen-address ${LISTEN_ADDRESS}"]
