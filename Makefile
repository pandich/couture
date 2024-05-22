# go setup
GO_MODULE 		= $(shell go list -m)
GOPATH			?= $(shell go env GOPATH)
GOHOSTOS		?= $(shell go env GOHOSTOS)
GOHOSTARCH		?= $(shell go env GOHOSTARCH)
APPLICATION 	= $(notdir $(GO_MODULE))
VERSION			= $(shell cat VERSION)

# how long asciicinema recordings should last
CAST_DURATION   ?= 5s

# TODO add a help target which shows all available targets

#
# Default
.PHONY: all
all: clean build

#
# External Commands
.PHONY: golangci-lint goreleaser gocmt scc gocomplete d2 asciinema asciicast2gif
staticcheck: #
	command -v $@ > /dev/null || go install honnef.co/go/tools/cmd/staticcheck@latest
goreleaser:
	command -v $@ > /dev/null || go install github.com/goreleaser/goreleaser
gocmt:
	command -v $@ > /dev/null || go install github.com/cuonglm/gocmt
scc:
	command -v $@ > /dev/null || go install github.com/boyter/scc
gocomplete:
	command -v $@ > /dev/null || go install github.com/posener/complete/v2/gocomplete
d2:
	command -v $@ > /dev/null || go install oss.terrastruct.com/d2@latest
asciinema:
	command -v $@ > /dev/null || pipx install asciinema
asciicast2gif:
	command -v $@ > /dev/null || npm install --global asciicast2gif


#
# Build
.PHONY: clean build release
clean:
	rm -rf dist/ build/
build: neat
	go build -o build/$(APPLICATION) .
release: goreleaser neat
	goreleaser build --snapshot --rm-dist


#
# Release
.PHONY: tag install uninstall
tag:
	git commit -am "tagging $(VERSION)"
	git tag $(VERSION)
	git push origin $(VERSION)
install: build
	go install .
	echo y | COMP_UNINSTALL=1 go run . > /dev/null || true
	echo y | COMP_INSTALL=1 go run . > /dev/null
uninstall:
	echo y | COMP_UNINSTALL=1 go run . > /dev/null || true
	go clean -i .


#
# Quality
# TODO modernize the code analysis
.PHONY: neat lint
neat:
	go mod tidy
	gofmt -l -s -w .
lint: staticcheck neat
	staticcheck ./...


#
# Documentation
.PHONY: docs diagrams casts
docs: diagrams casts
casts: docs/casts/*.cast.sh
diagrams: docs/diagrams/*.d2

%.cast.sh: build asciinema asciicast2gif
	@echo casting $@
	@timeout --foreground $(CAST_DURATION) asciinema rec --overwrite --command="sh $@" $*.cast || true
	@asciicast2gif -t monokai $*.cast $*.cast.gif

%.d2: d2
	@echo $@ '->' $*.png
	@d2 --theme=200 --font-regular=docs/resources/SimplyMono-Bold.ttf --font-bold=docs/resources/SimplyMono-Bold.ttf --layout=elk $@ $*.png
