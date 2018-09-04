export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOBIN?=$(BIN)
export GO=$(shell which go)
export GOPATH=$(abspath $(ROOT)/../../..)
export BUILD=cd $(ROOT) && $(GO) install -v

all:
	$(BUILD) ./cmd/...