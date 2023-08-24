SHELL=/bin/bash -o pipefail

GO ?= go
BUILD_ARGS ?= build
HELM ?= helm
CHART_SOURCE_DIR ?= $(shell pwd)/../charts/falco-exporter

TEST_FLAGS ?= -v -race

.PHONY: falco-exporter
falco-exporter:
	$(GO) $(BUILD_ARGS) ./cmd/falco-exporter

.PHONY: deploy/k8s
deploy/k8s:
	rm -rf $@/falco-exporter/templates/*
	$(HELM) template falco-exporter ${CHART_SOURCE_DIR} \
		--set skipHelm=true \
		--output-dir $@

.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test ${TEST_FLAGS} ./...
