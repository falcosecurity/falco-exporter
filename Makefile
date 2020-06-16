SHELL=/bin/bash -o pipefail

GO ?= go
HELM ?= helm

TEST_FLAGS ?= -v -race

.PHONY: falco-exporter
falco-exporter:
	$(GO) build ./cmd/falco-exporter

.PHONY: deploy/k8s
deploy/k8s:
	rm -rf $@/falco-exporter/templates/*
	$(HELM) template falco-exporter deploy/helm/falco-exporter \
		--set skipHelm=true \
		--output-dir $@

.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test ${TEST_FLAGS} ./...