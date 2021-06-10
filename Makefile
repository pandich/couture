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

.PHONY: golangci-lint goreleaser gocmt scc statik gocomplete
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
gocomplete:
	@command -v $@ > /dev/null || $(GO_GET) github.com/posener/complete/v2/gocomplete
#
# Targets

# Build
.PHONY: all clean build assets
all: clean build
clean:
	@echo cleaning
	@rm -rf dist/
	@find $(SOURCES) -name statik.go -exec rm {} \;
build: neat assets
	@echo building
	@$(GO) build -o dist/$(APPLICATION) $(COMMAND)
assets: statik
	@echo assets
	@statik -dest=internal/pkg -src=assets -f -p assets

# Release
.PHONY: install uninstall release
install: build
	@echo installing
	@$(GO) install $(APPLICATION)
uninstall:
	@echo uninstalling
	@$(GO) clean -i $(APPLICATION)
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
.PHONY: setup-env install-shell-completions
setup-env: golangci-lint goreleaser scc gocmt statik gocomplete
	@git config --local core.hooksPath .githooks
install-shell-completions: gocomplete
	@echo installing completions
	@echo y | COMP_UNINSTALL=1 $(APPLICATION) > /dev/null
	@echo y | COMP_INSTALL=1 $(APPLICATION) > /dev/null

# Documentation
.PHONY: record-examples
record-examples: example-fake-multi-line example-fake-single-line
.PHONY: example-fake-multi-line
example-fake-multi-line:
	@asciinema rec --overwrite --command="couture --rate-limit=5 --highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa --expand --multiline @@fake" docs/$@.cast
	@make docs/$@.gif
.PHONY: example-fake-single-line
example-fake-single-line:
	@asciinema rec --overwrite --command="couture --rate-limit=5 --highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa @@fake" docs/$@.cast
	@make docs/$@.gif
.PHONY: %.gif
%.gif:
	@asciicast2gif -t monokai $*.cast $@
