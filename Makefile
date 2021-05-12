.PHONY: clean neat build lint install

#
# Core
#

APPLICATION 			=  couture
SOURCE_DIRS 			=  cmd pkg internal
INSTALL_DIR				?= $(HOME)/bin
NATIVE_BINARY			= dist/$(APPLICATION)_$(GOHOSTOS)_$(GOHOSTARCH)/$(APPLICATION)


#
# Go Environment
#

GOPATH					?= $(shell $(GO) env GOPATH)
GOHOSTOS				?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH				?= $(shell $(GO) env GOHOSTARCH)
GO 						= go
GO_GET 					= $(GO) get -u
GORELEASER_ARGS 		?= --snapshot --rm-dist


all: clean build

clean:
	@echo cleaning
	@rm -rf dist/

#
# Build
#

build: neat
	@echo goreleaser building
	@goreleaser build $(GORELEASER_ARGS) --single-target
build-all: goreleaser neat
	@echo building
	@goreleaser build $(GORELEASER_ARGS)

#
# Release
#

install: build
	# TODO do this the proper go way, and not a direct copy
	@mkdir -p $(INSTALL_DIR) && cp $(NATIVE_BINARY) $(INSTALL_DIR)/

release: goreleaser
	@echo releasing
	# TODO configure this properly
	@goreleaser release $(GORELEASER_ARGS)


#
# Quality
#

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

#
# Commands
#

golangci-lint:
	@command -v golangci-lint > /dev/null || $(GO_GET) github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0
goreleaser:
	@command -v goreleaser > /dev/null || $(GO_GET) github.com/goreleaser/goreleaser
gocmt:
	@command -v gocmt > /dev/null || $(GO_GET) github.com/cuonglm/gocmt
