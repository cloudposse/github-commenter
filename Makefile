SHELL = /bin/bash

export DOCKER_ORG ?= cloudposse
export DOCKER_IMAGE ?= $(DOCKER_ORG)/github-commenter
export DOCKER_TAG ?= latest
export DOCKER_IMAGE_NAME ?= $(DOCKER_IMAGE):$(DOCKER_TAG)
export DOCKER_BUILD_FLAGS =

PATH:=$(PATH):$(GOPATH)/bin

include $(shell curl --silent -o .build-harness "https://raw.githubusercontent.com/cloudposse/build-harness/master/templates/Makefile.build-harness"; echo .build-harness)


.PHONY : go-get
go-get:
	go get


.PHONY : go-build
go-build: go-get
	CGO_ENABLED=0 go build -v -o "./dist/bin/github-commenter" *.go
