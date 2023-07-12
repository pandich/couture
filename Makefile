APPLICATION 	= $(shell go list -m | sed 's/^.*\///')

GO_MODULE 		= $(shell go list -m)
GOPATH			?= $(shell go env GOPATH)
GOHOSTOS		?= $(shell go env GOHOSTOS)
GOHOSTARCH		?= $(shell go env GOHOSTARCH)

.PHONY: all
all: clean build

#
# External Commands

.PHONY: golangci-lint goreleaser gocmt scc gocomplete
staticcheck:
	@command -v $@ > /dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
goreleaser:
	@command -v $@ > /dev/null || go install github.com/goreleaser/goreleaser
gocmt:
	@command -v $@ > /dev/null || go install github.com/cuonglm/gocmt
scc:
	@command -v $@ > /dev/null || go install github.com/boyter/scc
gocomplete:
	@command -v $@ > /dev/null || go install github.com/posener/complete/v2/gocomplete

#
# Targets

# Build
.PHONY: clean build pre-build
build: pre-build
	@go build -o build/$(APPLICATION) .
clean:
	@rm -rf dist/
pre-build: neat

# Release
.PHONY: install uninstall release
install:
	@go install .
uninstall:
	@go clean -i .
release: goreleaser pre-build
	@goreleaser build --snapshot --rm-dist

# Code Quality
.PHONY: neat lint metrics
neat:
	@echo tidying
	@go mod tidy
	@echo formatting
	@gofmt -l -s -w .
lint: staticcheck neat
	@echo vetting
	@go vet
	@echo linting
	@staticcheck ./...
metrics: scc
	@scc --wide --by-file --sort code --include-ext go

# Utility
.PHONY: setup-env install-shell-completions
setup-env: staticcheck goreleaser scc gocmt gocomplete
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
	@asciinema rec --overwrite --command="$(GO_MODULE) --rate-limit=5 --highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa --expand --multiline @@fake" docs/$@.cast
	@make docs/$@.gif
.PHONY: example-fake-single-line
example-fake-single-line:
	@asciinema rec --overwrite --command="$(GO_MODULE) --rate-limit=5 --highlight --filter=+distincto --filter=+'\"first_name\"\s*:\s*\"B' --filter=+quinoa @@fake" docs/$@.cast
	@make docs/$@.gif
.PHONY: %.gif
%.gif:
	@asciicast2gif -t monokai $*.cast $@
