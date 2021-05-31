APPLICATION = couture

#
# Variables

# Core
COMMAND		= main.go
SOURCES	= $(COMMAND) cmd internal

# Go
GO 				= go
GOPATH			?= $(shell $(GO) env GOPATH)
GOHOSTOS		?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH		?= $(shell $(GO) env GOHOSTARCH)
GO_GET 			= $(GO) get -u

.PHONY: default
default: all

#
# External Commands

.PHONY: golangci-lint goreleaser gocmt scc statik
golangci-lint:
	@command -v $@ > /dev/null || $(GO_GET) github.com/golangci/golangci-lint/cmd/golangci-lint
goreleaser:
	@command -v $@ > /dev/null || $(GO_GET) github.com/goreleaser/goreleaser
gocmt:
	@command -v $@ > /dev/null || $(GO_GET) github.com/cuonglm/gocmt
scc:
	@command -v $@ > /dev/null || $(GO_GET) github.com/boyter/scc
statik:
	@command -v $@ > /dev/null || $(GO_GET) github.com/rakyll/statik
#
# Targets

# Build
.PHONY: all clean build
all: clean build
clean:
	@echo cleaning
	@rm -rf dist/
build: neat assets
	@echo building
	@$(GO) build -o dist/couture $(COMMAND)
assets: statik
	@statik -dest=internal/pkg -src=assets -f -p assets

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
neat:
	@echo tidying
	@go mod tidy
	@echo commenting
	@find $(SOURCES) -type d -exec gocmt -p -i -d {} \; 2> /dev/null
	@echo formatting
	@gofmt -l -s -w $(SOURCES)
lint: golangci-lint neat
	@echo linting
	@golangci-lint run
metrics: scc
	@scc --wide --by-file --no-gen --sort lines $(SOURCES)

# Utility
setup-env: golangci-lint goreleaser scc gocmt statik
	@git config --local core.hooksPath .githooks
