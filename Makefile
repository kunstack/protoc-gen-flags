# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools.
CONTAINER_TOOL ?= docker
GOBIN ?= $(shell go env GOBIN)

# Version information for build-time injection
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse --short HEAD)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build-time ldflags for version injection
LDFLAGS = -X musecut.com/musecut/pkg/version.gitVersion=$(VERSION)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(GOBIN))
	GOBIN = $(shell go env GOPATH)/bin
endif

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

## Tool Binaries
BUF_GEN ?= $(GOBIN)/buf
CI_LINT ?= $(GOBIN)/golangci-lint

.PHONY: fmt
fmt: ## Run go fmt to gofmt (reformat) package sources.
	go fmt ./...
	$(BUF_GEN) format -w

.PHONY: vet
vet: ## Run go vet to report likely mistakes in packages.
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint to check code quality.
	$(GOBIN)/golangci-lint run ./...

.PHONY: tidy
tidy: ## Run go mod tidy to clean up go.mod and go.sum.
	go mod tidy

.PHONY: test
test: ## Run tests.
	go test -v ./...

.PHONY: clean
clean: ## Clean build artifacts.
	go clean -testcache

.PHONY: generate
generate: deps ## Run generate command to generate Go files and docs by processing source.
	$(BUF_GEN) generate
	go generate ./...

##@ Publishing

.PHONY: push
push: generate ## Push protobuf definitions to buf.build registry
	$(BUF_GEN) push --tag $(VERSION)

##@ Build Dependencies

## Tool Versions
BUF_GEN_VERSION ?= v1.54.0
CI_LINT_VERSION ?= v2.3.0

## Proto tool versions
PROTOC_GEN_GO_VERSION ?= v1.33.0

.PHONY: deps
deps: ## Download all dependencies locally if necessary. if not installed, installation will proceed.
	@echo "Installing tools..."
	@test -s $(GOBIN)/protoc-gen-go || GOBIN=$(GOBIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION)
	@test -s $(GOBIN)/buf || GOBIN=$(GOBIN) go install github.com/bufbuild/buf/cmd/buf@$(BUF_GEN_VERSION)
	@test -s $(GOBIN)/golangci-lint || GOBIN=$(GOBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(CI_LINT_VERSION)
	@echo "All tools installed successfully."
