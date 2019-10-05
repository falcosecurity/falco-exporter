SHELL=/bin/bash -o pipefail

GO ?= go

TEST_FLAGS ?= -v -race

export GO111MODULE=on

.PHONY: falco-exporter
falco-exporter:
	$(GO) build ./cmd/falco-exporter

.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test ${TEST_FLAGS} ./...