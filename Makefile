.PHONY: default
default: all

#
# Variables

# Core
APPLICATION	= couture
SOURCE_DIRS	= cmd internal
COMMAND		= cmd/$(APPLICATION).go

# Go
GO 				= go
GOPATH			?= $(shell $(GO) env GOPATH)
GOHOSTOS		?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH		?= $(shell $(GO) env GOHOSTARCH)
GO_GET 			= $(GO) get -u

#
# External Commands

.PHONY: golangci-lint goreleaser gocmt scc
golangci-lint:;
	@command -v golangci-lint > /dev/null || $(GO_GET) github.com/golangci/golangci-lint/cmd/golangci-lint
goreleaser:
	@command -v goreleaser > /dev/null || $(GO_GET) github.com/goreleaser/goreleaser
gocmt:
	@command -v gocmt > /dev/null || $(GO_GET) github.com/cuonglm/gocmt
scc:
	@command -v scc > /dev/null || $(GO_GET) github.com/boyter/scc

#
# Targets

# Build
.PHONY: all clean build
all: clean build
clean:
	@echo cleaning
	@rm -rf dist/
build: neat
	@echo building
	@$(GO) build -o dist/couture $(COMMAND)

# Release
.PHONY: install uninstall release
install: build
	@echo installing
	@$(GO) install $(COMMAND)
uninstall:
	@echo uninstalling
	@$(GO) clean -i $(COMMAND)
release: goreleaser build
	@echo releasing
	@goreleaser build --snapshot --rm-dist

# Code Quality
.PHONY: neat lint metrics
neat: gocmt
	@echo tidying
	@go mod tidy
	@echo commenting
	@find $(SOURCE_DIRS) -type d -exec gocmt -p -i -d {} \; 2> /dev/null
	@echo formatting
	@gofmt -l -s -w $(SOURCE_DIRS)
lint: golangci-lint neat
	@echo linting
	@golangci-lint run
metrics: scc
	@scc --wide --by-file --no-gen --sort lines $(SOURCE_DIRS)
