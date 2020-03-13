SHELL=/bin/bash -o pipefail

DOCKER ?= docker
GO ?= go
HELM ?= helm

COMMIT_NO := $(shell git rev-parse HEAD 2> /dev/null || true)
GIT_COMMIT := $(if $(shell git status --porcelain --untracked-files=no),${COMMIT_NO}-dirty,${COMMIT_NO})
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
GIT_BRANCH_CLEAN := $(shell echo $(GIT_BRANCH) | sed -e "s/[^[:alnum:]]/-/g")

IMAGE_NAME_BUILDER_BASE ?= docker.io/falcosecurity/falco-exporter

IMAGE_NAME_BUILDER_BASE_BRANCH := $(IMAGE_NAME_BUILDER_BASE):$(GIT_BRANCH_CLEAN)
IMAGE_NAME_BUILDER_BASE_COMMIT := $(IMAGE_NAME_BUILDER_BASE):$(GIT_COMMIT)
IMAGE_NAME_BUILDER_BASE_LATEST := $(IMAGE_NAME_BUILDER_BASE):latest

TEST_FLAGS ?= -v -race

export GO111MODULE=on

.PHONY: falco-exporter
falco-exporter:
	$(GO) build ./cmd/falco-exporter

.PHONY: image/build
image/build:
	$(DOCKER) build \
		-t "$(IMAGE_NAME_BUILDER_BASE_BRANCH)" \
		-f build/Dockerfile .
	$(DOCKER) tag $(IMAGE_NAME_BUILDER_BASE_BRANCH) $(IMAGE_NAME_BUILDER_BASE_COMMIT)
	$(DOCKER) tag "$(IMAGE_NAME_BUILDER_BASE_BRANCH)" $(IMAGE_NAME_BUILDER_BASE_COMMIT)


.PHONY: image/push
image/push:
	$(DOCKER) push $(IMAGE_NAME_BUILDER_BASE_BRANCH)
	$(DOCKER) push $(IMAGE_NAME_BUILDER_BASE_COMMIT)

.PHONY: image/latest
image/latest:
	$(DOCKER) tag $(IMAGE_NAME_BUILDER_BASE_COMMIT) $(IMAGE_NAME_BUILDER_BASE_LATEST)
	$(DOCKER) push $(IMAGE_NAME_BUILDER_BASE_LATEST)

.PHONY: deploy/k8s
deploy/k8s:
	rm -rf deploy/k8s/*
	$(HELM) template falco-exporter deploy/helm/falco-exporter \
		--set skipHelm=true \
		--output-dir deploy/k8s

.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test ${TEST_FLAGS} ./...