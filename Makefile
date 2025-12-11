VERSION ?= $(shell git describe --tags --always --dirty)
EMBED_VERSION ?= $(shell go list -m -f '{{.Version}}' github.com/workpi-ai/model-registry 2>/dev/null || echo "v0.0.0")

.PHONY: build test clean install

build:
	@echo "Building with embed version: $(EMBED_VERSION)"
	go build -ldflags "\
		-X github.com/workpi-ai/model-registry-go.EmbedVersion=$(EMBED_VERSION)" \
		./...

test:
	go test -v ./...

clean:
	go clean
	rm -rf bin/

install:
	go mod download
	go mod tidy
