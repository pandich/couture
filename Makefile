APPLICATION = couture
SOURCE_DIRS = cmd pkg internal
GOPATH		?= $(shell $(GO) env GOPATH)
GOHOSTOS	?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH	?= $(shell $(GO) env GOHOSTARCH)
INSTALL_DIR	?= $(HOME)/bin

BINARY		= dist/$(APPLICATION)_$(GOHOSTOS)_$(GOHOSTARCH)/$(APPLICATION)

GO 			= go
GO_GET 		= $(GO) get -u
FIND_PEGS 	= find $(SOURCE_DIRS) -name '*.peg' -type f

.PHONY: clean neat generate build lint install

all: clean build

install: build
	@mkdir -p $(INSTALL_DIR) && cp $(BINARY) $(INSTALL_DIR)/

clean:
	@echo cleaning
	@rm -rf dist/
	@$(FIND_PEGS) -exec rm -f {}.go \;

build: generate neat
	@echo building
	@command -v goreleaser > /dev/null || $(GO_GET) github.com/goreleaser/goreleaser
	@goreleaser build --snapshot --rm-dist

generate:
	@echo generating
	@command -v pigeon > /dev/null || $(GO_GET) github.com/mna/pigeon
	@$(FIND_PEGS) -exec pigeon -o {}.go {} \;

neat:
	@echo tidying
	@go mod tidy

	@echo commenting
	@command -v gocmt > /dev/null || $(GO_GET) github.com/cuonglm/gocmt
	@find $(SOURCE_DIRS) -type d -exec gocmt -p -i -d {} \; 2> /dev/null

	@echo formatting
	@gofmt -l -s -w $(SOURCE_DIRS)

lint:
	@echo linting
	@command -v golangci-lint > /dev/null || $(GO_GET) github.com/golangci/golangci-lint/cmd/golangci-lint@v1.39.0
	@golangci-lint run

